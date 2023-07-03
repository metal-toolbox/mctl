package get

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
)

// this is a pure side-effect
func writeResults(data ...any) {
	switch output {
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
