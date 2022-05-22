# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth <s@ricketyspace.net>

MOD=ricketyspace.net/peach

peach: fmt
	go build ${BUILD_OPTS}

fmt:
	go fmt ${MOD} ${MOD}/client ${MOD}/nws

test:
	go test ${MOD}/client ${MOD}/nws ${ARGS}
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
