package ecc

import "math/big"

var (
	Zero = big.NewInt(0)

	// 2 ^ 256 - 2 ^ 32 - 977
	BitcoinOrder = big.NewInt(0).
			Sub(
			big.NewInt(0).
				Sub(
					big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
					big.NewInt(0).Exp(big.NewInt(2), big.NewInt(32), nil),
				),
			big.NewInt(977),
		)

	BitcoinN = func() *big.Int {
		n := big.NewInt(0)
		n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
		return n
	}()

	BitcoinGenX = func() *big.Int {
		gx := big.NewInt(0)
		gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
		return gx
	}()

	BitcoinGenY = func() *big.Int {
		gy := big.NewInt(0)
		gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
		return gy
	}()

	BitcoingGenPoint = S256Point(BitcoinGenX, BitcoinGenY)
)
