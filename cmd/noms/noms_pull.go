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
)

func nomsPull(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("pull", "Fetch and integrate changes from a remote repository.")
	remote := cmd.Arg("remote", "remote name").Default("origin").String()
	branch := cmd.Arg("branch", "branch to pull").String()

	return cmd, func(input string) int {
		return pull(*remote, *branch, outputFormat)
	}
}

func pull(remote, branch, format string) int {
	url, ok := remotes[remote]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: remote '%s' not found\n", remote)
		return 1
	}

	// TODO: Actually pull from remote URL
	// For now, just simulate

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"status": "pulled",
			"remote": remote,
			"url":    url,
			"branch": branch,
		})
	} else {
		fmt.Printf("Pulled from %s (%s)\n", remote, url)
	}
	return 0
}
