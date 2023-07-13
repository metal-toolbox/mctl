package attributes

import (
	"encoding/json"

	bmc "github.com/bmc-toolbox/common"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

var (
	alloyNamespace = "sh.hollow.alloy.outofband.status"
)

type firmwareAttribute struct {
	Firmware *bmc.Firmware `json:"firmware,omitempty"`
}

type ComponentWithFirmware struct {
	Name     string        `json:"name"`
	Vendor   string        `json:"vendor"`
	Model    string        `json:"model"`
	Firmware *bmc.Firmware `json:"firmware"`
}

func FirmwareFromComponents(cs []ss.ServerComponent) []*ComponentWithFirmware {
	var set []*ComponentWithFirmware
	for idx := range cs {
		c := cs[idx]
		for _, attr := range c.VersionedAttributes {
			attr := attr
			if ok, fw := isFirmwareAttribute(&attr); ok {
				set = append(set, &ComponentWithFirmware{
					Name:     c.Name,
					Vendor:   c.Vendor,
					Model:    c.Model,
					Firmware: fw,
				})
			}
		}
	}
	return set
}

func isFirmwareAttribute(attr *ss.VersionedAttributes) (bool, *bmc.Firmware) {
	if attr.Namespace != alloyNamespace {
		return false, nil
	}

	var fw firmwareAttribute
	err := json.Unmarshal(attr.Data, &fw)

	if err != nil || fw.Firmware == nil {
		return false, nil
	}

	return true, fw.Firmware
}
