# Autoload .env if it exists
set dotenv-load

VERSION := `grep "const Version " pkg/version/version.go | sed -E 's/.*"(.+)"$$/\1/'`
GIT_COMMIT := `git rev-parse HEAD`
BUILD_DATE := `date '+%Y-%m-%d'`
VERSION_PATH := "github.com/multisig-labs/gogotools/pkg/version"
LDFLAGS := "-X " + VERSION_PATH + ".BuildDate=" + BUILD_DATE + " -X " + VERSION_PATH + ".Version=" + VERSION + " -X " + VERSION_PATH + ".GitCommit=" + GIT_COMMIT

export KEYSTORE_PASSWORD := "jhgGJ4hg4"

default:
  @just --list --unsorted

build:
	CGO_ENABLED=1 go build -ldflags "{{LDFLAGS}}" -o bin/ggt main.go

info:
	xh -b :9650/ext/info    id=1 jsonrpc="2.0" "params[username]=admin" "params[password]={{KEYSTORE_PASSWORD}}" method="info.getNodeVersion"
	xh -b :9650/ext/info    id=1 jsonrpc="2.0" "params[username]=admin" "params[password]={{KEYSTORE_PASSWORD}}" method="info.getVMs"
	xh -b :9650/ext/bc/P    id=1 jsonrpc="2.0" "params[username]=admin" "params[password]={{KEYSTORE_PASSWORD}}" method="platform.getSubnets"
	xh -b :9650/ext/bc/P    id=1 jsonrpc="2.0" "params[username]=admin" "params[password]={{KEYSTORE_PASSWORD}}" method="platform.getBlockchains"
	xh -b :9650/ext/admin   id=1 jsonrpc="2.0" "params[username]=admin" "params[password]={{KEYSTORE_PASSWORD}}" method="admin.getChainAliases" "params[chain]=nYKSQqw15js3YmQxX1Z8aC13jr3EA39V9x1DQpep9yKz9TYPc"
