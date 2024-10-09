package sse

import "testing"

func TestSM4SBOXWithAESNI(t *testing.T) {
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		SM4SBOXWithAESNI(dst)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != sm4_sbox[k+i] {
				t.Fatalf("SM4 SBOX(%v) = %x; want %x", k+i, dst.Bytes()[i], sm4_sbox[k+i])
			}
		}
	}	
}

func TestZUCSBOXWithAESNI(t *testing.T) {
	dst := &XMM{}
	for k := 0; k < 256; k += 16 {
		SetBytes(dst, []byte{byte(k), byte(k + 1), byte(k + 2), byte(k + 3), byte(k + 4), byte(k + 5), byte(k + 6), byte(k + 7), byte(k + 8), byte(k + 9), byte(k + 10), byte(k + 11), byte(k + 12), byte(k + 13), byte(k + 14), byte(k + 15)})
		ZUCSBOXWithAESNI(dst)
		for i := 0; i < 16; i++ {
			if dst.bytes[i] != zuc_sbox[k+i] {
				t.Fatalf("ZUC SBOX(%v) = %x; want %x", k+i, dst.Bytes()[i], zuc_sbox[k+i])
			}
		}
	}	
}
