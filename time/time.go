// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// ISO 8601 utility functions
package time

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ISO 8601 Duration regex for matching duration in PT3H4M60S format.
var durationRegex = regexp.MustCompile(`PT(([0-9]{0,2})?H)?(([0-9]{0,2})?M)?(([0-9]{0,2}?)S)?`)

// Converts ISO 8601 duration[1] to time.Duration
//
// Recognizes durations in this format: PT3H4M60S
//
// [1]: https://en.wikipedia.org/wiki/ISO_8601#Durations
func Duration(duration string) (time.Duration, error) {
	m := durationRegex.FindStringSubmatch(duration)
	if m == nil || len(m) == 0 {
		return 0, fmt.Errorf("duration invalid: %v", duration)
	}
	hours, err := strconv.Atoi(m[2])
	if err != nil {
		hours = 0
	}
	mins, err := strconv.Atoi(m[4])
	if err != nil {
		mins = 0
	}
	secs, err := strconv.Atoi(m[6])
	if err != nil {
		secs = 0
	}

	// Add 'em all together.
	secs += hours * 3600
	secs += mins * 60

	// Convert seconds to time.Duration.
	d, err := time.ParseDuration(fmt.Sprintf("%ds", secs))

	return d, nil
}
}
