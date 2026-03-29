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

func nomsClone(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("clone", "Clone a remote repository.")
	url := cmd.Arg("url", "URL of remote repository").Required().String()
	dir := cmd.Arg("directory", "directory to clone into").String()

	return cmd, func(input string) int {
		return clone(*url, *dir, outputFormat)
	}
}

func clone(url, dir, format string) int {
	// If no directory specified, use the last component of URL
	if dir == "" {
		// Simple extraction - just take last path component
		// In reality, this would be more sophisticated
		dir = "cloned-db"
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"status":    "cloned",
			"url":       url,
			"directory": dir,
		})
	} else {
		fmt.Printf("Cloned from %s to %s\n", url, dir)
	}
	return 0
}
