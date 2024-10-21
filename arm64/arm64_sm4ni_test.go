package arm64

import (
	"bytes"
	"testing"
)

func TestExpandKeySM4(t *testing.T) {
	enc := make([]uint32, 32)
	expected := []uint32{0xf12186f9, 0x41662b61, 0x5a6ab19a, 0x7ba92077, 0x367360f4, 0x776a0c61, 0xb6bb89b3, 0x24763151, 0xa520307c, 0xb7584dbd, 0xc30753ed, 0x7ee55b57, 0x6988608c, 0x30d895b7, 0x44ba14af, 0x104495a1, 0xd120b428, 0x73b55fa3, 0xcc874966, 0x92244439, 0xe89e641f, 0x98ca015a, 0xc7159060, 0x99e1fd2e, 0xb79bd80c, 0x1d2115b0, 0xe228aeb, 0xf1780c81, 0x428d3654, 0x62293496, 0x1cf72e5, 0x9124a012}
	ck := &Vector128{}
	fk := &Vector128{}
	key := &Vector128{}
	VLD1_16B([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}, key)
	VREV32_B(key, key)
	VLD1_2D([]uint64{0x56aa3350a3b1bac6, 0xb27022dc677d9197}, fk)
	VEOR(fk, key, key)

	for i := 0; i < 32; i += 4 {
		VLD1_4S(sm4_ck[i:], ck)
		SM4EKEY(ck, key, key)
		VST1_4S(key, enc[i:])
	}
	for i := 0; i < 32; i++ {
		if expected[i] != enc[i] {
			t.Errorf("expected[%d] = %x; got %x", i, expected[i], enc[i])
		}
	}
}

func TestSM4(t *testing.T) {
	src := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
	encKey := []uint32{0xf12186f9, 0x41662b61, 0x5a6ab19a, 0x7ba92077, 0x367360f4, 0x776a0c61, 0xb6bb89b3, 0x24763151, 0xa520307c, 0xb7584dbd, 0xc30753ed, 0x7ee55b57, 0x6988608c, 0x30d895b7, 0x44ba14af, 0x104495a1, 0xd120b428, 0x73b55fa3, 0xcc874966, 0x92244439, 0xe89e641f, 0x98ca015a, 0xc7159060, 0x99e1fd2e, 0xb79bd80c, 0x1d2115b0, 0xe228aeb, 0xf1780c81, 0x428d3654, 0x62293496, 0x1cf72e5, 0x9124a012}
	expected := []byte{0x68, 0x1e, 0xdf, 0x34, 0xd2, 0x06, 0x96, 0x5e, 0x86, 0xb3, 0xe9, 0x4f, 0x53, 0x6e, 0x42, 0x46}
	data := &Vector128{}
	rk := &Vector128{}
	VLD1_16B(src, data)
	VREV32_B(data, data)
	for i := 0; i < 32; i += 4 {
		VLD1_4S(encKey[i:], rk)
		SM4E(rk, data)
	}
	VREV64_B(data, data)
	VEXT(8, data, data, data)
	result := make([]byte, 16)
	VST1_16B(data, result)
	if !bytes.Equal(result, expected) {
		t.Errorf("expected = %x; got %x", expected, result)
	}
	// decrypt
	decKey := make([]uint32, 32)
	for i := 0; i < 32; i++ {
		decKey[i] = encKey[31-i]
	}
	VLD1_16B(expected, data)
	VREV32_B(data, data)
	for i := 0; i < 32; i += 4 {
		VLD1_4S(decKey[i:], rk)
		SM4E(rk, data)
	}
	VREV64_B(data, data)
	VEXT(8, data, data, data)
	VST1_16B(data, result)
	if !bytes.Equal(result, src) {
		t.Errorf("expected = %x; got %x", src, result)
	}
}
