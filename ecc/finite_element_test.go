package ecc_test

import (
	"ecc"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFiniteElement(t *testing.T) {
	order := big.NewInt(57)

	f44 := ecc.NewFieldElement(order, big.NewInt(44))
	f33 := ecc.NewFieldElement(order, big.NewInt(33))

	require.True(t, f44.Add(f33).EqualTo(ecc.NewFieldElement(order, big.NewInt(20))))

	require.True(t, f44.Add(f33).Negate().EqualTo(ecc.NewFieldElement(order, big.NewInt(37))))

	require.True(t, f44.Substract(f33).EqualTo(ecc.NewFieldElement(order, big.NewInt(11))))

	require.True(t, f33.Substract(f44).EqualTo(ecc.NewFieldElement(order, big.NewInt(46))))

	scalar := big.NewInt(46)
	require.Equal(t,
		ecc.NewFieldElement(order, scalar).ScalarMul(scalar),
		ecc.NewFieldElement(order, scalar).Power(big.NewInt(2)),
	)

	f46 := ecc.NewFieldElement(order, big.NewInt(46))
	require.True(t,
		f46.Multiply(f46).EqualTo(
			f46.Power(big.NewInt(58)),
		),
	)

	order19 := big.NewInt(19)

	div := ecc.NewFieldElement(order19, big.NewInt(2)).Divide(ecc.NewFieldElement(order19, big.NewInt(3)))
	mul := ecc.NewFieldElement(order19, big.NewInt(3)).Multiply(ecc.NewFieldElement(order19, big.NewInt(7)))

	require.True(t, div.EqualTo(ecc.NewFieldElement(order19, big.NewInt(7))))
	require.True(t, mul.EqualTo(ecc.NewFieldElement(order19, big.NewInt(2))))

}

func TestFinitSetOf19thOrder(t *testing.T) {
	var elements []*ecc.FieldElement
	for num := 0; num < 19; num++ {
		out := ecc.NewFieldElement(big.NewInt(19), big.NewInt(int64(num))).ScalarMul(big.NewInt(2))
		elements = append(elements, out)
	}

	for idx, e := range elements {
		fmt.Println("#", idx, e.String())
	}
}

func TestFinitSetOf19thOrderToThePwrOfOrderMinusOne(t *testing.T) {
	var elements []*ecc.FieldElement
	for num := 1; num < 19; num++ {
		out := ecc.NewFieldElement(big.NewInt(19), big.NewInt(int64(num))).Power(big.NewInt(18))
		elements = append(elements, out)

		bN := big.NewInt(int64(num))
		r := big.NewInt(0).Mod(big.NewInt(0).Mul(bN, big.NewInt(0).Exp(bN, big.NewInt(17), nil)), big.NewInt(19))

		fmt.Println("#", num, r.String())
	}

	for idx, e := range elements {
		fmt.Println("#", idx, e.String())
	}
}

func TestSqrt(t *testing.T) {
	x := new(big.Int)
	x.SetString("49fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a", 16)
	//y^2 = x^3 + 7
	y2 := (ecc.S256Field(x).Power(big.NewInt(int64(3)))).Add(ecc.S256Field(big.NewInt(int64(7))))
	y := y2.Sqrt()

	fmt.Printf("y value for gien x is %s\n", y)

	//check x,y on the curve
	p := ecc.NewPoint(ecc.S256Field(x), y, ecc.S256Field(big.NewInt(int64(0))),
		ecc.S256Field(big.NewInt(int64(7))))

	require.NotNil(t, p)

	fmt.Println("point from the SEC compressed format is on the curve")
}
