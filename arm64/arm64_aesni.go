package arm64

import "github.com/emmansun/simd/alg/aes"

func AESE(rk, state *Vector128) {
	var (
		V0 = &Vector128{}
		V1 = &Vector128{}
		V2 = &Vector128{}
		V3 = &Vector128{}
		V4 = &Vector128{}
		V5 = &Vector128{}
		V6 = &Vector128{}
		V7 = &Vector128{}
		V8 = &Vector128{}
		V9 = &Vector128{}
		V10 = &Vector128{}
		V11 = &Vector128{}
		V12 = &Vector128{}
		V13 = &Vector128{}
		V14 = &Vector128{}
		V15 = &Vector128{}
	)
	shiftRow := &Vector128{}
	VLD1_2D([]uint64{0x030E09040F0A0500, 0x0B06010C07020D08}, shiftRow)
	// State XOR RoundKey
	VEOR(rk, state, state)
	// ShiftRows
	VTBL_B(shiftRow, []*Vector128{state}, state)
	// SubBytes
	tmp := &Vector128{}
	r4 := &Vector128{}
	VLD1_16B([]byte{0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40, 0x40}, r4)
	VLD1_16B(aes.SBOX[:], V0)
	VLD1_16B(aes.SBOX[16:], V1)
	VLD1_16B(aes.SBOX[32:], V2)
	VLD1_16B(aes.SBOX[48:], V3)
	VLD1_16B(aes.SBOX[64:], V4)
	VLD1_16B(aes.SBOX[80:], V5)
	VLD1_16B(aes.SBOX[96:], V6)
	VLD1_16B(aes.SBOX[112:], V7)
	VLD1_16B(aes.SBOX[128:], V8)
	VLD1_16B(aes.SBOX[144:], V9)
	VLD1_16B(aes.SBOX[160:], V10)
	VLD1_16B(aes.SBOX[176:], V11)
	VLD1_16B(aes.SBOX[192:], V12)
	VLD1_16B(aes.SBOX[208:], V13)
	VLD1_16B(aes.SBOX[224:], V14)
	VLD1_16B(aes.SBOX[240:], V15)
	VTBL_B(state, []*Vector128{V0, V1, V2, V3}, tmp)
	VSUB_B(r4, state, state)
	VTBX_B(state, []*Vector128{V4, V5, V6, V7}, tmp)
	VSUB_B(r4, state, state)
	VTBX_B(state, []*Vector128{V8, V9, V10, V11}, tmp)
	VSUB_B(r4, state, state)
	VTBX_B(state, []*Vector128{V12, V13, V14, V15}, tmp)
	copy(state.bytes[:], tmp.bytes[:])
}

func SboxWithAESNI(m1l, m1h, m2l, m2h, x *Vector128) {
	var (
		nibble_mask     = &Vector128{}
		inverseShiftRow = &Vector128{}
		y               = &Vector128{}
		z               = &Vector128{}
		zero            = &Vector128{}
	)
	VDUP_BYTE(0, zero)
	VDUP_BYTE(0x0f, nibble_mask)
	VLD1_2D([]uint64{0x0B0E0104070A0D00, 0x0306090C0F020508}, inverseShiftRow)

	VAND(x, nibble_mask, z)
	VTBL_B(z, []*Vector128{m1l}, y)
	VUSHR_D(4, x, x)
	VAND(x, nibble_mask, z)
	VTBL_B(z, []*Vector128{m1h}, z)
	VEOR(y, z, x)

	VTBL_B(inverseShiftRow, []*Vector128{x}, x)

	AESE(zero, x)

	VAND(x, nibble_mask, z)
	VTBL_B(z, []*Vector128{m2l}, y)
	VUSHR_D(4, x, x)
	VAND(x, nibble_mask, z)
	VTBL_B(z, []*Vector128{m2h}, z)
	VEOR(y, z, x)
}

func GenLookupTable(m uint64, c byte, ltl, lth *Vector128) {
	mb := &Vector128{}
	VLD1_2D([]uint64{m, m}, mb)
	for i := 0; i < 16; i++ {
		ltl.bytes[i] = affineByte(mb.bytes[:8], byte(i), c)
		lth.bytes[i] = affineByte(mb.bytes[:8], byte(i*16), 0)
	}
}

// parity(x) = 1 if x has an odd number of 1s in it, and 0 otherwise.
func parity(x byte) byte {
	var t byte
	for i := 0; i < 8; i++ {
		t ^= x & 1
		x >>= 1
	}
	return t
}

func affineByte(tsrc2qw []byte, src1byte, imm byte) byte {
	var retbyte byte
	for i := 0; i < 8; i++ {
		retbyte |= ((parity(tsrc2qw[7-i] & src1byte)) ^ ((imm >> i) & 1)) << i
	}
	return retbyte
}
