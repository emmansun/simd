package ppc64

import "encoding/binary"

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

func VPMSUMD(src1, src2, dst *Vector128) {
	hi1 := binary.BigEndian.Uint64(src1.bytes[:])
	lo1 := binary.BigEndian.Uint64(src1.bytes[8:])
	hi2 := binary.BigEndian.Uint64(src2.bytes[:])
	lo2 := binary.BigEndian.Uint64(src2.bytes[8:])

	hi, lo := clmul(hi1, hi2)
	hi3, lo3 := clmul(lo1, lo2)
	hi ^= hi3
	lo ^= lo3
	binary.BigEndian.PutUint64(dst.bytes[:], hi)
	binary.BigEndian.PutUint64(dst.bytes[8:], lo)
}

func VMULOUB(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 2 {
		a := uint16(src1.bytes[i+1])
		b := uint16(src2.bytes[i+1])
		binary.BigEndian.PutUint16(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULEUB(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 2 {
		a := uint16(src1.bytes[i])
		b := uint16(src2.bytes[i])
		binary.BigEndian.PutUint16(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULOSB(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 2 {
		a := int16(int8(src1.bytes[i+1]))
		b := int16(int8(src2.bytes[i+1]))
		binary.BigEndian.PutUint16(tmp.bytes[i:], uint16(a*b))
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULESB(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 2 {
		a := int16(int8(src1.bytes[i]))
		b := int16(int8(src2.bytes[i]))
		binary.BigEndian.PutUint16(tmp.bytes[i:], uint16(a*b))
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULOUH(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 4 {
		a := uint32(binary.BigEndian.Uint16(src1.bytes[i+2:]))
		b := uint32(binary.BigEndian.Uint16(src2.bytes[i+2:]))
		binary.BigEndian.PutUint32(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULEUH(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 4 {
		a := uint32(binary.BigEndian.Uint16(src1.bytes[i:]))
		b := uint32(binary.BigEndian.Uint16(src2.bytes[i:]))
		binary.BigEndian.PutUint32(tmp.bytes[i:], a*b)
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULOSH(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 4 {
		a := int32(int16(binary.BigEndian.Uint16(src1.bytes[i+2:])))
		b := int32(int16(binary.BigEndian.Uint16(src2.bytes[i+2:])))
		binary.BigEndian.PutUint32(tmp.bytes[i:], uint32(a*b))
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VMULESH(src1, src2, dst *Vector128) {
	tmp := &Vector128{}
	for i := 0; i < 16; i += 4 {
		a := int32(int16(binary.BigEndian.Uint16(src1.bytes[i:])))
		b := int32(int16(binary.BigEndian.Uint16(src2.bytes[i:])))
		binary.BigEndian.PutUint32(tmp.bytes[i:], uint32(a*b))
	}
	copy(dst.bytes[:], tmp.bytes[:])
}

func VADDUBM(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = src1.bytes[i] + src2.bytes[i]
	}
}

func VADDUHM(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 2 {
		binary.BigEndian.PutUint16(dst.bytes[i:], binary.BigEndian.Uint16(src1.bytes[i:])+binary.BigEndian.Uint16(src2.bytes[i:]))
	}
}

func VADDUWM(src1, src2, dst *Vector128) {
	for i := 0; i < 16; i += 4 {
		binary.BigEndian.PutUint32(dst.bytes[i:], binary.BigEndian.Uint32(src1.bytes[i:])+binary.BigEndian.Uint32(src2.bytes[i:]))
	}
}

func VSUBUBM(vB, vA, dst *Vector128) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = vA.bytes[i] - vB.bytes[i]
	}
}

func VSUBUBS(vB, vA, dst *Vector128) {
	for i := 0; i < 16; i++ {
		if vA.bytes[i] < vB.bytes[i] {
			dst.bytes[i] = 0
		} else {
			dst.bytes[i] = vA.bytes[i] - vB.bytes[i]
		}
	}
}
