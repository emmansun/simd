package s390x

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

func VL(rawbytes []byte, dst *Vector128) {
	copy(dst.bytes[:], rawbytes)
}

func VL_UINT64(ints []uint64, dst *Vector128) {
	binary.BigEndian.PutUint64(dst.bytes[:], ints[0])
	binary.BigEndian.PutUint64(dst.bytes[8:], ints[1])
}

func VL_UINT32(ints []uint32, dst *Vector128) {
	binary.BigEndian.PutUint32(dst.bytes[:], ints[0])
	binary.BigEndian.PutUint32(dst.bytes[4:], ints[1])
	binary.BigEndian.PutUint32(dst.bytes[8:], ints[2])
	binary.BigEndian.PutUint32(dst.bytes[12:], ints[3])
}

func VST(src *Vector128, dst []byte) {
	copy(dst, src.bytes[:])
}

// AND
func VN(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] & src2.bytes[i]
	}
}

// XOR
func VX(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] ^ src2.bytes[i]
	}
}

func VO(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] | src2.bytes[i]
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

func VREPIB(imm uint8, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = imm
	}
}

func VREPIH(imm uint16, dst *Vector128) {
	for i := 0; i < 16; i += 2 {
		binary.BigEndian.PutUint16(dst.bytes[i:], imm)
	}
}

func VREPIF(imm uint32, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		binary.BigEndian.PutUint32(dst.bytes[i:], imm)
	}
}

func VREPIG(imm uint64, dst *Vector128) {
	for i := 0; i < 16; i += 8 {
		binary.BigEndian.PutUint64(dst.bytes[i:], imm)
	}
}

func VREPF(idx uint8, src, dst *Vector128) {
	idx = idx & 0x03
	w := binary.BigEndian.Uint32(src.bytes[idx*4:])
	for i := 0; i < 16; i += 4 {
		binary.BigEndian.PutUint32(dst.bytes[i:], w)
	}
}

func VZERO(dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = 0
	}
}

// Vector Shift Left
func VSL(indicator, src, dst *Vector128) {
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

// Vector Shift Right Logical
func VSRL(indicator, src, dst *Vector128) {
	sh := indicator.bytes[15] & 0x03
	tmp := Vector128{}
	for i := 0; i < 16; i++ {
		ind := indicator.bytes[i] & 0x03
		if sh != ind {
			panic("VSRL: shift amount must be the same for all bytes")
		}
		if i > 0 {
			tmp.bytes[i] = src.bytes[i]>>sh | src.bytes[i-1]<<(8-sh)
		} else {
			tmp.bytes[i] = src.bytes[i] >> sh
		}
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

// Vector Element(Byte) Shift Right Arithmetic
func VESRAB(imm uint8, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	for i := 0; i < 16; i++ {
		tmp := int8(src.bytes[i])
		tmp >>= imm
		dst.bytes[i] = byte(tmp)
	}
}

// Vector Element(Half Word) Shift Right Arithmetic
func VESRAH(imm uint8, src, dst *Vector128) {
	if imm > 15 {
		imm = 15
	}
	for i := 0; i < 8; i++ {
		tmp := int16(binary.BigEndian.Uint16(src.bytes[i*2:]))
		tmp >>= imm
		binary.BigEndian.PutUint16(dst.bytes[i*2:], uint16(tmp))
	}
}

// Vector Element(Word) Shift Right Arithmetic
func VESRAF(imm uint8, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	for i := 0; i < 4; i++ {
		tmp := int32(binary.BigEndian.Uint32(src.bytes[i*4:]))
		tmp >>= imm
		binary.BigEndian.PutUint32(dst.bytes[i*4:], uint32(tmp))
	}
}

// Vector Element(Double Word) Shift Right Arithmetic
func VESRAG(imm uint8, src, dst *Vector128) {
	if imm > 63 {
		imm = 63
	}
	for i := 0; i < 2; i++ {
		tmp := int64(binary.BigEndian.Uint64(src.bytes[i*2:]))
		tmp >>= imm
		binary.BigEndian.PutUint64(dst.bytes[i*4:], uint64(tmp))
	}
}

// Vector Element(Byte) Shift Right Logical
func VESRLB(imm uint8, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src.bytes[i] >> imm
	}
}

// Vector Element(Halfword) Shift Right Logical
func VESRLH(imm uint8, src, dst *Vector128) {
	if imm > 15 {
		imm = 15
	}
	for i := 0; i < 8; i++ {
		tmp := binary.BigEndian.Uint16(src.bytes[i*2:])
		tmp >>= imm
		binary.BigEndian.PutUint16(dst.bytes[i*2:], tmp)
	}
}

// Vector Element(Word) Shift Right Logical
func VESRLF(imm uint8, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	for i := 0; i < 4; i++ {
		tmp := binary.BigEndian.Uint32(src.bytes[i*4:])
		tmp >>= imm
		binary.BigEndian.PutUint32(dst.bytes[i*4:], tmp)
	}
}

// Vector Element(Double Word) Shift Right Logical
func VESRLG(imm uint8, src, dst *Vector128) {
	if imm > 63 {
		imm = 63
	}
	for i := 0; i < 2; i++ {
		tmp := binary.BigEndian.Uint64(src.bytes[i*8:])
		tmp >>= imm
		binary.BigEndian.PutUint64(dst.bytes[i*8:], tmp)
	}
}

// Vector Element(Byte) Shift Left
func VESLB(imm uint8, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src.bytes[i] << imm
	}
}

// Vector Element(Halfword) Shift Left
func VESLH(imm uint8, src, dst *Vector128) {
	if imm > 15 {
		imm = 15
	}
	for i := 0; i < 8; i++ {
		tmp := binary.BigEndian.Uint16(src.bytes[i*2:])
		tmp <<= imm
		binary.BigEndian.PutUint16(dst.bytes[i*2:], tmp)
	}
}

// Vector Element(Word) Shift Left
func VESLF(imm uint8, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	for i := 0; i < 4; i++ {
		tmp := binary.BigEndian.Uint32(src.bytes[i*4:])
		tmp <<= imm
		binary.BigEndian.PutUint32(dst.bytes[i*4:], tmp)
	}
}

// Vector Element(Double Word) Shift Left
func VESLG(imm uint8, src, dst *Vector128) {
	if imm > 63 {
		imm = 63
	}
	for i := 0; i < 2; i++ {
		tmp := binary.BigEndian.Uint64(src.bytes[i*8:])
		tmp <<= imm
		binary.BigEndian.PutUint64(dst.bytes[i*8:], tmp)
	}
}

// Vector Element(Byte) Rotate Shift Left Logical
func VERLLB(imm uint8, src, dst *Vector128) {
	if imm > 7 {
		imm = 7
	}
	for i := 0; i < 16; i++ {
		tmp := src.bytes[i]
		tmp = tmp>>(8-imm) | tmp<<imm
		dst.bytes[i] = tmp
	}
}

// Vector Element(Half Word) Rotate Shift Left Logical
func VERLLH(imm uint8, src, dst *Vector128) {
	if imm > 15 {
		imm = 15
	}
	for i := 0; i < 8; i++ {
		tmp := binary.BigEndian.Uint16(src.bytes[i*2:])
		tmp = tmp>>(16-imm) | tmp<<imm
		binary.BigEndian.PutUint16(dst.bytes[i*2:], tmp)
	}
}

// Vector Element(Word) Rotate Shift Left Logical
func VERLLF(imm uint8, src, dst *Vector128) {
	if imm > 31 {
		imm = 31
	}
	for i := 0; i < 4; i++ {
		tmp := binary.BigEndian.Uint32(src.bytes[i*4:])
		tmp = tmp>>(32-imm) | tmp<<imm
		binary.BigEndian.PutUint32(dst.bytes[i*4:], tmp)
	}
}

// Vector Element(Double Word) Rotate Shift Left Logical
func VERLLG(imm uint8, src, dst *Vector128) {
	if imm > 63 {
		imm = 63
	}
	for i := 0; i < 2; i++ {
		tmp := binary.BigEndian.Uint64(src.bytes[i*8:])
		tmp = tmp>>(64-imm) | tmp<<imm
		binary.BigEndian.PutUint64(dst.bytes[i*8:], tmp)
	}
}

// Vector Load Element Immediate (Byte)
func VLEIB(idx, value uint8, dst *Vector128) {
	if idx > 15 {
		idx = 15
	}
	dst.bytes[idx] = value
}

// Vector Load Element Immediate (Halfword)
func VLEIH(idx uint8, value uint16, dst *Vector128) {
	if idx > 7 {
		idx = 7
	}
	binary.BigEndian.PutUint16(dst.bytes[idx:], value)
}

// Vector Load Element Immediate (Word)
func VLEIF(idx uint8, value uint32, dst *Vector128) {
	if idx > 3 {
		idx = 3
	}
	binary.BigEndian.PutUint32(dst.bytes[idx:], value)
}

// Vector Load Element Immediate (Doubleword)
func VLEIG(idx uint8, value uint64, dst *Vector128) {
	if idx > 1 {
		idx = 1
	}
	binary.BigEndian.PutUint64(dst.bytes[idx:], value)
}

// Vector Subsctraction
func VSB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src2.bytes[i] - src1.bytes[i]
	}
}

// Vector Maximum
func VMXB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		if int8(src2.bytes[i]) > int8(src1.bytes[i]) {
			dst.bytes[i] = src2.bytes[i]
		} else {
			dst.bytes[i] = src1.bytes[i]
		}
	}
}

// Vector Maximum Logical
func VMXLB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		if src2.bytes[i] > src1.bytes[i] {
			dst.bytes[i] = src2.bytes[i]
		} else {
			dst.bytes[i] = src1.bytes[i]
		}
	}
}

func VMLHH(src1, src2, dst *Vector128) {
	for i := 0; i < 8; i++ {
		e0 := binary.BigEndian.Uint16(src1.bytes[i*2:])
		e1 := binary.BigEndian.Uint16(src2.bytes[i*2:])
		binary.BigEndian.PutUint16(dst.bytes[i*2:], uint16((uint32(e0)*uint32(e1))>>16))
	}
}

func VMLHW(src1, src2, dst *Vector128) {
	for i := 0; i < 8; i++ {
		e0 := binary.BigEndian.Uint16(src1.bytes[i*2:])
		e1 := binary.BigEndian.Uint16(src2.bytes[i*2:])
		binary.BigEndian.PutUint16(dst.bytes[i*2:], uint16(uint32(e0)*uint32(e1)))
	}
}

func VECB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		a := int8(src1.bytes[i])
		b := int8(src2.bytes[i])
		if b > a {
			dst.bytes[i] = 1
		} else if a == b {
			dst.bytes[i] = 0
		} else {
			dst.bytes[i] = 0xff
		}
	}
}

func VECLB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		a := src1.bytes[i]
		b := src2.bytes[i]
		if b > a {
			dst.bytes[i] = 1
		} else if a == b {
			dst.bytes[i] = 0
		} else {
			dst.bytes[i] = 0xff
		}
	}
}

func VAB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		a := src1.bytes[i]
		b := src2.bytes[i]
		dst.bytes[i] = byte(a + b)
	}
}

func VAH(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 2 {
		a := binary.BigEndian.Uint16(src1.bytes[i:])
		b := binary.BigEndian.Uint16(src2.bytes[i:])
		binary.BigEndian.PutUint16(dst.bytes[i:], a+b)
	}
}

func VAF(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		a := binary.BigEndian.Uint32(src1.bytes[i:])
		b := binary.BigEndian.Uint32(src2.bytes[i:])
		binary.BigEndian.PutUint32(dst.bytes[i:], a+b)
	}
}

func VCEQB(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		if src1.bytes[i] == src2.bytes[i] {
			dst.bytes[i] = 0xff
		} else {
			dst.bytes[i] = 0
		}
	}
}

func VMLOB(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 2 {
		a := uint16(src1.bytes[i+1])
		b := uint16(src2.bytes[i+1])
		binary.BigEndian.PutUint16(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMLEB(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 2 {
		a := uint16(src1.bytes[i])
		b := uint16(src2.bytes[i])
		binary.BigEndian.PutUint16(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMLOH(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 4 {
		a := uint32(binary.BigEndian.Uint16(src1.bytes[i+2:]))
		b := uint32(binary.BigEndian.Uint16(src2.bytes[i+2:]))
		binary.BigEndian.PutUint32(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMLEH(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 4 {
		a := uint32(binary.BigEndian.Uint16(src1.bytes[i:]))
		b := uint32(binary.BigEndian.Uint16(src2.bytes[i:]))
		binary.BigEndian.PutUint32(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}
