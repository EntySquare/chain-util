package pkg

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"regexp"
)

// EthAddressFromPrivateKey 从私钥生成地址
func EthAddressFromPrivateKey(privateKeyHex string) (string, error) {
	publicKey, err := AddressFromPrivateKey(privateKeyHex)
	if err != nil {
		return "", err
	}
	address := crypto.PubkeyToAddress(*publicKey)
	// 返回地址
	return address.Hex(), nil
}

// EthGenerateAddressFromMnemonic Eth 根据助记词生成地址和私钥
func EthGenerateAddressFromMnemonic(mnemonic string) (address string, privateKey string, err error) {
	// 生成私钥  这里可以选择传入指定密码或者空字符串，不同密码生成的助记词不同
	seed := bip39.NewSeed(mnemonic, "")
	// 使用种子生成 BIP32 主密钥
	wallet, err := EthNewFromSeed(seed)
	if err != nil {
		return "", "", err
	}
	path, err := MustParseDerivationPath(EthDerivationPath) //最后一位是同一个助记词的地址id，从0开始，相同助记词可以生产无限个地址
	if err != nil {
		return "", "", err
	}
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", err
	}
	address = account.Address.Hex()
	privateKey, err = wallet.PrivateKeyHex(account)
	if err != nil {
		return "", "", err
	}

	return address, privateKey, nil
}

// EthCreateWallet ETH   创建钱包 生成地址和私钥 助记词
func EthCreateWallet() (address string, privateKey string, mnemonic string, err error) {

	mnemonic, err = GetMnemonicBy256()
	if err != nil {
		return "", "", "", err
	}

	address, privateKey, err = EthGenerateAddressFromMnemonic(mnemonic)
	if err != nil {
		return "", "", "", err
	}

	return address, privateKey, mnemonic, nil
}

func IsValidEthAddress(address string) bool {
	// 定义 BSC 地址的正则表达式
	bscAddressPattern := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return bscAddressPattern.MatchString(address)
}
