// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Thin HTTP client wrapper.
package client

import (
	"net/http"

	"ricketyspace.net/peach/version"
)

// HTTP client.
var client = http.Client{}

// Make a HTTP GET request.
func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(buildHeaders(req))
}

// Add default headers for the peach http client.
func buildHeaders(req *http.Request) *http.Request {
	req.Header.Set("User-Agent", "peach/"+version.Version+
		" ricketyspace.net/peach")
	return req
}
