// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package datas

import (
	"time"

	"github.com/attic-labs/noms/go/hash"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/util"
)

// ValueAt resolves a reference string to a value
// ref can be: commit hash, branch name, or time expression
func (db *database) ValueAt(ref string) (types.Value, error) {
	refType, t, err := util.ParseRef(ref)
	if err != nil {
		return nil, err
	}

	switch refType {
	case util.HashRef, util.ShortHashRef:
		h := hash.Parse(ref)
		return db.ReadValue(h), nil

	case util.BranchRef:
		ds := db.GetDataset(ref)
		if head, ok := ds.MaybeHeadRef(); ok {
			return db.ReadValue(head.TargetHash()), nil
		}
		return nil, nil

	case util.TimeRef:
		return db.ValueAtTime(ref, t)

	default:
		return nil, nil
	}
}

// ValueAtTime returns the value at a specific time by walking commit history
func (db *database) ValueAtTime(datasetID string, t time.Time) (types.Value, error) {
	datasets := db.Datasets()

	var targetRef hash.Hash
	var found bool

	datasets.IterAll(func(k, v types.Value) {
		if datasetID != "" && string(k.(types.String)) != datasetID {
			return
		}
		ref := v.(types.Ref)
		targetRef = ref.TargetHash()
		found = true
	})

	if !found {
		return nil, nil
	}

	// Walk back through commits to find the one at or before the target time
	visited := hash.HashSet{}
	currentHash := targetRef

	for !currentHash.IsEmpty() {
		if visited.Has(currentHash) {
			break
		}
		visited.Insert(currentHash)

		commit := db.ReadValue(currentHash)
		if commit == nil {
			break
		}

		commitStruct := commit.(types.Struct)

		// Check the date of this commit
		meta := commitStruct.Get(MetaField)
		if metaStruct, ok := meta.(types.Struct); ok {
			if dateField, ok := metaStruct.MaybeGet("date"); ok {
				if dateStr, ok := dateField.(types.String); ok {
					if commitTime, err := time.Parse(time.RFC3339, string(dateStr)); err == nil {
						if !commitTime.After(t) {
							// Found the commit at or before target time
							return db.ReadValue(currentHash), nil
						}
					}
				}
			}
		}

		// Move to parent
		parents := commitStruct.Get(ParentsField).(types.Set)
		var parentHash hash.Hash
		parents.IterAll(func(p types.Value) {
			if parentHash.IsEmpty() {
				parentHash = p.(types.Ref).TargetHash()
			}
		})
		currentHash = parentHash
	}

	// Return the earliest commit if we went past the target time
	return db.ReadValue(currentHash), nil
}

// CommitAt returns the commit hash at a specific time
func (db *database) CommitAt(datasetID string, t time.Time) (hash.Hash, error) {
	datasets := db.Datasets()

	var targetRef hash.Hash
	var found bool

	datasets.IterAll(func(k, v types.Value) {
		if datasetID != "" && string(k.(types.String)) != datasetID {
			return
		}
		ref := v.(types.Ref)
		targetRef = ref.TargetHash()
		found = true
	})

	if !found {
		return hash.Hash{}, nil
	}

	// Walk back through commits to find the one at or before the target time
	visited := hash.HashSet{}
	currentHash := targetRef

	for !currentHash.IsEmpty() {
		if visited.Has(currentHash) {
			break
		}
		visited.Insert(currentHash)

		commit := db.ReadValue(currentHash)
		if commit == nil {
			break
		}

		commitStruct := commit.(types.Struct)

		// Check the date of this commit
		meta := commitStruct.Get(MetaField)
		if metaStruct, ok := meta.(types.Struct); ok {
			if dateField, ok := metaStruct.MaybeGet("date"); ok {
				if dateStr, ok := dateField.(types.String); ok {
					if commitTime, err := time.Parse(time.RFC3339, string(dateStr)); err == nil {
						if !commitTime.After(t) {
							return currentHash, nil
						}
					}
				}
			}
		}

		// Move to parent
		parents := commitStruct.Get(ParentsField).(types.Set)
		var parentHash hash.Hash
		parents.IterAll(func(p types.Value) {
			if parentHash.IsEmpty() {
				parentHash = p.(types.Ref).TargetHash()
			}
		})
		currentHash = parentHash
	}

	return currentHash, nil
}
