package configs

import (
	_ "embed"
)

const (
	NodeConfigFilename   = "node-config.json"
	CChainConfigFilename = "cchain-config.json"
	XChainConfigFilename = "xchain-config.json"
	AvaGenesisFilename   = "ava-genesis.json"
	ChainConfigFilename  = "config.json"
	AliasConfigFilename  = "aliases.json"
	BashScriptFilename   = "start.sh"
	AccountsFilename     = "accounts.json"
	ContractsFilename    = "contracts.json"
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

//go:embed ava-genesis.json
var AvaGenesis string
