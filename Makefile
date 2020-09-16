#@IgnoreInspection BashAddShebang

export APP=nsm

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export LDFLAGS="-w -s"

all: format lint build

############################################################
# Build and Run
############################################################

build:
	CGO_ENABLED=1 go build -o ${APP} -ldflags $(LDFLAGS)  .

# Please do not use `make run` on production. There is a performance hit due to existence of -race flag.
run:
	go run -race -ldflags $(LDFLAGS) . server

############################################################
# Format and Lint
############################################################

check-formatter:
	which gofumpt || GO111MODULE=off go get -u mvdan.cc/gofumpt
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	which godot || GO111MODULE=off go get -u github.com/tetafro/godot/cmd/godot

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofumpt -s -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R godot -w R

check-linter:
	which golangci-lint || GO111MODULE=off curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.30.0

lint: check-linter
	golangci-lint run --deadline 10m $(ROOT)/...

############################################################
# Test
############################################################

test:
	go test -v -race -p 1 ./...

ci-test:
	go test -v -race -p 1 -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func coverage.txt

############################################################
# Development Environment
############################################################

up:
	docker-compose up -d

down:
	docker-compose down
	rm -rf data
