# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth ravikumar <s@ricketyspace.net>

MOD=ricketyspace.net/peach

peach: fmt
	go build ${BUILD_OPTS}

fmt:
	go fmt ${MOD} ${MOD}/client ${MOD}/nws ${MOD}/photon

test:
	go test ${MOD}/client ${MOD}/nws ${MOD}/photon ${ARGS}
.PHONY: test

clean:
	go clean
.PHONY: clean

image:
	./bin/image
.PHONY: image

image-push:
	./bin/image-push
.PHONY: image-push
