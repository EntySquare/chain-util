package keys_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/EntySquare/chain-util/pkg/tron/keys"
	"github.com/stretchr/testify/assert"
)

func Test_mnemonic_to_pk(t *testing.T) {
	// Hardcoded index of 0 for brandnew account.
	private, _ := keys.FromMnemonicSeedAndPassphrase("bag car educate river behind lumber fee seminar spend air stuff phrase mango basket fine crystal number strong eight what spawn impact crater surprise", "", 0)
	pk_bytes := private.Serialize()
	fmt.Println("private=", private)
	fmt.Println("pk_bytes=", pk_bytes)

	println("Privatekey: ", hex.EncodeToString(pk_bytes))
	assert.Equal(t, hex.EncodeToString(pk_bytes), "b5a4cea271ff424d7c31dc12a3e43e401df7a40d7412a15750f3f0b6b5449a28")
}
