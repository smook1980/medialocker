VERSION = $(shell git describe --always --tags)
COMMIT = $(shell git rev-parse --short=8 HEAD)

SOURCES := $(shell find . -name '*.go')
PKGS := $(shell glide novendor)

LDFLAGS=-ldflags "-s -X medialocker.Version=${VERSION} -X medialocker.Commit=${COMMIT}"
DEV_LDFLAGS=-ldflags "-s -X medialocker.Version=snapshot -X medialocker.Commit=${COMMIT}"
BINARY=./bin/locker
BINARY-DEV=./bin/locker-dev

GLIDE := $(shell realpath ${GOPATH}/bin/glide)
GO-BINDATA := $(shell realpath ${GOPATH}/bin/go-bindata)

release: dep assets test ${BINARY}

install: dep assets test

build: assets test ${BINARY}

build-dev: assets test ${BINARY-DEV}

run: bindata
	go run -tags release ${LDFLAGS} ./cmd/locker/main.go

run-dev: bindata
	go run -tags dev ${DEV_LDFLAGS} ./cmd/locker/main.go

$(BINARY-DEV): $(SOURCES)
	go build -tags dev -o $(BINARY-DEV) ${DEV_LDFLAGS} ./cmd/locker/main.go

${BINARY}: $(SOURCES)
	go build -tags release -o ${BINARY} ${LDFLAGS} ./cmd/locker/main.go

assets: js bindata

js:
	cd ui && npm run build

bindata: ${GO-BINDATA}
	go generate -x .

${GLIDE}:
	curl https://glide.sh/get | sh

${GO-BINDATA}:
	go get -u github.com/jteeuwen/go-bindata/...

dep: jsdep godep

godep: ${GLIDE}
	glide install

jsdep:
	cd ui && yarn install

test: jstest gotest

gotest:
	go test -i -timeout 10s -race ${PKGS}

jstest:
	cd ui && yarn run test

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf ui/dist || true

uikit_update:
	@echo "Using git subtree to pull latest changes to ui/vendor/uikit"
	git subtree pull --prefix ./ui/vendor/uikit uikit/uikit master --squash

.PHONY: uikit_update clean jstest gotesttrace gotest test jsdep godep dep js bindata assets run-dev run
