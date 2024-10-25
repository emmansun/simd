package ghash

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
		copy(v[:], m.key[:])
		// H * P^8
		for k := 0; k < i; k++ {
			double8(&v)
		}
		copy(m.m[i][0][:], z[:])
		copy(m.m[i][reflect8Bits(1)][:], v[:])

		for j := 2; j < 256; j += 2 {
			copy(v[:], m.m[i][reflect8Bits(j/2)][:])
			double(&v)
			copy(m.m[i][reflect8Bits(j)][:], v[:])
			xor(&v, &m.m[i][reflect8Bits(1)], &v)
			copy(m.m[i][reflect8Bits(j+1)][:], v[:])
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
	for i := 0; i < 32; i ++ {
		copy(v[:], m.key[:])
		// H * P^4
		for k := 0; k < i; k++ {
			double4(&v)
		}
		copy(m.m[i][0][:], z[:])
		copy(m.m[i][reflect4Bits(1)][:], v[:])

		for j := 2; j < 16; j += 2 {
			copy(v[:], m.m[i][reflect4Bits(j/2)][:])
			double(&v)
			copy(m.m[i][reflect4Bits(j)][:], v[:])
			xor(&v, &m.m[i][reflect4Bits(1)], &v)
			copy(m.m[i][reflect4Bits(j+1)][:], v[:])
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
