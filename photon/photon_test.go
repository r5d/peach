// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package photon

import (
	"os"
	"testing"
)

func TestEnabled(t *testing.T) {
	if Enabled() {
		t.Errorf("geo is enabled")
		return
	}

	os.Setenv("PEACH_PHOTON_URL", "https://photon.komoot.io")
	if !Enabled() {
		t.Errorf("geo is not enabled")
		return
	}
}

func TestPhotonUrl(t *testing.T) {
	os.Setenv("PEACH_PHOTON_URL", "")
	pUrl, err := Url()
	if err == nil {
		t.Errorf("url: %v", pUrl)
		return
	}

	os.Setenv("PEACH_PHOTON_URL", "https://photon.komoot.io")
	pUrl, err = Url()
	if err != nil {
		t.Errorf("url: %v", err)
		return
	}
	if pUrl.String() != "https://photon.komoot.io/api" {
		t.Errorf("url: %v", pUrl)
		return
	}
}

func TestGeocode(t *testing.T) {
	os.Setenv("PEACH_PHOTON_URL", "https://photon.komoot.io")
	mCoords, err := Geocode("Tiffin,OH")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if len(mCoords) < 0 {
		t.Errorf("no matching coordinates")
		return
	}
	if mCoords[0].Lat != 41.114485 {
		t.Errorf("lat: %v", mCoords[0].Lat)
		return
	}
	if mCoords[0].Lng != -83.1779537 {
		t.Errorf("lng: %v", mCoords[0].Lat)
		return
	}
	if mCoords[0].Name != "Tiffin, Ohio" {
		t.Errorf("name: %v", mCoords[0].Name)
		return
	}
}
