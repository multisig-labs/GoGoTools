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
	f := avacrypto.Factory{}
	pkbytes := ethcrypto.FromECDSA(h.PK)
	avapk, _ := f.ToPrivateKey(pkbytes)
	addr, _ := address.Format(chain, hrp, avapk.PublicKey().Address().Bytes())
	return addr
}

func (h HDKey) EthPrivKey() string {
	pkb := ethcrypto.FromECDSA(h.PK)
	return common.Bytes2Hex(pkb)
}

func (h HDKey) AvaPrivKey() string {
	f := avacrypto.Factory{}
	pkbytes := ethcrypto.FromECDSA(h.PK)
	avapk, _ := f.ToPrivateKey(pkbytes)
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

// func forEth(pk *ecdsa.PrivateKey) (privKeyHex string, addr string) {
// 	pkb := ethcrypto.FromECDSA(pk)
// 	privKeyHex = common.Bytes2Hex(pkb)
// 	addr = ethcrypto.PubkeyToAddress(pk.PublicKey).String()
// 	return
// }

// func forAva(pk *ecdsa.PrivateKey) (privKeyHex string, addr string) {
// 	f := avacrypto.Factory{}
// 	pkbytes := ethcrypto.FromECDSA(pk)
// 	avapk, _ := f.ToPrivateKey(pkbytes)
// 	privKeyHex = avapk.String()
// 	addr, _ = address.Format("P", "avax", avapk.PublicKey().Address().Bytes())
// 	return
// }

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
	// publicKey, _ := key.ECPubKey()
	privateKeyECDSA := privateKey.ToECDSA()
	return privateKeyECDSA, nil
}

// pk, adr := forEth(privateKeyECDSA)
// fmt.Printf("%s %s\n\n", adr, pk)

// pk, adr = forAva(privateKeyECDSA)
// fmt.Printf("%s %s\n\n", adr, pk)
// publicKeyECDSA := publicKey.ToECDSA()

// x := ethcrypto.FromECDSA(privateKeyECDSA)
// fmt.Printf("pk: %x", x)
// addr := ethcrypto.PubkeyToAddress(*publicKeyECDSA)
// fmt.Printf("addr: %x", addr)

// f := avacrypto.Factory{}

// ethcrypto.FromECDSAPub(publicKeyECDSA)
// pkbytes := ethcrypto.FromECDSA(privateKeyECDSA)

// avapk, err := f.ToPrivateKey(pkbytes)
// if err != nil {
// 	fmt.Printf("err %s", err)
// }
// fmt.Printf("avapk %x", avapk)
// a, err := address.Format("P", "avax", avapk.PublicKey().Address().Bytes())
// if err != nil {
// 	fmt.Printf("err %s", err)
// }
// fmt.Printf("a %s", a)
