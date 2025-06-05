package hd

import (
	"crypto/ecdsa"
	"fmt"

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

type HDPubKey struct {
	PubKey *ecdsa.PublicKey
	Path   string
}

func (h HDPubKey) EthAddr() string {
	return ethcrypto.PubkeyToAddress(*h.PubKey).String()
}

func (h HDPubKey) AvaAddr(chain string, hrp string) string {
	pkbytes := ethcrypto.CompressPubkey(h.PubKey)
	avapk, err := avacrypto.ToPublicKey(pkbytes)
	if err != nil {
		panic(err)
	}
	addr, err := address.Format(chain, hrp, avapk.Address().Bytes())
	if err != nil {
		panic(err)
	}
	return addr
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
				return // Skip failed derivations silently
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

// Get the Extended Public Key (xpubkey) for the given mnemonic and derivation path
func XPubKey(mnemonic string, path accounts.DerivationPath) (*hdkeychain.ExtendedKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	// Derive to the account level using private key (for hardened derivations)
	// For Ethereum: m/44'/60'/0'
	accountKey := masterKey
	for i := range 3 { // First 3 components are hardened
		accountKey, err = accountKey.Derive(path[i])
		if err != nil {
			return nil, err
		}
	}

	// Convert to public key at account level
	accountPubKey, err := accountKey.Neuter()
	if err != nil {
		return nil, err
	}

	return accountPubKey, nil
}

func DerivePubKeys(xpub *hdkeychain.ExtendedKey, path accounts.DerivationPath, numKeys int) ([]HDPubKey, error) {
	hdkeys := []HDPubKey{}

	// Only derive the non-hardened part of the path (last 2 components for standard paths)
	// For Ethereum: 0/0, 0/1, 0/2, etc.
	nonHardenedPath := path[3:] // Skip the first 3 hardened components

	var derive = func(limit int, next func() accounts.DerivationPath) {
		for i := range limit {
			// Create path for this index: 0/i
			currentPath := make(accounts.DerivationPath, len(nonHardenedPath))
			copy(currentPath, nonHardenedPath)
			currentPath[len(currentPath)-1] = uint32(i) // Set the address index

			if pk, err := derivePublicKey(xpub, currentPath); err != nil {
				return // Skip failed derivations silently
			} else {
				// Reconstruct full path for display
				fullPath := make(accounts.DerivationPath, len(path))
				copy(fullPath, path)
				fullPath[len(fullPath)-1] = uint32(i)

				hdk := HDPubKey{
					PubKey: pk,
					Path:   fullPath.String(),
				}
				hdkeys = append(hdkeys, hdk)
			}
		}
	}

	derive(numKeys, func() accounts.DerivationPath { return nonHardenedPath })
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
	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key at index %d: %w", n, err)
		}
	}

	pubkey, err := key.ECPubKey()
	if err != nil {
		return nil, err
	}
	return pubkey.ToECDSA(), nil
}
