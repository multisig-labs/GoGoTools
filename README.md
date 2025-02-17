<h1 align="center">GoGoTools 🎈</h1>
<p align="center">A (growing) collection of useful tools for Avalanche developers.</p>

## Usage

<pre><code>
Usage:
  ggt <command> ...

Commands:
  bech32-decode            Decode a bech32 address
  cb58-decode              Decode a value from CB58 (ID or NodeID)
  cb58-decode-sig          Decode a signature (r,s,v) from CB58
  cb58-encode              Encode a value to CB58
  completion               Generate shell completion scripts
  help                     Help about any command
  inspect-tx-p             Inspect a P-Chain transaction
  mnemonic-addrs           Show addresses for a BIP39 mnemonic
  mnemonic-generate        Generate a random BIP39 mnemonic
  mnemonic-insecure        Generate an INSECURE test BIP39 mnemonic
  mnemonic-keys            Show keys and addresses for a BIP39 mnemonic
  msgdigest                Generate a hash digest for an Avalanche Signed Message (ERC-191)
  pk                       Show various address encodings of a private key
  random-bls               Generate a random BLS key
  random-id                Generate a random ID
  random-nodeid            Generate a random node ID
  rpc                      Ergonomic access to avalanche node RPC APIs
  verify-bls               Verify a BLS Proof of Possession
  version                  Version
  vmid                     Given a vmName, try to encode the ASCII name as a vmID
  vmname                   Given a vmID, try to decode the ASCII name
  warp-construct-uptime    Construct an uptime message
  warp-get                 Get a warp message from a transaction ID
  warp-parse               Parse a warp message
</code></pre>

## Mnemonics

Avalanche P-Chain and C-Chain use different address formats, and `ggt` provides utilities to help with this.

```sh
❯ bin/ggt mnemonic-keys "test test test test test test test test test test test junk"
=== BIP39 Mnemonic ===
test test test test test test test test test test test junk

=== Ethereum Derivation Path ===
Path               EVM Addr     Ava Addr           EVM Private Key      Ava Private Key
m/44'/60'/0'/0/0   0xf39Fd6...  P-avax15428vq2...  ac0974bec39a17e3...  PrivateKey-2JmTFo8knhffGK32...
m/44'/60'/0'/0/1   0x709979...  P-avax1cjzphr6...  59c6995e998f97a5...  PrivateKey-gYChgv9KmCAaRH47...
m/44'/60'/0'/0/2   0x3C44Cd...  P-avax1sj3m3zu...  5de4111afa1a4b94...  PrivateKey-iMKNGkysaBKiThvd...
m/44'/60'/0'/0/3   0x90F79b...  P-avax1y3fgnts...  7c852118294e51e6...  PrivateKey-wqhSPTuv3JB9YPix...
m/44'/60'/0'/0/4   0x15d34A...  P-avax1jft4f7x...  47e179ec19748859...  PrivateKey-Yf6fhJUE97QgwBkb...
m/44'/60'/0'/0/5   0x996550...  P-avax1syp9y2m...  8b3a350cf5c34c91...  PrivateKey-24KNq4HhQo5BL6xM...
m/44'/60'/0'/0/6   0x976EA7...  P-avax1zyf8ga3...  92db14e403b83dfe...  PrivateKey-27gEa9Qudm22UaxC...
m/44'/60'/0'/0/7   0x14dC79...  P-avax1mer4xr0...  4bbbf85ce3377467...  PrivateKey-aMXjshukTmhmsn2P...
m/44'/60'/0'/0/8   0x23618e...  P-avax1l5rrv44...  dbda1821b80551c9...  PrivateKey-2fpphtBdokVfG6uA...
m/44'/60'/0'/0/9   0xa0Ee7A...  P-avax14flvw0x...  2a871d0798f97d79...  PrivateKey-KjKHSmNQ9bFN3bK2...

=== Avalanche Derivation Path ===
Path               EVM Addr     Ava Addr           EVM Private Key      Ava Private Key
m/44'/9000'/0'/0/0 0x5a299B...  P-avax1yljhuvj...  211cdc80c23ccc8e...  PrivateKey-Fapb8hTUMABpZc9z...
m/44'/9000'/0'/0/1 0xFf9bc6...  P-avax18wvaf02...  9f8799874aeb19dc...  PrivateKey-2DFyMtm5iENeShhP...
m/44'/9000'/0'/0/2 0x9Fa1E0...  P-avax1kk4tuwm...  9ac048b0ccc9a3d9...  PrivateKey-2B9ukB5wfRqS8vnJ...
m/44'/9000'/0'/0/3 0x8E7046...  P-avax1au4cssw...  2119b8ec008c7599...  PrivateKey-FaWRW4Fwpy8dgPq1...
m/44'/9000'/0'/0/4 0xA7a060...  P-avax1s8kxj6a...  edf2bd01bbd1b1c3...  PrivateKey-2oo4uTtBJfU64mB2...
m/44'/9000'/0'/0/5 0xd667fe...  P-avax12r4ys4s...  36bc4c0b5e9b13a3...  PrivateKey-R79TRuBkoJympQD3...
m/44'/9000'/0'/0/6 0xaEF349...  P-avax1sdn9w8t...  e8c61325266f87b6...  PrivateKey-2mWtmQ2iQxEdGwjh...
m/44'/9000'/0'/0/7 0x8F80c9...  P-avax1hvdsp7z...  dad79f166f5359a5...  PrivateKey-2fP2rzMCgN6LvKzK...
m/44'/9000'/0'/0/8 0x6cB771...  P-avax1935pmeg...  2b15356991756b70...  PrivateKey-KyVd5iu1CJ32Qryn...
m/44'/9000'/0'/0/9 0x468E01...  P-avax1yj4kuns...  248c879bf5a22274...  PrivateKey-H6bSJB5xj2FbMuV9...

```

## Installation

Requires [Go](https://golang.org/doc/install) version >= 1.23

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

which will create the binary `bin/ggt` (Make sure you add it to your $PATH)
