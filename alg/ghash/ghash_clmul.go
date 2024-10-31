// The theoretical proof behind it is still not very clear.
package ghash

import (
	"encoding/binary"

	"github.com/emmansun/simd/amd64/sse"
	"github.com/emmansun/simd/arm64"
	"github.com/emmansun/simd/ppc64"
)

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
		sse.MOVOU(&T0, &B2)
		sse.MOVOU(&T1, &B2)
		sse.MOVOU(&T2, &B3)
		sse.PCLMULQDQ(&T0, &B0, 0x00)
		sse.PCLMULQDQ(&T1, &B0, 0x11)
		sse.PCLMULQDQ(&T2, &B1, 0x00)

		g.processClmulResult(&T0, &T2, &T1, &B4)
		g.reduceRound(&T0, &B4, &POLY)
		g.reduceRound(&T0, &B4, &POLY)
		sse.PXOR(&T0, &T1)
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
		g.reduceRound(&ACC0, &T0, &POLY)
		g.reduceRound(&ACC0, &T0, &POLY)
		sse.PXOR(&ACC0, &ACC1) // ACC0 holds the result

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
	// 2H
	sse.MOVOU(ACC0, TH)
	sse.MOVOU(ACCM, THM)
	sse.MOVOU(ACC1, TH)

	// X * 2H
	sse.PSHUFD(T0, X0, 78)
	sse.PXOR(T0, X0)
	sse.PCLMULQDQ(ACC0, X0, 0x00)
	sse.PCLMULQDQ(ACC1, X0, 0x11)
	sse.PCLMULQDQ(ACCM, T0, 0x00)

	g.processClmulResult(ACC0, ACCM, ACC1, T0)
	g.reduceRound(ACC0, T0, POLY)
	g.reduceRound(ACC0, T0, POLY)
	sse.PXOR(ACC0, ACC1) // ACC0 holds the result
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

// Reduce X by POLY and store the result in X
func (g *clmulAMD64Ghash) reduceRound(X, T0, POLY *sse.XMM) {
	sse.MOVOU(T0, POLY)
	sse.PCLMULQDQ(T0, X, 0x01) // High POLY * Low X
	sse.PSHUFD(X, X, 78)       // High X -> Low X, Low X -> High X
	sse.PXOR(X, T0)
}

type clmulARM64Ghash struct {
	bytesProductTable [256]byte
}

func NewClmulARM64Ghash(h []byte) *clmulARM64Ghash {
	g := &clmulARM64Ghash{}
	var (
		B0   = &arm64.Vector128{}
		B1   = &arm64.Vector128{}
		B2   = &arm64.Vector128{}
		B3   = &arm64.Vector128{}
		T0   = &arm64.Vector128{}
		T1   = &arm64.Vector128{}
		T2   = &arm64.Vector128{}
		T3   = &arm64.Vector128{}
		POLY = &arm64.Vector128{}
		ZERO = &arm64.Vector128{}
	)
	arm64.VLD1_16B(h, B0)
	arm64.VREV64_B(B0, B0) // B0.D[0] = High part, B0.D[1] = Low part
	arm64.VEOR(ZERO, ZERO, ZERO)
	arm64.VLD1_2D([]uint64{0xc200000000000000, 0x0000000000000001}, POLY)

	// Multiply by 2 modulo P
	I := int64(B0.Uint64s()[0])
	I = I >> 63
	arm64.VLD1_2D([]uint64{uint64(I), uint64(I)}, T1)
	arm64.VAND(POLY, T1, T1)
	arm64.VUSHR_D(63, B0, T2)
	arm64.VEXT(8, ZERO, T2, T2)
	arm64.VSLI_D(1, B0, T2)
	arm64.VEOR(T1, T2, B0)

	// Karatsuba pre-computations
	arm64.VEXT(8, B0, B0, B1) // B1.D[0] = B0.D[1], B1.D[1] = B0.D[0]
	arm64.VEOR(B1, B0, B1)    // B1.D[0] = B0.D[1] ^ B0.D[0], B1.D[1] = B0.D[0] ^ B0.D[1]
	arm64.VST1_16B(B0, g.bytesProductTable[14*16:])
	arm64.VST1_16B(B1, g.bytesProductTable[15*16:])

	arm64.VMOV(B0, B2)
	arm64.VMOV(B1, B3)

	for i := 6; i >= 0; i-- {
		arm64.VPMULL(B0, B2, T1)  // T1 = ACC1 = B0.D[0] * B2.D[0]
		arm64.VPMULL2(B0, B2, T0) // T0 = ACC0 = B0.D[1] * B2.D[1]
		arm64.VPMULL(B1, B3, T2)  // T2 = ACCM = B1.D[0] * B3.D[0]
		g.processClmulResult(T0, T2, T1, ZERO, T3)
		g.reduceRound(T0, T2, POLY)
		g.reduceRound(T0, T2, POLY)
		arm64.VEOR(T0, T1, B2)
		arm64.VMOV(B2, B3)
		arm64.VEXT(8, B2, B2, B2) // B2.D[0] = B2.D[1], B2.D[1] = B2.D[0]
		arm64.VEOR(B2, B3, B3)    // B3.D[0] = B2.D[1] ^ B2.D[0], B3.D[1] = B2.D[0] ^ B2.D[1]
		arm64.VST1_16B(B2, g.bytesProductTable[(2*i)*16:])
		arm64.VST1_16B(B3, g.bytesProductTable[(2*i+1)*16:])
	}
	return g
}

func (g *clmulARM64Ghash) Hash(T *[16]byte, data []byte) {
	var (
		ACC0 = &arm64.Vector128{}
		ACC1 = &arm64.Vector128{}
		ACCM = &arm64.Vector128{}
		T0   = &arm64.Vector128{}
		T1   = &arm64.Vector128{}
		T2   = &arm64.Vector128{}
		T3   = &arm64.Vector128{}
		B0   = &arm64.Vector128{}
		B1   = &arm64.Vector128{}
		B2   = &arm64.Vector128{}
		B3   = &arm64.Vector128{}
		B4   = &arm64.Vector128{}
		B5   = &arm64.Vector128{}
		B6   = &arm64.Vector128{}
		B7   = &arm64.Vector128{}
		POLY = &arm64.Vector128{}
		ZERO = &arm64.Vector128{}
	)
	arm64.VEOR(ACC0, ACC0, ACC0)
	arm64.VEOR(ZERO, ZERO, ZERO)
	arm64.VLD1_2D([]uint64{0xc200000000000000, 0x0000000000000001}, POLY)
	// handle 8 blocks at a time
	for len(data) >= 128 {
		// load precomputed values
		arm64.VLD1_16B(g.bytesProductTable[16*0:], T1)
		arm64.VLD1_16B(g.bytesProductTable[16*1:], T2)
		// load 8 blocks
		arm64.VLD1_16B(data, B0)
		arm64.VLD1_16B(data[16:], B1)
		arm64.VLD1_16B(data[32:], B2)
		arm64.VLD1_16B(data[48:], B3)
		arm64.VLD1_16B(data[64:], B4)
		arm64.VLD1_16B(data[80:], B5)
		arm64.VLD1_16B(data[96:], B6)
		arm64.VLD1_16B(data[112:], B7)

		// process first block
		// prepare data for multiplication
		arm64.VREV64_B(B0, B0)
		arm64.VEOR(B0, ACC0, B0)
		arm64.VEXT(8, B0, B0, T0)
		arm64.VEOR(B0, T0, T0)
		// multiply
		arm64.VPMULL(B0, T1, ACC1)
		arm64.VPMULL2(B0, T1, ACC0)
		arm64.VPMULL(T0, T2, ACCM)

		g.mulRoundAAD(B1, T0, T1, T2, T3, ACC0, ACC1, ACCM, 1)
		g.mulRoundAAD(B2, T0, T1, T2, T3, ACC0, ACC1, ACCM, 2)
		g.mulRoundAAD(B3, T0, T1, T2, T3, ACC0, ACC1, ACCM, 3)
		g.mulRoundAAD(B4, T0, T1, T2, T3, ACC0, ACC1, ACCM, 4)
		g.mulRoundAAD(B5, T0, T1, T2, T3, ACC0, ACC1, ACCM, 5)
		g.mulRoundAAD(B6, T0, T1, T2, T3, ACC0, ACC1, ACCM, 6)
		g.mulRoundAAD(B7, T0, T1, T2, T3, ACC0, ACC1, ACCM, 7)

		g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
		g.reduceRound(ACC0, T0, POLY)
		g.reduceRound(ACC0, T0, POLY)
		arm64.VEOR(ACC0, ACC1, ACC0)    // ACC0 holds the result
		arm64.VEXT(8, ACC0, ACC0, ACC0) // ACC0.D[0] = ACC0.D[1], ACC0.D[1] = ACC0.D[0]

		data = data[128:]
	}

	// load precomputed values
	arm64.VLD1_16B(g.bytesProductTable[16*14:], T1)
	arm64.VLD1_16B(g.bytesProductTable[16*15:], T2)

	// handle one block at a time
	for len(data) >= 16 {
		// load 1 block
		arm64.VLD1_16B(data, B0)
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, T0, T1, T2, POLY, ZERO)

		data = data[16:]
	}
	if len(data) > 0 {
		var partialBlock [16]byte
		copy(partialBlock[:], data)
		arm64.VLD1_16B(partialBlock[:], B0)
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, T0, T1, T2, POLY, ZERO)
	}
	arm64.VREV64_B(ACC0, ACC0)
	arm64.VST1_16B(ACC0, T[:])
}

func (g *clmulARM64Ghash) mulOneBlock(B0, ACC0, ACCM, ACC1, T0, T1, T2, POLY, ZERO *arm64.Vector128) {
	arm64.VREV64_B(B0, B0)
	arm64.VEOR(B0, ACC0, B0)
	arm64.VEXT(8, B0, B0, T0)
	arm64.VEOR(B0, T0, T0)
	// multiply
	arm64.VPMULL(B0, T1, ACC1)
	arm64.VPMULL2(B0, T1, ACC0)
	arm64.VPMULL(T0, T2, ACCM)

	g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
	g.reduceRound(ACC0, T0, POLY)
	g.reduceRound(ACC0, T0, POLY)
	arm64.VEOR(ACC0, ACC1, ACC0)    // ACC0 holds the result
	arm64.VEXT(8, ACC0, ACC0, ACC0) // ACC0.D[0] = ACC0.D[1], ACC0.D[1] = ACC0.D[0]

}

// Multiply X by 2H^(8-i) and accumulate the result in ACC0, ACC1, ACCM
// Y = [ACC0, ACCM, ACC1]
// Y = X * (2H)^(8-i) + Y
// T0, T1, T2, T3 are temporary registers
func (g *clmulARM64Ghash) mulRoundAAD(X, T0, T1, T2, T3, ACC0, ACC1, ACCM *arm64.Vector128, i int) {
	// prepare data for multiplication
	arm64.VREV64_B(X, X)
	arm64.VEXT(8, X, X, T0)
	arm64.VEOR(X, T0, T0)
	// load precomputed values
	arm64.VLD1_16B(g.bytesProductTable[16*(i*2):], T1)
	arm64.VLD1_16B(g.bytesProductTable[16*((i*2)+1):], T2)
	// multiply and accumulate
	arm64.VPMULL(X, T1, T3)    // T3 = new ACC1
	arm64.VEOR(ACC1, T3, ACC1) // ACC1 = ACC1 ^ T3
	arm64.VPMULL2(X, T1, T3)   // T3 = new ACC0
	arm64.VEOR(ACC0, T3, ACC0) // ACC0 = ACC0 ^ T3
	arm64.VPMULL(T0, T2, T3)   // T3 = new ACCM
	arm64.VEOR(ACCM, T3, ACCM) // ACCM = ACCM ^ T3
}

func (g *clmulARM64Ghash) processClmulResult(ACC0, ACCM, ACC1, ZERO, T *arm64.Vector128) {
	arm64.VEOR(ACC0, ACCM, ACCM)
	arm64.VEOR(ACC1, ACCM, ACCM)
	arm64.VEXT(8, ZERO, ACCM, T)
	arm64.VEXT(8, ACCM, ZERO, ACCM)
	arm64.VEOR(ACCM, ACC0, ACC0) // ACC0
	arm64.VEOR(T, ACC1, ACC1)    // ACC1
}

// Reduce X by POLY and store the result in X
// T is a temporary register
func (g *clmulARM64Ghash) reduceRound(D, T, POLY *arm64.Vector128) {
	arm64.VPMULL(D, POLY, T)
	arm64.VEXT(8, D, D, D)
	arm64.VEOR(D, T, D)
}

type clmulPPC64Ghash struct {
	isPPC64LE         bool
	bytesProductTable [400]byte // (8 * 3 + 1) * 16
}

func NewClmulPPC64Ghash(h []byte, isPPC64LE bool) *clmulPPC64Ghash {
	g := &clmulPPC64Ghash{}
	g.isPPC64LE = isPPC64LE
	var (
		XC2  = &ppc64.Vector128{}
		T0   = &ppc64.Vector128{}
		T1   = &ppc64.Vector128{}
		T2   = &ppc64.Vector128{}
		ZERO = &ppc64.Vector128{}
		H    = &ppc64.Vector128{}
		HL   = &ppc64.Vector128{}
		HH   = &ppc64.Vector128{}
		IN   = &ppc64.Vector128{}
		XL   = &ppc64.Vector128{}
		XM   = &ppc64.Vector128{}
		XH   = &ppc64.Vector128{}
		H2   = &ppc64.Vector128{}
		H2L  = &ppc64.Vector128{}
		H2H  = &ppc64.Vector128{}
	)

	var h1, h2 uint64
	if isPPC64LE {
		// can use VPERM to convert from little-endian to big-endian
		var v [16]byte
		h1 = binary.LittleEndian.Uint64(h[:8])
		h2 = binary.LittleEndian.Uint64(h[8:])
		binary.BigEndian.PutUint64(v[:], h1)
		binary.BigEndian.PutUint64(v[8:], h2)
		ppc64.LXVD2X_PPC64LE(v[:], H)
	} else {
		ppc64.LXVD2X(h, H)
	}

	ppc64.VXOR(ZERO, ZERO, ZERO)
	g.initPoly(XC2, ZERO, T0, T1)

	// Multiply by 2 modulo P
	ppc64.VSPLTISB(7, T2)
	ppc64.VSPLTB(0, H, T1)  // most significant byte
	ppc64.VSL(H, T0, H)     // H<<=1
	ppc64.VSRAB(T1, T2, T1) // broadcast carry bit
	ppc64.VAND(T1, XC2, T1)
	ppc64.VXOR(H, T1, IN)        // twisted H
	ppc64.VSLDOI(8, IN, IN, H)   // twist even more ...
	ppc64.VSLDOI(8, ZERO, H, HL) // ... and split
	ppc64.VSLDOI(8, H, ZERO, HH)

	ppc64.VSLDOI(8, ZERO, XC2, XC2) // 0xc2.0
	if isPPC64LE {
		ppc64.STXVD2X_PPC64LE(XC2, g.bytesProductTable[24*16:])
		ppc64.STXVD2X_PPC64LE(HL, g.bytesProductTable[21*16:])
		ppc64.STXVD2X_PPC64LE(H, g.bytesProductTable[22*16:])
		ppc64.STXVD2X_PPC64LE(HH, g.bytesProductTable[23*16:])
	} else {
		ppc64.STXVD2X(XC2, g.bytesProductTable[24*16:])
		ppc64.STXVD2X(HL, g.bytesProductTable[21*16:])
		ppc64.STXVD2X(H, g.bytesProductTable[22*16:])
		ppc64.STXVD2X(HH, g.bytesProductTable[23*16:])
	}

	for i := 6; i >= 0; i-- {
		ppc64.VPMSUMD(IN, HL, XL) // H.lo·H.lo
		ppc64.VPMSUMD(IN, H, XM)  // H.hi·H.lo+H.lo·H.hi
		ppc64.VPMSUMD(IN, HH, XH) // H.hi·H.hi

		g.processClmulResult(XL, XM, XH, ZERO, T0)
		g.reduceRound(XL, T2, XC2)
		g.reduceRound(XL, T2, XC2)

		ppc64.VXOR(XL, XH, IN)

		ppc64.VSLDOI(8, IN, IN, H2)
		ppc64.VSLDOI(8, ZERO, H2, H2L)
		ppc64.VSLDOI(8, H2, ZERO, H2H)

		if isPPC64LE {
			ppc64.STXVD2X_PPC64LE(H2L, g.bytesProductTable[(3*i)*16:])
			ppc64.STXVD2X_PPC64LE(H2, g.bytesProductTable[(3*i+1)*16:])
			ppc64.STXVD2X_PPC64LE(H2H, g.bytesProductTable[(3*i+2)*16:])
		} else {
			ppc64.STXVD2X(H2L, g.bytesProductTable[(3*i)*16:])
			ppc64.STXVD2X(H2, g.bytesProductTable[(3*i+1)*16:])
			ppc64.STXVD2X(H2H, g.bytesProductTable[(3*i+2)*16:])
		}
	}

	return g
}

func (g *clmulPPC64Ghash) Hash(T *[16]byte, data []byte) {
	var (
		XC2   = &ppc64.Vector128{}
		VPERM = &ppc64.Vector128{}
		T0    = &ppc64.Vector128{}
		ZERO  = &ppc64.Vector128{}
		H     = &ppc64.Vector128{}
		HL    = &ppc64.Vector128{}
		HH    = &ppc64.Vector128{}
		ACC0  = &ppc64.Vector128{}
		ACC1  = &ppc64.Vector128{}
		ACCM  = &ppc64.Vector128{}
		B0    = &ppc64.Vector128{}
		B1    = &ppc64.Vector128{}
		B2    = &ppc64.Vector128{}
		B3    = &ppc64.Vector128{}
		B4    = &ppc64.Vector128{}
		B5    = &ppc64.Vector128{}
		B6    = &ppc64.Vector128{}
		B7    = &ppc64.Vector128{}
	)
	if g.isPPC64LE {
		ppc64.LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, VPERM)
	}
	ppc64.VXOR(ZERO, ZERO, ZERO)
	ppc64.VXOR(ACC0, ACC0, ACC0)
	g.loadPrecomputed(24, XC2)

	// handle 8 blocks at a time
	for len(data) >= 128 {
		// load 8 blocks
		g.load8Blocks(data, B0, B1, B2, B3, B4, B5, B6, B7, VPERM)

		// load precomputed values
		g.loadPrecomputed(0, HL)
		g.loadPrecomputed(1, H)
		g.loadPrecomputed(2, HH)

		// process first block
		// prepare data for multiplication
		ppc64.VXOR(B0, ACC0, B0)
		ppc64.VPMSUMD(B0, HL, ACC0)
		ppc64.VPMSUMD(B0, H, ACCM)
		ppc64.VPMSUMD(B0, HH, ACC1)

		g.mulRoundAAD(B1, ACC0, ACC1, ACCM, HL, H, HH, T0, 1)
		g.mulRoundAAD(B2, ACC0, ACC1, ACCM, HL, H, HH, T0, 2)
		g.mulRoundAAD(B3, ACC0, ACC1, ACCM, HL, H, HH, T0, 3)
		g.mulRoundAAD(B4, ACC0, ACC1, ACCM, HL, H, HH, T0, 4)
		g.mulRoundAAD(B5, ACC0, ACC1, ACCM, HL, H, HH, T0, 5)
		g.mulRoundAAD(B6, ACC0, ACC1, ACCM, HL, H, HH, T0, 6)
		g.mulRoundAAD(B7, ACC0, ACC1, ACCM, HL, H, HH, T0, 7)

		g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
		g.reduceRound(ACC0, T0, XC2)
		g.reduceRound(ACC0, T0, XC2)
		ppc64.VXOR(ACC0, ACC1, ACC0) // ACC0 holds the result

		data = data[128:]
	}

	// load precomputed values
	g.loadPrecomputed(21, HL)
	g.loadPrecomputed(22, H)
	g.loadPrecomputed(23, HH)

	// handle one block at a time
	for len(data) >= 16 {
		// load 1 block
		g.loadData(data, 0, B0)
		if g.isPPC64LE {
			ppc64.VPERM(B0, B0, VPERM, B0)
		}
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, HL, H, HH, T0, XC2, ZERO)
		data = data[16:]
	}
	if len(data) > 0 {
		var partialBlock [16]byte
		copy(partialBlock[:], data)
		g.loadData(partialBlock[:], 0, B0)
		if g.isPPC64LE {
			ppc64.VPERM(B0, B0, VPERM, B0)
		}
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, HL, H, HH, T0, XC2, ZERO)
	}
	if g.isPPC64LE {
		ppc64.VPERM(ACC0, ACC0, VPERM, ACC0)
		ppc64.STXVD2X_PPC64LE(ACC0, T[:])
	} else {
		ppc64.STXVD2X(ACC0, T[:])
	}
}

func (g *clmulPPC64Ghash) load8Blocks(data []byte, B0, B1, B2, B3, B4, B5, B6, B7, VPERM *ppc64.Vector128) {
	g.loadData(data, 0, B0)
	g.loadData(data, 1, B1)
	g.loadData(data, 2, B2)
	g.loadData(data, 3, B3)
	g.loadData(data, 4, B4)
	g.loadData(data, 5, B5)
	g.loadData(data, 6, B6)
	g.loadData(data, 7, B7)

	if g.isPPC64LE {
		ppc64.VPERM(B0, B0, VPERM, B0)
		ppc64.VPERM(B1, B1, VPERM, B1)
		ppc64.VPERM(B2, B2, VPERM, B2)
		ppc64.VPERM(B3, B3, VPERM, B3)
		ppc64.VPERM(B4, B4, VPERM, B4)
		ppc64.VPERM(B5, B5, VPERM, B5)
		ppc64.VPERM(B6, B6, VPERM, B6)
		ppc64.VPERM(B7, B7, VPERM, B7)
	}
}

func (g *clmulPPC64Ghash) mulOneBlock(B0, ACC0, ACCM, ACC1, HL, H, HH, T0, XC2, ZERO *ppc64.Vector128) {
	ppc64.VXOR(B0, ACC0, B0)
	ppc64.VPMSUMD(B0, HL, ACC0)
	ppc64.VPMSUMD(B0, H, ACCM)
	ppc64.VPMSUMD(B0, HH, ACC1)
	g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
	g.reduceRound(ACC0, T0, XC2)
	g.reduceRound(ACC0, T0, XC2)
	ppc64.VXOR(ACC0, ACC1, ACC0) // ACC0 holds the result
}

func (g *clmulPPC64Ghash) mulRoundAAD(IN, ACC0, ACC1, ACCM, T0, T1, T2, T3 *ppc64.Vector128, i int) {
	// load precomputed values
	g.loadPrecomputed(3*i, T0)   // H2L
	g.loadPrecomputed(3*i+1, T1) // H2
	g.loadPrecomputed(3*i+2, T2) // H2H

	ppc64.VPMSUMD(IN, T0, T3)  // H2L
	ppc64.VXOR(ACC0, T3, ACC0) // ACC0 = ACC0 ^ H2L
	ppc64.VPMSUMD(IN, T1, T3)
	ppc64.VXOR(ACCM, T3, ACCM) // ACCM = ACCM ^ H2
	ppc64.VPMSUMD(IN, T2, T3)
	ppc64.VXOR(ACC1, T3, ACC1) // ACC1 = ACC1 ^ H2H
}

func (g *clmulPPC64Ghash) loadData(data []byte, idx int, T *ppc64.Vector128) {
	if g.isPPC64LE {
		ppc64.LXVD2X_PPC64LE(data[idx*16:], T)
	} else {
		ppc64.LXVD2X(data[idx*16:], T)
	}
}

func (g *clmulPPC64Ghash) loadPrecomputed(i int, T *ppc64.Vector128) {
	if g.isPPC64LE {
		ppc64.LXVD2X_PPC64LE(g.bytesProductTable[i*16:], T)
	} else {
		ppc64.LXVD2X(g.bytesProductTable[i*16:], T)
	}
}

func (g *clmulPPC64Ghash) initPoly(POLY, ZERO, T0, T1 *ppc64.Vector128) {
	// compute POLY = 0xc2000000000000000000000000000001
	ppc64.VSPLTISB(0x10, POLY)         // 0xf0
	ppc64.VSPLTISB(1, T0)              // one
	ppc64.VADDUBM(POLY, POLY, POLY)    // 0xe0
	ppc64.VOR(POLY, T0, POLY)          // 0xe1
	ppc64.VSLDOI(15, POLY, ZERO, POLY) // 0xe1...
	ppc64.VSLDOI(1, ZERO, T0, T1)      // ...1
	ppc64.VADDUBM(POLY, POLY, POLY)    // 0xc2...
	ppc64.VOR(POLY, T1, POLY)          // 0xc2....01
}

func (g *clmulPPC64Ghash) processClmulResult(ACC0, ACCM, ACC1, ZERO, T *ppc64.Vector128) {
	ppc64.VSLDOI(8, ACCM, ZERO, T) // T = ACCM.hi || zero
	ppc64.VXOR(ACC0, T, ACC0)      // ACC0 = ACC0 ^ (ACCM.hi || zero) = ACCM.hi || ACC0.lo
	ppc64.VSLDOI(8, ZERO, ACCM, T) // T = zero || ACCM.lo
	ppc64.VXOR(ACC1, T, ACC1)      // ACC1 = ACC1 ^ (zero || ACCM.lo) = ACC1.hi || ACCM.lo
}

func (g *clmulPPC64Ghash) reduceRound(ACC0, T, POLY *ppc64.Vector128) {
	ppc64.VPMSUMD(ACC0, POLY, T) // POLY.hi = 0
	ppc64.VSLDOI(8, ACC0, ACC0, ACC0)
	ppc64.VXOR(ACC0, T, ACC0)
}
