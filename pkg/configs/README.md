# GoGoTools

## Commands Scratchpad

Some example of the kinds of commands you might run during a dev workflow.

```sh
ggt utils init
rm -rf NodeV1
ggt node prepare NodeV1 --ava-bin=$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego --vm-name=subnetevm --vm-bin=$GOPATH/src/github.com/ava-labs/subnet-evm/build/subnet-evm
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
