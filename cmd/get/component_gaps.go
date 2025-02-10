//go:build staff
// +build staff

// We use build tag to avoid mctl infrastructure changes as staff is a private repo
// which lint and gendoc cannot access.
// To download the staff module, please run 2 commands below:
// export GOPRIVATE=go.equinixmetal.net/*
// git config --global url.ssh://git@github.com/equinixmetal.insteadOf https://github.com/equinixmetal
// See https://github.com/equinixmetal/go-staff
package get

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.equinixmetal.net/staff"

	mctl "github.com/metal-toolbox/mctl/cmd"
	"github.com/metal-toolbox/mctl/internal/app"
)

type getComponentGapsFlags struct {
	id string
}

var (
	flagsDefinedGetComponentGaps *getComponentGapsFlags
)

// CompareComponents compares two objects and prints differences, even if they're not directly structs
func CompareComponents(obj1, obj2 interface{}, path string) {
	if reflect.DeepEqual(obj1, obj2) {
		return
	}

	v1 := reflect.ValueOf(obj1)
	v2 := reflect.ValueOf(obj2)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}
	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if v1.Kind() != reflect.Struct && v1.Kind() != reflect.Slice {
		if v1.Interface() != v2.Interface() {
			fmt.Printf("Component %s: %v != %v\n", path, v1.Interface(), v2.Interface())
		}
		return
	}

	typeOfObj := v1.Type()
	for i := 0; i < v1.NumField(); i++ {
		fieldName := typeOfObj.Field(i).Name
		value1 := v1.Field(i)
		value2 := v2.Field(i)

		// Build the current path for output (Component + Field)
		currentPath := path + "." + fieldName
		if value1.Kind() == reflect.Slice {
			compareSlices(currentPath, value1, value2, typeOfObj)
		} else if value1.Kind() == reflect.Struct {
			CompareComponents(value1.Interface(), value2.Interface(), currentPath)
		} else {
			if value1.Interface() != value2.Interface() {
				fmt.Printf("Component %s Field %s: %v != %v\n", path, fieldName, value1.Interface(), value2.Interface())
			}
		}
	}
}

func compareSlices(path string, slice1, slice2 reflect.Value, typeOfObj reflect.Type) {
	if slice1.Len() != slice2.Len() {
		fmt.Printf("Component %s Field %s: lengths differ (%d != %d)\n", path, typeOfObj, slice1.Len(), slice2.Len())
	}

	// Slice item order may be different in fleetDB and EMAPI, eg {CPU1, CPU2} and {CPU2, CPU1}.
	// Sort the slices by "Serial" field
	sort.SliceStable(slice1.Interface(), func(i, j int) bool {
		return slice1.Index(i).FieldByName("Serial").String() < slice1.Index(j).FieldByName("Serial").String()
	})
	sort.SliceStable(slice2.Interface(), func(i, j int) bool {
		return slice2.Index(i).FieldByName("Serial").String() < slice2.Index(j).FieldByName("Serial").String()
	})

	extraElem := "EMAPI EXTRA"
	if slice1.Len() > slice2.Len() {
		extraElem = "FleetDB EXTRA"
		temp := slice2
		slice2 = slice1
		slice1 = temp
	}

	for i := 0; i < slice1.Len(); i++ {
		elem1 := slice1.Index(i)
		elem2 := slice2.Index(i)

		if elem1.Kind() == reflect.Struct {
			CompareComponents(elem1.Interface(), elem2.Interface(), fmt.Sprintf("%s[%d]", path, i))
		} else {
			if elem1.Interface() != elem2.Interface() {
				fmt.Printf("Component %s Field %s[%d]: %v != %v\n", path, typeOfObj, i, elem1.Interface(), elem2.Interface())
			}
		}
	}

	for i := slice1.Len(); i < slice2.Len(); i++ {
		elem := slice2.Index(i)
		fmt.Printf("Component %s [%v] Field %s[%d]: %v\n", path, extraElem, typeOfObj, i, elem.Interface())
	}
}

// Get firmware info
var getComponentGaps = &cobra.Command{
	Use:   "component_gaps",
	Short: "Get gaps between FleetDB inventory and EMAPI hardware",
	Run: func(cmd *cobra.Command, _ []string) {
		theApp := mctl.MustCreateApp(cmd.Context())

		ctx, cancel := context.WithTimeout(cmd.Context(), mctl.CmdTimeout)
		defer cancel()

		fleetClient, err := app.NewFleetDBAPIClient(cmd.Context(), theApp.Config.FleetDBAPI, theApp.Reauth)
		if err != nil {
			log.Fatal(err)
		}

		serverID, err := uuid.Parse(flagsDefinedGetComponentGaps.id)
		if err != nil {
			log.Fatal(err)
		}

		fleetInventory, _, err := fleetClient.GetServerInventory(ctx, serverID, false)
		if err != nil {
			log.Fatal("fleetdb API client returned error: ", err)
		}

		emapiClient, err := staff.NewClient()
		if err != nil {
			log.Fatal(err)
		}

		opts := &staff.ListOptions{
			Excludes: []string{"firmware_version"},
		}
		emapiHardware, _, err := emapiClient.Staff.Components.ListByHardwareID(serverID.String(), opts)
		if err != nil {
			log.Fatal(err)
		}

		fleetComponents, err := convertFleetDBServer(fleetInventory)
		if err != nil {
			log.Fatal(err)
		}

		emapiComponents, err := convertEMAPIServer(emapiHardware)
		if err != nil {
			log.Fatal(err)
		}

		CompareComponents(fleetComponents, emapiComponents, "Component")

		os.Exit(0)
	},
}

func init() {
	cmdGet.AddCommand(getComponentGaps)

	flagsDefinedGetComponentGaps = &getComponentGapsFlags{}
	mctl.AddServerFlag(getComponentGaps, &flagsDefinedGetComponentGaps.id)
	mctl.RequireFlag(getComponentGaps, mctl.ServerFlag)
}
