package s390x

import (
	"encoding/hex"
	"testing"
)

func TestTransposeMatrix(t *testing.T) {
	t0 := &Vector128{}
	t1 := &Vector128{}
	t2 := &Vector128{}
	t3 := &Vector128{}

	VL_UINT64([]uint64{0x2222222211111111, 0x4444444433333333}, t0)
	VL_UINT64([]uint64{0x6666666655555555, 0x8888888877777777}, t1)
	VL_UINT64([]uint64{0xaaaaaaaa99999999, 0xccccccccbbbbbbbb}, t2)
	VL_UINT64([]uint64{0xeeeeeeeedddddddd, 0x00000000ffffffff}, t3)

	TransposeMatrix(t0, t1, t2, t3)

	got0 := hex.EncodeToString(t0.Bytes())
	if got0 != "2222222266666666aaaaaaaaeeeeeeee" {
		t.Errorf("t0 = %v; want 2222222266666666aaaaaaaaeeeeeeee", got0)
	}
	got1 := hex.EncodeToString(t1.Bytes())
	if got1 != "111111115555555599999999dddddddd" {
		t.Errorf("t1 = %v; want 111111115555555599999999dddddddd", got1)
	}
	got2 := hex.EncodeToString(t2.Bytes())
	if got2 != "4444444488888888cccccccc00000000" {
		t.Errorf("t2 = %v; want 4444444488888888cccccccc00000000", got2)
	}
	got3 := hex.EncodeToString(t3.Bytes())
	if got3 != "3333333377777777bbbbbbbbffffffff" {
		t.Errorf("t3 = %v; want 3333333377777777bbbbbbbbffffffff", got3)
	}
}
