package utils

import (
	"crypto/rand"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/signer"
)

func RandomID() (ids.ID, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ids.ID{}, err
	}
	i, _ := ids.ToID(b)
	return i, nil
}

func RandomNodeID() (ids.NodeID, error) {
	r := make([]byte, 20)
	_, err := rand.Read(r)
	if err != nil {
		return ids.NodeID{}, err
	}
	zeroSlice := make([]byte, 10)
	copy(r, zeroSlice)
	nodeid := ids.NodeID(r)
	return nodeid, nil
}

func RandomBLS() (*bls.LocalSigner, *signer.ProofOfPossession, error) {
	pop := &signer.ProofOfPossession{}
	sk, err := bls.NewSigner()
	if err != nil {
		return nil, nil, err
	}
	pop = signer.NewProofOfPossession(sk)
	err = pop.Verify()
	if err != nil {
		return nil, nil, err
	}

	return sk, pop, nil
}
