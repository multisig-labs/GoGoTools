# GoGoTools

See [GoGoTools](https://github.com/multisig-labs/GoGoTools) repo for more info.

This directory was initialized with `ggt utils init` which created some default config files for you. Feel free to change them or leave them as default. They will be copied into the right place in each node directory you create with `ggt node prepare`.

## Genesis

Hardhat and other tools use a test mnemonic to seed accounts on their EVMs, so we have taken the first key and added it to the Cchain and SubnetEVM genesis files with an initial balance, to reduce friction when developing.

Mnemonic: test test test test test test test test test test test junk
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

## Commands Scratchpad

Some example of the kinds of commands you might run during a dev workflow.

```sh
ggt utils init
ggt init v1.9.7 v0.4.8  # This will download avalanchego and subnet-evm from GitHub
# Prepare a node with just avalanchego and no custom VMs
ggt node prepare NodeV1 --ava-bin=$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego
# Prepare a node with avalanchego and a custom VM
ggt node prepare NodeV1 --ava-bin=$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego --vm-name=subnetevm --vm-bin=$GOPATH/src/github.com/ava-labs/avalanchego/build/plugins/srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy

rm -rf NodeV1 # Remove a node directory

ggt node run NodeV1 --clear-logs

ggt wallet create-chain NodeV1 MyChain subnetevm

ggt node info | jq
ggt node health | jq
ggt node explorer MyChain

export ETH_RPC_URL=`ggt node info | jq -r '.rpcs.MyChain'`
echo $ETH_RPC_URL
cast chain-id

ggt cast balances | jq
ggt cast send-eth owner alice 1ether | jq
ggt cast send owner NativeMinter "mintNativeCoin(address,uint256)" alice 1ether | jq
ggt cast call owner TxAllowList "readAllowList(address)" bob
ggt cast send owner TxAllowList "setEnabled(address)" bob | jq
ggt cast call owner FeeConfigManager "getFeeConfigLastChangedAt()"
ggt cast call owner FeeConfigManager "getFeeConfig()" | xargs cast --abi-decode "getFeeConfig()(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)"
```
