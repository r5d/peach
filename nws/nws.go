// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Functions for accessing the National Weather Service API.
package nws

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"ricketyspace.net/peach/client"
)

type PointLocationProperties struct {
	City  string
	State string
}

type PointLocation struct {
	Properties PointLocationProperties
}

type PointProperties struct {
	GridId           string
	GridX            int
	GridY            int
	Forecast         string
	ForecastHourly   string
	RelativeLocation PointLocation
}

type Point struct {
	Properties PointProperties
}

type ForecastPeriod struct {
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

type ForecastProperties struct {
	Periods []ForecastPeriod
}

type Forecast struct {
	Properties ForecastProperties
}

type Error struct {
	Title  string
	Type   string
	Status int
	Detail string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s: %s", e.Status, e.Type, e.Detail)
}

// NWS `/points` endpoint.
//
// TODO: return Error instead of error
func Points(lat, lng float32) (*Point, error) {
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
		perr := new(Error)
		err := json.Unmarshal(body, perr)
		if err != nil {
			return nil, fmt.Errorf("points: json: %v", err)
		}
		return nil, fmt.Errorf("points: %v", perr)
	}

	// Unmarshal.
	point := new(Point)
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
// TODO: return Error instead of error.
func GetForecast(point *Point) (*Forecast, error) {
	if point == nil {
		return nil, fmt.Errorf("forecast: point nil")
	}
	if len(point.Properties.Forecast) == 0 {
		return nil, fmt.Errorf("forecast: link empty")
	}

	// Get the forecast
	body, nwsErr := get(point.Properties.Forecast)
	if nwsErr != nil {
		return nil, fmt.Errorf("forecast: %v", nwsErr)
	}

	// Unmarshal.
	forecast := new(Forecast)
	err := json.Unmarshal(body, forecast)
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
// TODO: return Error instead of error
func GetForecastHourly(point *Point) (*Forecast, error) {
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
		perr := new(Error)
		err := json.Unmarshal(body, perr)
		if err != nil {
			return nil, fmt.Errorf("forecast hourly: json: %v", err)
		}
		return nil, fmt.Errorf("forecast: %v", perr)
	}

	// Unmarshal.
	forecast := new(Forecast)
	err = json.Unmarshal(body, forecast)
	if err != nil {
		return nil, fmt.Errorf("forecast hourly: decode: %v", err)
	}
	if len(forecast.Properties.Periods) == 0 {
		return nil, fmt.Errorf("forecast hourly: periods empty")
	}
	return forecast, nil
}

// HTTP GET a NWS endpoint.
func get(url string) ([]byte, *Error) {
	tries := 5
	retryDelay := 100 * time.Millisecond
	for {
		resp, err := client.Get(url)
		if err != nil {
			return nil, &Error{
				Title:  fmt.Sprintf("http get failed: %v", url),
				Type:   "http-get",
				Status: 500,
				Detail: err.Error(),
			}
		}
		if tries > 0 && resp.StatusCode != 200 {
			tries -= 1

			// Wait before re-try.
			time.Sleep(retryDelay)

			retryDelay *= 2 // Exponential back-off delay.
			continue        // Re-try
		}

		// Parse response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &Error{
				Title:  fmt.Sprintf("parsing body: %v", url),
				Type:   "response-body",
				Status: 500,
				Detail: err.Error(),
			}
		}

		// Check if the request failed.
		if resp.StatusCode != 200 {
			nwsErr := Error{}
			err := json.Unmarshal(body, &nwsErr)
			if err != nil {
				return nil, &Error{
					Title:  fmt.Sprintf("json decode: %v", url),
					Type:   "json-decode",
					Status: 500,
					Detail: err.Error(),
				}
			}
			return nil, &nwsErr
		}
		// Response OK.
		return body, nil
	}
}
