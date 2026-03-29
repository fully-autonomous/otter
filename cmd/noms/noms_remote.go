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

func loadRemotesConfig() *config.RemoteConfig {
	rc, err := config.LoadRemotes(config.RemoteConfigPath())
	if err != nil {
		return config.DefaultRemoteConfig()
	}
	return rc
}

func saveRemotesConfig(rc *config.RemoteConfig) {
	if err := config.SaveRemotes(config.RemoteConfigPath(), rc); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save remote config: %v\n", err)
	}
}

func listRemotes(format string) int {
	rc := loadRemotesConfig()
	remotes := rc.ListRemotes()

	if format == "json" {
		remotesList := []map[string]string{}
		for _, r := range remotes {
			remotesList = append(remotesList, map[string]string{
				"name": r.Name,
				"url":  r.URL,
			})
		}
		json.NewEncoder(os.Stdout).Encode(remotesList)
	} else {
		if len(remotes) == 0 {
			fmt.Println("No remotes configured")
			fmt.Println("\nTo add a remote:")
			fmt.Println("  noms remote --add origin https://example.com/db")
			return 0
		}
		for _, r := range remotes {
			fmt.Printf("  %s -> %s\n", r.Name, r.URL)
		}
	}
	return 0
}

func addRemote(name, url string, format string) int {
	rc := loadRemotesConfig()

	if err := rc.AddRemote(name, url); err != nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": err.Error(),
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return 1
	}

	saveRemotesConfig(rc)

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
	rc := loadRemotesConfig()

	if err := rc.RemoveRemote(name); err != nil {
		if format == "json" {
			json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": err.Error(),
			})
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return 1
	}

	saveRemotesConfig(rc)

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

func getRemote(name string) (*config.Remote, error) {
	rc := loadRemotesConfig()
	remote, err := rc.GetRemote(name)
	if err != nil {
		return nil, err
	}
	return &remote, nil
}
