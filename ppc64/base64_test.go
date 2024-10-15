package ppc64

import (
	"bytes"
	"errors"
	"testing"
)

func encode(src, lut []byte, isPPC64LE bool) (dst []byte) {
	X0 := &Vector128{}
	rev_bytes := &Vector128{}

	if isPPC64LE {
		LXVD2X_PPC64LE(src, X0)
		LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, rev_bytes)
		reshuffle_mask := &Vector128{}
		LXVD2X_UINT64([]uint64{0x0d0c0e0d000f0100, 0x0302040306050706}, reshuffle_mask)
		VPERM(X0, X0, reshuffle_mask, X0)
	} else {
		LXVD2X(src, X0)
		LXVD2X_UINT64([]uint64{0x0f0e0d0c0b0a0908, 0x0706050403020100}, rev_bytes)
		reshuffle_mask := &Vector128{}
		LXVD2X_UINT64([]uint64{0x0a0b090a07080607, 0x0405030401020001}, reshuffle_mask)
		VPERM(X0, X0, reshuffle_mask, X0)
	}

	X1 := &Vector128{}
	VOR(X0, X0, X1)

	mulhi_mask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x0fc0fc000fc0fc00, 0x0fc0fc000fc0fc00}, mulhi_mask)
	VAND(X1, mulhi_mask, X1)
	shiftrightMask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x0006000a0006000a, 0x0006000a0006000a}, shiftrightMask)
	VSRH(X1, shiftrightMask, X1)
	mullo_mask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x003F03F0003F03F0, 0x003F03F0003F03F0}, mullo_mask)
	VAND(X0, mullo_mask, X0)
	shiftleftMask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x0008000400080004, 0x0008000400080004}, shiftleftMask)
	VSLH(X0, shiftleftMask, X0)
	VOR(X1, X0, X0)
	range_1_end := &Vector128{}
	LXVD2X_UINT64([]uint64{0x3333333333333333, 0x3333333333333333}, range_1_end)
	VSUBUBS(range_1_end, X0, X1)
	range_0_end := &Vector128{}
	LXVD2X_UINT64([]uint64{0x1919191919191919, 0x1919191919191919}, range_0_end)
	X2 := &Vector128{}
	VCMPGTUB(range_0_end, X0, X2)
	VSUBUBM(X2, X1, X1)

	if isPPC64LE {
		LXVD2X_PPC64LE(lut, X2)
		VPERM(X2, X2, rev_bytes, X2)
	} else {
		LXVD2X(lut, X2)
	}

	VPERM(X2, X2, X1, X2)
	VADDUBM(X2, X0, X0)
	dst = make([]byte, 16)

	if isPPC64LE {
		XXPERMDI(X0, X0, 2, X0)
		STXVD2X_PPC64LE(X0, dst)
	} else {
		VPERM(X0, X0, rev_bytes, X0)
		STXVD2X(X0, dst)
	}

	return
}

func encodeSTD(src []byte, isPPC64LE bool) ([]byte) {
	return encode(src, []byte{65, 71, 252, 252, 252, 252, 252, 252, 252, 252, 252, 252, 237, 240, 0, 0}, isPPC64LE)
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
		ret := encodeSTD(c.in, false)
		if string(ret) != c.exptected {
			t.Errorf("encodeSTD() = %v; want %v", string(ret), c.exptected)
		}
		ret = encodeSTD(c.in, true)
		if string(ret) != c.exptected {
			t.Errorf("encodeSTD() = %v; want %v", string(ret), c.exptected)
		}
	}
}

func encodeURL(src []byte, isPPC64LE bool) ([]byte) {
	return encode(src, []byte{65, 71, 252, 252, 252, 252, 252, 252, 252, 252, 252, 252, 239, 32, 0, 0}, isPPC64LE)
}

func TestEncodeURL(t *testing.T) {
	cases := []struct {
		in        []byte
		exptected string
	}{
		{[]byte("!?$*&()'-=@~0000"), "IT8kKiYoKSctPUB-"},
		{[]byte("\x2b\xf7\xcc\x27\x01\xfe\x43\x97\xb4\x9e\xbe\xed\x5a\xcc\x70\x90"), "K_fMJwH-Q5e0nr7t"},
	}
	for _, c := range cases {
		ret := encodeURL(c.in, false)
		if string(ret) != c.exptected {
			t.Errorf("encodeURL() = %v; want %v", string(ret), c.exptected)
		}
		ret = encodeURL(c.in, true)
		if string(ret) != c.exptected {
			t.Errorf("encodeURL() = %v; want %v", string(ret), c.exptected)
		}
	}
}

func decodeSTD(src []byte, isPPC64LE bool) (dst []byte, err error) {
	X0 := &Vector128{}
	if isPPC64LE {
		LXVD2X_PPC64LE(src, X0)
		rev_bytes := &Vector128{}
		LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, rev_bytes)
		VPERM(X0, X0, rev_bytes, X0)
	} else {
		LXVD2X(src, X0)
	}
	nibble_mask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x0F0F0F0F0F0F0F0F, 0x0F0F0F0F0F0F0F0F}, nibble_mask)
	FOUR := &Vector128{}
	VSPLTISB(4, FOUR)
	X1 := &Vector128{}
	X2 := &Vector128{}
	X3 := &Vector128{}
	VSRW(X0, FOUR, X1)
	VAND(X1, nibble_mask, X1)
	VAND(X0, nibble_mask, X2)
	stddec_lut_hi := &Vector128{}
	stddec_lut_lo := &Vector128{}
	LXVD2X_UINT64([]uint64{0x1010010204080408, 0x1010101010101010}, stddec_lut_hi)
	LXVD2X_UINT64([]uint64{0x1511111111111111, 0x1111131A1B1B1B1A}, stddec_lut_lo)
	VPERM(stddec_lut_hi, stddec_lut_hi, X1, X3)
	VPERM(stddec_lut_lo, stddec_lut_lo, X2, X2)
	VAND(X2, X3, X2)

	// check if the input is valid
	// should use VCMPEQUBCC
	for i := 0; i < 16; i++ {
		if X2.bytes[i] != 0 {
			return nil, errors.New("invalid input")
		}
	}

	base64_nibble_mask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x2F2F2F2F2F2F2F2F, 0x2F2F2F2F2F2F2F2F}, base64_nibble_mask)
	VCMPEQUB(X0, base64_nibble_mask, X2)
	VADDUBM(X1, X2, X1)

	stddec_lut_roll := &Vector128{}
	LXVD2X_UINT64([]uint64{0x00101304BFBFB9B9, 0x0000000000000000}, stddec_lut_roll)
	VPERM(stddec_lut_roll, stddec_lut_roll, X1, X1)
	VADDUBM(X0, X1, X0)

	dec_reshuffle_const0 := &Vector128{}
	LXVD2X_UINT64([]uint64{0x4001400140014001, 0x4001400140014001}, dec_reshuffle_const0)
	VMULEUB(X0, dec_reshuffle_const0, X1)
	VMULOUB(X0, dec_reshuffle_const0, X2)
	VADDUHM(X1, X2, X0)
	dec_reshuffle_const1 := &Vector128{}
	LXVD2X_UINT64([]uint64{0x1000000110000001, 0x1000000110000001}, dec_reshuffle_const1)
	VMULEUH(X0, dec_reshuffle_const1, X1)
	VMULOUH(X0, dec_reshuffle_const1, X2)
	VADDUWM(X1, X2, X0)

	dst = make([]byte, 16)
	dec_reshuffle_mask := &Vector128{}
	if isPPC64LE {
		LXVD2X_UINT64([]uint64{0x0A09070605030201, 0x000000000f0e0d0b}, dec_reshuffle_mask)
		VPERM(X0, X0, dec_reshuffle_mask, X0)
		STXVD2X_PPC64LE(X0, dst)
		dst = dst[:12]
	} else {
		LXVD2X_UINT64([]uint64{0x010203050607090a, 0x0b0d0e0f00000000}, dec_reshuffle_mask)
		VPERM(X0, X0, dec_reshuffle_mask, X0)
		STXVD2X(X0, dst)
		dst = dst[:12]
	}

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
		ret, err := decodeSTD([]byte(c.in), false)
		if err != nil {
			t.Errorf("decodeSTD() = %v; want nil", err)
		}
		if !bytes.Equal(ret, c.out) {
			t.Errorf("decodeSTD() = %x; want %x", ret, c.out)
		}
		ret, err = decodeSTD([]byte(c.in), false)
		if err != nil {
			t.Errorf("decodeSTD() = %v; want nil", err)
		}
		if !bytes.Equal(ret, c.out) {
			t.Errorf("decodeSTD() = %x; want %x", ret, c.out)
		}
	}
}

func decodeURL(src []byte, isPPC64LE bool) (dst []byte, err error) {
	X0 := &Vector128{}
	if isPPC64LE {
		LXVD2X_PPC64LE(src, X0)
		rev_bytes := &Vector128{}
		LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, rev_bytes)
		VPERM(X0, X0, rev_bytes, X0)
	} else {
		LXVD2X(src, X0)
	}
	nibble_mask := &Vector128{}
	LXVD2X_UINT64([]uint64{0x0F0F0F0F0F0F0F0F, 0x0F0F0F0F0F0F0F0F}, nibble_mask)
	FOUR := &Vector128{}
	VSPLTISB(4, FOUR)
	X1 := &Vector128{}
	X2 := &Vector128{}
	X3 := &Vector128{}
	VSRW(X0, FOUR, X1)
	VAND(X1, nibble_mask, X1)
	VAND(X0, nibble_mask, X2)
	dec_lut_hi := &Vector128{}
	dec_lut_lo := &Vector128{}
	LXVD2X_UINT64([]uint64{0x1010010204080428, 0x1010101010101010}, dec_lut_hi)
	LXVD2X_UINT64([]uint64{0x1511111111111111, 0x1111131B1B1A1B33}, dec_lut_lo)
	VPERM(dec_lut_hi, dec_lut_hi, X1, X3)
	VPERM(dec_lut_lo, dec_lut_lo, X2, X2)
	VAND(X2, X3, X2)

	// check if the input is valid
	// should use VCMPEQUBCC
	for i := 0; i < 16; i++ {
		if X2.bytes[i] != 0 {
			return nil, errors.New("invalid input")
		}
	}

	url_const_5e := &Vector128{}
	LXVD2X_UINT64([]uint64{0x5E5E5E5E5E5E5E5E, 0x5E5E5E5E5E5E5E5E}, url_const_5e)
	VCMPGTUB(url_const_5e, X0, X2)
	VSUBUBM(X2, X1, X1)

	dec_lut_roll := &Vector128{}
	LXVD2X_UINT64([]uint64{0x00001104BFBFE0B9, 0xB900000000000000}, dec_lut_roll)
	VPERM(dec_lut_roll, dec_lut_roll, X1, X1)
	VADDUBM(X0, X1, X0)

	dec_reshuffle_const0 := &Vector128{}
	LXVD2X_UINT64([]uint64{0x4001400140014001, 0x4001400140014001}, dec_reshuffle_const0)
	VMULEUB(X0, dec_reshuffle_const0, X1)
	VMULOUB(X0, dec_reshuffle_const0, X2)
	VADDUHM(X1, X2, X0)
	dec_reshuffle_const1 := &Vector128{}
	LXVD2X_UINT64([]uint64{0x1000000110000001, 0x1000000110000001}, dec_reshuffle_const1)
	VMULEUH(X0, dec_reshuffle_const1, X1)
	VMULOUH(X0, dec_reshuffle_const1, X2)
	VADDUWM(X1, X2, X0)

	dst = make([]byte, 16)
	dec_reshuffle_mask := &Vector128{}
	if isPPC64LE {
		LXVD2X_UINT64([]uint64{0x0A09070605030201, 0x000000000f0e0d0b}, dec_reshuffle_mask)
		VPERM(X0, X0, dec_reshuffle_mask, X0)
		STXVD2X_PPC64LE(X0, dst)
		dst = dst[:12]
	} else {
		LXVD2X_UINT64([]uint64{0x010203050607090a, 0x0b0d0e0f00000000}, dec_reshuffle_mask)
		VPERM(X0, X0, dec_reshuffle_mask, X0)
		STXVD2X(X0, dst)
		dst = dst[:12]
	}

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
		ret, err := decodeURL([]byte(c.in), false)
		if err != nil {
			t.Errorf("decodeURL() = %v; want nil", err)
		}
		if !bytes.Equal(ret, c.out) {
			t.Errorf("decodeURL() = %x; want %x", ret, c.out)
		}
		ret, err = decodeURL([]byte(c.in), false)
		if err != nil {
			t.Errorf("decodeURL() = %v; want nil", err)
		}
		if !bytes.Equal(ret, c.out) {
			t.Errorf("decodeURL() = %x; want %x", ret, c.out)
		}
	}
}
