// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ricketyspace.net/peach/nws"
)

// peach port. defaults to 8151
var peachPort = flag.Int("p", 8151, "Port to run peach on")

// peach listen address. set during init.
var peachAddr = ""

// holds static content.
//go:embed templates static
var peachFS embed.FS

// html templates.
var peachTemplates = template.Must(template.ParseFS(peachFS, "templates/*.tmpl"))

// lat,long regex.
var latLngRegex = regexp.MustCompile(`/(-?[0-9]+\.?[0-9]+?),(-?[0-9]+\.?[0-9]+)`)

type Weather struct {
	Location string
	Now      WeatherNow
	Period   WeatherPeriod
	Timeline WeatherTimeline
}

type WeatherNow struct {
	Temperature     int
	TemperatureUnit string
	Forecast        string
	WindSpeed       string
	WindDirection   string
}

type WeatherPeriod struct {
	Forecast        string
	Hour            int
	Temperature     int
	TemperatureUnit string
}

type WeatherTimeline struct {
	Periods []WeatherPeriod
}

func init() {
	flag.Parse()
	if *peachPort < 80 {
		log.Fatalf("port number is invalid: %v", *peachPort)
	}
	peachAddr = fmt.Sprintf(":%d", *peachPort)
}

func main() {
	// static files handler.
	http.Handle("/static/", http.FileServer(http.FS(peachFS)))

	// default handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/41.115,-83.177", 302)
			return
		}

		m := latLngRegex.FindStringSubmatch(r.URL.Path)
		if len(m) != 3 || m[0] != r.URL.Path {
			http.NotFound(w, r)
			return
		}
		lat, err := strconv.ParseFloat(m[1], 32)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		lng, err := strconv.ParseFloat(m[2], 32)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		showWeather(w, float32(lat), float32(lng))
	})
	log.Fatal(http.ListenAndServe(peachAddr, nil))
}

func showWeather(w http.ResponseWriter, lat, lng float32) {
	point, err := nws.Points(lat, lng)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// get forecast
	f, err := nws.Forecast(point)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fh, err := nws.ForecastHourly(point)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// make weather
	weather, err := NewWeather(point, f, fh)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// render.
	err = peachTemplates.ExecuteTemplate(w, "weather.tmpl", weather)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func NewWeather(point *nws.NWSPoint, f, fh *nws.NWSForecast) (*Weather, error) {
	w := new(Weather)
	w.Location = fmt.Sprintf("%s, %s",
		strings.ToLower(point.Properties.RelativeLocation.Properties.City),
		strings.ToLower(point.Properties.RelativeLocation.Properties.State),
	)
	w.Now = WeatherNow{
		Temperature:     fh.Properties.Periods[0].Temperature,
		TemperatureUnit: fh.Properties.Periods[0].TemperatureUnit,
		Forecast:        fh.Properties.Periods[0].ShortForecast,
		WindSpeed:       fh.Properties.Periods[0].WindSpeed,
		WindDirection:   fh.Properties.Periods[0].WindDirection,
	}
	w.Period = WeatherPeriod{
		Forecast: f.Properties.Periods[0].DetailedForecast,
	}

	// build timeline.
	periods := []WeatherPeriod{}
	max := 6
	for i, period := range fh.Properties.Periods {
		if i%2 != 0 {
			continue // take every other period
		}
		t, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			return nil, err
		}
		p := WeatherPeriod{
			Forecast:        period.DetailedForecast,
			Hour:            t.Hour(),
			Temperature:     period.Temperature,
			TemperatureUnit: period.TemperatureUnit,
		}
		periods = append(periods, p)
		if len(periods) == max {
			break
		}
	}
	w.Timeline = WeatherTimeline{
		Periods: periods,
	}

	return w, nil
}
