package arm64

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
