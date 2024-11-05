package sm3

import "math/bits"

const (
	CONST0 = 0x79cc4519
	CONST1 = 0x7a879d8a
)

var IV = [8]uint32{0x7380166f, 0x4914b2b9, 0x172442d7, 0xda8a0600, 0xa96f30bc, 0x163138aa, 0xe38dee4d, 0xb0fb0e4e}

func P0(x uint32) uint32 {
	return x ^ bits.RotateLeft32(x, 9) ^ bits.RotateLeft32(x, 17)
}

func P1(x uint32) uint32 {
	return x ^ (x<<15 | x>>17) ^ (x<<23 | x>>9)
}

func FF(round byte, x, y, z uint32) uint32 {
	if round < 16 {
		return x ^ y ^ z
	}
	return (x & y) | (x & z) | (y & z)
}

func GG(round byte, x, y, z uint32) uint32 {
	if round < 16 {
		return x ^ y ^ z
	}
	return (x & y) | (^x & z)
}

func FF16(x, y, z uint32) uint32 {
	return (x & y) | (x & z) | (y & z)
}

func GG16(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}
