package ecc_test

import (
	"ecc"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDerEncoding(t *testing.T) {
	r := new(big.Int)
	r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	rField := ecc.S256Field(r)
	s := new(big.Int)
	s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sField := ecc.S256Field(s)
	sig := ecc.NewSignature(rField, sField)

	der := sig.Der()

	expected := big.NewInt(0)
	expected.SetString("3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)

	require.Equal(t, expected.Bytes(), der)
}
