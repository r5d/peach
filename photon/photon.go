// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package photon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"ricketyspace.net/peach/client"
)

// Coordinates.
type Coordinates struct {
	Lat  float32
	Lng  float32
	Name string
}

// Represents a response from the Photon API.
type Response struct {
	Features []Feature
}

// Represents features object in the Response.
type Feature struct {
	Geometry   Geometry
	Properties Properties
}

// Represents geometry object in the Response.
type Geometry struct {
	Coordinates []float32
}

// Represents properties object in the Response.
type Properties struct {
	CountryCode string
	Name        string
	State       string
}

// Returns true of geocoding is possible.
func Enabled() bool {
	return len(os.Getenv("PEACH_PHOTON_URL")) > 0
}

func Url() (*url.URL, error) {
	if !Enabled() {
		return nil, fmt.Errorf("geocoding not enabled")
	}

	pu, err := url.Parse(os.Getenv("PEACH_PHOTON_URL"))
	if err != nil {
		return nil, err
	}
	if len(pu.Path) < 1 || pu.Path == "/" {
		pu.Path = "/api"
	}
	return pu, nil
}

// Returns a list of matching Coordinates for a given location.
func Geocode(location string) ([]Coordinates, error) {
	mCoords := []Coordinates{} // Matching coordinates
	location = strings.TrimSpace(location)
	if len(location) < 2 {
		return mCoords, fmt.Errorf("geocode: location invalid")
	}

	// Construct request.
	u, err := Url()
	if err != nil {
		return mCoords, fmt.Errorf("geocode: %v", err)
	}
	q := url.Values{}
	q.Add("q", location)
	q.Add("osm_tag", "place:city")
	q.Add("limit", "5")
	u.RawQuery = q.Encode()

	// Make request.
	resp, err := client.Get(u.String())
	if err != nil {
		return mCoords, fmt.Errorf("geocode: get: %v", err)
	}

	// Parse response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mCoords, fmt.Errorf("geocode: body: %v", err)
	}

	// Check if the request failed.
	if resp.StatusCode != 200 {
		return mCoords, fmt.Errorf("geocode: %s", body)
	}

	// Unmarshal
	r := new(Response)
	err = json.Unmarshal(body, r)
	if err != nil {
		return mCoords, fmt.Errorf("geocode: decode: %v", err)
	}

	// Make matching coordinates list.
	for _, feature := range r.Features {
		if feature.Properties.CountryCode != "US" {
			continue // skip
		}

		c := Coordinates{}
		c.Lat = feature.Geometry.Coordinates[1]
		c.Lng = feature.Geometry.Coordinates[0]
		c.Name = fmt.Sprintf("%s, %s",
			feature.Properties.Name,
			feature.Properties.State,
		)
		mCoords = append(mCoords, c)
	}
	return mCoords, nil
}
