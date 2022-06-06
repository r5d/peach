// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package cache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Errorf("cache is nil")
		return
	}
	if c.store == nil {
		t.Errorf("cache.store is nil")
		return
	}
	// Try manually adding an item.
	c.store["foo"] = item{
		value:   "bar",
		expires: time.Now().Add(time.Second * 10),
	}
	if c.store["foo"].value != "bar" {
		t.Errorf("cache.store['foo'] is not bar")
		return
	}
}

func TestCacheSet(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Errorf("cache is nil")
		return
	}
	exp := time.Now().Add(time.Second * 10)
	c.Set("foo", "bar", exp)
	if c.store["foo"].value != "bar" {
		t.Errorf("cache.store['foo'] is not bar")
		return
	}
	if c.store["foo"].expires.Unix() != exp.Unix() {
		t.Errorf("cache.store['foo'].expires no set correctly")
		return
	}
}

func TestCacheGet(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Errorf("cache is nil")
		return
	}

	// Test 1
	exp := time.Now().Add(time.Second * 10)
	c.Set("foo", "bar", exp)
	if c.Get("foo") != "bar" {
		t.Errorf("cache.Get(foo) is not bar")
		return
	}

	// Test 2
	exp = time.Now().Add(time.Second * -10)
	c.Set("sna", "fu", exp)
	if c.Get("sna") != "" {
		t.Errorf("cache.Get(sna) is not empty: %s", c.Get("sna"))
		return
	}
}
