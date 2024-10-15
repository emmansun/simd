package amd64

import (
	"errors"
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
var stddec_lut_hi = sse.Set64(0x1010101010101010, 0x0804080402011010)
var stddec_lut_lo = sse.Set64(0x1A1B1B1B1A131111, 0x1111111111111115)
var stddec_lut_roll = sse.Set64(0x0000000000000000, 0xB9B9BFBF04131000)
var dec_reshuffle_const0 = sse.Set64(0x0140014001400140, 0x0140014001400140)
var dec_reshuffle_const1 = sse.Set64(0x0001100000011000, 0x0001100000011000)
var dec_reshuffle_mask = sse.Set64(0xFFFFFFFF0C0D0E08, 0x090A040506000102)
var base64_nibble_mask = sse.Set64(0x2F2F2F2F2F2F2F2F, 0x2F2F2F2F2F2F2F2F)

func encode(src []byte) (dst []byte) {
	x0 := &sse.XMM{}
	sse.SetBytes(x0, src)
	x1 := &sse.XMM{}
	sse.SetBytes(x1, encodeStdLut[:])
	x2 := &sse.XMM{}
	x3 := &sse.XMM{}

	// enc reshuffle
	// Input, bytes MSB to LSB:
	// 0 0 0 0 l k j i h g f e d c b a
	sse.PSHUFB(x0, &reshuffle_mask)
	// x0, bytes MSB to LSB:
	// k l j k
	// h i g h
	// e f d e
	// b c a b
	sse.MOVOU(x2, x0)
	sse.PAND(x2, &mulhi_mask)
	// bits, upper case are most significant bits, lower case are least significant bits
	// 0000kkkk LL000000 JJJJJJ00 00000000
	// 0000hhhh II000000 GGGGGG00 00000000
	// 0000eeee FF000000 DDDDDD00 00000000
	// 0000bbbb CC000000 AAAAAA00 00000000
	sse.PMULHUW(x2, &mulhi_const) // shift right high 16 bits by 6 and low 16 bits by 10 bits
	// 00000000 00kkkkLL 00000000 00JJJJJJ
	// 00000000 00hhhhII 00000000 00GGGGGG
	// 00000000 00eeeeFF 00000000 00DDDDDD
	// 00000000 00bbbbCC 00000000 00AAAAAA
	sse.PAND(x0, &mullo_mask)
	// 00000000 00llllll 000000jj KKKK0000
	// 00000000 00iiiiii 000000gg HHHH0000
	// 00000000 00ffffff 000000dd EEEE0000
	// 00000000 00cccccc 000000aa BBBB0000
	sse.PMULLW(x0, &mullo_const) // shift left high 16 bits by 8 bits, and low 16 bits by 4 bits
	// 00llllll 00000000 00jjKKKK 00000000
	// 00iiiiii 00000000 00ggHHHH 00000000
	// 00ffffff 00000000 00ddEEEE 00000000
	// 00cccccc 00000000 00aaBBBB 00000000
	sse.POR(x0, x2)
	// 00llllll 00kkkkLL 00jjKKKK 00JJJJJJ
	// 00iiiiii 00hhhhII 00ggHHHH 00GGGGGG
	// 00ffffff 00eeeeFF 00ddEEEE 00DDDDDD
	// 00cccccc 00bbbbCC 00aaBBBB 00AAAAAA
	sse.MOVOU(x2, x0)
	sse.MOVOU(x3, x0)
	sse.PSUBUSB(x3, &range_1_end) // Create LUT indices from the input. The index for range #0 is right, others are 1 less than expected.
	sse.PCMPGTB(x2, &range_0_end) // mask is 0xFF (-1) for range #[1..4] and 0x00 for range #0.
	sse.PSUBB(x3, x2)             // Subtract -1, so add 1 to indices for range #[1..4]. All indices are now correct.
	sse.MOVOU(x2, x1)
	sse.PSHUFB(x2, x3)            // get offsets from LUT.
	sse.PADDB(x0, x2)             // add offsets to input value.
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

func decode(src []byte) (dst []byte, err error) {
	x0 := &sse.XMM{}
	x1 := &sse.XMM{}
	x2 := &sse.XMM{}
	x3 := &sse.XMM{}
	x4 := &sse.XMM{}
	zero := sse.Set64(0, 0)
	// validate input
	sse.MOVOU(x3, &stddec_lut_hi)
	sse.MOVOU(x4, &stddec_lut_lo)
	sse.SetBytes(x0, src)
	sse.MOVOU(x1, x0)
	sse.MOVOU(x2, x0)
	sse.PSRLW(x1, 4)
	sse.PAND(x1, &base64_nibble_mask)
	sse.PAND(x2, &base64_nibble_mask)
	sse.PSHUFB(x3, x1)
	sse.PSHUFB(x4, x2)
	sse.PAND(x4, x3)
	sse.PCMPGTB(x4, &zero)
	ret := sse.PMOVMSKB(x4)
	if ret != 0 {
		return nil, errors.New("invalid input")
	}
	// uses PMADDUBSW (_mm_maddubs_epi16) / PMADDWD (_mm_madd_epi16) and PSHUFB to reshuffle bits.
	//
	// in bits, upper case are most significant bits, lower case are least significant bits
	// 00llllll 00kkkkLL 00jjKKKK 00JJJJJJ
	// 00iiiiii 00hhhhII 00ggHHHH 00GGGGGG
	// 00ffffff 00eeeeFF 00ddEEEE 00DDDDDD
	// 00cccccc 00bbbbCC 00aaBBBB 00AAAAAA
	//
	// out bits, upper case are most significant bits, lower case are least significant bits:
	// 00000000 00000000 00000000 00000000
	// LLllllll KKKKkkkk JJJJJJjj IIiiiiii
	// HHHHhhhh GGGGGGgg FFffffff EEEEeeee
	// DDDDDDdd CCcccccc BBBBbbbb AAAAAAaa
	sse.MOVOU(x2, &base64_nibble_mask)
	sse.PCMPEQB(x2, x0)
	sse.PADDB(x1, x2)
	sse.MOVOU(x2, &stddec_lut_roll)
	sse.PSHUFB(x2, x1)
	sse.PADDB(x0, x2)
	sse.PMADDUBSW(x0, &dec_reshuffle_const0)
	sse.PMADDWD(x0, &dec_reshuffle_const1)
	sse.PSHUFB(x0, &dec_reshuffle_mask)
	dst = make([]byte, 12)
	copy(dst, x0.Bytes()[:12])
	return
}

func TestDecode(t *testing.T) {
	ret, err := decode([]byte("YWJjZGVmZ2hpamtsYWJjZGVmZ2hpamts"))
	if err != nil {
		t.Errorf("decode() = %v; want nil", err)
	}
	if string(ret) != "abcdefghijkl" {
		t.Errorf("encode() = %v; want abcdefghijkl", string(ret))
	}
}
