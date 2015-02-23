#!/bin/sh

set -e

# Grab dependencies for coveralls.io integration
go get -u github.com/axw/gocov/gocov
go get -u github.com/mattn/goveralls
go get -u golang.org/x/tools/cmd/cover
# Grab all project dependencies
go get -t -v ./...
go get
go build
go test -v ./...

go tool cover -func=coverage
