package logindex

import (
	"context"

	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/logindex/splunk"
	"github.com/metal-toolbox/mctl/pkg/model"
)

type Queryor interface {
	SearchByAssetID(ctx context.Context, assetID, conditionID uuid.UUID) error
}

func NewQueryor(cfg *model.ConfigLogIndex) (Queryor, error) {
	return splunk.NewSplunkClient(cfg)
}
