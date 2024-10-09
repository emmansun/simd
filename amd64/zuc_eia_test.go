package amd64

import (
	"encoding/hex"
	"testing"
)

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
		got := EIA16Bytes(data, tt.keys)
		if got != tt.want {
			t.Errorf("EIA16Bytes() = %v; want %v", got, tt.want)
		}
	}
}

var test64Vectors = []struct {
	dataHex string
	keys    []uint32
	want    uint64
}{
	{
		"11111111111111111111111111111111",
		[]uint32{0x3d9caf57, 0x8a89937c, 0x176b36fd, 0x11d75481, 0xfdfd3376, 0xf854d429, 0x44d97210, 0x59dbc2bb},
		0xd42c65f566559766,
	},
	{
		"983b41d47d780c9e1ad11d7eb70391b1",
		[]uint32{0x3d9caf57, 0x8a89937c, 0x176b36fd, 0x11d75481, 0xfdfd3376, 0xf854d429, 0x44d97210, 0x59dbc2bb},
		0x2e6a65f9d9cda19a,
	},
}

func TestEIA256_64(t *testing.T) {
	for _, tt := range test64Vectors {
		data, _ := hex.DecodeString(tt.dataHex)
		got := EIA256RoundTag8(data, tt.keys)
		if got != tt.want {
			t.Errorf("EIA256RoundTag8() = %x; want %x", got, tt.want)
		}
	}
}

var test128Vectors = []struct {
	dataHex string
	keys    []uint32
	want1    uint64
	want2    uint64
}{
	{
		"11111111111111111111111111111111",
		[]uint32{0x17c8fa3d, 0x4342534c , 0xca2c1aaf , 0xfe44033d , 0x8058b02 , 0xcda8ecbf , 0xc26c7761 , 0xf9fd0fc3},
		0x8f8d87816971a2b4,
		0x30d19f879dffca43,
	},
	{
		"983b41d47d780c9e1ad11d7eb70391b1",
		[]uint32{0x3d9caf57, 0x8a89937c, 0x176b36fd, 0x11d75481, 0xfdfd3376, 0xf854d429, 0x44d97210, 0x59dbc2bb},
		0x2e6a65f9d9cda19a,
		0xeafdf0e5fd66e534,
	},
}

func TestEIA256_128(t *testing.T) {
	for _, tt := range test128Vectors {
		data, _ := hex.DecodeString(tt.dataHex)
		got1, got2 := EIA256RoundTag16(data, tt.keys)
		if got1 != tt.want1 {
			t.Errorf("EIA256RoundTag16() = %x; want %x", got1, tt.want1)
		}
		if got2 != tt.want2 {
			t.Errorf("EIA256RoundTag16() = %x; want %x", got2, tt.want2)
		}
	}
}
