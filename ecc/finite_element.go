package ecc

import (
	"fmt"
	"math/big"
)

type FieldElement struct {
	order *big.Int
	num   *big.Int
}

func NewFieldElement(order, num *big.Int) *FieldElement {
	if num.Cmp(order) >= 0 {
		panic(fmt.Sprintf("num cannot be greater than %s", order.String()))
	}

	return &FieldElement{
		order: order,
		num:   num,
	}
}

func (f *FieldElement) String() string {
	return fmt.Sprintf("FieldElement{order: %s, num: %s}", f.order.String(), f.num.String())
}

func (f *FieldElement) EqualTo(other *FieldElement) bool {
	return f.order.Cmp(other.order) == 0 && f.num.Cmp(other.num) == 0
}

func (f *FieldElement) checkOrder(other *FieldElement) {
	if f.order.Cmp(other.order) != 0 {
		panic("field elements does not have the same order")
	}
}

func (f *FieldElement) Add(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	return NewFieldElement(f.order,
		big.NewInt(0).Mod(big.NewInt(0).Add(f.num, other.num), f.order))
}

func (f *FieldElement) Negate() *FieldElement {
	return NewFieldElement(f.order, big.NewInt(0).Sub(f.order, f.num))
}

func (f *FieldElement) Substract(other *FieldElement) *FieldElement {
	f.checkOrder(other)
	return f.Add(other.Negate())
}

func (f *FieldElement) Multiply(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	return NewFieldElement(f.order,
		big.NewInt(0).Mod(big.NewInt(0).Mul(f.num, other.num), f.order))
}

func (f *FieldElement) Power(pwr *big.Int) *FieldElement {
	t := big.NewInt(0).Mod(pwr, big.NewInt(0).Sub(f.order, big.NewInt(1)))
	res := big.NewInt(0).Exp(f.num, t, f.order)
	modRes := big.NewInt(0).Mod(res, f.order)
	return NewFieldElement(f.order, modRes)
}

func (f *FieldElement) ScalarMul(v *big.Int) *FieldElement {
	return NewFieldElement(f.order, big.NewInt(0).Mod(big.NewInt(0).Mul(f.num, v), f.order))
}

func (f *FieldElement) Divide(other *FieldElement) *FieldElement {
	f.checkOrder(other)

	// c * b = a
	// a / b = c

	// a / b == a * b ^ (p - 2)
	reverseOfOther := other.Power(big.NewInt(0).Sub(f.order, big.NewInt(2)))
	return f.Multiply(reverseOfOther)
}

func (f *FieldElement) Inverse() *FieldElement {
	return f.Power(big.NewInt(0).Sub(f.order, big.NewInt(2)))
}

func S256Field(num *big.Int) *FieldElement {
	return NewFieldElement(BitcoinOrder, num)
}
