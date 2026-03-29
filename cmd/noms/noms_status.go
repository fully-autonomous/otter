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

func nomsStatus(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("status", "Show the working tree status.")

	return cmd, func(input string) int {
		cfg := config.NewResolver()
		db, err := cfg.GetDatabase("")
		d.CheckErrorNoUsage(err)
		defer db.Close()

		return showStatus(db, outputFormat)
	}
}

func showStatus(db datas.Database, format string) int {
	datasets := db.Datasets()

	branches := []string{}
	staged := []string{}
	modified := []string{}

	datasets.IterAll(func(k, v types.Value) {
		if k, ok := k.(types.String); ok {
			branches = append(branches, string(k))
			modified = append(modified, string(k)) // All datasets are "modified" if they exist
		}
	})

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"current_branch": currentBranch(db),
			"branches":       branches,
			"staged":         staged,
			"modified":       modified,
			"clean":          len(branches) == 0,
		})
	} else {
		fmt.Printf("Current branch: %s\n\n", currentBranch(db))

		if len(branches) == 0 {
			fmt.Println("No commits yet")
			fmt.Println("\n(use 'noms init' to initialize)")
			return 0
		}

		fmt.Println("Changes to be committed:")
		if len(staged) == 0 {
			fmt.Println("  (none)")
		} else {
			for _, s := range staged {
				fmt.Printf("  %s\n", s)
			}
		}

		fmt.Println("\nChanges not staged for commit:")
		if len(modified) == 0 {
			fmt.Println("  (none)")
		} else {
			for _, m := range modified {
				fmt.Printf("  %s\n", m)
			}
		}
	}
	return 0
}

func currentBranch(db datas.Database) string {
	// Noms doesn't have a "current branch" concept like Git
	// It treats datasets as branches. For now, show first dataset.
	datasets := db.Datasets()
	var first string
	datasets.IterAll(func(k, v types.Value) {
		if first == "" {
			if k, ok := k.(types.String); ok {
				first = string(k)
			}
		}
	})
	if first == "" {
		return "(no branches)"
	}
	return first
}
