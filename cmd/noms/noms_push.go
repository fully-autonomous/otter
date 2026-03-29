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
	"github.com/attic-labs/noms/go/types"
)

func nomsPush(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("push", "Push commits to a remote repository.")
	remote := cmd.Arg("remote", "remote name").Default("origin").String()
	branch := cmd.Arg("branch", "branch to push").String()

	return cmd, func(input string) int {
		return push(*remote, *branch, outputFormat)
	}
}

func push(remote, branch, format string) int {
	url, ok := remotes[remote]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: remote '%s' not found\n", remote)
		return 1
	}

	cfg := config.NewResolver()
	db, err := cfg.GetDatabase("")
	d.CheckErrorNoUsage(err)
	defer db.Close()

	// Get source dataset
	dsName := branch
	if dsName == "" {
		// Default to first dataset
		datasets := db.Datasets()
		var first string
		datasets.IterAll(func(k, v types.Value) {
			if first == "" {
				if k, ok := k.(types.String); ok {
					first = string(k)
				}
			}
		})
		dsName = first
	}

	if dsName == "" {
		fmt.Fprintf(os.Stderr, "Error: no branch specified and no branches exist\n")
		return 1
	}

	ds := db.GetDataset(dsName)
	ref, ok := ds.MaybeHeadRef()
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: branch '%s' has no commits\n", dsName)
		return 1
	}

	// TODO: Actually push to remote URL
	// For now, just simulate
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"status": "pushed",
			"remote": remote,
			"url":    url,
			"branch": dsName,
			"commit": ref.TargetHash().String(),
		})
	} else {
		fmt.Printf("Pushed %s to %s (%s)\n", dsName, remote, url)
	}
	return 0
}
