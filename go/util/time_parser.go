// Copyright 2025 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	relativePatterns = []*regexp.Regexp{
		regexp.MustCompile(`^(\d+)\s*second(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*minute(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*hour(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*day(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*week(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*month(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*year(s)?\s+ago$`),
		regexp.MustCompile(`^(\d+)\s*hour(s)?\s+ago$`),
	}

	dateFormats = []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02/2006",
		"02-01-2006",
	}
)

// ParseTime parses natural language time expressions
// Examples:
//   - "1 hour ago"
//   - "2 days ago"
//   - "yesterday"
//   - "2024-01-15"
//   - "2024-01-15T10:30:00Z"
//   - "1 week ago"
//   - "3 months ago"
func ParseTime(expr string) (time.Time, error) {
	expr = strings.TrimSpace(expr)
	lower := strings.ToLower(expr)

	// Handle special cases
	if lower == "now" {
		return time.Now(), nil
	}

	if lower == "yesterday" {
		return time.Now().AddDate(0, 0, -1), nil
	}

	if lower == "today" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}

	if lower == "tomorrow" {
		return time.Now().AddDate(0, 0, 1), nil
	}

	if lower == "last week" {
		return time.Now().AddDate(0, 0, -7), nil
	}

	if lower == "last month" {
		return time.Now().AddDate(0, -1, 0), nil
	}

	if lower == "last year" {
		return time.Now().AddDate(-1, 0, 0), nil
	}

	// Try relative patterns
	for _, pattern := range relativePatterns {
		match := pattern.FindStringSubmatch(lower)
		if match != nil {
			var duration time.Duration
			var value int
			fmt.Sscanf(match[1], "%d", &value)

			switch {
			case strings.Contains(pattern.String(), "second"):
				duration = time.Duration(value) * time.Second
			case strings.Contains(pattern.String(), "minute"):
				duration = time.Duration(value) * time.Minute
			case strings.Contains(pattern.String(), "hour"):
				duration = time.Duration(value) * time.Hour
			case strings.Contains(pattern.String(), "day"):
				duration = time.Duration(value) * 24 * time.Hour
			case strings.Contains(pattern.String(), "week"):
				duration = time.Duration(value) * 7 * 24 * time.Hour
			case strings.Contains(pattern.String(), "month"):
				duration = time.Duration(value) * 30 * 24 * time.Hour
			case strings.Contains(pattern.String(), "year"):
				duration = time.Duration(value) * 365 * 24 * time.Hour
			}

			return time.Now().Add(-duration), nil
		}
	}

	// Try date formats
	for _, format := range dateFormats {
		if t, err := time.Parse(format, expr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time expression: %s", expr)
}

// ParseRange parses time ranges like "2024-01-01..2024-02-01" or "1 week ago..now"
func ParseRange(expr string) (start, end time.Time, err error) {
	parts := strings.Split(expr, "..")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid range format: %s", expr)
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	if startStr == "" {
		start = time.Time{}
	} else {
		start, err = ParseTime(startStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if endStr == "" || endStr == "now" {
		end = time.Now()
	} else {
		end, err = ParseTime(endStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return start, end, nil
}

// ParseRef parses a reference string that could be a hash, branch, or time expression
func ParseRef(ref string) (RefType, time.Time, error) {
	// Try as time first
	if t, err := ParseTime(ref); err == nil {
		return TimeRef, t, nil
	}

	// Check if it looks like a hash (40 hex chars)
	if matched, _ := regexp.MatchString(`^[a-f0-9]{40}$`, ref); matched {
		return HashRef, time.Time{}, nil
	}

	// Short hash
	if matched, _ := regexp.MatchString(`^[a-f0-9]{4,39}$`, ref); matched {
		return ShortHashRef, time.Time{}, nil
	}

	// Assume it's a branch name
	return BranchRef, time.Time{}, nil
}

// RefType represents the type of reference
type RefType int

const (
	UnknownRef RefType = iota
	HashRef
	ShortHashRef
	BranchRef
	TimeRef
)

// FormatTime formats a time in a human-readable way
func FormatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, plural(minutes))
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, plural(days))
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / (7 * 24))
		return fmt.Sprintf("%d week%s ago", weeks, plural(weeks))
	default:
		return t.Format("2006-01-02")
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
