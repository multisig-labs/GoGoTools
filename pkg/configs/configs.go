package configs

import (
	_ "embed"
)

//go:embed accounts.json
var Accounts string

//go:embed contracts.json
var Contracts string

//go:embed node-config.json
var NodeConfig string

//go:embed subnetevm-genesis.json
var SubnetEVMGenesis string

//go:embed subnetevm-config.json
var SubnetEVMConfig string

//go:embed cchain-config.json
var CChainConfig string

//go:embed xchain-config.json
var XChainConfig string

//go:embed README.md
var Readme string

//go:embed start.sh
var StartBash string
