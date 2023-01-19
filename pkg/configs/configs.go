package configs

import (
	_ "embed"
)

//go:embed node-config.json
var NodeConfig string

//go:embed genesis-subnetevm.json
var genesisSubnetEVM string
