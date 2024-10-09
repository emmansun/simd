package ppc64

import (
	"encoding/hex"
	"testing"
)

func TestGeneraeGCMPoly(t *testing.T) {
	var (
		XC2  = Vector128{}
		T0   = Vector128{}
		T1   = Vector128{}
		T2   = Vector128{}
		ZERO = Vector128{}
	)
	VSPLTISB(0x10, &XC2)      // 0xf0
	VSPLTISB(1, &T0)          // one
	VADDUBM(&XC2, &XC2, &XC2) // 0xe0
	VXOR(&ZERO, &ZERO, &ZERO)
	VOR(&XC2, &T0, &XC2)          // 0xe1
	VSLDOI(15, &XC2, &ZERO, &XC2) // 0xe1...
	VSLDOI(1, &ZERO, &T0, &T1)    // ...1
	VADDUBM(&XC2, &XC2, &XC2)     // 0xc2...
	VSPLTISB(7, &T2)
	VOR(&XC2, &T1, &XC2) // 0xc2....01
	if hex.EncodeToString(XC2.Bytes()) != "c2000000000000000000000000000001" {
		t.Errorf("XC2 = %v; want c2000000000000000000000000000001", hex.EncodeToString(XC2.Bytes()))
	}
}

func TestVSRAB(t *testing.T) {
	var (
		XC2 = Vector128{}
		T0  = Vector128{}
	)
	XC2.bytes[0] = 0xc2
	XC2.bytes[15] = 1
	VSPLTISB(1, &T0) // one
	VSRAB(&XC2, &T0, &XC2)
	if hex.EncodeToString(XC2.Bytes()) != "e1000000000000000000000000000000" {
		t.Errorf("XC2 = %v; want e1000000000000000000000000000000", hex.EncodeToString(XC2.Bytes()))
	}
	VSPLTISB(7, &T0) // 7
	VSRAB(&XC2, &T0, &XC2)
	if hex.EncodeToString(XC2.Bytes()) != "ff000000000000000000000000000000" {
		t.Errorf("XC2 = %v; want ff000000000000000000000000000000", hex.EncodeToString(XC2.Bytes()))
	}
}

func TestVSL(t *testing.T) {
	var (
		XC2 = Vector128{}
		T0  = Vector128{}
	)
	XC2.bytes[0] = 0xc2
	XC2.bytes[15] = 1
	VSPLTISB(1, &T0) // one
	VSL(&XC2, &T0, &XC2)
	if hex.EncodeToString(XC2.Bytes()) != "84000000000000000000000000000002" {
		t.Errorf("XC2 = %v; want 84000000000000000000000000000002", hex.EncodeToString(XC2.Bytes()))
	}
}

func TestXTSMul2(t *testing.T) {
	var (
		B0     = Vector128{}
		ESPERM = Vector128{}
		POLY   = Vector128{}
		T0     = Vector128{}
		T1     = Vector128{}
	)
	rawBytes, _ := hex.DecodeString("F0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF")
	LXVD2X(rawBytes, &B0)
	LXVD2X_UINT64([]uint64{0x0f0e0d0c0b0a0908, 0x0706050403020100}, &ESPERM)
	LXVD2X_UINT64([]uint64{0x0000000000000000, 0x0000000000000087}, &POLY)
	VPERM(&B0, &B0, &ESPERM, &B0)
	VSPLTB(0, &B0, &T0)
	VSPLTISB(7, &T1)
	VSRAB(&T0, &T1, &T0)
	VAND(&POLY, &T0, &T0)
	VSPLTISB(1, &T1)
	VSL(&B0, &T1, &T1)
	VXOR(&T0, &T1, &B0)
	VPERM(&B0, &B0, &ESPERM, &B0)
	if hex.EncodeToString(B0.Bytes()) != "67e3e5e7e9ebedeff1f3f5f7f9fbfdff" {
		t.Errorf("B0 = %v; want 67e3e5e7e9ebedeff1f3f5f7f9fbfdff", hex.EncodeToString(B0.Bytes()))
	}
}

func TestXTSMul2GB(t *testing.T) {
	var (
		B0     = Vector128{}
		ESPERM = Vector128{}
		POLY   = Vector128{}
		T0     = Vector128{}
		T1     = Vector128{}
	)
	rawBytes, _ := hex.DecodeString("F0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF")
	LXVD2X(rawBytes, &B0)
	LXVD2X_UINT64([]uint64{0x0f0e0d0c0b0a0908, 0x0706050403020100}, &ESPERM)
	LXVD2X_UINT64([]uint64{0xe100000000000000, 0x0000000000000000}, &POLY)
	VSPLTB(15, &B0, &T0)
	VSPLTISB(7, &T1)
	VSLB(&T0, &T1, &T0)
	VSRAB(&T0, &T1, &T0)
	VAND(&POLY, &T0, &T0)
	VSPLTISB(1, &T1)
	VSR(&B0, &T1, &B0)
	VXOR(&T0, &B0, &B0)
	if hex.EncodeToString(B0.Bytes()) != "9978f979fa7afb7bfc7cfd7dfe7eff7f" {
		t.Errorf("B0 = %v; want 9978f979fa7afb7bfc7cfd7dfe7eff7f", hex.EncodeToString(B0.Bytes()))
	}
}
