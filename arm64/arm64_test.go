package arm64

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestPreTransposeMatrix(t *testing.T) {
	t0 := &Vector128{}
	t1 := &Vector128{}
	t2 := &Vector128{}
	t3 := &Vector128{}

	VLD1_2D([]uint64{0x2222222211111111, 0x4444444433333333}, t0)
	VLD1_2D([]uint64{0x6666666655555555, 0x8888888877777777}, t1)
	VLD1_2D([]uint64{0xaaaaaaaa99999999, 0xccccccccbbbbbbbb}, t2)
	VLD1_2D([]uint64{0xeeeeeeeedddddddd, 0x00000000ffffffff}, t3)

	PRE_TRANSPOSE_S(t0, t1, t2, t3)

	got0 := hex.EncodeToString(t0.Bytes())
	if got0 != "111111115555555599999999dddddddd" {
		t.Errorf("t0 = %v; want 111111115555555599999999dddddddd", got0)
	}
	got1 := hex.EncodeToString(t1.Bytes())
	if got1 != "2222222266666666aaaaaaaaeeeeeeee" {
		t.Errorf("t1 = %v; want 2222222266666666aaaaaaaaeeeeeeee", got1)
	}
	got2 := hex.EncodeToString(t2.Bytes())
	if got2 != "3333333377777777bbbbbbbbffffffff" {
		t.Errorf("t2 = %v; want 3333333377777777bbbbbbbbffffffff", got2)
	}
	got3 := hex.EncodeToString(t3.Bytes())
	if got3 != "4444444488888888cccccccc00000000" {
		t.Errorf("t3 = %v; want 4444444488888888cccccccc00000000", got3)
	}
}

func TestPreTransposeMatrix2(t *testing.T) {
	t0 := &Vector128{}
	t1 := &Vector128{}
	t2 := &Vector128{}
	t3 := &Vector128{}

	VLD1_2D([]uint64{0x2222222211111111, 0x4444444433333333}, t0)
	VLD1_2D([]uint64{0x6666666655555555, 0x8888888877777777}, t1)
	VLD1_2D([]uint64{0xaaaaaaaa99999999, 0xccccccccbbbbbbbb}, t2)
	VLD1_2D([]uint64{0xeeeeeeeedddddddd, 0x00000000ffffffff}, t3)

	PRE_TRANSPOSE_S2(t0, t1, t2, t3)

	got0 := hex.EncodeToString(t0.Bytes())
	if got0 != "111111115555555599999999dddddddd" {
		t.Errorf("t0 = %v; want 111111115555555599999999dddddddd", got0)
	}
	got1 := hex.EncodeToString(t1.Bytes())
	if got1 != "2222222266666666aaaaaaaaeeeeeeee" {
		t.Errorf("t1 = %v; want 2222222266666666aaaaaaaaeeeeeeee", got1)
	}
	got2 := hex.EncodeToString(t2.Bytes())
	if got2 != "3333333377777777bbbbbbbbffffffff" {
		t.Errorf("t2 = %v; want 3333333377777777bbbbbbbbffffffff", got2)
	}
	got3 := hex.EncodeToString(t3.Bytes())
	if got3 != "4444444488888888cccccccc00000000" {
		t.Errorf("t3 = %v; want 4444444488888888cccccccc00000000", got3)
	}
}

func TestTransposeMatrix(t *testing.T) {
	t0 := &Vector128{}
	t1 := &Vector128{}
	t2 := &Vector128{}
	t3 := &Vector128{}

	VLD1_2D([]uint64{0x2222222211111111, 0x4444444433333333}, t0)
	VLD1_2D([]uint64{0x6666666655555555, 0x8888888877777777}, t1)
	VLD1_2D([]uint64{0xaaaaaaaa99999999, 0xccccccccbbbbbbbb}, t2)
	VLD1_2D([]uint64{0xeeeeeeeedddddddd, 0x00000000ffffffff}, t3)

	TRANSPOSE_S(t0, t1, t2, t3)

	got0 := hex.EncodeToString(t0.Bytes())
	if got0 != "dddddddd999999995555555511111111" {
		t.Errorf("t0 = %v; want dddddddd999999995555555511111111", got0)
	}
	got1 := hex.EncodeToString(t1.Bytes())
	if got1 != "eeeeeeeeaaaaaaaa6666666622222222" {
		t.Errorf("t1 = %v; want eeeeeeeeaaaaaaaa6666666622222222", got1)
	}
	got2 := hex.EncodeToString(t2.Bytes())
	if got2 != "ffffffffbbbbbbbb7777777733333333" {
		t.Errorf("t2 = %v; want ffffffffbbbbbbbb7777777733333333", got2)
	}
	got3 := hex.EncodeToString(t3.Bytes())
	if got3 != "00000000cccccccc8888888844444444" {
		t.Errorf("t3 = %v; want 00000000cccccccc8888888844444444", got3)
	}
}

func TestTransposeMatrix2(t *testing.T) {
	t0 := &Vector128{}
	t1 := &Vector128{}
	t2 := &Vector128{}
	t3 := &Vector128{}

	VLD1_2D([]uint64{0x2222222211111111, 0x4444444433333333}, t0)
	VLD1_2D([]uint64{0x6666666655555555, 0x8888888877777777}, t1)
	VLD1_2D([]uint64{0xaaaaaaaa99999999, 0xccccccccbbbbbbbb}, t2)
	VLD1_2D([]uint64{0xeeeeeeeedddddddd, 0x00000000ffffffff}, t3)

	TRANSPOSE_S2(t0, t1, t2, t3)

	got0 := hex.EncodeToString(t0.Bytes())
	if got0 != "dddddddd999999995555555511111111" {
		t.Errorf("t0 = %v; want dddddddd999999995555555511111111", got0)
	}
	got1 := hex.EncodeToString(t1.Bytes())
	if got1 != "eeeeeeeeaaaaaaaa6666666622222222" {
		t.Errorf("t1 = %v; want eeeeeeeeaaaaaaaa6666666622222222", got1)
	}
	got2 := hex.EncodeToString(t2.Bytes())
	if got2 != "ffffffffbbbbbbbb7777777733333333" {
		t.Errorf("t2 = %v; want ffffffffbbbbbbbb7777777733333333", got2)
	}
	got3 := hex.EncodeToString(t3.Bytes())
	if got3 != "00000000cccccccc8888888844444444" {
		t.Errorf("t3 = %v; want 00000000cccccccc8888888844444444", got3)
	}
}

func TestPMULLW(t *testing.T) {
	cases := []struct {
		a, b, expected string
	}{
		{"ffffffffffffffffffffffffffffffff", "ffffffffffffffffffffffffffffffff", "01000100010001000100010001000100"},
		{"ffffffffffffffffffffffffffffffff", "ff7fff7fff7fff7fff7fff7fff7fff7f", "01800180018001800180018001800180"},
		{"ff7fff7fff7fff7fff7fff7fff7fff7f", "ffffffffffffffffffffffffffffffff", "01800180018001800180018001800180"},
		{"ffffffffffffffffffffffffffffffff", "80808080808080808080808080808080", "807f807f807f807f807f807f807f807f"},
		{"80808080808080808080808080808080", "ffffffffffffffffffffffffffffffff", "807f807f807f807f807f807f807f807f"},
		{"80808080808080808080808080808080", "80808080808080808080808080808080", "00400040004000400040004000400040"},
		{"000102030405060708090a0b0c0d0e0f", "101112131415161718191a1b1c1d1e1f", "0010245c50b88424c0a0042d50c9a475"},
		{"00100100001001000010010000100100", "101112131415161718191a1b1c1d1e1f", "000012130040161700801a1b00c01e1f"},
	}
	for _, c := range cases {
		tmp1 := &Vector128{}
		tmp2 := &Vector128{}
		tmp := &Vector128{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		VLD1_16B(aa, tmp1)
		VLD1_16B(bb, tmp2)
		VMUL_H(tmp1, tmp2, tmp)
		if fmt.Sprintf("%x", tmp.Bytes()) != c.expected {
			t.Errorf("PMULLW() = %v; want %v", fmt.Sprintf("%x", tmp.Bytes()), c.expected)
		}
	}
}

func TestPMULHUW(t *testing.T) {
	cases := []struct {
		a, b, expected string
	}{
		{"ffffffffffffffffffffffffffffffff", "ffffffffffffffffffffffffffffffff", "fefffefffefffefffefffefffefffeff"},
		{"ffffffffffffffffffffffffffffffff", "ff7fff7fff7fff7fff7fff7fff7fff7f", "fe7ffe7ffe7ffe7ffe7ffe7ffe7ffe7f"},
		{"ff7fff7fff7fff7fff7fff7fff7fff7f", "ffffffffffffffffffffffffffffffff", "fe7ffe7ffe7ffe7ffe7ffe7ffe7ffe7f"},
		{"ffffffffffffffffffffffffffffffff", "80808080808080808080808080808080", "7f807f807f807f807f807f807f807f80"},
		{"80808080808080808080808080808080", "ffffffffffffffffffffffffffffffff", "7f807f807f807f807f807f807f807f80"},
		{"80808080808080808080808080808080", "80808080808080808080808080808080", "80408040804080408040804080408040"},
		{"000102030405060708090a0b0c0d0e0f", "101112131415161718191a1b1c1d1e1f", "110039006900a200e2002b017b01d401"},
		{"00100100001001000010010000100100", "101112131415161718191a1b1c1d1e1f", "110100005101000091010000d1010000"},
	}
	tmp5 := &Vector128{}
	VLD1_4S([]uint32{0x07060302, 0x0f0e0b0a, 0x17161312, 0x1f1e1b1a}, tmp5)
	for _, c := range cases {
		tmp1 := &Vector128{}
		tmp2 := &Vector128{}
		tmp3 := &Vector128{}
		tmp4 := &Vector128{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		VLD1_16B(aa, tmp1)
		VLD1_16B(bb, tmp2)
		UMULL_H(tmp1, tmp2, tmp3)
		UMULL2_H(tmp1, tmp2, tmp4)
		VTBL_B(tmp5, []*Vector128{tmp3, tmp4}, tmp1)
		if fmt.Sprintf("%x", tmp1.Bytes()) != c.expected {
			t.Errorf("PMULHUW() = %v; want %v", fmt.Sprintf("%x", tmp1.Bytes()), c.expected)
		}
	}
}
