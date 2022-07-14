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
	ForecastGridData string
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
	GeneratedAt string
	Periods     []ForecastPeriod
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
	Id         string
	Properties FeatureProperties
}

type FeatureCollection struct {
	Features []Feature
}

type ForecastGrid struct {
	Properties GridProperties
}

type GridProperties struct {
	RelativeHumidity GridHumidity
}

type GridHumidity struct {
	Uom    string
	Values []GridHumidityValue
}

type GridHumidityValue struct {
	ValidTime string
	Value     int
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
	Alerts         *FeatureCollection
}

var pCache *cache.Cache
var fCache *cache.Cache
var fhCache *cache.Cache
var fgCache *cache.Cache
var aCache *cache.Cache
var baseUrl *url.URL

func init() {
	var err error

	pCache = cache.NewCache()
	fCache = cache.NewCache()
	fhCache = cache.NewCache()
	fgCache = cache.NewCache()
	aCache = cache.NewCache()

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
	p, nwsErr := Points(lat, lng)
	if nwsErr != nil {
		return nil, nwsErr
	}

	f, nwsErr := GetForecast(p)
	if nwsErr != nil {
		return nil, nwsErr
	}

	fh, nwsErr := GetForecastHourly(p)
	if nwsErr != nil {
		return nil, nwsErr
	}

	a, nwsErr := GetAlerts(lat, lng)
	if nwsErr != nil {
		return nil, nwsErr
	}

	return &ForecastBundle{
		Point:          p,
		Forecast:       f,
		ForecastHourly: fh,
		Alerts:         a,
	}, nil
}

// NWS `/points` endpoint.
func Points(lat, lng float32) (*Point, *Error) {
	var nwsErr *Error
	var expires time.Time
	var body []byte

	ll := fmt.Sprintf("%.4f,%.4f", lat, lng)
	if body = pCache.Get(ll); len(body) == 0 {
		url := fmt.Sprintf("https://api.weather.gov/points/%s", ll)
		body, expires, nwsErr = get(url)
		if nwsErr != nil {
			return nil, nwsErr
		}
		// Cache it.
		pCache.Set(ll, body, expires)
	}

	// Unmarshal.
	point := new(Point)
	err := json.Unmarshal(body, point)
	if err != nil {
		return nil, &Error{
			Title:  "unable json unmarshal",
			Type:   "points-json-error",
			Status: 500,
			Detail: err.Error(),
		}
	}
	if point.Properties.Forecast == "" {
		return nil, &Error{
			Title:  "forecast empty",
			Type:   "points-forecast-error",
			Status: 500,
			Detail: "forecast is empty",
		}
	}
	if point.Properties.ForecastHourly == "" {
		return nil, &Error{
			Title:  "forecast hourly empty",
			Type:   "points-forecast-error",
			Status: 500,
			Detail: "forecast  hourly is empty",
		}
	}
	return point, nil
}

// NWS forecast endpoint.
func GetForecast(point *Point) (*Forecast, *Error) {
	var nwsErr *Error
	var expires time.Time
	var body []byte

	if point == nil {
		return nil, &Error{
			Title:  "point is nil",
			Type:   "forecast-points-invalid",
			Status: 500,
			Detail: "point is nil",
		}
	}
	if len(point.Properties.Forecast) == 0 {
		return nil, &Error{
			Title:  "forecast link is empty",
			Type:   "forecast-link-invalid",
			Status: 500,
			Detail: "forecast link is empty",
		}
	}

	if body = fCache.Get(point.Properties.Forecast); len(body) == 0 {
		// Get the forecast
		body, expires, nwsErr = get(point.Properties.Forecast)
		if nwsErr != nil {
			return nil, nwsErr
		}
		// Cache it.
		fCache.Set(point.Properties.Forecast, body, expires)
	}

	// Unmarshal.
	forecast := new(Forecast)
	err := json.Unmarshal(body, forecast)
	if err != nil {
		return nil, &Error{
			Title:  "forecast json unmarshal failed",
			Type:   "forecast-json-error",
			Status: 500,
			Detail: "forecast json unmarshal failed",
		}
	}
	if len(forecast.Properties.Periods) == 0 {
		return nil, &Error{
			Title:  "forecast has no periods",
			Type:   "forecast-periods-empty",
			Status: 500,
			Detail: "forecast has no periods",
		}
	}
	return forecast, nil
}

// NWS forecast hourly endpoint.
func GetForecastHourly(point *Point) (*Forecast, *Error) {
	var nwsErr *Error
	var expires time.Time
	body := []byte{}

	if point == nil {
		return nil, &Error{
			Title:  "point is nil",
			Type:   "forecast-hourly--points-invalid",
			Status: 500,
			Detail: "point is nil",
		}
	}
	if len(point.Properties.ForecastHourly) == 0 {
		return nil, &Error{
			Title:  "forecast hourly link is empty",
			Type:   "forecast-hourly-link-invalid",
			Status: 500,
			Detail: "forecast hourly link is empty",
		}
	}

	if body = fhCache.Get(point.Properties.ForecastHourly); len(body) == 0 {
		// Get the hourly forecast.
		body, expires, nwsErr = get(point.Properties.ForecastHourly)
		if nwsErr != nil {
			return nil, nwsErr
		}
		// Cache it.
		fhCache.Set(point.Properties.ForecastHourly, body, expires)
	}

	// Unmarshal.
	forecast := new(Forecast)
	err := json.Unmarshal(body, forecast)
	if err != nil {
		return nil, &Error{
			Title:  "forecast hourly json unmarshal failed",
			Type:   "forecast-hourly-json-error",
			Status: 500,
			Detail: "forecast hourly json unmarshal failed",
		}
	}
	if len(forecast.Properties.Periods) == 0 {
		return nil, &Error{
			Title:  "forecast hourly has no periods",
			Type:   "forecast-hourly-periods-empty",
			Status: 500,
			Detail: "forecast hourly has no periods",
		}
	}

	// Check for staleness.
	genAt, tErr := time.Parse(time.RFC3339, forecast.Properties.GeneratedAt)
	if tErr != nil {
		return nil, &Error{
			Title:  "forecast hourly time parsing error",
			Type:   "forecast-hourly-time-parse-error",
			Status: 500,
			Detail: "forecast hourly generated time parsing failed",
		}
	}
	if time.Since(genAt).Seconds() > 86400 {
		fhCache.Set(point.Properties.ForecastHourly,
			[]byte{}, time.Now()) // Invalidate cache.
		return nil, &Error{
			Title:  "forecast hourly is stale",
			Type:   "forecast-hourly-stale-data",
			Status: 500,
			Detail: fmt.Sprintf("stale data from weather.gov from %v",
				forecast.Properties.GeneratedAt),
		}
	}
	return forecast, nil
}

// NWS forecast grid data endpoint.
func GetForecastGridData(point *Point) (*ForecastGrid, *Error) {
	var nwsErr *Error
	var expires time.Time
	var body []byte

	if point == nil {
		return nil, &Error{
			Title:  "point is nil",
			Type:   "griddata-points-invalid",
			Status: 500,
			Detail: "point is nil",
		}
	}
	if len(point.Properties.ForecastGridData) == 0 {
		return nil, &Error{
			Title:  "forecast grid data link is empty",
			Type:   "forecast-griddata-link-invalid",
			Status: 500,
			Detail: "forecast grid data link is empty",
		}
	}

	if body = fgCache.Get(point.Properties.ForecastGridData); len(body) == 0 {
		// Get the forecast grid data
		body, expires, nwsErr = get(point.Properties.ForecastGridData)
		if nwsErr != nil {
			return nil, nwsErr
		}
		// Cache it.
		fgCache.Set(point.Properties.ForecastGridData, body, expires)
	}

	// Unmarshal.
	grid := new(ForecastGrid)
	err := json.Unmarshal(body, grid)
	if err != nil {
		return nil, &Error{
			Title:  "forecast grid data json unmarshal failed",
			Type:   "forecast-griddata-json-error",
			Status: 500,
			Detail: "forecast grid data json unmarshal failed",
		}
	}
	return grid, nil
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

	// Point: lat,lng.
	ll := fmt.Sprintf("%.4f,%.4f", lat, lng)

	// Build query.
	q := url.Values{}
	q.Add("status", "actual")
	q.Add("message_type", "alert")
	q.Add("point", ll)
	q.Add("urgency", "Immediate,Expected")
	q.Add("certainty", "Observed,Likely,Possible")
	u.RawQuery = q.Encode()

	// Hit it.
	var expires time.Time
	var body []byte
	if body = aCache.Get(ll); len(body) == 0 {
		body, expires, err = get(u.String())
		if err != nil {
			return
		}
		// Cache it.
		aCache.Set(ll, body, expires)
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
