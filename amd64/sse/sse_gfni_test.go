package sse

import "testing"

func TestAESSBOXWithGFNI(t *testing.T) {
	m2 := Set64(0xf1e3c78f1f3e7cf8, 0xf1e3c78f1f3e7cf8)
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		// AES SBOX
		GF2P8AFFINEINVQB(dst, &m2, 0x63)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != aes_sbox[k+i] {
				t.Fatalf("AES SBOX(%v) = %x; want %x", k+i, dst.Bytes()[i], aes_sbox[k+i])
			}
		}
	}
}

func TestAESSBOXWithGFNI2(t *testing.T) {
	m1 := Set64(0x0102040810204080, 0x0102040810204080)
	m2 := Set64(0xf1e3c78f1f3e7cf8, 0xf1e3c78f1f3e7cf8)
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		// AES SBOX
		GF2P8AFFINEQB(dst, &m1, 0x00)
		GF2P8AFFINEINVQB(dst, &m2, 0x63)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != aes_sbox[k+i] {
				t.Fatalf("AES SBOX(%v) = %x; want %x", k+i, dst.Bytes()[i], aes_sbox[k+i])
			}
		}
	}
}

func testSBOXWithGFNI(t *testing.T, idx int, m1, m2 *XMM, c1, c2 byte, sbox *[256]byte) {
	t.Helper()
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		// SBOX
		GF2P8AFFINEQB(dst, m1, c1)
		GF2P8AFFINEINVQB(dst, m2, c2)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != sbox[k+i] {
				t.Fatalf("Case %v SBOX(%v) = %x; want %x", idx, k+i, dst.Bytes()[i], sbox[k+i])
			}
		}
	}
}

func TestSM4SBOXWithGFNI(t *testing.T) {
	cases := []struct {
		m1, m2 uint64
		c1     byte
	}{
		{
			0xa7ac65de3de94796,
			0x75f1228d6c1e85c9,
			0x69,
		},
		{
			0x34ac259e022dbc52,
			0xd72d8e511e6c8b19,
			0x65,
		},
		{
			0x4c287db91a22505d,
			0xf3ab34a974a6b589,
			0x3e,
		},
		{
			0x242842865a99abe6,
			0x2f09380ba6746587,
			0x8e,
		},
		{
			0xddec4505ceae37d1,
			0x33a1047152fe3b63,
			0x86,
		},
		{
			0x8aec81c17591b3ee,
			0x9dd1d601fe524761,
			0xd6,
		},
		{
			0xd517b18efe321f4d,
			0xdfe3c2ed969ab135,
			0xce,
		},
		{
			0x06170a353a729b0d,
			0xaf4db0439a96b349,
			0x23,
		},
	}
	for i, c := range cases {
		m1 := Set64(c.m1, c.m1)
		m2 := Set64(c.m2, c.m2)
		testSBOXWithGFNI(t, i+1, &m1, &m2, c.c1, 0xd3, &sm4_sbox)
	}
}

func TestZUCSBOXWithGFNI(t *testing.T) {
	cases := []struct {
		m1, m2 uint64
	}{
		{
			0xF33E408A76F65828,
			0x9581FB0653B61C09,
		},
		{
			0x2B3E78B290E2AA3C,
			0xE95DF5D48F166EAB,
		},
		{
			0x95124E5A9E18ACC6,
			0xC305CBCC771ADAF1,
		},
		{
			0xBF12A8BCA6D25E0C,
			0xC1A71BBED5BA082D,
		},
		{
			0xDD06C8F01EAE7C70,
			0xB903E5360F14F0E3,
		},
		{
			0x1106DCE4D4485096,
			0x6973993A7FB45C4D,
		},
		{
			0x47F4C026028C6E52,
			0x293F6F5E93664AD1,
		},
		{
			0xA7F40AEC16B4426A,
			0x27916DF23DC646A1,
		},
	}
	for i, c := range cases {
		m1 := Set64(c.m1, c.m1)
		m2 := Set64(c.m2, c.m2)
		testSBOXWithGFNI(t, i+1, &m1, &m2, 0, 0x55, &zuc_sbox)
	}
}
