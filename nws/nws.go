// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Functions for accessing the National Weather Service API.
package nws

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"ricketyspace.net/peach/cache"
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

var pCache *cache.Cache
var fCache *cache.Cache
var fhCache *cache.Cache

func init() {
	pCache = cache.NewCache()
	fCache = cache.NewCache()
	fhCache = cache.NewCache()
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s: %s", e.Status, e.Type, e.Detail)
}

func CacheWeather(lat, lng float32) {
	p, err := Points(lat, lng)
	if err != nil {
		return
	}
	GetForecast(p)
	GetForecastHourly(p)
}

// NWS `/points` endpoint.
//
// TODO: return Error instead of error
func Points(lat, lng float32) (*Point, error) {
	var nwsErr *Error
	var expires time.Time
	var body []byte

	ll := fmt.Sprintf("%.4f,%.4f", lat, lng)
	if body = pCache.Get(ll); len(body) == 0 {
		url := fmt.Sprintf("https://api.weather.gov/points/%s", ll)
		body, expires, nwsErr = get(url)
		if nwsErr != nil {
			return nil, fmt.Errorf("points: %v", nwsErr)
		}
		// Cache it.
		pCache.Set(ll, body, expires)
	}

	// Unmarshal.
	point := new(Point)
	err := json.Unmarshal(body, point)
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
	var nwsErr *Error
	var expires time.Time
	var body []byte

	if point == nil {
		return nil, fmt.Errorf("forecast: point nil")
	}
	if len(point.Properties.Forecast) == 0 {
		return nil, fmt.Errorf("forecast: link empty")
	}

	if body = fCache.Get(point.Properties.Forecast); len(body) == 0 {
		// Get the forecast
		body, expires, nwsErr = get(point.Properties.Forecast)
		if nwsErr != nil {
			return nil, fmt.Errorf("forecast: %v", nwsErr)
		}
		// Cache it.
		fCache.Set(point.Properties.Forecast, body, expires)
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
	var nwsErr *Error
	var expires time.Time
	body := []byte{}

	if point == nil {
		return nil, fmt.Errorf("forecast hourly: point nil")
	}
	if len(point.Properties.ForecastHourly) == 0 {
		return nil, fmt.Errorf("forecast hourly: link empty")
	}

	if body = fhCache.Get(point.Properties.ForecastHourly); len(body) == 0 {
		// Get the hourly forecast.
		body, expires, nwsErr = get(point.Properties.ForecastHourly)
		if nwsErr != nil {
			return nil, fmt.Errorf("forecast hourly: %v", nwsErr)
		}
		// Cache it.
		fhCache.Set(point.Properties.ForecastHourly, body, expires)
	}

	// Unmarshal.
	forecast := new(Forecast)
	err := json.Unmarshal(body, forecast)
	if err != nil {
		return nil, fmt.Errorf("forecast hourly: decode: %v", err)
	}
	if len(forecast.Properties.Periods) == 0 {
		return nil, fmt.Errorf("forecast hourly: periods empty")
	}
	return forecast, nil
}

// HTTP GET a NWS endpoint.
func get(url string) ([]byte, time.Time, *Error) {
	// Default response expiration time
	expires := time.Now()

	tries := 5
	retryDelay := 100 * time.Millisecond
	for {
		resp, err := client.Get(url)
		if err != nil {
			return nil, expires, &Error{
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
			return nil, expires, &Error{
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
				return nil, expires, &Error{
					Title:  fmt.Sprintf("json decode: %v", url),
					Type:   "json-decode",
					Status: 500,
					Detail: err.Error(),
				}
			}
			return nil, expires, &nwsErr
		}

		// Parse expiration time of the response.
		expiresHeader := resp.Header.Get("expires")
		if len(expiresHeader) < 1 {
			return nil, expires, &Error{
				Title:  "expiration header empty",
				Type:   "expiration-header",
				Status: 500,
				Detail: "response expiration header is empty",
			}
		}
		expires, err := time.Parse(time.RFC1123, expiresHeader)
		if err != nil {
			return nil, expires, &Error{
				Title:  "expiration header could not be parsed",
				Type:   "expiration-header-parse-failed",
				Status: 500,
				Detail: err.Error(),
			}
		}
		// Response OK.
		return body, expires, nil
	}
}
