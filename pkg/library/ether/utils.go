package ether

import "math/big"

func hex2BigInt(hex string) *big.Int {
	num, _ := new(big.Int).SetString(hex[2:], 16)
	return num
}
