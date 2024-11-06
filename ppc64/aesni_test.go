package ppc64

import (
	"testing"

	"github.com/emmansun/simd/alg/sm4"
	"github.com/emmansun/simd/alg/zuc"
)

func testSboxWithAESNI(t *testing.T, idx int, m1l, m1h, m2l, m2h *Vector128, sbox *[256]byte) {
	t.Helper()
	dst := &Vector128{}
	for k := 0; k < 256; k += 16 {
		LXVD2X([]byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)}, dst)
		SboxWithAESNI(m1l, m1h, m2l, m2h, dst)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != sbox[k+i] {
				t.Fatalf("Case %v SBOX(%v) = %x; want %x", idx, k+i, dst.Bytes()[i], sbox[k+i])
			}
		}
	}
}

func TestSM4SBOXWithAESNI(t *testing.T) {
	cases := []struct {
		m1, m2 uint64
		c1, c2 byte
	}{
		{
			0xa7ac65de3de94796,
			0xc101dd410ab464fa,
			0x69,
			0x61,
		},
		{
			0x34ac259e022dbc52,
			0x4e87acc7b40a9acb,
			0x65,
			0x2f,
		},
		{
			0x4c287db91a22505d,
			0x480e4c47651dbad3,
			0x3e,
			0x6c,
		},
		{
			0x242842865a99abe6,
			0xce81fbc81d658b2d,
			0x8e,
			0xe9,
		},
		{
			0xddec4505ceae37d1,
			0x336292532a5b1650,
			0x86,
			0x3c,
		},
		{
			0x8aec81c17591b3ee,
			0x0b95eaa45b2a5619,
			0xd6,
			0x4d,
		},
		{
			0xd517b18efe321f4d,
			0x6b0232fcc37428e8,
			0xce,
			0x81,
		},
		{
			0x06170a353a729b0d,
			0x9c3a8cc474c361a8,
			0x23,
			0x3b,
		},
		{
			0x669b0d608a162e14,
			0xe2acf9f7ddaf54fe,
			0x01,
			0x34,
		},
	}
	for i, c := range cases {
		m1l := &Vector128{}
		m1h := &Vector128{}
		m2l := &Vector128{}
		m2h := &Vector128{}
		GenLookupTable(c.m1, c.c1, m1l, m1h)
		GenLookupTable(c.m2, c.c2, m2l, m2h)
		testSboxWithAESNI(t, i+1, m1l, m1h, m2l, m2h, &sm4.SBOX)
	}
}

func TestZUCSBOXWithAESNI(t *testing.T) {
	cases := []struct {
		m1, m2 uint64
		c1, c2 byte
	}{
		{
			0xdd06c8f01eae7c70,
			0x0dedd9055ad8a502,
			0x00,
			0xfe,
		},
		{
			0x1106dce4d4485096,
			0x3c1a99b2ad1ed43a,
			0x00,
			0x32,
		},
		{
			0x2b3e78b290e2aa3c,
			0x6e7093a30891430e,
			0x00,
			0xbc,
		},
		{
			0xf33e408a76f65828,
			0x2ef66ddb8e57fd81,
			0x00,
			0xab,
		},
		{
			0x95124e5a9e18acc6,
			0x9636b3cc88265d01,
			0x00,
			0xd8,
		},
		{
			0xbf12a8bca6d25e0c,
			0xdfb9827207e02587,
			0x00,
			0x58,
		},
		{
			0xa7f40aec16b4426a,
			0xebbcaeeccda0f262,
			0x00,
			0xb7,
		},
		{
			0x47f4c026028c6e52,
			0x1584e79df5664595,
			0x00,
			0xec,
		},
	}
	for i, c := range cases {
		m1l := &Vector128{}
		m1h := &Vector128{}
		m2l := &Vector128{}
		m2h := &Vector128{}
		GenLookupTable(c.m1, c.c1, m1l, m1h)
		GenLookupTable(c.m2, c.c2, m2l, m2h)
		testSboxWithAESNI(t, i+1, m1l, m1h, m2l, m2h, &zuc.SBOX)
	}
}
