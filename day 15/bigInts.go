package main

import "math/big"

type bigInt = *big.Int

func leBig(a, b bigInt) bool {
	return a.Cmp(b) == -1
}

func makeBigInt(val int) bigInt {
	return big.NewInt(int64(val))
}

func copyBigInt(val bigInt) bigInt {
	return new(big.Int).Set(val)
}

func addBig(a, b bigInt) bigInt {
	return makeBigInt(0).Add(a, b)
}

func equalsBig(a,b bigInt) bool {
	return a.Cmp(b) == 0
}
