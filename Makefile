GIT_HASH = $(shell git rev-parse HEAD | tr -d "\n")
VERSION = $(shell git describe --tags --always --dirty --match=*.*.*)
GO_PKGS= \
    github.com/rafalkorepta/UmFmYWxLb3JlcHRhcmVjcnVpdG1lbnQgdGFzaw-/internal/...

define linker_flags
-X main.Version=$(VERSION) \
-X main.Commit=$(GIT_HASH)
endef

all: backend
.PHONY: all

init:
	go get -u github.com/axw/gocov/...c
	go get -u github.com/AlekSi/gocov-xml
.PHONY: init

backend: lint test-backend build-backend
.PHONY: backend

lint:
	golangci-lint run
.PHONY: lint

test-backend:
	echo "mode: set" > coverage-all.out
	$(foreach pkg,$(GO_PKGS),\
		go test -v -race -coverprofile=coverage.out $(pkg) | tee -a test-results.out || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out || exit 1;)
	go tool cover -func=coverage-all.out
.PHONY: test-backend

build-locally: build-linux-backend build-container-locally
.PHONY: build-locally

build-container-locally:
	docker build -t $(DOCKER_USERNAME)/weather-backend:local-latest .
.PHONY: build-container-locally

build-backend:
	go build -ldflags '$(linker_flags) -s' -o dist/backend main.go
.PHONY: build-backend

build-linux-backend:
	env GOOS=linux GOARCH=amd64 go build -ldflags '$(linker_flags) -s' -o dist/backend main.go
.PHONY: build-linux-backend

deploy:
	docker build -f Dockerfile -t $(DOCKER_USERNAME)/weather-backend:$(VERSION) .
	docker push $(DOCKER_USERNAME)/weather-backend:$(VERSION)
	docker logout
.PHONY: deploy
