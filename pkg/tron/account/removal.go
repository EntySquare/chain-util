package account

import (
	"fmt"
	"github.com/EntySquare/chain-util/pkg/tron/common"
	"github.com/EntySquare/chain-util/pkg/tron/store"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

// RemoveAccount - removes an account from the keystore
func RemoveAccount(name string) error {
	accountExists := store.DoesNamedAccountExist(name)

	if !accountExists {
		return fmt.Errorf("account %s doesn't exist", name)
	}

	uDir, _ := homedir.Dir()
	tronCTLDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	accountDir := fmt.Sprintf("%s/%s", tronCTLDir, name)
	os.RemoveAll(accountDir)

	return nil
}
