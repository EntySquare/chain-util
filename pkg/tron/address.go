package tron

import (
	"bytes"
	"crypto/ecdsa"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"github.com/EntySquare/chain-util/pkg"

	"github.com/EntySquare/chain-util/pkg/tron/common"
	ethereumCommon "github.com/ethereum/go-ethereum/common"
	"github.com/tyler-smith/go-bip39"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 21
	// AddressLengthBase58 is the expected length of the address in base58format
	AddressLengthBase58 = 34
	// TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(0x41)
)

// Address represents the 21 byte address of an Tron account.
type Address []byte

// Bytes get bytes from address
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex get bytes from address in string
func (a Address) Hex() string {
	return common.ToHex(a[:])
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address {
	id := b.Bytes()
	base := bytes.Repeat([]byte{0}, AddressLength-len(id))
	return append(base, id...)
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address {
	addr, err := common.FromHex(s)
	if err != nil {
		return nil
	}
	return addr
}

// Base58ToAddress returns Address with byte values of s.
func Base58ToAddress(s string) (Address, error) {
	addr, err := common.DecodeCheck(s)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// Base64ToAddress returns Address with byte values of s.
func Base64ToAddress(s string) (Address, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return Address(decoded), nil
}

// String implements fmt.Stringer.
func (a Address) String() string {
	if len(a) == 0 {
		return ""
	}

	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}
	return common.EncodeCheck(a.Bytes())
}

// PublicKeyToAddress returns address from ecdsa public key
func PublicKeyToAddress(p ecdsa.PublicKey) Address {
	address := crypto.PubkeyToAddress(p)

	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return addressTron
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	*a = Address(srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return []byte(a), nil
}

// AddressFromPrivateKey 从私钥生成地址
func AddressFromPrivateKey(privateKeyHex string) (string, error) {
	publicKey, err := pkg.AddressFromPrivateKey(privateKeyHex)
	if err != nil {
		return "", err
	}
	tronAddress := PublicKeyToAddress(*publicKey)
	// 返回地址
	return tronAddress.String(), nil
}

// CreateWallet Tron 创建钱包 生成地址和私钥 助记词
func CreateWallet() (address string, privateKey string, mnemonic string, err error) {

	mnemonic, err = pkg.GetMnemonicBy256()
	if err != nil {
		return "", "", "", err
	}

	address, privateKey, err = GenerateAddressFromMnemonic(mnemonic)
	if err != nil {
		return "", "", "", err
	}

	return address, privateKey, mnemonic, nil
}

func IsValidAddress(address string) bool {
	// 定义 TRON 地址的正则表达式
	tronAddressPattern := regexp.MustCompile("^T[0-9A-HJ-NP-Za-km-z]{33}$")
	return tronAddressPattern.MatchString(address)
}

// GenerateAddressFromMnemonic Tron 根据助记词生成地址和私钥
func GenerateAddressFromMnemonic(mnemonic string) (address string, privateKey string, err error) {
	// 生成私钥  这里可以选择传入指定密码或者空字符串，不同密码生成的助记词不同
	seed := bip39.NewSeed(mnemonic, "")
	// 使用种子生成 BIP32 主密钥
	wallet, err := pkg.TronNewFromSeed(seed)
	if err != nil {
		return "", "", err
	}
	path, err := pkg.MustParseDerivationPath(pkg.TronDerivationPath) //最后一位是同一个助记词的地址id，从0开始，相同助记词可以生产无限个地址
	if err != nil {
		return "", "", err
	}
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", err
	}
	privateKey, err = wallet.PrivateKeyHex(account)
	if err != nil {
		return "", "", err
	}
	tronAddress := AddressToTronAddress(account.Address)
	return tronAddress.String(), privateKey, nil
}

// AddressToTronAddress returns address from ecdsa public key
func AddressToTronAddress(address ethereumCommon.Address) Address {
	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return addressTron
}
