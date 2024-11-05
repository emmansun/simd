package avx

import (
	"encoding/binary"

	"github.com/emmansun/simd/amd64/sse"
)

func VPXOR(dst, src1, src2 *sse.XMM) {
	for i := 0; i < 16; i++ {
		dst.Bytes()[i] = src1.Bytes()[i] ^ src2.Bytes()[i]
	}
}

func VMOVDQU(dst *sse.XMM, src *sse.XMM) {
	copy(dst.Bytes(), src.Bytes())
}

func VMOVDQU_L16B(dst *sse.XMM, src []byte) {
	copy(dst.Bytes(), src)
}

func VMOVEDQU_S16B(dst []byte, src *sse.XMM) {
	copy(dst, src.Bytes())
}

func VMOVDQU_L4S(dst *sse.XMM, src []uint32) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(dst.Bytes()[i*4:], src[i])
	}
}

func VMOVDQU_S4S(dst []uint32, src *sse.XMM) {
	for i := 0; i < 4; i++ {
		dst[i] = binary.LittleEndian.Uint32(src.Bytes()[i*4:])
	}
}

func VPSHUFB(dst, src1, src2 *sse.XMM) {
	tmp := sse.XMM{}
	tmpBytes := tmp.Bytes()
	src1Bytes := src1.Bytes()
	maskBytes := src2.Bytes()
	for i := 0; i < 16; i++ {
		if maskBytes[i]&0x80 == 0x80 {
			tmpBytes[i] = 0
		} else {
			idx := maskBytes[i] & 0x0f
			tmpBytes[i] = src1Bytes[idx]
		}
	}
	VMOVDQU(dst, &tmp)
}

func VPSHUFD(dst, src *sse.XMM, imm uint) {
	tmp := &sse.XMM{}
	srcBytes := src.Bytes()
	tmpBytes := tmp.Bytes()
	for i := 0; i < 4; i++ {
		idx := (imm >> (i * 2)) & 0x03
		copy(tmpBytes[i*4:], srcBytes[idx*4:])
	}
	VMOVDQU(dst, tmp)
}

func VPUNPCKLQDQ(dst, src1, src2 *sse.XMM) {
	tmp := &sse.XMM{}
	src1Bytes := src1.Bytes()
	src2Bytes := src2.Bytes()
	tmpBytes := tmp.Bytes()
	for i := 0; i < 8; i++ {
		tmpBytes[i] = src1Bytes[i]
		tmpBytes[i+8] = src2Bytes[i]
	}
	VMOVDQU(dst, tmp)
}

func VPUNPCKHQDQ(dst, src1, src2 *sse.XMM) {
	tmp := &sse.XMM{}
	src1Bytes := src1.Bytes()
	src2Bytes := src2.Bytes()
	tmpBytes := tmp.Bytes()
	for i := 0; i < 8; i++ {
		tmpBytes[i] = src1Bytes[i+8]
		tmpBytes[i+8] = src2Bytes[i+8]
	}
	VMOVDQU(dst, tmp)
}

func VPSRLD(dst, src *sse.XMM, imm uint) {
	e0 := binary.LittleEndian.Uint32(src.Bytes()[:])
	e1 := binary.LittleEndian.Uint32(src.Bytes()[4:])
	e2 := binary.LittleEndian.Uint32(src.Bytes()[8:])
	e3 := binary.LittleEndian.Uint32(src.Bytes()[12:])
	if imm > 31 {
		e0 = 0
		e1 = 0
		e2 = 0
		e3 = 0
	} else {
		e0 = e0 >> imm
		e1 = e1 >> imm
		e2 = e2 >> imm
		e3 = e3 >> imm
	}
	binary.LittleEndian.PutUint32(dst.Bytes()[:], e0)
	binary.LittleEndian.PutUint32(dst.Bytes()[4:], e1)
	binary.LittleEndian.PutUint32(dst.Bytes()[8:], e2)
	binary.LittleEndian.PutUint32(dst.Bytes()[12:], e3)
}

func VPSLLD(dst, src *sse.XMM, imm uint) {
	e0 := binary.LittleEndian.Uint32(src.Bytes()[:])
	e1 := binary.LittleEndian.Uint32(src.Bytes()[4:])
	e2 := binary.LittleEndian.Uint32(src.Bytes()[8:])
	e3 := binary.LittleEndian.Uint32(src.Bytes()[12:])
	if imm > 31 {
		e0 = 0
		e1 = 0
		e2 = 0
		e3 = 0
	} else {
		e0 = e0 << imm
		e1 = e1 << imm
		e2 = e2 << imm
		e3 = e3 << imm
	}
	binary.LittleEndian.PutUint32(dst.Bytes()[:], e0)
	binary.LittleEndian.PutUint32(dst.Bytes()[4:], e1)
	binary.LittleEndian.PutUint32(dst.Bytes()[8:], e2)
	binary.LittleEndian.PutUint32(dst.Bytes()[12:], e3)
}

func VPBLENDD(dst, src1, src2 *sse.XMM, imm uint) {
	src1Words := src1.Uint32s()
	src2Words := src2.Uint32s()
	dstBytes := dst.Bytes()

	for i := 0; i < 4; i++ {
		if (imm & (1 << i)) != 0 {
			binary.LittleEndian.PutUint32(dstBytes[i*4:], src2Words[i])
		} else {
			binary.LittleEndian.PutUint32(dstBytes[i*4:], src1Words[i])
		}
	}
}

func VPALIGNR(dst, src1, src2 *sse.XMM, imm8 byte) {
	tmp := &sse.XMM{}
	src1Bytes := src1.Bytes()
	src2Bytes := src2.Bytes()
	tmpBytes := tmp.Bytes()

	for i := imm8; i < 16; i++ {
		tmpBytes[i-imm8] = src2Bytes[i]
	}
	for i := 0; i < int(imm8); i++ {
		tmpBytes[16-int(imm8)+i] = src1Bytes[i]
	}
	VMOVDQU(dst, tmp)
}

func VPSRLDQ(dst, src *sse.XMM, imm8 byte) {
	tmp := &sse.XMM{}
	srcBytes := src.Bytes()
	tmpBytes := tmp.Bytes()
	if imm8 > 16 {
		imm8 = 16
	}
	for i := 0; i < 16-int(imm8); i++ {
		tmpBytes[i] = srcBytes[i+int(imm8)]
	}
	VMOVDQU(dst, tmp)
}
