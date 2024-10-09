package amd64

import (
	"encoding/binary"

	"github.com/emmansun/simd/amd64/sse"
)

var bitReverseL = sse.Set64(0x0f070b030d050901, 0x0e060a020c040800)
var bitReverseH = sse.Set64(0xf070b030d0509010, 0xe060a020c0408000)
var nibbleMask = sse.Set64(0x0f0f0f0f0f0f0f0f, 0x0f0f0f0f0f0f0f0f)
var shufMaskDW0_0_DW1_0 = sse.Set64(0xffffffff07060504, 0xffffffff03020100)
var shufMaskDW2_0_DW3_0 = sse.Set64(0xffffffff0f0e0d0c, 0xffffffff0b0a0908)
var bits_32_63 = sse.Set64(0x0000000000000000, 0xffffffff00000000)
var shuf_mask_0_0_dw1_0 = sse.Set64(0xffffffff07060504, 0xffffffffffffffff)
var shuf_mask_0_0_0_dw1 = sse.Set64(0x07060504ffffffff, 0xffffffffffffffff)

func EIA16Bytes(data []byte, keys []uint32) uint32 {
	var (
		XTMP1   = sse.XMM{}
		XTMP2   = sse.XMM{}
		XTMP3   = sse.XMM{}
		XTMP4   = sse.XMM{}
		XDATA   = sse.XMM{}
		XDIGEST = sse.XMM{}
		KS_L    = sse.XMM{}
		KS_M1   = sse.XMM{}
	)
	sse.SetBytes(&XDATA, data)
	sse.MOVOU(&XTMP4, &nibbleMask)
	sse.MOVOU(&XTMP2, &XDATA)
	sse.PAND(&XTMP2, &XTMP4)

	sse.PANDN(&XTMP4, &XDATA)
	sse.PSRLQ(&XTMP4, 4)

	sse.MOVOU(&XTMP3, &bitReverseH)
	sse.PSHUFB(&XTMP3, &XTMP2)
	sse.MOVOU(&XTMP1, &bitReverseL)
	sse.PSHUFB(&XTMP1, &XTMP4)

	sse.PXOR(&XTMP3, &XTMP1) // XTMP3 - bit reverse data bytes

	// ZUC authentication part, 4x32 data bits
	// setup KS
	k1 := sse.SetEpi32(keys[0], keys[1], keys[2], keys[3])
	k2 := sse.SetEpi32(keys[2], keys[3], keys[4], keys[5])
	sse.MOVOU(&XTMP1, &k1)
	sse.MOVOU(&XTMP2, &k2)
	sse.PSHUFD(&KS_L, &XTMP1, 0x61)
	sse.PSHUFD(&KS_M1, &XTMP2, 0x61)

	// setup data
	sse.MOVOU(&XTMP1, &XTMP3)
	sse.PSHUFB(&XTMP1, &shufMaskDW0_0_DW1_0)
	sse.MOVOU(&XTMP2, &XTMP1)
	sse.PSHUFB(&XTMP3, &shufMaskDW2_0_DW3_0)
	sse.MOVOU(&XDIGEST, &XTMP3)

	// clmul
	// xor the results from 4 32-bit words together
	// Calculate lower 32 bits of tag
	sse.PCLMULQDQ(&XTMP1, &KS_L, 0x00)
	sse.PCLMULQDQ(&XTMP2, &KS_L, 0x11)
	sse.PCLMULQDQ(&XDIGEST, &KS_M1, 0x00)
	sse.PCLMULQDQ(&XTMP3, &KS_M1, 0x11)

	// XOR all products and move 32-bits to lower 32 bits
	sse.PXOR(&XTMP2, &XTMP1)
	sse.PXOR(&XDIGEST, &XTMP3)
	sse.PXOR(&XDIGEST, &XTMP2)
	sse.PSRLDQ(&XDIGEST, 4)

	return binary.LittleEndian.Uint32(XDIGEST.Bytes())
}

func EIA256RoundTag8(data []byte, keys []uint32) uint64 {
	var (
		XTMP1   = sse.XMM{}
		XTMP2   = sse.XMM{}
		XTMP3   = sse.XMM{}
		XTMP4   = sse.XMM{}
		XTMP5   = sse.XMM{}
		XTMP6   = sse.XMM{}
		XDATA   = sse.XMM{}
		XDIGEST = sse.XMM{}
		KS_L    = sse.XMM{}
		KS_M1   = sse.XMM{}
		KS_M2   = sse.XMM{}
	)
	sse.SetBytes(&XDATA, data)
	sse.MOVOU(&XTMP4, &nibbleMask)
	sse.MOVOU(&XTMP2, &XDATA)
	sse.PAND(&XTMP2, &XTMP4)

	sse.PANDN(&XTMP4, &XDATA)
	sse.PSRLQ(&XTMP4, 4)

	sse.MOVOU(&XTMP3, &bitReverseH)
	sse.PSHUFB(&XTMP3, &XTMP2)
	sse.MOVOU(&XTMP1, &bitReverseL)
	sse.PSHUFB(&XTMP1, &XTMP4)

	sse.PXOR(&XTMP3, &XTMP1) // XTMP3 - bit reverse data bytes

	// ZUC authentication part, 4x32 data bits
	// setup KS
	k1 := sse.SetEpi32(keys[0], keys[1], keys[2], keys[3])
	k2 := sse.SetEpi32(keys[2], keys[3], keys[4], keys[5])
	k3 := sse.SetEpi32(keys[4], keys[5], keys[6], keys[7])
	sse.MOVOU(&XTMP1, &k1)
	sse.MOVOU(&XTMP2, &k2)
	sse.MOVOU(&XTMP4, &k3)
	sse.PSHUFD(&KS_L, &XTMP1, 0x61)
	sse.PSHUFD(&KS_M1, &XTMP2, 0x61)
	sse.PSHUFD(&KS_M2, &XTMP4, 0x61)

	// setup data
	sse.MOVOU(&XTMP1, &XTMP3)
	sse.PSHUFB(&XTMP1, &shufMaskDW0_0_DW1_0)
	sse.MOVOU(&XTMP2, &XTMP1)
	sse.PSHUFB(&XTMP3, &shufMaskDW2_0_DW3_0)
	sse.MOVOU(&XDIGEST, &XTMP3)

	// clmul
	// xor the results from 4 32-bit words together
	// Save data for following products
	sse.MOVOU(&XTMP5, &XTMP2) //  Data bits [31:0 0s 63:32 0s]
	sse.MOVOU(&XTMP6, &XTMP3) //  Data bits [95:64 0s 127:96 0s]

	sse.PCLMULQDQ(&XTMP1, &KS_L, 0x00)
	sse.PCLMULQDQ(&XTMP2, &KS_L, 0x11)
	sse.PCLMULQDQ(&XDIGEST, &KS_M1, 0x00)
	sse.PCLMULQDQ(&XTMP3, &KS_M1, 0x11)

	sse.PXOR(&XTMP2, &XTMP1)
	sse.PXOR(&XDIGEST, &XTMP3)
	sse.PXOR(&XDIGEST, &XTMP2)
	sse.MOVOU_U64(&XDIGEST, 0, XDIGEST.Uint64s()[0]) // Clear top 64 bits
	sse.PSRLDQ(&XDIGEST, 4)

	// Calculate upper 32 bits of tag
	sse.MOVOU(&XTMP1, &XTMP5)
	sse.MOVOU(&XTMP2, &XTMP5)
	sse.MOVOU(&XTMP3, &XTMP6)
	sse.MOVOU(&XTMP4, &XTMP6)

	sse.PCLMULQDQ(&XTMP1, &KS_L, 0x10)
	sse.PCLMULQDQ(&XTMP2, &KS_M1, 0x01)
	sse.PCLMULQDQ(&XTMP3, &KS_M1, 0x10)
	sse.PCLMULQDQ(&XTMP4, &KS_M2, 0x01)

	// XOR all the products and keep only bits 63-32
	sse.PXOR(&XTMP1, &XTMP2)
	sse.PXOR(&XTMP3, &XTMP4)
	sse.PXOR(&XTMP1, &XTMP3)
	sse.PAND(&XTMP1, &bits_32_63)

	sse.POR(&XDIGEST, &XTMP1)

	return binary.LittleEndian.Uint64(XDIGEST.Bytes())
}

func EIA256RoundTag16(data []byte, keys []uint32) (uint64, uint64) {
	var (
		XTMP1   = sse.XMM{}
		XTMP2   = sse.XMM{}
		XTMP3   = sse.XMM{}
		XTMP4   = sse.XMM{}
		XTMP5   = sse.XMM{}
		XTMP6   = sse.XMM{}
		XDATA   = sse.XMM{}
		XDIGEST = sse.XMM{}
		KS_L    = sse.XMM{}
		KS_M1   = sse.XMM{}
		KS_M2   = sse.XMM{}
		KS_H    = sse.XMM{}
	)
	sse.SetBytes(&XDATA, data)
	sse.MOVOU(&XTMP4, &nibbleMask)
	sse.MOVOU(&XTMP2, &XDATA)
	sse.PAND(&XTMP2, &XTMP4)

	sse.PANDN(&XTMP4, &XDATA)
	sse.PSRLQ(&XTMP4, 4)

	sse.MOVOU(&XTMP3, &bitReverseH)
	sse.PSHUFB(&XTMP3, &XTMP2)
	sse.MOVOU(&XTMP1, &bitReverseL)
	sse.PSHUFB(&XTMP1, &XTMP4)

	sse.PXOR(&XTMP3, &XTMP1) // XTMP3 - bit reverse data bytes

	// ZUC authentication part, 4x32 data bits
	// setup KS
	k1 := sse.SetEpi32(keys[0], keys[1], keys[2], keys[3])
	k2 := sse.SetEpi32(keys[2], keys[3], keys[4], keys[5])
	k3 := sse.SetEpi32(keys[4], keys[5], keys[6], keys[7])
	sse.MOVOU(&XTMP1, &k1)
	sse.MOVOU(&XTMP2, &k2)
	sse.MOVOU(&XTMP4, &k3)
	sse.PSHUFD(&KS_L, &XTMP1, 0x61)
	sse.PSHUFD(&KS_M1, &XTMP2, 0x61)
	sse.PSHUFD(&KS_M2, &XTMP4, 0x61)
	sse.PSHUFD(&KS_H, &XTMP4, 0xbb)

	// setup data
	sse.MOVOU(&XTMP1, &XTMP3)
	sse.PSHUFB(&XTMP1, &shufMaskDW0_0_DW1_0)
	sse.MOVOU(&XTMP2, &XTMP1)
	sse.PSHUFB(&XTMP3, &shufMaskDW2_0_DW3_0)
	sse.MOVOU(&XDIGEST, &XTMP3)

	// clmul
	// xor the results from 4 32-bit words together
	// Save data for following products
	sse.MOVOU(&XTMP5, &XTMP2) //  Data bits [31:0 0s 63:32 0s]
	sse.MOVOU(&XTMP6, &XTMP3) //  Data bits [95:64 0s 127:96 0s]

	sse.PCLMULQDQ(&XTMP1, &KS_L, 0x00)
	sse.PCLMULQDQ(&XTMP2, &KS_L, 0x11)
	sse.PCLMULQDQ(&XDIGEST, &KS_M1, 0x00)
	sse.PCLMULQDQ(&XTMP3, &KS_M1, 0x11)

	sse.PXOR(&XTMP2, &XTMP1)
	sse.PXOR(&XDIGEST, &XTMP3)
	sse.PXOR(&XDIGEST, &XTMP2)
	sse.MOVOU_U64(&XDIGEST, 0, XDIGEST.Uint64s()[0]) // Clear top 64 bits
	sse.PSRLDQ(&XDIGEST, 4)

	// Calculate upper 32 bits of tag
	sse.MOVOU(&XTMP1, &XTMP5)
	sse.MOVOU(&XTMP2, &XTMP5)
	sse.MOVOU(&XTMP3, &XTMP6)
	sse.MOVOU(&XTMP4, &XTMP6)

	sse.PCLMULQDQ(&XTMP1, &KS_L, 0x10)
	sse.PCLMULQDQ(&XTMP2, &KS_M1, 0x01)
	sse.PCLMULQDQ(&XTMP3, &KS_M1, 0x10)
	sse.PCLMULQDQ(&XTMP4, &KS_M2, 0x01)

	// XOR all the products and keep only bits 63-32
	sse.PXOR(&XTMP1, &XTMP2)
	sse.PXOR(&XTMP3, &XTMP4)
	sse.PXOR(&XTMP1, &XTMP3)
	sse.PAND(&XTMP1, &bits_32_63)

	sse.POR(&XDIGEST, &XTMP1)
	
	// Prepare data and calculate bits 95-64 of tag
	sse.MOVOU(&XTMP1, &XTMP5)
	sse.MOVOU(&XTMP2, &XTMP5)
	sse.MOVOU(&XTMP3, &XTMP6)
	sse.MOVOU(&XTMP4, &XTMP6)

	sse.PCLMULQDQ(&XTMP1, &KS_M1, 0x00)
	sse.PCLMULQDQ(&XTMP2, &KS_M1, 0x11)
	sse.PCLMULQDQ(&XTMP3, &KS_M2, 0x00)
	sse.PCLMULQDQ(&XTMP4, &KS_M2, 0x11)

	sse.PXOR(&XTMP1, &XTMP2)
	sse.PXOR(&XTMP3, &XTMP4)
	sse.PXOR(&XTMP1, &XTMP3)
	sse.PSHUFB(&XTMP1, &shuf_mask_0_0_dw1_0)

	sse.POR(&XDIGEST, &XTMP1)

	// Prepare data and calculate bits 127-96 of tag
	sse.MOVOU(&XTMP1, &XTMP5)
	sse.MOVOU(&XTMP2, &XTMP5)
	sse.MOVOU(&XTMP3, &XTMP6)
	sse.MOVOU(&XTMP4, &XTMP6)

	sse.PCLMULQDQ(&XTMP1, &KS_M1, 0x10)
	sse.PCLMULQDQ(&XTMP2, &KS_M2, 0x01)
	sse.PCLMULQDQ(&XTMP3, &KS_M2, 0x10)
	sse.PCLMULQDQ(&XTMP4, &KS_H, 0x01)

	sse.PXOR(&XTMP1, &XTMP2)
	sse.PXOR(&XTMP3, &XTMP4)
	sse.PXOR(&XTMP1, &XTMP3)
	sse.PSHUFB(&XTMP1, &shuf_mask_0_0_0_dw1)

	sse.POR(&XDIGEST, &XTMP1)

	return XDIGEST.Uint64s()[0], XDIGEST.Uint64s()[1]
}
