package account

import (
	"github.com/EntySquare/chain-util/pkg"
	"github.com/EntySquare/chain-util/pkg/tron/keys"
	"github.com/EntySquare/chain-util/pkg/tron/store"
)

// Creation struct for account
type Creation struct {
	Name               string
	Passphrase         string
	Mnemonic           string
	MnemonicPassphrase string
	HdAccountNumber    *uint32
	HdIndexNumber      *uint32
}

// CreateNewLocalAccount assumes all the inputs are valid, legitmate
func CreateNewLocalAccount(candidate *Creation) error {
	ks := store.FromAccountName(candidate.Name)
	if candidate.Mnemonic == "" {
		candidate.Mnemonic, _ = pkg.GetMnemonicBy256()
	}
	// Hardcoded index of 0 for brandnew account.
	private, _ := keys.FromMnemonicSeedAndPassphrase(candidate.Mnemonic, candidate.MnemonicPassphrase, 0)
	_, err := ks.ImportECDSA(private.ToECDSA(), candidate.Passphrase)
	if err != nil {
		return err
	}
	return nil
}
