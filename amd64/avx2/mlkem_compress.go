// References
// ==========
//
// - [REF_AVX2]
//   CRYSTALS-Kyber optimized AVX2 implementation
//   Bos, Ducas, Kiltz, Lepoint, Lyubashevsky, Schanck, Schwabe, Seiler, Stehlé
//   https://github.com/pq-crystals/kyber/tree/main/avx2
//
//
// This file is derived from the public domain
// AVX2 Kyber implementation @[REF_AVX2].

package avx2

import (
	"math/big"
	"strconv"

	"github.com/emmansun/simd/amd64/sse"
)

const (
	barrettMultiplier = 5039 // 2¹² * 2¹² / q
	barrettShift      = 24   // log₂(2¹² * 2¹²)
	n                 = 256
	q                 = 3329
)

// compress maps a field element uniformly to the range 0 to 2ᵈ-1, according to
// FIPS 203, Definition 4.7.
func compress(x uint16, d uint8) uint16 {
	// We want to compute (x * 2ᵈ) / q, rounded to nearest integer, with 1/2
	// rounding up (see FIPS 203, Section 2.3).

	// Barrett reduction produces a quotient and a remainder in the range [0, 2q),
	// such that dividend = quotient * q + remainder.
	dividend := uint32(x) << d // x * 2ᵈ
	quotient := uint32(uint64(dividend) * barrettMultiplier >> barrettShift)
	remainder := dividend - quotient*q

	// Since the remainder is in the range [0, 2q), not [0, q), we need to
	// portion it into three spans for rounding.
	//
	//     [ 0,       q/2     ) -> round to 0
	//     [ q/2,     q + q/2 ) -> round to 1
	//     [ q + q/2, 2q      ) -> round to 2
	//
	// We can convert that to the following logic: add 1 if remainder > q/2,
	// then add 1 again if remainder > q + q/2.
	//
	// Note that if remainder > x, then ⌊x⌋ - remainder underflows, and the top
	// bit of the difference will be set.
	quotient += (q/2 - remainder) >> 31 & 1
	quotient += (q + q/2 - remainder) >> 31 & 1

	// quotient might have overflowed at this point, so reduce it by masking.
	var mask uint32 = (1 << d) - 1
	return uint16(quotient & mask)
}

// decompress maps a number x between 0 and 2ᵈ-1 uniformly to the full range of
// field elements, according to FIPS 203, Definition 4.8.
func decompress(y uint16, d uint8) uint16 {
	// We want to compute (y * q) / 2ᵈ, rounded to nearest integer, with 1/2
	// rounding up (see FIPS 203, Section 2.3).

	dividend := uint32(y) * q
	quotient := dividend >> d // (y * q) / 2ᵈ

	// The d'th least-significant bit of the dividend (the most significant bit
	// of the remainder) is 1 for the top half of the values that divide to the
	// same quotient, which are the ones that round up.
	quotient += dividend >> (d - 1) & 1

	// quotient is at most (2¹¹-1) * q / 2¹¹ + 1 = 3328, so it didn't overflow.
	return uint16(quotient)
}

func CompressRat(x uint16, d uint8) uint16 {
	if x >= q {
		panic("x out of range")
	}
	if d <= 0 || d >= 12 {
		panic("d out of range")
	}

	precise := big.NewRat((1<<d)*int64(x), q) // (2ᵈ / q) * x == (2ᵈ * x) / q

	// FloatString rounds halves away from 0, and our result should always be positive,
	// so it should work as we expect. (There's no direct way to round a Rat.)
	rounded, err := strconv.ParseInt(precise.FloatString(0), 10, 64)
	if err != nil {
		panic(err)
	}

	// If we rounded up, `rounded` may be equal to 2ᵈ, so we perform a final reduction.
	return uint16(rounded % (1 << d))
}

func DecompressRat(y uint16, d uint8) uint16 {
	if y >= 1<<d {
		panic("y out of range")
	}
	if d <= 0 || d >= 12 {
		panic("d out of range")
	}

	precise := big.NewRat(q*int64(y), 1<<d) // (q / 2ᵈ) * y  ==  (q * y) / 2ᵈ

	// FloatString rounds halves away from 0, and our result should always be positive,
	// so it should work as we expect. (There's no direct way to round a Rat.)
	rounded, err := strconv.ParseInt(precise.FloatString(0), 10, 64)
	if err != nil {
		panic(err)
	}

	// If we rounded up, `rounded` may be equal to q, so we perform a final reduction.
	return uint16(rounded % q)
}
func ringCompressAndEncode(s []byte, f [n]uint16, d uint8, compressFunc func(x uint16, d uint8) uint16) []byte {
	var b byte
	var bIdx uint8
	for i := range n {
		c := compressFunc(f[i], d)
		var cIdx uint8
		for cIdx < d {
			b |= byte(c>>cIdx) << bIdx
			bits := min(8-bIdx, d-cIdx)
			bIdx += bits
			cIdx += bits
			if bIdx == 8 {
				s = append(s, b)
				b = 0
				bIdx = 0
			}
		}
	}
	if bIdx != 0 {
		panic("mlkem: internal error: bitsFilled != 0")
	}
	return s
}

func ringDecodeAndDecompress(b []byte, d uint8, decompressFunc func(y uint16, d uint8) uint16) [n]uint16 {
	var f [n]uint16
	var bIdx uint8
	for i := 0; i < n; i++ {
		var c uint16
		var cIdx uint8
		for cIdx < d {
			c |= uint16(b[0]>>bIdx) << cIdx
			c &= (1 << d) - 1
			bits := min(8-bIdx, d-cIdx)
			bIdx += bits
			cIdx += bits
			if bIdx == 8 {
				b = b[1:]
				bIdx = 0
			}
		}
		f[i] = uint16(decompressFunc(c, d))
	}
	if len(b) != 0 {
		panic("mlkem: internal error: leftover bytes")
	}
	return f
}

func ringCompressAndEncode4(out []byte, f [256]uint16) {
	var f0, f1, f2, f3, v, shift1, mask, shift2, permdidx YMM
	SetOneInt16(&shift1, 1<<9)
	SetOneInt16(&mask, 0x0F)
	SetOneInt16(&shift2, 16<<8+1)
	SetOneInt16(&v, 20159) // floor(2^26/q + 0.5)
	VMOVDQU_Luint32(&permdidx, []uint32{
		0, 4, 1, 5, 2, 6, 3, 7,
	})
	for i := 0; i < n/64; i++ {
		VMOVDQU_Luint16(&f0, f[i*64+0:])
		VMOVDQU_Luint16(&f1, f[i*64+16:])
		VMOVDQU_Luint16(&f2, f[i*64+32:])
		VMOVDQU_Luint16(&f3, f[i*64+48:])
		VPMULHW(&f0, &f0, &v)
		VPMULHW(&f1, &f1, &v)
		VPMULHW(&f2, &f2, &v)
		VPMULHW(&f3, &f3, &v)
		VPMULHRS(&f0, &f0, &shift1)
		VPMULHRS(&f1, &f1, &shift1)
		VPMULHRS(&f2, &f2, &shift1)
		VPMULHRS(&f3, &f3, &shift1)
		VPAND(&f0, &f0, &mask)
		VPAND(&f1, &f1, &mask)
		VPAND(&f2, &f2, &mask)
		VPAND(&f3, &f3, &mask)
		VPACKUSWB(&f0, &f0, &f1)
		VPACKUSWB(&f2, &f2, &f3)
		VPMADDUBSW(&f0, &f0, &shift2)
		VPMADDUBSW(&f2, &f2, &shift2)
		VPACKUSWB(&f0, &f0, &f2)
		VPERMD(&f0, &f0, &permdidx)
		VMOVEDQU_Suint8(out[i*32:], &f0)
	}
}

func ringDecodeAndDecompress4(out *[256]uint16, in []byte) {
	var q1, shufbidx, mask, shift, f YMM
	SetOneInt16(&q1, q) // 3329
	VMOVDQU_Luint8(&shufbidx, []byte{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7})
	SetOneInt32(&mask, 0x00F0000F)
	SetOneInt32(&shift, (128<<16)+2048)

	for i := 0; i < n/16; i++ {
		copy(f.Bytes(), in[i*8:(i+1)*8])
		copy(f.Bytes()[8:], []byte{0, 0, 0, 0, 0, 0, 0, 0})
		copy(f.Bytes()[16:], f.Bytes()[0:16])
		VPSHUFB(&f, &f, &shufbidx)
		VPAND(&f, &f, &mask)
		VPMULLW(&f, &f, &shift)
		VPMULHRS(&f, &f, &q1)
		VMOVEDQU_Suint16(out[i*16:], &f)
	}
}

func ringCompressAndEncode5(out []byte, f [256]uint16) {
	var v, shift1, mask, shift2, shift3, sllvdidx, shufbidx, f0, f1 YMM
	var t0, t1, t3 sse.XMM
	SetOneInt16(&v, 20159) // floor(2^26/q + 0.5)
	SetOneInt16(&shift1, 1<<10)
	SetOneInt16(&mask, 0x1F)
	SetOneInt16(&shift2, 32<<8+1)
	SetOneInt32(&shift3, 1024<<16+1)
	SetOneInt64(&sllvdidx, 12)
	VMOVDQU_Luint8(&shufbidx, []byte{0, 1, 2, 3, 4, 255, 255, 255, 255, 255, 8, 9, 10, 11, 12, 255, 9, 10, 11, 12, 255, 0, 1, 2, 3, 4, 255, 255, 255, 255, 255, 8})

	for i := 0; i < n/32; i++ {
		VMOVDQU_Luint16(&f0, f[i*32:])
		VMOVDQU_Luint16(&f1, f[i*32+16:])

		VPMULHW(&f0, &f0, &v)
		VPMULHRS(&f0, &f0, &shift1)
		VPAND(&f0, &f0, &mask)

		VPMULHW(&f1, &f1, &v)
		VPMULHRS(&f1, &f1, &shift1)
		VPAND(&f1, &f1, &mask)

		VPACKUSWB(&f0, &f0, &f1)

		VPMADDUBSW(&f0, &f0, &shift2)

		VPMADDWD(&f0, &f0, &shift3)

		VPSLLVD(&f0, &f0, &sllvdidx)
		VPSRLVQ(&f0, &f0, &sllvdidx)

		VPSHUFB(&f0, &f0, &shufbidx)

		CastToXMM(&t0, &f0)
		ExtractXMM(&t1, &f0, 1)
		CastToXMM(&t3, &shufbidx)
		sse.PBLENDVB(&t0, &t0, &t1, &t3)

		copy(out[i*20:], t0.Bytes())
		copy(out[i*20+16:], t1.Bytes()[:4])
	}
}

func ringDecodeAndDecompress5(out *[256]uint16, in []byte) {
	var q1, shufbidx, shift, mask, f YMM
	SetOneInt16(&q1, q)
	VMOVDQU_Luint8(&shufbidx, []byte{0, 0, 0, 1, 1, 1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 8, 8, 8, 8, 9, 9, 9})
	VMOVDQU_Luint16(&mask, []uint16{31, 992, 124, 3968, 496, 62, 1984, 248, 31, 992, 124, 3968, 496, 62, 1984, 248})
	VMOVDQU_Luint16(&shift, []uint16{1024, 32, 256, 8, 64, 512, 16, 128, 1024, 32, 256, 8, 64, 512, 16, 128})

	for i := 0; i < n/16; i++ {
		copy(f.Bytes(), in[i*10:(i+1)*10])
		copy(f.Bytes()[10:], []byte{0, 0, 0, 0, 0, 0})
		copy(f.Bytes()[16:], f.Bytes()[0:16])

		VPSHUFB(&f, &f, &shufbidx) // f = _mm256_shuffle_epi8(f, shufbidx);
		VPAND(&f, &f, &mask)       // f = _mm256_and_si256(f, mask);
		VPMULLW(&f, &f, &shift)    // f = _mm256_mullo_epi16(f, shift);
		VPMULHRS(&f, &f, &q1)      // f = _mm256_mulhrs_epi16(f, q);

		VMOVEDQU_Suint16(out[i*16:], &f)
	}
}

func ringCompressAndEncode10(out []byte, f [256]uint16) {
	var f0, f1, f2, v, v8, off, shift1, mask, shift2, sllvdidx, shufbidx YMM
	var t0, t1 sse.XMM
	SetOneInt16(&v, 20159) // floor(2^26/q + 0.5)
	VPSLLW(&v8, &v, byte(3))
	SetOneInt16(&off, 0xf)
	SetOneInt16(&shift1, 1<<12)
	SetOneInt16(&mask, 1023)
	SetOneInt64(&shift2, int64(uint64(1024)<<48+uint64(1)<<32+uint64(1024)<<16+uint64(1)))
	SetOneInt64(&sllvdidx, 12)
	VMOVDQU_Luint8(&shufbidx, []byte{0, 1, 2, 3, 4, 8, 9, 10, 11, 12, 255, 255, 255, 255, 255, 255, 9, 10, 11, 12, 255, 255, 255, 255, 255, 255, 0, 1, 2, 3, 4, 8})
	for i := 0; i < n/16; i++ {
		VMOVDQU_Luint16(&f0, f[i*16:])
		VPMULLW(&f1, &f0, &v8)           // f1 = _mm256_mullo_epi16(f0, v8);
		VPADDW(&f2, &f0, &off)           // f2 = _mm256_add_epi16(f0, off);
		VPSLLW(&f0, &f0, byte(3))        // f0 = _mm256_slli_epi16(f0, 3);
		VPMULHW(&f0, &f0, &v)            // f0 = _mm256_mulhi_epi16(f0, v);
		VPSUBW(&f2, &f1, &f2)            // f2 = _mm256_sub_epi16(f1, f2);
		VPANDN(&f1, &f1, &f2)            // f1 = _mm256_andnot_si256(f1, f2);
		VPSRLW(&f1, &f1, byte(15))       // f1 = _mm256_srli_epi16(f1, 15);
		VPSUBW(&f0, &f0, &f1)            // f0 = _mm256_sub_epi16(f0, f1);
		VPMULHRS(&f0, &f0, &shift1)      // f0 = _mm256_mulhrs_epi16(f0, shift1);
		VPAND(&f0, &f0, &mask)           // f0 = _mm256_and_si256(f0, mask);
		VPMADDWD(&f0, &f0, &shift2)      // f0 = _mm256_madd_epi16(f0, shift2);
		VPSLLVD(&f0, &f0, &sllvdidx)     // f0 = _mm256_sllv_epi32(f0, sllvdidx);
		VPSRLQ(&f0, &f0, byte(12))       // f0 = _mm256_srli_epi64(f0, 12);
		VPSHUFB(&f0, &f0, &shufbidx)     // f0 = _mm256_shuffle_epi8(f0, shufbidx);
		CastToXMM(&t0, &f0)              // t0 = _mm256_castsi256_si128(f0);
		ExtractXMM(&t1, &f0, 1)          // t1 = _mm256_extracti128_si256(f0, 1);
		sse.PBLENDW(&t0, &t0, &t1, 0xE0) // t0 = _mm_blend_epi16(t0, t1, 0xE0);
		copy(out[i*20:], t0.Bytes())
		copy(out[i*20+16:], t1.Bytes()[:4])
	}
}

func ringDecodeAndDecompress10(out *[256]uint16, in []byte) {
	var q1, shufbidx, sllvdidx, mask, f YMM
	SetOneInt32(&q1, 4*q+q<<16) // 4*q + q*2^16 = 3329*65537
	VMOVDQU_Luint8(&shufbidx, []byte{0, 1, 1, 2, 2, 3, 3, 4, 5, 6, 6, 7, 7, 8, 8, 9, 2, 3, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10, 10, 11})
	SetOneInt64(&sllvdidx, 4)
	SetOneInt32(&mask, 32736<<16+8184)

	for i := 0; i < n/16-1; i++ {
		VMOVDQU_Luint8(&f, in[i*20:])
		VPERMQ(&f, &f, 0x94)
		VPSHUFB(&f, &f, &shufbidx)
		VPSLLVD(&f, &f, &sllvdidx)
		VPSRLW(&f, &f, byte(1))
		VPAND(&f, &f, &mask)
		VPMULHRS(&f, &f, &q1)
		VMOVEDQU_Suint16(out[i*16:], &f)
	}
	/* Handle load in last iteration especially to avoid buffer overflow */
	var rest [32]byte
	copy(rest[:], in[300:])
	VMOVDQU_Luint8(&f, rest[:])
	VPERMQ(&f, &f, 0x94)
	VPSHUFB(&f, &f, &shufbidx)
	VPSLLVD(&f, &f, &sllvdidx)
	VPSRLW(&f, &f, byte(1))
	VPAND(&f, &f, &mask)
	VPMULHRS(&f, &f, &q1)
	VMOVEDQU_Suint16(out[15*16:], &f)
}

func ringCompressAndEncode11(out []byte, f [256]uint16) {
	var f0, f1, f2 YMM
	var t0, t1, t2 sse.XMM
	var v, v8, off, shift1, mask, shift2, sllvdidx, srlvqidx, shufbidx YMM
	SetOneInt16(&v, 20159)
	VPSLLW(&v8, &v, byte(3))
	SetOneInt16(&off, 36)
	SetOneInt16(&shift1, 1<<13)
	SetOneInt16(&mask, 2047)
	SetOneInt64(&shift2, int64(uint64(2048)<<48+uint64(1)<<32+uint64(2048)<<16+uint64(1)))
	SetOneInt64(&sllvdidx, 10)
	VMOVDQU_Luint64(&srlvqidx, []uint64{10, 30, 10, 30})
	VMOVDQU_Luint8(&shufbidx, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 255, 255, 255, 255, 255, 5, 6, 7, 8, 9, 10, 255, 255, 255, 255, 0, 0, 1, 2, 3, 4})

	for i := 0; i < n/16-1; i++ {
		VMOVDQU_Luint16(&f0, f[i*16:])
		VPMULLW(&f1, &f0, &v8)       // f1 = _mm256_mullo_epi16(f0, v8);
		VPADDW(&f2, &f0, &off)       // f2 = _mm256_add_epi16(f0, off);
		VPSLLW(&f0, &f0, byte(3))    // f0 = _mm256_slli_epi16(f0, 3);
		VPMULHW(&f0, &f0, &v)        // f0 = _mm256_mulhi_epi16(f0, v);
		VPSUBW(&f2, &f1, &f2)        // f2 = _mm256_sub_epi16(f1, f2);
		VPANDN(&f1, &f1, &f2)        // f1 = _mm256_andnot_si256(f1, f2);
		VPSRLW(&f1, &f1, byte(15))   // f1 = _mm256_srli_epi16(f1, 15);
		VPSUBW(&f0, &f0, &f1)        // f0 = _mm256_sub_epi16(f0, f1);
		VPMULHRS(&f0, &f0, &shift1)  // f0 = _mm256_mulhrs_epi16(f0, shift1);
		VPAND(&f0, &f0, &mask)       // f0 = _mm256_and_si256(f0, mask);
		VPMADDWD(&f0, &f0, &shift2)  // f0 = _mm256_madd_epi16(f0, shift2);
		VPSLLVD(&f0, &f0, &sllvdidx) // f0 = _mm256_sllv_epi32(f0, sllvdidx);
		VPSRLDQ(&f1, &f0, 8)         // f1 = _mm256_bsrli_epi128(f0, 8);
		VPSRLVQ(&f0, &f0, &srlvqidx) // f0 = _mm256_srlv_epi64(f0, srlvqidx);
		VPSLLQ(&f1, &f1, byte(34))   // f1 = _mm256_slli_epi64(f1, 34);
		VPADDQ(&f0, &f0, &f1)        // f0 = _mm256_add_epi64(f0, f1);
		VPSHUFB(&f0, &f0, &shufbidx) // f0 = _mm256_shuffle_epi8(f0, shufbidx);
		CastToXMM(&t0, &f0)          // t0 = _mm256_castsi256_si128(f0);
		ExtractXMM(&t1, &f0, 1)      // t1 = _mm256_extracti128_si256(f0, 1);
		CastToXMM(&t2, &shufbidx)
		sse.PBLENDVB(&t0, &t0, &t1, &t2) // t0 = _mm_blendv_epi8(t0, t1, _mm256_castsi256_si128(shufbidx));
		copy(out[i*22:], t0.Bytes())
		copy(out[i*22+16:], t1.Bytes()[:8])
	}
	VMOVDQU_Luint16(&f0, f[240:])
	VPMULLW(&f1, &f0, &v8)       // f1 = _mm256_mullo_epi16(f0, v8);
	VPADDW(&f2, &f0, &off)       // f2 = _mm256_add_epi16(f0, off);
	VPSLLW(&f0, &f0, byte(3))    // f0 = _mm256_slli_epi16(f0, 3);
	VPMULHW(&f0, &f0, &v)        // f0 = _mm256_mulhi_epi16(f0, v);
	VPSUBW(&f2, &f1, &f2)        // f2 = _mm256_sub_epi16(f1, f2);
	VPANDN(&f1, &f1, &f2)        // f1 = _mm256_andnot_si256(f1, f2);
	VPSRLW(&f1, &f1, byte(15))   // f1 = _mm256_srli_epi16(f1, 15);
	VPSUBW(&f0, &f0, &f1)        // f0 = _mm256_sub_epi16(f0, f1);
	VPMULHRS(&f0, &f0, &shift1)  // f0 = _mm256_mulhrs_epi16(f0, shift1);
	VPAND(&f0, &f0, &mask)       // f0 = _mm256_and_si256(f0, mask);
	VPMADDWD(&f0, &f0, &shift2)  // f0 = _mm256_madd_epi16(f0, shift2);
	VPSLLVD(&f0, &f0, &sllvdidx) // f0 = _mm256_sllv_epi32(f0, sllvdidx);
	VPSRLDQ(&f1, &f0, 8)         // f1 = _mm256_bsrli_epi128(f0, 8);
	VPSRLVQ(&f0, &f0, &srlvqidx) // f0 = _mm256_srlv_epi64(f0, srlvqidx);
	VPSLLQ(&f1, &f1, byte(34))   // f1 = _mm256_slli_epi64(f1, 34);
	VPADDQ(&f0, &f0, &f1)        // f0 = _mm256_add_epi64(f0, f1);
	VPSHUFB(&f0, &f0, &shufbidx) // f0 = _mm256_shuffle_epi8(f0, shufbidx);
	CastToXMM(&t0, &f0)          // t0 = _mm256_castsi256_si128(f0);
	ExtractXMM(&t1, &f0, 1)      // t1 = _mm256_extracti128_si256(f0, 1);
	CastToXMM(&t2, &shufbidx)
	sse.PBLENDVB(&t0, &t0, &t1, &t2) // t0 = _mm_blendv_epi8(t0, t1, _mm256_castsi256_si128(shufbidx));
	copy(out[330:], t0.Bytes())
	copy(out[330+16:], t1.Bytes()[:6])
}

func ringDecodeAndDecompress11(out *[256]uint16, in []byte) {
	var q1, shufbidx, srlvdidx, srlvqidx, shift, mask, f YMM
	SetOneInt16(&q1, q)
	VMOVDQU_Luint8(&shufbidx, []byte{0, 1, 1, 2, 2, 3, 4, 5, 5, 6, 6, 7, 8, 9, 9, 10, 3, 4, 4, 5, 5, 6, 7, 8, 8, 9, 9, 10, 11, 12, 12, 13})
	VMOVDQU_Luint32(&srlvdidx, []uint32{0, 1, 0, 0, 0, 1, 0, 0})
	VMOVDQU_Luint64(&srlvqidx, []uint64{0, 2, 0, 2})
	VMOVDQU_Luint16(&shift, []uint16{32, 4, 1, 32, 8, 1, 32, 4, 32, 4, 1, 32, 8, 1, 32, 4})
	SetOneInt16(&mask, 32752)

	for i := 0; i < n/16-1; i++ {
		VMOVDQU_Luint8(&f, in[i*22:])
		VPERMQ(&f, &f, 0x94)       // f = _mm256_permute4x64_epi64(f, 0x94);
		VPSHUFB(&f, &f, &shufbidx) // f = _mm256_shuffle_epi8(f, shufbidx);
		VPSRLVD(&f, &f, &srlvdidx) // f = _mm256_srlv_epi32(f, srlvdidx);
		VPSRLVQ(&f, &f, &srlvqidx) // f = _mm256_srlv_epi64(f, srlvqidx);
		VPMULLW(&f, &f, &shift)    // f = _mm256_mullo_epi16(f, shift);
		VPSRLW(&f, &f, byte(1))    // f = _mm256_srli_epi16(f, 1);
		VPAND(&f, &f, &mask)       // f = _mm256_and_si256(f, mask);
		VPMULHRS(&f, &f, &q1)      // f = _mm256_mulhrs_epi16(f, q);
		VMOVEDQU_Suint16(out[i*16:], &f)
	}
	var rest [32]byte
	copy(rest[:], in[330:])
	VMOVDQU_Luint8(&f, rest[:])
	VPERMQ(&f, &f, 0x94)       // f = _mm256_permute4x64_epi64(f, 0x94);
	VPSHUFB(&f, &f, &shufbidx) // f = _mm256_shuffle_epi8(f, shufbidx);
	VPSRLVD(&f, &f, &srlvdidx) // f = _mm256_srlv_epi32(f, srlvdidx);
	VPSRLVQ(&f, &f, &srlvqidx) // f = _mm256_srlv_epi64(f, srlvqidx);
	VPMULLW(&f, &f, &shift)    // f = _mm256_mullo_epi16(f, shift);
	VPSRLW(&f, &f, byte(1))    // f = _mm256
	VPAND(&f, &f, &mask)       // f = _mm256_and_si256(f, mask);
	VPMULHRS(&f, &f, &q1)      // f = _mm256_mulhrs_epi16(f, q);
	VMOVEDQU_Suint16(out[15*16:], &f)
}
