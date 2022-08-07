// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package time

import "testing"

func TestDurationToSeconds(t *testing.T) {
	d, err := Duration("PT3H4M60S")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if d.Seconds() != 11100 {
		t.Errorf("duration in seconds incorrect: %v", d)
		return
	}

	d, err = Duration("PT4M60S")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if d.Seconds() != 300 {
		t.Errorf("duration in seconds incorrect: %v", d)
		return
	}

	d, err = Duration("PT12H")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if d.Seconds() != 43200 {
		t.Errorf("duration in seconds incorrect: %v", d)
		return
	}

	d, err = Duration("PT1H")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if d.Seconds() != 3600 {
		t.Errorf("duration in seconds incorrect: %v", d)
		return
	}

	d, err = Duration("PT2H")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if d.Seconds() != 7200 {
		t.Errorf("duration in seconds incorrect: %v", d)
		return
	}

	d, err = Duration("PT45M")
	if err != nil {
		t.Errorf("failed: %v", err)
		return
	}
	if d.Seconds() != 2700 {
		t.Errorf("duration in seconds incorrect: %v", d)
		return
	}
}
