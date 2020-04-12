BUILD_DIRECTORY=_build
PROGRAM_NAME=meet

all: fmt test lint build

ci: tools dependencies generate fmt check-repository-unchanged test lint build

build-directory:
	mkdir -p ./${BUILD_DIRECTORY}

build: build-directory
	go build -o ./${BUILD_DIRECTORY}/${PROGRAM_NAME} ./cmd/${PROGRAM_NAME}

build-race: build-directory
	go build -race -o ./${BUILD_DIRECTORY}/${PROGRAM_NAME} ./cmd/${PROGRAM_NAME}

frontend:
	./_tools/build_frontend.sh

check-repository-unchanged:
	./_tools/check_repository_unchanged.sh

tools:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/google/wire/cmd/wire
	go get -u golang.org/x/tools/cmd/goimports

dependencies:
	go get ./...

generate:
	 go generate ./...

lint: 
	go vet ./...
	staticcheck ./...

doc:
	@echo "http://localhost:6060/pkg/github.com/boreq/${PROGRAM_NAME}/"
	godoc -http=:6060

fmt:
	goimports -l -w adapters/ application/ cmd/ domain/ internal/ ports/

test:
	go test ./...

test-verbose:
	go test -v ./...

clean:
	rm -rf ./${BUILD_DIRECTORY}

.PHONY: all build build-directory frontend check-repository-unchanged build-race tools dependencies lint doc fmt test test-verbose clean
