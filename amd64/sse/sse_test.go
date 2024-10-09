package sse

import (
	"fmt"
	"testing"
)

func TestPSRAW(t *testing.T) {
	tmp1 := Set64(0xf070b030d0509010, 0xe060a020c0408000)
	tmp := XMM{}
	MOVOU(&tmp, &tmp1)
	PSRAW(&tmp1, 32)
	if fmt.Sprintf("%x", tmp1.Bytes()) != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("PSRAW() = %v; want ffffffff...", fmt.Sprintf("%x", tmp1.Bytes()))
	}
	PSRAW(&tmp, 4)
	if fmt.Sprintf("%x", tmp.Bytes()) != "000804fc020a06fe010905fd030b07ff" {
		t.Errorf("PSRAW() = %v; want 000804fc020a06fe010905fd030b07ff", fmt.Sprintf("%x", tmp.Bytes()))
	}
}
