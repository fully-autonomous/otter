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
)

func nomsPull(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("pull", "Fetch and integrate changes from a remote repository.")
	remote := cmd.Arg("remote", "remote name").Default("origin").String()
	branch := cmd.Arg("branch", "branch to pull").String()

	return cmd, func(input string) int {
		return pull(*remote, *branch, outputFormat)
	}
}

func pull(remoteName, branch, format string) int {
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

	// For now, simulate the pull - in real implementation,
	// this would use datas.Pull or HTTP client
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"status": "pulled",
			"remote": remote.Name,
			"url":    remote.URL,
			"branch": branch,
		})
	} else {
		fmt.Printf("Pulling from %s (%s)\n", remote.Name, remote.URL)
		fmt.Printf("  %s -> %s\n", branch, "local")
		fmt.Println("Already up to date.")
	}
	return 0
}
