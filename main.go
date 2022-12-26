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

	"ricketyspace.net/peach/search"
	"ricketyspace.net/peach/version"
	"ricketyspace.net/peach/weather"
)

// Peach port. Defaults to 8151
var peachPort = flag.Int("p", 8151, "Port to run peach on")

// Peach listen address. Set during init.
var peachAddr = ""

// Holds static content.
//
//go:embed templates static/peach.min.css
//go:embed static/font/roboto-flex.ttf
//go:embed static/logo/peach-*.png
var peachFS embed.FS

// HTML templates.
var peachTemplates = template.Must(template.ParseFS(peachFS, "templates/*.tmpl"))

// Lat,Long regex.
var latLngRegex = regexp.MustCompile(`/(-?[0-9]+\.?[0-9]+?),(-?[0-9]+\.?[0-9]+)`)

func init() {
	flag.Parse()
	if *peachPort < 80 {
		log.Fatalf("port number is invalid: %v", *peachPort)
	}
	peachAddr = fmt.Sprintf(":%d", *peachPort)
}

func main() {
	// Default handler.
	http.HandleFunc("/", defaultHandler)

	// Static files handler.
	http.HandleFunc("/static/", serveStaticFile)

	// Search handler.
	http.HandleFunc("/search", showSearch)

	// Meta handler.
	http.HandleFunc("/about", showMeta)

	// Start server
	log.Fatal(http.ListenAndServe(peachAddr, nil))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
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

func showMeta(w http.ResponseWriter, r *http.Request) {
	// Make meta info.
	type Meta struct {
		Version string
		Title   string
	}
	m := new(Meta)
	m.Version = version.Version
	m.Title = "about"

	// Render.
	err := peachTemplates.ExecuteTemplate(w, "about.tmpl", m)
	if err != nil {
		log.Printf("weather: template: %v", err)
		return
	}
}

func serveStaticFile(w http.ResponseWriter, r *http.Request) {
	// Add Cache-Control header
	w.Header().Set("Cache-Control", "max-age=604800")

	// Serve.
	server := http.FileServer(http.FS(peachFS))
	server.ServeHTTP(w, r)
}

func showSearch(w http.ResponseWriter, r *http.Request) {
	search, err, status := search.NewSearch(r)
	if err != nil && status == 404 {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	err = peachTemplates.ExecuteTemplate(w, "search.tmpl", search)
	if err != nil {
		log.Printf("search: template: %v", err)
		return
	}
}
