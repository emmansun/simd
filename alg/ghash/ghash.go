package ghash

import "encoding/binary"

// ref: GCM revised spec.
// Each element is a vector of 128 bits. The ith bit of an element X is denoted as Xi. The leftmost
// bit is X0, and the rightmost bit is X127.
// So, bit 0 of X is X[0] >> 7, bit 1 is X[0] >> 6, ..., bit 7 is X[0] & 1, bit 8 is X[1] >> 7, bit 127 is X[15] & 1.

// rawMethod represents a raw GHASH method, no optimzation.
type rawMethod struct {
	key [16]byte
}

func double(v *[16]byte) {
	var carryIn byte
	for j := range v {
		carryOut := (v[j] << 7) & 0x80
		v[j] = (v[j] >> 1) + carryIn
		carryIn = carryOut
	}
	if carryIn != 0 {
		v[0] ^= 0xE1 //  1<<7 | 1<<6 | 1<<5 | 1
	}
}

func double4(v *[16]byte) {
	for i := 0; i < 4; i++ {
		double(v)
	}
}

func double8(v *[16]byte) {
	for i := 0; i < 8; i++ {
		double(v)
	}
}

func xor(dst, z, y *[16]byte) {
	for i := range dst {
		dst[i] = z[i] ^ y[i]
	}
}

// Mul sets y to y * m.key.
func (m *rawMethod) Mul(y *[16]byte) {
	var z [16]byte
	for _, i := range m.key {
		for k := 0; k < 8; k++ {
			if (i>>(7-k))&1 == 1 {
				xor(&z, &z, y) // z ^= y
			}
			double(y)
		}
	}
	copy(y[:], z[:])
}

func clear(v *[16]byte) {
	for i := range v {
		v[i] = 0
	}
}

func reflect4Bits(i int) int {
	i = ((i << 2) & 0xc) | ((i >> 2) & 0x3)
	i = ((i << 1) & 0xa) | ((i >> 1) & 0x5)
	return i
}

// reflect8Bits returns the reflection of the 8-bit value i.
func reflect8Bits(i int) int {
	i = reflect4Bits(i&0xf)<<4 | reflect4Bits(i>>4)
	return i
}

type simpleMethod8Bits struct {
	key [16]byte
	m   [16][256][16]byte
}

// NewSimpleMethod8Bits returns a new GHASH method using the optimized method
// Simple, 8-bit tables.
func NewSimpleMethod8Bits(key []byte) *simpleMethod8Bits {
	m := new(simpleMethod8Bits)
	copy(m.key[:], key)
	var z, v [16]byte
	clear(&z)
	for i := 0; i < 16; i++ {
		copy(m.m[i][0][:], z[:]) // This is NOT needed.
		// reflect8Bits(1) == 128
		if i == 0 {
			copy(m.m[i][128][:], m.key[:])
		} else {
			copy(v[:], m.m[i-1][128][:])
			double8(&v)
			copy(m.m[i][128][:], v[:])
		}

		for j := 64; j > 0; j /= 2 {
			copy(v[:], m.m[i][2*j][:])
			double(&v)
			copy(m.m[i][j][:], v[:])
		}

		for j := 2; j < 256; j *= 2 {
			for k := 1; k < j; k++ {
				xor(&v, &m.m[i][j], &m.m[i][k])
				copy(m.m[i][j+k][:], v[:])
			}
		}

	}
	return m
}

func (m *simpleMethod8Bits) Mul(y *[16]byte) {
	var z [16]byte
	for i := 0; i < 16; i++ {
		xor(&z, &z, &m.m[i][y[i]])
	}
	copy(y[:], z[:])
}

type simpleMethod4Bits struct {
	key [16]byte
	m   [32][16][16]byte
}

// NewSimpleMethod4Bits returns a new GHASH method using the optimized method
// Simple, 4-bit tables.
func NewSimpleMethod4Bits(key []byte) *simpleMethod4Bits {
	m := new(simpleMethod4Bits)
	copy(m.key[:], key)
	var z, v [16]byte
	clear(&z)
	for i := 0; i < 32; i++ {
		copy(m.m[i][0][:], z[:]) // This is NOT needed.
		// reflect4Bits(1) == 8
		if i == 0 {
			copy(m.m[i][8][:], m.key[:])
		} else {
			copy(v[:], m.m[i-1][8][:])
			double4(&v)
			copy(m.m[i][8][:], v[:])
		}

		for j := 4; j > 0; j /= 2 {
			copy(v[:], m.m[i][2*j][:])
			double(&v)
			copy(m.m[i][j][:], v[:])
		}

		for j := 2; j < 16; j *= 2 {
			for k := 1; k < j; k++ {
				xor(&v, &m.m[i][j], &m.m[i][k])
				copy(m.m[i][j+k][:], v[:])
			}
		}
	}
	return m
}

func (m *simpleMethod4Bits) Mul(y *[16]byte) {
	var z [16]byte
	for i := 0; i < 32; i += 2 {
		xor(&z, &z, &m.m[i][y[i/2]>>4])
		xor(&z, &z, &m.m[i+1][y[i/2]&0xf])
	}
	copy(y[:], z[:])
}

type shoupMethod8Bits struct {
	key [16]byte
	m0  [256][16]byte // *H precomputed table
	r   [256][2]byte  // reduction table
}

func NewShoupMethod8Bits(key []byte) *shoupMethod8Bits {
	var v [16]byte
	m := new(shoupMethod8Bits)
	copy(m.key[:], key)
	clear(&m.m0[0]) // This is NOT needed.
	// reflect8Bits(1) == 128
	copy(m.m0[128][:], m.key[:])

	for j := 64; j > 0; j /= 2 {
		copy(v[:], m.m0[2*j][:])
		double(&v)
		copy(m.m0[j][:], v[:])
	}

	for j := 2; j < 256; j *= 2 {
		for k := 1; k < j; k++ {
			xor(&v, &m.m0[j], &m.m0[k])
			copy(m.m0[j+k][:], v[:])
		}
	}

	for i := 1; i < 256; i++ {
		clear(&v)
		v[15] = byte(i)
		double8(&v)
		copy(m.r[i][:], v[:2])
	}

	return m
}

func (m *shoupMethod8Bits) Mul(y *[16]byte) {
	var z [16]byte
	var a byte
	for i := 15; i > 0; i-- {
		// z = y * H
		xor(&z, &z, &m.m0[y[i]])

		// z = z * P^8
		a = z[15]
		for j := 15; j > 0; j-- {
			z[j] = z[j-1]
		}
		z[0] = m.r[a][0]
		z[1] ^= m.r[a][1]
	}
	xor(&z, &z, &m.m0[y[0]])
	copy(y[:], z[:])
}

func (m *shoupMethod8Bits) MulImpl2(y *[16]byte) {
	var z [16]byte
	var a byte
	for i := 15; i >= 0; i-- {
		// z = z * P^8
		a = z[15]
		for j := 15; j > 0; j-- {
			z[j] = z[j-1]
		}
		z[0] = m.r[a][0]
		z[1] ^= m.r[a][1]

		// z = y * H
		xor(&z, &z, &m.m0[y[i]])

	}
	copy(y[:], z[:])
}

type shoupMethod4Bits struct {
	key [16]byte
	m0  [16][16]byte // *H precomputed table
	r   [16][2]byte  // reduction table
}

func NewShoupMethod4Bits(key []byte) *shoupMethod4Bits {
	var v [16]byte
	m := new(shoupMethod4Bits)
	copy(m.key[:], key)
	clear(&m.m0[0]) // This is NOT needed.
	// reflect4Bits(1) == 8
	copy(m.m0[8][:], m.key[:])

	for j := 4; j > 0; j /= 2 {
		copy(v[:], m.m0[2*j][:])
		double(&v)
		copy(m.m0[j][:], v[:])
	}

	for j := 2; j < 16; j *= 2 {
		for k := 1; k < j; k++ {
			xor(&v, &m.m0[j], &m.m0[k])
			copy(m.m0[j+k][:], v[:])
		}
	}

	for i := 1; i < 16; i++ {
		clear(&v)
		v[15] = byte(i)
		double4(&v)
		copy(m.r[i][:], v[:2])
	}

	return m
}

func (m *shoupMethod4Bits) double4(v *[16]byte) {
	a := v[15] & 0xf
	for j := 15; j > 0; j-- {
		v[j] = (v[j] >> 4) | (v[j-1] << 4)
	}
	v[0] = v[0] >> 4
	v[0] ^= m.r[a][0]
	v[1] ^= m.r[a][1]
}

// We can call double4 twice to implement double8.
// But the following implementation is more efficient.
func (m *shoupMethod4Bits) double8(v *[16]byte) {
	a := v[15] & 0xf
	b := v[15] >> 4
	for j := 15; j > 0; j-- {
		v[j] = v[j-1]
	}
	v[0] = (m.r[a][0] >> 4) ^ m.r[b][0]
	v[1] ^= (m.r[a][1]>>4 | m.r[a][0]<<4) ^ m.r[b][1]
	//m.r[a][1] << 4 == 0
	//v[2] ^= (m.r[a][1] << 4)
}

func (m *shoupMethod4Bits) mulH(yi byte, z *[16]byte) {
	var zl [16]byte

	yh := yi >> 4
	yl := yi & 0xf

	xor(z, z, &m.m0[yh])
	copy(zl[:], m.m0[yl][:])
	m.double4(&zl)
	xor(z, z, &zl)
}

func (m *shoupMethod4Bits) MulImpl3(y *[16]byte) {
	var z [16]byte
	for i := 15; i > 0; i-- {
		m.mulH(y[i], &z)
		m.double8(&z)
	}
	m.mulH(y[0], &z)
	copy(y[:], z[:])
}

func (m *shoupMethod4Bits) Mul(y *[16]byte) {
	var z [16]byte
	for i := 15; i >= 1; i-- {
		w := y[i]
		for j := 0; j < 2; j++ {
			xor(&z, &z, &m.m0[w&0xf])
			m.double4(&z)
			w >>= 4
		}
	}
	xor(&z, &z, &m.m0[y[0]&0xf])
	double4(&z)
	xor(&z, &z, &m.m0[y[0]>>4])

	copy(y[:], z[:])
}

func (m *shoupMethod4Bits) MulImpl2(y *[16]byte) {
	var z [16]byte
	for i := 15; i >= 0; i-- {
		w := y[i]
		for j := 0; j < 2; j++ {
			m.double4(&z)
			xor(&z, &z, &m.m0[w&0xf])
			w >>= 4
		}
	}
	copy(y[:], z[:])
}

// gcmFieldElement represents a value in GF(2¹²⁸). In order to reflect the GCM
// standard and make binary.BigEndian suitable for marshaling these values, the
// bits are stored in big endian order. For example:
//
//	the coefficient of x⁰ can be obtained by v.low >> 63.
//	the coefficient of x⁶³ can be obtained by v.low & 1.
//	the coefficient of x⁶⁴ can be obtained by v.high >> 63.
//	the coefficient of x¹²⁷ can be obtained by v.high & 1.
type gcmFieldElement struct {
	low, high uint64
}

// It's similar to the Shoup method.
type gcmMethod struct {
	productTable [16]gcmFieldElement
}

func NewGCMMethod(key []byte) *gcmMethod {
	// We precompute 16 multiples of |key|. However, when we do lookups
	// into this table we'll be using bits from a field element and
	// therefore the bits will be in the reverse order. So normally one
	// would expect, say, 4*key to be in index 4 of the table but due to
	// this bit ordering it will actually be in index 0010 (base 2) = 2.
	x := gcmFieldElement{
		binary.BigEndian.Uint64(key[:8]),
		binary.BigEndian.Uint64(key[8:]),
	}
	m := new(gcmMethod)
	m.productTable[8] = x

	for j := 4; j > 0; j /= 2 {
		m.productTable[j] = gcmDouble(&m.productTable[j*2])
	}

	for j := 2; j < 16; j *= 2 {
		for k := 1; k < j; k++ {
			m.productTable[j+k] = gcmAdd(&m.productTable[j], &m.productTable[k])
		}
	}
	return m
}

// gcmAdd adds two elements of GF(2¹²⁸) and returns the sum.
func gcmAdd(x, y *gcmFieldElement) gcmFieldElement {
	// Addition in a characteristic 2 field is just XOR.
	return gcmFieldElement{x.low ^ y.low, x.high ^ y.high}
}

// gcmDouble returns the result of doubling an element of GF(2¹²⁸).
func gcmDouble(x *gcmFieldElement) (double gcmFieldElement) {
	msbSet := x.high&1 == 1

	// Because of the bit-ordering, doubling is actually a right shift.
	double.high = x.high >> 1
	double.high |= x.low << 63
	double.low = x.low >> 1

	// If the most-significant bit was set before shifting then it,
	// conceptually, becomes a term of x^128. This is greater than the
	// irreducible polynomial so the result has to be reduced. The
	// irreducible polynomial is 1+x+x^2+x^7+x^128. We can subtract that to
	// eliminate the term at x^128 which also means subtracting the other
	// four terms. In characteristic 2 fields, subtraction == addition ==
	// XOR.
	if msbSet {
		double.low ^= 0xe100000000000000
	}

	return
}

func (m *gcmMethod) Mul(y *[16]byte) {
	yField := gcmFieldElement{
		binary.BigEndian.Uint64(y[:8]),
		binary.BigEndian.Uint64(y[8:]),
	}
	var z gcmFieldElement

	for i := 0; i < 2; i++ {
		word := yField.high
		if i == 1 {
			word = yField.low
		}

		// Multiplication works by multiplying z by 16 and adding in
		// one of the precomputed multiples of H.
		for j := 0; j < 64; j += 4 {
			msw := z.high & 0xf
			z.high >>= 4
			z.high |= z.low << 60
			z.low >>= 4
			z.low ^= uint64(gcmReductionTable[msw]) << 48

			// the values in |table| are ordered for
			// little-endian bit positions. See the comment
			// in NewGCMWithNonceSize.
			t := &m.productTable[word&0xf]

			z.low ^= t.low
			z.high ^= t.high
			word >>= 4
		}
	}

	binary.BigEndian.PutUint64(y[:8], z.low)
	binary.BigEndian.PutUint64(y[8:], z.high)
}

// gcmReductionTable is stored irreducible polynomial's double & add precomputed results.
// 0000 - 0
// 0001 - irreducible polynomial >> 3
// 0010 - irreducible polynomial >> 2
// 0011 - (irreducible polynomial >> 3 xor irreducible polynomial >> 2)
// ...
// 1000 - just the irreducible polynomial
var gcmReductionTable = []uint16{
	0x0000, 0x1c20, 0x3840, 0x2460, 0x7080, 0x6ca0, 0x48c0, 0x54e0,
	0xe100, 0xfd20, 0xd940, 0xc560, 0x9180, 0x8da0, 0xa9c0, 0xb5e0,
}
