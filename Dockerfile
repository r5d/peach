# SPDX-License-Identifier: ISC
# Copyright © 2022 siddharth <s@ricketyspace.net>

FROM golang:1.18.2

WORKDIR /usr/src/peach

COPY . .
RUN make BUILD_OPTS="-v -o /usr/local/bin/peach"

ENTRYPOINT ["peach"]
