package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/metal-toolbox/mctl/pkg/model"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/skratchdot/open-golang/open"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// The oauth, pkce handling code here was adapted for mctl from an internal project.

const (
	keyringService = "sh.hollow.mctl"
	// tokenNamePrefix = "serverservice"
	// TODO: can this be generalized to localhost - the Oauth provider needs to allow the changed callback URL.
	// pkceCallbackURL = ""
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
func AccessToken(ctx context.Context, apiKind model.APIKind, cfg *model.ConfigOIDC) (string, error) {
	authenticator := newOIDCAuthenticator(apiKind, cfg)

	token, err := authenticator.refreshToken(ctx)
	if err != nil {
		token, err = authenticator.getOAuth2Token(ctx)
		if err != nil {
			return "", err
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

	token, err := a.authCodePKCE(oauthConfig, a.audienceEndpoint)
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

	// scopes := []string{"offline_access", "read"}

	// return oauth configuration
	return &oauth2.Config{
		ClientID:    a.clientID,
		RedirectURL: a.pkceCallbackURL,
		Endpoint:    provider.Endpoint(),
		Scopes:      a.scopes,
	}, nil
}

func (a *authenticator) refreshToken(ctx context.Context) (*oauth2.Token, error) {
	oauthConfig, err := a.oauth2Config(ctx)
	if err != nil {
		return nil, err
	}

	authToken, err := keyring.Get(keyringService, fmt.Sprintf("%s_token", a.tokenNamePrefix))
	if err != nil {
		return nil, err
	}

	refToken, err := keyring.Get(keyringService, fmt.Sprintf("%s_refresh_token", a.tokenNamePrefix))
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
	tok, err := jwt.ParseSigned(rawAccess)
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
	err := keyring.Set(keyringService, fmt.Sprintf("%s_token", a.tokenNamePrefix), token.AccessToken)
	if err != nil {
		return err
	}

	return keyring.Set(keyringService, fmt.Sprintf("%s_refresh_token", a.tokenNamePrefix), token.RefreshToken)
}

// authCodePKCE starts a server and listens for an oauth2 callback and will
// return the API token to the caller
func (a *authenticator) authCodePKCE(oauthConfig *oauth2.Config, audience string) (*oauth2.Token, error) {
	tc := make(chan *oauth2.Token)

	// nolint:gomnd // state string is limited to 20 random characters
	c := &authClient{
		oauthConfig: oauthConfig,
		state:       randStr(20),
	}

	c.codeVerifier, _ = cv.CreateCodeVerifier()

	// nolint:gomnd // read header timeout is set to 30s
	server := &http.Server{Addr: ":18000", ReadHeaderTimeout: time.Second * 30}

	http.HandleFunc("/identity/callback", func(w http.ResponseWriter, r *http.Request) {
		c.handlePKCECallback(w, r, tc)
	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			fmt.Printf("ERROR: %s\n", err.Error())
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
		fmt.Printf("Failed to open browser automatically, please visit %s to complete auth\n\n", authURL)
	}

	token := <-tc

	ctx, cancel := context.WithTimeout(context.Background(), callbackTimeout)
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

func (c *authClient) handlePKCECallback(w http.ResponseWriter, r *http.Request, tc chan *oauth2.Token) {
	state := r.URL.Query().Get("state")
	if state != c.state {
		fmt.Println("ERROR: oauth state doesn't match")
		w.WriteHeader(http.StatusBadRequest)
		tc <- nil
	}

	code := r.URL.Query().Get("code")

	token, err := c.oauthConfig.Exchange(context.Background(), code,
		oauth2.SetAuthURLParam("code_verifier", c.codeVerifier.String()),
	)

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		tc <- nil
	}

	w.Write([]byte("Success. You can now close this window.")) //nolint
	tc <- token
}
