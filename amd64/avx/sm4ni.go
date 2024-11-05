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
