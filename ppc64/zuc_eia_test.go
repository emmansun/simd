package ppc64

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func eia16Bytes(data []byte, keys []uint32) uint32 {
	var (
		XTMP1         Vector128
		XTMP2         Vector128
		XTMP3         Vector128
		XTMP4         Vector128
		XDATA         Vector128
		XDIGEST       Vector128
		KS_L          Vector128
		KS_M1         Vector128
		BIT_REV_TAB_L Vector128
		BIT_REV_TAB_H Vector128
	)
	LXVD2X_UINT64([]uint64{0x0008040c020a060e, 0x0109050d030b070f}, &BIT_REV_TAB_L)
	VSPLTISB(4, &XTMP2)
	VSLB(&BIT_REV_TAB_L, &XTMP2, &BIT_REV_TAB_H)
	//LXVD2X_UINT64([]uint64{0x008040c020a060e0, 0x109050d030b070f0}, &BIT_REV_TAB_H)

	LXVD2X(data, &XDATA)
	// Change byte order (for PPC64)
	LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, &XTMP1)
	VPERM(&XDATA, &XDATA, &XTMP1, &XDATA)

	VPERMXOR(&BIT_REV_TAB_L, &BIT_REV_TAB_H, &XDATA, &XTMP3)

	// ZUC authentication part, 4x32 data bits
	// setup data
	VSPLTISW(0, &XTMP2)
	LXVD2X_UINT64([]uint64{0x0000000010111213, 0x0000000014151617}, &XTMP4)
	VPERM(&XTMP2, &XTMP3, &XTMP4, &XTMP1)
	LXVD2X_UINT64([]uint64{0x0000000018191a1b, 0x000000001c1d1e1f}, &XTMP4)
	VPERM(&XTMP2, &XTMP3, &XTMP4, &XTMP2)

	fmt.Printf("XTMP1: %x\n", XTMP1.Uint32s())
	fmt.Printf("XTMP2: %x\n", XTMP2.Uint32s())

	// setup KS
	LXVW4X_UINT32(keys, &KS_L)
	LXVD2X_UINT64([]uint64{0x0405060708090a0b, 0x0001020304050607}, &XTMP4)
	VPERM(&KS_L, &KS_L, &XTMP4, &KS_L)
	LXVW4X_UINT32(keys[2:], &KS_M1)
	VPERM(&KS_M1, &KS_M1, &XTMP4, &KS_M1)
	fmt.Printf("KS_L: %x\n", KS_L.Uint32s())
	fmt.Printf("KS_M1: %x\n", KS_M1.Uint32s())

	// clmul
	// xor the results from 4 32-bit words together
	// Calculate lower 32 bits of tag
	VPMSUMD(&XTMP1, &KS_L, &XTMP3)
	VPMSUMD(&XTMP2, &KS_M1, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XTMP3)
	VSPLTW(2, &XTMP3, &XDIGEST)

	fmt.Printf("XDIGEST: %x\n", XDIGEST.Uint32s())
	// use MFVSRWZ to get the lower 32 bits
	return XDIGEST.Uint32s()[3]
}

var testVectors = []struct {
	dataHex string
	keys    []uint32
	want    uint32
}{
	{

		"983b41d47d780c9e1ad11d7eb70391b1",
		[]uint32{0xa10eb178, 0xd2758cfc, 0x7b86b39d, 0x1ef5b475, 0x1902e017, 0x9820fb9c, 0xac9485e2, 0x1072e635},
		0xe27354df,
	},
	{
		"de0b35da2dc62f83e7b78d6306ca0ea0",
		[]uint32{0x1902e017, 0x9820fb9c, 0xac9485e2, 0x1072e635, 0xda0126c1, 0xb2168f8c, 0x4be50389, 0x185ce9fa},
		0x985490ae,
	},
	{
		"7e941b7be91348f9fcb170e2217fecd9",
		[]uint32{0xda0126c1, 0xb2168f8c, 0x4be50389, 0x185ce9fa, 0xa47d64c6, 0x28d03e82, 0xb8505ba7, 0x217a99b1},
		0x7387e168,
	},
	{
		"7f9f68adb16e5d7d21e569d280ed775c",
		[]uint32{0xa47d64c6, 0x28d03e82, 0xb8505ba7, 0x217a99b1, 0xc2fb807, 0x5bbbc219, 0x17f1a3fa, 0x4cd31ce0},
		0x2e9ed291,
	},
}

func TestEIA16Bytes(t *testing.T) {
	for _, tt := range testVectors {
		data, _ := hex.DecodeString(tt.dataHex)
		got := eia16Bytes(data, tt.keys)
		if got != tt.want {
			t.Errorf("EIA16Bytes() = %x; want %x", got, tt.want)
		}
	}
}

func eia256RoundTag8(data []byte, keys []uint32) uint64 {
	var (
		XTMP1         Vector128
		XTMP2         Vector128
		XTMP3         Vector128
		XTMP4         Vector128
		XTMP5         Vector128
		XTMP6         Vector128
		XDATA         Vector128
		XDIGEST       Vector128
		ZERO          Vector128
		KS_L          Vector128
		KS_M1         Vector128
		KS_M2         Vector128
		BIT_REV_TAB_L Vector128
		BIT_REV_TAB_H Vector128
	)
	LXVD2X_UINT64([]uint64{0x0008040c020a060e, 0x0109050d030b070f}, &BIT_REV_TAB_L)
	VSPLTISB(4, &XTMP2)
	VSLB(&BIT_REV_TAB_L, &XTMP2, &BIT_REV_TAB_H)
	//LXVD2X_UINT64([]uint64{0x008040c020a060e0, 0x109050d030b070f0}, &BIT_REV_TAB_H)

	LXVD2X(data, &XDATA)
	// Change byte order (for PPC64)
	LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, &XTMP1)
	VPERM(&XDATA, &XDATA, &XTMP1, &XDATA)

	VPERMXOR(&BIT_REV_TAB_L, &BIT_REV_TAB_H, &XDATA, &XTMP3)

	// ZUC authentication part, 4x32 data bits
	// setup data
	VSPLTISW(0, &ZERO)
	LXVD2X_UINT64([]uint64{0x0000000010111213, 0x0000000014151617}, &XTMP4)
	VPERM(&ZERO, &XTMP3, &XTMP4, &XTMP1)
	LXVD2X_UINT64([]uint64{0x0000000018191a1b, 0x000000001c1d1e1f}, &XTMP4)
	VPERM(&ZERO, &XTMP3, &XTMP4, &XTMP2)

	VOR(&XTMP1, &XTMP1, &XTMP5)
	VOR(&XTMP2, &XTMP2, &XTMP6)

	// setup KS
	LXVW4X_UINT32(keys, &KS_L)
	LXVD2X_UINT64([]uint64{0x0405060708090a0b, 0x0001020304050607}, &XTMP4)
	VPERM(&KS_L, &KS_L, &XTMP4, &KS_L)
	LXVW4X_UINT32(keys[2:], &KS_M1)
	VPERM(&KS_M1, &KS_M1, &XTMP4, &KS_M1)
	LXVW4X_UINT32(keys[4:], &KS_M2)
	VPERM(&KS_M2, &KS_M2, &XTMP4, &KS_M2)

	// clmul
	// xor the results from 4 32-bit words together
	// Calculate lower 32 bits of tag
	VPMSUMD(&XTMP1, &KS_L, &XTMP3)
	VPMSUMD(&XTMP2, &KS_M1, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XDIGEST)
	VSLDOI(12, &XDIGEST, &XDIGEST, &XDIGEST)

	// Calculate upper 32 bits of tag
	VOR(&XTMP5, &XTMP5, &XTMP1)
	VOR(&XTMP6, &XTMP6, &XTMP2)

	VSLDOI(8, &KS_M1, &KS_L, &KS_L)
	VPMSUMD(&XTMP1, &KS_L, &XTMP3)
	VSLDOI(8, &KS_M2, &KS_M1, &KS_M1)
	VPMSUMD(&XTMP2, &KS_M1, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XTMP3)
	VSLDOI(8, &XTMP3, &XTMP3, &XTMP3)
	VSLDOI(4, &XDIGEST, &XTMP3, &XDIGEST)

	return XDIGEST.Uint64s()[1]
}

var test64Vectors = []struct {
	dataHex string
	keys    []uint32
	want    uint64
}{
	{
		"11111111111111111111111111111111",
		[]uint32{0x3d9caf57, 0x8a89937c, 0x176b36fd, 0x11d75481, 0xfdfd3376, 0xf854d429, 0x44d97210, 0x59dbc2bb},
		0x66559766d42c65f5,
	},
	{
		"983b41d47d780c9e1ad11d7eb70391b1",
		[]uint32{0x3d9caf57, 0x8a89937c, 0x176b36fd, 0x11d75481, 0xfdfd3376, 0xf854d429, 0x44d97210, 0x59dbc2bb},
		0xd9cda19a2e6a65f9,
	},
}

func TestEIA256_64(t *testing.T) {
	for _, tt := range test64Vectors {
		data, _ := hex.DecodeString(tt.dataHex)
		got := eia256RoundTag8(data, tt.keys)
		if got != tt.want {
			t.Errorf("eia256RoundTag8() = %x; want %x", got, tt.want)
		}
	}
}

func eia256RoundTag16(data []byte, keys []uint32) (uint64, uint64) {
	var (
		XTMP1         Vector128
		XTMP2         Vector128
		XTMP3         Vector128
		XTMP4         Vector128
		XTMP5         Vector128
		XTMP6         Vector128
		XDATA         Vector128
		XDIGEST       Vector128
		ZERO          Vector128
		KS_L          Vector128
		KS_M1         Vector128
		KS_M2         Vector128
		KS_H          Vector128
		BIT_REV_TAB_L Vector128
		BIT_REV_TAB_H Vector128
	)
	LXVD2X_UINT64([]uint64{0x0008040c020a060e, 0x0109050d030b070f}, &BIT_REV_TAB_L)
	VSPLTISB(4, &XTMP2)
	VSLB(&BIT_REV_TAB_L, &XTMP2, &BIT_REV_TAB_H)
	//LXVD2X_UINT64([]uint64{0x008040c020a060e0, 0x109050d030b070f0}, &BIT_REV_TAB_H)

	LXVD2X(data, &XDATA)
	// Change byte order (for PPC64)
	LXVD2X_UINT64([]uint64{0x0706050403020100, 0x0f0e0d0c0b0a0908}, &XTMP1)
	VPERM(&XDATA, &XDATA, &XTMP1, &XDATA)

	VPERMXOR(&BIT_REV_TAB_L, &BIT_REV_TAB_H, &XDATA, &XTMP3)

	// ZUC authentication part, 4x32 data bits
	// setup data
	VSPLTISW(0, &ZERO)
	LXVD2X_UINT64([]uint64{0x0000000010111213, 0x0000000014151617}, &XTMP4)
	VPERM(&ZERO, &XTMP3, &XTMP4, &XTMP1)
	LXVD2X_UINT64([]uint64{0x0000000018191a1b, 0x000000001c1d1e1f}, &XTMP4)
	VPERM(&ZERO, &XTMP3, &XTMP4, &XTMP2)

	VOR(&XTMP1, &XTMP1, &XTMP5)
	VOR(&XTMP2, &XTMP2, &XTMP6)

	// setup KS
	LXVW4X_UINT32(keys, &KS_L)
	LXVD2X_UINT64([]uint64{0x0405060708090a0b, 0x0001020304050607}, &XTMP4)
	VPERM(&KS_L, &KS_L, &XTMP4, &KS_L)
	LXVW4X_UINT32(keys[2:], &KS_M1)
	VPERM(&KS_M1, &KS_M1, &XTMP4, &KS_M1)
	LXVW4X_UINT32(keys[4:], &KS_M2)
	VOR(&KS_M2, &KS_M2, &KS_H)
	//VSLDOI(8, &KS_H, &KS_H, &KS_H)
	//VSLDOI(8, &ZERO, &KS_H, &KS_H)
	VPERM(&KS_M2, &KS_M2, &XTMP4, &KS_M2)
	fmt.Printf("KS_L: %x\n", KS_L.Uint32s())
	fmt.Printf("KS_M1: %x\n", KS_M1.Uint32s())
	fmt.Printf("KS_M2: %x\n", KS_M2.Uint32s())
	fmt.Printf("KS_H: %x\n", KS_H.Uint32s())
	// clmul
	// xor the results from 4 32-bit words together
	// Calculate lower 32 bits of tag
	VPMSUMD(&XTMP1, &KS_L, &XTMP3)
	VPMSUMD(&XTMP2, &KS_M1, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XDIGEST)
	VSLDOI(12, &XDIGEST, &XDIGEST, &XDIGEST)

	// Calculate upper 32 bits of tag
	VOR(&XTMP5, &XTMP5, &XTMP1)
	VOR(&XTMP6, &XTMP6, &XTMP2)

	VSLDOI(8, &KS_M1, &KS_L, &KS_L)
	VPMSUMD(&XTMP1, &KS_L, &XTMP3)
	VSLDOI(8, &KS_M2, &KS_M1, &XTMP1)
	VPMSUMD(&XTMP2, &XTMP1, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XTMP3)

	VSLDOI(8, &XTMP3, &XTMP3, &XTMP3)
	VSLDOI(4, &XDIGEST, &XTMP3, &XDIGEST)

	// Prepare data and calculate bits 95-64 of tag
	VOR(&XTMP5, &XTMP5, &XTMP1)
	VOR(&XTMP6, &XTMP6, &XTMP2)
	VPMSUMD(&XTMP1, &KS_M1, &XTMP3)
	VPMSUMD(&XTMP2, &KS_M2, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XTMP3)
	VSLDOI(8, &XTMP3, &XTMP3, &XTMP3)
	VSLDOI(4, &XDIGEST, &XTMP3, &XDIGEST)

	// Prepare data and calculate bits 127-96 of tag
	VOR(&XTMP5, &XTMP5, &XTMP1)
	VOR(&XTMP6, &XTMP6, &XTMP2)
	VSLDOI(8, &KS_M2, &KS_M1, &KS_M1)
	VPMSUMD(&XTMP1, &KS_M1, &XTMP3)
	VSLDOI(8, &KS_H, &KS_M2, &KS_M2)
	VPMSUMD(&XTMP2, &KS_M2, &XTMP4)
	VXOR(&XTMP3, &XTMP4, &XTMP3)
	VSLDOI(8, &XTMP3, &XTMP3, &XTMP3)
	VSLDOI(4, &XDIGEST, &XTMP3, &XDIGEST)

	return XDIGEST.Uint64s()[0], XDIGEST.Uint64s()[1]
}

var test128Vectors = []struct {
	dataHex string
	keys    []uint32
	want1   uint64
	want2   uint64
}{
	{
		"11111111111111111111111111111111",
		[]uint32{0x17c8fa3d, 0x4342534c, 0xca2c1aaf, 0xfe44033d, 0x8058b02, 0xcda8ecbf, 0xc26c7761, 0xf9fd0fc3},
		0x6971a2b48f8d8781,
		0x9dffca4330d19f87,
	},
	{
		"983b41d47d780c9e1ad11d7eb70391b1",
		[]uint32{0x3d9caf57, 0x8a89937c, 0x176b36fd, 0x11d75481, 0xfdfd3376, 0xf854d429, 0x44d97210, 0x59dbc2bb},
		0xd9cda19a2e6a65f9,
		0xfd66e534eafdf0e5,
	},
}

func TestEIA256_128(t *testing.T) {
	for _, tt := range test128Vectors {
		data, _ := hex.DecodeString(tt.dataHex)
		got1, got2 := eia256RoundTag16(data, tt.keys)
		if got1 != tt.want1 {
			t.Errorf("eia256RoundTag16() = %x; want %x", got1, tt.want1)
		}
		if got2 != tt.want2 {
			t.Errorf("eia256RoundTag16() = %x; want %x", got2, tt.want2)
		}
	}
}
