package list

import (
	"encoding/json"
	"fmt"
	"log"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func findAttribute(ns string, attributes []fleetdbapi.Attributes) *fleetdbapi.Attributes {
	for _, attribute := range attributes {
		if attribute.Namespace == ns {
			return &attribute
		}
	}

	return nil
}

func printJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
