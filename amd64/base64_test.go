package amd64

import (
	"testing"

	"github.com/emmansun/simd/amd64/sse"
)

var encodeStdLut = [16]byte{65, 71, 252, 252, 252, 252, 252, 252, 252, 252, 252, 252, 237, 240, 0, 0}
var reshuffle_mask = sse.Set64(0x0a0b090a07080607, 0x0405030401020001)
var mulhi_mask = sse.Set64(0x0FC0FC000FC0FC00, 0x0FC0FC000FC0FC00)
var mulhi_const = sse.Set64(0x0400004004000040, 0x0400004004000040)
var mullo_mask = sse.Set64(0x003F03F0003F03F0, 0x003F03F0003F03F0)
var mullo_const = sse.Set64(0x0100001001000010, 0x0100001001000010)
var range_1_end = sse.Set64(0x3333333333333333, 0x3333333333333333)
var range_0_end = sse.Set64(0x1919191919191919, 0x1919191919191919)

func encode(src []byte) (dst []byte) {
	x0 := &sse.XMM{}
	sse.SetBytes(x0, src)
	x1 := &sse.XMM{}
	sse.SetBytes(x1, encodeStdLut[:])
	x2 := &sse.XMM{}
	x3 := &sse.XMM{}

	sse.PSHUFB(x0, &reshuffle_mask)
	sse.MOVOU(x2, x0)
	sse.PAND(x2, &mulhi_mask)
	sse.PMULHUW(x2, &mulhi_const)
	sse.PAND(x0, &mullo_mask)
	sse.PMULLW(x0, &mullo_const)
	sse.POR(x0, x2)
	sse.MOVOU(x2, x0)
	sse.MOVOU(x3, x0)
	sse.PSUBUSB(x3, &range_1_end)
	sse.PCMPGTB(x2, &range_0_end)
	sse.PSUBB(x3, x2)
	sse.MOVOU(x2, x1)
	sse.PSHUFB(x2, x3)
	sse.PADDB(x0, x2)
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
