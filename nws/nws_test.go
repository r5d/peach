// Copyright © 2022 siddharth <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package nws

import "testing"

func TestPoints(t *testing.T) {
	// Test valid lat,lng.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("points: %v", err)
	}
	if np.Properties.ForecastLink != "https://api.weather.gov/gridpoints/CLE/33,42/forecast" {
		t.Errorf("points: forcecast link: '%v'", np.Properties.ForecastLink)
	}
	if np.Properties.ForecastHourlyLink != "https://api.weather.gov/gridpoints/CLE/33,42/forecast/hourly" {
		t.Errorf("points: forcecast link: '%v'", np.Properties.ForecastHourlyLink)
	}

	// Test invalid lat,lng
	np, err = Points(115.0, -83.177)
	if err == nil {
		t.Errorf("points: %v", np)
	}
}
