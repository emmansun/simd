package mldsa

import (
	"encoding/binary"
	"fmt"
	mathrand "math/rand/v2"
	"testing"

	"github.com/emmansun/simd/amd64/avx2"
)

func ymmUint32s(src *avx2.YMM) [8]uint32 {
	var out [8]uint32
	for i := range out {
		out[i] = binary.LittleEndian.Uint32(src.Bytes()[i*4:])
	}
	return out
}

func TestFieldsMulAVX2(t *testing.T) {
	testcases := []struct {
		name string
		a    [8]uint32
		b    [8]uint32
	}{
		{
			name: "all zeros",
		},
		{
			name: "identity",
			a:    [8]uint32{1, 1, 1, 1, 1, 1, 1, 1},
			b:    [8]uint32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "max reduced",
			a:    [8]uint32{q - 1, q - 1, q - 1, q - 1, q - 1, q - 1, q - 1, q - 1},
			b:    [8]uint32{q - 1, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "mixed lanes",
			a:    [8]uint32{0, 1, 2, 3, q/2 - 1, q / 2, q - 2, q - 1},
			b:    [8]uint32{q - 1, q - 2, q/2 + 1, q/2 - 1, 3, 2, 1, 0},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var aYMM, bYMM, outYMM avx2.YMM
			avx2.VMOVDQU_Luint32(&aYMM, tc.a[:])
			avx2.VMOVDQU_Luint32(&bYMM, tc.b[:])

			fieldsMulAVX2(&aYMM, &bYMM, &outYMM)

			got := ymmUint32s(&outYMM)
			for i := range got {
				want := uint32(fieldMul(fieldElement(tc.a[i]), fieldElement(tc.b[i])))
				if got[i] != want {
					t.Fatalf("lane %d mismatch: got %d want %d (a=%d b=%d)", i, got[i], want, tc.a[i], tc.b[i])
				}
			}
		})
	}

	prng := mathrand.New(mathrand.NewPCG(1, 2))
	for i := 0; i < 128; i++ {
		var aVals, bVals [8]uint32
		for j := range aVals {
			aVals[j] = uint32(prng.IntN(int(q)))
			bVals[j] = uint32(prng.IntN(int(q)))
		}

		t.Run("random", func(t *testing.T) {
			var aYMM, bYMM, outYMM avx2.YMM
			avx2.VMOVDQU_Luint32(&aYMM, aVals[:])
			avx2.VMOVDQU_Luint32(&bYMM, bVals[:])

			fieldsMulAVX2(&aYMM, &bYMM, &outYMM)

			got := ymmUint32s(&outYMM)
			for j := range got {
				want := uint32(fieldMul(fieldElement(aVals[j]), fieldElement(bVals[j])))
				if got[j] != want {
					t.Fatalf("case %d lane %d mismatch: got %d want %d (a=%d b=%d)", i, j, got[j], want, aVals[j], bVals[j])
				}
			}
		})
	}
}

func TestButterflyAVX2ReductionAtQ(t *testing.T) {
	// even=0, odd=0, zeta=1 -> t=0, so outOdd starts from q and must reduce to 0.
	evenVals := [8]uint32{}
	oddVals := [8]uint32{}
	zetasVals := [8]uint32{1, 1, 1, 1, 1, 1, 1, 1}

	var evenYMM, oddYMM, zetasYMM, outEven, outOdd avx2.YMM
	avx2.VMOVDQU_Luint32(&evenYMM, evenVals[:])
	avx2.VMOVDQU_Luint32(&oddYMM, oddVals[:])
	avx2.VMOVDQU_Luint32(&zetasYMM, zetasVals[:])

	butterflyAVX2(&evenYMM, &oddYMM, &zetasYMM, &zetasYMM, &outEven, &outOdd)

	wantEven := [8]uint32{}
	wantOdd := [8]uint32{}
	if got := ymmUint32s(&outEven); got != wantEven {
		t.Fatalf("outEven mismatch: got %v want %v", got, wantEven)
	}
	if got := ymmUint32s(&outOdd); got != wantOdd {
		t.Fatalf("outOdd mismatch: got %v want %v", got, wantOdd)
	}
}

func TestShuffle8(t *testing.T) {
	in0 := [8]uint32{0, 1, 2, 3, 4, 5, 6, 7}
	in1 := [8]uint32{10, 11, 12, 13, 14, 15, 16, 17}

	var in0YMM, in1YMM, out0, out1 avx2.YMM
	avx2.VMOVDQU_Luint32(&in0YMM, in0[:])
	avx2.VMOVDQU_Luint32(&in1YMM, in1[:])

	shuffle8(&in0YMM, &in1YMM, &out0, &out1)

	want0 := [8]uint32{0, 1, 2, 3, 10, 11, 12, 13}
	want1 := [8]uint32{4, 5, 6, 7, 14, 15, 16, 17}
	if got := ymmUint32s(&out0); got != want0 {
		t.Fatalf("out0 mismatch: got %v want %v", got, want0)
	}
	if got := ymmUint32s(&out1); got != want1 {
		t.Fatalf("out1 mismatch: got %v want %v", got, want1)
	}
}

func TestShuffle4(t *testing.T) {
	in0 := [8]uint32{0, 1, 2, 3, 4, 5, 6, 7}
	in1 := [8]uint32{10, 11, 12, 13, 14, 15, 16, 17}

	var in0YMM, in1YMM, out0, out1 avx2.YMM
	avx2.VMOVDQU_Luint32(&in0YMM, in0[:])
	avx2.VMOVDQU_Luint32(&in1YMM, in1[:])

	shuffle4(&in0YMM, &in1YMM, &out0, &out1)

	want0 := [8]uint32{0, 1, 10, 11, 4, 5, 14, 15}
	want1 := [8]uint32{2, 3, 12, 13, 6, 7, 16, 17}
	if got := ymmUint32s(&out0); got != want0 {
		t.Fatalf("out0 mismatch: got %v want %v", got, want0)
	}
	if got := ymmUint32s(&out1); got != want1 {
		t.Fatalf("out1 mismatch: got %v want %v", got, want1)
	}
}

func TestShuffle2(t *testing.T) {
	in0 := [8]uint32{0, 1, 2, 3, 4, 5, 6, 7}
	in1 := [8]uint32{10, 11, 12, 13, 14, 15, 16, 17}

	var in0YMM, in1YMM, out0, out1 avx2.YMM
	avx2.VMOVDQU_Luint32(&in0YMM, in0[:])
	avx2.VMOVDQU_Luint32(&in1YMM, in1[:])

	shuffle2(&in0YMM, &in1YMM, &out0, &out1)

	want0 := [8]uint32{0, 10, 2, 12, 4, 14, 6, 16}
	want1 := [8]uint32{1, 11, 3, 13, 5, 15, 7, 17}

	if got := ymmUint32s(&out0); got != want0 {
		t.Fatalf("out0 mismatch: got %v want %v", got, want0)
	}
	if got := ymmUint32s(&out1); got != want1 {
		t.Fatalf("out1 mismatch: got %v want %v", got, want1)
	}
}

func randomRingElement() ringElement {
	var r ringElement
	for i := range r {
		r[i] = fieldElement(mathrand.IntN(q))
	}
	return r
}

// print a ring element in a human-friendly way, for debugging.
func printRingElement(r nttElement) {
	for i := range 32 {
		for j := range 8 {
			fmt.Printf("%d ", r[i*8+j])
		}
		fmt.Println()
	}
}

func TestNTTAVX2(t *testing.T) {
	// Fixed-pattern cases.
	testcases := []struct {
		name string
		f    ringElement
	}{
		{
			name: "all zeros",
		},
		{
			name: "all ones",
			f: func() ringElement {
				var r ringElement
				for i := range r {
					r[i] = 1
				}
				return r
			}(),
		},
		{
			name: "all q-1",
			f: func() ringElement {
				var r ringElement
				for i := range r {
					r[i] = q - 1
				}
				return r
			}(),
		},
		{
			name: "identity basis e0",
			f: func() ringElement {
				var r ringElement
				r[0] = 1
				return r
			}(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			want := ntt(tc.f)
			got := nttAVX2(tc.f)
			if got != want {
				/*
					for i := range got {
						if got[i] != want[i] {
							t.Errorf("coefficient[%d]: got %d want %d", i, got[i], want[i])
						}
					}
					t.Fatalf("nttAVX2 result mismatch for %q (first diff shown above)", tc.name)
				*/
				printRingElement(want)
				println()
				printRingElement(got)
				t.Fatalf("nttAVX2 result mismatch for %q (first diff shown above)", tc.name)
			}
		})
	}

	// Randomised cases.
	for range 64 {
		f := randomRingElement()
		t.Run("random", func(t *testing.T) {
			want := ntt(f)
			got := nttAVX2(f)
			if got != want {
				for i := range got {
					if got[i] != want[i] {
						t.Errorf("coefficient[%d]: got %d want %d", i, got[i], want[i])
					}
				}
				t.Fatalf("nttAVX2 random result mismatch (first diff shown above)")
			}
		})
	}
}

// inverseButterflyScalar is the scalar baseline for inverseButterflyAVX2.
// It computes:
//
//	outEven = (even + odd) mod q
//	outOdd = ((q + even - odd) * zeta) mod q
func inverseButterflyScalar(even, odd, zeta fieldElement) (outEven, outOdd fieldElement) {
	outEven = fieldAdd(even, odd)
	// Compute (q + even - odd) first to avoid underflow, then multiply by zeta
	temp := fieldAdd(fieldElement(q), even)
	temp = fieldSub(temp, odd)
	outOdd = fieldMul(temp, zeta)
	return outEven, outOdd
}

func TestInverseButterflyAVX2(t *testing.T) {
	testcases := []struct {
		name         string
		evenVals     [8]uint32
		oddVals      [8]uint32
		zetaVals     [8]uint32
		zetaHighVals [8]uint32
	}{
		{
			name: "all zeros",
		},
		{
			name:     "identity",
			evenVals: [8]uint32{1, 2, 3, 4, 5, 6, 7, 8},
			oddVals:  [8]uint32{1, 1, 1, 1, 1, 1, 1, 1},
			zetaVals: [8]uint32{1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:     "q-1 values",
			evenVals: [8]uint32{q - 1, q - 1, q - 1, q - 1, q - 1, q - 1, q - 1, q - 1},
			oddVals:  [8]uint32{q - 2, q - 2, q - 2, q - 2, q - 2, q - 2, q - 2, q - 2},
			zetaVals: [8]uint32{uint32(zetasMontgomery[1]), uint32(zetasMontgomery[1]), uint32(zetasMontgomery[1]), uint32(zetasMontgomery[1]),
				uint32(zetasMontgomery[1]), uint32(zetasMontgomery[1]), uint32(zetasMontgomery[1]), uint32(zetasMontgomery[1])},
		},
		{
			name:     "mixed zetas",
			evenVals: [8]uint32{100, 200, 300, 400, 500, 600, 700, 800},
			oddVals:  [8]uint32{50, 75, 100, 125, 150, 175, 200, 225},
			zetaVals: [8]uint32{
				uint32(zetasMontgomery[1]), uint32(zetasMontgomery[2]), uint32(zetasMontgomery[3]), uint32(zetasMontgomery[4]),
				uint32(zetasMontgomery[5]), uint32(zetasMontgomery[6]), uint32(zetasMontgomery[7]), uint32(zetasMontgomery[8]),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup YMM registers
			var evenYMM, oddYMM, zetasYMM, zetasHYMM, outEvenAVX2, outOddAVX2 avx2.YMM
			avx2.VMOVDQU_Luint32(&evenYMM, tc.evenVals[:])
			avx2.VMOVDQU_Luint32(&oddYMM, tc.oddVals[:])
			avx2.VMOVDQU_Luint32(&zetasYMM, tc.zetaVals[:])

			// For zetasH, we extract high 32 bits of each 64-bit pair.
			// Since fieldsMulEvenOddAVX2 uses VPSRLQ to extract odd dwords,
			// we'll compute it the same way or pass precomputed values.
			// For simplicity, compute it from zetaVals on the fly:
			var zetasHighComputed [8]uint32
			for i := 0; i < 4; i++ {
				zetasHighComputed[2*i] = tc.zetaVals[2*i+1]
				zetasHighComputed[2*i+1] = tc.zetaVals[2*i+1]
			}
			avx2.VMOVDQU_Luint32(&zetasHYMM, zetasHighComputed[:])

			// Call AVX2 butterfly
			inverseButterflyAVX2(&evenYMM, &oddYMM, &zetasYMM, &zetasHYMM, &outEvenAVX2, &outOddAVX2)

			// Verify each lane against scalar version
			gotEven := ymmUint32s(&outEvenAVX2)
			gotOdd := ymmUint32s(&outOddAVX2)

			for lane := range 8 {
				wantEven, wantOdd := inverseButterflyScalar(
					fieldElement(tc.evenVals[lane]),
					fieldElement(tc.oddVals[lane]),
					fieldElement(tc.zetaVals[lane]),
				)

				if gotEven[lane] != uint32(wantEven) {
					t.Errorf("lane %d outEven mismatch: got %d want %d (even=%d odd=%d zeta=%d)",
						lane, gotEven[lane], wantEven, tc.evenVals[lane], tc.oddVals[lane], tc.zetaVals[lane])
				}
				if gotOdd[lane] != uint32(wantOdd) {
					t.Errorf("lane %d outOdd mismatch: got %d want %d (even=%d odd=%d zeta=%d)",
						lane, gotOdd[lane], wantOdd, tc.evenVals[lane], tc.oddVals[lane], tc.zetaVals[lane])
				}
			}
		})
	}

	// Random test cases
	prng := mathrand.New(mathrand.NewPCG(42, 43))
	for testCase := 0; testCase < 64; testCase++ {
		var evenVals, oddVals, zetaVals [8]uint32
		for i := range evenVals {
			evenVals[i] = uint32(prng.IntN(int(q)))
			oddVals[i] = uint32(prng.IntN(int(q)))
			// Use real zetasMontgomery values for realistic test
			zetaVals[i] = uint32(zetasMontgomery[prng.IntN(n)])
		}

		t.Run(fmt.Sprintf("random_%d", testCase), func(t *testing.T) {
			var evenYMM, oddYMM, zetasYMM, zetasHYMM, outEvenAVX2, outOddAVX2 avx2.YMM
			avx2.VMOVDQU_Luint32(&evenYMM, evenVals[:])
			avx2.VMOVDQU_Luint32(&oddYMM, oddVals[:])
			avx2.VMOVDQU_Luint32(&zetasYMM, zetaVals[:])

			var zetasHighComputed [8]uint32
			for i := 0; i < 4; i++ {
				zetasHighComputed[2*i] = zetaVals[2*i+1]
				zetasHighComputed[2*i+1] = zetaVals[2*i+1]
			}
			avx2.VMOVDQU_Luint32(&zetasHYMM, zetasHighComputed[:])

			inverseButterflyAVX2(&evenYMM, &oddYMM, &zetasYMM, &zetasHYMM, &outEvenAVX2, &outOddAVX2)

			gotEven := ymmUint32s(&outEvenAVX2)
			gotOdd := ymmUint32s(&outOddAVX2)

			for lane := range 8 {
				wantEven, wantOdd := inverseButterflyScalar(
					fieldElement(evenVals[lane]),
					fieldElement(oddVals[lane]),
					fieldElement(zetaVals[lane]),
				)

				if gotEven[lane] != uint32(wantEven) {
					t.Errorf("case %d lane %d outEven mismatch: got %d want %d", testCase, lane, gotEven[lane], wantEven)
				}
				if gotOdd[lane] != uint32(wantOdd) {
					t.Errorf("case %d lane %d outOdd mismatch: got %d want %d", testCase, lane, gotOdd[lane], wantOdd)
				}
			}
		})
	}
}

func TestInverseNTTAVX2(t *testing.T) {
	assertEqual := func(t *testing.T, got, want ringElement, context string) {
		t.Helper()
		if got == want {
			return
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("%s mismatch at coefficient[%d]: got %d want %d", context, i, got[i], want[i])
			}
		}
		t.Fatalf("%s mismatch", context)
	}

	// Fixed-pattern cases in the NTT domain.
	testcases := []struct {
		name string
		f    nttElement
	}{
		{
			name: "all zeros",
		},
		{
			name: "all ones",
			f: func() nttElement {
				var r nttElement
				for i := range r {
					r[i] = 1
				}
				return r
			}(),
		},
		{
			name: "all q-1",
			f: func() nttElement {
				var r nttElement
				for i := range r {
					r[i] = q - 1
				}
				return r
			}(),
		},
		{
			name: "single non-zero",
			f: func() nttElement {
				var r nttElement
				r[0] = 1
				return r
			}(),
		},
		{
			name: "from forward ntt",
			f: func() nttElement {
				return ntt(randomRingElement())
			}(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			want := inverseNTT(tc.f)
			got := inverseNTTAVX2(tc.f)
			assertEqual(t, got, want, fmt.Sprintf("inverseNTTAVX2 result for %q", tc.name))
		})
	}

	// Randomised cases.
	for range 64 {
		f := nttElement(randomRingElement())
		t.Run("random", func(t *testing.T) {
			want := inverseNTT(f)
			got := inverseNTTAVX2(f)
			assertEqual(t, got, want, "inverseNTTAVX2 random result")
		})
	}
}
