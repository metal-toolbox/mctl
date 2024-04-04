package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/spf13/cobra"
)

type flagDetails struct {
	name  string
	short string
}

func (f *flagDetails) Name() string {
	return f.name
}

var (
	ConfigFileFlag                    = &flagDetails{name: "config", short: "c"}
	ReAuthFlag                        = &flagDetails{name: "reauth"}
	ServerFlag                        = &flagDetails{name: "server", short: "s"}
	SkipFWStatusFlag                  = &flagDetails{name: "skip-fw-status"}
	SkipBiosConfigFlag                = &flagDetails{name: "skip-bios-config"}
	FromFileFlag                      = &flagDetails{name: "from-file", short: "F"}
	FirmwareIDsFlag                   = &flagDetails{name: "firmware-ids", short: "U"}
	NameFlag                          = &flagDetails{name: "name", short: "n"}
	LabelsFlag                        = &flagDetails{name: "labels", short: "l"}
	FacilityFlag                      = &flagDetails{name: "facility"}
	BMCAddressFlag                    = &flagDetails{name: "bmc-addr", short: "a"}
	BMCUsernameFlag                   = &flagDetails{name: "bmc-user", short: "u"}
	BMCPasswordFlag                   = &flagDetails{name: "bmc-pass", short: "p"}
	FirmwareIDFlag                    = &flagDetails{name: "firmware-id", short: "f"}
	FirmwareVersionFlag               = &flagDetails{name: "firmware-version", short: "V"}
	FirmwareSetFlag                   = &flagDetails{name: "set-id"}
	FirmwareAddFlag                   = &flagDetails{name: "add-firmware-ids"}
	FirmwareRemoveFlag                = &flagDetails{name: "remove-firmware-ids"}
	MacAOCFlag                        = &flagDetails{name: "aoc-mac"}
	MacBMCFlag                        = &flagDetails{name: "bmc-mac"}
	OutputFlag                        = &flagDetails{name: "output", short: "o"}
	ForceFlag                         = &flagDetails{name: "force"}
	DryRunFlag                        = &flagDetails{name: "dry-run"}
	SkipBmcResetFlag                  = &flagDetails{name: "skip-bmc-reset"}
	PowerOffRequiredFlag              = &flagDetails{name: "power-off-required"}
	ModelFlag                         = &flagDetails{name: "model", short: "m"}
	VendorFlag                        = &flagDetails{name: "vendor", short: "v"}
	SlugFlag                          = &flagDetails{name: "slug"}
	WithRecordsFlag                   = &flagDetails{name: "with-records"}
	PageFlag                          = &flagDetails{name: "page"}
	LimitFlag                         = &flagDetails{name: "limit"}
	WithBMCErrorsFlag                 = &flagDetails{name: "with-bmc-errors"}
	WithCredsFlag                     = &flagDetails{name: "with-creds"}
	PrintTableFlag                    = &flagDetails{name: "table", short: "t"}
	BiosConfigFlag                    = &flagDetails{name: "bios-config"}
	ListComponentsFlag                = &flagDetails{name: "list-components"}
	ServerSerialFlag                  = &flagDetails{name: "serial"}
	ServerActionPowerActionFlag       = &flagDetails{name: "action"}
	ServerActionPowerActionStatusFlag = &flagDetails{name: "action-status"}

	OutputTypeJSON outputType = "json"
	OutputTypeText outputType = "text"
)

var (
	errOutputType = errors.New("unsupported output type")
)

type outputType string

func (o *outputType) String() string {
	return string(*o)
}

func (o *outputType) Type() string {
	return "outputType"
}

func (o *outputType) Set(value string) error {
	value = strings.ToLower(value)

	switch value {
	case OutputTypeJSON.String(), OutputTypeText.String():
		*o = outputType(value)
		return nil
	default:
		return errOutputType
	}
}

//nolint:staticcheck // SA5011 log.Fatalf will make sure we don't continue if flag is nil
func RequireFlag(cmd *cobra.Command, flagDetail *flagDetails) {
	flag := cmd.PersistentFlags().Lookup(flagDetail.name)
	if flag == nil {
		log.Fatalf("no flag with name '%s'", flagDetail.name)
	}

	if err := cmd.MarkPersistentFlagRequired(flag.Name); err != nil {
		log.Fatal(err)
	}

	var flagUsage string
	if flag.Shorthand != "" {
		flagUsage = "-" + flag.Shorthand
	} else {
		flagUsage = "--" + flag.Name
	}

	valueType := strings.ToUpper(strings.ReplaceAll(flag.Name, "-", ""))

	flag.Usage = "[required] " + flag.Usage
	cmd.Use = fmt.Sprintf("%s %s %s", cmd.Use, flagUsage, valueType)
}

func MutuallyExclusiveFlags(cmd *cobra.Command, flags ...*flagDetails) {
	flagNames := make([]string, len(flags))

	for i, f := range flags {
		flagNames[i] = f.name
	}

	cmd.MarkFlagsMutuallyExclusive(flagNames...)
}

func RequireOneFlag(cmd *cobra.Command, flags ...*flagDetails) {
	flagNames := make([]string, len(flags))

	for i, f := range flags {
		flagNames[i] = f.name
	}

	cmd.MarkFlagsOneRequired(flagNames...)
}

func AddConfigFileFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, ConfigFileFlag.name, ConfigFileFlag.short, "",
		"config file (default is $XDG_CONFIG_HOME/mctl/config.yml)")
}

func AddReAuthFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, ReAuthFlag.name, false, "re-authenticate with oauth services")
}

func AddServerFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, ServerFlag.name, ServerFlag.short, "", "ID of the server")
}

func AddSkipFWStatusFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, SkipFWStatusFlag.name, false, "Skip firmware status data collection")
}

func AddSkipBiosConfigFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, SkipBiosConfigFlag.name, false, "Skip BIOS configuration data collection")
}

func AddFromFileFlag(cmd *cobra.Command, ptr *string, usage string) {
	cmd.PersistentFlags().StringVarP(ptr, FromFileFlag.name, FromFileFlag.short, "", usage)
}

func AddFirmwareIDsFlag(cmd *cobra.Command, ptr *[]string) {
	usage := "comma separated list of firmware IDs"
	cmd.PersistentFlags().StringSliceVarP(ptr, FirmwareIDsFlag.name, FirmwareIDsFlag.short, []string{}, usage)
}

func AddNameFlag(cmd *cobra.Command, ptr *string, usage string) {
	cmd.PersistentFlags().StringVarP(ptr, NameFlag.name, NameFlag.short, "", usage)
}

//nolint:gocritic // ptrToRefParam we need the pointer map argument
func AddLabelsFlag(cmd *cobra.Command, ptr *map[string]string, usage string) {
	cmd.PersistentFlags().StringToStringVarP(ptr, LabelsFlag.name, LabelsFlag.short, nil, usage)
}

func AddFacilityFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVar(ptr, FacilityFlag.name, "", "facility name")
}

func AddBMCAddressFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, BMCAddressFlag.name, BMCAddressFlag.short, "", "address of the bmc")
}

func AddBMCUsernameFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, BMCUsernameFlag.name, BMCUsernameFlag.short, "", "username of the bmc user")
}

func AddBMCPasswordFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, BMCPasswordFlag.name, BMCPasswordFlag.short, "", "password of the bmc user")
}

func AddFirmwareIDFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, FirmwareIDFlag.name, FirmwareIDFlag.short, "", "ID of the firmware")
}

func AddFirmwareVersionFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, FirmwareVersionFlag.name, FirmwareVersionFlag.short, "", "firmware version")
}

func AddFirmwareSetFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVar(ptr, FirmwareSetFlag.name, "", "ID of the firmware set")
}

func AddFirmwareAddIDsFlag(cmd *cobra.Command, ptr *[]string) {
	usage := "comma separated list of firmware IDs to be added"
	cmd.PersistentFlags().StringSliceVar(ptr, FirmwareAddFlag.name, []string{}, usage)
}

func AddFirmwareRemoveIDsFlag(cmd *cobra.Command, ptr *[]string) {
	usage := "comma separated list of firmware IDs to be removed"
	cmd.PersistentFlags().StringSliceVar(ptr, FirmwareRemoveFlag.name, []string{}, usage)
}

func AddMacAOCFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVar(ptr, MacAOCFlag.name, "", "aoc mac address")
}

func AddMacBMCFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVar(ptr, MacBMCFlag.name, "", "bmc mac address")
}

func AddOutputFlag(cmd *cobra.Command, ptr *string) {
	*ptr = OutputTypeJSON.String() // default value
	outputFlag := (*outputType)(ptr)
	cmd.PersistentFlags().VarP(outputFlag, OutputFlag.name, OutputFlag.short, "{json|text}")
}

func AddForceFlag(cmd *cobra.Command, ptr *bool, usage string) {
	cmd.PersistentFlags().BoolVar(ptr, ForceFlag.name, false, usage)
}

func AddDryRunFlag(cmd *cobra.Command, ptr *bool, usage string) {
	cmd.PersistentFlags().BoolVar(ptr, DryRunFlag.name, false, usage)
}

func AddSkipBmcResetFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, SkipBmcResetFlag.name, false, "skip BMC reset before firmware install")
}

func AddPowerOffRequiredFlag(cmd *cobra.Command, ptr *bool, usage string) {
	cmd.PersistentFlags().BoolVar(ptr, PowerOffRequiredFlag.name, false, usage)
}

func AddModelFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, ModelFlag.name, ModelFlag.short, "", "filter by model")
}

func AddVendorFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVarP(ptr, VendorFlag.name, VendorFlag.short, "", "filter by vendor")
}

func AddSlugFlag(cmd *cobra.Command, ptr *string, usage string) {
	cmd.PersistentFlags().StringVar(ptr, SlugFlag.name, "", usage)
}

func AddWithRecordsFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, WithRecordsFlag.name, false,
		"print record count found with pagination info and return")
}

func AddPageFlag(cmd *cobra.Command, ptr *int) {
	cmd.PersistentFlags().IntVar(ptr, PageFlag.name, 0, "limit results to page (for use with --limit)")
}

func AddPageLimitFlag(cmd *cobra.Command, ptr *int) {
	defaultLimit := 10
	usage := fmt.Sprintf("limit results returned. Max value is %d (hard limit set in fleetdb). To list more than %d, you must query each page (with '--page') individually", fleetdbapi.MaxPaginationSize, fleetdbapi.MaxPaginationSize)

	cmd.PersistentFlags().IntVar(ptr, LimitFlag.name, defaultLimit, usage)
}

func AddWithBMCErrorsFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, WithBMCErrorsFlag.name, false, "include BMC errors")
}

func AddWithCredsFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, WithCredsFlag.name, false, "include credentials")
}

func AddPrintTableFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, PrintTableFlag.name, false, "format output in a table")
}

func AddBIOSConfigFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, BiosConfigFlag.name, false, "print bios configuration")
}

func AddListComponentsFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(ptr, ListComponentsFlag.name, false, "include component data")
}

func AddServerSerialFlag(cmd *cobra.Command, ptr *string) {
	cmd.PersistentFlags().StringVar(ptr, ServerSerialFlag.name, "", "filter by server serial")
}

func AddServerPowerActionFlag(cmd *cobra.Command, ptr *string, params []string) {
	cmd.PersistentFlags().StringVar(
		ptr,
		ServerActionPowerActionFlag.name,
		"",
		fmt.Sprintf("run a server power action [%s]", strings.Join(params, "|")),
	)
}

func AddServerPowerActionStatusFlag(cmd *cobra.Command, ptr *bool) {
	cmd.PersistentFlags().BoolVar(
		ptr,
		ServerActionPowerActionStatusFlag.name,
		false,
		"Query the last power action status/response",
	)
}
