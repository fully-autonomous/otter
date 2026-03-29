// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/attic-labs/kingpin"
	"github.com/attic-labs/noms/cmd/util"
)

func nomsTui(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("tui", "Launch interactive terminal UI.")
	return cmd, func(input string) int {
		return launchTui()
	}
}

func launchTui() int {
	// Check if bun is available
	_, err := exec.LookPath("bun")
	if err != nil {
		fmt.Println("Error: bun is required to run the TUI")
		fmt.Println("Install bun: https://bun.sh")
		return 1
	}

	// Check if tui directory exists
	if _, err := os.Stat("tui"); os.IsNotExist(err) {
		fmt.Println("Error: tui directory not found")
		return 1
	}

	// Run the TUI
	cmd := exec.Command("bun", "run", "tui/src/index.ts")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		return 1
	}

	return 0
}
