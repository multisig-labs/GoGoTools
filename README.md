<h1 align="center">GoGoTools 🎈</h1>
<p align="center">A (growing) collection of useful tools for Avalanche developers.</p>

## Usage

<pre><code>
Usage:
  ggt <command> ...

Commands:
Commands:
  balance                      Get the balance of an address
  balance-pk                   Get the balance of an address for a private key
  bech32-decode                Decode a bech32 address
  cb58-decode                  Decode a value from CB58 (ID or NodeID)
  cb58-decode-sig              Decode a signature (r,s,v) from CB58
  cb58-encode                  Encode a value to CB58
  completion                   Generate shell completion scripts
  cross-chain-tx               Transfer assets from C-Chain to P-Chain
  help                         Help about any command
  inspect-tx-p                 Inspect a P-Chain transaction
  l1-validators                Get current validators from a L1 validator RPC endpoint
  mnemonic-addrs               Show addresses for a BIP39 mnemonic
  mnemonic-generate            Generate a random BIP39 mnemonic
  mnemonic-insecure            Generate an INSECURE test BIP39 mnemonic
  mnemonic-keys                Show keys and addresses for a BIP39 mnemonic
  msgdigest                    Generate a hash digest for an Avalanche Signed Message (ERC-191)
  pk                           Show various address encodings of a private key
  random-bls                   Generate a random BLS key
  random-id                    Generate a random ID
  random-nodeid                Generate a random node ID
  revert-reason                Find revert reason for a failed tx hash
  rpc                          Ergonomic access to avalanche node RPC APIs
  verify-bls                   Verify a BLS Proof of Possession
  version                      Version
  vmid                         Given a vmName, try to encode the ASCII name as a vmID
  vmname                       Given a vmID, try to decode the ASCII name
  warp-aggregate-signatures    Aggregate signatures for a warp message
  warp-construct-l1-validator-registration Construct an unsigned L1ValidatorRegistration msg
  warp-construct-l1-weight     Construct an unsigned msg to change weight on P-Chain
  warp-construct-uptime        Construct an unsgined uptime message
  warp-get                     Get a warp message from a transaction ID
  warp-parse                   Parse a warp message
  xpub                         Show xpub for a BIP39 mnemonic and derivation path
  xpub-addrs                   Show addresses for an xpub key and derivation path
</code></pre>

## Mnemonics

Avalanche P-Chain and C-Chain use different address formats, and `ggt` provides utilities to help with this.

```sh
❯ ggt mnemonic-keys "test test test test test test test test test test test junk"
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

# Show various address encodings of a private key
❯ ggt pk PrivateKey-Fapb8hTUMABpZc9zWurmPR7why34LQNshss5RZHgCAgiY5n83

PrivKey Hex:   0x211cdc80c23ccc8eceab5d6903312391e656366a7a553e2c501b06add1729816
PrivKey CB58:  PrivateKey-Fapb8hTUMABpZc9zWurmPR7why34LQNshss5RZHgCAgiY5n83
Eth addr:      0x5a299B0010BAc9c0339B6EF600B1f2943131b1e7
Ava addr:      P-avax1yljhuvjkmtu0y5ls6kf4exsdd8gea9mp8jd32r
Ava addr:      P-fuji1yljhuvjkmtu0y5ls6kf4exsdd8gea9mptqfwxu
Ava addr:      P-local1yljhuvjkmtu0y5ls6kf4exsdd8gea9mp7pshft

❯ ggt pk 0x211cdc80c23ccc8eceab5d6903312391e656366a7a553e2c501b06add1729816

PrivKey Hex:   0x211cdc80c23ccc8eceab5d6903312391e656366a7a553e2c501b06add1729816
PrivKey CB58:  PrivateKey-Fapb8hTUMABpZc9zWurmPR7why34LQNshss5RZHgCAgiY5n83
Eth addr:      0x5a299B0010BAc9c0339B6EF600B1f2943131b1e7
Ava addr:      P-avax1yljhuvjkmtu0y5ls6kf4exsdd8gea9mp8jd32r
Ava addr:      P-fuji1yljhuvjkmtu0y5ls6kf4exsdd8gea9mptqfwxu
Ava addr:      P-local1yljhuvjkmtu0y5ls6kf4exsdd8gea9mp7pshft

```

## xpub

In Bitcoin, there is a concept of an `xpub` key (Extended Public Key), which allows you to derive all public keys (NOT private keys) from an HD mnemonic. GoGoTools lets you do this for evm and P-Chain addresses. Useful for services that need to monitor many addresses for activity, or gas funds, etc.

```sh
# Generate an xpub for a mnemonic
❯ ggt xpub "test test test test test test test test test test test junk"
xpub6Ce9NcJvTk36xtLSrJLZqE7wtgA5deCeYs7rSQtreh4cj6ByPtrg9sD7V2FNFLPnf8heNP3FGkeV9qwfzvZNSd54JoNXVsXFYSYwHsnJxqP

# Use the xpub to derive addresses
❯ ggt xpub-addrs xpub6Ce9NcJvTk36xtLSrJLZqE7wtgA5deCeYs7rSQtreh4cj6ByPtrg9sD7V2FNFLPnf8heNP3FGkeV9qwfzvZNSd54JoNXVsXFYSYwHsnJxqP --json | jq
{
  "xpub": "xpub6Ce9NcJvTk36xtLSrJLZqE7wtgA5deCeYs7rSQtreh4cj6ByPtrg9sD7V2FNFLPnf8heNP3FGkeV9qwfzvZNSd54JoNXVsXFYSYwHsnJxqP",
  "path": "m/44'/60'/0'/0/*",
  "evm_addrs": [
    "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
    "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
    "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
    "0x90F79bf6EB2c4f870365E785982E1f101E93b906",
    "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
    "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc",
    "0x976EA74026E726554dB657fA54763abd0C3a0aa9",
    "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955",
    "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f",
    "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720"
  ],
  "ava_addrs": [
    "P-avax15428vq2uzwhm3taey9sr9x5vm6tk78ewvkckhy",
    "P-avax1cjzphr67dug28rw9ueewrqllmxlqe5f09qqqnc",
    "P-avax1sj3m3zudqkw4plsvm7p08h9q89lsy4wp9c809f",
    "P-avax1y3fgntsaf2cmylsnu3880g7wp6aj73zun9t5wa",
    "P-avax1jft4f7x5cajyw6w8vfj9708mx9jfswq0ydvr0f",
    "P-avax1syp9y2mlfrrjwdaj5ytezh9wuerwuhnctaelxn",
    "P-avax1zyf8ga3p89xmjaxje2ckvtvhjzv7g24sgu0uqg",
    "P-avax1mer4xr0wqlq7v8m678vvx4zeyep20v8rwarugf",
    "P-avax1l5rrv44ufqq97psu2jqx909mdnjapwv0fwudyx",
    "P-avax14flvw0x8fstzly79tacgsulxvkpv858quqds6w"
  ]
}
```

## Warp

```sh
# Set up some env vars for convenience
export ETH_RPC_URL=https://api.avax.network/ext/bc/C/rpc
export AVA_RPC_URL=https://api.avax.network

# Get a warp message from a transaction ID
❯ ggt warp-get 0x7b0a220154bde7cb8e603dd5b88d418c55fadf10d09f6e43c399f94088f83bca

UnsignedMessage(NetworkID = 1, SourceChainID = 2q9e4r6Mu3U68nU1fYjgbR6JvwrRx36CohpAX5UQxse55x1Q5, Payload = 000000000001000000141424aef0d5272373beb69b2a860bd1da078df67f000000b60000000000010ad6355dc6b82cd375e3914badb3e2f8d907d0856f8e679b2db46f8938a2f01200000014000000000000000000000e8bd2300dfc723e53d38d543b279b9bd69c5b6754a09bce1eab2de2d9135eff7e391e42583fca4c19c6007e864971c2baba777dfa312ca7994e0000000067b3d5bf00000001000000016cc54e2d13e91e29867851238a8af6c53ca4a9bf00000001000000016cc54e2d13e91e29867851238a8af6c53ca4a9bf0000000000000002)

Payload (*message.RegisterL1Validator): {
  "subnetID": "5moznRzaAEhzWkNTQVdT1U4Kb9EU7dbsKZQNmHwtN5MGVQRyT",
  "nodeID": "NodeID-11111111116MAZMYBTpRxmWciKx4L",
  "blsPublicKey": "8d543b279b9bd69c5b6754a09bce1eab2de2d9135eff7e391e42583fca4c19c6007e864971c2baba777dfa312ca7994e",
  "expiry": 1739838911,
  "remainingBalanceOwner": {
    "threshold": 1,
    "addresses": [
      "P-avax1dnz5utgnay0znpnc2y3c4zhkc572f2dlkwapg8"
    ]
  },
  "disableOwner": {
    "threshold": 1,
    "addresses": [
      "P-avax1dnz5utgnay0znpnc2y3c4zhkc572f2dlkwapg8"
    ]
  },
  "weight": 2
}


# Parse a hex-encoded warp message
❯ ggt warp-parse 0x0000000000010000000000000000000000000000000000000000000000000000000000000000000000350000000000010000000000000027000000000002ffa2cb7c97396dc67fc06e9f4dbf03a667f671c5cadaf55d11251d3172d2662501

WarpMessage(UnsignedMessage(NetworkID = 1, SourceChainID = 11111111111111111111111111111111LpoYY, Payload = 0000000000010000000000000027000000000002ffa2cb7c97396dc67fc06e9f4dbf03a667f671c5cadaf55d11251d3172d2662501))

Payload (*message.L1ValidatorRegistration): {"validationID":"2wasqFU3CptbuuWgJg5awBJHar9bm1ANwp123rFSUQJ8txiQmJ","registered":true}
```

## Balances

```sh
❯ ggt balance 0x746189b6b6C2162C6BeAB83d4dB76f8C96A5381C
0.207649043134145176 ETH

❯ ggt balance avax19zfygxaf59stehzedhxjesads0p5jdvfeedal0
13.164163522 AVAX

# If ETH_FROM is set,
❯ ggt balance
0.207649043134145176 ETH


❯ export PRIVATE_KEY=010101...
❯ ggt balance-pk
0.2076490431 ETH   0x746189b6b6C2162C6BeAB83d4dB76f8C96A5381C
13.164163522 AVAX  avax19zfygxaf59stehzedhxjesads0p5jdvfeedal0



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
