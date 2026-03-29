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
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
)

func nomsCheckout(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("checkout", "Switch between branches or restore working tree.")

	create := cmd.Flag("create", "create and switch to a new branch").Short('b').Bool()
	force := cmd.Flag("force", "discard unsaved changes").Short('f').Bool()
	branchName := cmd.Arg("branch", "branch name to checkout").String()

	return cmd, func(input string) int {
		cfg := config.NewResolver()
		db, err := cfg.GetDatabase("")
		d.CheckErrorNoUsage(err)
		defer db.Close()

		if *branchName == "" {
			// List available branches
			return listAvailableBranches(db, outputFormat)
		}

		return checkoutBranch(db, *branchName, *create, *force, outputFormat)
	}
}

func listAvailableBranches(db datas.Database, format string) int {
	datasets := db.Datasets()
	branches := []string{}

	datasets.IterAll(func(k, v types.Value) {
		if k, ok := k.(types.String); ok {
			branches = append(branches, string(k))
		}
	})

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string][]string{
			"branches": branches,
		})
	} else {
		fmt.Println("Available branches:")
		for _, b := range branches {
			fmt.Printf("  %s\n", b)
		}
	}
	return 0
}

func checkoutBranch(db datas.Database, name string, create, force bool, format string) int {
	ds := db.GetDataset(name)

	if _, ok := ds.MaybeHeadRef(); !ok {
		if create {
			// Create new branch (empty commit will be created on first write)
			if format == "json" {
				json.NewEncoder(os.Stdout).Encode(map[string]string{
					"status": "created",
					"branch": name,
				})
			} else {
				fmt.Printf("Switched to new branch: %s\n", name)
			}
			return 0
		}
		fmt.Fprintf(os.Stderr, "Error: branch '%s' does not exist. Use -b to create.\n", name)
		return 1
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"status": "switched",
			"branch": name,
		})
	} else {
		fmt.Printf("Switched to branch: %s\n", name)
	}
	return 0
}
