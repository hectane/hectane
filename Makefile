CWD = $(shell pwd)
PKG = github.com/hectane/hectane
CMD = hectane

UID = $(shell id -u)
GID = $(shell id -g)

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')
BINDATA = $(shell find server/static server/templates)

all: dist/${CMD}

dist/${CMD}: ${SOURCES} server/ab0x.go | cache dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib ${PKG}/cmd/${CMD}

cache:
	@mkdir cache

dist:
	@mkdir dist

server/ab0x.go: ${BINDATA} | dist/fileb0x
	dist/fileb0x b0x.yaml

dist/fileb0x: | dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib github.com/UnnoTed/fileb0x

clean:
	@rm -rf cache dist

.PHONY: clean
