//go:build staff
// +build staff

package get

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/bmc-toolbox/common"
	"github.com/go-playground/validator"
	"github.com/metal-toolbox/rivets/v2/types"
	"go.equinixmetal.net/staff"
)

type commonComponents struct {
	Chassis chassis `json:"chassis" validate:"required"`
}

type chassis struct {
	BMC         BMC         `json:"bmc" validate:"required"`
	Motherboard Motherboard `json:"motherboard" validate:"required"`
	CPU         []CPU       `json:"cpu" validate:"required"`
	Memory      []Memory    `json:"memory" validate:"required"`
	Disk        []Disk      `json:"disk" validate:"required"`
	NIC         []NIC       `json:"nic" validate:"required"`
}

type BMC struct {
	Manufacturer string   `json:"manufacturer" validate:"required"`
	Model        string   `json:"model" validate:"required"`
	Firmware     Firmware `json:"firmware" validate:"required"`
}

type Motherboard struct {
	Manufacturer string   `json:"manufacturer" validate:"required"`
	Model        string   `json:"model" validate:"required"`
	Firmware     Firmware `json:"firmware" validate:"required"`
	Serial       string   `json:"serial"`
}

type CPU struct {
	Manufacturer string   `json:"manufacturer" validate:"required"`
	Model        string   `json:"model" validate:"required"`
	Cores        int      `json:"cores" validate:"required"`
	ClockSpeed   int64    `json:"clock_speed" validate:"required"`
	Firmware     Firmware `json:"firmware" validate:"required"`
}

type Memory struct {
	Manufacturer string   `json:"manufacturer" validate:"required"`
	Model        string   `json:"model" validate:"required"`
	Capacity     int64    `json:"capacity" validate:"required"`
	Speed        int64    `json:"speed" validate:"required"`
	Firmware     Firmware `json:"firmware" validate:"required"`
	Serial       string   `json:"serial"`
}

type Disk struct {
	Manufacturer string   `json:"manufacturer" validate:"required"`
	Model        string   `json:"model" validate:"required"`
	Capacity     int64    `json:"capacity" validate:"required"`
	Firmware     Firmware `json:"firmware" validate:"required"`
	Serial       string   `json:"serial"`
}

type NIC struct {
	Manufacturer string   `json:"manufacturer" validate:"required"`
	Model        string   `json:"model" validate:"required"`
	Speed        int64    `json:"speed" validate:"required"`
	Firmware     Firmware `json:"firmware" validate:"required"`
	MACs         []string `json:"macs" validate:"required"`
}

type Firmware struct {
	Version string `json:"version" validate:"required"`
}

const (
	BMCType         = "ManagementControllerComponent"
	MotherboardType = "MotherboardComponent"
	CPUType         = "ProcessorComponent"
	MemoryType      = "MemoryComponent"
	DiskType        = "DiskComponent"
	NICType         = "NetworkComponent"

	NoFirmware = "N/A"
)

func convertFleetDBServer(server *types.Server) (*commonComponents, error) {
	commonData := &commonComponents{}

	for _, component := range server.Components {
		firmware := Firmware{}
		if component.Firmware != nil {
			firmware = Firmware{component.Firmware.Installed}
		}
		attribute := component.Attributes
		if attribute == nil {
			attribute = &types.ComponentAttributes{}
		}

		switch component.Name {
		case common.SlugPhysicalMem:
			commonData.Chassis.Memory = append(commonData.Chassis.Memory, Memory{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        attribute.PartNumber,
				Speed:        attribute.ClockSpeedHz,
				Capacity:     attribute.SizeBytes,
				Firmware:     firmware,
				Serial:       component.Serial,
			})

		case common.SlugCPU:
			commonData.Chassis.CPU = append(commonData.Chassis.CPU, CPU{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Cores:        attribute.Cores,
				ClockSpeed:   attribute.ClockSpeedHz,
				Firmware:     firmware,
			})

		case common.SlugNIC:
			commonData.Chassis.NIC = append(commonData.Chassis.NIC, NIC{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				// TO BE FILLED. It's is inside the `data`(json) in attribute where GetInventory hasn't parsed
				// or copied to the response.
				Speed:    -1,
				Firmware: firmware,
				// TO BE FILLED. It's is inside the `data`(json) in attribute where GetInventory hasn't parsed
				// or copied to the response.
				MACs: []string{},
			})

		case common.SlugBMC:
			commonData.Chassis.BMC = BMC{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Firmware:     firmware,
			}

		case common.SlugDrive:
			var bytes int64 = attribute.CapacityBytes

			commonData.Chassis.Disk = append(commonData.Chassis.Disk, Disk{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				// Sometimes GetInventory will drop this stat while it is in the database.
				Capacity: bytes,
				Firmware: firmware,
				Serial:   component.Serial,
			})

		case common.SlugBIOS:
			commonData.Chassis.Motherboard = Motherboard{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Firmware:     firmware,
				Serial:       component.Serial,
			}
		}
	}

	validate := validator.New()
	err := validate.Struct(commonData)
	if err != nil {
		fmt.Printf("[ERROR] EMAPI is missing required fields: %v\n", err)
	}

	return commonData, nil
}

func convertEMAPIServer(comps []staff.Component) (*commonComponents, error) {
	commonData := &commonComponents{}

	for _, component := range comps {
		firmware := Firmware{component.FirmwareVersion}
		switch component.Type {
		case BMCType:
			commonData.Chassis.BMC = BMC{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Firmware:     firmware,
			}

		case MotherboardType:
			commonData.Chassis.Motherboard = Motherboard{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Firmware:     firmware,
				Serial:       component.Serial,
			}

		case CPUType:
			cores, _ := component.Data["cores"].(float64)
			clock, _ := component.Data["clock"].(float64)
			commonData.Chassis.CPU = append(commonData.Chassis.CPU, CPU{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Cores:        int(cores),
				ClockSpeed:   int64(clock),
				Firmware:     firmware,
			})

		case DiskType:
			bytesStr := component.Data["size"]
			bytes, _ := parseSizeWithUnit(bytesStr.(string))
			commonData.Chassis.Disk = append(commonData.Chassis.Disk, Disk{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Capacity:     bytes,
				Firmware:     firmware,
				Serial:       component.Serial,
			})

		case MemoryType:
			capStr := component.Data["size"].(string)
			cap, _ := parseSizeWithUnit(capStr)
			commonData.Chassis.Memory = append(commonData.Chassis.Memory, Memory{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				// GetServerComponent and GetInventory don't have this stats currently.
				Capacity: cap,
				Speed:    int64(math.Round(component.Data["clock"].(float64) * 1000000)),
				Firmware: firmware,
				Serial:   component.Serial,
			})

		case NICType:
			var rate int64
			if component.Data["rate"] != nil {
				if rateStr, ok := component.Data["rate"].(string); ok {
					rate, _ = strconv.ParseInt(rateStr, 10, 64)
				}
			}

			commonData.Chassis.NIC = append(commonData.Chassis.NIC, NIC{
				Manufacturer: strings.ToLower(component.Vendor),
				Model:        component.Model,
				Speed:        rate,
				Firmware:     firmware,
				MACs:         []string{component.Serial},
			})
		}
	}

	validate := validator.New()
	err := validate.Struct(commonData)
	if err != nil {
		fmt.Printf("[ERROR] EMAPI is missing required fields: %v\n", err)
	}
	return commonData, nil
}

func parseSizeWithUnit(size string) (int64, error) {
	if size == "" {
		return 0, nil
	}

	unitMap := map[string]uint64{
		"":   1,
		"B":  1,
		"K":  1024,
		"KB": 1024,
		"M":  1024 * 1024,
		"MB": 1024 * 1024,
		"G":  1024 * 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"T":  1024 * 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	size = strings.ToUpper(size)
	re := regexp.MustCompile("[A-Za-z]+")
	unit := re.FindString(size)
	valueStr := size[:len(size)-2]

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}

	if multiplier, ok := unitMap[unit]; ok {
		bytes := int64(value * float64(multiplier))
		return bytes, nil
	}

	return 0, fmt.Errorf("invalid unit: %s", unit)
}
