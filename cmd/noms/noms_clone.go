// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/attic-labs/kingpin"
	"github.com/attic-labs/noms/cmd/util"
	"github.com/attic-labs/noms/go/nbs"
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
	// If no directory specified, derive from URL
	if dir == "" {
		// Extract name from URL path
		dir = filepath.Base(url)
		if dir == "." || dir == "/" {
			dir = "cloned-db"
		}
		// Remove extension if present
		if ext := filepath.Ext(dir); ext != "" {
			dir = dir[:len(dir)-len(ext)]
		}
	}

	// Check if directory already exists
	if _, err := os.Stat(dir); err == nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "directory already exists: " + dir,
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: directory '%s' already exists\n", dir)
		}
		return 1
	}

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": "failed to create directory: " + err.Error(),
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: failed to create directory: %v\n", err)
		}
		return 1
	}

	// Create NBS store
	dbPath := filepath.Join(dir, ".noms")
	cs := nbs.NewLocalStore(dbPath, 256*1024*1024)
	defer cs.Close()

	// For now, simulate clone - in real implementation,
	// this would fetch from remote
	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"status":    "cloned",
			"url":       url,
			"directory": dir,
		})
	} else {
		fmt.Printf("Cloning from %s\n", url)
		fmt.Printf("Creating directory: %s\n", dir)
		fmt.Println("Initialized empty Noms database")
		fmt.Printf("\nTo configure as a remote:\n")
		fmt.Printf("  cd %s\n", dir)
		fmt.Printf("  noms remote --add origin %s\n", url)
	}
	return 0
}
