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

func VLD1_16B(rawbytes []byte, dst *Vector128) {
	copy(dst.bytes[:], rawbytes)
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

func VST1_16B(src *Vector128, dst []byte) {
	copy(dst, src.bytes[:])
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

func VSHL_B(imm byte, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	for i := 0; i < 16; i += 1 {
		dst.bytes[i] = src.bytes[i] << imm
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
