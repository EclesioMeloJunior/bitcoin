package ecc

import "math/big"

type Op int

const (
	Add Op = iota
	Sub
	Mul
	Div
	Exp
)

type Point struct {
	// coefficients of the curve
	a *big.Int
	b *big.Int

	// points on the curve
	x *big.Int
	y *big.Int
}

func performOp(x, y *big.Int, op Op) *big.Int {
	var op *big.Int

	switch op {
	case Add:
		return op.Add(x, y)
	case Sub:
		return op.Sub(x, y)
	case Mul:
		return op.Mul(x, y)
	case Div:
		return op.Div(x, y)
	case Exp:
		return op.Exp(x, y, nil)
	}
}
