// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// ISO 8601 utility functions
package time

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ISO 8601 Duration regex for matching duration in P18DT3H4M60S format.
var durationRegex = regexp.MustCompile(`P(([0-9]{1,})?D)?T(([0-9]{1,2})?H)?(([0-9]{1,2})?M)?(([0-9]{1,2}?)S)?`)

// Converts ISO 8601 duration[1] to time.Duration
//
// Recognizes durations in this format: P18DT3H4M60S
//
// [1]: https://en.wikipedia.org/wiki/ISO_8601#Durations
func Duration(duration string) (time.Duration, error) {
	m := durationRegex.FindStringSubmatch(duration)
	if m == nil || len(m) == 0 {
		return 0, fmt.Errorf("duration invalid: %v", duration)
	}
	days, err := strconv.Atoi(m[2])
	if err != nil {
		days = 0
	}
	hours, err := strconv.Atoi(m[4])
	if err != nil {
		hours = 0
	}
	mins, err := strconv.Atoi(m[6])
	if err != nil {
		mins = 0
	}
	secs, err := strconv.Atoi(m[8])
	if err != nil {
		secs = 0
	}

	// Add 'em all together.
	secs += days * 86400
	secs += hours * 3600
	secs += mins * 60

	// Convert seconds to time.Duration.
	d, err := time.ParseDuration(fmt.Sprintf("%ds", secs))

	return d, nil
}

// Checks if the given ISO 8601 time is within the current time
// period.
//
// `t` must be in 2022-08-07T02:00:00+00:00/PT1H format
//
// Returns true if the time `t` is the current time period; false
// otherwise.
func IsCurrent(t string) (bool, error) {
	parts := strings.Split(t, "/")
	if len(parts) != 2 {
		return false, fmt.Errorf("time invalid")
	}

	// Parse time `t` into time intervals t1 and t2.
	t1, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return false, fmt.Errorf("time invalid: %s", err)
	}
	d, err := Duration(parts[1])
	if err != nil {
		return false, fmt.Errorf("time invalid: %s", err)
	}
	t2 := t1.Add(d)

	// Time `t` is in the current time period if current time is
	// within the interval t1 and t2.
	now := time.Now()
	if t1.Before(now) && now.Before(t2) {
		return true, nil
	}
	if t1.Equal(now) || t2.Equal(now) {
		return true, nil
	}
	return false, nil
}
