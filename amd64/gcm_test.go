package amd64

import (
	"encoding/hex"
	"testing"

	"github.com/emmansun/simd/amd64/sse"
)

func TestGCM(t *testing.T) {
	data, _ := hex.DecodeString("d38f531e5330944d3811556fff7b1f9f")
	var (
		B0 = sse.XMM{}
		B1 = sse.XMM{}
		B2 = sse.XMM{}
		B3 = sse.XMM{}
		B4 = sse.XMM{}
		T0 = sse.XMM{}
		T1 = sse.XMM{}
		T2 = sse.XMM{}
	)
	POLY := sse.Set64(0xc200000000000000, 0x0000000000000001)
	sse.SetBytes(&B0, data)
	sse.PSHUFD(&T0, &B0, 0xff)
	sse.MOVOU(&T1, &B0)
	sse.PSRAW(&T0, 31)
	sse.PAND(&T0, &POLY)
	sse.PSRLW(&T1, 31)
	sse.PSLLDQ(&T1, 4)
	sse.PSLLW(&B0, 1)
	sse.PXOR(&B0, &T0)
	sse.PXOR(&B0, &T1)
	if hex.EncodeToString(B0.Bytes()) != "a71fa73ca660289b7022aadefef73efc" {
		t.Errorf("B0 = %v; want a71fa73ca660289b7022aadefef73efc", hex.EncodeToString(B0.Bytes()))
	}
	// Karatsuba pre-computations
	sse.PSHUFD(&B1, &B0, 78)
	sse.PXOR(&B1, &B0)
	if hex.EncodeToString(B1.Bytes()) != "d73d0de258971667d73d0de258971667" {
		t.Errorf("B1 = %v; want d73d0de258971667d73d0de258971667", hex.EncodeToString(B1.Bytes()))
	}
	sse.MOVOU(&B2, &B0)
	sse.MOVOU(&B3, &B1)

	sse.MOVOU(&T0, &B2)
	sse.MOVOU(&T1, &B2)
	sse.MOVOU(&T2, &B3)
	sse.PCLMULQDQ(&T0, &B0, 0x00)
	sse.PCLMULQDQ(&T1, &B0, 0x11)
	sse.PCLMULQDQ(&T2, &B1, 0x00)

	sse.PXOR(&T2, &T0)
	sse.PXOR(&T2, &T1)
	sse.MOVOU(&B4, &T2)
	sse.PSLLDQ(&B4, 8)
	sse.PSRLDQ(&T2, 8)
	sse.PXOR(&T0, &B4)
	sse.PXOR(&T1, &T2)

	// Fast reduction
	// 1st reduction
	sse.MOVOU(&B2, &POLY)
	sse.PCLMULQDQ(&B2, &T0, 0x01)
	sse.PSHUFD(&T0, &T0, 78)
	sse.PXOR(&T0, &B2)
	// 2nd reduction
	sse.MOVOU(&B2, &POLY)
	sse.PCLMULQDQ(&B2, &T0, 0x01)
	sse.PSHUFD(&T0, &T0, 78)
	sse.PXOR(&B2, &T0)
	sse.PXOR(&B2, &T1)
	if hex.EncodeToString(B2.Bytes()) != "9208acefd693f27fc7223dce2c483080" {
		t.Errorf("B2 = %v; want 9208acefd693f27fc7223dce2c483080", hex.EncodeToString(B2.Bytes()))
	}

	sse.PSHUFD(&B3, &B2, 78)
	sse.PXOR(&B3, &B2)
	if hex.EncodeToString(B3.Bytes()) != "552a9121fadbc2ff552a9121fadbc2ff" {
		t.Errorf("B3 = %v; want 552a9121fadbc2ff552a9121fadbc2ff", hex.EncodeToString(B3.Bytes()))
	}
}
