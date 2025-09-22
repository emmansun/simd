package avx2

import (
	"encoding/binary"

	"github.com/emmansun/simd/amd64/sse"
)

type YMM struct {
	bytes [32]byte
}

func (m *YMM) Bytes() []byte {
	return m.bytes[:]
}

func VMOVDQU_Luint8(dst *YMM, src []byte) {
	copy(dst.Bytes(), src)
}

func VMOVDQU_Luint16(dst *YMM, src []uint16) {
	_ = src[15] // bounds check hint to compiler; see golang.org/issue/14808
	binary.LittleEndian.PutUint16(dst.Bytes()[0:], src[0])
	binary.LittleEndian.PutUint16(dst.Bytes()[2:], src[1])
	binary.LittleEndian.PutUint16(dst.Bytes()[4:], src[2])
	binary.LittleEndian.PutUint16(dst.Bytes()[6:], src[3])
	binary.LittleEndian.PutUint16(dst.Bytes()[8:], src[4])
	binary.LittleEndian.PutUint16(dst.Bytes()[10:], src[5])
	binary.LittleEndian.PutUint16(dst.Bytes()[12:], src[6])
	binary.LittleEndian.PutUint16(dst.Bytes()[14:], src[7])
	binary.LittleEndian.PutUint16(dst.Bytes()[16:], src[8])
	binary.LittleEndian.PutUint16(dst.Bytes()[18:], src[9])
	binary.LittleEndian.PutUint16(dst.Bytes()[20:], src[10])
	binary.LittleEndian.PutUint16(dst.Bytes()[22:], src[11])
	binary.LittleEndian.PutUint16(dst.Bytes()[24:], src[12])
	binary.LittleEndian.PutUint16(dst.Bytes()[26:], src[13])
	binary.LittleEndian.PutUint16(dst.Bytes()[28:], src[14])
	binary.LittleEndian.PutUint16(dst.Bytes()[30:], src[15])
}

func VMOVDQU_Luint32(dst *YMM, src []uint32) {
	_ = src[7]
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], src[i])
	}
}

func VMOVDQU_Luint64(dst *YMM, src []uint64) {
	_ = src[3]
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], src[i])
	}
}

func VMOVEDQU_Suint8(dst []byte, src *YMM) {
	copy(dst, src.Bytes())
}

func VMOVEDQU_Suint16(dst []uint16, src *YMM) {
	_ = dst[15] // bounds check hint to compiler; see golang.org/issue/14808
	for i := 0; i < 16; i++ {
		dst[i] = binary.LittleEndian.Uint16(src.Bytes()[i*2:])
	}
}

func sign16Extend32(b uint16) int32 {
	a := uint32(b)
	a <<= 16
	return int32(a) >> 16
}

// _mm256_mulhi_epi16
// Multiply the packed signed 16-bit integers in a and b, producing intermediate 32-bit integers,
// and store the high 16 bits of the intermediate integers in dst.
func VPMULHW(dst, src1, src2 *YMM) {
	var temp [4]byte
	result := &YMM{}
	for i := 0; i < 16; i++ {
		a := binary.LittleEndian.Uint16(src1.Bytes()[i*2:])
		b := binary.LittleEndian.Uint16(src2.Bytes()[i*2:])
		prod := sign16Extend32(a) * sign16Extend32(b)

		binary.LittleEndian.PutUint32(temp[:], uint32(prod))
		copy(result.Bytes()[i*2:], temp[2:])
	}
	copy(dst.Bytes(), result.Bytes())
}

// _mm256_mullo_epi16
// Multiply the packed signed 16-bit integers in a and b, producing intermediate 32-bit integers,
// and store the low 16 bits of the intermediate integers in dst.
func VPMULLW(dst, src1, src2 *YMM) {
	result := &YMM{}
	for i := 0; i < 16; i++ {
		a := binary.LittleEndian.Uint16(src1.Bytes()[i*2:])
		b := binary.LittleEndian.Uint16(src2.Bytes()[i*2:])
		prod := sign16Extend32(a) * sign16Extend32(b)

		binary.LittleEndian.PutUint16(result.Bytes()[i*2:], uint16(prod))
	}
	copy(dst.Bytes(), result.Bytes())
}

// _mm256_mulhrs_epi16
// Multiply packed signed 16-bit integers in a and b, producing intermediate signed 32-bit integers.
// Truncate each intermediate integer to the 18 most significant bits, round by adding 1, and store bits [16:1] to dst.
func VPMULHRS(dst, src1, src2 *YMM) {
	result := &YMM{}
	for i := 0; i < 16; i++ {
		a := binary.LittleEndian.Uint16(src1.Bytes()[i*2:])
		b := binary.LittleEndian.Uint16(src2.Bytes()[i*2:])
		prod := sign16Extend32(a) * sign16Extend32(b)

		prod >>= 14
		prod += 1
		prod >>= 1

		binary.LittleEndian.PutUint16(result.Bytes()[i*2:], uint16(prod))
	}
	copy(dst.Bytes(), result.Bytes())
}

func VPAND(dst, src1, src2 *YMM) {
	for i := 0; i < 32; i++ {
		dst.Bytes()[i] = src1.Bytes()[i] & src2.Bytes()[i]
	}
}

func saturateU8(a int16) uint8 {
	if a < 0 {
		return 0
	} else if a > 255 {
		return 255
	}
	return uint8(a)
}

// _mm256_packus_epi16
// Convert packed signed 16-bit integers from a and b to packed 8-bit integers using unsigned saturation,
// and store the results in dst.
func VPACKUSWB(dst, src1, src2 *YMM) {
	result := &YMM{}
	for i := 0; i < 8; i++ {
		a := int16(binary.LittleEndian.Uint16(src1.Bytes()[i*2:]))
		b := int16(binary.LittleEndian.Uint16(src2.Bytes()[i*2:]))
		result.Bytes()[i] = saturateU8(a)
		result.Bytes()[i+8] = saturateU8(b)
	}
	for i := 8; i < 16; i++ {
		a := int16(binary.LittleEndian.Uint16(src1.Bytes()[i*2:]))
		b := int16(binary.LittleEndian.Uint16(src2.Bytes()[i*2:]))
		result.Bytes()[i+8] = saturateU8(a)
		result.Bytes()[i+16] = saturateU8(b)
	}
	copy(dst.Bytes(), result.Bytes())
}

func saturate16(a int32) int16 {
	if a < -32768 {
		return -32768
	} else if a > 32767 {
		return 32767
	}

	return int16(a)
}

// _mm256_maddubs_epi16
// Vertically multiply each unsigned 8-bit integer from a with the corresponding signed 8-bit integer from b,
// producing intermediate signed 16-bit integers. Horizontally add adjacent pairs of intermediate signed 16-bit integers,
// and pack the saturated results in dst.
func VPMADDUBSW(dst, src1, src2 *YMM) {
	for i := 0; i < 16; i++ {
		a0 := uint8(src1.Bytes()[i*2])
		a1 := uint8(src1.Bytes()[i*2+1])
		b0 := int8(src2.Bytes()[i*2])
		b1 := int8(src2.Bytes()[i*2+1])

		// Multiply and add
		sum := int32(a0)*int32(b0) + int32(a1)*int32(b1)
		binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], uint16(saturate16(sum)))
	}
}

// _mm256_permutevar8x32_epi32
// Permute 32-bit integers in a within 256-bit blocks using the control in idx, and store the results in dst.
func VPERMD(dst, src, idx *YMM) {
	result := &YMM{}
	for i := 0; i < 8; i++ {
		j := binary.LittleEndian.Uint32(idx.Bytes()[i*4:]) & 0x07
		copy(result.Bytes()[i*4:], src.Bytes()[j*4:])
	}
	copy(dst.Bytes(), result.Bytes())
}

func VPSHUFB(dst, src1, src2 *YMM) {
	tmp := YMM{}
	tmpBytes := tmp.Bytes()
	src1Bytes := src1.Bytes()
	maskBytes := src2.Bytes()
	for i := 0; i < 16; i++ {
		if maskBytes[i]&0x80 == 0x80 {
			tmpBytes[i] = 0
		} else {
			idx := maskBytes[i] & 0xf
			tmpBytes[i] = src1Bytes[idx]
		}
	}
	for i := 16; i < 32; i++ {
		if maskBytes[i]&0x80 == 0x80 {
			tmpBytes[i] = 0
		} else {
			idx := maskBytes[i] & 0xf
			tmpBytes[i] = src1Bytes[16+idx]
		}
	}
	copy(dst.Bytes(), tmp.Bytes())
}

func SetOneInt16(dst *YMM, a int16) {
	for i := 0; i < 16; i++ {
		binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], uint16(a))
	}
}

func SetOneInt32(dst *YMM, a int32) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], uint32(a))
	}
}

func SetOneInt64(dst *YMM, a int64) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], uint64(a))
	}
}

// _mm256_add_epi16
// Add packed 16-bit integers in a and b, and store the results in dst.
func VPADDW(dst, src1, src2 *YMM) {
	for i := 0; i < 16; i++ {
		a := binary.LittleEndian.Uint16(src1.Bytes()[i*2:])
		b := binary.LittleEndian.Uint16(src2.Bytes()[i*2:])
		binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], a+b)
	}
}

// _mm256_sub_epi16
// Subtract packed 16-bit integers in b from packed 16-bit integers in a, and store the results in dst.
func VPSUBW(dst, src1, src2 *YMM) {
	for i := 0; i < 16; i++ {
		a := int16(binary.LittleEndian.Uint16(src1.Bytes()[i*2:]))
		b := int16(binary.LittleEndian.Uint16(src2.Bytes()[i*2:]))
		binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], uint16(a-b))
	}
}

// _mm256_slli_epi16
// Shift packed 16-bit integers in a left by imm bits while shifting in zeros, and store the results in dst.
func VPSLLW(dst, src *YMM, imm byte) {
	for i := 0; i < 16; i++ {
		w := binary.LittleEndian.Uint16(src.Bytes()[i*2:])
		if imm > 15 {
			dst.Bytes()[i*2] = 0
		}
		binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], w<<imm)
	}
}

// _mm256_srli_epi16
// Shift packed 16-bit integers in a right by imm bits while shifting in zeros, and store the results in dst.
func VPSRLW(dst, src *YMM, imm byte) {
	for i := 0; i < 16; i++ {
		w := binary.LittleEndian.Uint16(src.Bytes()[i*2:])
		if imm > 15 {
			binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], 0)
		}
		binary.LittleEndian.PutUint16(dst.Bytes()[i*2:], w>>imm)
	}
}

// _mm256_andnot_si256
// Compute the bitwise NOT of packed 32-bit integers in a, AND with packed 32-bit integers in b, and store the results in dst.
func VPANDN(dst, src1, src2 *YMM) {
	for i := 0; i < 32; i++ {
		dst.Bytes()[i] = ^src1.Bytes()[i] & src2.Bytes()[i]
	}
}

// _mm256_madd_epi16
// Multiply packed signed 16-bit integers in a and b, producing intermediate signed 32-bit integers.
// Horizontally add adjacent pairs of intermediate signed 32-bit integers, and store the results in dst.
func VPMADDWD(dst, src1, src2 *YMM) {
	for i := 0; i < 8; i++ {
		a0 := int16(binary.LittleEndian.Uint16(src1.Bytes()[i*4:]))
		a1 := int16(binary.LittleEndian.Uint16(src1.Bytes()[i*4+2:]))
		b0 := int16(binary.LittleEndian.Uint16(src2.Bytes()[i*4:]))
		b1 := int16(binary.LittleEndian.Uint16(src2.Bytes()[i*4+2:]))
		// Multiply and add
		sum := int32(a0)*int32(b0) + int32(a1)*int32(b1)
		binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], uint32(sum))
	}
}

// _mm256_sllv_epi32
// Shift packed 32-bit integers in a left by the amount specified in the corresponding 32-bit integer in count, and store the results in dst.
func VPSLLVD(dst, src1, src2 *YMM) {
	for i := 0; i < 8; i++ {
		a := binary.LittleEndian.Uint32(src1.Bytes()[i*4:])
		b := binary.LittleEndian.Uint32(src2.Bytes()[i*4:])
		if b >= 32 {
			binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], 0)
		} else {
			binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], a<<b)
		}
	}
}

// _mm256_srlv_epi32
// Shift packed 32-bit integers in a right by the amount specified by the corresponding element in count while shifting in zeros, and store the results in dst.
func VPSRLVD(dst, a, count *YMM) {
	for i := 0; i < 8; i++ {
		ai := binary.LittleEndian.Uint32(a.Bytes()[i*4:])
		bi := binary.LittleEndian.Uint32(count.Bytes()[i*4:])
		if bi >= 32 {
			binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], 0)
		} else {
			binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], ai>>bi)
		}
	}
}

// _mm256_srlv_epi64
// Shift packed 64-bit integers in a right by the amount specified by the corresponding element in count while shifting in zeros, and store the results in dst.
func VPSRLVQ(dst, src1, src2 *YMM) {
	for i := 0; i < 4; i++ {
		a := binary.LittleEndian.Uint64(src1.Bytes()[i*8:])
		b := binary.LittleEndian.Uint64(src2.Bytes()[i*8:])
		if b >= 64 {
			binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], 0)
		} else {
			binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], a>>b)
		}
	}
}

// _mm256_srli_epi64
// Shift packed 64-bit integers in a right by imm bits while shifting in zeros, and store the results in dst.
func VPSRLQ(dst, src *YMM, imm byte) {
	for i := 0; i < 4; i++ {
		q := binary.LittleEndian.Uint64(src.Bytes()[i*8:])
		if imm > 63 {
			binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], 0)
		} else {
			binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], q>>imm)
		}
	}
}

// _mm256_slli_epi64
// Shift packed 64-bit integers in a left by imm8 while shifting in zeros, and store the results in dst.
func VPSLLQ(dst, src *YMM, imm8 byte) {
	for i := 0; i < 4; i++ {
		q := binary.LittleEndian.Uint64(src.Bytes()[i*8:])
		if imm8 > 63 {
			binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], 0)
		} else {
			binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], q<<imm8)
		}
	}
}

func CastToXMM(dst *sse.XMM, src *YMM) {
	copy(dst.Bytes(), src.Bytes()[:16])
}

func ExtractXMM(dst *sse.XMM, src *YMM, imm uint) {
	if imm&1 == 0 {
		copy(dst.Bytes(), src.Bytes()[:16])
	} else {
		copy(dst.Bytes(), src.Bytes()[16:])
	}
}

// _mm256_permute4x64_epi64
// Permute 64-bit integers in a using the control in imm8, and store the results in dst.
func VPERMQ(dst, a *YMM, imm8 byte) {
	result := &YMM{}
	for i := 0; i < 4; i++ {
		j := (imm8 >> (i * 2)) & 0x03
		copy(result.Bytes()[i*8:], a.Bytes()[j*8:])
	}
	copy(dst.Bytes(), result.Bytes())
}

// _mm256_bsrli_epi128
// Shift 128-bit lanes in a right by imm8 bytes while shifting in zeros, and store the results in dst.
func VPSRLDQ(dst, src *YMM, imm8 byte) {
	if imm8 > 15 {
		imm8 = 16
	}
	if imm8 == 16 {
		for i := 0; i < 32; i++ {
			dst.Bytes()[i] = 0
		}
	} else {
		tmp := &YMM{}
		for i := imm8; i < 16; i++ {
			tmp.Bytes()[i-imm8] = src.Bytes()[i]
			tmp.Bytes()[i-imm8+16] = src.Bytes()[i+16]
		}
		copy(dst.Bytes(), tmp.Bytes())
	}
}

// _mm256_add_epi64
// Add packed 64-bit integers in a and b, and store the results in dst.
func VPADDQ(dst, a, b *YMM) {
	for i := 0; i < 4; i++ {
		ai := binary.LittleEndian.Uint64(a.Bytes()[i*8:])
		bi := binary.LittleEndian.Uint64(b.Bytes()[i*8:])
		binary.LittleEndian.PutUint64(dst.Bytes()[i*8:], ai+bi)
	}
}
