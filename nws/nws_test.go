// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package nws

import "testing"

func TestPoints(t *testing.T) {
	// Test valid lat,lng.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("points: %v", err)
		return
	}
	if np.Properties.Forecast != "https://api.weather.gov/gridpoints/CLE/33,42/forecast" {
		t.Errorf("points: forcecast link: '%v'", np.Properties.Forecast)
	}
	if np.Properties.ForecastHourly != "https://api.weather.gov/gridpoints/CLE/33,42/forecast/hourly" {
		t.Errorf("points: forcecast link: '%v'", np.Properties.ForecastHourly)
	}
	if np.Properties.GridId != "CLE" {
		t.Errorf("points: gridid: %v", np.Properties.GridId)
	}
	if np.Properties.GridX != 33 {
		t.Errorf("points: gridx: %v", np.Properties.GridX)
	}
	if np.Properties.GridY != 42 {
		t.Errorf("points: gridy: %v", np.Properties.GridY)
	}

	if np.Properties.RelativeLocation.Properties.City != "Tiffin" {
		t.Errorf("points: location: city: %v", np.Properties)
	}
	if np.Properties.RelativeLocation.Properties.State != "OH" {
		t.Errorf("points: location: state: %v", np.Properties)
	}

	// Test invalid lat,lng
	np, err = Points(115.0, -83.177)
	if err == nil {
		t.Errorf("points: %v", np)
	}
}

func TestForecast(t *testing.T) {
	// Get point.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Get forecast.
	fc, err := Forecast(np)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Verify periods.
	for i, period := range fc.Properties.Periods {
		if period.Number < 1 {
			t.Errorf("period: %d: number invalid: %v", i, period.Number)
		}
		if len(period.Name) < 1 {
			t.Errorf("period: %d: name invalid: %v", i, period.Name)
		}
		if len(period.StartTime) < 1 {
			t.Errorf("period: %d: start time invalid: %v", i,
				period.StartTime)
		}
		if len(period.EndTime) < 1 {
			t.Errorf("period: %d: end time invalid: %v", i,
				period.EndTime)
		}
		if len(period.TemperatureUnit) < 1 {
			t.Errorf("period: %d: temperature unit invalid: %v",
				i, period.TemperatureUnit)
		}
		if len(period.WindSpeed) < 1 {
			t.Errorf("period: %d: wind speed invalid: %v",
				i, period.WindSpeed)
		}
		if len(period.WindDirection) < 1 {
			t.Errorf("period: %d: wind direction invalid: %v",
				i, period.WindDirection)
		}
		if len(period.ShortForecast) < 1 {
			t.Errorf("period: %d: short forecast invalid: %v",
				i, period.ShortForecast)
		}
		if len(period.DetailedForecast) < 1 {
			t.Errorf("period: %d: detailed forecast invalid: %v",
				i, period.DetailedForecast)
		}
	}
}

func TestForecastHourly(t *testing.T) {
	// Get point.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Get forecast hourly.
	fc, err := ForecastHourly(np)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Verify periods.
	for i, period := range fc.Properties.Periods {
		if period.Number < 1 {
			t.Errorf("period: %d: number invalid: %v", i, period.Number)
		}
		if len(period.StartTime) < 1 {
			t.Errorf("period: %d: start time invalid: %v", i,
				period.StartTime)
		}
		if len(period.EndTime) < 1 {
			t.Errorf("period: %d: end time invalid: %v", i,
				period.EndTime)
		}
		if len(period.TemperatureUnit) < 1 {
			t.Errorf("period: %d: temperature unit invalid: %v",
				i, period.TemperatureUnit)
		}
		if len(period.WindSpeed) < 1 {
			t.Errorf("period: %d: wind speed invalid: %v",
				i, period.WindSpeed)
		}
		if len(period.WindDirection) < 1 {
			t.Errorf("period: %d: wind direction invalid: %v",
				i, period.WindDirection)
		}
		if len(period.ShortForecast) < 1 {
			t.Errorf("period: %d: short forecast invalid: %v",
				i, period.ShortForecast)
		}
	}
}
