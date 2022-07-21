// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package time

import "testing"

func TestDurationToSeconds(t *testing.T) {
	secs, err := durationToSeconds("PT3H4M60S")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if secs != 11100 {
		t.Errorf("duration in seconds incorrect: %v", secs)
		return
	}

	secs, err = durationToSeconds("PT4M60S")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if secs != 300 {
		t.Errorf("duration in seconds incorrect: %v", secs)
		return
	}

	secs, err = durationToSeconds("PT12H")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if secs != 43200 {
		t.Errorf("duration in seconds incorrect: %v", secs)
		return
	}

	secs, err = durationToSeconds("PT1H")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if secs != 3600 {
		t.Errorf("duration in seconds incorrect: %v", secs)
		return
	}

	secs, err = durationToSeconds("PT2H")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if secs != 7200 {
		t.Errorf("duration in seconds incorrect: %v", secs)
		return
	}

	secs, err = durationToSeconds("PT45M")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if secs != 2700 {
		t.Errorf("duration in seconds incorrect: %v", secs)
		return
	}
}
