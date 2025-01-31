package validatormanagercmd

import (
	"encoding/json"

	"github.com/ava-labs/avalanchego/ids"
	evmtypes "github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/multisig-labs/gogotools/pkg/glacier"
)

type L1Config struct {
	Name           string         `json:"name"`
	Network        string         `json:"network"`
	SubnetID       ids.ID         `json:"subnet_id"`
	BlockchainID   ids.ID         `json:"blockchain_id"`
	VMSubnetID     ids.ID         `json:"vm_subnet_id"`
	VMBlockchainID ids.ID         `json:"vm_blockchain_id"`
	VMAddress      common.Address `json:"vm_address"`
	ValidatorsURL  string         `json:"validators_url"`
	AvaURL         string         `json:"ava_url"`
	EvmURL         string         `json:"evm_url"`
}

type Info struct {
	Config             L1Config                    `json:"config"`
	PchainValidators   []glacier.Validator         `json:"pchain_validators"`
	ContractValidators []Validator                 `json:"contract_validators"`
	Uptime             []evmtypes.CurrentValidator `json:"uptime"`
	ContractSettings   ContractSettings            `json:"contract_settings"`
}

// Maps to the ValidatorManager contract
type Validator struct {
	Status         uint8
	NodeID         []byte
	StartingWeight uint64
	MessageNonce   uint64
	Weight         uint64
	StartedAt      uint64
	EndedAt        uint64
	ValidationID   ids.ID
}

// Do the dance to get NodeID to output as the correct format
func (v *Validator) MarshalJSON() ([]byte, error) {
	nodeID, err := ids.ToNodeID(v.NodeID)
	if err != nil {
		return nil, err
	}
	sj := struct {
		Status         uint8
		NodeID         string
		StartingWeight uint64
		MessageNonce   uint64
		Weight         uint64
		StartedAt      uint64
		EndedAt        uint64
		ValidationID   ids.ID
	}{
		Status:         v.Status,
		NodeID:         nodeID.String(),
		StartingWeight: v.StartingWeight,
		MessageNonce:   v.MessageNonce,
		Weight:         v.Weight,
		StartedAt:      v.StartedAt,
		EndedAt:        v.EndedAt,
		ValidationID:   v.ValidationID,
	}
	return json.Marshal(sj)
}
