// The theoretical proof behind it is still not very clear.
package arm64

type clmulARM64Ghash struct {
	bytesProductTable [256]byte
}

func NewClmulARM64Ghash(h []byte) *clmulARM64Ghash {
	g := &clmulARM64Ghash{}
	var (
		B0   = &Vector128{}
		B1   = &Vector128{}
		B2   = &Vector128{}
		B3   = &Vector128{}
		T0   = &Vector128{}
		T1   = &Vector128{}
		T2   = &Vector128{}
		T3   = &Vector128{}
		POLY = &Vector128{}
		ZERO = &Vector128{}
	)
	VLD1_16B(h, B0)
	VREV64_B(B0, B0) // B0.D[0] = High part, B0.D[1] = Low part
	VEOR(ZERO, ZERO, ZERO)
	VLD1_2D([]uint64{0xc200000000000000, 0x0000000000000001}, POLY)

	// Multiply by 2 modulo P
	I := int64(B0.Uint64s()[0])
	I = I >> 63
	VLD1_2D([]uint64{uint64(I), uint64(I)}, T1)
	VAND(POLY, T1, T1)
	VUSHR_D(63, B0, T2)
	VEXT(8, ZERO, T2, T2)
	VSLI_D(1, B0, T2)
	VEOR(T1, T2, B0)

	VEXT(8, B0, B0, B1) // B1.D[0] = B0.D[1], B1.D[1] = B0.D[0]
	VEOR(B1, B0, B1)    // B1.D[0] = B0.D[1] ^ B0.D[0], B1.D[1] = B0.D[0] ^ B0.D[1]
	VST1_16B(B0, g.bytesProductTable[14*16:])
	VST1_16B(B1, g.bytesProductTable[15*16:])

	VMOV(B0, B2)
	VMOV(B1, B3)

	for i := 6; i >= 0; i-- {
		// Karatsuba multiplication
		VPMULL(B0, B2, T1)  // T1 = ACC1 = B0.D[0] * B2.D[0]
		VPMULL2(B0, B2, T0) // T0 = ACC0 = B0.D[1] * B2.D[1]
		VPMULL(B1, B3, T2)  // T2 = ACCM = B1.D[0] * B3.D[0]
		g.processClmulResult(T0, T2, T1, ZERO, T3)
		g.fastReduction(B2, T1, T0, T3, POLY)
		VMOV(B2, B3)
		VEXT(8, B2, B2, B2) // B2.D[0] = B2.D[1], B2.D[1] = B2.D[0]
		VEOR(B2, B3, B3)    // B3.D[0] = B2.D[1] ^ B2.D[0], B3.D[1] = B2.D[0] ^ B2.D[1]
		VST1_16B(B2, g.bytesProductTable[(2*i)*16:])
		VST1_16B(B3, g.bytesProductTable[(2*i+1)*16:])
	}
	return g
}

func (g *clmulARM64Ghash) Hash(T *[16]byte, data []byte) {
	var (
		ACC0 = &Vector128{}
		ACC1 = &Vector128{}
		ACCM = &Vector128{}
		T0   = &Vector128{}
		T1   = &Vector128{}
		T2   = &Vector128{}
		T3   = &Vector128{}
		B0   = &Vector128{}
		B1   = &Vector128{}
		B2   = &Vector128{}
		B3   = &Vector128{}
		B4   = &Vector128{}
		B5   = &Vector128{}
		B6   = &Vector128{}
		B7   = &Vector128{}
		POLY = &Vector128{}
		ZERO = &Vector128{}
	)
	VEOR(ACC0, ACC0, ACC0)
	VEOR(ZERO, ZERO, ZERO)
	VLD1_2D([]uint64{0xc200000000000000, 0x0000000000000001}, POLY)
	// handle 8 blocks at a time
	for len(data) >= 128 {
		// load precomputed values
		VLD1_16B(g.bytesProductTable[16*0:], T1)
		VLD1_16B(g.bytesProductTable[16*1:], T2)
		// load 8 blocks
		VLD1_16B(data, B0)
		VLD1_16B(data[16:], B1)
		VLD1_16B(data[32:], B2)
		VLD1_16B(data[48:], B3)
		VLD1_16B(data[64:], B4)
		VLD1_16B(data[80:], B5)
		VLD1_16B(data[96:], B6)
		VLD1_16B(data[112:], B7)

		// process first block
		// prepare data for multiplication
		VREV64_B(B0, B0)
		VEOR(B0, ACC0, B0)
		VEXT(8, B0, B0, T0)
		VEOR(B0, T0, T0)
		// Karatsuba multiplication
		VPMULL(B0, T1, ACC1)
		VPMULL2(B0, T1, ACC0)
		VPMULL(T0, T2, ACCM)

		g.mulRoundAAD(B1, T0, T1, T2, T3, ACC0, ACC1, ACCM, 1)
		g.mulRoundAAD(B2, T0, T1, T2, T3, ACC0, ACC1, ACCM, 2)
		g.mulRoundAAD(B3, T0, T1, T2, T3, ACC0, ACC1, ACCM, 3)
		g.mulRoundAAD(B4, T0, T1, T2, T3, ACC0, ACC1, ACCM, 4)
		g.mulRoundAAD(B5, T0, T1, T2, T3, ACC0, ACC1, ACCM, 5)
		g.mulRoundAAD(B6, T0, T1, T2, T3, ACC0, ACC1, ACCM, 6)
		g.mulRoundAAD(B7, T0, T1, T2, T3, ACC0, ACC1, ACCM, 7)

		// delayed reduction
		g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
		g.fastReduction(ACC0, ACC1, ACC0, T0, POLY)
		VEXT(8, ACC0, ACC0, ACC0) // ACC0.D[0] = ACC0.D[1], ACC0.D[1] = ACC0.D[0]

		data = data[128:]
	}

	// load precomputed values
	VLD1_16B(g.bytesProductTable[16*14:], T1)
	VLD1_16B(g.bytesProductTable[16*15:], T2)

	// handle one block at a time
	for len(data) >= 16 {
		// load 1 block
		VLD1_16B(data, B0)
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, T0, T1, T2, POLY, ZERO)

		data = data[16:]
	}
	if len(data) > 0 {
		var partialBlock [16]byte
		copy(partialBlock[:], data)
		VLD1_16B(partialBlock[:], B0)
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, T0, T1, T2, POLY, ZERO)
	}
	VREV64_B(ACC0, ACC0)
	VST1_16B(ACC0, T[:])
}

func (g *clmulARM64Ghash) mulOneBlock(B0, ACC0, ACCM, ACC1, T0, T1, T2, POLY, ZERO *Vector128) {
	VREV64_B(B0, B0)
	VEOR(B0, ACC0, B0)
	VEXT(8, B0, B0, T0)
	VEOR(B0, T0, T0)
	// Karatsuba multiplication
	VPMULL(B0, T1, ACC1)
	VPMULL2(B0, T1, ACC0)
	VPMULL(T0, T2, ACCM)

	g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
	g.fastReduction(ACC0, ACC1, ACC0, T0, POLY)
	VEXT(8, ACC0, ACC0, ACC0) // ACC0.D[0] = ACC0.D[1], ACC0.D[1] = ACC0.D[0]

}

// Multiply X by 2H^(8-i) and accumulate the result in ACC0, ACC1, ACCM
// Y = [ACC0, ACCM, ACC1]
// Y = X * (2H)^(8-i) + Y
// T0, T1, T2, T3 are temporary registers
func (g *clmulARM64Ghash) mulRoundAAD(X, T0, T1, T2, T3, ACC0, ACC1, ACCM *Vector128, i int) {
	// prepare data for multiplication
	VREV64_B(X, X)
	VEXT(8, X, X, T0)
	VEOR(X, T0, T0)
	// load precomputed values
	VLD1_16B(g.bytesProductTable[16*(i*2):], T1)
	VLD1_16B(g.bytesProductTable[16*((i*2)+1):], T2)
	// Karatsuba multiplication and accumulate
	VPMULL(X, T1, T3)    // T3 = new ACC1
	VEOR(ACC1, T3, ACC1) // ACC1 = ACC1 ^ T3
	VPMULL2(X, T1, T3)   // T3 = new ACC0
	VEOR(ACC0, T3, ACC0) // ACC0 = ACC0 ^ T3
	VPMULL(T0, T2, T3)   // T3 = new ACCM
	VEOR(ACCM, T3, ACCM) // ACCM = ACCM ^ T3
}

func (g *clmulARM64Ghash) processClmulResult(ACC0, ACCM, ACC1, ZERO, T *Vector128) {
	VEOR(ACC0, ACCM, ACCM)
	VEOR(ACC1, ACCM, ACCM)
	VEXT(8, ZERO, ACCM, T)
	VEXT(8, ACCM, ZERO, ACCM)
	VEOR(ACCM, ACC0, ACC0) // ACC0
	VEOR(T, ACC1, ACC1)    // ACC1
}

// Fast reduction for [ACC1:ACC0] by POLY and store the result in TARGET
// T is a temporary register
func (g *clmulARM64Ghash) fastReduction(TARGET, ACC1, ACC0, T, POLY *Vector128) {
	VPMULL(ACC0, POLY, T)
	VEXT(8, ACC0, ACC0, ACC0)
	VEOR(ACC0, T, ACC0)
	VPMULL(ACC0, POLY, T)
	VEXT(8, ACC0, ACC0, ACC0)
	VEOR(ACC0, T, ACC0)
	VEOR(ACC0, ACC1, TARGET)
}
