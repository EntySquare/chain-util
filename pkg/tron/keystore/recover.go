//go:build !windows
// +build !windows

package keystore

import (
	"github.com/EntySquare/chain-util/pkg/tron"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func RecoverPubkey(hash []byte, signature []byte) (tron.Address, error) {

	if signature[64] >= 27 {
		signature[64] -= 27
	}

	sigPublicKey, err := secp256k1.RecoverPubkey(hash, signature)
	if err != nil {
		return nil, err
	}
	pubKey, err := UnmarshalPublic(sigPublicKey)
	if err != nil {
		return nil, err
	}

	addr := tron.PublicKeyToAddress(*pubKey)
	return addr, nil
}
