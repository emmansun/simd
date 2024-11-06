package ppc64

import "github.com/emmansun/simd/alg/aes"

func VSBOX(src, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i++ {
		tmp.bytes[i] = aes.SBOX[src.bytes[i]]
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func SboxWithAESNI(m1l, m1h, m2l, m2h, x *Vector128) {
	var (
		nibble_mask	 = &Vector128{}
		v_four		 = &Vector128{}
		y = &Vector128{}
		z = &Vector128{}
	)
	LXVD2X_UINT64([]uint64{0x0F0F0F0F0F0F0F0F, 0x0F0F0F0F0F0F0F0F}, nibble_mask)
	VSPLTISB(4, v_four)
	VAND(x, nibble_mask, z)
	VPERM(m1l, m1l, z, y)
	VSRD(x, v_four, x) 
	VAND(x, nibble_mask, z)
	VPERM(m1h, m1h, z, x)
	VXOR(y, x, x)

	VSBOX(x, x)

	VAND(x, nibble_mask, z)
	VPERM(m2l, m2l, z, y)
	VSRD(x, v_four, x)
	VAND(x, nibble_mask, z)
	VPERM(m2h, m2h, z, x)
	VXOR(y, x, x)
}

func GenLookupTable(m uint64, c byte, ltl, lth *Vector128) {
	mb := &Vector128{}
	LXVD2X_UINT64([]uint64{m, m}, mb)
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
		retbyte |= ((parity(tsrc2qw[i] & src1byte)) ^ ((imm >> i) & 1)) << i
	}
	return retbyte
}
