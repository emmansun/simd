// https://cdrdv2-public.intel.com/782879/architecture-instruction-set-extensions-programming-reference.pdf
package avx

import (
	"math/bits"

	"github.com/emmansun/simd/alg/sm3"
	"github.com/emmansun/simd/amd64/sse"
)

// Perform Initial Calculation for the Next Four SM3 Message Words
// The VSM3MSG1 instruction is one of the two SM3 message scheduling instructions.
// The instruction performs an initial calculation for the next four SM3 message words.
func VSM3MSG1(srcdst, src1, src2 *sse.XMM) {
	w0_3 := src2.Uint32s()
	w7_10 := srcdst.Uint32s()
	w13_15 := src1.Uint32s()

	tmp0 := w7_10[0] ^ w0_3[0] ^ bits.RotateLeft32(w13_15[0], 15)
	tmp1 := w7_10[1] ^ w0_3[1] ^ bits.RotateLeft32(w13_15[1], 15)
	tmp2 := w7_10[2] ^ w0_3[2] ^ bits.RotateLeft32(w13_15[2], 15)
	tmp3 := w7_10[3] ^ w0_3[3]

	VMOVDQU_L4S(srcdst, []uint32{sm3.P1(tmp0), sm3.P1(tmp1), sm3.P1(tmp2), sm3.P1(tmp3)})
}

// Perform Final Calculation for the Next Four SM3 Message Words
// The VSM3MSG2 instruction is one of the two SM3 message scheduling instructions.
// The instruction performs the final calculation for the next four SM3 message words
func VSM3MSG2(srcdst, src1, src2 *sse.XMM) {
	wtmp := srcdst.Uint32s()
	w3_6 := src1.Uint32s()
	w10_13 := src2.Uint32s()

	w16 := bits.RotateLeft32(w3_6[0], 7) ^ w10_13[0] ^ wtmp[0]
	w17 := bits.RotateLeft32(w3_6[1], 7) ^ w10_13[1] ^ wtmp[1]
	w18 := bits.RotateLeft32(w3_6[2], 7) ^ w10_13[2] ^ wtmp[2]
	w19 := bits.RotateLeft32(w3_6[3], 7) ^ w10_13[3] ^ wtmp[3]

	w19 ^= bits.RotateLeft32(w16, 6) ^ bits.RotateLeft32(w16, 15) ^ bits.RotateLeft32(w16, 30)

	VMOVDQU_L4S(srcdst, []uint32{w16, w17, w18, w19})
}

func SM3MSG(out, w0, w1, w2, w3, t0, t1 *sse.XMM) {
	// prepare data for VSM3MSG1
	VPALIGNR(out, w2, w1, 12) // out = W10 W9 W8 W7
	VPSRLDQ(t0, w3, 4)        // t0 = 0 W15 W14 W13
	VSM3MSG1(out, t0, w0)     // out = WTMP3 WTMP2 WTMP1 WTMP0
	// prepare data for VSM3MSG2
	VPALIGNR(t0, w1, w0, 12) // t0 = W6 W5 W4 W3
	VPALIGNR(t1, w3, w2, 8)  // t1 = W13 W12 W11 W10
	VSM3MSG2(out, t0, t1)    // out = W19 W18 W17 W16
}

// Perform Two Rounds of SM3 Operation
// The VSM3RNDS2 instruction performs two rounds of SM3 operation using initial SM3 state (C, D, G, H) from the
// first operand, an initial SM3 states (A, B, E, F) from the second operand and a pre-computed words from the third
// operand. The first operand with initial SM3 state of (C, D, G, H) assumes input of non-rotated left variables from
// previous state. The updated SM3 state (A, B, E, F) is written to the first operand.
// The imm8 should contain the even round number for the first of the two rounds computed by this instruction. The
// computation masks the imm8 value by AND’ing it with 0x3E so that only even round numbers from 0 through 62
// are used for this operation.
func VSM3RNDS2(srcdst, src1, src2 *sse.XMM, imm8 byte) {
	var A, B, C, D, E, F, G, H [3]uint32
	var W [6]uint32
	FEBA := src1.Uint32s()
	A[0] = FEBA[3]
	B[0] = FEBA[2]
	E[0] = FEBA[1]
	F[0] = FEBA[0]
	HGDC := srcdst.Uint32s()
	C[0] = HGDC[3]
	D[0] = HGDC[2]
	G[0] = HGDC[1]
	H[0] = HGDC[0]
	w := src2.Uint32s() // W0, W1, W4, W5
	W[0] = w[0]
	W[1] = w[1]
	W[4] = w[2]
	W[5] = w[3]

	C[0] = bits.RotateLeft32(C[0], 9)
	D[0] = bits.RotateLeft32(D[0], 9)
	G[0] = bits.RotateLeft32(G[0], 19)
	H[0] = bits.RotateLeft32(H[0], 19)

	var CONST uint32
	ROUND := imm8 & 0x3E // even round number 0...62
	if ROUND < 16 {
		CONST = sm3.CONST0
	} else {
		CONST = sm3.CONST1
	}
	CONST = bits.RotateLeft32(CONST, int(ROUND))

	for i := 0; i < 2; i++ {
		S1 := bits.RotateLeft32(bits.RotateLeft32(A[i], 12)+E[i]+CONST, 7)
		S2 := S1 ^ bits.RotateLeft32(A[i], 12)
		T1 := sm3.FF(ROUND, A[i], B[i], C[i]) + D[i] + S2 + (W[i] ^ W[i+4])
		T2 := sm3.GG(ROUND, E[i], F[i], G[i]) + H[i] + S1 + W[i]
		D[i+1] = C[i]
		C[i+1] = bits.RotateLeft32(B[i], 9)
		B[i+1] = A[i]
		A[i+1] = T1
		H[i+1] = G[i]
		G[i+1] = bits.RotateLeft32(F[i], 19)
		F[i+1] = E[i]
		E[i+1] = sm3.P0(T2)
		CONST = bits.RotateLeft32(CONST, 1)
	}
	VMOVDQU_L4S(srcdst, []uint32{F[2], E[2], B[2], A[2]})
}

func SM3RNDS4(stateABEF, stateCDGH, w0, w4, t *sse.XMM, imm8 byte) {
	VPUNPCKLQDQ(t, w0, w4) // t = W5 W4 W1 W0
	VSM3RNDS2(stateCDGH, stateABEF, t, imm8)
	VPUNPCKHQDQ(t, w0, w4) // t = W7 W6 W3 W2
	VSM3RNDS2(stateABEF, stateCDGH, t, imm8+2)
}

func prepareStates(X0, X1, X2, X3, X4 *sse.XMM, state *[8]uint32) {
	VMOVDQU_L4S(X1, state[:4]) // X1 = D C B A
	VMOVDQU_L4S(X2, state[4:]) // X2 = H G F E
	VPSHUFD(X1, X1, 0xB1)      // X1 = C D A B
	VPSHUFD(X2, X2, 0xB1)      // X2 = G H E F
	VPUNPCKLQDQ(X0, X2, X1)    // X0 = A B E F
	VPUNPCKHQDQ(X1, X2, X1)    // X1 = C D G H
	VPSRLD(X2, X1, 9)
	VPSLLD(X3, X1, 23)
	VPXOR(X2, X2, X3) // X2 = ROR32(CDGH, 9)
	VPSRLD(X3, X1, 19)
	VPSLLD(X4, X1, 13)
	VPXOR(X3, X3, X4)         // X3 = ROR32(CDGH, 19)
	VPBLENDD(X1, X2, X3, 0x3) // X1 = ROR32(C, 9) ROR32(D, 9) ROR32(G, 19) ROR32(H, 19)
}

// X0 = A B E F, X1 = C D G H
func saveStates(X0, X1, X2, X3, X4 *sse.XMM, state *[8]uint32) {
	VPSLLD(X2, X1, 9)
	VPSRLD(X3, X1, 23)
	VPXOR(X2, X2, X3) // X2 = ROL32(CDGH, 9)

	VPSLLD(X3, X1, 19)
	VPSRLD(X4, X1, 13)
	VPXOR(X3, X3, X4) // X3 = ROL32(CDGH, 19)

	VPBLENDD(X1, X2, X3, 0x3) // X1 = ROL32(C, 9) ROL32(D, 9) ROL32(G, 19) ROL32(H, 19)

	VPSHUFD(X0, X0, 0xB1) // X0 = B A F E
	VPSHUFD(X1, X1, 0xB1) // X1 = D C H G

	VPUNPCKHQDQ(X2, X0, X1) // X2 = D C B A
	VPUNPCKLQDQ(X3, X0, X1) // X3 = H G F E

	VMOVDQU_S4S(state[:4], X2)
	VMOVDQU_S4S(state[4:], X3)
}

func sm3block(state *[8]uint32, p []byte) {
	var (
		X0        = &sse.XMM{}
		X1        = &sse.XMM{}
		X2        = &sse.XMM{}
		X3        = &sse.XMM{}
		X4        = &sse.XMM{}
		X5        = &sse.XMM{}
		X6        = &sse.XMM{}
		X7        = &sse.XMM{}
		X8        = &sse.XMM{}
		X10       = &sse.XMM{}
		X11       = &sse.XMM{}
		flip_mask = &sse.XMM{}
	)
	VMOVDQU_L16B(flip_mask, []byte{0x03, 0x02, 0x01, 0x00, 0x07, 0x06, 0x05, 0x04, 0x0b, 0x0a, 0x09, 0x08, 0x0f, 0x0e, 0x0d, 0x0c})
	// prepare initial state to X0, X1
	prepareStates(X0, X1, X2, X3, X4, state)

	for len(p) >= 64 {
		// save state
		VMOVDQU(X10, X0)
		VMOVDQU(X11, X1)
		// load data
		VMOVDQU_L16B(X2, p)
		VMOVDQU_L16B(X3, p[16:])
		VMOVDQU_L16B(X4, p[32:])
		VMOVDQU_L16B(X5, p[48:])
		VPSHUFB(X2, X2, flip_mask)
		VPSHUFB(X3, X3, flip_mask)
		VPSHUFB(X4, X4, flip_mask)
		VPSHUFB(X5, X5, flip_mask)

		// message schedule & compress
		SM3MSG(X6, X2, X3, X4, X5, X7, X8) // X6 = W19 W18 W17 W16
		SM3RNDS4(X0, X1, X2, X3, X7, 0)

		SM3MSG(X2, X3, X4, X5, X6, X7, X8) // X2 = W23 W22 W21 W20
		SM3RNDS4(X0, X1, X3, X4, X7, 4)

		SM3MSG(X3, X4, X5, X6, X2, X7, X8) // X3 = W27 W26 W25 W24
		SM3RNDS4(X0, X1, X4, X5, X7, 8)

		SM3MSG(X4, X5, X6, X2, X3, X7, X8) // X4 = W31 W30 W29 W28
		SM3RNDS4(X0, X1, X5, X6, X7, 12)

		SM3MSG(X5, X6, X2, X3, X4, X7, X8) // X5 = W35 W34 W33 W32
		SM3RNDS4(X0, X1, X6, X2, X7, 16)

		SM3MSG(X6, X2, X3, X4, X5, X7, X8) // X6 = W39 W38 W37 W36
		SM3RNDS4(X0, X1, X2, X3, X7, 20)

		SM3MSG(X2, X3, X4, X5, X6, X7, X8) // X2 = W43 W42 W41 W40
		SM3RNDS4(X0, X1, X3, X4, X7, 24)

		SM3MSG(X3, X4, X5, X6, X2, X7, X8) // X3 = W47 W46 W45 W44
		SM3RNDS4(X0, X1, X4, X5, X7, 28)

		SM3MSG(X4, X5, X6, X2, X3, X7, X8) // X4 = W51 W50 W49 W48
		SM3RNDS4(X0, X1, X5, X6, X7, 32)

		SM3MSG(X5, X6, X2, X3, X4, X7, X8) // X5 = W55 W54 W53 W52
		SM3RNDS4(X0, X1, X6, X2, X7, 36)

		SM3MSG(X6, X2, X3, X4, X5, X7, X8) // X6 = W59 W58 W57 W56
		SM3RNDS4(X0, X1, X2, X3, X7, 40)

		SM3MSG(X2, X3, X4, X5, X6, X7, X8) // X2 = W63 W62 W61 W60
		SM3RNDS4(X0, X1, X3, X4, X7, 44)

		SM3MSG(X3, X4, X5, X6, X2, X7, X8) // X3 = W67 W66 W65 W64
		SM3RNDS4(X0, X1, X4, X5, X7, 48)

		SM3RNDS4(X0, X1, X5, X6, X7, 52)
		SM3RNDS4(X0, X1, X6, X2, X7, 56)
		SM3RNDS4(X0, X1, X2, X3, X7, 60)

		// update state
		VPXOR(X0, X0, X10)
		VPXOR(X1, X1, X11)

		p = p[64:]
	}
	saveStates(X0, X1, X2, X3, X4, state)
}
