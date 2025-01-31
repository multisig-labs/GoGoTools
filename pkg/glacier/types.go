package glacier

import "github.com/ava-labs/avalanchego/ids"

type Pageable interface {
	NextPage() string
}

type Validator struct {
	ValidationID          ids.ID         `json:"validationId"`
	CreationTimestamp     uint64         `json:"creationTimestamp"`
	NodeID                ids.NodeID     `json:"nodeId"`
	SubnetID              ids.ID         `json:"subnetId"`
	Weight                uint64         `json:"weight"`
	RemainingBalance      uint64         `json:"remainingBalance"`
	BlsCredentials        BlsCredentials `json:"blsCredentials"`
	RemainingBalanceOwner Owner          `json:"remainingBalanceOwner"`
	DeactivationOwner     Owner          `json:"deactivationOwner"`
}

type BlsCredentials struct {
	PublicKey         string `json:"publicKey"`
	ProofOfPossession string `json:"proofOfPossession"`
}

type Owner struct {
	Addresses []string `json:"addresses"`
	Threshold uint64   `json:"threshold"`
}

type ListValidators struct {
	Validators    []Validator `json:"validators"`
	NextPageToken string      `json:"nextPageToken"`
}

func (r ListValidators) NextPage() string {
	return r.NextPageToken
}
