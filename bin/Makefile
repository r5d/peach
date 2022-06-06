SH_FILES=minify

all: check fmt
.PHONY: all

check:
	shellcheck -x ${SH_FILES}
.PHONY: check

fmt:
	shfmt -w -i 4 -s ${SH_FILES}
.PHONY: fmt
