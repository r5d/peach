// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"ricketyspace.net/peach/version"
)

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check user-agent header.
		expectedUA := fmt.Sprintf("peach/%s peach.ricketyspace.net",
			version.Version)
		if r.Header.Get("User-Agent") != expectedUA {
			t.Errorf("header: user agent: %v != %v",
				r.Header.Get("User-Agent"), expectedUA)
			return
		}

		// Check cache-control header.
		if r.Header.Get("Cache-Control") != "max-age=0" {
			t.Errorf("header: cache control: %v != max-age=0",
				r.Header.Get("Cache-Control"))
			return
		}

		// Check accept header.
		if r.Header.Get("Accept") != "application/geo+json" {
			t.Errorf("header: accept: %v != application/geo+json",
				r.Header.Get("Accept"))
			return
		}
		fmt.Fprint(w, "OK")
	}))
	defer ts.Close()

	res, err := Get(ts.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
		return
	}
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("response read failed: %v", err)
		return
	}
}
