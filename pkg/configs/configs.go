package configs

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
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
var AvaGenesisRaw []byte

// We need to munge the timestamps from the raw file
var AvaGenesis string

func init() {
	// difference between unlock schedule locktime and startime in original genesis
	const genesisLocktimeStartimeDelta = 2836800

	var (
		genesisMap map[string]interface{}
		err        error
	)

	if err = json.Unmarshal(AvaGenesisRaw, &genesisMap); err != nil {
		log.Fatalf("error reading ava-genesis.json %s", err)
	}

	startTime := time.Now().Unix()
	lockTime := startTime + genesisLocktimeStartimeDelta
	genesisMap["startTime"] = float64(startTime)
	allocations, ok := genesisMap["allocations"].([]interface{})
	if !ok {
		panic(errors.New("could not get allocations in genesis"))
	}
	for _, allocIntf := range allocations {
		alloc, ok := allocIntf.(map[string]interface{})
		if !ok {
			panic(fmt.Errorf("unexpected type for allocation in genesis. got %T", allocIntf))
		}
		unlockSchedule, ok := alloc["unlockSchedule"].([]interface{})
		if !ok {
			panic(errors.New("could not get unlockSchedule in allocation"))
		}
		for _, schedIntf := range unlockSchedule {
			sched, ok := schedIntf.(map[string]interface{})
			if !ok {
				panic(fmt.Errorf("unexpected type for unlockSchedule elem in genesis. got %T", schedIntf))
			}
			if _, ok := sched["locktime"]; ok {
				sched["locktime"] = float64(lockTime)
			}
		}
	}

	// now we can marshal the *whole* thing into bytes
	updatedGenesis, err := json.Marshal(genesisMap)
	if err != nil {
		panic(err)
	}

	AvaGenesis = string(updatedGenesis)
}
