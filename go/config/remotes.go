// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// RemoteType represents the type of remote
type RemoteType string

const (
	RemoteTypeHTTP  RemoteType = "http"
	RemoteTypeHTTPS RemoteType = "https"
	RemoteTypeFile  RemoteType = "file"
	RemoteTypeS3    RemoteType = "s3"
)

// Remote represents a remote repository
type Remote struct {
	Name string     `json:"name"`
	URL  string     `json:"url"`
	Type RemoteType `json:"type"`
}

// RemoteConfig holds all remotes for a database
type RemoteConfig struct {
	Remotes map[string]Remote `json:"remotes"`
}

// DefaultRemoteConfig returns a new RemoteConfig
func DefaultRemoteConfig() *RemoteConfig {
	return &RemoteConfig{
		Remotes: make(map[string]Remote),
	}
}

// AddRemote adds a remote to the config
func (rc *RemoteConfig) AddRemote(name, url string) error {
	remoteType := detectRemoteType(url)

	rc.Remotes[name] = Remote{
		Name: name,
		URL:  url,
		Type: remoteType,
	}
	return nil
}

// RemoveRemote removes a remote by name
func (rc *RemoteConfig) RemoveRemote(name string) error {
	if _, ok := rc.Remotes[name]; !ok {
		return fmt.Errorf("remote '%s' not found", name)
	}
	delete(rc.Remotes, name)
	return nil
}

// GetRemote returns a remote by name
func (rc *RemoteConfig) GetRemote(name string) (Remote, error) {
	if r, ok := rc.Remotes[name]; ok {
		return r, nil
	}
	return Remote{}, fmt.Errorf("remote '%s' not found", name)
}

// ListRemotes returns all remotes
func (rc *RemoteConfig) ListRemotes() []Remote {
	remotes := make([]Remote, 0, len(rc.Remotes))
	for _, r := range rc.Remotes {
		remotes = append(remotes, r)
	}
	return remotes
}

// LoadRemotes loads remotes from the config file
func LoadRemotes(configFile string) (*RemoteConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return DefaultRemoteConfig(), nil
	}

	var rc RemoteConfig
	if err := json.Unmarshal(data, &rc); err != nil {
		return nil, err
	}

	if rc.Remotes == nil {
		rc.Remotes = make(map[string]Remote)
	}

	return &rc, nil
}

// SaveRemotes saves remotes to the config file
func SaveRemotes(configFile string, rc *RemoteConfig) error {
	data, err := json.MarshalIndent(rc, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// detectRemoteType determines the remote type from URL
func detectRemoteType(url string) RemoteType {
	if len(url) >= 5 && url[:5] == "https" {
		return RemoteTypeHTTPS
	}
	if len(url) >= 4 && url[:4] == "http" {
		return RemoteTypeHTTP
	}
	if len(url) >= 3 && url[:3] == "s3:" {
		return RemoteTypeS3
	}
	return RemoteTypeFile
}

// RemoteConfigPath returns the path to the remotes config file
func RemoteConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".nomsremotes"
	}
	return filepath.Join(home, ".noms", "remotes.json")
}
