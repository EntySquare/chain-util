//go:build windows
// +build windows

package keystore

import (
	"fmt"
	"github.com/EntySquare/chain-util/pkg/tron"
)

func RecoverPubkey(hash []byte, signature []byte) (tron.Address, error) {
	return nil, fmt.Errorf("not implemented")
}
