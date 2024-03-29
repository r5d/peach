// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

// A simple in-memory cache store.
package cache

import "time"

// An item in the key-value cache store.
type item struct {
	value   []byte
	expires time.Time // Time when the key-value expires
}

// A key-value cache store.
type Cache struct {
	sema  chan int // Semaphore for read/write access to cache.
	store map[string]item
}

// Returns a new empty cache store.
func NewCache() *Cache {
	c := new(Cache)
	c.sema = make(chan int, 1)
	c.store = make(map[string]item)
	return c
}

// Set the (key,value) item to the cache store. This item will be
// considered expired after time `expires`.
//
// Cache.Get will return an empty string once `expires` is past the
// current time.
func (c *Cache) Set(key string, value []byte, expires time.Time) {
	// Get sema token before accessing the cache.
	c.sema <- 1
	defer func() {
		// Give up sema token.
		<-c.sema
	}()
	c.store[key] = item{
		value:   value,
		expires: expires,
	}
}

// Get an (key,value) item from the cache store by key.
//
// An empty []byte will be returned when if the key does not exist or
// if the item corresponding to the key has expired.
func (c *Cache) Get(key string) []byte {
	// Get sema token before accessing the cache.
	c.sema <- 1
	defer func() {
		// Give up sema token.
		<-c.sema
	}()

	if _, ok := c.store[key]; !ok {
		return []byte{}
	}
	// Check if the item expired.
	if time.Until(c.store[key].expires).Seconds() < 0 {
		return []byte{}
	}
	return c.store[key].value
}
