package arm64

import (
	"encoding/binary"

	"github.com/emmansun/simd/alg/sm4"
)

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
		intval := b0 ^ sm4.T(b1^b2^b3^rk)
		binary.LittleEndian.PutUint32(roundresult.bytes[:], b1)
		binary.LittleEndian.PutUint32(roundresult.bytes[4:], b2)
		binary.LittleEndian.PutUint32(roundresult.bytes[8:], b3)
		binary.LittleEndian.PutUint32(roundresult.bytes[12:], intval)
	}
	copy(Vd.bytes[:], roundresult.bytes[:])
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
		intval := b0 ^ sm4.T_KEY(b1^b2^b3^ck)
		binary.LittleEndian.PutUint32(roundresult.bytes[:], b1)
		binary.LittleEndian.PutUint32(roundresult.bytes[4:], b2)
		binary.LittleEndian.PutUint32(roundresult.bytes[8:], b3)
		binary.LittleEndian.PutUint32(roundresult.bytes[12:], intval)
	}
	copy(Vd.bytes[:], roundresult.bytes[:])
}
