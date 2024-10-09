package ppc64

import (
	"encoding/hex"
	"testing"
)

func TestGCM(t *testing.T) {
	var (
		XC2  = Vector128{}
		T0   = Vector128{}
		T1   = Vector128{}
		T2   = Vector128{}
		ZERO = Vector128{}
		H    = Vector128{}
		HL   = Vector128{}
		HH   = Vector128{}
		IN   = Vector128{}
		XL   = Vector128{}
		XM   = Vector128{}
		XH   = Vector128{}
		IN1  = Vector128{}
		H2   = Vector128{}
		H2L  = Vector128{}
		H2H  = Vector128{}
	)
	LXVD2X_UINT64([]uint64{0x9f1f7bff6f551138, 0x4d9430531e538fd3}, &H)
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
	VSPLTB(0, &H, &T1)   // most significant byte
	VSL(&H, &T0, &H)     // H<<=1
	VSRAB(&T1, &T2, &T1) // broadcast carry bit
	VAND(&T1, &XC2, &T1)
	VXOR(&H, &T1, &IN)           // twisted H
	VSLDOI(8, &IN, &IN, &H)      // twist even more ...
	VSLDOI(8, &ZERO, &XC2, &XC2) // 0xc2.0
	VSLDOI(8, &ZERO, &H, &HL)    // ... and split
	VSLDOI(8, &H, &ZERO, &HH)
	if hex.EncodeToString(HL.Bytes()) != "00000000000000009b2860a63ca71fa7" {
		t.Errorf("HL = %v; want 00000000000000009b2860a63ca71fa7", hex.EncodeToString(HL.Bytes()))
	}
	if hex.EncodeToString(H.Bytes()) != "9b2860a63ca71fa7fc3ef7fedeaa2270" {
		t.Errorf("H = %v; want 9b2860a63ca71fa7fc3ef7fedeaa2270", hex.EncodeToString(H.Bytes()))
	}
	if hex.EncodeToString(HH.Bytes()) != "fc3ef7fedeaa22700000000000000000" {
		t.Errorf("HH = %v; want fc3ef7fedeaa22700000000000000000", hex.EncodeToString(HH.Bytes()))
	}
	VPMSUMD(&IN, &HL, &XL) // H.lo·H.lo
	VPMSUMD(&IN, &H, &XM)  // H.hi·H.lo+H.lo·H.hi
	VPMSUMD(&IN, &HH, &XH) // H.hi·H.hi

	VPMSUMD(&XL, &XC2, &T2) // 1st reduction phase

	VSLDOI(8, &XM, &ZERO, &T0)
	VSLDOI(8, &ZERO, &XM, &T1)
	VXOR(&XL, &T0, &XL)
	VXOR(&XH, &T1, &XH)

	VSLDOI(8, &XL, &XL, &XL)
	VXOR(&XL, &T2, &XL)

	VSLDOI(8, &XL, &XL, &T1) // 2nd reduction phase
	VPMSUMD(&XL, &XC2, &XL)
	VXOR(&T1, &XH, &T1)
	VXOR(&XL, &T1, &IN1)

	VSLDOI(8, &IN1, &IN1, &H2)
	VSLDOI(8, &ZERO, &H2, &H2L)
	VSLDOI(8, &H2, &ZERO, &H2H)

	if hex.EncodeToString(H2.Bytes()) != "7ff293d6efac08928030482cce3d22c7" {
		t.Errorf("H2 = %v; want 7ff293d6efac08928030482cce3d22c7", hex.EncodeToString(H2.Bytes()))
	}

	if hex.EncodeToString(H2L.Bytes()) != "00000000000000007ff293d6efac0892" {
		t.Errorf("H2L = %v; want 00000000000000007ff293d6efac0892", hex.EncodeToString(H2L.Bytes()))
	}

	if hex.EncodeToString(H2H.Bytes()) != "8030482cce3d22c70000000000000000" {
		t.Errorf("H2H = %v; want 8030482cce3d22c70000000000000000", hex.EncodeToString(H2H.Bytes()))
	}
}
