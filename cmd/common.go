package cmd

import (
	"encoding/json"
	"log"
	"time"

	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	ErrAttributeFromLabel = errors.New("error creating Attribute from Label")
	ErrLabelFromAttribute = errors.New("error creating Label from Attribute")
)

const (
	CmdTimeout = 20 * time.Second
)

func MustCreateApp(ctx context.Context) *app.App {
	mctl, err := app.New(ctx, cfgFile, reAuth)
	if err != nil {
		log.Fatal(err)
	}

	return mctl
}

func AttributeFromLabels(ns string, labels map[string]string) (*serverservice.Attributes, error) {
	data, err := json.Marshal(labels)
	if err != nil {
		return nil, errors.Wrap(ErrAttributeFromLabel, err.Error())
	}

	return &serverservice.Attributes{
		Namespace: ns,
		Data:      data,
	}, nil
}
