// Copyright © 2022 siddharth <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// Thin HTTP client wrapper.
package client

import "net/http"

// HTTP client.
var client = http.Client{}

// Make a HTTP GET request.
func Get(url string) (resp *http.Response, err error) {
	return nil, nil
}

// Add default headers for the peach http client.
func (req *http.Request) buildHeaders() {
	req.Header.Set("User-Agent", "peach/"+peach.Version+" ricketyspace.net/contact")
}
