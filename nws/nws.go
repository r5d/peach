// Copyright © 2022 siddharth <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Functions for accessing the National Weather Service API.
// TODO: remove NWS prefix from all types.
package nws

import (
	"encoding/json"
	"fmt"
	"io"

	"ricketyspace.net/peach/client"
)

type PointLocationProperties struct {
	City  string
	State string
}

type PointLocation struct {
	Properties PointLocationProperties
}

type NWSPointProperties struct {
	GridId           string
	GridX            int
	GridY            int
	Forecast         string
	ForecastHourly   string
	RelativeLocation PointLocation
}

type NWSPoint struct {
	Properties NWSPointProperties
}

type NWSForecastPeriod struct {
	Number           int
	Name             string
	StartTime        string
	EndTime          string
	IsDayTime        bool
	Temperature      int
	TemperatureUnit  string
	TemperatureTrend string
	WindSpeed        string
	WindDirection    string
	ShortForecast    string
	DetailedForecast string
}

type NWSForecastProperties struct {
	Periods []NWSForecastPeriod
}

type NWSForecast struct {
	Properties NWSForecastProperties
}

type NWSError struct {
	Title  string
	Type   string
	Status int
	Detail string
}

func (e NWSError) Error() string {
	return fmt.Sprintf("%d: %s: %s", e.Status, e.Type, e.Detail)
}

// NWS `/points` endpoint.
//
// TODO: return NWSError instead of error
func Points(lat, lng float32) (*NWSPoint, error) {
	url := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", lat, lng)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("points: http get: %v", err)
	}

	// Parse response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("points: body: %v", err)
	}

	// Check if the request failed.
	if resp.StatusCode != 200 {
		perr := new(NWSError)
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
		return nil, fmt.Errorf("points: decode: %v", err)
	}
	if point.Properties.Forecast == "" {
		return nil, fmt.Errorf("points: forecast empty")
	}
	if point.Properties.ForecastHourly == "" {
		return nil, fmt.Errorf("points: forecasthourly empty")
	}
	return point, nil
}

// NWS forecast endpoint.
//
// TODO: return NWSError instead of error.
func Forecast(point *NWSPoint) (*NWSForecast, error) {
	if point == nil {
		return nil, fmt.Errorf("forecast: point nil")
	}
	if len(point.Properties.Forecast) == 0 {
		return nil, fmt.Errorf("forecast: link empty")
	}

	// Get the forecast
	resp, err := client.Get(point.Properties.Forecast)
	if err != nil {
		return nil, fmt.Errorf("forecast: get: %v", err)
	}

	// Parse response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("forecast: body: %v", err)
	}

	// Check if the request failed.
	if resp.StatusCode != 200 {
		perr := new(NWSError)
		err := json.Unmarshal(body, perr)
		if err != nil {
			return nil, fmt.Errorf("forecast: json: %v", err)
		}
		return nil, fmt.Errorf("forecast: %v", perr)
	}

	// Unmarshal.
	forecast := new(NWSForecast)
	err = json.Unmarshal(body, forecast)
	if err != nil {
		return nil, fmt.Errorf("forecast: decode: %v", err)
	}
	if len(forecast.Properties.Periods) == 0 {
		return nil, fmt.Errorf("forecast: periods empty")
	}
	return forecast, nil
}

// NWS forecast hourly endpoint.
//
// TODO: return NWSError instead of error
func ForecastHourly(point *NWSPoint) (*NWSForecast, error) {
	if point == nil {
		return nil, fmt.Errorf("forecast hourly: point nil")
	}
	if len(point.Properties.ForecastHourly) == 0 {
		return nil, fmt.Errorf("forecast hourly: link empty")
	}

	// Get the forecast
	resp, err := client.Get(point.Properties.ForecastHourly)
	if err != nil {
		return nil, fmt.Errorf("forecast hourly: get: %v", err)
	}

	// Parse response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("forecast hourly: body: %v", err)
	}

	// Check if the request failed.
	if resp.StatusCode != 200 {
		perr := new(NWSError)
		err := json.Unmarshal(body, perr)
		if err != nil {
			return nil, fmt.Errorf("forecast hourly: json: %v", err)
		}
		return nil, fmt.Errorf("forecast: %v", perr)
	}

	// Unmarshal.
	forecast := new(NWSForecast)
	err = json.Unmarshal(body, forecast)
	if err != nil {
		return nil, fmt.Errorf("forecast hourly: decode: %v", err)
	}
	if len(forecast.Properties.Periods) == 0 {
		return nil, fmt.Errorf("forecast hourly: periods empty")
	}
	return forecast, nil
}
