package s390x

import (
	"encoding/hex"
	"testing"
)

func TestXTSMul2(t *testing.T) {
	var (
		B0     = Vector128{}
		ESPERM = Vector128{}
		POLY   = Vector128{}
		T0     = Vector128{}
		T1     = Vector128{}
	)
	rawBytes, _ := hex.DecodeString("F0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF")
	VL(rawBytes, &B0)
	VL_UINT64([]uint64{0x0f0e0d0c0b0a0908, 0x0706050403020100}, &ESPERM)
	VZERO(&POLY)
	VLEIB(15, 0x87, &POLY)

	VPERM(&B0, &B0, &ESPERM, &B0)
	VESRAF(31, &B0, &T0)
	VREPF(0, &T0, &T0)
	VN(&POLY, &T0, &T0)
	VREPIB(1, &T1)
	VSL(&T1, &B0, &T1)
	VX(&T0, &T1, &B0)
	VPERM(&B0, &B0, &ESPERM, &B0)

	if hex.EncodeToString(B0.Bytes()) != "67e3e5e7e9ebedeff1f3f5f7f9fbfdff" {
		t.Errorf("B0 = %v; want 67e3e5e7e9ebedeff1f3f5f7f9fbfdff", hex.EncodeToString(B0.Bytes()))
	}
}

func TestXTSMul2GB(t *testing.T) {
	var (
		B0     = Vector128{}
		POLY   = Vector128{}
		T0     = Vector128{}
		T1     = Vector128{}
	)
	rawBytes, _ := hex.DecodeString("F0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF")
	VL(rawBytes, &B0)
	VZERO(&POLY)
	VLEIB(0, 0xe1, &POLY)

	VESLF(31, &B0, &T0)
	VESRAF(31, &T0, &T0)
	VREPF(3, &T0, &T0)
	VN(&POLY, &T0, &T0)
	VREPIB(1, &T1)
	VSRL(&T1, &B0, &T1)
	VX(&T0, &T1, &B0)

	if hex.EncodeToString(B0.Bytes()) != "9978f979fa7afb7bfc7cfd7dfe7eff7f" {
		t.Errorf("B0 = %v; want 9978f979fa7afb7bfc7cfd7dfe7eff7f", hex.EncodeToString(B0.Bytes()))
	}
}
