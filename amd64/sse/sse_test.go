package sse

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestPSRAW(t *testing.T) {
	tmp1 := Set64(0xf070b030d0509010, 0xe060a020c0408000)
	tmp := XMM{}
	MOVOU(&tmp, &tmp1)
	PSRAW(&tmp1, 32)
	if fmt.Sprintf("%x", tmp1.Bytes()) != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("PSRAW() = %v; want ffffffff...", fmt.Sprintf("%x", tmp1.Bytes()))
	}
	PSRAW(&tmp, 4)
	if fmt.Sprintf("%x", tmp.Bytes()) != "000804fc020a06fe010905fd030b07ff" {
		t.Errorf("PSRAW() = %v; want 000804fc020a06fe010905fd030b07ff", fmt.Sprintf("%x", tmp.Bytes()))
	}
}

func TestPMADDUBSW(t *testing.T) {
	cases := []struct {
		a, b, expected string
	}{
		{"ffffffffffffffffffffffffffffffff", "ffffffffffffffffffffffffffffffff", "02fe02fe02fe02fe02fe02fe02fe02fe"},
		{"ffffffffffffffffffffffffffffffff", "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f", "ff7fff7fff7fff7fff7fff7fff7fff7f"},
		{"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f", "ffffffffffffffffffffffffffffffff", "02ff02ff02ff02ff02ff02ff02ff02ff"},
		{"ffffffffffffffffffffffffffffffff", "80808080808080808080808080808080", "00800080008000800080008000800080"},
		{"80808080808080808080808080808080", "ffffffffffffffffffffffffffffffff", "00ff00ff00ff00ff00ff00ff00ff00ff"},
		{"000102030405060708090a0b0c0d0e0f", "101112131415161718191a1b1c1d1e1f", "11005d00b9002501a1012d02c9027503"},
	}
	for _, c := range cases {
		tmp1 := &XMM{}
		tmp2 := &XMM{}
		tmp := &XMM{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		SetBytes(tmp, aa)
		SetBytes(tmp1, aa)
		SetBytes(tmp2, bb)
		PMADDUBSW(tmp1, tmp2)
		if fmt.Sprintf("%x", tmp1.Bytes()) != c.expected {
			t.Errorf("PMADDUBSW() = %v; want %v", fmt.Sprintf("%x", tmp1.Bytes()), c.expected)
		}
	}
}

func TestPSUBUSB(t *testing.T) {
	cases := []struct {
		a, b, expected string
	}{
		{"ffffffffffffffffffffffffffffffff", "ffffffffffffffffffffffffffffffff", "00000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffff", "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f", "80808080808080808080808080808080"},
		{"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f", "ffffffffffffffffffffffffffffffff", "00000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffff", "80808080808080808080808080808080", "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f"},
		{"80808080808080808080808080808080", "ffffffffffffffffffffffffffffffff", "00000000000000000000000000000000"},
		{"40014001400140014001400140014001", "101112131415161718191a1b1c1d1e1f", "30002e002c002a002800260024002200"},
	}
	for _, c := range cases {
		tmp1 := &XMM{}
		tmp2 := &XMM{}
		tmp := &XMM{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		SetBytes(tmp, aa)
		SetBytes(tmp1, aa)
		SetBytes(tmp2, bb)
		PSUBUSB(tmp1, tmp2)
		if fmt.Sprintf("%x", tmp1.Bytes()) != c.expected {
			t.Errorf("PSUBUSB() = %v; want %v", fmt.Sprintf("%x", tmp1.Bytes()), c.expected)
		}
	}
}

func TestPMADDWD(t *testing.T) {
	cases := []struct {
		a, b, expected string
	}{
		{"ffffffffffffffffffffffffffffffff", "ffffffffffffffffffffffffffffffff", "02000000020000000200000002000000"},
		{"ffffffffffffffffffffffffffffffff", "ff7fff7fff7fff7fff7fff7fff7fff7f", "0200ffff0200ffff0200ffff0200ffff"},
		{"ff7fff7fff7fff7fff7fff7fff7fff7f", "ffffffffffffffffffffffffffffffff", "0200ffff0200ffff0200ffff0200ffff"},
		{"ffffffffffffffffffffffffffffffff", "80808080808080808080808080808080", "00ff000000ff000000ff000000ff0000"},
		{"80808080808080808080808080808080", "ffffffffffffffffffffffffffffffff", "00ff000000ff000000ff000000ff0000"},
		{"80808080808080808080808080808080", "80808080808080808080808080808080", "0080007f0080007f0080007f0080007f"},
		{"000102030405060708090a0b0c0d0e0f", "101112131415161718191a1b1c1d1e1f", "246c4a00d4dc0b01c4cd0d02f43e5003"},
		{"00100100001001000010010000100100", "101112131415161718191a1b1c1d1e1f", "12131101165751011a9b91011edfd101"},
	}
	for _, c := range cases {
		tmp1 := &XMM{}
		tmp2 := &XMM{}
		tmp := &XMM{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		SetBytes(tmp, aa)
		SetBytes(tmp1, aa)
		SetBytes(tmp2, bb)
		PMADDWD(tmp1, tmp2)
		if fmt.Sprintf("%x", tmp1.Bytes()) != c.expected {
			t.Errorf("PMADDWD() = %v; want %v", fmt.Sprintf("%x", tmp1.Bytes()), c.expected)
		}
	}
}

func TestPMOVMSKB(t *testing.T) {
	cases := []struct {
		a string
		expected uint64
	}{
		{"ffffffffffffffffffffffffffffffff", 65535},
		{"ff7fff7fff7fff7fff7fff7fff7fff7f", 21845},
		{"a0b112131415161718191a1b1c1d1e1f", 3},
		{"101112131415161718191a1b1c1d8e8f", 49152},
	}
	for _, c := range cases {
		tmp1 := &XMM{}
		aa, _ := hex.DecodeString(c.a)
		SetBytes(tmp1, aa)
		ret1 := PMOVMSKB(tmp1)
		if ret1 != c.expected {
			t.Errorf("PMOVMSKB() = %v; want %v", ret1, c.expected)
		}
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
		tmp1 := &XMM{}
		tmp2 := &XMM{}
		tmp := &XMM{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		SetBytes(tmp, aa)
		SetBytes(tmp1, aa)
		SetBytes(tmp2, bb)
		PMULLW(tmp1, tmp2)
		if fmt.Sprintf("%x", tmp1.Bytes()) != c.expected {
			t.Errorf("PMULLW() = %v; want %v", fmt.Sprintf("%x", tmp1.Bytes()), c.expected)
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
	for _, c := range cases {
		tmp1 := &XMM{}
		tmp2 := &XMM{}
		tmp := &XMM{}
		aa, _ := hex.DecodeString(c.a)
		bb, _ := hex.DecodeString(c.b)
		SetBytes(tmp, aa)
		SetBytes(tmp1, aa)
		SetBytes(tmp2, bb)
		PMULHUW(tmp1, tmp2)
		if fmt.Sprintf("%x", tmp1.Bytes()) != c.expected {
			t.Errorf("PMULHUW() = %v; want %v", fmt.Sprintf("%x", tmp1.Bytes()), c.expected)
		}
	}
}
