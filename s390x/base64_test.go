package s390x

import (
	"bytes"
	"errors"
	"testing"
)

func encode(src, lut []byte) (dst []byte) {
	var (
		X0             = &Vector128{}
		X1             = &Vector128{}
		X2             = &Vector128{}
		REV_BYTES      = &Vector128{}
		reshuffle_mask = &Vector128{}
		MULHI_MASK     = &Vector128{}
		LUT            = &Vector128{}
		MULHI_CONST    = &Vector128{}
		MULLO_MASK     = &Vector128{}
		MULLO_CONST    = &Vector128{}
		RANGE1_END     = &Vector128{}
		RANGE0_END     = &Vector128{}
		ZERO           = &Vector128{}
	)
	VL(lut, LUT)
	VL(src, X0)
	VL_UINT64([]uint64{0x0f0e0d0c0b0a0908, 0x0706050403020100}, REV_BYTES)
	VL_UINT64([]uint64{0x0a0b090a07080607, 0x0405030401020001}, reshuffle_mask)
	VPERM(X0, X0, reshuffle_mask, X0)
	VREPIF(0x0fc0fc00, MULHI_MASK)
	VREPIF(0x04000040, MULHI_CONST)
	VN(X0, MULHI_MASK, X1)
	VMLHH(X1, MULHI_CONST, X1)
	VREPIF(0x003F03F0, MULLO_MASK)
	VREPIF(0x01000010, MULLO_CONST)
	VN(X0, MULLO_MASK, X2)
	VMLHW(X2, MULLO_CONST, X2)
	VO(X1, X2, X0)

	VZERO(ZERO)
	VREPIB(0x33, RANGE1_END)
	VSB(RANGE1_END, X0, X1)
	VMXB(ZERO, X1, X1)

	VREPIB(0x19, RANGE0_END)
	VCGTB(X0, RANGE0_END, X2)
	VSB(X2, X1, X1)

	VPERM(LUT, LUT, X1, X2)
	VAB(X2, X0, X0)

	dst = make([]byte, 16)
	VPERM(X0, X0, REV_BYTES, X0)
	VST(X0, dst)

	return
}

func encodeSTD(src []byte) []byte {
	return encode(src, []byte{65, 71, 252, 252, 252, 252, 252, 252, 252, 252, 252, 252, 237, 240, 0, 0})
}

func TestEncodeSTD(t *testing.T) {
	cases := []struct {
		in        []byte
		exptected string
	}{
		{[]byte("abcdefghijkl0000"), "YWJjZGVmZ2hpamts"},
		{[]byte("\x2b\xf7\xcc\x27\x01\xfe\x43\x97\xb4\x9e\xbe\xed\x5a\xcc\x70\x90"), "K/fMJwH+Q5e0nr7t"},
	}
	for _, c := range cases {
		ret := encodeSTD(c.in)
		if string(ret) != c.exptected {
			t.Errorf("encodeSTD() = %v; want %v", string(ret), c.exptected)
		}
	}
}

func decodeSTD(src []byte) (dst []byte, err error) {
	var (
		X0                   = &Vector128{}
		X1                   = &Vector128{}
		X2                   = &Vector128{}
		X3                   = &Vector128{}
		nibble_mask          = &Vector128{}
		stddec_lut_hi        = &Vector128{}
		stddec_lut_lo        = &Vector128{}
		base64_nibble_mask   = &Vector128{}
		stddec_lut_roll      = &Vector128{}
		dec_reshuffle_const0 = &Vector128{}
		dec_reshuffle_const1 = &Vector128{}
		dec_reshuffle_mask   = &Vector128{}
	)
	VL(src, X0)
	VREPIB(0x0f, nibble_mask)
	VESRLF(4, X0, X1)
	VN(X1, nibble_mask, X1)
	VN(X0, nibble_mask, X2)

	VL_UINT64([]uint64{0x1010010204080408, 0x1010101010101010}, stddec_lut_hi)
	VL_UINT64([]uint64{0x1511111111111111, 0x1111131A1B1B1B1A}, stddec_lut_lo)
	VPERM(stddec_lut_hi, stddec_lut_hi, X1, X3)
	VPERM(stddec_lut_lo, stddec_lut_lo, X2, X2)
	VN(X2, X3, X2)

	// check if the input is valid
	for i := 0; i < 16; i++ {
		if X2.bytes[i] != 0 {
			return nil, errors.New("invalid input")
		}
	}

	VREPIB(0x2f, base64_nibble_mask)
	VCEQB(X0, base64_nibble_mask, X2)
	VAB(X2, X1, X1)

	VL_UINT64([]uint64{0x00101304BFBFB9B9, 0x0000000000000000}, stddec_lut_roll)
	VPERM(stddec_lut_roll, stddec_lut_roll, X1, X2)
	VAB(X0, X2, X0)

	VREPIF(0x40014001, dec_reshuffle_const0)
	VREPIF(0x10000001, dec_reshuffle_const1)
	VMLEB(X0, dec_reshuffle_const0, X1)
	VMLOB(X0, dec_reshuffle_const0, X2)
	VAH(X1, X2, X0)
	VMLEH(X0, dec_reshuffle_const1, X1)
	VMLOH(X0, dec_reshuffle_const1, X2)
	VAF(X1, X2, X0)

	dst = make([]byte, 16)
	VL_UINT64([]uint64{0x010203050607090a, 0x0b0d0e0f00000000}, dec_reshuffle_mask)
	VPERM(X0, X0, dec_reshuffle_mask, X0)
	VST(X0, dst)
	dst = dst[:12]
	return
}

func TestDecodeSTD(t *testing.T) {
	cases := []struct {
		in  string
		out []byte
	}{
		{"YWJjZGVmZ2hpamtsYWJjZGVmZ2hpamts", []byte("abcdefghijkl")},
		{"K/fMJwH+Q5e0nr7tK/fMJwH+Q5e0nr7t", []byte("\x2b\xf7\xcc\x27\x01\xfe\x43\x97\xb4\x9e\xbe\xed")},
	}
	for _, c := range cases {
		ret, err := decodeSTD([]byte(c.in))
		if err != nil {
			t.Errorf("decodeSTD() = %v; want nil", err)
		}
		if !bytes.Equal(ret, c.out) {
			t.Errorf("decodeSTD() = %x; want %x", ret, c.out)
		}
	}
}

func decodeURL(src []byte) (dst []byte, err error) {
	var (
		X0                   = &Vector128{}
		X1                   = &Vector128{}
		X2                   = &Vector128{}
		X3                   = &Vector128{}
		nibble_mask          = &Vector128{}
		dec_lut_hi           = &Vector128{}
		dec_lut_lo           = &Vector128{}
		base64_nibble_mask   = &Vector128{}
		dec_lut_roll         = &Vector128{}
		dec_reshuffle_const0 = &Vector128{}
		dec_reshuffle_const1 = &Vector128{}
		dec_reshuffle_mask   = &Vector128{}
	)
	VL(src, X0)
	VREPIB(0x0f, nibble_mask)
	VESRLF(4, X0, X1)
	VN(X1, nibble_mask, X1)
	VN(X0, nibble_mask, X2)

	VL_UINT64([]uint64{0x1010010204080428, 0x1010101010101010}, dec_lut_hi)
	VL_UINT64([]uint64{0x1511111111111111, 0x1111131B1B1A1B33}, dec_lut_lo)
	VPERM(dec_lut_hi, dec_lut_hi, X1, X3)
	VPERM(dec_lut_lo, dec_lut_lo, X2, X2)
	VN(X2, X3, X2)

	// check if the input is valid
	for i := 0; i < 16; i++ {
		if X2.bytes[i] != 0 {
			return nil, errors.New("invalid input")
		}
	}

	VREPIB(0x5e, base64_nibble_mask)
	VCGTB(X0, base64_nibble_mask, X2)
	VSB(X2, X1, X1)

	VL_UINT64([]uint64{0x00001104BFBFE0B9, 0xB900000000000000}, dec_lut_roll)
	VPERM(dec_lut_roll, dec_lut_roll, X1, X2)
	VAB(X0, X2, X0)

	VREPIF(0x40014001, dec_reshuffle_const0)
	VREPIF(0x10000001, dec_reshuffle_const1)
	VMLEB(X0, dec_reshuffle_const0, X1)
	VMLOB(X0, dec_reshuffle_const0, X2)
	VAH(X1, X2, X0)
	VMLEH(X0, dec_reshuffle_const1, X1)
	VMLOH(X0, dec_reshuffle_const1, X2)
	VAF(X1, X2, X0)

	dst = make([]byte, 16)
	VL_UINT64([]uint64{0x010203050607090a, 0x0b0d0e0f00000000}, dec_reshuffle_mask)
	VPERM(X0, X0, dec_reshuffle_mask, X0)
	VST(X0, dst)
	dst = dst[:12]
	return
}

func TestDecodeURL(t *testing.T) {
	cases := []struct {
		in  string
		out []byte
	}{
		{"IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-", []byte("!?$*&()'-=@~")},
		{"K_fMJwH-Q5e0nr7tK_fMJwH-Q5e0nr7t", []byte("\x2b\xf7\xcc\x27\x01\xfe\x43\x97\xb4\x9e\xbe\xed")},
	}
	for _, c := range cases {
		ret, err := decodeURL([]byte(c.in))
		if err != nil {
			t.Errorf("decodeURL() = %v; want nil", err)
		}
		if !bytes.Equal(ret, c.out) {
			t.Errorf("decodeURL() = %x; want %x", ret, c.out)
		}
	}
}
