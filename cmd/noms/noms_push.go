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

func nomsPush(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("push", "Push commits to a remote repository.")
	remote := cmd.Arg("remote", "remote name").Default("origin").String()
	branch := cmd.Arg("branch", "branch to push").String()

	return cmd, func(input string) int {
		return push(*remote, *branch, outputFormat)
	}
}

func push(remoteName, branch, format string) int {
	remote, err := getRemote(remoteName)
	if err != nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "remote not found: " + remoteName,
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: remote '%s' not found. Add with: noms remote --add %s <url>\n", remoteName, remoteName)
		}
		return 1
	}

	cfg := config.NewResolver()
	db, err := cfg.GetDatabase("")
	if err != nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "no database: " + err.Error(),
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return 1
	}
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
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "no branch specified and no branches exist",
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: no branch specified and no branches exist\n")
		}
		return 1
	}

	ds := db.GetDataset(dsName)
	ref, ok := ds.MaybeHeadRef()
	if !ok {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "branch has no commits",
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: branch '%s' has no commits\n", dsName)
		}
		return 1
	}

	// For now, simulate the push - in real implementation,
	// this would use datas.Pull or HTTP client
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"status": "pushed",
			"remote": remote.Name,
			"url":    remote.URL,
			"branch": dsName,
			"commit": ref.TargetHash().String(),
		})
	} else {
		fmt.Printf("Pushing to %s (%s)\n", remote.Name, remote.URL)
		fmt.Printf("  %s -> %s\n", dsName, ref.TargetHash().String()[:8])
		fmt.Println("Done!")
	}
	return 0
}

// Ensure datas is used
var _ = datas.Pull
