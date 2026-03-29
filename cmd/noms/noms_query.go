// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/attic-labs/kingpin"
	"github.com/attic-labs/noms/cmd/util"
	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/datas"
)

func nomsQuery(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("query", "Query system tables (noms_log, noms_branches, noms_datasets).")
	table := cmd.Arg("table", "table to query (noms_log, noms_branches, noms_datasets)").Default("noms_log").String()

	return cmd, func(input string) int {
		return runQuery(*table, outputFormat)
	}
}

func runQuery(table, format string) int {
	cfg := config.NewResolver()
	db, err := cfg.GetDatabase("")
	if err != nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "no database initialized: " + err.Error(),
			})
		} else {
			fmt.Println("No Noms database found. Run 'noms init' to initialize.")
		}
		return 1
	}
	defer db.Close()

	// Use reflection to access the concrete *database type
	dbVal := reflect.ValueOf(db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}

	switch table {
	case "noms_log":
		// Try to call GetLog method
		method := dbVal.MethodByName("GetLog")
		if !method.IsValid() {
			return showLogFallback(format)
		}
		results := method.Call([]reflect.Value{reflect.ValueOf(""), reflect.ValueOf(50)})
		if len(results) == 0 {
			return showLogFallback(format)
		}
		logs, ok := results[0].Interface().([]datas.NomsLogEntry)
		if !ok {
			return showLogFallback(format)
		}
		return printLogEntries(logs, format)

	case "noms_branches":
		method := dbVal.MethodByName("GetBranches")
		if !method.IsValid() {
			return showBranchesFallback(format)
		}
		results := method.Call(nil)
		if len(results) == 0 {
			return showBranchesFallback(format)
		}
		branches, ok := results[0].Interface().([]datas.NomsBranchEntry)
		if !ok {
			return showBranchesFallback(format)
		}
		return printBranchEntries(branches, format)

	case "noms_datasets":
		method := dbVal.MethodByName("GetDatasets")
		if !method.IsValid() {
			return showDatasetsFallback(format)
		}
		results := method.Call(nil)
		if len(results) == 0 {
			return showDatasetsFallback(format)
		}
		datasets, ok := results[0].Interface().([]datas.NomsDatasetEntry)
		if !ok {
			return showDatasetsFallback(format)
		}
		return printDatasetEntries(datasets, format)

	default:
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "unknown table: " + table,
			})
		} else {
			fmt.Printf("Unknown table: %s\n", table)
			fmt.Println("Available tables: noms_log, noms_branches, noms_datasets")
		}
		return 1
	}
}

func printLogEntries(logs []datas.NomsLogEntry, format string) int {
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"noms_log": logs,
		})
	} else {
		if len(logs) == 0 {
			fmt.Println("No commits yet")
			return 0
		}
		fmt.Println("Commit History:")
		fmt.Println("================")
		for _, entry := range logs {
			hash := entry.Hash
			if len(hash) > 8 {
				hash = hash[:8]
			}
			fmt.Printf("\nHash: %s\n", hash)
			fmt.Printf("Dataset: %s\n", entry.Dataset)
			fmt.Printf("Author: %s\n", entry.Author)
			fmt.Printf("Date: %s\n", entry.Date.Format("2006-01-02 15:04:05"))
			fmt.Printf("Message: %s\n", entry.Message)
		}
	}
	return 0
}

func printBranchEntries(branches []datas.NomsBranchEntry, format string) int {
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"branches": branches,
		})
	} else {
		fmt.Println("Branches:")
		fmt.Println("=========")
		for _, b := range branches {
			head := b.Head
			if len(head) > 8 {
				head = head[:8]
			}
			fmt.Printf("  %s -> %s\n", b.Name, head)
		}
		if len(branches) == 0 {
			fmt.Println("  (none)")
		}
	}
	return 0
}

func printDatasetEntries(datasets []datas.NomsDatasetEntry, format string) int {
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"datasets": datasets,
		})
	} else {
		fmt.Println("Datasets:")
		fmt.Println("=========")
		for _, d := range datasets {
			head := d.Head
			if len(head) > 8 {
				head = head[:8]
			}
			fmt.Printf("  %s -> %s\n", d.Name, head)
		}
		if len(datasets) == 0 {
			fmt.Println("  (none)")
		}
	}
	return 0
}

func showLogFallback(format string) int {
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"note": "use 'noms log' to view commit history",
		})
	} else {
		fmt.Println("Use 'noms log' to view commit history")
	}
	return 0
}

func showBranchesFallback(format string) int {
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"note": "use 'noms branch --list' to view branches",
		})
	} else {
		fmt.Println("Use 'noms branch --list' to view branches")
	}
	return 0
}

func showDatasetsFallback(format string) int {
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"note": "use 'noms branch --list' to view datasets",
		})
	} else {
		fmt.Println("Use 'noms branch --list' to view datasets")
	}
	return 0
}
