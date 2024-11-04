package amd64

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/emmansun/simd/alg/ghash"
)

var ghashCases = []struct {
	key  string
	data string
}{
	{
		"66e94bd4ef8a2c3b884cfa59ca342b2e",
		"48af93501fa62adbcd414cce6034d8",
	},
	{
		"66e94bd4ef8a2c3b884cfa59ca342b2e",
		"48af93501fa62adbcd414cce6034d89587",
	},
	{
		"66e94bd4ef8a2c3b884cfa59ca342b2e",
		"48af93501fa62adbcd414cce6034d895dda1bf8f132f042098661572e7483094fd12e518ce062c98acee28d95df4416bed31a2f04476c18bb40c84a74b97dc5b16842d4fa186f56ab33256971fa110f4",
	},
	{
		"66e94bd4ef8a2c3b884cfa59ca342b2e",
		"48af93501fa62adbcd414cce6034d895dda1bf8f132f042098661572e7483094fd12e518ce062c98acee28d95df4416bed31a2f04476c18bb40c84a74b97dc5b16842d4fa186f56ab33256971fa110f448af93501fa62adbcd414cce6034d895dda1bf8f132f042098661572e7483094fd12e518ce062c98acee28d95df4416bed31a2f04476c18bb40c84a74b97dc5b16842d4fa186f56ab33256971fa110f4abcd",
	},
}

func TestAMD64GHash(t *testing.T) {
	for i, c := range ghashCases {
		key, _ := hex.DecodeString(c.key)
		g1 := NewClmulAMD64Ghash(key)
		g2 := ghash.NewGCMMethod(key)
		var T1, T2 [16]byte
		data, _ := hex.DecodeString(c.data)
		g1.Hash(&T1, data)
		g2.Hash(&T2, data)
		if !bytes.Equal(T1[:], T2[:]) {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T1[:]), hex.EncodeToString(T2[:]))
		}
	}
}
