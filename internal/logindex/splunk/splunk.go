package splunk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	gosplunk "github.com/georgestarcher/querysplunk/goSplunk"
	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/pkg/model"
)

const (
	timeout = 120 * time.Second
	index   = "flasher"
)

var (
	ErrTimeout = errors.New("timeout waiting for search results")
)

type Splunk struct {
	client *gosplunk.SplunkConnection
}

func NewSplunkClient(cfg *model.ConfigLogIndex) (*Splunk, error) {
	// setup splunk connection structure
	client := gosplunk.SplunkConnection{
		Authtoken: cfg.Token,
		BaseURL:   cfg.Endpoint,
		TLSverify: true,
		Timeout:   timeout,
	}

	return &Splunk{client}, nil
}

func (s *Splunk) SearchByAssetID(ctx context.Context, serverID, conditionID uuid.UUID) error {
	err := s.client.Login()
	if err != nil {
		log.Fatalf("ERROR: Couldn't login to splunk: %s", err)
	}

	query := fmt.Sprintf(`index=%s earliest=-1h@h and latest=@m | search %s | fields "_time", "assetFacilityCode", "assetID", "bmc", "msg"`, index, serverID.String())
	gosplunk.SplunkQuery{Query: query}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = conn.DispatchQuery(&splunkQuery, outputfile)
	}()

	wg.Wait()

	if err != nil {
		log.Fatalf("ERROR: %s", err)
	} else {
		log.Print("SUCCESS: Query Completed")
	}

	return nil
}
