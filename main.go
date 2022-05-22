// Copyright © 2022 siddharth <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"ricketyspace.net/peach/nws"
)

// Peach port. Defaults to 8151
var peachPort = flag.Int("p", 8151, "Port to run peach on")

// Peach listen address. Set during init.
var peachAddr = ""

// Holds static content.
//go:embed templates static
var peachFS embed.FS

// HTML templates.
var peachTemplates = template.Must(template.ParseFS(peachFS, "templates/*.tmpl"))

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
	http.Handle("/static/", http.FileServer(http.FS(peachFS)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path[1:]) != 0 {
			http.NotFound(w, r)
			return
		}
		showWeather(w, 41.115, -83.177)
	})
	log.Fatal(http.ListenAndServe(peachAddr, nil))
}

func showWeather(w http.ResponseWriter, lat, lng float32) {
	point, err := nws.Points(lat, lng)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get forecast
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

	// Make weather
	weather, err := NewWeather(point, f, fh)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = peachTemplates.ExecuteTemplate(w, "weather.tmpl", weather)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func NewWeather(point *nws.NWSPoint, f, fh *nws.NWSForecast) (*Weather, error) {
	w := new(Weather)
	w.Location = fmt.Sprintf("%s, %s",
		point.Properties.RelativeLocation.Properties.City,
		point.Properties.RelativeLocation.Properties.State,
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

	// Build timeline.
	periods := []WeatherPeriod{}
	max := 12
	for i, period := range fh.Properties.Periods {
		if i%2 != 0 {
			continue // Take every other period
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
