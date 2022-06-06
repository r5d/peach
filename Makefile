# SPDX-License-Identifier: ISC
# Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>

MOD=ricketyspace.net/peach
PKGS=${MOD}/client ${MOD}/nws ${MOD}/photon

peach: fix fmt
	go build ${BUILD_OPTS}

fmt:
	go fmt ${MOD} ${PKGS}
.PHONY: fmt

fix:
	go fix ${MOD} ${PKGS}
.PHONY: fix

test:
	go test ${PKGS} ${ARGS}
.PHONY: test

clean:
	go clean
.PHONY: clean
