# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth <s@ricketyspace.net>

MOD=ricketyspace.net/peach

peach: fmt
	go build

fmt:
	go fmt ${MOD} ${MOD}/http/client

