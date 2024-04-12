package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/metal-toolbox/mctl/pkg/model"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/skratchdot/open-golang/open"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

// The oauth, pkce handling code here was adapted for mctl from an internal project.

const (
	keyringService = "sh.hollow.mctl"
)

var (
	callbackTimeout = time.Second * 6
	// ErrNoToken is returned when a token isn't returned from the auth flow
	ErrNoToken = errors.New("failed to get a token")
)

type authenticator struct {
	disable          bool
	tokenNamePrefix  string
	pkceCallbackURL  string
	clientID         string
	audienceEndpoint string
	issuerEndpoint   string
	scopes           []string
}

func newOIDCAuthenticator(apiKind model.APIKind, cfg *model.ConfigOIDC) *authenticator {
	return &authenticator{
		disable:          cfg.Disable,
		tokenNamePrefix:  string(apiKind),
		pkceCallbackURL:  cfg.PkceCallbackURL,
		clientID:         cfg.ClientID,
		audienceEndpoint: cfg.AudienceEndpoint,
		issuerEndpoint:   cfg.IssuerEndpoint,
		scopes:           cfg.Scopes,
	}
}

// AccessToken looks up the keyring for the service access token, if none is found, it fetches a new one.
func AccessToken(ctx context.Context, apiKind model.APIKind, cfg *model.ConfigOIDC, reauth bool) (string, error) {
	authenticator := newOIDCAuthenticator(apiKind, cfg)

	var token *oauth2.Token
	var err error

	if reauth {
		token, err = authenticator.getOAuth2Token(ctx)
		if err != nil {
			return "", err
		}
	} else {
		token, err = authenticator.refreshToken(ctx)
		if err != nil {
			token, err = authenticator.getOAuth2Token(ctx)
			if err != nil {
				return "", err
			}
		}
	}

	return token.AccessToken, nil
}

// GetOAuth2Token retrieves the OAuth2 token from the issuer and stores it in the local keyring with the given name.
func (a *authenticator) getOAuth2Token(ctx context.Context) (*oauth2.Token, error) {
	oauthConfig, err := a.oauth2Config(ctx)
	if err != nil {
		return nil, err
	}

	token, err := a.authCodePKCE(ctx, oauthConfig, a.audienceEndpoint)
	if err != nil {
		return nil, err
	}

	if err := a.keyringStoreToken(token); err != nil {
		return nil, err
	}

	return token, nil
}

func (a *authenticator) oauth2Config(ctx context.Context) (*oauth2.Config, error) {
	// setup oidc provider
	provider, err := oidc.NewProvider(ctx, a.issuerEndpoint)
	if err != nil {
		return nil, err
	}

	// return oauth configuration
	return &oauth2.Config{
		ClientID:    a.clientID,
		RedirectURL: a.pkceCallbackURL,
		Endpoint:    provider.Endpoint(),
		Scopes:      a.scopes,
	}, nil
}

func (a *authenticator) keyringNameRefreshToken() string {
	return fmt.Sprintf("%s_%s_refresh_token", a.clientID, a.tokenNamePrefix)
}

func (a *authenticator) keyringNameToken() string {
	return fmt.Sprintf("%s_%s_token", a.clientID, a.tokenNamePrefix)
}

func (a *authenticator) refreshToken(ctx context.Context) (*oauth2.Token, error) {
	oauthConfig, err := a.oauth2Config(ctx)
	if err != nil {
		return nil, err
	}

	authToken, err := keyring.Get(keyringService, a.keyringNameToken())
	if err != nil {
		return nil, err
	}

	refToken, err := keyring.Get(keyringService, a.keyringNameRefreshToken())
	if err != nil {
		return nil, err
	}

	token, err := a.tokenFromRaw(authToken, refToken)
	if err != nil {
		return nil, err
	}

	ts := oauthConfig.TokenSource(ctx, token)

	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}

	// if the token was refreshed we need to save the new token
	if newToken.AccessToken != token.AccessToken {
		if err := a.keyringStoreToken(newToken); err != nil {
			return nil, err
		}
	}

	return newToken, nil
}

// tokenFromRaw will take a access and refresh token string and convert them into a proper token
func (a *authenticator) tokenFromRaw(rawAccess, refresh string) (*oauth2.Token, error) {
	tok, err := jwt.ParseSigned(rawAccess, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, err
	}

	cl := jwt.Claims{}

	if err := tok.UnsafeClaimsWithoutVerification(&cl); err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  rawAccess,
		RefreshToken: refresh,
		Expiry:       cl.Expiry.Time(),
	}, nil
}

func (a *authenticator) keyringStoreToken(token *oauth2.Token) error {
	err := keyring.Set(keyringService, a.keyringNameToken(), token.AccessToken)
	if err != nil {
		return err
	}

	return keyring.Set(keyringService, a.keyringNameRefreshToken(), token.RefreshToken)
}

// authCodePKCE starts a server and listens for an oauth2 callback and will
// return the API token to the caller
func (a *authenticator) authCodePKCE(ctx context.Context, oauthConfig *oauth2.Config, audience string) (*oauth2.Token, error) {
	tc := make(chan *oauth2.Token)

	// nolint:gomnd // state string is limited to 20 random characters
	c := &authClient{
		oauthConfig: oauthConfig,
		state:       randStr(20),
	}

	c.codeVerifier, _ = cv.CreateCodeVerifier()

	mux := http.NewServeMux()
	mux.HandleFunc("/identity/callback", func(w http.ResponseWriter, r *http.Request) {
		c.handlePKCECallback(ctx, w, r, tc)
	})

	// nolint:gomnd // read header timeout is set to 30s
	server := &http.Server{Addr: ":18000", ReadHeaderTimeout: time.Second * 30, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			log.Printf("ERROR: %s\n", err.Error())
			tc <- nil
		}
	}()

	// Create code_challenge with S256 method
	codeChallenge := c.codeVerifier.CodeChallengeS256()
	authURL := oauthConfig.AuthCodeURL(c.state,
		oauth2.SetAuthURLParam("audience", audience),
		oauth2.SetAuthURLParam("key", "value"),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	if err := open.Start(authURL); err != nil {
		log.Printf("Failed to open browser automatically, please visit %s to complete auth\n\n", authURL)
	}

	token := <-tc

	ctx, cancel := context.WithTimeout(ctx, callbackTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return nil, err
	}

	if token == nil {
		return nil, ErrNoToken
	}

	return token, nil
}

func randStr(length int) string {
	buff := make([]byte, length)
	_, _ = rand.Read(buff)

	return base64.StdEncoding.EncodeToString(buff)[:length]
}

type authClient struct {
	oauthConfig  *oauth2.Config
	codeVerifier *cv.CodeVerifier
	state        string
}

func (c *authClient) handlePKCECallback(ctx context.Context, w http.ResponseWriter, r *http.Request, tc chan *oauth2.Token) {
	state := r.URL.Query().Get("state")
	if state != c.state {
		log.Printf("ERROR: oauth state doesn't match")
		w.WriteHeader(http.StatusBadRequest)
		tc <- nil
	}

	code := r.URL.Query().Get("code")

	token, err := c.oauthConfig.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", c.codeVerifier.String()),
	)
	if err != nil {
		log.Printf("ERROR in token exchange: %s\n", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		tc <- nil
	}

	w.Write(successHTML()) //nolint
	tc <- token
}

func successHTML() []byte {
	return []byte(`
	<div class="info"><p>Auth success, you may close this page.</p></div>
	<div class="auto"></div>

	<script>
	let auto = document.querySelector('.auto');
	var timeout = 5;
	setTimeout(function() { this.close(); }, timeout*1000);
	var id = setInterval(function() {
		if (timeout < 1) {
			clearInterval(id);
			auto.innerHTML = "Unable to close page automatically";
			return;
		}

		timeout--;
		auto.innerHTML = "This page will close automatically in " + timeout + " seconds...";
	}, 1000);
	</script>
`)
}
