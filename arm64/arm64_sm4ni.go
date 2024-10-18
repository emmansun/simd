package arm64

import "encoding/binary"

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
