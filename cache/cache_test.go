// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package cache

import (
	"bytes"
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
		value:   []byte("bar"),
		expires: time.Now().Add(time.Second * 10),
	}
	if bytes.Compare(c.store["foo"].value, []byte("bar")) != 0 {
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
	c.Set("foo", []byte("bar"), exp)
	if bytes.Compare(c.store["foo"].value, []byte("bar")) != 0 {
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
	c.Set("foo", []byte("bar"), exp)
	if bytes.Compare(c.Get("foo"), []byte("bar")) != 0 {
		t.Errorf("cache.Get(foo) is not bar")
		return
	}

	// Test 2
	exp = time.Now().Add(time.Second * -10)
	c.Set("sna", []byte("fu"), exp)
	if bytes.Compare(c.Get("sna"), []byte{}) != 0 {
		t.Errorf("cache.Get(sna) is not empty: %s", c.Get("sna"))
		return
	}
}
