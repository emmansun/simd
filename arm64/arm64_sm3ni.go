package arm64

import (
	"encoding/binary"
	"math/bits"
)

const (
	_T0 = 0x79cc4519
	_T1 = 0x7a879d8a
)

var IV = [8]uint32{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600, 0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e}

func p0(x uint32) uint32 {
	return x ^ bits.RotateLeft32(x, 9) ^ bits.RotateLeft32(x, 17)
}

func p1(x uint32) uint32 {
	return x ^ (x<<15 | x>>17) ^ (x<<23 | x>>9)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3PARTW1--SM3PARTW1-?lang=en
func SM3PARTW1(Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	tmp := &Vector128{}
	VEOR(Vd, Vn, tmp)
	for i := 4; i < 16; i += 4 {
		v := binary.LittleEndian.Uint32(Vm.bytes[i:])
		v = bits.RotateLeft32(v, 15)
		binary.LittleEndian.PutUint32(result.bytes[i-4:], v)
	}
	VEOR(tmp, result, result)
	for i := 0; i < 16; i += 4 {
		if i == 12 {
			v := binary.LittleEndian.Uint32(tmp.bytes[i:])
			v ^= bits.RotateLeft32(binary.LittleEndian.Uint32(result.bytes[:]), 15)
			binary.LittleEndian.PutUint32(result.bytes[i:], v)
		}
		v := binary.LittleEndian.Uint32(result.bytes[i:])
		v = p1(v)
		binary.LittleEndian.PutUint32(result.bytes[i:], v)
	}
	copy(Vd.bytes[:], result.bytes[:])
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3PARTW2--SM3PARTW2-?lang=en
func SM3PARTW2(Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	tmp := &Vector128{}

	for i := 0; i < 16; i += 4 {
		v := binary.LittleEndian.Uint32(Vm.bytes[i:])
		v = bits.RotateLeft32(v, 7)
		binary.LittleEndian.PutUint32(tmp.bytes[i:], v)
	}
	VEOR(Vn, tmp, tmp)
	VEOR(Vd, tmp, result)

	tmp2 := bits.RotateLeft32(binary.LittleEndian.Uint32(tmp.bytes[:]), 15)
	tmp2 = p1(tmp2)
	tmp3 := binary.LittleEndian.Uint32(result.bytes[12:])
	binary.LittleEndian.PutUint32(result.bytes[12:], tmp2^tmp3)
	copy(Vd.bytes[:], result.bytes[:])
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3SS1--SM3SS1-?lang=en
// Va.S[3]: place T constant
// Vm.S[3]: sm3 state word E
// Vn.S[3]: sm3 state word A
func SM3SS1(Va, Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	for i := 0; i < 12; i++ {
		result.bytes[i] = 0
	}
	t1 := binary.LittleEndian.Uint32(Va.bytes[12:])
	t2 := binary.LittleEndian.Uint32(Vm.bytes[12:])
	t3 := binary.LittleEndian.Uint32(Vn.bytes[12:])
	t1 += t2 + bits.RotateLeft32(t3, 12)
	t1 = bits.RotateLeft32(t1, 7)
	binary.LittleEndian.PutUint32(result.bytes[12:], t1)
	copy(Vd.bytes[:], result.bytes[:])
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3TT1A--SM3TT1A-?lang=en
// Vm: Wj'
// Vn: SS1
// Vd: state
func SM3TT1A(imm byte, Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	imm = imm & 0x3
	WjPrime := binary.LittleEndian.Uint32(Vm.bytes[imm*4:])

	SS2 := binary.LittleEndian.Uint32(Vd.bytes[12:])
	SS2 = bits.RotateLeft32(SS2, 12)
	SS2 ^= binary.LittleEndian.Uint32(Vn.bytes[12:])

	TT1 := binary.LittleEndian.Uint32(Vd.bytes[4:]) ^ binary.LittleEndian.Uint32(Vd.bytes[8:]) ^ binary.LittleEndian.Uint32(Vd.bytes[12:])
	TT1 += SS2 + WjPrime + binary.LittleEndian.Uint32(Vd.bytes[:])

	binary.LittleEndian.PutUint32(result.bytes[:], binary.LittleEndian.Uint32(Vd.bytes[4:]))
	binary.LittleEndian.PutUint32(result.bytes[4:], bits.RotateLeft32(binary.LittleEndian.Uint32(Vd.bytes[8:]), 9))
	binary.LittleEndian.PutUint32(result.bytes[8:], binary.LittleEndian.Uint32(Vd.bytes[12:]))
	binary.LittleEndian.PutUint32(result.bytes[12:], TT1)
	copy(Vd.bytes[:], result.bytes[:])
}

func ff(x, y, z uint32) uint32 {
	return (x & y) | (x & z) | (y & z)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3TT1B--SM3TT1B-?lang=en
// Vm: Wj'
// Vn: SS1
// Vd: state
func SM3TT1B(imm byte, Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	imm = imm & 0x3
	WjPrime := binary.LittleEndian.Uint32(Vm.bytes[imm*4:])

	SS2 := binary.LittleEndian.Uint32(Vd.bytes[12:])
	SS2 = bits.RotateLeft32(SS2, 12)
	SS2 ^= binary.LittleEndian.Uint32(Vn.bytes[12:])

	d1 := binary.LittleEndian.Uint32(Vd.bytes[4:])
	d2 := binary.LittleEndian.Uint32(Vd.bytes[8:])
	d3 := binary.LittleEndian.Uint32(Vd.bytes[12:])
	TT1 := ff(d3, d2, d1)
	TT1 += SS2 + WjPrime + binary.LittleEndian.Uint32(Vd.bytes[:])

	binary.LittleEndian.PutUint32(result.bytes[:], binary.LittleEndian.Uint32(Vd.bytes[4:]))
	binary.LittleEndian.PutUint32(result.bytes[4:], bits.RotateLeft32(binary.LittleEndian.Uint32(Vd.bytes[8:]), 9))
	binary.LittleEndian.PutUint32(result.bytes[8:], binary.LittleEndian.Uint32(Vd.bytes[12:]))
	binary.LittleEndian.PutUint32(result.bytes[12:], TT1)
	copy(Vd.bytes[:], result.bytes[:])
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3TT2A--SM3TT2A-?lang=en
// Vm: Wj
// Vn: SS1
// Vd: state
func SM3TT2A(imm byte, Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	imm = imm & 0x3
	Wj := binary.LittleEndian.Uint32(Vm.bytes[imm*4:])

	TT2 := binary.LittleEndian.Uint32(Vd.bytes[4:]) ^ binary.LittleEndian.Uint32(Vd.bytes[8:]) ^ binary.LittleEndian.Uint32(Vd.bytes[12:])
	TT2 += Wj + binary.LittleEndian.Uint32(Vd.bytes[:]) + binary.LittleEndian.Uint32(Vn.bytes[12:])

	binary.LittleEndian.PutUint32(result.bytes[:], binary.LittleEndian.Uint32(Vd.bytes[4:]))
	binary.LittleEndian.PutUint32(result.bytes[4:], bits.RotateLeft32(binary.LittleEndian.Uint32(Vd.bytes[8:]), 19))
	binary.LittleEndian.PutUint32(result.bytes[8:], binary.LittleEndian.Uint32(Vd.bytes[12:]))
	binary.LittleEndian.PutUint32(result.bytes[12:], p0(TT2))
	copy(Vd.bytes[:], result.bytes[:])
}

func gg(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM3TT2B--SM3TT2B-?lang=en
// Vm: Wj
// Vn: SS1
// Vd: state
func SM3TT2B(imm byte, Vm, Vn, Vd *Vector128) {
	result := &Vector128{}
	imm = imm & 0x3
	Wj := binary.LittleEndian.Uint32(Vm.bytes[imm*4:])

	d1 := binary.LittleEndian.Uint32(Vd.bytes[4:])
	d2 := binary.LittleEndian.Uint32(Vd.bytes[8:])
	d3 := binary.LittleEndian.Uint32(Vd.bytes[12:])
	TT2 := gg(d3, d2, d1)
	TT2 += Wj + binary.LittleEndian.Uint32(Vd.bytes[:]) + binary.LittleEndian.Uint32(Vn.bytes[12:])

	binary.LittleEndian.PutUint32(result.bytes[:], binary.LittleEndian.Uint32(Vd.bytes[4:]))
	binary.LittleEndian.PutUint32(result.bytes[4:], bits.RotateLeft32(binary.LittleEndian.Uint32(Vd.bytes[8:]), 19))
	binary.LittleEndian.PutUint32(result.bytes[8:], binary.LittleEndian.Uint32(Vd.bytes[12:]))
	binary.LittleEndian.PutUint32(result.bytes[12:], p0(TT2))
	copy(Vd.bytes[:], result.bytes[:])
}

func sm3block(state *[8]uint32, p []byte) {
	var (
		V0 = &Vector128{}
		V1 = &Vector128{}
		V2 = &Vector128{}
		V3 = &Vector128{}
		V4 = &Vector128{}
		// for sm3 state
		V8  = &Vector128{}
		V9  = &Vector128{}
		V15 = &Vector128{}
		V16 = &Vector128{}
		// for T constants
		V11 = &Vector128{}
		V12 = &Vector128{}
	)

	// load sate
	VLD1_4S(state[:], V8)
	VLD1_4S(state[4:], V9)
	VREV64_S(V8, V8)
	VEXT(8, V8, V8, V8)
	VREV64_S(V9, V9)
	VEXT(8, V9, V9, V9)

	for len(p) >= 64 {
		// save last state
		VMOV(V8, V15)
		VMOV(V9, V16)
		//load data
		VLD1_16B(p, V0)
		VREV32_B(V0, V0)
		VLD1_16B(p[16:], V1)
		VREV32_B(V1, V1)
		VLD1_16B(p[32:], V2)
		VREV32_B(V2, V2)
		VLD1_16B(p[48:], V3)
		VREV32_B(V3, V3)

		// first 16 rounds
		// load T constants
		VLD1_4S([]uint32{0, 0, 0, _T0}, V11)
		qroundA(V11, V12, V8, V9, V0, V1, V2, V3, V4)
		qroundA(V11, V12, V8, V9, V1, V2, V3, V4, V0)
		qroundA(V11, V12, V8, V9, V2, V3, V4, V0, V1)
		qroundA(V11, V12, V8, V9, V3, V4, V0, V1, V2)

		// last 48 rounds
		// load T constants
		VLD1_4S([]uint32{0, 0, 0, 0x9d8a7a87}, V11)
		qroundB(V11, V12, V8, V9, V4, V0, V1, V2, V3)
		qroundB(V11, V12, V8, V9, V0, V1, V2, V3, V4)
		qroundB(V11, V12, V8, V9, V1, V2, V3, V4, V0)
		qroundB(V11, V12, V8, V9, V2, V3, V4, V0, V1)
		qroundB(V11, V12, V8, V9, V3, V4, V0, V1, V2)
		qroundB(V11, V12, V8, V9, V4, V0, V1, V2, V3)
		qroundB(V11, V12, V8, V9, V0, V1, V2, V3, V4)
		qroundB(V11, V12, V8, V9, V1, V2, V3, V4, V0)
		qroundB(V11, V12, V8, V9, V2, V3, V4, V0, V1)
		qroundC(V11, V12, V8, V9, V3, V4)
		qroundC(V11, V12, V8, V9, V4, V0)
		qroundC(V11, V12, V8, V9, V0, V1)

		VEOR(V8, V15, V8)
		VEOR(V9, V16, V9)

		p = p[64:]
	}
	VREV64_S(V8, V8)
	VEXT(8, V8, V8, V8)
	VREV64_S(V9, V9)
	VEXT(8, V9, V9, V9)
	VST1_4S(V8, state[:])
	VST1_4S(V9, state[4:])
}

// Compress 4 words and generate 4 words, used v6, v7, v10 as temp registers
// s4, used to store next 4 words
// s0, W(4i) W(4i+1) W(4i+2) W(4i+3)
// s1, W(4i+4) W(4i+5) W(4i+6) W(4i+7)
// s2, W(4i+8) W(4i+9) W(4i+10) W(4i+11)
// s3, W(4i+12) W(4i+13) W(4i+14) W(4i+15)
// t0, t1, T constant
// st1, st2, sm3 state
func qroundA(t0, t1, st1, st2, s0, s1, s2, s3, s4 *Vector128) {
	var (
		V6  = &Vector128{}
		V7  = &Vector128{}
		V10 = &Vector128{}
	)

	// Extension
	VEXT(3*4, s2, s1, s4) // w7,w8,w9,w10
	VEXT(3*4, s1, s0, V6) // w3,w4,w5,w6
	VEXT(2*4, s3, s2, V7) // w10,w11,w12,w13
	SM3PARTW1(s3, s0, s4)
	SM3PARTW2(V6, V7, s4) // s4 include W16, W17, W18, W19

	VEOR(s0, s1, V10) //v10 is W'

	// compression
	roundA(0, t0, t1, st1, st2, s0, V10)
	roundA(1, t1, t0, st1, st2, s0, V10)
	roundA(2, t0, t1, st1, st2, s0, V10)
	roundA(3, t1, t0, st1, st2, s0, V10)
}

// t0, t1, T constant
// st1, st2, sm3 state
// w, W(4i) W(4i+1) W(4i+2) W(4i+3)
// wt, W(4i)' W(4i+1)' W(4i+2)' W(4i+3)'
func roundA(i byte, t0, t1, st1, st2, w, wt *Vector128) {
	V5 := &Vector128{}
	SM3SS1(t0, st2, st1, V5) // V5 is SS1
	VSHL_S(1, t0, t1)
	VSRI_S(31, t0, t1)
	SM3TT1A(i, wt, V5, st1) // TT1
	SM3TT2A(i, w, V5, st2)  // TT2
}

func qroundB(t0, t1, st1, st2, s0, s1, s2, s3, s4 *Vector128) {
	var (
		V6  = &Vector128{}
		V7  = &Vector128{}
		V10 = &Vector128{}
	)

	// Extension
	VEXT(3*4, s2, s1, s4) // w7,w8,w9,w10
	VEXT(3*4, s1, s0, V6) // w3,w4,w5,w6
	VEXT(2*4, s3, s2, V7) // w10,w11,w12,w13
	SM3PARTW1(s3, s0, s4)
	SM3PARTW2(V6, V7, s4) // s4 include W16, W17, W18, W19

	VEOR(s0, s1, V10) //v10 is W'

	// compression
	roundB(0, t0, t1, st1, st2, s0, V10)
	roundB(1, t1, t0, st1, st2, s0, V10)
	roundB(2, t0, t1, st1, st2, s0, V10)
	roundB(3, t1, t0, st1, st2, s0, V10)
}

func qroundC(t0, t1, st1, st2, s0, s1 *Vector128) {
	var (
		V10 = &Vector128{}
	)

	VEOR(s0, s1, V10) //v10 is W'
	// compression
	roundB(0, t0, t1, st1, st2, s0, V10)
	roundB(1, t1, t0, st1, st2, s0, V10)
	roundB(2, t0, t1, st1, st2, s0, V10)
	roundB(3, t1, t0, st1, st2, s0, V10)
}

// t0, t1, T constant
// st1, st2, sm3 state
// w, W(4i) W(4i+1) W(4i+2) W(4i+3)
// wt, W(4i)' W(4i+1)' W(4i+2)' W(4i+3)'
func roundB(i byte, t0, t1, st1, st2, w, wt *Vector128) {
	V5 := &Vector128{}
	SM3SS1(t0, st2, st1, V5) // V5 is SS1
	VSHL_S(1, t0, t1)
	VSRI_S(31, t0, t1)
	SM3TT1B(i, wt, V5, st1) // TT1
	SM3TT2B(i, w, V5, st2)  // TT2
}
