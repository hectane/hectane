COLORON = tput setaf 3
COLOROFF = tput init

CWD = $(shell pwd)
PKG = github.com/hectane/hectane
CMD = hectane

UID = $(shell id -u)
GID = $(shell id -g)

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')
UIFILES = $(shell find ui/public ui/src) ui/package-lock.json

all: dist/${CMD}

dist/${CMD}: ${SOURCES} server/ab0x.go | cache dist
	@$(COLORON)
	@echo "Building Hectane binary..."
	@$(COLOROFF)
	@docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib ${PKG}/cmd/${CMD}
	@touch dist/hectane

server/ab0x.go: dist/fileb0x b0x.yaml .dep-static
	@$(COLORON)
	@echo "Generating Go source from static files..."
	@$(COLOROFF)
	@dist/fileb0x b0x.yaml

dist/fileb0x: | cache dist
	@$(COLORON)
	@echo "Building fileb0x binary..."
	@$(COLOROFF)
	@docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib github.com/UnnoTed/fileb0x

.dep-static: ${UIFILES} .dep-node_modules
	@$(COLORON)
	@echo "Building web UI..."
	@$(COLOROFF)
	@docker run \
	    --rm \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/ui:/usr/src/ui \
	    -w /usr/src/ui \
	    nathanosman/betternode \
	    npm run build
	@cp -rT ui/build server/www
	@touch .dep-static

.dep-node_modules: ui/package.json
	@$(COLORON)
	@echo "Fetching Node packages..."
	@$(COLOROFF)
	@docker run \
	    --rm \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/ui:/usr/src/ui \
	    -w /usr/src/ui \
	    nathanosman/betternode \
	    npm install
	@touch .dep-node_modules

cache:
	@mkdir cache

dist:
	@mkdir dist

clean:
	@rm -f .dep-* server/ab0x.go
	@rm -rf cache dist server/www ui/build ui/node_modules

.PHONY: clean
