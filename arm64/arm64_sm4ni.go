package arm64

import "encoding/binary"

var sm4_ck = [32]uint32{
	0x00070e15, 0x1c232a31, 0x383f464d, 0x545b6269, 0x70777e85, 0x8c939aa1, 0xa8afb6bd, 0xc4cbd2d9,
	0xe0e7eef5, 0xfc030a11, 0x181f262d, 0x343b4249, 0x50575e65, 0x6c737a81, 0x888f969d, 0xa4abb2b9,
	0xc0c7ced5, 0xdce3eaf1, 0xf8ff060d, 0x141b2229, 0x30373e45, 0x4c535a61, 0x686f767d, 0x848b9299,
	0xa0a7aeb5, 0xbcc3cad1, 0xd8dfe6ed, 0xf4fb0209, 0x10171e25, 0x2c333a41, 0x484f565d, 0x646b7279,
}

// T
func t(in uint32) uint32 {
	var b uint32

	b = uint32(sm4_sbox[in&0xff])
	b |= uint32(sm4_sbox[in>>8&0xff]) << 8
	b |= uint32(sm4_sbox[in>>16&0xff]) << 16
	b |= uint32(sm4_sbox[in>>24&0xff]) << 24

	// L
	return b ^ (b<<2 | b>>30) ^ (b<<10 | b>>22) ^ (b<<18 | b>>14) ^ (b<<24 | b>>8)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM4E--SM4-encode-?lang=en
// SM4E handle 4 round keys
// Vn: round key
// Vd: in out data
func SM4E(Vn, Vd *Vector128) {
	roundresult := &Vector128{}
	copy(roundresult.bytes[:], Vd.bytes[:])
	for i := 0; i < 16; i += 4 {
		rk := binary.LittleEndian.Uint32(Vn.bytes[i:])
		b0 := binary.LittleEndian.Uint32(roundresult.bytes[:])
		b1 := binary.LittleEndian.Uint32(roundresult.bytes[4:])
		b2 := binary.LittleEndian.Uint32(roundresult.bytes[8:])
		b3 := binary.LittleEndian.Uint32(roundresult.bytes[12:])
		intval := b0 ^ t(b1^b2^b3^rk)
		binary.LittleEndian.PutUint32(roundresult.bytes[:], b1)
		binary.LittleEndian.PutUint32(roundresult.bytes[4:], b2)
		binary.LittleEndian.PutUint32(roundresult.bytes[8:], b3)
		binary.LittleEndian.PutUint32(roundresult.bytes[12:], intval)
	}
	copy(Vd.bytes[:], roundresult.bytes[:])
}

func t2(in uint32) uint32 {
	var b uint32

	b = uint32(sm4_sbox[in&0xff])
	b |= uint32(sm4_sbox[in>>8&0xff]) << 8
	b |= uint32(sm4_sbox[in>>16&0xff]) << 16
	b |= uint32(sm4_sbox[in>>24&0xff]) << 24

	// L2
	return b ^ (b<<13 | b>>19) ^ (b<<23 | b>>9)
}

// https://developer.arm.com/documentation/ddi0602/2024-09/SIMD-FP-Instructions/SM4EKEY--SM4-key-?lang=en
// Vm: ck
// Vn: key
// Vd: result
func SM4EKEY(Vm, Vn, Vd *Vector128) {
	roundresult := &Vector128{}
	copy(roundresult.bytes[:], Vn.bytes[:])
	for i := 0; i < 16; i += 4 {
		ck := binary.LittleEndian.Uint32(Vm.bytes[i:])
		b0 := binary.LittleEndian.Uint32(roundresult.bytes[:])
		b1 := binary.LittleEndian.Uint32(roundresult.bytes[4:])
		b2 := binary.LittleEndian.Uint32(roundresult.bytes[8:])
		b3 := binary.LittleEndian.Uint32(roundresult.bytes[12:])
		intval := b0 ^ t2(b1^b2^b3^ck)
		binary.LittleEndian.PutUint32(roundresult.bytes[:], b1)
		binary.LittleEndian.PutUint32(roundresult.bytes[4:], b2)
		binary.LittleEndian.PutUint32(roundresult.bytes[8:], b3)
		binary.LittleEndian.PutUint32(roundresult.bytes[12:], intval)
	}
	copy(Vd.bytes[:], roundresult.bytes[:])
}
