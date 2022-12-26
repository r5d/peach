# SPDX-License-Identifier: ISC
# Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>

MOD=ricketyspace.net/peach
PKGS=${MOD}/cache ${MOD}/client ${MOD}/nws ${MOD}/photon ${MOD}/time
CSS=static/peach.min.css

peach: vet fix fmt ${CSS}
	go build -race ${BUILD_OPTS}

fmt:
	go fmt ${MOD} ${PKGS}
.PHONY: fmt

fix:
	go fix ${MOD} ${PKGS}
.PHONY: fix

vet:
	go vet ${MOD} ${PKGS}
.PHONY: vet

test:
	go test -race ${PKGS} ${ARGS}
.PHONY: test

${CSS}: static/peach.css
	./bin/minify

clean:
	go clean
.PHONY: clean
