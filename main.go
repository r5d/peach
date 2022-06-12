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

	"ricketyspace.net/peach/photon"
	"ricketyspace.net/peach/version"
	"ricketyspace.net/peach/weather"
)

// Peach port. Defaults to 8151
var peachPort = flag.Int("p", 8151, "Port to run peach on")

// Peach listen address. Set during init.
var peachAddr = ""

// Holds static content.
//go:embed templates static/*.min.css static/font/*
var peachFS embed.FS

// HTML templates.
var peachTemplates = template.Must(template.ParseFS(peachFS, "templates/*.tmpl"))

// Lat,Long regex.
var latLngRegex = regexp.MustCompile(`/(-?[0-9]+\.?[0-9]+?),(-?[0-9]+\.?[0-9]+)`)

type Search struct {
	Title          string
	Version        string
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
	// Search handler.
	http.HandleFunc("/search", showSearch)

	// Default handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		if r.URL.Path == "/" {
			http.Redirect(w, r, "/41.115,-83.177", 302)
			return
		}
		if r.URL.Path == "/version" {
			fmt.Fprintf(w, "v%s\n", version.Version)
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

	// Static files handler.
	http.HandleFunc("/static/", serveStaticFile)

	// Start server
	log.Fatal(http.ListenAndServe(peachAddr, nil))
}

func showWeather(w http.ResponseWriter, lat, lng float32) {
	// Make weather
	weather, err, status := weather.NewWeather(lat, lng)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// Render.
	err = peachTemplates.ExecuteTemplate(w, "weather.tmpl", weather)
	if err != nil {
		log.Printf("weather: template: %v", err)
		return
	}
}

func showSearch(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

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

func serveStaticFile(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// Add Cache-Control header
	w.Header().Set("Cache-Control", "max-age=604800")

	// Serve.
	server := http.FileServer(http.FS(peachFS))
	server.ServeHTTP(w, r)
}

func NewSearch(r *http.Request) (*Search, error) {
	s := new(Search)
	s.Title = "search"
	s.Version = version.Version

	if r.Method == "GET" {
		return s, nil
	}

	// Get location.
	err := r.ParseForm()
	if err != nil {
		return s, fmt.Errorf("form: %v", err)
	}
	location := strings.TrimSpace(r.PostForm.Get("location"))
	s.Location = location
	if len(location) < 2 {
		s.Message = "location invalid"
	}

	// Try to fetch matching coordinates.
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

func logRequest(r *http.Request) {
	addr := r.RemoteAddr
	if len(r.Header.Get("X-Forwarded-For")) > 0 {
		addr = r.Header.Get("X-Forwarded-For")
	}
	log.Printf("%v - %v - %v", addr, r.URL, r.Header.Get("User-agent"))
}
