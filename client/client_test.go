// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package client

import (
	"io"
	"testing"
)

func TestGet(t *testing.T) {
	res, err := Get("https://plan.cat/~s")
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
