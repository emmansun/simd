// The theoretical proof behind it is still not very clear.
package amd64

import "github.com/emmansun/simd/amd64/sse"

type clmulAMD64Ghash struct {
	bytesProductTable [256]byte
}

func NewClmulAMD64Ghash(h []byte) *clmulAMD64Ghash {
	g := &clmulAMD64Ghash{}
	// H * 2
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
	BSWAP := sse.Set64(0x0001020304050607, 0x08090a0b0c0d0e0f)
	sse.SetBytes(&B0, h)
	sse.PSHUFB(&B0, &BSWAP)
	sse.PSHUFD(&T0, &B0, 0xff)
	sse.MOVOU(&T1, &B0)
	sse.PSRAW(&T0, 31)
	sse.PAND(&T0, &POLY)
	sse.PSRLW(&T1, 31)
	sse.PSLLDQ(&T1, 4)
	sse.PSLLW(&B0, 1)
	sse.PXOR(&B0, &T0)
	sse.PXOR(&B0, &T1)

	// Karatsuba pre-computations
	copy(g.bytesProductTable[14*16:], B0.Bytes())
	sse.PSHUFD(&B1, &B0, 78)
	sse.PXOR(&B1, &B0)
	copy(g.bytesProductTable[15*16:], B1.Bytes())
	sse.MOVOU(&B2, &B0)
	sse.MOVOU(&B3, &B1)

	for i := 6; i >= 0; i-- {
		// Karatsuba multiplication
		sse.MOVOU(&T0, &B2)
		sse.MOVOU(&T1, &B2)
		sse.MOVOU(&T2, &B3)
		sse.PCLMULQDQ(&T0, &B0, 0x00)
		sse.PCLMULQDQ(&T1, &B0, 0x11)
		sse.PCLMULQDQ(&T2, &B1, 0x00)
		g.processClmulResult(&T0, &T2, &T1, &B4)
		g.fastReduction(&T1, &T0, &B4, &POLY)
		sse.MOVOU(&B2, &T0)

		copy(g.bytesProductTable[(i*2)*16:], B2.Bytes())
		sse.PSHUFD(&B3, &B2, 78)
		sse.PXOR(&B3, &B2)
		copy(g.bytesProductTable[((i*2)+1)*16:], B3.Bytes())
	}
	return g
}

func (g *clmulAMD64Ghash) Hash(T *[16]byte, data []byte) {
	var (
		ACC0 = sse.XMM{}
		ACC1 = sse.XMM{}
		ACCM = sse.XMM{}
		T0   = sse.XMM{}
		T1   = sse.XMM{}
		T2   = sse.XMM{}
		X0   = sse.XMM{}
		X1   = sse.XMM{}
		X2   = sse.XMM{}
		X3   = sse.XMM{}
		X4   = sse.XMM{}
		X5   = sse.XMM{}
		X6   = sse.XMM{}
		X7   = sse.XMM{}
	)
	POLY := sse.Set64(0xc200000000000000, 0x0000000000000001)
	BSWAP := sse.Set64(0x0001020304050607, 0x08090a0b0c0d0e0f)
	sse.PXOR(&ACC0, &ACC0)

	// handle 8 blocks at a time
	for len(data) >= 128 {
		// load 8 blocks
		sse.SetBytes(&X0, data)
		sse.SetBytes(&X1, data[16:])
		sse.SetBytes(&X2, data[32:])
		sse.SetBytes(&X3, data[48:])
		sse.SetBytes(&X4, data[64:])
		sse.SetBytes(&X5, data[80:])
		sse.SetBytes(&X6, data[96:])
		sse.SetBytes(&X7, data[112:])
		sse.PSHUFB(&X0, &BSWAP)
		sse.PSHUFB(&X1, &BSWAP)
		sse.PSHUFB(&X2, &BSWAP)
		sse.PSHUFB(&X3, &BSWAP)
		sse.PSHUFB(&X4, &BSWAP)
		sse.PSHUFB(&X5, &BSWAP)
		sse.PSHUFB(&X6, &BSWAP)
		sse.PSHUFB(&X7, &BSWAP)
		// add previous result
		sse.PXOR(&X0, &ACC0)

		sse.SetBytes(&ACC0, g.bytesProductTable[16*0:])
		sse.SetBytes(&ACCM, g.bytesProductTable[16*1:])
		sse.MOVOU(&ACC1, &ACC0)

		// Karatsuba multiplication
		sse.PSHUFD(&T1, &X0, 78)
		sse.PXOR(&T1, &X0)
		sse.PCLMULQDQ(&ACC0, &X0, 0x00)
		sse.PCLMULQDQ(&ACC1, &X0, 0x11)
		sse.PCLMULQDQ(&ACCM, &T1, 0x00)

		g.mulRoundAAD(&X1, &T1, &T2, &ACC0, &ACC1, &ACCM, 1)
		g.mulRoundAAD(&X2, &T1, &T2, &ACC0, &ACC1, &ACCM, 2)
		g.mulRoundAAD(&X3, &T1, &T2, &ACC0, &ACC1, &ACCM, 3)
		g.mulRoundAAD(&X4, &T1, &T2, &ACC0, &ACC1, &ACCM, 4)
		g.mulRoundAAD(&X5, &T1, &T2, &ACC0, &ACC1, &ACCM, 5)
		g.mulRoundAAD(&X6, &T1, &T2, &ACC0, &ACC1, &ACCM, 6)
		g.mulRoundAAD(&X7, &T1, &T2, &ACC0, &ACC1, &ACCM, 7)

		g.processClmulResult(&ACC0, &ACCM, &ACC1, &T0)
		// postponed reduction
		g.fastReduction(&ACC1, &ACC0, &T0, &POLY)
		data = data[128:]
	}
	sse.SetBytes(&T1, g.bytesProductTable[16*14:])
	sse.SetBytes(&T2, g.bytesProductTable[16*15:])
	// handle one block at a time
	for len(data) >= 16 {
		// load 1 block
		sse.SetBytes(&X0, data)
		sse.PSHUFB(&X0, &BSWAP)
		g.mulOneBlock(&X0, &ACC0, &ACCM, &ACC1, &T0, &T1, &T2, &POLY)
		data = data[16:]
	}
	if len(data) > 0 {
		var partialBlock [16]byte
		copy(partialBlock[:], data)
		sse.SetBytes(&X0, partialBlock[:])
		sse.PSHUFB(&X0, &BSWAP)
		g.mulOneBlock(&X0, &ACC0, &ACCM, &ACC1, &T0, &T1, &T2, &POLY)
	}
	sse.PSHUFB(&ACC0, &BSWAP)
	copy(T[:], ACC0.Bytes())
}

// Multiply X by 2H and accumulate the result in ACC0
// ACC0 = X * 2H + ACC0
// X0 is the input block
// ACCM, ACC1, T0 are temporary registers
// TH, THM are 2H for Karatsuba representation respectively
func (g *clmulAMD64Ghash) mulOneBlock(X0, ACC0, ACCM, ACC1, T0, TH, THM, POLY *sse.XMM) {
	// add previous result
	sse.PXOR(X0, ACC0)
	sse.MOVOU(ACC0, TH)
	sse.MOVOU(ACCM, THM)
	sse.MOVOU(ACC1, TH)

	// Karatsuba multiplication
	sse.PSHUFD(T0, X0, 78)
	sse.PXOR(T0, X0)
	sse.PCLMULQDQ(ACC0, X0, 0x00)
	sse.PCLMULQDQ(ACC1, X0, 0x11)
	sse.PCLMULQDQ(ACCM, T0, 0x00)

	g.processClmulResult(ACC0, ACCM, ACC1, T0)
	g.fastReduction(ACC1, ACC0, T0, POLY)
}

// handle Karatsuba results in ACC0, ACCM, ACC1 to get the final result in ACC0, ACC1
// T is a temporary register
func (g *clmulAMD64Ghash) processClmulResult(ACC0, ACCM, ACC1, T *sse.XMM) {
	sse.PXOR(ACCM, ACC0)
	sse.PXOR(ACCM, ACC1)
	sse.MOVOU(T, ACCM)
	sse.PSRLDQ(ACCM, 8)
	sse.PSLLDQ(T, 8)
	sse.PXOR(ACC1, ACCM)
	sse.PXOR(ACC0, T)
}

// Multiply X by 2H^(8-i) and accumulate the result in ACC0, ACC1, ACCM
// Y = [ACC0, ACCM, ACC1]
// Y = X * (2H)^(8-i) + Y
// T1, T2 are temporary registers
func (g *clmulAMD64Ghash) mulRoundAAD(X, T1, T2, ACC0, ACC1, ACCM *sse.XMM, i int) {
	sse.SetBytes(T1, g.bytesProductTable[16*(i*2):])
	sse.MOVOU(T2, T1)
	sse.PCLMULQDQ(T1, X, 0x00)
	sse.PCLMULQDQ(T2, X, 0x11)

	sse.PXOR(ACC0, T1)
	sse.PXOR(ACC1, T2)

	sse.PSHUFD(T1, X, 78)
	sse.PXOR(X, T1)
	sse.SetBytes(T1, g.bytesProductTable[16*((i*2)+1):])
	sse.PCLMULQDQ(T1, X, 0x00)
	sse.PXOR(ACCM, T1)
}

// Fast reduction for [X1:X0] by POLY and store the result in X0
// T is a temporary register
func (g *clmulAMD64Ghash) fastReduction(X1, X0, T, POLY *sse.XMM) {
	sse.MOVOU(T, POLY)
	sse.PCLMULQDQ(T, X0, 0x01) // High POLY * Low X
	sse.PSHUFD(X0, X0, 78)     // High X -> Low X, Low X -> High X
	sse.PXOR(X0, T)
	sse.MOVOU(T, POLY)
	sse.PCLMULQDQ(T, X0, 0x01) // High POLY * Low X
	sse.PSHUFD(X0, X0, 78)     // High X -> Low X, Low X -> High X
	sse.PXOR(X0, T)
	sse.PXOR(X0, X1)
}
