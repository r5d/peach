// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Functions for accessing the National Weather Service API.
package nws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
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

type FeatureProperties struct {
	Event       string
	Severity    string
	Description string
	Instruction string
}

type Feature struct {
	Properties FeatureProperties
}

type FeatureCollection struct {
	Features []Feature
}

type Error struct {
	Title  string
	Type   string
	Status int
	Detail string
}

// NWS Forecast bundle.
type ForecastBundle struct {
	Point          *Point
	Forecast       *Forecast
	ForecastHourly *Forecast
}

var pCache *cache.Cache
var fCache *cache.Cache
var fhCache *cache.Cache
var baseUrl *url.URL

func init() {
	var err error

	pCache = cache.NewCache()
	fCache = cache.NewCache()
	fhCache = cache.NewCache()

	// Parse NWS base url.
	baseUrl, err = url.Parse("https://api.weather.gov")
	if err != nil {
		panic(`url parse: nws base url: ` + err.Error())
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s: %s", e.Status, e.Type, e.Detail)
}

// Gets NWS's forecast and hourly forecast.
func GetForecastBundle(lat, lng float32) (*ForecastBundle, *Error) {
	p, err := Points(lat, lng)
	if err != nil {
		return nil, &Error{
			Title:  "unable get points",
			Type:   "points-failed",
			Status: 500,
			Detail: err.Error(),
		}
	}
	f, err := GetForecast(p)
	if err != nil {
		return nil, &Error{
			Title:  "unable get forecast",
			Type:   "forecast-failed",
			Status: 500,
			Detail: err.Error(),
		}
	}
	fh, err := GetForecastHourly(p)
	if err != nil {
		return nil, &Error{
			Title:  "unable get hourly forecast",
			Type:   "forecast-hourly-failed",
			Status: 500,
			Detail: err.Error(),
		}
	}
	return &ForecastBundle{
		Point:          p,
		Forecast:       f,
		ForecastHourly: fh,
	}, nil
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

// NWS active alerts endpoint.
func GetAlerts(lat, lng float32) (fc *FeatureCollection, err *Error) {
	// Alerts endpoint.
	u, uErr := baseUrl.Parse("/alerts/active")
	if uErr != nil {
		err = &Error{
			Title:  "alerts url parsing failed",
			Type:   "url-parse-error",
			Status: 500,
			Detail: uErr.Error(),
		}
		return
	}

	// Build query.
	q := url.Values{}
	q.Add("status", "actual")
	q.Add("message_type", "alert")
	q.Add("point", fmt.Sprintf("%.4f,%.4f", lat, lng))
	q.Add("urgency", "Immediate,Expected")
	q.Add("certainty", "Observed,Likely,Possible")
	u.RawQuery = q.Encode()

	// Hit it.
	body, _, err := get(u.String())
	if err != nil {
		return
	}

	// Unmarshal.
	fc = new(FeatureCollection)
	jErr := json.Unmarshal(body, fc)
	if jErr != nil {
		err = &Error{
			Title:  "feature collection decode failed",
			Type:   "json-decode-error",
			Status: 500,
			Detail: jErr.Error(),
		}
		return
	}
	return
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
