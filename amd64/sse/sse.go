package sse

import (
	"encoding/binary"
	"math"
)

type XMM struct {
	bytes [16]byte
}

func (m *XMM) Bytes() []byte {
	return m.bytes[:]
}

func (m *XMM) Uint64s() []uint64 {
	return []uint64{binary.LittleEndian.Uint64(m.bytes[:]), binary.LittleEndian.Uint64(m.bytes[8:])}
}

func (m *XMM) Uint32s() []uint32 {
	return []uint32{binary.LittleEndian.Uint32(m.bytes[:]), binary.LittleEndian.Uint32(m.bytes[4:]), binary.LittleEndian.Uint32(m.bytes[8:]), binary.LittleEndian.Uint32(m.bytes[12:])}
}

func (m *XMM) Uint16s() []uint16 {
	return []uint16{binary.LittleEndian.Uint16(m.bytes[:]), binary.LittleEndian.Uint16(m.bytes[2:]), binary.LittleEndian.Uint16(m.bytes[4:]), binary.LittleEndian.Uint16(m.bytes[6:]), binary.LittleEndian.Uint16(m.bytes[8:]), binary.LittleEndian.Uint16(m.bytes[10:]), binary.LittleEndian.Uint16(m.bytes[12:]), binary.LittleEndian.Uint16(m.bytes[14:])}
}

func MOVOU_U64(dst *XMM, hi, lo uint64) {
	binary.LittleEndian.PutUint64(dst.bytes[:], lo)
	binary.LittleEndian.PutUint64(dst.bytes[8:], hi)
}

func MOVOU(dst, src *XMM) {
	copy(dst.bytes[:], src.bytes[:])
}

func SetBytes(dst *XMM, b []byte) {
	copy(dst.bytes[:], b)
}

func Set64(hi, lo uint64) (m XMM) {
	binary.LittleEndian.PutUint64(m.bytes[:], lo)
	binary.LittleEndian.PutUint64(m.bytes[8:], hi)
	return
}

func SetEpi32(e0, e1, e2, e3 uint32) (m XMM) {
	binary.LittleEndian.PutUint32(m.bytes[:], e0)
	binary.LittleEndian.PutUint32(m.bytes[4:], e1)
	binary.LittleEndian.PutUint32(m.bytes[8:], e2)
	binary.LittleEndian.PutUint32(m.bytes[12:], e3)
	return
}

func mm_and_si128(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = dst.bytes[i] & src.bytes[i]
	}
}

func PAND(dst, src *XMM) {
	mm_and_si128(dst, src)
}

func mm_or_si128(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = dst.bytes[i] | src.bytes[i]
	}
}

func POR(dst, src *XMM) {
	mm_or_si128(dst, src)
}

func mm_xor_si128(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = dst.bytes[i] ^ src.bytes[i]
	}
}

func PXOR(dst, src *XMM) {
	mm_xor_si128(dst, src)
}

func mm_andnot_si128(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		dst.bytes[i] = (^dst.bytes[i]) & src.bytes[i]
	}
}

func PANDN(dst, src *XMM) {
	mm_andnot_si128(dst, src)
}

func mm_shuffle_epi8(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 16; i++ {
		if src.bytes[i]&0x80 == 0x80 {
			tmp.bytes[i] = 0
		} else {
			idx := src.bytes[i] & 0x0f
			tmp.bytes[i] = dst.bytes[idx]
		}
	}
	MOVOU(dst, &tmp)
}

func PSHUFB(dst, src *XMM) {
	mm_shuffle_epi8(dst, src)
}

func mm_srli_epi32(dst *XMM, imm uint) {
	e0 := binary.LittleEndian.Uint32(dst.bytes[:])
	e1 := binary.LittleEndian.Uint32(dst.bytes[4:])
	e2 := binary.LittleEndian.Uint32(dst.bytes[8:])
	e3 := binary.LittleEndian.Uint32(dst.bytes[12:])
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
	binary.LittleEndian.PutUint32(dst.bytes[:], e0)
	binary.LittleEndian.PutUint32(dst.bytes[4:], e1)
	binary.LittleEndian.PutUint32(dst.bytes[8:], e2)
	binary.LittleEndian.PutUint32(dst.bytes[12:], e3)
}

func PSRLW(dst *XMM, imm uint) {
	mm_srli_epi32(dst, imm)
}

func mm_slli_epi32(dst *XMM, imm uint) {
	e0 := binary.LittleEndian.Uint32(dst.bytes[:])
	e1 := binary.LittleEndian.Uint32(dst.bytes[4:])
	e2 := binary.LittleEndian.Uint32(dst.bytes[8:])
	e3 := binary.LittleEndian.Uint32(dst.bytes[12:])
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
	binary.LittleEndian.PutUint32(dst.bytes[:], e0)
	binary.LittleEndian.PutUint32(dst.bytes[4:], e1)
	binary.LittleEndian.PutUint32(dst.bytes[8:], e2)
	binary.LittleEndian.PutUint32(dst.bytes[12:], e3)
}

func PSLLW(dst *XMM, imm uint) {
	mm_slli_epi32(dst, imm)
}

func mm_srli_epi64(dst *XMM, imm uint) {
	lo := binary.LittleEndian.Uint64(dst.bytes[:])
	hi := binary.LittleEndian.Uint64(dst.bytes[8:])
	if imm > 63 {
		lo = 0
		hi = 0
	} else {
		lo = lo >> imm
		hi = hi >> imm
	}
	binary.LittleEndian.PutUint64(dst.bytes[:], lo)
	binary.LittleEndian.PutUint64(dst.bytes[8:], hi)
}

func PSRLD(dst *XMM, imm uint) {
	mm_srli_epi64(dst, imm)
}

func PSRLQ(dst *XMM, imm uint) {
	mm_srli_epi64(dst, imm)
}

func mm_slli_epi64(dst *XMM, imm uint) {
	lo := binary.LittleEndian.Uint64(dst.bytes[:])
	hi := binary.LittleEndian.Uint64(dst.bytes[8:])
	if imm > 63 {
		lo = 0
		hi = 0
	} else {
		lo = lo << imm
		hi = hi << imm
	}
	binary.LittleEndian.PutUint64(dst.bytes[:], lo)
	binary.LittleEndian.PutUint64(dst.bytes[8:], hi)
}

func PSLLD(dst *XMM, imm uint) {
	mm_slli_epi64(dst, imm)
}

func PSLLQ(dst *XMM, imm uint) {
	mm_slli_epi64(dst, imm)
}

func PSHUFD(dst, src *XMM, imm uint) {
	tmp := XMM{}
	for i := 0; i < 4; i++ {
		idx := (imm >> (i * 2)) & 0x03
		copy(tmp.bytes[i*4:], src.bytes[idx*4:])
	}
	MOVOU(dst, &tmp)
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

func PCLMULQDQ(dst, src *XMM, imm uint) {
	var tmp1, tmp2 uint64
	if imm&1 == 0 {
		tmp1 = binary.LittleEndian.Uint64(dst.bytes[:])
	} else {
		tmp1 = binary.LittleEndian.Uint64(dst.bytes[8:])
	}
	if (imm>>4)&1 == 0 {
		tmp2 = binary.LittleEndian.Uint64(src.bytes[:])
	} else {
		tmp2 = binary.LittleEndian.Uint64(src.bytes[8:])
	}

	hi, lo := clmul(tmp1, tmp2)
	binary.LittleEndian.PutUint64(dst.bytes[:], lo)
	binary.LittleEndian.PutUint64(dst.bytes[8:], hi)
}

func PSRLDQ(dst *XMM, imm uint) {
	tmp := XMM{}
	if imm > 15 {
		for i := 0; i < 16; i++ {
			tmp.bytes[i] = 0
		}
	} else {
		for i := 0; i < 16-int(imm); i++ {
			tmp.bytes[i] = dst.bytes[i+int(imm)]
		}
		for i := 16 - int(imm); i < 16; i++ {
			tmp.bytes[i] = 0
		}
	}
	MOVOU(dst, &tmp)
}

func PSRAW(dst *XMM, imm byte) {
	for i := 0; i < 4; i++ {
		w := int32(binary.LittleEndian.Uint32(dst.bytes[i*4:]))
		if imm > 31 {
			imm = 31
		}
		binary.LittleEndian.PutUint32(dst.bytes[i*4:], uint32(w>>imm))
	}
}

func PSLLDQ(dst *XMM, imm byte) {
	tmp := XMM{}
	if imm > 15 {
		for i := 0; i < 16; i++ {
			tmp.bytes[i] = 0
		}
	} else {
		for i := 15; i >= int(imm); i-- {
			tmp.bytes[i] = dst.bytes[i-int(imm)]
		}
		for i := 0; i < int(imm); i++ {
			tmp.bytes[i] = 0
		}
	}
	MOVOU(dst, &tmp)
}

func PMULHUW(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 8; i++ {
		e0 := binary.LittleEndian.Uint16(dst.bytes[i*2:])
		e1 := binary.LittleEndian.Uint16(src.bytes[i*2:])
		binary.LittleEndian.PutUint16(tmp.bytes[i*2:], uint16((uint32(e0)*uint32(e1))>>16))
	}
	MOVOU(dst, &tmp)
}

func PMULLW(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 8; i++ {
		e0 := int16(binary.LittleEndian.Uint16(dst.bytes[i*2:]))
		e1 := int16(binary.LittleEndian.Uint16(src.bytes[i*2:]))
		binary.LittleEndian.PutUint16(tmp.bytes[i*2:], uint16(e0*e1))
	}
	MOVOU(dst, &tmp)
}

func PSUBUSB(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		a := dst.bytes[i]
		b := src.bytes[i]
		// SaturateToUnsignedByte
		if a < b {
			dst.bytes[i] = 0
		} else {
			dst.bytes[i] = a - b
		}
	}
}

func PSUBB(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		a := dst.bytes[i]
		b := src.bytes[i]
		dst.bytes[i] = a - b
	}
}

func PADDB(dst, src *XMM) {
	for i := 0; i < 16; i++ {
		a := dst.bytes[i]
		b := src.bytes[i]
		dst.bytes[i] = a + b
	}
}

func PCMPGTB(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 16; i++ {
		if dst.bytes[i] > src.bytes[i] {
			tmp.bytes[i] = 0xff
		} else {
			tmp.bytes[i] = 0
		}
	}
	MOVOU(dst, &tmp)
}

func PCMPEQB(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 16; i++ {
		if dst.bytes[i] == src.bytes[i] {
			tmp.bytes[i] = 0xff
		} else {
			tmp.bytes[i] = 0
		}
	}
	MOVOU(dst, &tmp)
}

func PMOVMSKB(src *XMM) uint64 {
	var ret uint64
	for i := 0; i < 16; i++ {
		ret |= uint64(src.bytes[i]>>7) << i
	}
	return ret
}

func SaturateAdd16(x, y int16) int16 {
	r := int32(x) + int32(y)
	if r > math.MaxInt16 {
		return math.MaxInt16
	} else if r < math.MinInt16 {
		return math.MinInt16
	}
	return int16(r)
}

func PMADDUBSW(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 8; i++ {
		// SaturateToSignedWord
		ret := SaturateAdd16(int16(dst.bytes[i*2])*int16(int8(src.bytes[i*2])), int16(dst.bytes[i*2+1])*int16(int8(src.bytes[i*2+1])))
		binary.LittleEndian.PutUint16(tmp.bytes[i*2:], uint16(ret))
	}
	MOVOU(dst, &tmp)
}

func PMADDWD(dst, src *XMM) {
	tmp := XMM{}
	for i := 0; i < 4; i++ {
		ret := int32(int16(binary.LittleEndian.Uint16(dst.bytes[i*4:]))) * int32(int16(binary.LittleEndian.Uint16(src.bytes[i*4:])))
		ret += int32(int16(binary.LittleEndian.Uint16(dst.bytes[i*4+2:]))) * int32(int16(binary.LittleEndian.Uint16(src.bytes[i*4+2:])))
		binary.LittleEndian.PutUint32(tmp.bytes[i*4:], uint32(ret))
	}
	MOVOU(dst, &tmp)
}
