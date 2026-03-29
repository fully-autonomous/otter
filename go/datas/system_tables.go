// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package datas

import (
	"fmt"
	"time"

	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/types"
)

// SystemTable types for querying metadata

// NomsLogEntry represents a single commit in the log
type NomsLogEntry struct {
	Hash    string    `json:"hash"`
	Parents string    `json:"parents"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	Dataset string    `json:"dataset"`
}

// NomsBranchEntry represents a branch reference
type NomsBranchEntry struct {
	Name      string    `json:"name"`
	Head      string    `json:"head"`
	CreatedAt time.Time `json:"created_at"`
}

// NomsDatasetEntry represents a dataset
type NomsDatasetEntry struct {
	Name     string `json:"name"`
	Head     string `json:"head"`
	RootHash string `json:"root_hash"`
}

// GetLog returns commit history for all or specific datasets
func (db *database) GetLog(datasetID string, limit int) ([]NomsLogEntry, error) {
	datasets := db.Datasets()
	var entries []NomsLogEntry

	datasets.IterAll(func(k, v types.Value) {
		dsName := string(k.(types.String))

		// Filter by dataset if specified
		if datasetID != "" && dsName != datasetID {
			return
		}

		ref := v.(types.Ref)
		commitHash := ref.TargetHash()

		// Walk the commit history
		visited := hash.HashSet{}
		count := 0

		for count < limit || limit == 0 {
			if visited.Has(commitHash) {
				break
			}
			visited.Insert(commitHash)

			commit := db.ReadValue(commitHash)
			if commit == nil {
				break
			}

			commitStruct := commit.(types.Struct)

			// Extract metadata
			var message, author string
			var date time.Time
			meta := commitStruct.Get(MetaField)
			if metaStruct, ok := meta.(types.Struct); ok {
				if m, ok := metaStruct.MaybeGet("message"); ok {
					message = string(m.(types.String))
				}
				if a, ok := metaStruct.MaybeGet("author"); ok {
					author = string(a.(types.String))
				}
				if t, ok := metaStruct.MaybeGet("date"); ok {
					if dateStr, ok := t.(types.String); ok {
						// Date is stored in RFC3339 format
						if parsed, err := time.Parse(time.RFC3339, string(dateStr)); err == nil {
							date = parsed
						}
					}
				}
			}

			// Get parents
			parents := commitStruct.Get(ParentsField).(types.Set)
			var parentStr string
			var parentHash hash.Hash
			parents.IterAll(func(p types.Value) {
				parentHash = p.(types.Ref).TargetHash()
				if parentStr != "" {
					parentStr += ","
				}
				parentStr += parentHash.String()
			})

			entries = append(entries, NomsLogEntry{
				Hash:    commitHash.String(),
				Parents: parentStr,
				Message: message,
				Author:  author,
				Date:    date,
				Dataset: dsName,
			})

			count++

			// Move to parent
			if parentHash.IsEmpty() {
				break
			}
			commitHash = parentHash
		}
	})

	return entries, nil
}

// GetBranches returns all branches (datasets)
func (db *database) GetBranches() ([]NomsBranchEntry, error) {
	datasets := db.Datasets()
	var branches []NomsBranchEntry

	datasets.IterAll(func(k, v types.Value) {
		name := string(k.(types.String))
		ref := v.(types.Ref)

		branches = append(branches, NomsBranchEntry{
			Name:      name,
			Head:      ref.TargetHash().String(),
			CreatedAt: time.Now(), // TODO: Track creation time
		})
	})

	return branches, nil
}

// GetDatasets returns all datasets with their current heads
func (db *database) GetDatasets() ([]NomsDatasetEntry, error) {
	datasets := db.Datasets()
	var entries []NomsDatasetEntry

	datasets.IterAll(func(k, v types.Value) {
		name := string(k.(types.String))
		ref := v.(types.Ref)

		entries = append(entries, NomsDatasetEntry{
			Name:     name,
			Head:     ref.TargetHash().String(),
			RootHash: ref.TargetHash().String(),
		})
	})

	return entries, nil
}

// QueryLog is a simple query interface for system tables
func (db *database) QueryLog(query string) ([]NomsLogEntry, error) {
	// Parse simple queries: "noms_log", "noms_log where dataset = 'x'"
	// For now, support basic filtering

	datasetID := ""
	fmt.Sscanf(query, "noms_log where dataset = '%s'", &datasetID)

	return db.GetLog(datasetID, 100)
}
