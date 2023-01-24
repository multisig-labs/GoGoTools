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

//go:embed genesis-subnetevm.json
var GenesisSubnetEVM string

//go:embed coreth-config.json
var CorethConfig string
