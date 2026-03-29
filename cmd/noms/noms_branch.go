// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/attic-labs/kingpin"
	"github.com/attic-labs/noms/cmd/util"
	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
)

func nomsBranch(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("branch", "List, create, or delete branches (datasets).")

	list := cmd.Flag("list", "list all branches").Short('l').Bool()
	delete := cmd.Flag("delete", "delete a branch").Short('d').Bool()
	create := cmd.Flag("create", "create a new branch").Short('c').Bool()
	branchName := cmd.Arg("name", "branch name").String()
	startPoint := cmd.Arg("start-point", "starting commit (default: current head)").String()

	return cmd, func(input string) int {
		cfg := config.NewResolver()
		db, err := cfg.GetDatabase("")
		if err != nil {
			if outputFormat == "json" {
				json.NewEncoder(os.Stdout).Encode(map[string]string{
					"error": "no database initialized: " + err.Error(),
				})
			} else {
				fmt.Println("No Noms database found. Run 'noms init' to initialize.")
			}
			return 0
		}
		defer db.Close()

		// List branches
		if *list {
			return listBranches(db, outputFormat)
		}

		// Delete branch
		if *delete {
			if *branchName == "" {
				fmt.Fprintf(os.Stderr, "Error: branch name required for delete\n")
				return 1
			}
			return deleteBranch(db, *branchName, outputFormat)
		}

		// Create branch
		if *create || *branchName != "" {
			if *branchName == "" {
				fmt.Fprintf(os.Stderr, "Error: branch name required\n")
				return 1
			}
			return createBranch(db, *branchName, *startPoint, outputFormat)
		}

		// No args, list branches
		return listBranches(db, outputFormat)
	}
}

func listBranches(db datas.Database, format string) int {
	datasets := db.Datasets()
	branches := []map[string]string{}

	datasets.IterAll(func(k, v types.Value) {
		if k, ok := k.(types.String); ok {
			if v, ok := v.(types.Ref); ok {
				branches = append(branches, map[string]string{
					"name": string(k),
					"hash": v.TargetHash().String(),
				})
			}
		}
	})

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(branches)
	} else {
		for _, b := range branches {
			fmt.Printf("  %s -> %s\n", b["name"], b["hash"][:8])
		}
		if len(branches) == 0 {
			fmt.Println("No branches")
		}
	}
	return 0
}

func createBranch(db datas.Database, name, startPoint string, format string) int {
	if startPoint != "" {
		fmt.Fprintf(os.Stderr, "Error: start-point not yet implemented\n")
		return 1
	}

	ds := db.GetDataset(name)

	if _, ok := ds.MaybeHeadRef(); ok {
		fmt.Fprintf(os.Stderr, "Error: branch '%s' already exists\n", name)
		return 1
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"status": "created",
			"name":   name,
		})
	} else {
		fmt.Printf("Created branch: %s\n", name)
	}
	return 0
}

func deleteBranch(db datas.Database, name, format string) int {
	ds := db.GetDataset(name)

	if _, ok := ds.MaybeHeadRef(); !ok {
		fmt.Fprintf(os.Stderr, "Error: branch '%s' does not exist\n", name)
		return 1
	}

	_, err := db.Delete(ds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting branch: %v\n", err)
		return 1
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"status": "deleted",
			"name":   name,
		})
	} else {
		fmt.Printf("Deleted branch: %s\n", name)
	}
	return 0
}
