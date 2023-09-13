package pkg

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

// 指定派生路径
const (
	EthDerivationPath = "m/44'/60'/0'/0/0"

	TronDerivationPath = "m/44'/195'/0'/0/0"
)

// GetMnemonic 生成助记词
func GetMnemonic(bitSize int) (string, error) {
	// 生成 256 位熵
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		panic(err)
	}

	// 使用熵生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func MustParseDerivationPath(path string) (accounts.DerivationPath, error) {
	parsed, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// AddressFromPrivateKey 从私钥生成地址
func AddressFromPrivateKey(privateKeyHex string) (*ecdsa.PublicKey, error) {
	// 如果 privateKeyHex 没有 "0x" 前缀, 添加一个
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] != "0x" {
		privateKeyHex = "0x" + privateKeyHex
	}

	// 从十六进制字符串解码私钥
	privateKeyBytes, err := hexutil.Decode(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 从字节创建 ECDSA 私钥
	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	publicKey := privateKeyECDSA.PublicKey
	// 返回地址
	return &publicKey, nil
}

// GetMnemonicBy256  生成助记词
func GetMnemonicBy256() (string, error) {
	mnemonic, err := GetMnemonic(256)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// GetMnemonicBy128  生成助记词
func GetMnemonicBy128() (string, error) {
	mnemonic, err := GetMnemonic(128)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}
