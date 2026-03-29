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

type remoteConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var remotes = map[string]string{}

func nomsRemote(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("remote", "Manage set of tracked remotes.")

	list := cmd.Flag("show-urls", "show remote URLs").Bool()
	add := cmd.Flag("add", "add a remote").Bool()
	remove := cmd.Flag("remove", "remove a remote").Bool()
	remoteName := cmd.Arg("name", "remote name").String()
	remoteURL := cmd.Arg("url", "remote URL").String()

	return cmd, func(input string) int {
		// List remotes
		if *list {
			return listRemotes(outputFormat)
		}

		// Add remote
		if *add {
			if *remoteName == "" || *remoteURL == "" {
				fmt.Fprintf(os.Stderr, "Error: both name and URL required for add\n")
				return 1
			}
			return addRemote(*remoteName, *remoteURL, outputFormat)
		}

		// Remove remote
		if *remove {
			if *remoteName == "" {
				fmt.Fprintf(os.Stderr, "Error: remote name required for remove\n")
				return 1
			}
			return removeRemote(*remoteName, outputFormat)
		}

		// No args, list remotes
		return listRemotes(outputFormat)
	}
}

func listRemotes(format string) int {
	names := []string{}
	urls := []string{}
	for name, url := range remotes {
		names = append(names, name)
		urls = append(urls, url)
	}

	if format == "json" {
		remotesList := []map[string]string{}
		for name, url := range remotes {
			remotesList = append(remotesList, map[string]string{
				"name": name,
				"url":  url,
			})
		}
		json.NewEncoder(os.Stdout).Encode(remotesList)
	} else {
		if len(remotes) == 0 {
			fmt.Println("No remotes configured")
			return 0
		}
		for name, url := range remotes {
			fmt.Printf("  %s -> %s\n", name, url)
		}
	}
	return 0
}

func addRemote(name, url string, format string) int {
	remotes[name] = url

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"status": "added",
			"name":   name,
			"url":    url,
		})
	} else {
		fmt.Printf("Added remote: %s -> %s\n", name, url)
	}
	return 0
}

func removeRemote(name string, format string) int {
	if _, ok := remotes[name]; !ok {
		fmt.Fprintf(os.Stderr, "Error: remote '%s' does not exist\n", name)
		return 1
	}

	delete(remotes, name)

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"status": "removed",
			"name":   name,
		})
	} else {
		fmt.Printf("Removed remote: %s\n", name)
	}
	return 0
}
