package btc

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
	"log"
)

// GenerateNestedSegWitAddressFromMnemonic Nested SegWit // P2SH-P2WPKH
func GenerateNestedSegWitAddressFromMnemonic(mnemonic string, params *chaincfg.Params) (string, string, error) {

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 创建主私钥
	masterKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		log.Fatal("Error generating master key:", err)
	}

	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 49) // BIP 49 (P2WPKH-P2SH)
	if err != nil {
		return "", "", err
	}

	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 0) // Bitcoin mainnet
	if err != nil {
		return "", "", err
	}

	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0) // Account 0
	if err != nil {
		return "", "", err
	}

	change, _ := account.Derive(0) // External change
	index := uint32(0)             // Address index

	privateKey, err := change.Derive(index)
	if err != nil {
		return "", "", err
	}

	publicKey, err := privateKey.ECPubKey()
	if err != nil {
		return "", "", err
	}

	p2shAddress, err := publicKeyToP2SHAddress(publicKey, params)
	if err != nil {
		return "", "", err
	}
	bb, _ := privateKey.ECPrivKey()
	wif, err := btcutil.NewWIF(bb, params, true)
	if err != nil {
		log.Fatal("Error generating WIF:", err)
		return "", "", err
	}
	//fmt.Println(wif.String())
	return p2shAddress.EncodeAddress(), wif.String(), nil
}

// GenerateNativeSegWitAddressFromMnemonic Native SegWit // P2WPKH
func GenerateNativeSegWitAddressFromMnemonic(mnemonic string, params *chaincfg.Params) (string, string, error) {

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 创建主私钥
	masterKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		log.Fatal("Error generating master key:", err)
	}

	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 84) // BIP 49 (P2WPKH-P2SH)
	if err != nil {
		return "", "", err
	}

	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 0) // Bitcoin mainnet
	if err != nil {
		return "", "", err
	}

	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0) // Account 0
	if err != nil {
		return "", "", err
	}

	change, _ := account.Derive(0) // External change
	index := uint32(0)             // Address index

	privateKey, err := change.Derive(index)
	if err != nil {
		return "", "", err
	}

	publicKey, err := privateKey.ECPubKey()
	if err != nil {
		return "", "", err
	}

	p2shAddress, err := publicKeyToP2WPKHAddress(publicKey, params)
	if err != nil {
		return "", "", err
	}
	bb, _ := privateKey.ECPrivKey()
	wif, err := btcutil.NewWIF(bb, params, true)
	if err != nil {
		log.Fatal("Error generating WIF:", err)
		return "", "", err
	}
	//fmt.Println(wif.String())
	return p2shAddress.EncodeAddress(), wif.String(), nil
}

// GenerateLegacyAddressFromMnemonic Legacy // Legacy
func GenerateLegacyAddressFromMnemonic(mnemonic string, params *chaincfg.Params) (string, string, error) {

	// 从助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 创建主私钥
	masterKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		log.Fatal("Error generating master key:", err)
	}

	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44) // BIP 49 (P2WPKH-P2SH)
	if err != nil {
		return "", "", err
	}

	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 0) // Bitcoin mainnet
	if err != nil {
		return "", "", err
	}

	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0) // Account 0
	if err != nil {
		return "", "", err
	}

	change, _ := account.Derive(0) // External change
	index := uint32(0)             // Address index

	privateKey, err := change.Derive(index)
	if err != nil {
		return "", "", err
	}

	publicKey, err := privateKey.ECPubKey()
	if err != nil {
		return "", "", err
	}

	p2shAddress, err := publicKeyToLegacyAddress(publicKey, params)
	if err != nil {
		return "", "", err
	}
	bb, _ := privateKey.ECPrivKey()
	wif, err := btcutil.NewWIF(bb, params, true)
	if err != nil {
		log.Fatal("Error generating WIF:", err)
		return "", "", err
	}
	return p2shAddress.EncodeAddress(), wif.String(), nil
}

// GenerateNestedSegWitAddressFromPrivateKey 根据 私钥生成 Nested SegWit // P2SH-P2WPKH
func GenerateNestedSegWitAddressFromPrivateKey(privateKeyString string, params *chaincfg.Params) (string, error) {
	// 将私钥的 Base58 编码字符串解码为字节切片
	privateKeyBytes, err := btcutil.DecodeWIF(privateKeyString)
	if err != nil {
		return "", err
	}

	// 从字节切片中获取私钥对象
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes.PrivKey.Serialize())

	// 从私钥生成公钥
	publicKey := privateKey.PubKey()

	// 生成 P2SH-P2WPKH 地址
	p2shAddress, err := publicKeyToP2SHAddress(publicKey, params)
	if err != nil {
		return "", err
	}

	return p2shAddress.EncodeAddress(), nil
}

// GenerateNativeSegWitAddressFromPrivateKey 根据 私钥生成 Native SegWit // P2WPKH
func GenerateNativeSegWitAddressFromPrivateKey(privateKeyString string, params *chaincfg.Params) (string, error) {
	// 将私钥的 Base58 编码字符串解码为字节切片
	privateKeyBytes, err := btcutil.DecodeWIF(privateKeyString)
	if err != nil {
		return "", err
	}

	// 从字节切片中获取私钥对象
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes.PrivKey.Serialize())
	// 从私钥生成公钥
	publicKey := privateKey.PubKey()
	// 生成 P2SH-P2WPKH 地址
	p2wpkhAddress, err := publicKeyToP2WPKHAddress(publicKey, params)
	if err != nil {
		return "", err
	}

	return p2wpkhAddress.EncodeAddress(), nil
}

// GenerateLegacyAddressFromPrivateKey 根据 私钥生成 Legacy // Legacy
func GenerateLegacyAddressFromPrivateKey(privateKeyString string, params *chaincfg.Params) (string, error) {
	// 将私钥的 Base58 编码字符串解码为字节切片
	privateKeyBytes, err := btcutil.DecodeWIF(privateKeyString)
	if err != nil {
		return "", err
	}

	// 从字节切片中获取私钥对象
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes.PrivKey.Serialize())

	// 从私钥生成公钥
	publicKey := privateKey.PubKey()

	// 生成 Legacy 地址
	legacyAddress, err := publicKeyToLegacyAddress(publicKey, params)
	if err != nil {
		return "", err
	}

	return legacyAddress.EncodeAddress(), nil
}

// 生成 P2SH-P2WPKH 地址
func publicKeyToP2SHAddress(publicKey *btcec.PublicKey, Params *chaincfg.Params) (btcutil.Address, error) {
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	script, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
	if err != nil {
		return nil, err
	}
	return btcutil.NewAddressScriptHash(script, Params)
}

// 生成 P2WPKH 地址
func publicKeyToP2WPKHAddress(publicKey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	return btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, params)
}

// 生成 Legacy 地址
func publicKeyToLegacyAddress(publicKey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	return btcutil.NewAddressPubKeyHash(pubKeyHash, params)
}

// 生成 Taproot 地址
func publicKeyToTaprootAddress(publicKey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	return btcutil.NewAddressWitnessScriptHash(pubKeyHash, params)
}
