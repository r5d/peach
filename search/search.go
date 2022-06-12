// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package search

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"ricketyspace.net/peach/photon"
	"ricketyspace.net/peach/version"
)

type Search struct {
	Title          string
	Version        string
	Location       string
	Message        string
	MatchingCoords []photon.Coordinates
	Enabled        bool
}

func NewSearch(r *http.Request) (*Search, error, int) {
	s := new(Search)
	s.Title = "search"
	s.Version = version.Version
	s.Enabled = photon.Enabled()

	if !s.Enabled {
		return s, fmt.Errorf("search disabled"), 404
	}
	if r.Method == "GET" {
		return s, nil, 200
	}

	// Get location.
	err := r.ParseForm()
	if err != nil {
		return s, fmt.Errorf("form: %v", err), 500
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
		return s, nil, 200
	}
	if len(s.MatchingCoords) < 1 {
		s.Message = "location not found"
		return s, nil, 200
	}
	return s, nil, 200
}
