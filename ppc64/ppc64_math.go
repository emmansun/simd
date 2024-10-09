package ppc64

import "encoding/binary"

func clmul(a, b uint64) (hi, lo uint64) {
	var temp uint64
	for i := 0; i < 64; i++ {
		temp = a & (b >> i) & 1
		for j := 1; j < i+1; j++ {
			temp ^= (a >> j) & (b >> (i - j)) & 1
		}
		lo |= temp << i
	}
	for i := 64; i < 127; i++ {
		temp = 0
		for j := i - 63; j < 64; j++ {
			temp ^= (a >> j) & (b >> (i - j)) & 1
		}
		hi |= temp << (i - 64)
	}
	return
}

func VPMSUMD(src1, src2, dst *Vector128) {
	hi1 := binary.BigEndian.Uint64(src1.bytes[:])
	lo1 := binary.BigEndian.Uint64(src1.bytes[8:])
	hi2 := binary.BigEndian.Uint64(src2.bytes[:])
	lo2 := binary.BigEndian.Uint64(src2.bytes[8:])

	hi, lo := clmul(hi1, hi2)
	hi3, lo3 := clmul(lo1, lo2)
	hi ^= hi3
	lo ^= lo3
	binary.BigEndian.PutUint64(dst.bytes[:], hi)
	binary.BigEndian.PutUint64(dst.bytes[8:], lo)
}
