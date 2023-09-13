package pkg

import (
	"crypto/ecdsa"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"sync"
)

type Wallet struct {
	mnemonic    string
	masterKey   *hdkeychain.ExtendedKey
	seed        []byte
	url         accounts.URL
	paths       map[common.Address]accounts.DerivationPath
	accounts    []accounts.Account
	stateLock   sync.RWMutex
	fixIssue172 bool
}

func newWallet(seed []byte, issue172 bool) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		masterKey:   masterKey,
		seed:        seed,
		accounts:    []accounts.Account{},
		paths:       map[common.Address]accounts.DerivationPath{},
		fixIssue172: issue172,
	}, nil
}

// EthNewFromMnemonic returns a new wallet from a BIP-39 mnemonic.
func EthNewFromMnemonic(mnemonic string) (*Wallet, error) {
	return NewFromMnemonic(mnemonic, false)
}

// TronNewFromMnemonic returns a new wallet from a BIP-39 mnemonic.
func TronNewFromMnemonic(mnemonic string) (*Wallet, error) {
	return NewFromMnemonic(mnemonic, true)
}

// NewFromMnemonic returns a new wallet from a BIP-39 mnemonic.
func NewFromMnemonic(mnemonic string, issue172 bool) (*Wallet, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("mnemonic is invalid")
	}

	seed, err := NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	wallet, err := newWallet(seed, issue172)
	if err != nil {
		return nil, err
	}
	wallet.mnemonic = mnemonic

	return wallet, nil
}

// TronNewFromSeed returns a new wallet from a BIP-39 seed.
func TronNewFromSeed(seed []byte) (*Wallet, error) {
	return newFromSeed(seed, false)
}

// EthNewFromSeed returns a new wallet from a BIP-39 seed.
func EthNewFromSeed(seed []byte) (*Wallet, error) {
	return newFromSeed(seed, false)
}

// newFromSeed returns a new wallet from a BIP-39 seed.
func newFromSeed(seed []byte, issue172 bool) (*Wallet, error) {
	if len(seed) == 0 {
		return nil, errors.New("seed is required")
	}

	return newWallet(seed, issue172)
}

// NewSeedFromMnemonic returns a BIP-39 seed based on a BIP-39 mnemonic.
func NewSeedFromMnemonic(mnemonic string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	return bip39.NewSeedWithErrorChecking(mnemonic, "")
}

// Derive implements accounts.Wallet, deriving a new account at the specific
// derivation path. If pin is set to true, the account will be added to the list
// of tracked accounts.
func (w *Wallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	// Try to derive the actual account and update its URL if successful
	w.stateLock.RLock() // Avoid device disappearing during derivation

	address, err := w.deriveAddress(path)

	w.stateLock.RUnlock()

	// If an error occurred or no pinning was requested, return
	if err != nil {
		return accounts.Account{}, err
	}

	account := accounts.Account{
		Address: address,
		URL: accounts.URL{
			Scheme: "",
			Path:   path.String(),
		},
	}

	if !pin {
		return account, nil
	}

	// Pinning needs to modify the state
	w.stateLock.Lock()
	defer w.stateLock.Unlock()

	if _, ok := w.paths[address]; !ok {
		w.accounts = append(w.accounts, account)
		w.paths[address] = path
	}

	return account, nil
}

// DeriveAddress derives the account address of the derivation path.
func (w *Wallet) deriveAddress(path accounts.DerivationPath) (common.Address, error) {
	publicKeyECDSA, err := w.derivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

// DerivePublicKey derives the public key of the derivation path.
func (w *Wallet) derivePublicKey(path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	privateKeyECDSA, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}

// DerivePrivateKey derives the private key of the derivation path.
func (w *Wallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	var err error
	key := w.masterKey
	for _, n := range path {
		if w.fixIssue172 && key.IsAffectedByIssue172() {
			key, err = key.Derive(n)
		} else {
			key, err = key.DeriveNonStandard(n)
		}
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}

	return privateKeyECDSA, nil
}

// PrivateKeyHex return the ECDSA private key in hex string format of the account.
func (w *Wallet) PrivateKeyHex(account accounts.Account) (string, error) {
	privateKeyBytes, err := w.PrivateKeyBytes(account)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(privateKeyBytes)[2:], nil
}

// PrivateKeyBytes returns the ECDSA private key in bytes format of the account.
func (w *Wallet) PrivateKeyBytes(account accounts.Account) ([]byte, error) {
	privateKey, err := w.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	return crypto.FromECDSA(privateKey), nil
}

// PrivateKey returns the ECDSA private key of the account.
func (w *Wallet) PrivateKey(account accounts.Account) (*ecdsa.PrivateKey, error) {
	path, err := ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}

	return w.derivePrivateKey(path)
}

// ParseDerivationPath parses the derivation path in string format into []uint32
func ParseDerivationPath(path string) (accounts.DerivationPath, error) {
	return accounts.ParseDerivationPath(path)
}
