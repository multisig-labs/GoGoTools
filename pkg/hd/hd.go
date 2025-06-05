package hd

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	avacrypto "github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

var EthDerivationPath = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}
var AvaDerivationPath = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 9000, 0x80000000 + 0, 0, 0}

type HDKey struct {
	PK   *ecdsa.PrivateKey
	Path string
}

func (h HDKey) EthAddr() string {
	return ethcrypto.PubkeyToAddress(h.PK.PublicKey).String()
}

func (h HDKey) AvaAddr(chain string, hrp string) string {
	pkbytes := ethcrypto.FromECDSA(h.PK)
	avapk, _ := avacrypto.ToPrivateKey(pkbytes)
	addr, _ := address.Format(chain, hrp, avapk.PublicKey().Address().Bytes())
	return addr
}

func (h HDKey) EthPrivKey() string {
	pkb := ethcrypto.FromECDSA(h.PK)
	return common.Bytes2Hex(pkb)
}

func (h HDKey) AvaPrivKey() string {
	pkbytes := ethcrypto.FromECDSA(h.PK)
	avapk, _ := avacrypto.ToPrivateKey(pkbytes)
	return avapk.String()
}

func DeriveHDKeys(mnemonic string, path accounts.DerivationPath, numKeys int) ([]HDKey, error) {
	// Generate seed from the mnemonic
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	// Generate master key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	hdkeys := []HDKey{}

	var derive = func(limit int, next func() accounts.DerivationPath) {
		for i := 0; i < limit; i++ {
			path := next()
			if pk, err := derivePrivateKey(masterKey, path); err != nil {
				fmt.Println("Account derivation failed", "error", err)
			} else {
				hdk := HDKey{
					PK:   pk,
					Path: path.String(),
				}
				hdkeys = append(hdkeys, hdk)
			}
		}
	}

	derive(numKeys, accounts.DefaultIterator(path))

	return hdkeys, nil
}

func XPubKey(mnemonic string, path accounts.DerivationPath) (*hdkeychain.ExtendedKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	masterPubKey, err := masterKey.Neuter()
	if err != nil {
		return nil, err
	}

	return masterPubKey, nil
}

func DerivePubKeys(xpub *hdkeychain.ExtendedKey, path accounts.DerivationPath, numKeys int) ([]HDKey, error) {
	hdkeys := []HDKey{}
	var derive = func(limit int, next func() accounts.DerivationPath) {
		for i := 0; i < limit; i++ {
			path := next()
			if pk, err := derivePublicKey(xpub, path); err != nil {
				fmt.Println("Account derivation failed", "error", err)
			} else {
				fmt.Println("PubKey", "path", path, "pk", ethcrypto.PubkeyToAddress(*pk).String())
				hdk := HDKey{
					PK: &ecdsa.PrivateKey{
						PublicKey: *pk,
						D:         big.NewInt(0),
					},
					Path: path.String(),
				}
				hdkeys = append(hdkeys, hdk)
			}
		}
	}

	derive(numKeys, accounts.DefaultIterator(path))
	return hdkeys, nil
}

func derivePrivateKey(masterKey *hdkeychain.ExtendedKey, path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	var err error
	key := masterKey
	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, _ := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	return privateKeyECDSA, nil
}

func derivePublicKey(masterKey *hdkeychain.ExtendedKey, path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	var err error
	key := masterKey
	fmt.Println("Begin deriving", "path", path)
	for _, n := range path {
		fmt.Println("Deriving", "path", path, "n", n)

		tmpKey, err := key.Derive(n)
		if err != nil {
			fmt.Println("Derive failed", "error", err)
		} else {
			key = tmpKey
		}
	}
	if key == nil {
		return nil, fmt.Errorf("key is nil")
	}
	pubkey, err := key.ECPubKey()
	if err != nil {
		return nil, err
	}
	return pubkey.ToECDSA(), nil
}
