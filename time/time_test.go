// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package time

import (
	"fmt"
	"testing"
	"time"
)

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

func TestIsCurrent(t *testing.T) {
	yes, err := IsCurrent("2022-08-07T01:00:00+00:00/PT1H")
	if yes || err != nil {
		t.Errorf("iscurrent failed: %v: %v", yes, err)
		return
	}

	h, err := time.ParseDuration("-3600s")
	if err != nil {
		t.Errorf("-3600s parsing duration: %v", err)
		return
	}
	ts := time.Now().Add(h).UTC().Format(time.RFC3339)
	yes, err = IsCurrent(fmt.Sprintf("%s+00:00/PT2H", ts[:len(ts)-1]))
	if !yes || err != nil {
		t.Errorf("iscurrent failed: %v: %v", yes, err)
		return
	}

}
