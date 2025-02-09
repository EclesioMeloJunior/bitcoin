package ecc

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	secret *big.Int
	pubKey *Point
}

func NewPrivateKey(secret *big.Int) *PrivateKey {
	return &PrivateKey{
		secret: secret,
		pubKey: BitcoingGenPoint.ScalarMul(secret),
	}
}

func (p *PrivateKey) String() string {
	return fmt.Sprintf("private key: {%s}", p.secret)
}

func (p *PrivateKey) PublicKey() *Point {
	return p.pubKey
}

func (p *PrivateKey) Sign(z *big.Int) *Signature {
	k, err := rand.Int(rand.Reader, BitcoinN)
	if err != nil {
		panic("failed to generate k")
	}

	r := BitcoingGenPoint.ScalarMul(k).x.num
	rField := NewFieldElement(BitcoinN, r)

	kField := NewFieldElement(BitcoinN, k)
	eField := NewFieldElement(BitcoinN, p.secret)
	zField := NewFieldElement(BitcoinN, z)

	// (z + r * e) / k
	sField := rField.Multiply(eField).Add(zField).Divide(kField)

	// s > n / 2 => s = n - s
	if sField.num.Cmp(big.NewInt(0).Div(BitcoinN, big.NewInt(2))) == 1 {
		sField = NewFieldElement(BitcoinN, big.NewInt(0).Sub(BitcoinN, sField.num))
	}

	return NewSignature(rField, sField)
}
