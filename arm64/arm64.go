package arm64

import "encoding/binary"

type Vector128 struct {
	bytes [16]byte
}

func (m *Vector128) Bytes() []byte {
	return m.bytes[:]
}

func (m *Vector128) Uint64s() []uint64 {
	return []uint64{binary.LittleEndian.Uint64(m.bytes[:]), binary.LittleEndian.Uint64(m.bytes[8:])}
}

func (m *Vector128) Uint32s() []uint32 {
	return []uint32{binary.LittleEndian.Uint32(m.bytes[:]), binary.LittleEndian.Uint32(m.bytes[4:]), binary.LittleEndian.Uint32(m.bytes[8:]), binary.LittleEndian.Uint32(m.bytes[12:])}
}

func (m *Vector128) Uint16s() []uint16 {
	return []uint16{binary.LittleEndian.Uint16(m.bytes[:]), binary.LittleEndian.Uint16(m.bytes[2:]), binary.LittleEndian.Uint16(m.bytes[4:]), binary.LittleEndian.Uint16(m.bytes[6:]), binary.LittleEndian.Uint16(m.bytes[8:]), binary.LittleEndian.Uint16(m.bytes[10:]), binary.LittleEndian.Uint16(m.bytes[12:]), binary.LittleEndian.Uint16(m.bytes[14:])}
}

func VMOV(src *Vector128, dst *Vector128) {
	copy(dst.bytes[:], src.bytes[:])
}

func VMOV_S(src, dst *Vector128, from, to byte) {
	from = from & 0x3
	to = to & 0x3
	binary.LittleEndian.PutUint32(dst.bytes[to*4:], binary.LittleEndian.Uint32(src.bytes[from*4:]))
}

func VLD1_16B(rawbytes []byte, dst *Vector128) {
	copy(dst.bytes[:], rawbytes)
}

func VLD1_8H(v []uint16, dst *Vector128) {
	binary.LittleEndian.PutUint16(dst.bytes[:], v[0])
	binary.LittleEndian.PutUint16(dst.bytes[2:], v[1])
	binary.LittleEndian.PutUint16(dst.bytes[4:], v[2])
	binary.LittleEndian.PutUint16(dst.bytes[6:], v[3])
	binary.LittleEndian.PutUint16(dst.bytes[8:], v[4])
	binary.LittleEndian.PutUint16(dst.bytes[10:], v[5])
	binary.LittleEndian.PutUint16(dst.bytes[12:], v[6])
	binary.LittleEndian.PutUint16(dst.bytes[14:], v[7])
}

func VLD1_4S(v []uint32, dst *Vector128) {
	binary.LittleEndian.PutUint32(dst.bytes[:], v[0])
	binary.LittleEndian.PutUint32(dst.bytes[4:], v[1])
	binary.LittleEndian.PutUint32(dst.bytes[8:], v[2])
	binary.LittleEndian.PutUint32(dst.bytes[12:], v[3])
}

func VLD1_2D(v []uint64, dst *Vector128) {
	binary.LittleEndian.PutUint64(dst.bytes[:], v[0])
	binary.LittleEndian.PutUint64(dst.bytes[8:], v[1])
}

func VLD2_16B(rawbytes []byte, dst1, dst2 *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst1.bytes[i] = rawbytes[2*i]
		dst2.bytes[i] = rawbytes[2*i+1]
	}
}

func VLD3_16B(rawbytes []byte, dst1, dst2, dst3 *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst1.bytes[i] = rawbytes[3*i]
		dst2.bytes[i] = rawbytes[3*i+1]
		dst3.bytes[i] = rawbytes[3*i+2]
	}
}

// vld4q_u8
func VLD4_16B(rawbytes []byte, dst1, dst2, dst3, dst4 *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst1.bytes[i] = rawbytes[4*i]
		dst2.bytes[i] = rawbytes[4*i+1]
		dst3.bytes[i] = rawbytes[4*i+2]
		dst4.bytes[i] = rawbytes[4*i+3]
	}
}

func VDUP_BYTE(src byte, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src
	}
}

func VDUP_S(src uint32, dst *Vector128) {
	binary.LittleEndian.PutUint32(dst.bytes[:], src)
	binary.LittleEndian.PutUint32(dst.bytes[4:], src)
	binary.LittleEndian.PutUint32(dst.bytes[8:], src)
	binary.LittleEndian.PutUint32(dst.bytes[12:], src)
}

func VST1_16B(src *Vector128, dst []byte) {
	copy(dst, src.bytes[:])
}

func VST1_4S(src *Vector128, dst []uint32) {
	dst[0] = binary.LittleEndian.Uint32(src.bytes[:])
	dst[1] = binary.LittleEndian.Uint32(src.bytes[4:])
	dst[2] = binary.LittleEndian.Uint32(src.bytes[8:])
	dst[3] = binary.LittleEndian.Uint32(src.bytes[12:])
}

func VST2_16B(src1, src2 *Vector128, dst []byte) {
	for i := 0; i < 16; i += 1 {
		dst[2*i] = src1.bytes[i]
		dst[2*i+1] = src2.bytes[i]
	}
}

func VST3_16B(src1, src2, src3 *Vector128, dst []byte) {
	for i := 0; i < 16; i += 1 {
		dst[3*i] = src1.bytes[i]
		dst[3*i+1] = src2.bytes[i]
		dst[3*i+2] = src3.bytes[i]
	}
}

func VST4_16B(src1, src2, src3, src4 *Vector128, dst []byte) {
	for i := 0; i < 16; i += 1 {
		dst[4*i] = src1.bytes[i]
		dst[4*i+1] = src2.bytes[i]
		dst[4*i+2] = src3.bytes[i]
		dst[4*i+3] = src4.bytes[i]
	}
}

func VREV16(src, dst *Vector128) {
	tmp := Vector128{}
	for i := 0; i < 16; i += 2 {
		tmp.bytes[i] = src.bytes[i+1]
		tmp.bytes[i+1] = src.bytes[i]
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VREV32_B(src, dst *Vector128) {
	tmp := Vector128{}
	for i := 0; i < 16; i += 4 {
		tmp.bytes[i] = src.bytes[i+3]
		tmp.bytes[i+1] = src.bytes[i+2]
		tmp.bytes[i+2] = src.bytes[i+1]
		tmp.bytes[i+3] = src.bytes[i]
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VREV64_B(src, dst *Vector128) {
	tmp := Vector128{}
	for i := 0; i < 16; i += 8 {
		tmp.bytes[i] = src.bytes[i+7]
		tmp.bytes[i+1] = src.bytes[i+6]
		tmp.bytes[i+2] = src.bytes[i+5]
		tmp.bytes[i+3] = src.bytes[i+4]
		tmp.bytes[i+4] = src.bytes[i+3]
		tmp.bytes[i+5] = src.bytes[i+2]
		tmp.bytes[i+6] = src.bytes[i+1]
		tmp.bytes[i+7] = src.bytes[i]
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VREV64_S(src, dst *Vector128) {
	tmp := Vector128{}
	for i := 0; i < 16; i += 8 {
		w1 := binary.LittleEndian.Uint32(src.bytes[i:])
		w2 := binary.LittleEndian.Uint32(src.bytes[i+4:])
		binary.LittleEndian.PutUint32(tmp.bytes[i:], w2)
		binary.LittleEndian.PutUint32(tmp.bytes[i+4:], w1)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

// Extract vector from pair of vectors
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/EXT--Extract-vector-from-pair-of-vectors-?lang=en
func VEXT(imm byte, Vm, Vn, Vd *Vector128) {
	imm = imm & 0xf
	start := 16 - int(imm)
	tmp := Vector128{}
	for i := 0; i < int(imm); i++ {
		tmp.bytes[start+i] = Vm.bytes[i]
	}
	for i := int(imm); i < 16; i++ {
		tmp.bytes[i-int(imm)] = Vn.bytes[i]
	}
	copy(Vd.bytes[:], tmp.bytes[:])
}

func VAND(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src1.bytes[i] & src2.bytes[i]
	}
}

func VORR(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src1.bytes[i] | src2.bytes[i]
	}
}

func VEOR(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src1.bytes[i] ^ src2.bytes[i]
	}
}

func VUSHR_B(imm byte, src, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src.bytes[i] >> imm
	}
}

func VUSHR_S(imm byte, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	v := src.Uint32s()
	VLD1_4S([]uint32{v[0] >> imm, v[1] >> imm, v[2] >> imm, v[3] >> imm}, dst)
}

func VUSHR_D(imm byte, src, dst *Vector128) {
	if imm > 63 {
		imm = 63
	}
	v := src.Uint64s()
	VLD1_2D([]uint64{v[0] >> imm, v[1] >> imm}, dst)
}

// Vector Shift Left and Insert
func VSLI_B(imm byte, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	mask := byte(0xff) >> (8 - imm)
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = (mask & dst.bytes[i]) | (src.bytes[i] << imm)
	}
}

func VSRI_S(imm byte, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	mask := uint32(0xffffffff) << (32 - imm)
	for i := 0; i < 16; i += 4 {
		v := binary.LittleEndian.Uint32(dst.bytes[i:])
		binary.LittleEndian.PutUint32(dst.bytes[i:], (mask&v)|(binary.LittleEndian.Uint32(src.bytes[i:])>>imm))
	}
}

func VSHL_B(imm byte, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src.bytes[i] << imm
	}
}

func VSHL_S(imm byte, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(dst.bytes[i:], binary.LittleEndian.Uint32(src.bytes[i:])<<imm)
	}
}

// Table vector Lookup.
// https://developer.arm.com/architectures/instruction-sets/intrinsics/#q=vqtbl4q_u8
// Architectures: A64
func VTBL_B(src *Vector128, table []*Vector128, dst *Vector128) {
	if len(table) > 4 || len(table) < 1 {
		panic("invalid table")
	}
	tmp := Vector128{}
	VDUP_BYTE(0, &tmp)

	switch len(table) {
	case 1:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}

	case 2:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] >= 16 && src.bytes[i] < 32 {
				tmp.bytes[i] = table[1].bytes[src.bytes[i]-16]
			} else if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}
	case 3:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] >= 32 && src.bytes[i] < 48 {
				tmp.bytes[i] = table[2].bytes[src.bytes[i]-32]
			} else if src.bytes[i] >= 16 && src.bytes[i] < 32 {
				tmp.bytes[i] = table[1].bytes[src.bytes[i]-16]
			} else if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}
	case 4:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] >= 48 && src.bytes[i] < 64 {
				tmp.bytes[i] = table[3].bytes[src.bytes[i]-48]
			} else if src.bytes[i] >= 32 && src.bytes[i] < 48 {
				tmp.bytes[i] = table[2].bytes[src.bytes[i]-32]
			} else if src.bytes[i] >= 16 && src.bytes[i] < 32 {
				tmp.bytes[i] = table[1].bytes[src.bytes[i]-16]
			} else if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

// Table vector lookup extension.
// https://developer.arm.com/architectures/instruction-sets/intrinsics/#q=vqtbx4q_u8
// Architectures: A64
func VTBX_B(src *Vector128, table []*Vector128, dst *Vector128) {
	if len(table) > 4 || len(table) < 1 {
		panic("invalid table")
	}
	tmp := Vector128{}
	copy(tmp.bytes[:], dst.bytes[:])

	switch len(table) {
	case 1:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}

	case 2:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] >= 16 && src.bytes[i] < 32 {
				tmp.bytes[i] = table[1].bytes[src.bytes[i]-16]
			} else if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}
	case 3:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] >= 32 && src.bytes[i] < 48 {
				tmp.bytes[i] = table[2].bytes[src.bytes[i]-32]
			} else if src.bytes[i] >= 16 && src.bytes[i] < 32 {
				tmp.bytes[i] = table[1].bytes[src.bytes[i]-16]
			} else if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}
	case 4:
		for i := 0; i < 16; i += 1 {
			if src.bytes[i] >= 48 && src.bytes[i] < 64 {
				tmp.bytes[i] = table[3].bytes[src.bytes[i]-48]
			} else if src.bytes[i] >= 32 && src.bytes[i] < 48 {
				tmp.bytes[i] = table[2].bytes[src.bytes[i]-32]
			} else if src.bytes[i] >= 16 && src.bytes[i] < 32 {
				tmp.bytes[i] = table[1].bytes[src.bytes[i]-16]
			} else if src.bytes[i] < 16 {
				tmp.bytes[i] = table[0].bytes[src.bytes[i]]
			}
		}
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

// Unsigned saturating Subtract.
// https://developer.arm.com/architectures/instruction-sets/intrinsics/vqsubq_u8
func VUQSUB_B(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		if src2.bytes[i] < src1.bytes[i] {
			dst.bytes[i] = 0
		} else {
			dst.bytes[i] = src2.bytes[i] - src1.bytes[i]
		}
	}
}

// Compare unsigned Higher (vector).
// https://developer.arm.com/architectures/instruction-sets/intrinsics/#q=vcgtq_u8
// https://developer.arm.com/documentation/ddi0596/2021-03/SIMD-FP-Instructions/CMHI--register---Compare-unsigned-Higher--vector--?lang=en
// Architectures: v7, A32, A64
func VCMHI_B(eq bool, src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		if eq {
			if src2.bytes[i] >= src1.bytes[i] {
				dst.bytes[i] = 0xff
			} else {
				dst.bytes[i] = 0
			}
		} else {
			if src2.bytes[i] > src1.bytes[i] {
				dst.bytes[i] = 0xff
			} else {
				dst.bytes[i] = 0
			}
		}
	}
}

// Compare unsigned higher or same (vector)
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/CMHS--register---Compare-unsigned-higher-or-same--vector--?lang=en
func VCMHS_B(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		if src2.bytes[i] >= src1.bytes[i] {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

// Compare bitwise equal (vector)
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/CMEQ--register---Compare-bitwise-equal--vector--?lang=en
func VCMEQ_B(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		if src2.bytes[i] == src1.bytes[i] {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

// Compare signed greater than (vector)
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/CMGT--register---Compare-signed-greater-than--vector--?lang=en
func VCMGT_B(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		if int8(src2.bytes[i]) > int8(src1.bytes[i]) {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

// Compare signed greater than or equal (vector)
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/CMGE--register---Compare-signed-greater-than-or-equal--vector--?lang=en
func VCMGE_B(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		if int8(src2.bytes[i]) >= int8(src1.bytes[i]) {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

// Compare bitwise test bits nonzero (vector)
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/CMTST--Compare-bitwise-test-bits-nonzero--vector--?lang=en
func VCMTST_B(Vm, Vn, Vd *Vector128) {
	for i := 0; i < 16; i += 1 {
		if Vm.bytes[i]&Vn.bytes[i] != 0 {
			Vd.bytes[i] = 0xff
		} else {
			Vd.bytes[i] = 0
		}
	}
}

// Unsigned Maximum across Vector.
// https://developer.arm.com/architectures/instruction-sets/intrinsics/#q=vmaxvq_u8
// Architectures: A64
func VUMAXV_B(max bool, src, dst *Vector128) {
	maxmin := src.bytes[0]
	for i := 1; i < 16; i += 1 {
		if max && src.bytes[i] > maxmin {
			maxmin = src.bytes[i]
		} else if !max && src.bytes[i] < maxmin {
			maxmin = src.bytes[i]
		}
	}
	dst.bytes[15] = maxmin
}

// https://developer.arm.com/architectures/instruction-sets/intrinsics/#q=vzip1q_u32
func VZIP1_S(Vm, Vn, dst *Vector128) {
	for i := 0; i < 2; i += 1 {
		a := binary.LittleEndian.Uint32(Vn.bytes[4*i:])
		b := binary.LittleEndian.Uint32(Vm.bytes[4*i:])
		binary.LittleEndian.PutUint32(dst.bytes[8*i:], a)
		binary.LittleEndian.PutUint32(dst.bytes[8*i+4:], b)
	}
}

func VZIP1_D(Vm, Vn, dst *Vector128) {
	a := Vn.Uint64s()
	b := Vm.Uint64s()
	VLD1_2D([]uint64{a[0], b[0]}, dst)
}

func VZIP2_S(Vm, Vn, dst *Vector128) {
	for i := 2; i < 4; i += 1 {
		a := binary.LittleEndian.Uint32(Vn.bytes[4*i:])
		b := binary.LittleEndian.Uint32(Vm.bytes[4*i:])
		binary.LittleEndian.PutUint32(dst.bytes[8*(i-2):], a)
		binary.LittleEndian.PutUint32(dst.bytes[8*(i-2)+4:], b)
	}
}

func VZIP2_D(Vm, Vn, dst *Vector128) {
	a := Vn.Uint64s()
	b := Vm.Uint64s()
	VLD1_2D([]uint64{a[1], b[1]}, dst)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/TRN1--Transpose-vectors--primary--?lang=en
func VTRN1_H(Vm, Vn, dst *Vector128) {
	a := Vn.Uint16s()
	b := Vm.Uint16s()
	VLD1_8H([]uint16{a[0], b[0], a[2], b[2], a[4], b[4], a[6], b[6]}, dst)
}

func VTRN1_S(Vm, Vn, dst *Vector128) {
	a := Vn.Uint32s()
	b := Vm.Uint32s()
	VLD1_4S([]uint32{a[0], b[0], a[2], b[2]}, dst)
}

func VTRN1_D(Vm, Vn, dst *Vector128) {
	a := Vn.Uint64s()
	b := Vm.Uint64s()
	VLD1_2D([]uint64{a[0], b[0]}, dst)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/TRN1--Transpose-vectors--primary--?lang=en
func VTRN2_H(Vm, Vn, dst *Vector128) {
	a := Vn.Uint16s()
	b := Vm.Uint16s()
	VLD1_8H([]uint16{a[1], b[1], a[3], b[3], a[5], b[5], a[7], b[7]}, dst)
}

func VTRN2_S(Vm, Vn, dst *Vector128) {
	a := Vn.Uint32s()
	b := Vm.Uint32s()
	VLD1_4S([]uint32{a[1], b[1], a[3], b[3]}, dst)
}

func VTRN2_D(Vm, Vn, dst *Vector128) {
	a := Vn.Uint64s()
	b := Vm.Uint64s()
	VLD1_2D([]uint64{a[1], b[1]}, dst)
}

// input: from high to low
// t0 = t0.S3, t0.S2, t0.S1, t0.S0
// t1 = t1.S3, t1.S2, t1.S1, t1.S0
// t2 = t2.S3, t2.S2, t2.S1, t2.S0
// t3 = t3.S3, t3.S2, t3.S1, t3.S0
// output: from high to low
// t0 = t3.S0, t2.S0, t1.S0, t0.S0
// t1 = t3.S1, t2.S1, t1.S1, t0.S1
// t2 = t3.S2, t2.S2, t1.S2, t0.S2
// t3 = t3.S3, t2.S3, t1.S3, t0.S3
func PRE_TRANSPOSE_S(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
	)
	VZIP1_S(t1, t0, tmp0)
	VZIP1_S(t3, t2, tmp1)
	VZIP2_S(t1, t0, tmp2)
	VZIP2_S(t3, t2, tmp3)
	VZIP1_D(tmp1, tmp0, t0)
	VZIP2_D(tmp1, tmp0, t1)
	VZIP1_D(tmp3, tmp2, t2)
	VZIP2_D(tmp3, tmp2, t3)
}

// Transpose Matrix with VTRN1/3
func PRE_TRANSPOSE_S2(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
	)
	VTRN1_S(t1, t0, tmp0)
	VTRN1_S(t3, t2, tmp1)
	VTRN2_S(t1, t0, tmp2)
	VTRN2_S(t3, t2, tmp3)
	VTRN1_D(tmp1, tmp0, t0)
	VTRN1_D(tmp3, tmp2, t1)
	VTRN2_D(tmp1, tmp0, t2)
	VTRN2_D(tmp3, tmp2, t3)
}

// input: from high to low
// t0 = t0.S3, t0.S2, t0.S1, t0.S0
// t1 = t1.S3, t1.S2, t1.S1, t1.S0
// t2 = t2.S3, t2.S2, t2.S1, t2.S0
// t3 = t3.S3, t3.S2, t3.S1, t3.S0
// output: from high to low
// t0 = t0.S0, t1.S0, t2.S0, t3.S0
// t1 = t0.S1, t1.S1, t2.S1, t3.S1
// t2 = t0.S2, t1.S2, t2.S2, t3.S2
// t3 = t0.S3, t1.S3, t2.S3, t3.S3
func TRANSPOSE_S(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
	)
	VZIP1_S(t0, t1, tmp0)
	VZIP1_S(t2, t3, tmp1)
	VZIP2_S(t0, t1, tmp2)
	VZIP2_S(t2, t3, tmp3)
	VZIP1_D(tmp0, tmp1, t0)
	VZIP2_D(tmp0, tmp1, t1)
	VZIP1_D(tmp2, tmp3, t2)
	VZIP2_D(tmp2, tmp3, t3)
}

// Transpose Matrix with VTRN1/3
func TRANSPOSE_S2(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
	)
	VTRN1_S(t0, t1, tmp0)
	VTRN1_S(t2, t3, tmp1)
	VTRN2_S(t0, t1, tmp2)
	VTRN2_S(t2, t3, tmp3)
	VTRN1_D(tmp0, tmp1, t0)
	VTRN1_D(tmp2, tmp3, t1)
	VTRN2_D(tmp0, tmp1, t2)
	VTRN2_D(tmp2, tmp3, t3)
}

func clmul(a, b uint64) (hi, lo uint64) {
	var temp uint64
	for i := 0; i < 64; i++ {
		temp = a & (b >> i) & 1
		for j := 1; j < i+1; j++ {
			temp ^= (a >> j) & (b >> (i - j)) & 1
		}
		lo |= temp << i
	}
	for i := 64; i < 127; i++ {
		temp = 0
		for j := i - 63; j < 64; j++ {
			temp ^= (a >> j) & (b >> (i - j)) & 1
		}
		hi |= temp << (i - 64)
	}
	return
}

// Polynomial multiply long
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/PMULL--PMULL2--Polynomial-multiply-long-?lang=en
func VPMULL(Vm, Vn, Vd *Vector128) {
	tmp1 := binary.LittleEndian.Uint64(Vm.bytes[:])
	tmp2 := binary.LittleEndian.Uint64(Vn.bytes[:])
	hi, lo := clmul(tmp1, tmp2)
	binary.LittleEndian.PutUint64(Vd.bytes[:], lo)
	binary.LittleEndian.PutUint64(Vd.bytes[8:], hi)
}

func VPMULL2(Vm, Vn, Vd *Vector128) {
	tmp1 := binary.LittleEndian.Uint64(Vm.bytes[8:])
	tmp2 := binary.LittleEndian.Uint64(Vn.bytes[8:])
	hi, lo := clmul(tmp1, tmp2)
	binary.LittleEndian.PutUint64(Vd.bytes[:], lo)
	binary.LittleEndian.PutUint64(Vd.bytes[8:], hi)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/ADD--vector---Add--vector--?lang=en
func VADD_B(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src1.bytes[i] + src2.bytes[i]
	}
}

func VADD_H(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 2 {
		a := binary.LittleEndian.Uint16(src1.bytes[i:])
		b := binary.LittleEndian.Uint16(src2.bytes[i:])
		binary.LittleEndian.PutUint16(dst.bytes[i:], a+b)
	}
}

func VADD_S(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		a := binary.LittleEndian.Uint32(src1.bytes[i:])
		b := binary.LittleEndian.Uint32(src2.bytes[i:])
		binary.LittleEndian.PutUint32(dst.bytes[i:], a+b)
	}
}

func VADD_D(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 8 {
		a := binary.LittleEndian.Uint64(src1.bytes[i:])
		b := binary.LittleEndian.Uint64(src2.bytes[i:])
		binary.LittleEndian.PutUint64(dst.bytes[i:], a+b)
	}
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SUB--vector---Subtract--vector--?lang=en
func VSUB_B(Vm, Vn, Vd *Vector128) {
	for i := 0; i < 16; i += 1 {
		Vd.bytes[i] = Vn.bytes[i] - Vm.bytes[i]
	}
}

func VSUB_H(Vm, Vn, Vd *Vector128) {
	for i := 0; i < 16; i += 2 {
		a := binary.LittleEndian.Uint16(Vm.bytes[i:])
		b := binary.LittleEndian.Uint16(Vn.bytes[i:])
		binary.LittleEndian.PutUint16(Vd.bytes[i:], b-a)
	}
}

func VSUB_S(Vm, Vn, Vd *Vector128) {
	for i := 0; i < 16; i += 4 {
		a := binary.LittleEndian.Uint32(Vm.bytes[i:])
		b := binary.LittleEndian.Uint32(Vn.bytes[i:])
		binary.LittleEndian.PutUint32(Vd.bytes[i:], b-a)
	}
}

func VSUB_D(Vm, Vn, Vd *Vector128) {
	for i := 0; i < 16; i += 8 {
		a := binary.LittleEndian.Uint64(Vm.bytes[i:])
		b := binary.LittleEndian.Uint64(Vn.bytes[i:])
		binary.LittleEndian.PutUint64(Vd.bytes[i:], b-a)
	}
}

func VMUL_H(Vm, Vn, Vd *Vector128) {
	for i := 0; i < 16; i += 2 {
		a := binary.LittleEndian.Uint16(Vm.bytes[i:])
		b := binary.LittleEndian.Uint16(Vn.bytes[i:])
		binary.LittleEndian.PutUint16(Vd.bytes[i:], a*b)
	}
}

// https://developer.arm.com/documentation/ddi0596/2021-03/SIMD-FP-Instructions/UMULL--UMULL2--vector---Unsigned-Multiply-long--vector--?lang=en
func UMULL_B(Vm, Vn, Vd *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 8; i++ {
		a := Vm.bytes[i]
		b := Vn.bytes[i]
		binary.LittleEndian.PutUint16(tmp.bytes[2*i:], uint16(a)*uint16(b))
	}
	copy(Vd.bytes[:], tmp.bytes[:])
}

func UMULL2_B(Vm, Vn, Vd *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 8; i++ {
		a := Vm.bytes[i+8]
		b := Vn.bytes[i+8]
		binary.LittleEndian.PutUint16(tmp.bytes[2*i:], uint16(a)*uint16(b))
	}
	copy(Vd.bytes[:], tmp.bytes[:])
}

func UMULL_H(Vm, Vn, Vd *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 8; i += 2 {
		a := binary.LittleEndian.Uint16(Vm.bytes[i:])
		b := binary.LittleEndian.Uint16(Vn.bytes[i:])
		binary.LittleEndian.PutUint32(tmp.bytes[2*i:], uint32(a)*uint32(b))
	}
	copy(Vd.bytes[:], tmp.bytes[:])
}

func UMULL2_H(Vm, Vn, Vd *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 8; i += 2 {
		a := binary.LittleEndian.Uint16(Vm.bytes[i+8:])
		b := binary.LittleEndian.Uint16(Vn.bytes[i+8:])
		binary.LittleEndian.PutUint32(tmp.bytes[2*i:], uint32(a)*uint32(b))
	}
	copy(Vd.bytes[:], tmp.bytes[:])
}

// Add pairwise (vector)
// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/ADDP--vector---Add-pairwise--vector--?lang=en
func VADDP_H(Vm, Vn, Vd *Vector128) {
	a := Vn.Uint16s()
	b := Vm.Uint16s()
	VLD1_8H([]uint16{a[0] + a[1], a[2] + a[3], a[4] + a[5], a[6] + a[7], b[0] + b[1], b[2] + b[3], b[4] + b[5], b[6] + b[7]}, Vd)
}

func VADDP_S(Vm, Vn, Vd *Vector128) {
	a := Vn.Uint32s()
	b := Vm.Uint32s()
	VLD1_4S([]uint32{a[0] + a[1], a[2] + a[3], b[0] + b[1], b[2] + b[3]}, Vd)
}
