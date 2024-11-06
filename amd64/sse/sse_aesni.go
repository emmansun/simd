package sse

import "github.com/emmansun/simd/alg/aes"

var shift_row = Set64(0x0B06010C07020D08, 0x030E09040F0A0500)
var shift_row_inv = Set64(0x0306090C0F020508, 0x0B0E0104070A0D00)
var const_0f = Set64(0x0F0F0F0F0F0F0F0F, 0x0F0F0F0F0F0F0F0F)

func mm_aesenclast_si128(state, rk *XMM) {
	// ShiftRows
	mm_shuffle_epi8(state, &shift_row)
	// SubBytes
	tmp := XMM{}
	for i := 0; i < 16; i++ {
		tmp.bytes[i] = aes.SBOX[state.bytes[i]]
	}
	// State XOR RoundKey
	mm_xor_si128(&tmp, rk)
	MOVOU(state, &tmp)
}

func AESENCLAST(state, rk *XMM) {
	mm_aesenclast_si128(state, rk)
}

func SboxWithAESNI(x, m1l, m1h, m2l, m2h *XMM) {
	y := &XMM{}
	z := &XMM{}
	MOVOU(z, x)
	PAND(z, &const_0f)
	MOVOU(y, m1l)
	PSHUFB(y, z)
	PSRLQ(x, 4)
	PAND(x, &const_0f)
	MOVOU(z, m1h)
	PSHUFB(z, x)
	MOVOU(x, z)
	PXOR(x, y)

	PSHUFB(x, &shift_row_inv)
	AESENCLAST(x, &const_0f)

	MOVOU(z, x)
	PANDN(z, &const_0f)
	MOVOU(y, m2l)
	PSHUFB(y, z)
	PSRLQ(x, 4)
	PAND(x, &const_0f)
	MOVOU(z, m2h)
	PSHUFB(z, x)
	MOVOU(x, z)
	PXOR(x, y)
}

func GenLookupTable(m uint64, c byte, ltl, lth *XMM) {
	mb := Set64(m, m)
	for i := 0; i < 16; i++ {
		ltl.bytes[i] = affineByte(mb.bytes[:8], byte(i), c)
		lth.bytes[i] = affineByte(mb.bytes[:8], byte(i*16), 0)
	}
}
