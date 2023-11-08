package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/dakota-xyz/recovery/cmd"
)

func main() {
	shard1 := flag.String("shard1", "", "Location of the file containing the first shard")
	shard2 := flag.String("shard2", "", "Location of the file containing the second shard")
	keyMap := flag.String("keymap", "", "Location of the JSON file containing the key map")
	targetCSV := flag.String("target", "keys.csv", "Target CSV file")
	flag.Parse()

	checkRequiredFlag(shard1, "shard1 is required")
	checkRequiredFlag(shard2, "shard2 is required")
	checkRequiredFlag(keyMap, "keyMap is required")

	slog.Info("Initiating recovery")
	slog.Debug("Paramaters in use",
		"shard1_location", *shard1,
		"shard2_location", *shard2,
		"keymap_location", *keyMap,
	)

	shard1File, err := os.Open(*shard1)
	checkError(err, "Failed to open shard1")
	defer shard1File.Close()
	shard2File, err := os.Open(*shard2)
	checkError(err, "Failed to open shard2")
	defer shard2File.Close()
	keyMapFile, err := os.Open(*keyMap)
	checkError(err, "Failed to open keymap")
	defer keyMapFile.Close()

	targetCSVFile, err := os.Create(*targetCSV)
	checkError(err, "Failed to open target CSV")
	defer targetCSVFile.Close()

	err = cmd.Recover(targetCSVFile, shard1File, shard2File, keyMapFile)
	checkError(err, "Failed to recover")

	slog.Info(fmt.Sprintf("Recovery complete. Results saved to %s", targetCSVFile.Name()))
}

func checkRequiredFlag(val *string, message string) {
	if val == nil || *val == "" {
		slog.Error(message)
		os.Exit(1)
	}
}

func checkError(err error, message string) {
	if err != nil {
		slog.Error(message, "error", err)
		os.Exit(2)
	}
}
