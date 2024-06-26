package parse

import (
	"encoding/json"
	"bufio"
	"log"
	"os"
	"time"
	"github.com/google/uuid"
	"errors"
	"context"
	"fmt"

	"github.com/spf13/cobra"

	mctl "github.com/metal-toolbox/mctl/cmd"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	rfleetdb "github.com/metal-toolbox/rivets/fleetdb"
	rt "github.com/metal-toolbox/rivets/types"
	"github.com/metal-toolbox/mctl/internal/app"
)

type SplunkEntry struct {
	Preview bool `json:"preview"`
	LastRow bool `json:"lastrow,omitempty"`
	Result  struct {
		Raw           string    `json:"_raw"`
		Time          time.Time `json:"_time"`
		Action        string    `json:"action"`
		AssetID       string    `json:"assetID"`
		Bmc           string    `json:"bmc"`
		ClusterClass  string    `json:"cluster_class"`
		ClusterEnv    string    `json:"cluster_env"`
		ClusterName   string    `json:"cluster_name"`
		Component     string    `json:"component"`
		ConditionID   string    `json:"conditionID"`
		ContainerID   string    `json:"container_id"`
		ContainerName string    `json:"container_name"`
		ControllerID  string    `json:"controllerID"`
		Eventtype     []string  `json:"eventtype"`
		File          string    `json:"file"`
		Function      string    `json:"function"`
		Fwversion     string    `json:"fwversion"`
		Host          string    `json:"host"`
		Index         string    `json:"index"`
		LabelK8SApp   string    `json:"label_k8s-app"`
		Level         string    `json:"level"`
		Line          string    `json:"line"`
		Linecount     string    `json:"linecount"`
		Msg           string    `json:"msg"`
		Namespace     string    `json:"namespace"`
		Pod           string    `json:"pod"`
		PodUID        string    `json:"pod_uid"`
		Punct         string    `json:"punct"`
		Source        string    `json:"source"`
		Sourcetype    string    `json:"sourcetype"`
		SplunkServer  string    `json:"splunk_server"`
		Time0         time.Time `json:"time"`
	} `json:"result"`
}

var flasherCmd = &cobra.Command{
	Use: "flasher",
	Short: "Parse splunk logs from flasher",
	Run: func(cmd *cobra.Command, _ []string) {
		flasherAction(cmd.Context())
	},
}

func flasherAction(ctx context.Context) {
	theApp := mctl.MustCreateApp(ctx)

	client, err := app.NewFleetDBAPIClient(ctx, theApp.Config.FleetDBAPI, theApp.Reauth)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readSplunkFile()
	if err != nil {
		log.Fatalf("Failed to parse:\n\n%s\n\n", err.Error())
	}

	fd, err := os.OpenFile(flagsDefinedParseAction.OutputCSVFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err.Error())
	}
	defer fd.Close()

	line := fmt.Sprintf("Server, ServerID, Vendor, Model, Facility, Mesage, Cluster, Component, ConditionID, BMC IP, Flasher Log Line, Time\n")
	_, err = fd.Write([]byte(line))
	if err != nil {
		log.Fatalf("Failed to write line: %s", err.Error())
	}

	entryCount := len(entries)
	for i, entry := range entries {
		if entry.Result.AssetID == "" {
			fmt.Printf("Skipping entry due to no AssetID: \n%+v\n\n", entry)
			continue
		}

		id, err := uuid.Parse(entry.Result.AssetID)
		if err != nil {
			log.Fatalf("Failed to parse uuid: %s", err.Error())
		}

		server, err := server(ctx, client, id)
		if err != nil {
			log.Fatalf("Failed to get server: %s", err.Error())
		}

		if entry.Result.Component == "" {
			entry.Result.Component = "NA"
		}

		line = fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s:%s, %s\n",
			server.Name,
			entry.Result.AssetID,
			server.Vendor,
			server.Model,
			server.Facility,
			entry.Result.Msg,
			entry.Result.ClusterName,
			entry.Result.Component,
			entry.Result.ConditionID,
			entry.Result.Bmc,
			entry.Result.File,
			entry.Result.Line,
			entry.Result.Time.Format("01-02-2006 15:04:05"))

		_, err = fd.Write([]byte(line))
		if err != nil {
			log.Fatalf("Failed to write line: %s", err.Error())
		}

		if i % 100 == 0 {
			fmt.Printf("Parsed %d/%d...\n", i, entryCount)
		}
	}
}

func extend(err error, new string) error {
	return errors.Join(err, errors.New(new))
}

func readSplunkFile() ([]SplunkEntry, error) {
	file, err := os.Open(flagsDefinedParseAction.JsonFileToParse)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var entries []SplunkEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		var entry SplunkEntry
		err := json.Unmarshal(scanner.Bytes(), &entry)
		if err != nil {
			return nil, err
			log.Fatalf("Failed to parse:\n\n%s\n\n - %s\n", err.Error(), scanner.Bytes())
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func server(ctx context.Context, client *fleetdbapi.Client, id uuid.UUID) (*rt.Server, error) {
	server, _, err := client.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	cserver := rfleetdb.ConvertServer(server)

	return cserver, nil
}