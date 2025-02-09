package ecc_test

import (
	"crypto/sha256"
	"ecc"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckPointIsOnCurve(t *testing.T) {
	order := big.NewInt(233)

	var a, b = ecc.NewFieldElement(order, big.NewInt(5)), ecc.NewFieldElement(order, big.NewInt(7))
	points := [...]struct {
		x        *ecc.FieldElement
		y        *ecc.FieldElement
		expected bool
	}{
		{ecc.NewFieldElement(order, big.NewInt(-1)), ecc.NewFieldElement(order, big.NewInt(-1)), true},
		{ecc.NewFieldElement(order, big.NewInt(-1)), ecc.NewFieldElement(order, big.NewInt(-2)), false},
		{ecc.NewFieldElement(order, big.NewInt(2)), ecc.NewFieldElement(order, big.NewInt(4)), false},
		{ecc.NewFieldElement(order, big.NewInt(18)), ecc.NewFieldElement(order, big.NewInt(77)), true},
		{ecc.NewFieldElement(order, big.NewInt(5)), ecc.NewFieldElement(order, big.NewInt(7)), false},
	}

	for _, p := range points {
		t.Run(fmt.Sprintf("testing_points_%s_%s", p.x.String(), p.y.String()), func(t *testing.T) {
			require.Equal(t, p.expected, ecc.CheckIsOnCurve(p.x, p.y, a, b))
		})
	}
}

func TestPointAddIdentity(t *testing.T) {
	order := big.NewInt(233)
	var a, b = ecc.NewFieldElement(order, big.NewInt(5)), ecc.NewFieldElement(order, big.NewInt(7))

	p := ecc.NewPoint(
		ecc.NewFieldElement(order, big.NewInt(-1)), ecc.NewFieldElement(order, big.NewInt(-1)), a, b)
	id := ecc.NewIdentityPoint(a, b)

	fmt.Println(p)

	require.Equal(t, p, p.Add(id))
}

func TestPointAddition(t *testing.T) {
	order := big.NewInt(223)
	var a, b = ecc.NewFieldElement(order, big.NewInt(0)), ecc.NewFieldElement(order, big.NewInt(7))

	// (192, 105)
	x1 := ecc.NewFieldElement(order, big.NewInt(192))
	y1 := ecc.NewFieldElement(order, big.NewInt(105))

	p := ecc.NewPoint(x1, y1, a, b)
	fmt.Println(p)

	yNeg := y1.Negate()
	p2 := ecc.NewPoint(x1, yNeg, a, b)
	res := p.Add(p2)

	// should result in a identity point
	require.Equal(t, ecc.NewIdentityPoint(a, b), res)
	fmt.Println(res)

	x2 := ecc.NewFieldElement(order, big.NewInt(17))
	y2 := ecc.NewFieldElement(order, big.NewInt(56))

	p3 := ecc.NewPoint(x2, y2, a, b)
	fmt.Println(p3)

	res = p.Add(p3)
	fmt.Println(res)
	require.Equal(
		t,
		ecc.NewPoint(
			ecc.NewFieldElement(order, big.NewInt(170)),
			ecc.NewFieldElement(order, big.NewInt(142)),
			a,
			b),
		res,
	)

	res = p.Add(p)
	fmt.Println(res)
	require.Equal(
		t,
		ecc.NewPoint(
			ecc.NewFieldElement(order, big.NewInt(49)),
			ecc.NewFieldElement(order, big.NewInt(71)),
			a,
			b),
		res,
	)
}

func TestPointScalarMul(t *testing.T) {
	order := big.NewInt(223)
	var a, b = ecc.NewFieldElement(order, big.NewInt(0)), ecc.NewFieldElement(order, big.NewInt(7))

	// 2 * (192, 105)
	x1 := ecc.NewFieldElement(order, big.NewInt(192))
	y1 := ecc.NewFieldElement(order, big.NewInt(105))

	p := ecc.NewPoint(x1, y1, a, b)
	res := p.ScalarMul(big.NewInt(2))
	fmt.Println(res)

	// 2 * (143, 98)
	x2 := ecc.NewFieldElement(order, big.NewInt(143))
	y2 := ecc.NewFieldElement(order, big.NewInt(98))

	p2 := ecc.NewPoint(x2, y2, a, b)
	res = p2.ScalarMul(big.NewInt(2))
	fmt.Println(res)
}

func TestPointWithBTCSetting(t *testing.T) {
	gx := big.NewInt(0)
	gx.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)

	gy := new(big.Int)
	gy.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	fst := ecc.S256Point(gx, gy)

	fmt.Println(fst)

	n := new(big.Int)
	n.SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	fmt.Println(fst.ScalarMul(n))
}

func TestSigVerification(t *testing.T) {
	x1 := big.NewInt(0)
	x1.SetString("4519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)

	y1 := big.NewInt(0)
	y1.SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)

	p1 := ecc.S256Point(x1, y1)

	fmt.Println(p1)

	z := big.NewInt(0)
	z.SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
	zField := ecc.NewFieldElement(ecc.BitcoinN, z)

	r := big.NewInt(0)
	r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	rField := ecc.NewFieldElement(ecc.BitcoinN, r)

	s := big.NewInt(0)
	s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sField := ecc.NewFieldElement(ecc.BitcoinN, s)

	sig := ecc.NewSignature(rField, sField)

	require.True(t, p1.Verify(zField, sig))
}

func TestPrivateKeySig(t *testing.T) {
	dummyPayload := sha256.Sum256([]byte("dummy_payload"))
	dummyPayload = sha256.Sum256(dummyPayload[:])

	z := big.NewInt(0)
	z.SetBytes(dummyPayload[:])
	zField := ecc.NewFieldElement(ecc.BitcoinN, z)

	privateK := ecc.NewPrivateKey(big.NewInt(12345))
	sig := privateK.Sign(z)

	fmt.Println(sig)

	pubK := privateK.PublicKey()

	require.True(t, pubK.Verify(zField, sig))

	// if using another public key
	anotherPrivate := ecc.NewPrivateKey(big.NewInt(4444))
	anotherPublic := anotherPrivate.PublicKey()
	require.False(t, anotherPublic.Verify(zField, sig))
}
