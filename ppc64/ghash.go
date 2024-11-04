// The theoretical proof behind it is still not very clear.
package ppc64

import "encoding/binary"

type clmulPPC64Ghash struct {
	isPPC64LE         bool
	bytesProductTable [400]byte // (8 * 3 + 1) * 16
}

func NewClmulPPC64Ghash(h []byte, isPPC64LE bool) *clmulPPC64Ghash {
	g := &clmulPPC64Ghash{}
	g.isPPC64LE = isPPC64LE
	var (
		XC2  = &Vector128{}
		T0   = &Vector128{}
		T1   = &Vector128{}
		T2   = &Vector128{}
		ZERO = &Vector128{}
		H    = &Vector128{}
		HL   = &Vector128{}
		HH   = &Vector128{}
		IN   = &Vector128{}
		XL   = &Vector128{}
		XM   = &Vector128{}
		XH   = &Vector128{}
		H2   = &Vector128{}
		H2L  = &Vector128{}
		H2H  = &Vector128{}
	)

	var h1, h2 uint64
	if isPPC64LE {
		// can use VPERM to convert from little-endian to big-endian
		var v [16]byte
		h1 = binary.LittleEndian.Uint64(h[:8])
		h2 = binary.LittleEndian.Uint64(h[8:])
		binary.BigEndian.PutUint64(v[:], h1)
		binary.BigEndian.PutUint64(v[8:], h2)
		LXVD2X_PPC64LE(v[:], H)
	} else {
		LXVD2X(h, H)
	}

	VXOR(ZERO, ZERO, ZERO)
	g.initPoly(XC2, ZERO, T0, T1)

	// Multiply by 2 modulo P
	VSPLTISB(7, T2)
	VSPLTB(0, H, T1)  // most significant byte
	VSL(H, T0, H)     // H<<=1
	VSRAB(T1, T2, T1) // broadcast carry bit
	VAND(T1, XC2, T1)
	VXOR(H, T1, IN)        // twisted H
	VSLDOI(8, IN, IN, H)   // twist even more ...
	VSLDOI(8, ZERO, H, HL) // ... and split
	VSLDOI(8, H, ZERO, HH)

	VSLDOI(8, ZERO, XC2, XC2) // 0xc2.0
	if isPPC64LE {
		STXVD2X_PPC64LE(XC2, g.bytesProductTable[24*16:])
		STXVD2X_PPC64LE(HL, g.bytesProductTable[21*16:])
		STXVD2X_PPC64LE(H, g.bytesProductTable[22*16:])
		STXVD2X_PPC64LE(HH, g.bytesProductTable[23*16:])
	} else {
		STXVD2X(XC2, g.bytesProductTable[24*16:])
		STXVD2X(HL, g.bytesProductTable[21*16:])
		STXVD2X(H, g.bytesProductTable[22*16:])
		STXVD2X(HH, g.bytesProductTable[23*16:])
	}

	for i := 6; i >= 0; i-- {
		// Multiplication
		VPMSUMD(IN, HL, XL) // H.lo·H.lo
		VPMSUMD(IN, H, XM)  // H.hi·H.lo+H.lo·H.hi
		VPMSUMD(IN, HH, XH) // H.hi·H.hi

		g.processClmulResult(XL, XM, XH, ZERO, T0)
		g.fastReduction(IN, XH, XL, T0, XC2)

		VSLDOI(8, IN, IN, H2)
		VSLDOI(8, ZERO, H2, H2L)
		VSLDOI(8, H2, ZERO, H2H)

		if isPPC64LE {
			STXVD2X_PPC64LE(H2L, g.bytesProductTable[(3*i)*16:])
			STXVD2X_PPC64LE(H2, g.bytesProductTable[(3*i+1)*16:])
			STXVD2X_PPC64LE(H2H, g.bytesProductTable[(3*i+2)*16:])
		} else {
			STXVD2X(H2L, g.bytesProductTable[(3*i)*16:])
			STXVD2X(H2, g.bytesProductTable[(3*i+1)*16:])
			STXVD2X(H2H, g.bytesProductTable[(3*i+2)*16:])
		}
	}

	return g
}

func (g *clmulPPC64Ghash) Hash(T *[16]byte, data []byte) {
	var (
		XC2   = &Vector128{}
		XPERM = &Vector128{}
		T0    = &Vector128{}
		ZERO  = &Vector128{}
		H     = &Vector128{}
		HL    = &Vector128{}
		HH    = &Vector128{}
		ACC0  = &Vector128{}
		ACC1  = &Vector128{}
		ACCM  = &Vector128{}
		B0    = &Vector128{}
		B1    = &Vector128{}
		B2    = &Vector128{}
		B3    = &Vector128{}
		B4    = &Vector128{}
		B5    = &Vector128{}
		B6    = &Vector128{}
		B7    = &Vector128{}
	)
	if g.isPPC64LE {
		LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, XPERM)
	}
	VXOR(ZERO, ZERO, ZERO)
	VXOR(ACC0, ACC0, ACC0)
	g.loadPrecomputed(24, XC2)

	// handle 8 blocks at a time
	for len(data) >= 128 {
		// load 8 blocks
		g.load8Blocks(data, B0, B1, B2, B3, B4, B5, B6, B7, XPERM)

		// load precomputed values
		g.loadPrecomputed(0, HL)
		g.loadPrecomputed(1, H)
		g.loadPrecomputed(2, HH)

		// process first block
		// add previous result
		VXOR(B0, ACC0, B0)
		// multiplication
		VPMSUMD(B0, HL, ACC0)
		VPMSUMD(B0, H, ACCM)
		VPMSUMD(B0, HH, ACC1)

		g.mulRoundAAD(B1, ACC0, ACC1, ACCM, HL, H, HH, T0, 1)
		g.mulRoundAAD(B2, ACC0, ACC1, ACCM, HL, H, HH, T0, 2)
		g.mulRoundAAD(B3, ACC0, ACC1, ACCM, HL, H, HH, T0, 3)
		g.mulRoundAAD(B4, ACC0, ACC1, ACCM, HL, H, HH, T0, 4)
		g.mulRoundAAD(B5, ACC0, ACC1, ACCM, HL, H, HH, T0, 5)
		g.mulRoundAAD(B6, ACC0, ACC1, ACCM, HL, H, HH, T0, 6)
		g.mulRoundAAD(B7, ACC0, ACC1, ACCM, HL, H, HH, T0, 7)

		g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
		g.fastReduction(ACC0, ACC1, ACC0, T0, XC2)

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
			VPERM(B0, B0, XPERM, B0)
		}
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, HL, H, HH, T0, XC2, ZERO)
		data = data[16:]
	}
	if len(data) > 0 {
		var partialBlock [16]byte
		copy(partialBlock[:], data)
		g.loadData(partialBlock[:], 0, B0)
		if g.isPPC64LE {
			VPERM(B0, B0, XPERM, B0)
		}
		g.mulOneBlock(B0, ACC0, ACCM, ACC1, HL, H, HH, T0, XC2, ZERO)
	}
	if g.isPPC64LE {
		VPERM(ACC0, ACC0, XPERM, ACC0)
		STXVD2X_PPC64LE(ACC0, T[:])
	} else {
		STXVD2X(ACC0, T[:])
	}
}

func (g *clmulPPC64Ghash) load8Blocks(data []byte, B0, B1, B2, B3, B4, B5, B6, B7, XPERM *Vector128) {
	g.loadData(data, 0, B0)
	g.loadData(data, 1, B1)
	g.loadData(data, 2, B2)
	g.loadData(data, 3, B3)
	g.loadData(data, 4, B4)
	g.loadData(data, 5, B5)
	g.loadData(data, 6, B6)
	g.loadData(data, 7, B7)

	if g.isPPC64LE {
		VPERM(B0, B0, XPERM, B0)
		VPERM(B1, B1, XPERM, B1)
		VPERM(B2, B2, XPERM, B2)
		VPERM(B3, B3, XPERM, B3)
		VPERM(B4, B4, XPERM, B4)
		VPERM(B5, B5, XPERM, B5)
		VPERM(B6, B6, XPERM, B6)
		VPERM(B7, B7, XPERM, B7)
	}
}

func (g *clmulPPC64Ghash) mulOneBlock(B0, ACC0, ACCM, ACC1, HL, H, HH, T0, XC2, ZERO *Vector128) {
	VXOR(B0, ACC0, B0)
	// Multiplication
	VPMSUMD(B0, HL, ACC0)
	VPMSUMD(B0, H, ACCM)
	VPMSUMD(B0, HH, ACC1)
	g.processClmulResult(ACC0, ACCM, ACC1, ZERO, T0)
	g.fastReduction(ACC0, ACC1, ACC0, T0, XC2)
}

func (g *clmulPPC64Ghash) mulRoundAAD(IN, ACC0, ACC1, ACCM, T0, T1, T2, T3 *Vector128, i int) {
	// load precomputed values
	g.loadPrecomputed(3*i, T0)   // H2L
	g.loadPrecomputed(3*i+1, T1) // H2
	g.loadPrecomputed(3*i+2, T2) // H2H

	// Multiplication and accumulate
	VPMSUMD(IN, T0, T3)  // H2L
	VXOR(ACC0, T3, ACC0) // ACC0 = ACC0 ^ H2L
	VPMSUMD(IN, T1, T3)
	VXOR(ACCM, T3, ACCM) // ACCM = ACCM ^ H2
	VPMSUMD(IN, T2, T3)
	VXOR(ACC1, T3, ACC1) // ACC1 = ACC1 ^ H2H
}

func (g *clmulPPC64Ghash) loadData(data []byte, idx int, T *Vector128) {
	if g.isPPC64LE {
		LXVD2X_PPC64LE(data[idx*16:], T)
	} else {
		LXVD2X(data[idx*16:], T)
	}
}

func (g *clmulPPC64Ghash) loadPrecomputed(i int, T *Vector128) {
	if g.isPPC64LE {
		LXVD2X_PPC64LE(g.bytesProductTable[i*16:], T)
	} else {
		LXVD2X(g.bytesProductTable[i*16:], T)
	}
}

func (g *clmulPPC64Ghash) initPoly(POLY, ZERO, T0, T1 *Vector128) {
	// compute POLY = 0xc2000000000000000000000000000001
	VSPLTISB(0x10, POLY)         // 0xf0
	VSPLTISB(1, T0)              // one
	VADDUBM(POLY, POLY, POLY)    // 0xe0
	VOR(POLY, T0, POLY)          // 0xe1
	VSLDOI(15, POLY, ZERO, POLY) // 0xe1...
	VSLDOI(1, ZERO, T0, T1)      // ...1
	VADDUBM(POLY, POLY, POLY)    // 0xc2...
	VOR(POLY, T1, POLY)          // 0xc2....01
}

func (g *clmulPPC64Ghash) processClmulResult(ACC0, ACCM, ACC1, ZERO, T *Vector128) {
	VSLDOI(8, ACCM, ZERO, T) // T = ACCM.hi || zero
	VXOR(ACC0, T, ACC0)      // ACC0 = ACC0 ^ (ACCM.hi || zero) = ACCM.hi || ACC0.lo
	VSLDOI(8, ZERO, ACCM, T) // T = zero || ACCM.lo
	VXOR(ACC1, T, ACC1)      // ACC1 = ACC1 ^ (zero || ACCM.lo) = ACC1.hi || ACCM.lo
}

// Fast reduction for [ACC1:ACC0] by POLY and store the result in TARGET
// T is a temporary register
func (g *clmulPPC64Ghash) fastReduction(TARGET, ACC1, ACC0, T, POLY *Vector128) {
	VPMSUMD(ACC0, POLY, T) // POLY.hi = 0
	VSLDOI(8, ACC0, ACC0, ACC0)
	VXOR(ACC0, T, ACC0)
	VPMSUMD(ACC0, POLY, T) // POLY.hi = 0
	VSLDOI(8, ACC0, ACC0, ACC0)
	VXOR(ACC0, T, ACC0)
	VXOR(ACC0, ACC1, TARGET)
}
