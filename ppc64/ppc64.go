package ppc64

import "encoding/binary"

type Vector128 struct {
	bytes [16]byte
}

func (m *Vector128) Bytes() []byte {
	return m.bytes[:]
}

func (m *Vector128) Uint64s() []uint64 {
	return []uint64{binary.BigEndian.Uint64(m.bytes[:]), binary.BigEndian.Uint64(m.bytes[8:])}
}

func (m *Vector128) Uint32s() []uint32 {
	return []uint32{binary.BigEndian.Uint32(m.bytes[:]), binary.BigEndian.Uint32(m.bytes[4:]), binary.BigEndian.Uint32(m.bytes[8:]), binary.BigEndian.Uint32(m.bytes[12:])}
}

func LVX(rawbytes []byte, dst *Vector128) {
	copy(dst.bytes[:], rawbytes)
}

func LVX_UINT64(ints []uint64, dst *Vector128) {
	binary.BigEndian.PutUint64(dst.bytes[:], ints[0])
	binary.BigEndian.PutUint64(dst.bytes[8:], ints[1])
}

func LXVD2X(rawbytes []byte, dst *Vector128) {
	copy(dst.bytes[:], rawbytes)
}

func LXVD2X_PPC64LE(rawbytes []byte, dst *Vector128) {
	for i := 0; i < 8; i++ {
		dst.bytes[i] = rawbytes[7-i]
	}
	for i := 8; i < 16; i++ {
		dst.bytes[i] = rawbytes[23-i]
	}
}

func STXVD2X(v *Vector128, dst []byte) {
	copy(dst, v.bytes[:])
}

func STXVD2X_PPC64LE(v *Vector128, dst []byte) {
	for i := 0; i < 8; i++ {
		dst[i] = v.bytes[7-i]
	}
	for i := 8; i < 16; i++ {
		dst[i] = v.bytes[23-i]
	}
}

func LXVD2X_UINT64(ints []uint64, dst *Vector128) {
	binary.BigEndian.PutUint64(dst.bytes[:], ints[0])
	binary.BigEndian.PutUint64(dst.bytes[8:], ints[1])
}

func LXVW4X_UINT32(ints []uint32, dst *Vector128) {
	binary.BigEndian.PutUint32(dst.bytes[:], ints[0])
	binary.BigEndian.PutUint32(dst.bytes[4:], ints[1])
	binary.BigEndian.PutUint32(dst.bytes[8:], ints[2])
	binary.BigEndian.PutUint32(dst.bytes[12:], ints[3])
}

func VAND(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] & src2.bytes[i]
	}
}

func VOR(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] | src2.bytes[i]
	}
}

func VXOR(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] ^ src2.bytes[i]
	}
}

func signExtend32(b byte) uint32 {
	if b&0x10 == 0 {
		return uint32(b & 0x1f)
	} else {
		return uint32(b&0x1f) | 0xffffffe0
	}
}

func VSPLTISW(src byte, dst *Vector128) {
	w := signExtend32(src)
	for i := 0; i < 16; i += 4 {
		binary.BigEndian.PutUint32(dst.bytes[i:], w)
	}
}

func VSPLTW(ind byte, src, dst *Vector128) {
	ind = ind & 0x03
	w := binary.BigEndian.Uint32(src.bytes[ind*4:])
	for i := 0; i < 16; i += 4 {
		binary.BigEndian.PutUint32(dst.bytes[i:], w)
	}
}

func signExtend8(b byte) byte {
	if b&0x10 == 0 {
		return byte(b & 0x1f)
	} else {
		return byte(b&0x1f) | 0xe0
	}
}

func VSPLTISB(b byte, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = signExtend8(b)
	}
}

func VSPLTB(ind byte, src, dst *Vector128) {
	ind = ind & 0x0f
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src.bytes[ind]
	}
}

func VSRB(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i++ {
		ind := indicator.bytes[i] & 0x7
		dst.bytes[i] = src.bytes[i] >> ind
	}
}

func VSRH(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i += 2 {
		ind := indicator.bytes[i+1] & 0x0f
		binary.BigEndian.PutUint16(dst.bytes[i:], binary.BigEndian.Uint16(src.bytes[i:])>>ind)
	}
}

func VSRW(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		ind := indicator.bytes[i+3] & 0x1f
		binary.BigEndian.PutUint32(dst.bytes[i:], binary.BigEndian.Uint32(src.bytes[i:])>>ind)
	}
}

func VSRD(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i += 8 {
		ind := indicator.bytes[i+7] & 0x3f
		binary.BigEndian.PutUint64(dst.bytes[i:], binary.BigEndian.Uint64(src.bytes[i:])>>ind)
	}
}

func VSLB(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i++ {
		ind := indicator.bytes[i] & 0x7
		dst.bytes[i] = src.bytes[i] << ind
	}
}

func VSLH(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i += 2 {
		ind := indicator.bytes[i+1] & 0x0f
		binary.BigEndian.PutUint16(dst.bytes[i:], binary.BigEndian.Uint16(src.bytes[i:])<<ind)
	}
}

func VSLW(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		ind := indicator.bytes[i+3] & 0x1f
		binary.BigEndian.PutUint32(dst.bytes[i:], binary.BigEndian.Uint32(src.bytes[i:])<<ind)
	}
}

func VSL(src, indicator, dst *Vector128) {
	sh := indicator.bytes[15] & 0x03
	tmp := Vector128{}
	for i := 0; i < 16; i++ {
		ind := indicator.bytes[i] & 0x03
		if sh != ind {
			panic("VSL: shift amount must be the same for all bytes")
		}
		if i < 15 {
			tmp.bytes[i] = src.bytes[i]<<sh | src.bytes[i+1]>>(8-sh)
		} else {
			tmp.bytes[i] = src.bytes[i] << sh
		}
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VSR(src, indicator, dst *Vector128) {
	sh := indicator.bytes[15] & 0x03
	tmp := Vector128{}
	for i := 0; i < 16; i++ {
		ind := indicator.bytes[i] & 0x03
		if sh != ind {
			panic("VSR: shift amount must be the same for all bytes")
		}
		if i > 0 {
			tmp.bytes[i] = src.bytes[i]>>sh | src.bytes[i-1]<<(8-sh)
		} else {
			tmp.bytes[i] = src.bytes[i] >> sh
		}
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VRLW(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		ind := indicator.bytes[i+3] & 0x1f
		s := binary.BigEndian.Uint32(src.bytes[i:])
		binary.BigEndian.PutUint32(dst.bytes[i:], (s>>ind)|(s<<(32-ind)))
	}
}

func VSLDOI(shb byte, vA, vB, vD *Vector128) {
	intShb := int(shb) & 0x0f
	tmp := Vector128{}
	for i := intShb; i < 16; i++ {
		tmp.bytes[i-intShb] = vA.bytes[i]
	}
	for i := 0; i < intShb; i++ {
		tmp.bytes[16-intShb+i] = vB.bytes[i]
	}
	copy(vD.bytes[:], tmp.bytes[:])
}

func VSRAB(src, indicator, dst *Vector128) {
	for i := 0; i < 16; i++ {
		ind := indicator.bytes[i] & 0x7
		dst.bytes[i] = byte(int8(src.bytes[i]) >> ind)
	}
}

func VPERM(src1, src2, perm, dst *Vector128) {
	tmp := Vector128{}
	for i := 0; i < 16; i++ {
		idx := perm.bytes[i] & 0x1f
		if idx < 16 {
			tmp.bytes[i] = src1.bytes[idx]
		} else {
			tmp.bytes[i] = src2.bytes[idx-16]
		}
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func XXPERMDI(vA, vB *Vector128, sh byte, dst *Vector128) {
	sh = sh & 0x0f
	tmp := Vector128{}
	d0 := binary.BigEndian.Uint64(vA.bytes[:])
	d1 := binary.BigEndian.Uint64(vA.bytes[8:])
	d2 := binary.BigEndian.Uint64(vB.bytes[:])
	d3 := binary.BigEndian.Uint64(vB.bytes[8:])

	switch sh {
	case 0:
		binary.BigEndian.PutUint64(tmp.bytes[:], d0)
		binary.BigEndian.PutUint64(tmp.bytes[8:], d2)
	case 1:
		binary.BigEndian.PutUint64(tmp.bytes[:], d0)
		binary.BigEndian.PutUint64(tmp.bytes[8:], d3)
	case 2:
		binary.BigEndian.PutUint64(tmp.bytes[:], d1)
		binary.BigEndian.PutUint64(tmp.bytes[8:], d2)
	case 3:
		binary.BigEndian.PutUint64(tmp.bytes[:], d1)
		binary.BigEndian.PutUint64(tmp.bytes[8:], d3)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VCMPGTUB(vA, vB, dst *Vector128) {
	for i := 0; i < 16; i++ {
		if vA.bytes[i] > vB.bytes[i] {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

func VCMPEQUB(vA, vB, dst *Vector128) {
	for i := 0; i < 16; i++ {
		if vA.bytes[i] == vB.bytes[i] {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

func VMRGEW(vA, vB, vD *Vector128) {
	s1 := vA.Uint32s()
	s2 := vB.Uint32s()

	binary.BigEndian.PutUint32(vD.bytes[:], s1[0])
	binary.BigEndian.PutUint32(vD.bytes[4:], s2[0])
	binary.BigEndian.PutUint32(vD.bytes[8:], s1[2])
	binary.BigEndian.PutUint32(vD.bytes[12:], s2[2])
}

func VMRGOW(vA, vB, vD *Vector128) {
	s1 := vA.Uint32s()
	s2 := vB.Uint32s()

	binary.BigEndian.PutUint32(vD.bytes[:], s1[1])
	binary.BigEndian.PutUint32(vD.bytes[4:], s2[1])
	binary.BigEndian.PutUint32(vD.bytes[8:], s1[3])
	binary.BigEndian.PutUint32(vD.bytes[12:], s2[3])
}

func TransposeMatrix1(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
		M0   = &Vector128{}
		M1   = &Vector128{}
		M2   = &Vector128{}
		M3   = &Vector128{}
	)
	LXVD2X_UINT64([]uint64{0x0001020310111213, 0x0405060714151617}, M0)
	LXVD2X_UINT64([]uint64{0x08090a0b18191a1b, 0x0c0d0e0f1c1d1e1f}, M1)
	LXVD2X_UINT64([]uint64{0x0001020304050607, 0x1011121314151617}, M2)
	LXVD2X_UINT64([]uint64{0x08090a0b0c0d0e0f, 0x18191a1b1c1d1e1f}, M3)

	VPERM(t1, t0, M0, tmp0)
	VPERM(t1, t0, M1, tmp1)
	VPERM(t3, t2, M0, tmp2)
	VPERM(t3, t2, M1, tmp3)
	VPERM(tmp2, tmp0, M2, t0)
	VPERM(tmp2, tmp0, M3, t1)
	VPERM(tmp3, tmp1, M2, t2)
	VPERM(tmp3, tmp1, M3, t3)
}

// TransposeMatrix2 transposes a 4x4 matrix with VMRGEW and VMRGOW.
func TransposeMatrix2(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
		M0   = &Vector128{}
		M1   = &Vector128{}
	)
	LXVD2X_UINT64([]uint64{0x0001020304050607, 0x1011121314151617}, M0)
	LXVD2X_UINT64([]uint64{0x08090a0b0c0d0e0f, 0x18191a1b1c1d1e1f}, M1)
	VMRGEW(t1, t0, tmp0)
	VMRGEW(t3, t2, tmp1)
	VMRGOW(t1, t0, tmp2)
	VMRGOW(t3, t2, tmp3)
	VPERM(tmp1, tmp0, M0, t0)
	VPERM(tmp1, tmp0, M1, t2)
	VPERM(tmp3, tmp2, M0, t1)
	VPERM(tmp3, tmp2, M1, t3)
}

// TransposeMatrix3 transposes a 4x4 matrix with VMRGEW/VMRGOW and XXPERMDI.
func TransposeMatrix3(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
	)
	VMRGEW(t1, t0, tmp0)
	VMRGEW(t3, t2, tmp1)
	VMRGOW(t1, t0, tmp2)
	VMRGOW(t3, t2, tmp3)
	XXPERMDI(tmp1, tmp0, 0, t0)
	XXPERMDI(tmp1, tmp0, 3, t2)
	XXPERMDI(tmp3, tmp2, 0, t1)
	XXPERMDI(tmp3, tmp2, 3, t3)
}

func PreTransposeMatrix1(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
		M0   = &Vector128{}
		M1   = &Vector128{}
		M2   = &Vector128{}
		M3   = &Vector128{}
	)
	LXVD2X_UINT64([]uint64{0x0001020310111213, 0x0405060714151617}, M0)
	LXVD2X_UINT64([]uint64{0x08090a0b18191a1b, 0x0c0d0e0f1c1d1e1f}, M1)
	LXVD2X_UINT64([]uint64{0x0001020304050607, 0x1011121314151617}, M2)
	LXVD2X_UINT64([]uint64{0x08090a0b0c0d0e0f, 0x18191a1b1c1d1e1f}, M3)

	VPERM(t0, t1, M0, tmp0)
	VPERM(t0, t1, M1, tmp1)
	VPERM(t2, t3, M0, tmp2)
	VPERM(t2, t3, M1, tmp3)
	VPERM(tmp0, tmp2, M2, t0)
	VPERM(tmp0, tmp2, M3, t1)
	VPERM(tmp1, tmp3, M2, t2)
	VPERM(tmp1, tmp3, M3, t3)
}

// PreTransposeMatrix2 transposes a 4x4 matrix with VMRGEW and VMRGOW.
func PreTransposeMatrix2(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
		M0   = &Vector128{}
		M1   = &Vector128{}
	)
	LXVD2X_UINT64([]uint64{0x0001020304050607, 0x1011121314151617}, M0)
	LXVD2X_UINT64([]uint64{0x08090a0b0c0d0e0f, 0x18191a1b1c1d1e1f}, M1)
	VMRGEW(t0, t1, tmp0)
	VMRGEW(t2, t3, tmp1)
	VMRGOW(t0, t1, tmp2)
	VMRGOW(t2, t3, tmp3)
	VPERM(tmp0, tmp1, M0, t0)
	VPERM(tmp0, tmp1, M1, t2)
	VPERM(tmp2, tmp3, M0, t1)
	VPERM(tmp2, tmp3, M1, t3)
}

// PreTransposeMatrix3 transposes a 4x4 matrix with VMRGEW/VMRGOW and XXPERMDI.
func PreTransposeMatrix3(t0, t1, t2, t3 *Vector128) {
	var (
		tmp0 = &Vector128{}
		tmp1 = &Vector128{}
		tmp2 = &Vector128{}
		tmp3 = &Vector128{}
	)
	VMRGEW(t0, t1, tmp0)
	VMRGEW(t2, t3, tmp1)
	VMRGOW(t0, t1, tmp2)
	VMRGOW(t2, t3, tmp3)
	XXPERMDI(tmp0, tmp1, 0, t0)
	XXPERMDI(tmp0, tmp1, 3, t2)
	XXPERMDI(tmp2, tmp3, 0, t1)
	XXPERMDI(tmp2, tmp3, 3, t3)
}
