package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	bmclibcomm "github.com/bmc-toolbox/common"
	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	rctypes "github.com/metal-toolbox/rivets/condition"
	rt "github.com/metal-toolbox/rivets/types"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
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
	vendor = strings.TrimSpace(vendor)
	if vendor == "" {
		return []serverservice.ComponentFirmwareSet{}, errors.Wrap(
			ErrFwSetByVendorModel,
			fmt.Sprintf("got empty vendor attribute"),
		)
	}

	model = strings.TrimSpace(model)
	if model == "" {
		return []serverservice.ComponentFirmwareSet{}, errors.Wrap(
			ErrFwSetByVendorModel,
			fmt.Sprintf("got empty model attribute"),
		)
	}

	// ?attr=sh.hollow.firmware_set.labels~vendor~eq~dell&attr=sh.hollow.firmware_set.labels~model~eq~r750&attr=sh.hollow.firmware_set.labels~latest~eq~false
	// list latest, default firmware sets by vendor, model attributes
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
				Keys:      []string{"latest"}, // latest indicates the most current revision of the firmware set.
				Operator:  "eq",
				Value:     "true",
			},
			{
				Namespace: FirmwareSetAttributeNS,
				Keys:      []string{"default"}, // default indicates the firmware set does not belong to an org/project.
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

type ErrUnexpectedResponse struct {
	statusCode int
	message    string
}

func (e *ErrUnexpectedResponse) Error() string {
	s := fmt.Sprintf("status code: %d", e.statusCode)

	if e.message != "" {
		s += " response message: " + e.message
	}

	return fmt.Sprintf("Unexpected response from Conditions API " + s)
}

func newErrUnexpectedResponse(statusCode int, message string) error {
	return &ErrUnexpectedResponse{statusCode, message}
}

// ConditionFromResponse returns a Condition object from the Condition API ServerResponse object
func ConditionFromResponse(response *coapiv1.ServerResponse) (rctypes.Condition, error) {

	if response.StatusCode != http.StatusOK {
		return rctypes.Condition{}, newErrUnexpectedResponse(response.StatusCode, response.Message)
	}

	if response.Records == nil || len(response.Records.Conditions) == 0 {
		return rctypes.Condition{}, errors.New("no record found for Condition")
	}

	return *response.Records.Conditions[0], nil
}

// conditionDisplay is the format in which the condition is printed to the user.
type conditionDisplay struct {
	ID         uuid.UUID       `json:"id"`
	Kind       rctypes.Kind    `json:"kind"`
	State      rctypes.State   `json:"state"`
	Parameters json.RawMessage `json:"parameters"`
	Status     json.RawMessage `json:"status"`
	UpdatedAt  time.Time       `json:"updated_at"`
	CreatedAt  time.Time       `json:"created_at"`
}

// FormatConditionResponse returns a prettyish JSON formatted output that can be printed to stdout.
func FormatConditionResponse(response *coapiv1.ServerResponse, kind rctypes.Kind) (string, error) {
	if response.StatusCode != http.StatusOK {
		return "", newErrUnexpectedResponse(response.StatusCode, response.Message)
	}

	if response.Records == nil {
		return "", errors.New("no records returned")
	}

	if len(response.Records.Conditions) == 0 {
		return "", errors.New("no record found for Condition")
	}

	var inc *rctypes.Condition
	for _, c := range response.Records.Conditions {
		if c.Kind == kind {
			inc = c
		}
	}

	if inc == nil {
		return "", fmt.Errorf("response contains no condition of type %s", string(kind))
	}

	display := &conditionDisplay{
		ID:         inc.ID,
		Kind:       inc.Kind,
		Parameters: inc.Parameters,
		State:      inc.State,
		Status:     inc.Status,
		UpdatedAt:  inc.UpdatedAt,
		CreatedAt:  inc.CreatedAt,
	}

	// XXX: seems highly unlikely that we get a response that deserializes cleanly and doesn't
	// re-serialize.
	b, err := json.MarshalIndent(display, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "bad json in response")
	}

	return string(b), nil
}

func PrintResults(format string, data ...any) {
	switch format {
	case "text":
		spew.Dump(data)
	case "json", "JSON":
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}

// Query server BMC credentials and update the given server object
func ServerBMCCredentials(ctx context.Context, client *serverservice.Client, server *rt.Server) error {
	cred, _, err := client.GetCredential(ctx, uuid.MustParse(server.ID), serverservice.ServerCredentialTypeBMC)
	if err != nil {
		// nolint:goerr113 // error is readable when formatted
		return fmt.Errorf("error in credential lookup for: %s, err: %s", server.ID, err.Error())
	}

	server.BMCUser = cred.Username
	server.BMCPassword = cred.Password

	return nil
}
