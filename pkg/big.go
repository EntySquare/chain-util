package pkg

import (
	"fmt"
	"math/big"
	"strconv"
)

func ConvertStringToFloat(amountStr string) (float64, error) {
	return strconv.ParseFloat(amountStr, 64)
}

// ConvertToFloatWithPrecision 根据币种精度将字符串金额 转换为浮点数 比如6位精度的币种 1000000 = 1 、 8位精度的币种 100000000 = 1
func ConvertToFloatWithPrecision(amountStr string, precision int64) (float64, error) {
	amountDec, success := new(big.Float).SetString(amountStr)
	if !success {
		return 0, fmt.Errorf("invalid amount format: %s", amountStr)
	}

	scale := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil))
	amountDec.Quo(amountDec, scale)

	amountFloat, _ := amountDec.Float64()
	return amountFloat, nil
}

// ConvertFloatToAmountString 根据币种精度将浮点数金额 转换为字符串 比如6位精度的币种 1 = 1000000 、 8位精度的币种 1 = 100000000
func ConvertFloatToAmountString(amountFloat float64, precision int64) string {
	amountBig := big.NewFloat(amountFloat)

	// 根据精度进行倍数调整
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil)
	amountBig.Mul(amountBig, new(big.Float).SetInt(multiplier))

	amountStr := amountBig.Text('f', 0)
	return amountStr
}

func BigIntAdd(num1, num2 string) *big.Int {
	bigNum1, _ := new(big.Int).SetString(num1, 10)
	if bigNum1 == nil {
		return nil
	}
	bigNum2, _ := new(big.Int).SetString(num2, 10)
	if bigNum2 == nil {
		return nil
	}
	sum := new(big.Int)
	sum.Add(bigNum1, bigNum2)
	return sum
}

// BigIntCmp
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func BigIntCmp(x, y string) int {
	intX := new(big.Int)
	intY := new(big.Int)
	intX.SetString(x, 10) // 10 表示十进制
	intY.SetString(y, 10)
	return intX.Cmp(intY)
}

func BigIntSub(num1, num2 string) *big.Int {
	bigNum1, _ := new(big.Int).SetString(num1, 10)
	if bigNum1 == nil {
		return nil
	}
	bigNum2, _ := new(big.Int).SetString(num2, 10)
	if bigNum2 == nil {
		return nil
	}
	sum := new(big.Int)
	sum.Sub(bigNum1, bigNum2)
	return sum
}

func Uint64ToBigInt(num uint64) *big.Int {
	result := new(big.Int)
	result.SetUint64(num)
	return result
}
