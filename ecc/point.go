package ecc

import (
	"fmt"
	"math/big"
)

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
	a *FieldElement
	b *FieldElement

	// points on the curve
	x *FieldElement
	y *FieldElement
}

func performOp(x, y *FieldElement, scalar *big.Int, op Op) *FieldElement {
	switch op {
	case Add:
		return x.Add(y)
	case Sub:
		return x.Substract(y)
	case Mul:
		if y != nil {
			return x.Multiply(y)
		}

		if scalar != nil {
			return x.ScalarMul(scalar)
		}

		panic("required a *FieldElement or a scalar")
	case Div:
		return x.Divide(y)
	case Exp:
		if scalar == nil {
			panic("scalar should not be nil")
		}
		return x.Power(scalar)
	default:
		panic("unreacheable")
	}
}

func CheckIsOnCurve(x, y, a, b *FieldElement) bool {
	// check if the point is on the curve
	// y ^ 2 = x ^ 3 + a * x + b

	// y ^ 2
	lhs := performOp(y, nil, big.NewInt(2), Exp)

	// x ^ 3
	xPwr3 := performOp(x, nil, big.NewInt(3), Exp)

	// a * x
	aTimesX := performOp(a, x, nil, Mul)

	rhs := performOp(performOp(xPwr3, aTimesX, nil, Add), b, nil, Add)

	return lhs.EqualTo(rhs)
}

func NewIdentityPoint(a, b *FieldElement) *Point {
	return &Point{
		a: a,
		b: b,
		x: nil,
		y: nil,
	}
}

func NewPoint(x, y, a, b *FieldElement) *Point {
	if !CheckIsOnCurve(x, y, a, b) {
		panic("point is not in the curve")
	}

	return &Point{
		a: a,
		b: b,
		x: x,
		y: y,
	}
}

func (p *Point) EqualTo(other *Point) bool {
	return p.a.EqualTo(other.a) &&
		p.b.EqualTo(other.b) &&
		p.x.EqualTo(other.x) &&
		p.y.EqualTo(other.y)
}

func (p *Point) NotEqual(other *Point) bool {
	return !p.EqualTo(other)
}

func (p *Point) Add(other *Point) *Point {
	if !p.a.EqualTo(other.a) || !p.b.EqualTo(other.b) {
		panic("points are not on the same curve")
	}

	if p.x == nil {
		return other
	} else if other.x == nil {
		return p
	}

	// points are on the vertical A(x, y) B (x, -y)
	zero := NewFieldElement(p.x.order, big.NewInt(0))
	if p.x.EqualTo(other.x) && performOp(p.y, other.y, nil, Add).EqualTo(zero) {
		return NewIdentityPoint(p.a, p.b)
	}

	// find the slope of line AB
	var (
		x1, y1                 = p.x, p.y
		x2, y2                 = other.x, other.y
		numerator, denominator = zero, zero
	)

	if x1.EqualTo(x2) && y1.EqualTo(y2) {
		// slope = [3 * (x ^ 2) + a ] / 2y
		numerator = performOp(
			performOp(performOp(x1, nil, big.NewInt(2), Exp), nil, big.NewInt(3), Mul),
			p.a, nil, Add)
		denominator = performOp(y1, nil, big.NewInt(2), Mul)
	} else {
		numerator = performOp(y2, y1, nil, Sub)
		denominator = performOp(x2, x1, nil, Sub)
	}

	slope := performOp(numerator, denominator, nil, Div)
	slopePwr2 := performOp(slope, nil, big.NewInt(2), Exp)
	x3 := performOp(performOp(slopePwr2, x1, nil, Sub), x2, nil, Sub)
	x3Subx1 := performOp(x3, x1, nil, Sub)
	y3 := performOp(performOp(slope, x3Subx1, nil, Mul), y1, nil, Add)
	minusY3 := performOp(y3, nil, big.NewInt(-1), Mul)

	return &Point{
		x: x3,
		y: minusY3,
		a: p.a,
		b: p.b,
	}
}

// ScalarMul uses binary expansion to execute a optimized
// multiplication mainly with big scalar values
func (p *Point) ScalarMul(s *big.Int) *Point {
	if s == nil {
		panic("scalar cannot be nil")
	}

	curr := p
	result := NewIdentityPoint(p.a, p.b)

	numBytesRepr := s.Bits()
	for i := 0; i < s.BitLen(); i++ {
		byteOffset := i / 64
		pos := uint(1 << (i % 64))

		if (uint(numBytesRepr[byteOffset]) & pos) == pos {
			result = result.Add(curr)
		}

		curr = curr.Add(curr)
	}

	return result
}

func (p *Point) Verify(z *FieldElement, sig *Signature) bool {
	sInv := sig.s.Inverse()
	u := z.Multiply(sInv)
	v := sig.r.Multiply(sInv)
	bigR := (BitcoingGenPoint.ScalarMul(u.num)).Add(p.ScalarMul(v.num))

	return bigR.x.num.Cmp(sig.r.num) == 0
}

func (p *Point) String() string {
	var (
		xString, yString = "nil", "nil"
	)

	if p.x != nil {
		xString = p.x.String()
	}
	if p.y != nil {
		yString = p.y.String()
	}

	return fmt.Sprintf("Point(x: %s, y: %s, a: %s, b: %s)", xString, yString, p.a.String(), p.b.String())
}

func S256Point(x, y *big.Int) *Point {
	a := S256Field(big.NewInt(0))
	b := S256Field(big.NewInt(7))
	if x == nil && y == nil {
		return NewIdentityPoint(a, b)
	}

	return NewPoint(S256Field(x), S256Field(y), a, b)
}
