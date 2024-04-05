package cmd

// firmware set command flags
type FirmwareSetFlags struct {
	// labels are key values
	Labels map[string]string
	// id is the firmware set id
	ID string
	// list of firmware UUIDs to be added to the set
	AddFirmwareUUIDs []string
	// list of firmware UUIDs to be removed from the set
	RemoveFirmwareUUIDs []string
	// name for the firmware set to be created/edited
	FirmwareSetName string
	// create firmware set from file
	CreateFromFile string
	// ignore any create errors
	IgnoreCreateErrors bool
}
