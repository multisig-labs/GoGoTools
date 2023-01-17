<h1 align="center">GoGoTools</h1>
<p align="center">A (growing) collection of useful tools for Avalanche developers.</p>

## Installation

Requires [Go](https://golang.org/doc/install) version >= 1.19

Clone Repository

```sh
git clone https://github.com/multisig-labs/gogotools.git
cd gogotools
```

Install [Just](https://github.com/casey/just) for your system, something like:

```sh
brew install just
  or
cargo install just
  or
apk add just
  etc
```

Then build with

```sh
just build
```

which will create the binary `bin/ggt`

## Usage

```
Usage:
  ggt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  msgdigest   Generate a hash digest of a message
  subnetid    Generate a valid SubnetID from a name string (max 32 chars)
  version     Show version

Flags:
  -h, --help   help for ggt

Use "ggt [command] --help" for more information about a command.
```
