package ghash

import (
	"encoding/hex"
	"testing"
)

var mulCases = []struct {
	key string
	y   string
	out string
}{
	{
		key: "66e94bd4ef8a2c3b884cfa59ca342b2e",
		y:   "00000000000000000000000000000002",
		out: "a549b97029ca95c365a4805fb8dd7092",
	},
	{
		key: "66e94bd4ef8a2c3b884cfa59ca342b2e",
		y:   "00000000000000000000000000000003",
		out: "f7ed65c83d2fdf22d776c07064b3c8db",
	},
	{
		key: "66e94bd4ef8a2c3b884cfa59ca342b2e",
		y:   "0388dace60b6a392f328c2b971b2fe78",
		out: "5e2ec746917062882c85b0685353deb7",
	},
}

func TestReflect4Bits(t *testing.T) {
	for i := 0; i < 16; i++ {
		b := i & 0x0f
		b = b&1<<3 | b&2<<1 | b&4>>1 | b&8>>3
		if reflect4Bits(i) != b {
			t.Errorf("for %d got %d, want %d", i, reflect4Bits(i), b)
		}
	}
}

func TestReflect8Bits(t *testing.T) {
	for i := 0; i < 256; i++ {
		b := i&1<<7 | i&2<<5 | i&4<<3 | i&8<<1 | i&16>>1 | i&32>>3 | i&64>>5 | i&128>>7
		if reflect8Bits(i) != b {
			t.Errorf("for %d got %d, want %d", i, reflect8Bits(i), b)
		}
	}
}

func TestRawMethodMul(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		var m rawMethod
		copy(m.key[:], key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.Mul(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestSimpleMethod8Bits(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewSimpleMethod8Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.Mul(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestSimpleMethod4Bits(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewSimpleMethod4Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.Mul(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestShoupMethod8BitsMul(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewShoupMethod8Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.Mul(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestShoupMethod8BitsMulImpl2(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewShoupMethod8Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.MulImpl2(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestShoupMethod4BitsMul(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewShoupMethod4Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.Mul(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestShoupMethod4BitsMulImpl2(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewShoupMethod4Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.MulImpl2(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestShoupMethod4BitsMulImpl3(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewShoupMethod4Bits(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.MulImpl3(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func TestGCMMethodMul(t *testing.T) {
	for i, c := range mulCases {
		key, _ := hex.DecodeString(c.key)
		m := NewGCMMethod(key)
		var y [16]byte
		y1, _ := hex.DecodeString(c.y)
		copy(y[:], y1)
		m.Mul(&y)
		if hex.EncodeToString(y[:]) != c.out {
			t.Errorf("case %d: got %v, want %v", i, hex.EncodeToString(y[:]), c.out)
		}
	}
}

func BenchmarkRawMethodMul(b *testing.B) {
	key, _ := hex.DecodeString("66e94bd4ef8a2c3b884cfa59ca342b2e")
	var m rawMethod
	copy(m.key[:], key)
	var y [16]byte
	y1, _ := hex.DecodeString("0388dace60b6a392f328c2b971b2fe78")
	copy(y[:], y1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(&y)
	}
}

func BenchmarkSimpleMethod8Bits(b *testing.B) {
	key, _ := hex.DecodeString("66e94bd4ef8a2c3b884cfa59ca342b2e")
	m := NewSimpleMethod8Bits(key)
	var y [16]byte
	y1, _ := hex.DecodeString("0388dace60b6a392f328c2b971b2fe78")
	copy(y[:], y1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(&y)
	}
}

func BenchmarkSimpleMethod4Bits(b *testing.B) {
	key, _ := hex.DecodeString("66e94bd4ef8a2c3b884cfa59ca342b2e")
	m := NewSimpleMethod4Bits(key)
	var y [16]byte
	y1, _ := hex.DecodeString("0388dace60b6a392f328c2b971b2fe78")
	copy(y[:], y1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(&y)
	}
}

func BenchmarkShoupMethod8Bits(b *testing.B) {
	key, _ := hex.DecodeString("66e94bd4ef8a2c3b884cfa59ca342b2e")
	m := NewShoupMethod8Bits(key)
	var y [16]byte
	y1, _ := hex.DecodeString("0388dace60b6a392f328c2b971b2fe78")
	copy(y[:], y1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(&y)
	}
}

func BenchmarkShoupMethod4Bits(b *testing.B) {
	key, _ := hex.DecodeString("66e94bd4ef8a2c3b884cfa59ca342b2e")
	m := NewShoupMethod4Bits(key)
	var y [16]byte
	y1, _ := hex.DecodeString("0388dace60b6a392f328c2b971b2fe78")
	copy(y[:], y1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(&y)
	}
}

func BenchmarkGCMMethodMul(b *testing.B) {
	key, _ := hex.DecodeString("66e94bd4ef8a2c3b884cfa59ca342b2e")
	m := NewGCMMethod(key)
	var y [16]byte
	y1, _ := hex.DecodeString("0388dace60b6a392f328c2b971b2fe78")
	copy(y[:], y1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(&y)
	}
}
