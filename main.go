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
	"ricketyspace.net/peach/photon"
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
	Title    string
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

type Search struct {
	Title          string
	Location       string
	Message        string
	MatchingCoords []photon.Coordinates
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

	// search handler.
	http.HandleFunc("/search", showSearch)

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

	// start server
	log.Fatal(http.ListenAndServe(peachAddr, nil))
}

func showWeather(w http.ResponseWriter, lat, lng float32) {
	point, err := nws.Points(lat, lng)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// get forecast
	f, err := nws.GetForecast(point)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fh, err := nws.GetForecastHourly(point)
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
		log.Printf("weather: template: %v", err)
		return
	}
}

func showSearch(w http.ResponseWriter, r *http.Request) {
	// Search is disabled if photon is not enabled.
	if !photon.Enabled() {
		http.NotFound(w, r)
		return
	}

	search, err := NewSearch(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = peachTemplates.ExecuteTemplate(w, "search.tmpl", search)
	if err != nil {
		log.Printf("search: template: %v", err)
		return
	}
}

func NewWeather(point *nws.Point, f, fh *nws.Forecast) (*Weather, error) {
	w := new(Weather)
	w.Location = fmt.Sprintf("%s, %s",
		strings.ToLower(point.Properties.RelativeLocation.Properties.City),
		strings.ToLower(point.Properties.RelativeLocation.Properties.State),
	)
	w.Title = w.Location
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

func NewSearch(r *http.Request) (*Search, error) {
	s := new(Search)
	s.Title = "search"

	if r.Method == "GET" {
		return s, nil
	}

	// get location.
	err := r.ParseForm()
	if err != nil {
		return s, fmt.Errorf("form: %v", err)
	}
	location := strings.TrimSpace(r.PostForm.Get("location"))
	s.Location = location
	if len(location) < 2 {
		s.Message = "location invalid"
	}

	// try to fetch matching coordinates.
	s.MatchingCoords, err = photon.Geocode(location)
	if err != nil {
		log.Printf("search: geocode: %v", err)
		s.Message = "unable to lookup location"
		return s, nil
	}
	if len(s.MatchingCoords) < 1 {
		s.Message = "location not found"
		return s, nil
	}
	return s, nil
}
