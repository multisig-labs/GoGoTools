<h1 align="center">GoGoTools 🎈</h1>
<p align="center">A (growing) collection of useful tools for Avalanche subnet developers.</p>

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

We are still trying to find the optimal workflows for doing this kind of dev work, but this is what we have as of now. Help wanted!

**Assumptions:**

- You have [Foundry](https://book.getfoundry.sh/getting-started/installation) installed and `cast` in your path
- You have a version of the `avalanchego` binary
- You have cloned the [subnet-evm](https://github.com/ava-labs/subnet-evm) repo and have a compiled evm binary
- You have this tool `ggt` compiled and in your $PATH

### Workflow

```sh
mkdir MySubnet
cd MySubnet
ggt utils init
ggt node prepare nodeV1 --ava-bin=/full/path/to/avalanchego --vm-name=subnetevm --vm-bin=/full/path/to/subnetevm
```

At this point you will have a directory called `nodeV1` which looks like this:

```
nodeV1
├── bin
│   ├── avalanchego -> /full/path/to/avalanchego
│   └── plugins
│       └── srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy -> /full/path/to/subnetevm
├── configs
│   ├── chains
│   │   ├── C
│   │   │   └── config.json
│   │   └── aliases.json
│   ├── node-config.json
│   └── vms
│       └── aliases.json
└── data
```

(Note that the `--vm-name=subnetevm` name you supplied for your VM (which can be any name) has been converted into an Avalanche `ids.ID` `srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy` and symlinked to your evm binary)

The default `node-config.json` configures `avalanchego` to be a single node with no staking.

```json
{
  "network-id": "local",
  "staking-enabled": false,
  "staking-ephemeral-cert-enabled": true,
  "staking-ephemeral-signer-enabled": true,
  "index-enabled": true,
  "api-keystore-enabled": true,
  "api-admin-enabled": true,
  "log-rotater-max-files": 1,
  "log-rotater-max-size": 1
}
```

Now we can start our node:

```sh
ggt node run nodeV1
```

This will start `avalanchego` from the `nodeV1` directory and you should see the `nodeV1/data` directory fill up with logs and data.

In another terminal, lets create our subnet (the `ggt utils init` cmd creates a sample genesis with all precompiles enabled):

```sh
ggt wallet create-chain MyChain subnetevm genesis.json
```

This will create a Subnet, and then inside that new subnet it will create a blockchain with the name `MyChain`, using the `subnetevm` virtual machine we registered earlier, and using the specified genesis file.

You should see an RPC URL printed to the terminal:

`http://localhost:9650/ext/bc/6SPgMtm5xfZrGGLJztaByMwKGJhrw4WzhKk6nGC5yfXqiJGuT/rpc`

You can now use this to issue commands to your EVM:

## Info

```sh
ggt node info | jq
```

(Assuming you have [jq](https://stedolan.github.io/jq/manual/) installed, and why wouldn't you!)

This collects node info from several different Avalanche API endpoints and gives you a single blob with all the data, including an `rpcs` key with the rpc url for each blockchain.

```json
{
  "nodeID": "NodeID-5VvkgcSJoUnRruSEMm9uR7P4Wxr1hQi1q",
  "networkID": 12345,
  "networkName": "local",
  "uptime": {
    "rewardingStakePercentage": "100.0000",
    "weightedAveragePercentage": "100.0000"
  },
  "getNodeVersion": {
    "version": "avalanche/1.9.7",
    "databaseVersion": "v1.4.5",
    "rpcProtocolVersion": "22",
    "gitCommit": "3e3e40f2f4658183d999807b724245023a13f5dc",
    "vmVersions": {
      "avm": "v1.9.7",
      "evm": "v0.11.5",
      "platform": "v1.9.7",
      "subnetevm": "v0.4.8@880ec774bf5746c6c6aceb6887d08b221ed565cd"
    }
  },
  "getVMs": {
    "vms": {
      "jvYyfQTxGMJLuGWa55kdP2p2zSUYsQ5Raupu4TW34ZAUBAbtq": ["avm"],
      "mgj786NP7uDwBCcq6YwThhaN8FLyybkCa4zBWTQbNgmK6k9A6": ["evm"],
      "qd2U4HDWUvMrVUeTcCHp6xH3Qpnn1XbU5MDdnBoiifFqvgXwT": ["nftfx"],
      "rWhpuQPF1kb72esV2momhMuTYGkEb1oL29pt2EBXWmSy4kxnT": ["platform"],
      "rXJsCSEYXg2TehWxCEEGj6JU2PWKTkd6cBdNLjoe2SpsKD9cy": ["propertyfx"],
      "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ": ["secp256k1fx"],
      "srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy": ["subnetevm"]
    }
  },
  "subnets": [
    {
      "id": "29uVeLPJB1eQJkzRemU8g8wZDw5uJRqpab5U2mX9euieVwiEbL",
      "controlKeys": ["P-local18jma8ppw3nhx5r4ap8clazz0dps7rv5u00z96u"],
      "threshold": "1"
    },
    {
      "id": "11111111111111111111111111111111LpoYY",
      "controlKeys": [],
      "threshold": "0"
    }
  ],
  "blockchains": [
    {
      "id": "SRq2ZdVwqyQcQqVwTtjZPTDttDWKTiUEg2vyy3AeobBjeS3z3",
      "name": "MyChain",
      "subnetID": "29uVeLPJB1eQJkzRemU8g8wZDw5uJRqpab5U2mX9euieVwiEbL",
      "vmID": "srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy"
    },
    {
      "id": "2CA6j5zYzasynPsFeNoqWkmTCt3VScMvXUZHbfDJ8k3oGzAPtU",
      "name": "C-Chain",
      "subnetID": "11111111111111111111111111111111LpoYY",
      "vmID": "mgj786NP7uDwBCcq6YwThhaN8FLyybkCa4zBWTQbNgmK6k9A6"
    },
    {
      "id": "2eNy1mUFdmaxXNj1eQHUe7Np4gju9sJsEtWQ4MX3ToiNKuADed",
      "name": "X-Chain",
      "subnetID": "11111111111111111111111111111111LpoYY",
      "vmID": "jvYyfQTxGMJLuGWa55kdP2p2zSUYsQ5Raupu4TW34ZAUBAbtq"
    }
  ],
  "aliases": {
    "blockchainAliases": {
      "SRq2ZdVwqyQcQqVwTtjZPTDttDWKTiUEg2vyy3AeobBjeS3z3": [
        "SRq2ZdVwqyQcQqVwTtjZPTDttDWKTiUEg2vyy3AeobBjeS3z3"
      ],
      "2CA6j5zYzasynPsFeNoqWkmTCt3VScMvXUZHbfDJ8k3oGzAPtU": [
        "C",
        "evm",
        "2CA6j5zYzasynPsFeNoqWkmTCt3VScMvXUZHbfDJ8k3oGzAPtU"
      ],
      "2eNy1mUFdmaxXNj1eQHUe7Np4gju9sJsEtWQ4MX3ToiNKuADed": [
        "X",
        "avm",
        "2eNy1mUFdmaxXNj1eQHUe7Np4gju9sJsEtWQ4MX3ToiNKuADed"
      ]
    }
  },
  "rpcs": {
    "MyChain": "http://localhost:9650/ext/bc/SRq2ZdVwqyQcQqVwTtjZPTDttDWKTiUEg2vyy3AeobBjeS3z3/rpc"
  }
}
```

This makes it easy to do something like this:

```sh
export ETH_RPC_URL=`ggt node info | jq -r '.rpcs.MyChain'`
cast call 0x0000000000000000000000000000000000000000 `cast sig "getCurrentBlockNumber()"`
```

## Subnet EVM Precompiles

The [Subnet-EVM](https://github.com/ava-labs/subnet-evm) repo has some nice example contracts you can use to interact with the default precompiles.

However, in the interest of getting as close to the metal as possible, to really understand how things are working, `ggt` has some convenience commands that wrap the (amazing!) `cast` command from Foundry. The `ggt utils init` command creates default `accounts.json` and `contracts.json` files, that you can modify with your particular info, and we use these to make issuing `cast` commands a little more ergonomic by using those files to resolve user and contract addresses.

Assuming you have your node running, and your `ETH_RPC_URL` pointing to it, you can do things like this:

```sh
# Balances of users in accounts.json
ggt cast balances | jq
# Send eth from one user to another
ggt cast send-eth owner alice 1ether | jq
# Call read-only contract methods
ggt cast call owner TxAllowList "readAllowList(address)" bob
ggt cast call owner FeeConfigManager "getFeeConfigLastChangedAt()"
# Send a signed tx to a contract / method
ggt cast send owner NativeMinter "mintNativeCoin(address,uint256)" alice 1ether | jq
ggt cast send owner TxAllowList "setEnabled(address)" bob | jq
ggt cast send owner TxAllowList "setNone(address)" bob | jq
```

Cast has tools to decode the output of a contract call, so for example to see the current fee configuration via the precompile we can do this:

```sh
export DATA=$(ggt cast call owner FeeConfigManager "getFeeConfig()")
cast --abi-decode "getFeeConfig()(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)" $DATA

8000000
2
25000000000
50000000
36
0
1000000
200000
```

Which maps to the struct returned by the precompile

```sol
function getFeeConfig()
    external
    view
    returns (
      uint256 gasLimit,
      uint256 targetBlockRate,
      uint256 minBaseFee,
      uint256 targetGas,
      uint256 baseFeeChangeDenominator,
      uint256 minBlockGasCost,
      uint256 maxBlockGasCost,
      uint256 blockGasCostStep
    );
```

<hr />

# 🚀 LFGG 🚀

This is version 0.01 of this tool, and we plan on making it so useful that it will be a core part of every Avalanche developer's toolbox. We welcome any and all idea / contributions on how to make the experience better.
