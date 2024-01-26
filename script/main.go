package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	coapiv1 "github.com/metal-toolbox/conditionorc/pkg/api/v1/types"
	"github.com/metal-toolbox/mctl/internal/app"
	"github.com/metal-toolbox/mctl/pkg/model"
)

type ParsedHardware struct {
	ID           string
	FacilityID   string
	IPMIAddress  string
	IPMIUsername string
	IPMIPassword string
}

func main() {
	// The provided text data
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	apps := app.App{
		Config: &model.Config{
			Conditions: &model.ConfigOIDC{
				Endpoint:        "http://localhost:9001",
				Disable:         true,
				ClientID:        "",
				IssuerEndpoint:  "",
				Scopes:          []string{},
				PkceCallbackURL: "",
			},
		},
		Reauth: false,
	}
	client, err := app.NewConditionsClient(context.Background(), apps.Config.Conditions, apps.Reauth)
	if err != nil {
		fmt.Printf("failed to connect to conditionorc %v", err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var hardwareMap map[string]interface{}

		// Parse JSON string into the map
		err := json.Unmarshal([]byte(line), &hardwareMap)
		if err != nil {
			fmt.Printf("Error parsing JSON: %w", err)
			return
		}

		if hardwareMap["type"] != "server" {
			continue
		}

		var enrollHardware = ParsedHardware{}
		enrollHardware.ID = hardwareMap["id"].(string)

		if facilityMap, ok := hardwareMap["facility"].(map[string]interface{}); ok {
			enrollHardware.FacilityID = facilityMap["id"].(string)
		}

		if dataMap, ok := hardwareMap["data"].(map[string]interface{}); ok {
			if ipmiMap, ok := dataMap["ipmi"].(map[string]interface{}); ok {
				enrollHardware.IPMIAddress = ipmiMap["address"].(string)
				enrollHardware.IPMIUsername = ipmiMap["username"].(string)
				enrollHardware.IPMIPassword = ipmiMap["password"].(string)
			}
		}
		fmt.Println()

		params, err := json.Marshal(coapiv1.AddServerParams{
			Facility: enrollHardware.FacilityID,
			IP:       enrollHardware.IPMIAddress,
			Username: enrollHardware.IPMIUsername,
			Password: enrollHardware.IPMIPassword,
		})
		if err != nil {
			fmt.Printf("Error marshal %v", err)
		}

		conditionCreate := coapiv1.ConditionCreate{
			Parameters: params,
		}
		fmt.Printf("conditionCreate %v\n", conditionCreate)

		_, err = client.ServerEnroll(context.Background(), enrollHardware.ID, conditionCreate)
		if err != nil {
			fmt.Printf("failed to enroll server %v\n", err)
			return
		}
	}
}
