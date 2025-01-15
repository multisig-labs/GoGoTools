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
	CGO_ENABLED=1 go build -ldflags "{{LDFLAGS}}" -o bin/ggt main.go

install: build
  mv bin/ggt $GOPATH/bin/ggt

# Whack, but it works
gen-contracts:
	#!/usr/bin/env bash -eo pipefail
	CORETH=0.13.2
	THISDIR=$PWD
	forge build
	binfile=$(mktemp)
	cat artifacts/erc20.sol/CustomERC20.json | jq -r '.bytecode.object' > ${binfile}
	echo "Generating GO code..."
	cd $GOPATH/pkg/mod/github.com/ava-labs/coreth@v${CORETH}
	cat $THISDIR/artifacts/erc20.sol/CustomERC20.json | jq '.abi' | go run ./cmd/abigen/  --bin ${binfile} --pkg erc20 --out $THISDIR/pkg/contracts/erc20/erc20.go --abi -
	cp $THISDIR/artifacts/erc20.sol/CustomERC20.json $THISDIR/pkg/contracts/erc20/erc20.json
