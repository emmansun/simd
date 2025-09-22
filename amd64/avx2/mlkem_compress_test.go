package avx2

import (
	"bytes"
	mathrand "math/rand/v2"
	"testing"
)

func randomRingElement() [n]uint16 {
	var r [n]uint16
	for i := range r {
		r[i] = uint16(mathrand.IntN(q))
	}
	return r
}

var testcases = []struct {
	name string
	f    [256]uint16
}{
	{ // all zeros
		name: "all zeros",
		f:    [256]uint16{},
	},
	{ // all q-1
		name: "all q-1",
		f: func() (r [256]uint16) {
			for i := range r {
				r[i] = q - 1
			}
			return
		}(),
	},
	{ // all (q-1)/2
		name: "all (q-1)/2",
		f: func() (r [256]uint16) {
			for i := range r {
				r[i] = (q - 1) / 2
			}
			return
		}(),
	},
	{ // increasing
		name: "increasing",
		f: func() (r [256]uint16) {
			for i := range r {
				r[i] = uint16(i * (q - 1) / 255)
			}
			return
		}(),
	},
	{ // decreasing
		name: "decreasing",
		f: func() (r [256]uint16) {
			for i := range r {
				r[i] = uint16((255 - i) * (q - 1) / 255)
			}
			return
		}(),
	},
	{ // random
		name: "random",
		f:    randomRingElement(),
	},
}

func TestRingCompressAndEncode4(t *testing.T) {
	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out3 [n / 2]byte
			out1 := ringCompressAndEncode(nil, tc.f, 4, CompressRat)
			out2 := ringCompressAndEncode(nil, tc.f, 4, compress)
			ringCompressAndEncode4(out3[:], tc.f)
			if !bytes.Equal(out1[:], out2[:]) {
				t.Errorf("testcase %d - %v: mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out2: %x", out2)
			}
			if !bytes.Equal(out1[:], out3[:]) {
				t.Errorf("testcase %d - %v: avx2 mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out3: %x", out3)
			}

			dec1 := ringDecodeAndDecompress(out1[:], 4, DecompressRat)
			dec2 := ringDecodeAndDecompress(out2[:], 4, decompress)
			if dec1 != dec2 {
				t.Errorf("testcase %d - %v: decompress1 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec1)
				t.Errorf("want: %v", tc.f)
			}
			var dec3 [n]uint16
			ringDecodeAndDecompress4(&dec3, out3[:])
			if dec1 != dec3 {
				t.Errorf("testcase %d - %v: decompress2 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec3)
				t.Errorf("want: %v", tc.f)
			}
		})
	}
}

func TestRingCompressAndEncode10(t *testing.T) {
	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out3 [320]byte
			out1 := ringCompressAndEncode(nil, tc.f, 10, CompressRat)
			out2 := ringCompressAndEncode(nil, tc.f, 10, compress)
			ringCompressAndEncode10(out3[:], tc.f)
			if !bytes.Equal(out1[:], out2[:]) {
				t.Errorf("testcase %d - %v: mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out2: %x", out2)
			}
			if !bytes.Equal(out1[:], out3[:]) {
				t.Errorf("testcase %d - %v: avx2 mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out3: %x", out3)
			}

			dec1 := ringDecodeAndDecompress(out1[:], 10, DecompressRat)
			dec2 := ringDecodeAndDecompress(out2[:], 10, decompress)
			if dec1 != dec2 {
				t.Errorf("testcase %d - %v: decompress1 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec1)
				t.Errorf("want: %v", tc.f)
			}
			var dec3 [n]uint16
			ringDecodeAndDecompress10(&dec3, out3[:])
			if dec1 != dec3 {
				t.Errorf("testcase %d - %v: decompress2 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec3)
				t.Errorf("want: %v", tc.f)
			}
		})
	}
}

func TestRingCompressAndEncode5(t *testing.T) {
	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out3 [160]byte
			out1 := ringCompressAndEncode(nil, tc.f, 5, CompressRat)
			out2 := ringCompressAndEncode(nil, tc.f, 5, compress)
			ringCompressAndEncode5(out3[:], tc.f)
			if !bytes.Equal(out1[:], out2[:]) {
				t.Errorf("testcase %d - %v: mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out2: %x", out2)
			}
			if !bytes.Equal(out1[:], out3[:]) {
				t.Errorf("testcase %d - %v: avx2 mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out3: %x", out3)
			}
			dec1 := ringDecodeAndDecompress(out1[:], 5, DecompressRat)
			dec2 := ringDecodeAndDecompress(out2[:], 5, decompress)
			if dec1 != dec2 {
				t.Errorf("testcase %d - %v: decompress1 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec1)
				t.Errorf("want: %v", tc.f)
			}
			var dec3 [n]uint16
			ringDecodeAndDecompress5(&dec3, out3[:])
			if dec1 != dec3 {
				t.Errorf("testcase %d - %v: decompress2 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec3)
				t.Errorf("want: %v", tc.f)
			}
		})
	}
}

func TestRingCompressAndEncode11(t *testing.T) {
	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var out3 [352]byte
			out1 := ringCompressAndEncode(nil, tc.f, 11, CompressRat)
			out2 := ringCompressAndEncode(nil, tc.f, 11, compress)
			ringCompressAndEncode11(out3[:], tc.f)
			if !bytes.Equal(out1[:], out2[:]) {
				t.Errorf("testcase %d - %v: mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out2: %x", out2)
			}
			if !bytes.Equal(out1[:], out3[:]) {
				t.Errorf("testcase %d - %v: avx2 mismatch", i, tc.name)
				t.Errorf("out1: %x", out1)
				t.Errorf("out3: %x", out3)
			}

			dec1 := ringDecodeAndDecompress(out1[:], 11, DecompressRat)
			dec2 := ringDecodeAndDecompress(out2[:], 11, decompress)
			if dec1 != dec2 {
				t.Errorf("testcase %d - %v: decompress1 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec1)
				t.Errorf("want: %v", tc.f)
			}
			var dec3 [n]uint16
			ringDecodeAndDecompress11(&dec3, out3[:])
			if dec1 != dec3 {
				t.Errorf("testcase %d - %v: decompress2 mismatch", i, tc.name)
				t.Errorf("got:  %v", dec3)
				t.Errorf("want: %v", tc.f)
			}
		})
	}
}
