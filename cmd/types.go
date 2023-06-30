package cmd

// firmware set command flags
type FirmwareSetFlags struct {
	// labels are key values
	Labels map[string]string
	// id is the firmware set id
	ID string
	// comma separated list of firmware UUIDs
	FirmwareUUIDs string
	// name for the firmware set to be created/edited
	FirmwareSetName string
}
