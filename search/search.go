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
