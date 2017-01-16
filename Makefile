VERSION ?= $(shell git describe --always --tags)
COMMIT ?= $(shell git rev-parse --short=8 HEAD)

SOURCES := $(shell find . -name '*.go')

LDFLAGS=-ldflags "-s -X medialocker.Version=${VERSION} -X medialocker.Commit=${COMMIT}"
DEV_LDFLAGS=-ldflags "-s -X medialocker.Version=snapshot -X medialocker.Commit=${COMMIT}"
BINARY=./bin/locker

GLIDE := $(shell realpath ${GOPATH}/bin/glide)
RICE := $(shell realpath ${GOPATH}/bin/rice)

default: dep build

build: assets ${BINARY}

dev-server:
	go run -tags dev ${DEV_LDFLAGS} ./cmd/locker/main.go

dev-bin:
	go build -tags dev -o ./bin/locker-dev ${DEV_LDFLAGS} ./cmd/locker/main.go

${BINARY}: $(SOURCES)
	go build -o ${BINARY} ${LDFLAGS} ./cmd/locker/main.go

# docker-${BINARY}: $(SOURCES)
#		CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o ${BINARY} ${LDFLAGS} \
#			./cmd/chronograf/main.go

# docker: dep assets docker-${BINARY}
#		docker build -t chronograf .

assets: js bindata

# dev-assets: dev-js bindata

bindata: ${RICE}
	go generate -x ./assets

js:
	cd ui && npm run build

# dev-js:
# 	cd ui && yarn build

dep: jsdep godep

${GLIDE}:
	curl https://glide.sh/get | sh

${RICE}:
	go get -u -v github.com/GeertJohan/go.rice/rice

godep: ${GLIDE}
	glide install

jsdep:
	cd ui && yarn install

# gen: bolt/internal/internal.proto
#		go generate -x ./bolt/internal

test: jstest gotest gotestrace

gotest:
	go test ./...

gotestrace:
	go test -race ./...

jstest:
	# cd ui && npm test

run: ${BINARY}
	./bin/locker
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	find . -type f -name rice-box.go -print0 | xargs -0 rm
	cd assets && rm -rf web

.PHONY: clean test jstest gotest run
