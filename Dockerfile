# SPDX-License-Identifier: ISC
# Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>

FROM golang:1.19.4

WORKDIR /usr/src/peach

COPY . .
RUN make BUILD_OPTS="-v -o /usr/local/bin/peach"

ENTRYPOINT ["peach"]
