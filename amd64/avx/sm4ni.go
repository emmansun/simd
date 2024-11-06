// https://cdrdv2-public.intel.com/782879/architecture-instruction-set-extensions-programming-reference.pdf
package avx

import (
	"encoding/binary"

	"github.com/emmansun/simd/alg/sm4"
	"github.com/emmansun/simd/amd64/sse"
)

// Perform Four Rounds of SM4 Key Expansion
// The VSM4KEY4 instruction performs four rounds of SM4 key expansion.
// The instruction operates on independent 128-bit lanes.
func VSM4KEY4(dst, src1, src2 *sse.XMM) {
	keyBytes := src2.Bytes()
	roundresult := &sse.XMM{}
	roundresultBytes := roundresult.Bytes()
	VMOVDQU(roundresult, src1)
	for i := 0; i < 16; i += 4 {
		ck := binary.LittleEndian.Uint32(keyBytes[i:])
		b0 := binary.LittleEndian.Uint32(roundresultBytes[:])
		b1 := binary.LittleEndian.Uint32(roundresultBytes[4:])
		b2 := binary.LittleEndian.Uint32(roundresultBytes[8:])
		b3 := binary.LittleEndian.Uint32(roundresultBytes[12:])
		intval := b0 ^ sm4.T_KEY(b1^b2^b3^ck)
		binary.LittleEndian.PutUint32(roundresultBytes[:], b1)
		binary.LittleEndian.PutUint32(roundresultBytes[4:], b2)
		binary.LittleEndian.PutUint32(roundresultBytes[8:], b3)
		binary.LittleEndian.PutUint32(roundresultBytes[12:], intval)
	}
	VMOVDQU(dst, roundresult)
}

// Performs Four Rounds of SM4 Encryption
// The SM4RNDS4 instruction performs four rounds of SM4 encryption.
// The instruction operates on independent 128-bit lanes.
func VSM4RNDS4(dst, src1, src2 *sse.XMM) {
	ckBytes := src2.Bytes()
	roundresult := &sse.XMM{}
	roundresultBytes := roundresult.Bytes()
	VMOVDQU(roundresult, src1)
	for i := 0; i < 16; i += 4 {
		rk := binary.LittleEndian.Uint32(ckBytes[i:])
		b0 := binary.LittleEndian.Uint32(roundresultBytes[:])
		b1 := binary.LittleEndian.Uint32(roundresultBytes[4:])
		b2 := binary.LittleEndian.Uint32(roundresultBytes[8:])
		b3 := binary.LittleEndian.Uint32(roundresultBytes[12:])
		intval := b0 ^ sm4.T(b1^b2^b3^rk)
		binary.LittleEndian.PutUint32(roundresultBytes[:], b1)
		binary.LittleEndian.PutUint32(roundresultBytes[4:], b2)
		binary.LittleEndian.PutUint32(roundresultBytes[8:], b3)
		binary.LittleEndian.PutUint32(roundresultBytes[12:], intval)
	}
	VMOVDQU(dst, roundresult)
}

func ExpandKey(out []uint32, key []byte) {
	_ = out[31]
	var (
		ck        = &sse.XMM{}
		fk        = &sse.XMM{}
		flip_mask = &sse.XMM{}
		keyXMM    = &sse.XMM{}
	)
	VMOVDQU_L16B(flip_mask, []byte{0x03, 0x02, 0x01, 0x00, 0x07, 0x06, 0x05, 0x04, 0x0b, 0x0a, 0x09, 0x08, 0x0f, 0x0e, 0x0d, 0x0c})
	VMOVDQU_L4S(fk, sm4.FK[:])

	VMOVDQU_L16B(keyXMM, key)
	VPSHUFB(keyXMM, keyXMM, flip_mask)
	VPXOR(keyXMM, keyXMM, fk)

	for i := 0; i < 32; i += 4 {
		VMOVDQU_L4S(ck, sm4.CK[i:])
		VSM4KEY4(keyXMM, keyXMM, ck)
		VMOVDQU_S4S(out[i:], keyXMM)
	}
}

func Encrypt(dst []byte, src []byte, key *[32]uint32) {
	_ = dst[15]
	_ = src[15]
	data := &sse.XMM{}
	rk := &sse.XMM{}
	flip_mask := &sse.XMM{}
	bswap_mask := &sse.XMM{}
	VMOVDQU_L16B(flip_mask, []byte{0x03, 0x02, 0x01, 0x00, 0x07, 0x06, 0x05, 0x04, 0x0b, 0x0a, 0x09, 0x08, 0x0f, 0x0e, 0x0d, 0x0c})
	VMOVDQU_L16B(bswap_mask, []byte{0x0f, 0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00})
	VMOVDQU_L16B(data, src)
	VPSHUFB(data, data, flip_mask)
	for i := 0; i < 32; i += 4 {
		VMOVDQU_L4S(rk, key[i:])
		VSM4RNDS4(data, data, rk)
	}
	VPSHUFB(data, data, bswap_mask)
	VMOVEDQU_S16B(dst, data)
}
