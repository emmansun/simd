package sse

import "testing"

var encodeStdLut = [16]byte{65, 71, 252, 252, 252, 252, 252, 252, 252, 252, 252, 252, 237, 240, 0, 0}
var reshuffle_mask = Set64(0x0a0b090a07080607, 0x0405030401020001)
var mulhi_mask = Set64(0x0FC0FC000FC0FC00, 0x0FC0FC000FC0FC00)
var mulhi_const = Set64(0x0400004004000040, 0x0400004004000040)
var mullo_mask = Set64(0x003F03F0003F03F0, 0x003F03F0003F03F0)
var mullo_const = Set64(0x0100001001000010, 0x0100001001000010)
var range_1_end = Set64(0x3333333333333333, 0x3333333333333333)
var range_0_end = Set64(0x1919191919191919, 0x1919191919191919)

func encode(src []byte) (dst []byte) {
	x0 := &XMM{}
	SetBytes(x0, src)
	x1 := &XMM{}
	SetBytes(x1, encodeStdLut[:])
	x2 := &XMM{}
	x3 := &XMM{}

	PSHUFB(x0, &reshuffle_mask)
	MOVOU(x2, x0)
	PAND(x2, &mulhi_mask)
	PMULHUW(x2, &mulhi_const)
	PAND(x0, &mullo_mask)
	PMULLW(x0, &mullo_const)
	POR(x0, x2)
	MOVOU(x2, x0)
	MOVOU(x3, x0)
	PSUBUSB(x3, &range_1_end)
	PCMPGTB(x2, &range_0_end)
	PSUBB(x3, x2)
	MOVOU(x2, x1)
	PSHUFB(x2, x3)
	PADDB(x0, x2)
	dst = make([]byte, 16)
	copy(dst, x0.Bytes())
	return
}

func TestEncode(t *testing.T) {
	ret := encode([]byte("abcdefghijkl0000"))
	if string(ret) != "YWJjZGVmZ2hpamts" {
		t.Errorf("encode() = %v; want YWJjZGVmZ2hpamts", string(ret))
	}
}
