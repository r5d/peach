// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package nws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestGetForecast(t *testing.T) {
	// Get point.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Get forecast.
	fc, err := GetForecast(np)
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

func TestGetForecastHourly(t *testing.T) {
	// Get point.
	np, err := Points(41.115, -83.177)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// Get forecast hourly.
	fc, err := GetForecastHourly(np)
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

func TestNWSGetWrapper(t *testing.T) {
	// Initialize test NWS server.
	fails := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fails > 0 {
			fails -= 1
			http.Error(w, `{"type":"urn:noaa:nws:api:UnexpectedProblem","title":"Unexpected Problem","status":500,"detail":"An unexpected problem has occurred.","instance":"urn:noaa:nws:api:request:493c3a1d-f87e-407f-ae2c-24483f5aab63","correlationId":"493c3a1d-f87e-407f-ae2c-24483f5aab63","additionalProp1":{}}`, 500)
			return
		}
		// Success.
		fmt.Fprintln(w, `{"@context":[],"properties":{"gridId":"CLE","gridX":82,"gridY":64,"forecast":"https://api.weather.gov/gridpoints/CLE/82,64/forecast","forecastHourly":"https://api.weather.gov/gridpoints/CLE/82,64/forecast/hourly","relativeLocation":{"properties":{"city":"Cleveland","state":"OH"}}}}`)
	}))
	defer ts.Close()

	// Test 1 - Server fails 5 times.
	fails = 5
	_, err := get(ts.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
		return
	}

	// Test 2 - Server fails 6 times.
	fails = 6
	respBody, err := get(ts.URL)
	if err == nil {
		t.Errorf("get did not fail: %s", respBody)
		return
	}
	if err != nil && respBody != nil {
		t.Errorf("body is not nil: %s", respBody)
	}
	if err.Title != "Unexpected Problem" {
		t.Errorf("err title: %s", err.Title)
		return
	}
	if err.Type != "urn:noaa:nws:api:UnexpectedProblem" {
		t.Errorf("err type: %s", err.Type)
		return
	}
	if err.Status != 500 {
		t.Errorf("err status: %d", err.Status)
		return
	}
	if err.Detail != "An unexpected problem has occurred." {
		t.Errorf("err detail: %s", err.Detail)
		return
	}

	// Test 3 - Server fails 1 time.
	fails = 1
	respBody, err = get(ts.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
		return
	}
	if respBody == nil {
		t.Errorf("body: %s", respBody)
		return
	}
	point := new(Point)
	jerr := json.Unmarshal(respBody, point)
	if jerr != nil {
		t.Errorf("points: decode: %v", jerr)
		return
	}
	if point.Properties.Forecast == "" {
		t.Errorf("points: forecast empty")
		return
	}
	if point.Properties.ForecastHourly == "" {
		t.Errorf("points: forecasthourly empty")
		return
	}
	if point.Properties.RelativeLocation.Properties.City == "" {
		t.Errorf("points: city empty")
		return
	}
	if point.Properties.RelativeLocation.Properties.State == "" {
		t.Errorf("points: state empty")
		return
	}
}
