// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/attic-labs/kingpin"
	"github.com/attic-labs/noms/cmd/util"
	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/nbs"
)

func nomsInit(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("init", "Initialize a new Noms database in the current directory.")
	force := cmd.Flag("force", "re-initialize an existing noms database").Short('f').Bool()

	return cmd, func(input string) int {
		// Check if already initialized
		_, err := config.FindNomsConfig()
		if err == nil && !*force {
			fmt.Fprintf(os.Stderr, "Error: already initialized. Use --force to reinitialize.\n")
			return 1
		}

		// Get current directory for database path
		cwd, err := os.Getwd()
		d.CheckErrorNoUsage(err)
		dbPath := filepath.Join(cwd, ".noms")

		// Create database directory
		err = os.MkdirAll(dbPath, 0755)
		d.CheckErrorNoUsage(err)

		// Create NBS store
		cs := nbs.NewLocalStore(dbPath, 256*1024*1024) // 256MB mem table
		defer cs.Close()

		// Create config file in current directory (local config)
		cfg := &config.Config{
			Db: map[string]config.DbConfig{
				"default": {
					Url: dbPath,
				},
			},
		}

		// Write config to current directory
		configPath := filepath.Join(cwd, config.NomsConfigFile)
		data := cfg.WriteableString()
		err = os.WriteFile(configPath, []byte(data), 0644)
		d.CheckErrorNoUsage(err)

		fmt.Printf("Initialized empty Noms database at: %s\n", dbPath)
		fmt.Printf("Config written to: %s\n", configPath)
		return 0
	}
}
