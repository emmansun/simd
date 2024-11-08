package arm64

import (
	"bytes"
	"testing"
)

var encodeStd = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
var encodeURL = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

var dencodeStdLut = [128]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 62, 255, 255, 255, 63,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 255, 255, 255, 255, 255, 255,
	255, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 255, 255, 255, 255,
	255, 255, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 255, 255, 255, 255, 255,
}

var dencodeUrlLut = [128]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 62, 255, 255,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 255, 255, 255, 255, 255, 255,
	255, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 255, 255, 255, 255,
	63, 255, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 255, 255, 255, 255, 255,
}

func encode(dst, src []byte, lut []byte) {
	var (
		V0  = &Vector128{}
		V1  = &Vector128{}
		V2  = &Vector128{}
		V3  = &Vector128{}
		V4  = &Vector128{}
		V5  = &Vector128{}
		V6  = &Vector128{}
		V7  = &Vector128{}
		V8  = &Vector128{}
		V9  = &Vector128{}
		V10 = &Vector128{}
		V11 = &Vector128{}
	)
	// load lookup table
	VLD1_16B(lut[:], V8)
	VLD1_16B(lut[16:], V9)
	VLD1_16B(lut[32:], V10)
	VLD1_16B(lut[48:], V11)

	VDUP_BYTE(0x3f, V7)
	// load input
	// V2=V S P M J G D A x u r o l i f c
	// V1=U R O L I F C z w t q n k h e b
	// V0=T Q N K H E B y v s p m j g d a
	VLD3_16B(src, V0, V1, V2)

	VUSHR_B(2, V0, V3) // V3 = .. 00DDDDDD 00AAAAAA
	VUSHR_B(4, V1, V4) // V4 = .. 0000EEEE 0000BBBB
	VUSHR_B(6, V2, V5) // V5 = .. 000000FF 000000CC
	VSLI_B(4, V0, V4)  // V4 = .. ddddEEEE aaaaBBBB
	VSLI_B(2, V1, V5)  // V5 = .. eeeeeeFF bbbbbbCC

	// Clear the high two bits in the second, third and fourth output.
	// V3 = .. 00DDDDDD 00AAAAAA
	VAND(V7, V4, V4) // V4 = .. 00ddEEEE 00aaBBBB
	VAND(V7, V5, V5) // V5 = .. 00eeeeFF 00bbbbCC
	VAND(V7, V2, V6) // V6 = .. 00ffffff 00cccccc

	// The bits have now been shifted to the right locations;
	// translate their values 0..63 to the Base64 alphabet.
	// Use a 64-byte table lookup:

	VTBL_B(V3, []*Vector128{V8, V9, V10, V11}, V3)
	VTBL_B(V4, []*Vector128{V8, V9, V10, V11}, V4)
	VTBL_B(V5, []*Vector128{V8, V9, V10, V11}, V5)
	VTBL_B(V6, []*Vector128{V8, V9, V10, V11}, V6)

	// Interleave and store output:
	VST4_16B(V3, V4, V5, V6, dst)
}

func encode16Bytes(dst, src []byte, lut []byte) {
	var (
		V0  = &Vector128{}
		V1  = &Vector128{}
		V2  = &Vector128{}
		V3  = &Vector128{}
		V4  = &Vector128{}
		V5  = &Vector128{}
		V6  = &Vector128{}
		V7  = &Vector128{}
		V8  = &Vector128{}
		V9  = &Vector128{}
		V10 = &Vector128{}
		V11 = &Vector128{}
		V12 = &Vector128{}
	)
	// load constant
	VLD1_2D([]uint64{0x0405030401020001, 0x0a0b090a07080607}, V3) // reshuffle_mask
	VLD1_2D([]uint64{0x0FC0FC000FC0FC00, 0x0FC0FC000FC0FC00}, V4) // mulhi_mask
	//VLD1_2D([]uint64{0x0400004004000040, 0x0400004004000040}, V5) // mulhi_const
	VLD1_2D([]uint64{0x003F03F0003F03F0, 0x003F03F0003F03F0}, V6)  // mullo_mask
	VLD1_2D([]uint64{0x0100001001000010, 0x0100001001000010}, V7)  // mullo_const
	VLD1_2D([]uint64{0x1f1e1b1a17161312, 0x0f0e0b0a07060302}, V12) // high part of word
	VSHL_S(2, V7, V5)

	// load lookup table
	VLD1_16B(lut[:], V8)
	VLD1_16B(lut[16:], V9)
	VLD1_16B(lut[32:], V10)
	VLD1_16B(lut[48:], V11)

	VLD1_16B(src, V0)
	VTBL_B(V3, []*Vector128{V0}, V0)
	VAND(V0, V4, V1)

	UMULL_H(V1, V5, V2)
	UMULL2_H(V1, V5, V1)
	VTBL_B(V12, []*Vector128{V1, V2}, V1)

	VAND(V0, V6, V0)
	VMUL_H(V0, V7, V0)
	VORR(V1, V0, V0)

	// The bits have now been shifted to the right locations;
	// translate their values 0..63 to the Base64 alphabet.
	// Use a 64-byte table lookup:
	VTBL_B(V0, []*Vector128{V8, V9, V10, V11}, V0)
	VST1_16B(V0, dst)
}

func TestEncodeSTD(t *testing.T) {
	cases := []struct {
		in       []byte
		expected string
	}{
		{[]byte("abcdefghijklabcdefghijklabcdefghijklabcdefghijkl"), "YWJjZGVmZ2hpamtsYWJjZGVmZ2hpamtsYWJjZGVmZ2hpamtsYWJjZGVmZ2hpamts"},
	}
	for _, c := range cases {
		ret := make([]byte, 64)
		encode(ret, c.in, encodeStd)
		if string(ret) != c.expected {
			t.Errorf("encode() = %v; want %v", string(ret), c.expected)
		}
	}
}

func TestEncodeURL(t *testing.T) {
	cases := []struct {
		in       []byte
		expected string
	}{
		{[]byte("!?$*&()'-=@~!?$*&()'-=@~!?$*&()'-=@~!?$*&()'-=@~"), "IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-"},
	}
	for _, c := range cases {
		ret := make([]byte, 64)
		encode(ret, c.in, encodeURL)
		if string(ret) != c.expected {
			t.Errorf("encode() = %v; want %v", string(ret), c.expected)
		}
	}
}

func TestEncode16BSTD(t *testing.T) {
	cases := []struct {
		in       []byte
		expected string
	}{
		{[]byte("abcdefghijkl0000"), "YWJjZGVmZ2hpamts"},
	}
	for _, c := range cases {
		ret := make([]byte, 16)
		encode16Bytes(ret, c.in, encodeStd)
		if string(ret) != c.expected {
			t.Errorf("encode() = %v; want %v", string(ret), c.expected)
		}
	}
}

func decode(dst, src []byte, lut *[128]byte) {
	var (
		V0  = &Vector128{}
		V1  = &Vector128{}
		V2  = &Vector128{}
		V3  = &Vector128{}
		V4  = &Vector128{}
		V5  = &Vector128{}
		V6  = &Vector128{}
		V7  = &Vector128{}
		V8  = &Vector128{}
		V9  = &Vector128{}
		V10 = &Vector128{}
		V11 = &Vector128{}
		V12 = &Vector128{}
		V13 = &Vector128{}
		V14 = &Vector128{}
		V15 = &Vector128{}
		V16 = &Vector128{}
		V17 = &Vector128{}
		V18 = &Vector128{}
		V19 = &Vector128{}
		V20 = &Vector128{}
		V21 = &Vector128{}
		V22 = &Vector128{}
		V23 = &Vector128{}
	)

	VLD1_16B(lut[:], V8)
	VLD1_16B(lut[16:], V9)
	VLD1_16B(lut[32:], V10)
	VLD1_16B(lut[48:], V11)
	VLD1_16B(lut[64:], V12)
	VLD1_16B(lut[80:], V13)
	VLD1_16B(lut[96:], V14)
	VLD1_16B(lut[112:], V15)

	VDUP_BYTE(0x40, V7)

	// load input
	VLD4_16B(src, V0, V1, V2, V3)

	// Get values from first LUT:
	VTBL_B(V0, []*Vector128{V8, V9, V10, V11}, V20)
	VTBL_B(V1, []*Vector128{V8, V9, V10, V11}, V21)
	VTBL_B(V2, []*Vector128{V8, V9, V10, V11}, V22)
	VTBL_B(V3, []*Vector128{V8, V9, V10, V11}, V23)

	// Get values from second LUT:
	VSUB_B(V7, V0, V0)
	VTBX_B(V0, []*Vector128{V12, V13, V14, V15}, V20)
	VSUB_B(V7, V1, V1)
	VTBX_B(V1, []*Vector128{V12, V13, V14, V15}, V21)
	VSUB_B(V7, V2, V2)
	VTBX_B(V2, []*Vector128{V12, V13, V14, V15}, V22)
	VSUB_B(V7, V3, V3)
	VTBX_B(V3, []*Vector128{V12, V13, V14, V15}, V23)

	// Check for invalid input, any value larger than 63:
	VCMHS_B(V7, V20, V16)
	VCMHS_B(V7, V21, V17)
	VCMHS_B(V7, V22, V18)
	VCMHS_B(V7, V23, V19)

	VORR(V16, V17, V16)
	VORR(V18, V19, V18)
	VORR(V16, V18, V16)

	// Check that all bits are zero:
	VUMAXV_B(true, V16, V17)
	if V17.bytes[15] != 0 {
		panic("invalid input")
	}

	// Compress four bytes into three:
	VSHL_B(2, V20, V4)
	VUSHR_B(4, V21, V16)
	VORR(V4, V16, V4)

	VSHL_B(4, V21, V5)
	VUSHR_B(2, V22, V16)
	VORR(V5, V16, V5)

	VSHL_B(6, V22, V16)
	VORR(V16, V23, V6)

	VST3_16B(V4, V5, V6, dst)
}

func TestDecodeSTD(t *testing.T) {
	cases := []struct {
		in       string
		expected []byte
	}{
		{"YWJjZGVmZ2hpamtsYWJjZGVmZ2hpamtsYWJjZGVmZ2hpamtsYWJjZGVmZ2hpamts", []byte("abcdefghijklabcdefghijklabcdefghijklabcdefghijkl")},
	}
	for _, c := range cases {
		ret := make([]byte, 48)
		decode(ret, []byte(c.in), &dencodeStdLut)
		if !bytes.Equal(ret, c.expected) {
			t.Errorf("decode() = %x; want %x", ret, c.expected)
		}
	}
}

func TestDecodeURL(t *testing.T) {
	cases := []struct {
		in       string
		expected []byte
	}{
		{"IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-IT8kKiYoKSctPUB-", []byte("!?$*&()'-=@~!?$*&()'-=@~!?$*&()'-=@~!?$*&()'-=@~")},
	}
	for _, c := range cases {
		ret := make([]byte, 48)
		decode(ret, []byte(c.in), &dencodeUrlLut)
		if !bytes.Equal(ret, c.expected) {
			t.Errorf("decode() = %x; want %x", ret, c.expected)
		}
	}
}

func decode16B(dst, src []byte, lut *[128]byte) {
	var (
		V0 = &Vector128{}

		V1 = &Vector128{}
		V2 = &Vector128{}
		V3 = &Vector128{}
		V4 = &Vector128{}
		//V5  = &Vector128{}
		//V6  = &Vector128{}
		V7  = &Vector128{}
		V8  = &Vector128{}
		V9  = &Vector128{}
		V10 = &Vector128{}
		V11 = &Vector128{}
		V12 = &Vector128{}
		V13 = &Vector128{}
		V14 = &Vector128{}
		V15 = &Vector128{}
		V16 = &Vector128{}
		V17 = &Vector128{}
		//V18 = &Vector128{}
		//V19 = &Vector128{}
		V20 = &Vector128{}
		//V21 = &Vector128{}
		//V22 = &Vector128{}
		//V23 = &Vector128{}
	)

	VLD1_2D([]uint64{0x0140014001400140, 0x0140014001400140}, V1)
	VLD1_2D([]uint64{0x0001100000011000, 0x0001100000011000}, V2)
	VLD1_2D([]uint64{0x090A040506000102, 0xFFFFFFFF0C0D0E08}, V3)
	//VLD1_2D([]uint64{0x0f0e0b0a07060302, 0x1f1e1b1a17161312}, V18) // high part of word
	//VLD1_2D([]uint64{0x0d0c090805040100, 0x1d1c191815141110}, V19) // low part of word
	//VLD1_2D([]uint64{0x0f0e0d0c07060504, 0x1f1e1d1c17161514}, V21) // high part of dword
	//VLD1_2D([]uint64{0x0b0a090803020100, 0x1b1a191813121110}, V22) // low part of dword

	VLD1_16B(lut[:], V8)
	VLD1_16B(lut[16:], V9)
	VLD1_16B(lut[32:], V10)
	VLD1_16B(lut[48:], V11)
	VLD1_16B(lut[64:], V12)
	VLD1_16B(lut[80:], V13)
	VLD1_16B(lut[96:], V14)
	VLD1_16B(lut[112:], V15)

	VDUP_BYTE(0x40, V7)

	// load input
	VLD1_16B(src, V0)

	// Get values from first LUT:
	VTBL_B(V0, []*Vector128{V8, V9, V10, V11}, V20)

	// Get values from second LUT:
	VSUB_B(V7, V0, V0)
	VTBX_B(V0, []*Vector128{V12, V13, V14, V15}, V20)

	// Check for invalid input, any value larger than 63:
	VCMHS_B(V7, V20, V16)

	// Check that all bits are zero:
	VUMAXV_B(true, V16, V17)
	if V17.bytes[15] != 0 {
		panic("invalid input")
	}

	// Compress four bytes into three:
	UMULL_B(V20, V1, V4)
	UMULL2_B(V20, V1, V0)
	VADDP_H(V0, V4, V0)
	//VTBL_B(V18, []*Vector128{V4, V0}, V5)
	//VTBL_B(V19, []*Vector128{V4, V0}, V6)
	//VADD_H(V5, V6, V0)

	UMULL_H(V0, V2, V4)
	UMULL2_H(V0, V2, V0)
	VADDP_S(V0, V4, V0)
	//VTBL_B(V21, []*Vector128{V4, V0}, V5)
	//VTBL_B(V22, []*Vector128{V4, V0}, V6)
	//VADD_S(V5, V6, V0)

	VTBL_B(V3, []*Vector128{V0}, V0)
	copy(dst, V0.Bytes()[:12])
}

func TestDecode16BSTD(t *testing.T) {
	cases := []struct {
		in       string
		expected []byte
	}{
		{"YWJjZGVmZ2hpamtsYWJjZGVmZ2hpamts", []byte("abcdefghijkl")},
	}
	for _, c := range cases {
		ret := make([]byte, 12)
		decode16B(ret, []byte(c.in), &dencodeStdLut)
		if !bytes.Equal(ret, c.expected) {
			t.Errorf("decode16B() = %x; want %x", ret, c.expected)
		}
	}
}

func TestVMUL(t *testing.T) {
	cases := []struct {
		size byte
		m    byte
		n    byte
		d    byte
		inst uint32
	}{
		{ // VMUL V0.B16, V1.B16, V2.B16
			0, 0, 1, 2, 0x4e209c22,
		},
		{ // VMUL V3.H8, V4.H8, V5.H8
			1, 3, 4, 5, 0x4e639c85,
		},
		{ // VMUL V0.H8, V7.H8, V0.H8
			1, 0, 7, 0, 0x4e609ce0,
		},
		{ // VMUL V6.S4, V7.S4, V8.S4
			2, 6, 7, 8, 0x4ea69ce8,
		},
	}
	for _, c := range cases {
		inst := uint32(0x4e209c00) | uint32(c.size&0x3)<<22 | uint32(c.d&0x1f) | uint32(c.n&0x1f)<<5 | (uint32(c.m&0x1f) << 16)
		if inst != c.inst {
			t.Errorf("VMUL = %x; want %x", inst, c.inst)
		}
	}
}

func TestUMULL(t *testing.T) {
	cases := []struct {
		q    byte
		size byte
		m    byte
		n    byte
		d    byte
		inst uint32
	}{
		{ // UMULL V0.B16, V1.B16, V4.H8
			0, 0, 0, 1, 4, 0x2e20c024,
		},
		{ // UMULL V0.H8, V2.H8, V4.S4
			0, 1, 0, 2, 4, 0x2e60c044,
		},
		{ // UMULL V6.S4, V7.S4, V8.D2
			0, 2, 6, 7, 8, 0x2ea6c0e8,
		},
		{ // UMULL V1.S4, V12.S4, V2.D2
			0, 1, 1, 12, 2, 0x2e61c182,
		},
		{ // UMULL2 V1.S4, V12.S4, V1.D2
			1, 1, 1, 12, 1, 0x6e61c181,
		},
		{ // UMULL2 V0.B16, V1.B16, V2.H8
			1, 0, 0, 1, 2, 0x6e20c022,
		},
		{ // UMULL2 V0.B16, V1.B16, V0.H8
			1, 0, 0, 1, 0, 0x6e20c020,
		},
		{ // UMULL2 V0.H8, V1.H8, V0.S4
			1, 1, 0, 1, 0, 0x6e60c020,
		},
		{ // UMULL2 V0.H8, V2.H8, V0.S4
			1, 1, 0, 2, 0, 0x6e60c040,
		},
		{ // UMULL2 V6.S4, V7.S4, V8.D2
			1, 2, 6, 7, 8, 0x6ea6c0e8,
		},
	}
	for _, c := range cases {
		inst := uint32(0x2e20c000) | uint32(c.q&1)<<30 | uint32(c.size&0x3)<<22 | uint32(c.d&0x1f) | uint32(c.n&0x1f)<<5 | (uint32(c.m&0x1f) << 16)
		if inst != c.inst {
			t.Errorf("UMULL = %x; want %x", inst, c.inst)
		}
	}
}

func vtbx(Vd, Vn, len, Vm byte) uint32 {
	inst := uint32(0x4e001000) | uint32(Vd&0x1f) | uint32(Vn&0x1f)<<5 | uint32(len&0x3)<<13 | uint32(Vm&0x1f)<<16
	return inst
}

func cmhi(Vd, Vn, Vm byte) uint32 {
	inst := uint32(0x6e203400) | uint32(Vd&0x1f) | uint32(Vn&0x1f)<<5 | uint32(Vm&0x1f)<<16
	return inst
}

func cmhs(Vd, Vn, Vm byte) uint32 {
	inst := uint32(0x6e203c00) | uint32(Vd&0x1f) | uint32(Vn&0x1f)<<5 | uint32(Vm&0x1f)<<16
	return inst
}

func TestInstructions(t *testing.T) {
	// VTBX V16.B16, [V12.B16, V13.B16, V14.B16, V15.B16], V16.B16
	inst := vtbx(16, 12, 3, 16)
	if inst != 0x4e107190 {
		t.Fatalf("got %x, expected 0x4e107190", inst)
	}
	// VTBX V20.B16, [V12.B16, V13.B16, V14.B16, V15.B16], V0.B16
	inst = vtbx(0, 12, 3, 20)
	if inst != 0x4e147180 {
		t.Fatalf("got %x, expected 0x4e147180", inst)
	}
	// VTBX V21.B16, [V12.B16, V13.B16, V14.B16, V15.B16], V1.B16
	inst = vtbx(1, 12, 3, 21)
	if inst != 0x4e157181 {
		t.Fatalf("got %x, expected 0x4e157181", inst)
	}
	// VTBX V22.B16, [V12.B16, V13.B16, V14.B16, V15.B16], V2.B16
	inst = vtbx(2, 12, 3, 22)
	if inst != 0x4e167182 {
		t.Fatalf("got %x, expected 0x4e167182", inst)
	}
	// VTBX V23.B16, [V12.B16, V13.B16, V14.B16, V15.B16], V3.B16
	inst = vtbx(3, 12, 3, 23)
	if inst != 0x4e177183 {
		t.Fatalf("got %x, expected 0x4e177183", inst)
	}
	// VCMHI V7.B16, V0.B16, V16.B16
	inst = cmhi(16, 0, 7)
	if inst != 0x6e273410 {
		t.Fatalf("got %x, expected 0x6e273410", inst)
	}
	// VCMHS V7.B16, V0.B16, V16.B16
	inst = cmhs(16, 0, 7)
	if inst != 0x6e273c10 {
		t.Fatalf("got %x, expected 0x6e273c10", inst)
	}
	// VCMHI V7.B16, V1.B16, V17.B16
	inst = cmhi(17, 1, 7)
	if inst != 0x6e273431 {
		t.Fatalf("got %x, expected 0x6e273431", inst)
	}
	// VCMHS V7.B16, V1.B16, V17.B16
	inst = cmhs(17, 1, 7)
	if inst != 0x6e273c31 {
		t.Fatalf("got %x, expected 0x6e273c31", inst)
	}
	// VCMHI V7.B16, V2.B16, V18.B16
	inst = cmhi(18, 2, 7)
	if inst != 0x6e273452 {
		t.Fatalf("got %x, expected 0x6e273452", inst)
	}
	// VCMHS V7.B16, V2.B16, V18.B16
	inst = cmhs(18, 2, 7)
	if inst != 0x6e273c52 {
		t.Fatalf("got %x, expected 0x6e273c52", inst)
	}
	// VCMHI V7.B16, V3.B16, V19.B16
	inst = cmhi(19, 3, 7)
	if inst != 0x6e273473 {
		t.Fatalf("got %x, expected 0x6e273473", inst)
	}
	// VCMHS V7.B16, V3.B16, V19.B16
	inst = cmhs(19, 3, 7)
	if inst != 0x6e273c73 {
		t.Fatalf("got %x, expected 0x6e273c73", inst)
	}
}
