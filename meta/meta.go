// Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
// SPDX-License-Identifier: ISC

package meta

import "ricketyspace.net/peach/version"

type Meta struct {
	Version string
}

func NewMeta() *Meta {
	m := new(Meta)
	m.Version = version.Version
	return m
}
