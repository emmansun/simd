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
	0, 255, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 255, 255, 255, 255,
	255, 255, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 255, 255, 255, 255,
}

var dencodeUrlLut = [128]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 62, 255, 255,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 255, 255, 255, 255, 255, 255,
	0, 255, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
	14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 255, 255, 255, 255,
	63, 255, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 255, 255, 255, 255,
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

func TestEncodeSTD(t *testing.T) {
	cases := []struct {
		in        []byte
		expected  string
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
		in        []byte
		expected  string
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

	VDUP_BYTE(0x3f, V7)

	// load input
	VLD4_16B(src, V0, V1, V2, V3)

	// Get indices for second LUT:
	VUQSUB_B(V7, V0, V16)
	VUQSUB_B(V7, V1, V17)
	VUQSUB_B(V7, V2, V18)
	VUQSUB_B(V7, V3, V19)

	// Get values from first LUT:
	VTBL_B(V0, []*Vector128{V8, V9, V10, V11}, V20)
	VTBL_B(V1, []*Vector128{V8, V9, V10, V11}, V21)
	VTBL_B(V2, []*Vector128{V8, V9, V10, V11}, V22)
	VTBL_B(V3, []*Vector128{V8, V9, V10, V11}, V23)

	// Get values from second LUT:
	VTBX_B(V16, []*Vector128{V12, V13, V14, V15}, V16)
	VTBX_B(V17, []*Vector128{V12, V13, V14, V15}, V17)
	VTBX_B(V18, []*Vector128{V12, V13, V14, V15}, V18)
	VTBX_B(V19, []*Vector128{V12, V13, V14, V15}, V19)

	// Get final values:
	VORR(V20, V16, V0)
	VORR(V21, V17, V1)
	VORR(V22, V18, V2)
	VORR(V23, V19, V3)

	// Check for invalid input, any value larger than 63:
	VCMHI_B(false, V7, V0, V16)
	VCMHI_B(false, V7, V1, V17)
	VCMHI_B(false, V7, V2, V18)
	VCMHI_B(false, V7, V3, V19)

	VORR(V16, V17, V16)
	VORR(V18, V19, V18)
	VORR(V16, V18, V16)

	// Check that all bits are zero:
	VUMAXV_B(true, V16, V17)
	if V17.bytes[15] != 0 {
		panic("invalid input")
	}

	// Compress four bytes into three:
	VSHL_B(2, V0, V4)
	VUSHR_B(4, V1, V16)
	VORR(V4, V16, V4)

	VSHL_B(4, V1, V5)
	VUSHR_B(2, V2, V16)
	VORR(V5, V16, V5)

	VSHL_B(6, V2, V16)
	VORR(V16, V3, V6)

	VST3_16B(V4, V5, V6, dst)
}

func TestDecodeSTD(t *testing.T) {
	cases := []struct {
		in        string
		expected  []byte
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
		in        string
		expected  []byte
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
