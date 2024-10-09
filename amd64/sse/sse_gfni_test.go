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

func TestSM4SBOXWithGFNI(t *testing.T) {
	m1 := Set64(0xa7ac65de3de94796, 0xa7ac65de3de94796)
	m2 := Set64(0x75f1228d6c1e85c9, 0x75f1228d6c1e85c9)
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		// SM4 SBOX
		GF2P8AFFINEQB(dst, &m1, 0x69)
		GF2P8AFFINEINVQB(dst, &m2, 0xd3)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != sm4_sbox[k+i] {
				t.Fatalf("SM4 SBOX(%v) = %x; want %x", k+i, dst.Bytes()[i], sm4_sbox[k+i])
			}
		}
	}
}

func TestZUCSBOXWithGFNI(t *testing.T) {
	m1 := Set64(0xF33E408A76F65828, 0xF33E408A76F65828)
	m2 := Set64(0x9581FB0653B61C09, 0x9581FB0653B61C09)
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		// ZUC SBOX
		GF2P8AFFINEQB(dst, &m1, 0x00)
		GF2P8AFFINEINVQB(dst, &m2, 0x55)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != zuc_sbox[k+i] {
				t.Fatalf("ZUC SBOX(%v) = %x; want %x", k+i, dst.Bytes()[i], zuc_sbox[k+i])
			}
		}
	}
}
