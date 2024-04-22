package test

import (
	"fmt"
	"github.com/EntySquare/chain-util/pkg/btc"
	"github.com/btcsuite/btcd/chaincfg"
	"log"
	"testing"
)

func TestGetBtc(t *testing.T) {
	const mnemonic = "orient equal because fox twin pizza exit vote glue cheap car ancient"

	address, privateKeyStr, err := btc.GenerateNestedSegWitAddressFromMnemonic(mnemonic, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("NestedSegWit ============")
	fmt.Println(address, privateKeyStr)

	// 调用生成地址的函数
	address, err = btc.GenerateNestedSegWitAddressFromPrivateKey(privateKeyStr, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal("Error generating address:", err)
	}
	// 打印生成的地址和 WIF 格式的私钥
	fmt.Println("NestedSegWit 私钥 Address:", address)

	addressNative, privateKeyNativeStr, err := btc.GenerateNativeSegWitAddressFromMnemonic(mnemonic, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("NativeSegWit ============")
	fmt.Println(addressNative, privateKeyNativeStr)

	// 调用生成地址的函数
	addressNative, err = btc.GenerateNativeSegWitAddressFromPrivateKey(privateKeyNativeStr, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal("Error generating address:", err)
	}
	// 打印生成的地址和 WIF 格式的私钥
	fmt.Println("NativeSegWit 私钥 Address:", addressNative)

	addressLegacy, privateKeyLegacyStr, err := btc.GenerateLegacyAddressFromMnemonic(mnemonic, &chaincfg.MainNetParams)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Legacy ============")
	fmt.Println(addressLegacy, privateKeyLegacyStr)

	// 调用生成地址的函数
	addressLegacy, err = btc.GenerateLegacyAddressFromPrivateKey(privateKeyLegacyStr, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal("Error generating address:", err)
	}
	// 打印生成的地址和 WIF 格式的私钥
	fmt.Println("Legacy 私钥 Address:", addressLegacy)

}
