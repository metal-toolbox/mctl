package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	bmclibcomm "github.com/bmc-toolbox/common"
	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
	"golang.org/x/net/context"
)

var (
	ErrAttributeFromLabel = errors.New("error creating Attribute from Label")
	ErrLabelFromAttribute = errors.New("error creating Label from Attribute")
	ErrFwSetByVendorModel = errors.New("error identifying firmware set by server vendor, model")
)

const (
	CmdTimeout = 20 * time.Second

	// TODO: merge constants along with the ones in Alloy into a separate library
	ServerVendorAttributeNS = "sh.hollow.alloy.server_vendor_attributes"
	FirmwareSetAttributeNS  = "sh.hollow.firmware_set.labels"
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

// AttributeByNamespace returns the serverservice attribute in the slice that matches the namespace
//
// TODO: move into common library and share with Alloy
func AttributeByNamespace(ns string, attributes []serverservice.Attributes) *serverservice.Attributes {
	for _, attribute := range attributes {
		if attribute.Namespace == ns {
			return &attribute
		}
	}

	return nil
}

// VendorModelFromAttrs unpacks the attributes payload to return the vendor, model attributes for a server
//
// TODO: move into common library and share with Alloy
func VendorModelFromAttrs(attrs []serverservice.Attributes) (vendor, model string) {
	attr := AttributeByNamespace(ServerVendorAttributeNS, attrs)
	if attr == nil {
		return "", ""
	}

	data := map[string]string{}
	if err := json.Unmarshal(attr.Data, &data); err != nil {
		return "", ""
	}

	return bmclibcomm.FormatVendorName(data["vendor"]), bmclibcomm.FormatProductName(data["model"])
}

// FirmwareSetIDByVendorModel returns the firmware set ID matched by the vendor, model attributes
//
// TODO: move into common library
func FirmwareSetIDByVendorModel(ctx context.Context, vendor, model string, client *serverservice.Client) (uuid.UUID, error) {
	fwSet, err := FirmwareSetByVendorModel(ctx, vendor, model, client)
	if err != nil {
		return uuid.Nil, err
	}

	log.Printf(
		"fw sets identified for vendor: %s, model: %s, fwset: %s\n",
		vendor,
		model,
		fwSet[0].UUID.String(),
	)

	return fwSet[0].UUID, nil
}

// FirmwareSetByVendorModel returns the firmware set matched by the vendor, model attributes
//
// TODO: move into common library
func FirmwareSetByVendorModel(ctx context.Context, vendor, model string, client *serverservice.Client) ([]serverservice.ComponentFirmwareSet, error) {
	// ?attr=sh.hollow.firmware_set.labels~vendor~eq~dell&attr=sh.hollow.firmware_set.labels~model~eq~r750&attr=sh.hollow.firmware_set.labels~latest~eq~false
	// list latest firmware sets by vendor, model attributes
	fwSetListparams := &serverservice.ComponentFirmwareSetListParams{
		AttributeListParams: []serverservice.AttributeListParams{
			{
				Namespace: FirmwareSetAttributeNS,
				Keys:      []string{"vendor"},
				Operator:  "eq",
				Value:     strings.ToLower(vendor),
			},
			{
				Namespace: FirmwareSetAttributeNS,
				Keys:      []string{"model"},
				Operator:  "like",
				Value:     strings.ToLower(model),
			},
			{
				Namespace: FirmwareSetAttributeNS,
				Keys:      []string{"latest"},
				Operator:  "eq",
				Value:     "true",
			},
		},
	}

	fwSet, _, err := client.ListServerComponentFirmwareSet(ctx, fwSetListparams)
	if err != nil {
		return []serverservice.ComponentFirmwareSet{}, errors.Wrap(ErrFwSetByVendorModel, err.Error())
	}

	if len(fwSet) == 0 {
		return []serverservice.ComponentFirmwareSet{}, errors.Wrap(
			ErrFwSetByVendorModel,
			fmt.Sprintf("no fw sets identified for vendor: %s, model: %s", vendor, model),
		)
	}

	return fwSet, nil
}
