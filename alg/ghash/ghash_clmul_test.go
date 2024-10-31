package ghash

import (
	"bytes"
	"encoding/hex"
	"testing"
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
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewClmulAMD64Ghash(key)
		var T [16]byte
		data, _ := hex.DecodeString(c.y)
		m.Hash(&T, data)
		if hex.EncodeToString(T[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T[:]), c.out)
		}
	}

	for i, c := range ghashCases {
		key, _ := hex.DecodeString(c.key)
		g1 := NewClmulAMD64Ghash(key)
		g2 := NewGCMMethod(key)
		var T1, T2 [16]byte
		data, _ := hex.DecodeString(c.data)
		g1.Hash(&T1, data)
		g2.Hash(&T2, data)
		if !bytes.Equal(T1[:], T2[:]) {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T1[:]), hex.EncodeToString(T2[:]))
		}
	}
}

func TestARM64GHash(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewClmulARM64Ghash(key)
		var T [16]byte
		data, _ := hex.DecodeString(c.y)
		m.Hash(&T, data)
		if hex.EncodeToString(T[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T[:]), c.out)
		}
	}

	for i, c := range ghashCases {
		key, _ := hex.DecodeString(c.key)
		g1 := NewClmulARM64Ghash(key)
		g2 := NewGCMMethod(key)
		var T1, T2 [16]byte
		data, _ := hex.DecodeString(c.data)
		g1.Hash(&T1, data)
		g2.Hash(&T2, data)
		if !bytes.Equal(T1[:], T2[:]) {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T1[:]), hex.EncodeToString(T2[:]))
		}
	}
}

func TestPPC64GHash(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewClmulPPC64Ghash(key, false)
		var T [16]byte
		data, _ := hex.DecodeString(c.y)
		m.Hash(&T, data)
		if hex.EncodeToString(T[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T[:]), c.out)
		}
	}

	for i, c := range ghashCases {
		key, _ := hex.DecodeString(c.key)
		g1 := NewClmulPPC64Ghash(key, false)
		g2 := NewGCMMethod(key)
		var T1, T2 [16]byte
		data, _ := hex.DecodeString(c.data)
		g1.Hash(&T1, data)
		g2.Hash(&T2, data)
		if !bytes.Equal(T1[:], T2[:]) {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T1[:]), hex.EncodeToString(T2[:]))
		}
	}
}

func TestPPC64LEGHash(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewClmulPPC64Ghash(key, true)
		var T [16]byte
		data, _ := hex.DecodeString(c.y)
		m.Hash(&T, data)
		if hex.EncodeToString(T[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T[:]), c.out)
		}
	}

	for i, c := range ghashCases {
		key, _ := hex.DecodeString(c.key)
		g1 := NewClmulPPC64Ghash(key, true)
		g2 := NewGCMMethod(key)
		var T1, T2 [16]byte
		data, _ := hex.DecodeString(c.data)
		g1.Hash(&T1, data)
		g2.Hash(&T2, data)
		if !bytes.Equal(T1[:], T2[:]) {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(T1[:]), hex.EncodeToString(T2[:]))
		}
	}
}
