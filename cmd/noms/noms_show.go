// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/attic-labs/kingpin"
	"github.com/attic-labs/noms/cmd/util"
	"github.com/attic-labs/noms/go/config"
	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	nomsutil "github.com/attic-labs/noms/go/util"
	"github.com/attic-labs/noms/go/util/datetime"
	"github.com/attic-labs/noms/go/util/outputpager"
)

func nomsShow(noms *kingpin.Application) (*kingpin.CmdClause, util.KingpinHandler) {
	cmd := noms.Command("show", "Print Noms values.")
	showRaw := cmd.Flag("raw", "dump the value in binary format").Bool()
	showStats := cmd.Flag("stats", "report statics related to the value").Bool()
	tzName := cmd.Flag("tz", "display formatted date comments in specified timezone, must be: local or utc").Default("local").String()
	asOf := cmd.Flag("as-of", "show value at a specific time (e.g., '1 hour ago', 'yesterday', '2024-01-15')").String()
	path := cmd.Arg("path", "value to display - see Spelling Values at https://github.com/attic-labs/noms/blob/master/doc/spelling.md").Required().String()

	return cmd, func(_ string) int {
		cfg := config.NewResolver()
		database, value, err := cfg.GetPath(*path)
		d.CheckErrorNoUsage(err)
		defer database.Close()

		// Handle time travel
		if *asOf != "" {
			targetTime, err := nomsutil.ParseTime(*asOf)
			if err == nil {
				// Try using the ValueAt method with the time
				db, dbErr := cfg.GetDatabase("")
				if dbErr == nil {
					defer db.Close()
					// Use reflection to call ValueAtTime
					val, err := tryValueAtTime(db, *path, targetTime)
					if err == nil && val != nil {
						value = val
					}
				}
			}
		}

		if value == nil {
			fmt.Fprintf(os.Stderr, "Value not found: %s\n", *path)
			return 0
		}

		if *showRaw && *showStats {
			fmt.Fprintln(os.Stderr, "--raw and --stats are mutually exclusive")
			return 0
		}

		if *showRaw {
			ch := types.EncodeValue(value)
			buf := bytes.NewBuffer(ch.Data())
			_, err = io.Copy(os.Stdout, buf)
			d.CheckError(err)
			return 0
		}

		if *showStats {
			types.WriteValueStats(os.Stdout, value, database)
			return 0
		}

		tz, _ := locationFromTimezoneArg(*tzName, nil)
		datetime.RegisterHRSCommenter(tz)

		pgr := outputpager.Start()
		defer pgr.Stop()

		types.WriteEncodedValue(pgr.Writer, value)
		fmt.Fprintln(pgr.Writer)
		return 0
	}
}

func tryValueAtTime(db datas.Database, path string, t time.Time) (types.Value, error) {
	dbVal := reflect.ValueOf(db)
	if dbVal.Kind() == reflect.Ptr {
		dbVal = dbVal.Elem()
	}

	method := dbVal.MethodByName("ValueAtTime")
	if !method.IsValid() {
		return nil, fmt.Errorf("ValueAtTime not available")
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(path), reflect.ValueOf(t)})
	if len(results) < 2 {
		return nil, nil
	}

	err, _ := results[1].Interface().(error)
	if err != nil {
		return nil, err
	}

	val, _ := results[0].Interface().(types.Value)
	return val, nil
}
