// Copyright © 2022 siddharth <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Functions for accessing the National Weather Service API.
package nws

import (
	"encoding/json"
	"fmt"
	"io"

	"ricketyspace.net/peach/client"
)

type NWSPointProperties struct {
	GridId             string
	GridX              int
	GridY              int
	ForecastLink       string `json:"forecast"`
	ForecastHourlyLink string `json:"forecastHourly"`
}

type NWSPoint struct {
	Properties NWSPointProperties
}

type NWSPointError struct {
	Title  string
	Type   string
	Status int
	Detail string
}

func (e NWSPointError) Error() string {
	return fmt.Sprintf("%d: %s: %s", e.Status, e.Type, e.Detail)
}

// NWS `/points` endpoint.
func Points(lat, lng float32) (*NWSPoint, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", lat, lng)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("points: http get: %v", err)
	}

	// Parse response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getting points: body: %v", err)
	}

	// Check if the request failed.
	if resp.StatusCode != 200 {
		perr := new(NWSPointError)
		err := json.Unmarshal(body, perr)
		if err != nil {
			return nil, fmt.Errorf("points: json: %v", err)
		}
		return nil, fmt.Errorf("points: %v", perr)
	}

	// Unmarshal.
	point := new(NWSPoint)
	err = json.Unmarshal(body, point)
	if err != nil {
		return nil, fmt.Errorf("points: json decode: %v", err)
	}
	if point.Properties.ForecastLink == "" {
		return nil, fmt.Errorf("points: json: %v", err)
	}
	if point.Properties.ForecastHourlyLink == "" {
		return nil, fmt.Errorf("points: json: %v", err)
	}
	return point, nil
}

