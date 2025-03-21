# Autoload .env if it exists
set dotenv-load

# Build vars for versioning the binary
VERSION := `grep "const Version " pkg/version/version.go | sed -E 's/.*"(.+)"$$/\1/'`
GIT_COMMIT := `git rev-parse HEAD`
BUILD_DATE := `date '+%Y-%m-%d'`
VERSION_PATH := "github.com/multisig-labs/gogotools/pkg/version"
LDFLAGS := "-X " + VERSION_PATH + ".BuildDate=" + BUILD_DATE + " -X " + VERSION_PATH + ".Version=" + VERSION + " -X " + VERSION_PATH + ".GitCommit=" + GIT_COMMIT

default:
  @just --list --unsorted

build:
    #!/usr/bin/env sh
    if [ "$(uname -m)" = "arm64" ]; then
        CGO_CFLAGS="-O2 -D__BLST_PORTABLE__" CGO_ENABLED=1 GOARCH=arm64 go build -ldflags "{{LDFLAGS}}" -o bin/ggt cmd/*
    else
        CGO_CFLAGS="-O2 -D__BLST_PORTABLE__" CGO_ENABLED=1 GOARCH=amd64 go build -ldflags "{{LDFLAGS}}" -o bin/ggt cmd/*
    fi

install: build
  mv bin/ggt $GOPATH/bin/ggt

