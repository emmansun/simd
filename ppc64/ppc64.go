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
