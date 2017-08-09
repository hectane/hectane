CWD = $(shell pwd)
PKG = github.com/hectane/hectane
CMD = hectane

UID = $(shell id -u)
GID = $(shell id -g)

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')

all: dist/${CMD}

dist/${CMD}: ${SOURCES} | cache dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    -w /go/src/${PKG} \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib ${PKG}/cmd/${CMD}

cache:
	@mkdir cache

dist:
	@mkdir dist

clean:
	@rm -rf cache dist

.PHONY: clean
