package install

import "github.com/google/uuid"

// This is a copy of the TaskParamters struct from flasher
//
// https://github.com/metal-toolbox/flasher/blob/4e663e45288bb63d0826d8f84a93ab8f58ea82ff/internal/model/task.go#L150C1-L177C2
type parameters struct {
	// Inventory identifier for the asset to install firmware on.
	AssetID uuid.UUID `json:"assetID"`

	// Reset device BMC before firmware install
	ResetBMCBeforeInstall bool `json:"resetBMCBeforeInstall,omitempty"`

	// Force install given firmware regardless of current firmware version.
	ForceInstall bool `json:"forceInstall,omitempty"`

	// FirmwareSetID specifies the firmware set to be applied.
	FirmwareSetID uuid.UUID `json:"firmwareSetID,omitempty"`
}
