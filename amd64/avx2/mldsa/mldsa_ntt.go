package mldsa

import (
	"github.com/emmansun/simd/amd64/avx2"
)

const (
	// ML-DSA global constants.
	n       = 256        // # of coefficients in the polynomials
	q       = 8380417    // 2²³ - 2¹³ + 1
	r       = 4193792    // 2³² mod q
	qNegInv = 4236238847 // -q⁻¹ mod r (q * qNegInv ≡ -1 mod r)
)

var zetasMontgomery = [n]fieldElement{
	4193792, 25847, 5771523, 7861508, 237124, 7602457, 7504169, 466468,
	1826347, 2353451, 8021166, 6288512, 3119733, 5495562, 3111497, 2680103,
	2725464, 1024112, 7300517, 3585928, 7830929, 7260833, 2619752, 6271868,
	6262231, 4520680, 6980856, 5102745, 1757237, 8360995, 4010497, 280005,
	2706023, 95776, 3077325, 3530437, 6718724, 4788269, 5842901, 3915439,
	4519302, 5336701, 3574422, 5512770, 3539968, 8079950, 2348700, 7841118,
	6681150, 6736599, 3505694, 4558682, 3507263, 6239768, 6779997, 3699596,
	811944, 531354, 954230, 3881043, 3900724, 5823537, 2071892, 5582638,
	4450022, 6851714, 4702672, 5339162, 6927966, 3475950, 2176455, 6795196,
	7122806, 1939314, 4296819, 7380215, 5190273, 5223087, 4747489, 126922,
	3412210, 7396998, 2147896, 2715295, 5412772, 4686924, 7969390, 5903370,
	7709315, 7151892, 8357436, 7072248, 7998430, 1349076, 1852771, 6949987,
	5037034, 264944, 508951, 3097992, 44288, 7280319, 904516, 3958618,
	4656075, 8371839, 1653064, 5130689, 2389356, 8169440, 759969, 7063561,
	189548, 4827145, 3159746, 6529015, 5971092, 8202977, 1315589, 1341330,
	1285669, 6795489, 7567685, 6940675, 5361315, 4499357, 4751448, 3839961,
	2091667, 3407706, 2316500, 3817976, 5037939, 2244091, 5933984, 4817955,
	266997, 2434439, 7144689, 3513181, 4860065, 4621053, 7183191, 5187039,
	900702, 1859098, 909542, 819034, 495491, 6767243, 8337157, 7857917,
	7725090, 5257975, 2031748, 3207046, 4823422, 7855319, 7611795, 4784579,
	342297, 286988, 5942594, 4108315, 3437287, 5038140, 1735879, 203044,
	2842341, 2691481, 5790267, 1265009, 4055324, 1247620, 2486353, 1595974,
	4613401, 1250494, 2635921, 4832145, 5386378, 1869119, 1903435, 7329447,
	7047359, 1237275, 5062207, 6950192, 7929317, 1312455, 3306115, 6417775,
	7100756, 1917081, 5834105, 7005614, 1500165, 777191, 2235880, 3406031,
	7838005, 5548557, 6709241, 6533464, 5796124, 4656147, 594136, 4603424,
	6366809, 2432395, 2454455, 8215696, 1957272, 3369112, 185531, 7173032,
	5196991, 162844, 1616392, 3014001, 810149, 1652634, 4686184, 6581310,
	5341501, 3523897, 3866901, 269760, 2213111, 7404533, 1717735, 472078,
	7953734, 1723600, 6577327, 1910376, 6712985, 7276084, 8119771, 4546524,
	5441381, 6144432, 7959518, 6094090, 183443, 7403526, 1612842, 4834730,
	7826001, 3919660, 8332111, 7018208, 3937738, 1400424, 7534263, 1976782,
}

/*
	var zetasAVX2Idx = [296]int{
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 8, 8, 8, 9, 9, 9, 9,
		10, 10, 10, 10, 11, 11, 11, 11,
		12, 12, 12, 12, 13, 13, 13, 13,
		14, 14, 14, 14, 15, 15, 15, 15,
		16, 16, 17, 17, 18, 18, 19, 19,
		20, 20, 21, 21, 22, 22, 23, 23,
		24, 24, 25, 25, 26, 26, 27, 27,
		28, 28, 29, 29, 30, 30, 31, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
		64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 92, 94,
		96, 98, 100, 102, 104, 106, 108, 110,
		112, 114, 116, 118, 120, 122, 124, 126,
		65, 67, 69, 71, 73, 75, 77, 79,
		81, 83, 85, 87, 89, 91, 93, 95,
		97, 99, 101, 103, 105, 107, 109, 111,
		113, 115, 117, 119, 121, 123, 125, 127,
		128, 132, 136, 140, 144, 148, 152, 156,
		160, 164, 168, 172, 176, 180, 184, 188,
		192, 196, 200, 204, 208, 212, 216, 220,
		224, 228, 232, 236, 240, 244, 248, 252,
		129, 133, 137, 141, 145, 149, 153, 157,
		161, 165, 169, 173, 177, 181, 185, 189,
		193, 197, 201, 205, 209, 213, 217, 221,
		225, 229, 233, 237, 241, 245, 249, 253,
		130, 134, 138, 142, 146, 150, 154, 158,
		162, 166, 170, 174, 178, 182, 186, 190,
		194, 198, 202, 206, 210, 214, 218, 222,
		226, 230, 234, 238, 242, 246, 250, 254,
		131, 135, 139, 143, 147, 151, 155, 159,
		163, 167, 171, 175, 179, 183, 187, 191,
		195, 199, 203, 207, 211, 215, 219, 223,
		227, 231, 235, 239, 243, 247, 251, 255,
	}
*/
var zetasMontgomeryAVX2 = [296]fieldElement{
	4193792, 25847, 5771523, 7861508, 237124, 7602457, 7504169, 466468,
	1826347, 1826347, 1826347, 1826347, 2353451, 2353451, 2353451, 2353451,
	8021166, 8021166, 8021166, 8021166, 6288512, 6288512, 6288512, 6288512,
	3119733, 3119733, 3119733, 3119733, 5495562, 5495562, 5495562, 5495562,
	3111497, 3111497, 3111497, 3111497, 2680103, 2680103, 2680103, 2680103,
	2725464, 2725464, 1024112, 1024112, 7300517, 7300517, 3585928, 3585928,
	7830929, 7830929, 7260833, 7260833, 2619752, 2619752, 6271868, 6271868,
	6262231, 6262231, 4520680, 4520680, 6980856, 6980856, 5102745, 5102745,
	1757237, 1757237, 8360995, 8360995, 4010497, 4010497, 280005, 280005,
	2706023, 95776, 3077325, 3530437, 6718724, 4788269, 5842901, 3915439,
	4519302, 5336701, 3574422, 5512770, 3539968, 8079950, 2348700, 7841118,
	6681150, 6736599, 3505694, 4558682, 3507263, 6239768, 6779997, 3699596,
	811944, 531354, 954230, 3881043, 3900724, 5823537, 2071892, 5582638,
	4450022, 4702672, 6927966, 2176455, 7122806, 4296819, 5190273, 4747489,
	3412210, 2147896, 5412772, 7969390, 7709315, 8357436, 7998430, 1852771,
	5037034, 508951, 44288, 904516, 4656075, 1653064, 2389356, 759969,
	189548, 3159746, 5971092, 1315589, 1285669, 7567685, 5361315, 4751448,
	6851714, 5339162, 3475950, 6795196, 1939314, 7380215, 5223087, 126922,
	7396998, 2715295, 4686924, 5903370, 7151892, 7072248, 1349076, 6949987,
	264944, 3097992, 7280319, 3958618, 8371839, 5130689, 8169440, 7063561,
	4827145, 6529015, 8202977, 1341330, 6795489, 6940675, 4499357, 3839961,
	2091667, 5037939, 266997, 4860065, 900702, 495491, 7725090, 4823422,
	342297, 3437287, 2842341, 4055324, 4613401, 5386378, 7047359, 7929317,
	7100756, 1500165, 7838005, 5796124, 6366809, 1957272, 5196991, 810149,
	5341501, 2213111, 7953734, 6712985, 5441381, 183443, 7826001, 3937738,
	3407706, 2244091, 2434439, 4621053, 1859098, 6767243, 5257975, 7855319,
	286988, 5038140, 2691481, 1247620, 1250494, 1869119, 1237275, 1312455,
	1917081, 777191, 5548557, 4656147, 2432395, 3369112, 162844, 1652634,
	3523897, 7404533, 1723600, 7276084, 6144432, 7403526, 3919660, 1400424,
	2316500, 5933984, 7144689, 7183191, 909542, 8337157, 2031748, 7611795,
	5942594, 1735879, 5790267, 2486353, 2635921, 1903435, 5062207, 3306115,
	5834105, 2235880, 6709241, 594136, 2454455, 185531, 1616392, 4686184,
	3866901, 1717735, 6577327, 8119771, 7959518, 1612842, 8332111, 7534263,
	3817976, 4817955, 3513181, 5187039, 819034, 7857917, 3207046, 4784579,
	4108315, 203044, 1265009, 1595974, 4832145, 7329447, 6950192, 6417775,
	7005614, 3406031, 6533464, 4603424, 8215696, 7173032, 3014001, 6581310,
	269760, 472078, 1910376, 4546524, 6094090, 4834730, 7018208, 1976782,
}

// fieldElement is an integer modulo q, an element of ℤ_q. It is always reduced.
type fieldElement uint32

// ringElement is a polynomial, an element of R_q, represented as an array.
type ringElement [n]fieldElement

// fieldReduceOnce reduces a value a < 2q.
// Also refer "A note on the implementation of the Number Theoretic Transform": https://eprint.iacr.org/2017/727.pdf .
func fieldReduceOnce(a uint32) fieldElement {
	x := a - q
	// If x underflowed, then x >= 2^32 - q > 2^31, so the top bit is set.
	x += (x >> 31) * q
	return fieldElement(x)
}

func fieldAdd(a, b fieldElement) fieldElement {
	x := uint32(a + b)
	return fieldReduceOnce(x)
}

func fieldSub(a, b fieldElement) fieldElement {
	x := uint32(a - b + q)
	return fieldReduceOnce(x)
}

func fieldReduce(a uint64) fieldElement {
	// See FIPS 204, Algorithm 49, MontgomeryReduce()
	t := uint32(a) * qNegInv
	return fieldReduceOnce(uint32((a + uint64(t)*q) >> 32))
}

func fieldMul(a, b fieldElement) fieldElement {
	x := uint64(a) * uint64(b)
	return fieldReduce(x)
}

// nttElement is an NTT representation, an element of T_q, represented as an array.
type nttElement [n]fieldElement

// ntt maps a ringElement to its nttElement representation.
//
// It implements NTT, according to FIPS 204, Algorithm 41.
// Also refer "A note on the implementation of the Number Theoretic Transform": https://eprint.iacr.org/2017/727.pdf .
func ntt(f ringElement) nttElement {
	k := 1
	// len: 128, 64, 32, ..., 1
	for len := 128; len >= 1; len /= 2 {
		// start
		for start := 0; start < n; start += 2 * len {
			zeta := zetasMontgomery[k]
			k++
			// Bounds check elimination hint.
			f, flen := f[start:start+len], f[start+len:start+len+len]
			for j := range len {
				t := fieldMul(zeta, flen[j])
				flen[j] = fieldSub(f[j], t)
				f[j] = fieldAdd(f[j], t)
			}
		}
	}
	return nttElement(f)
}

func fieldsMulAVX2(a *avx2.YMM, b *avx2.YMM, out *avx2.YMM) {
	var t0, t1, t2, qNegInvYMM, qYMM avx2.YMM
	avx2.VMOVDQU_Luint32(&qNegInvYMM, []uint32{qNegInv, qNegInv, qNegInv, qNegInv, qNegInv, qNegInv, qNegInv, qNegInv})
	avx2.VMOVDQU_Luint32(&qYMM, []uint32{q, q, q, q, q, q, q, q})
	// Multiply a and b
	avx2.VPMULUDQ(&t0, a, b) // multiply even indexes of a and b
	avx2.VPSRLQ(&t1, a, 32)
	//avx2.VMOVSHDUP(&t1, a)
	avx2.VPSRLQ(&t2, b, 32)
	//avx2.VMOVSHDUP(&t2, b)
	avx2.VPMULUDQ(&t1, &t1, &t2) // multiply odd indexes of a and b

	// Montgomery reduction: t1 = a * b * qNegInv mod r
	avx2.VPMULUDQ(out, &t0, &qNegInvYMM)
	avx2.VPMULUDQ(&t2, &t1, &qNegInvYMM)

	avx2.VPMULUDQ(out, out, &qYMM)
	avx2.VPMULUDQ(&t2, &t2, &qYMM)

	avx2.VPADDQ(out, out, &t0)
	avx2.VPADDQ(&t1, &t1, &t2)

	avx2.VPSRLQ(out, out, 32)
	// avx2.VMOVSHDUP(out, out)
	avx2.VPBLENDD(out, out, &t1, 0xAA)

	// Final reduction: if out >= q, subtract q
	avx2.VPCMPGTD(&t0, &qYMM, out)
	avx2.VPANDN(&t0, &t0, &qYMM)
	avx2.VPSUBD(out, out, &t0)
}

func fieldsMulEvenOddAVX2(a *avx2.YMM, even, odd *avx2.YMM, out *avx2.YMM) {
	var t0, t1, t2, qNegInvYMM, qYMM avx2.YMM
	avx2.VMOVDQU_Luint32(&qNegInvYMM, []uint32{qNegInv, qNegInv, qNegInv, qNegInv, qNegInv, qNegInv, qNegInv, qNegInv})
	avx2.VMOVDQU_Luint32(&qYMM, []uint32{q, q, q, q, q, q, q, q})
	// Multiply a and b
	avx2.VPMULUDQ(&t0, a, even) // multiply even indexes of a and even
	avx2.VPSRLQ(&t1, a, 32)
	//avx2.VMOVSHDUP(&t1, a)
	avx2.VPMULUDQ(&t1, &t1, odd) // multiply odd indexes of a and odd

	// Montgomery reduction: t1 = a * b * qNegInv mod r
	avx2.VPMULUDQ(out, &t0, &qNegInvYMM)
	avx2.VPMULUDQ(&t2, &t1, &qNegInvYMM)

	avx2.VPMULUDQ(out, out, &qYMM)
	avx2.VPMULUDQ(&t2, &t2, &qYMM)

	avx2.VPADDQ(out, out, &t0)
	avx2.VPADDQ(&t1, &t1, &t2)

	avx2.VPSRLQ(out, out, 32)
	// avx2.VMOVSHDUP(out, out)
	avx2.VPBLENDD(out, out, &t1, 0xAA)

	// Final reduction: if out >= q, subtract q
	avx2.VPCMPGTD(&t0, &qYMM, out)
	avx2.VPANDN(&t0, &t0, &qYMM)
	avx2.VPSUBD(out, out, &t0)
}

func butterflyAVX2(even, odd, zetasL, zetasH, outEven, outOdd *avx2.YMM) {
	var t0, t1, t2, qYMM avx2.YMM
	avx2.VMOVDQU_Luint32(&qYMM, []uint32{q, q, q, q, q, q, q, q})

	fieldsMulEvenOddAVX2(odd, zetasL, zetasH, &t0)

	avx2.VPADDD(&t1, &qYMM, even)
	avx2.VPSUBD(outOdd, &t1, &t0)
	// Final reduction: if outOdd >= q, subtract q
	avx2.VPCMPGTD(&t2, &qYMM, outOdd)
	avx2.VPANDN(&t2, &t2, &qYMM)
	avx2.VPSUBD(outOdd, outOdd, &t2)

	avx2.VPADDD(outEven, even, &t0)
	// Final reduction: if outEven >= q, subtract q
	avx2.VPCMPGTD(&t2, &qYMM, outEven)
	avx2.VPANDN(&t2, &t2, &qYMM)
	avx2.VPSUBD(outEven, outEven, &t2)
}

func nttLevel0to1AVX2(f []fieldElement, off int) {
	var ymm0, ymm1, ymm2, ymm3, ymm4, ymm5, ymm6, ymm7, zetasYMM avx2.YMM
	fUint32 := make([]uint32, len(f))
	for i, v := range f {
		fUint32[i] = uint32(v)
	}
	avx2.VMOVDQU_Luint32(&ymm0, fUint32[off*8+0:])
	avx2.VMOVDQU_Luint32(&ymm1, fUint32[off*8+1*32:])
	avx2.VMOVDQU_Luint32(&ymm2, fUint32[off*8+2*32:])
	avx2.VMOVDQU_Luint32(&ymm3, fUint32[off*8+3*32:])
	avx2.VMOVDQU_Luint32(&ymm4, fUint32[off*8+4*32:])
	avx2.VMOVDQU_Luint32(&ymm5, fUint32[off*8+5*32:])
	avx2.VMOVDQU_Luint32(&ymm6, fUint32[off*8+6*32:])
	avx2.VMOVDQU_Luint32(&ymm7, fUint32[off*8+7*32:])

	zeta1 := uint32(zetasMontgomeryAVX2[1])
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{zeta1, zeta1, zeta1, zeta1, zeta1, zeta1, zeta1, zeta1})

	butterflyAVX2(&ymm0, &ymm4, &zetasYMM, &zetasYMM, &ymm0, &ymm4)
	butterflyAVX2(&ymm1, &ymm5, &zetasYMM, &zetasYMM, &ymm1, &ymm5)
	butterflyAVX2(&ymm2, &ymm6, &zetasYMM, &zetasYMM, &ymm2, &ymm6)
	butterflyAVX2(&ymm3, &ymm7, &zetasYMM, &zetasYMM, &ymm3, &ymm7)

	// level 1: offset = 64, step = 2
	// zeta indexes = 2, 3
	zeta2 := uint32(zetasMontgomeryAVX2[2])
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{zeta2, zeta2, zeta2, zeta2, zeta2, zeta2, zeta2, zeta2})

	butterflyAVX2(&ymm0, &ymm2, &zetasYMM, &zetasYMM, &ymm0, &ymm2)
	butterflyAVX2(&ymm1, &ymm3, &zetasYMM, &zetasYMM, &ymm1, &ymm3)

	zeta3 := uint32(zetasMontgomeryAVX2[3])
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{zeta3, zeta3, zeta3, zeta3, zeta3, zeta3, zeta3, zeta3})
	butterflyAVX2(&ymm4, &ymm6, &zetasYMM, &zetasYMM, &ymm4, &ymm6)
	butterflyAVX2(&ymm5, &ymm7, &zetasYMM, &zetasYMM, &ymm5, &ymm7)

	avx2.VMOVEDQU_Suint32(fUint32[off*8+0:], &ymm0)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+1*32:], &ymm1)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+2*32:], &ymm2)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+3*32:], &ymm3)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+4*32:], &ymm4)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+5*32:], &ymm5)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+6*32:], &ymm6)
	avx2.VMOVEDQU_Suint32(fUint32[off*8+7*32:], &ymm7)

	for i, v := range fUint32 {
		f[i] = fieldElement(v)
	}
}

func shuffle8(in0, in1, out0, out1 *avx2.YMM) {
	avx2.VPERM2I128(out0, in0, in1, 0x20)
	avx2.VPERM2I128(out1, in0, in1, 0x31)
}

func shuffle4(in0, in1, out0, out1 *avx2.YMM) {
	avx2.VPUNPCKLQDQ(out0, in0, in1)
	avx2.VPUNPCKHQDQ(out1, in0, in1)
}

// shuffle2 will change the value in in0
func shuffle2(in0, in1, out0, out1 *avx2.YMM) {
	avx2.VPSLLQ(out0, in1, 32)
	//avx2.VMOVSLDUP(out0, in1)
	avx2.VPBLENDD(out0, in0, out0, 0xAA)
	avx2.VPSRLQ(in0, in0, 32)
	//avx2.VMOVSHDUP(in0, in0)
	avx2.VPBLENDD(out1, in0, in1, 0xAA)
}

func nttLevel2to7AVX2(f []fieldElement, off int) {
	var ymm0, ymm1, ymm2, ymm3, ymm4, ymm5, ymm6, ymm7, ymm8, zetasYMM, zetasHYMM avx2.YMM
	fUint32 := make([]uint32, len(f))
	for i, v := range f {
		fUint32[i] = uint32(v)
	}
	avx2.VMOVDQU_Luint32(&ymm0, fUint32[off*64+0:])
	avx2.VMOVDQU_Luint32(&ymm1, fUint32[off*64+1*8:])
	avx2.VMOVDQU_Luint32(&ymm2, fUint32[off*64+2*8:])
	avx2.VMOVDQU_Luint32(&ymm3, fUint32[off*64+3*8:])
	avx2.VMOVDQU_Luint32(&ymm4, fUint32[off*64+4*8:])
	avx2.VMOVDQU_Luint32(&ymm5, fUint32[off*64+5*8:])
	avx2.VMOVDQU_Luint32(&ymm6, fUint32[off*64+6*8:])
	avx2.VMOVDQU_Luint32(&ymm7, fUint32[off*64+7*8:])

	// level 2: offset = 32, step = 4
	// zeta indexes = 4, 5, 6, 7
	zeta4 := uint32(zetasMontgomeryAVX2[4+off])
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{zeta4, zeta4, zeta4, zeta4, zeta4, zeta4, zeta4, zeta4})
	butterflyAVX2(&ymm0, &ymm4, &zetasYMM, &zetasYMM, &ymm0, &ymm4)
	butterflyAVX2(&ymm1, &ymm5, &zetasYMM, &zetasYMM, &ymm1, &ymm5)
	butterflyAVX2(&ymm2, &ymm6, &zetasYMM, &zetasYMM, &ymm2, &ymm6)
	butterflyAVX2(&ymm3, &ymm7, &zetasYMM, &zetasYMM, &ymm3, &ymm7)

	// Input dword layout:
	//   ymm0 = [ 0  1  2  3 |  4  5  6  7]
	//   ymm1 = [ 8  9 10 11 | 12 13 14 15]
	//   ymm2 = [ 16 17 18 19 | 20 21 22 23]
	//   ymm3 = [ 24 25 26 27 | 28 29 30 31]
	//   ymm4 = [ 32 33 34 35 | 36 37 38 39]
	//   ymm5 = [ 40 41 42 43 | 44 45 46 47]
	//   ymm6 = [ 48 49 50 51 | 52 53 54 55]
	//   ymm7 = [ 56 57 58 59 | 60 61 62 63]
	// Required dword layout:
	//   ymm8 = [ 0  1  2  3 | 32 33 34 35]
	//   ymm4 = [ 4  5  6  7 | 36 37 38 39]
	//   ymm0 = [ 8  9 10 11 | 40 41 42 43]
	//   ymm5 = [12 13 14 15 | 44 45 46 47]
	//   ymm1 = [16 17 18 19 | 48 49 50 51]
	//   ymm6 = [20 21 22 23 | 52 53 54 55]
	//   ymm2 = [24 25 26 27 | 56 57 58 59]
	//   ymm7 = [28 29 30 31 | 60 61 62 63]
	shuffle8(&ymm0, &ymm4, &ymm8, &ymm4)
	shuffle8(&ymm1, &ymm5, &ymm0, &ymm5)
	shuffle8(&ymm2, &ymm6, &ymm1, &ymm6)
	shuffle8(&ymm3, &ymm7, &ymm2, &ymm7)

	// level 3: offset = 16, step = 8
	// zeta indexes = 8, 9, 10, 11, 12, 13, 14, 15
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[8+8*off]), uint32(zetasMontgomeryAVX2[8+8*off+1]), uint32(zetasMontgomeryAVX2[8+8*off+2]), uint32(zetasMontgomeryAVX2[8+8*off+3]),
		uint32(zetasMontgomeryAVX2[8+8*off+4]), uint32(zetasMontgomeryAVX2[8+8*off+5]), uint32(zetasMontgomeryAVX2[8+8*off+6]), uint32(zetasMontgomeryAVX2[8+8*off+7]),
	})
	butterflyAVX2(&ymm8, &ymm1, &zetasYMM, &zetasYMM, &ymm8, &ymm1)
	butterflyAVX2(&ymm4, &ymm6, &zetasYMM, &zetasYMM, &ymm4, &ymm6)
	butterflyAVX2(&ymm0, &ymm2, &zetasYMM, &zetasYMM, &ymm0, &ymm2)
	butterflyAVX2(&ymm5, &ymm7, &zetasYMM, &zetasYMM, &ymm5, &ymm7)

	// Input dword layout:
	//   ymm8 = [ 0  1  2  3 | 32 33 34 35]
	//   ymm4 = [ 4  5  6  7 | 36 37 38 39]
	//   ymm0 = [ 8  9 10 11 | 40 41 42 43]
	//   ymm5 = [12 13 14 15 | 44 45 46 47]
	//   ymm1 = [16 17 18 19 | 48 49 50 51]
	//   ymm6 = [20 21 22 23 | 52 53 54 55]
	//   ymm2 = [24 25 26 27 | 56 57 58 59]
	//   ymm7 = [28 29 30 31 | 60 61 62 63]
	// Required dword layout:
	//   ymm3 = [0 1 16 17 | 32 33 48 49]
	//   ymm1 = [2 3 18 19 | 34 35 50 51]
	//   ymm8 = [4 5 20 21 | 36 37 52 53]
	//   ymm6 = [6 7 22 23 | 38 39 54 55]
	//   ymm4 = [8 9 24 25 | 40 41 56 57]
	//   ymm2 = [10 11 26 27 | 42 43 58 59]
	//   ymm0 = [12 13 28 29 | 44 45 60 61]
	//   ymm7 = [14 15 30 31 | 46 47 62 63]
	shuffle4(&ymm8, &ymm1, &ymm3, &ymm1)
	shuffle4(&ymm4, &ymm6, &ymm8, &ymm6)
	shuffle4(&ymm0, &ymm2, &ymm4, &ymm2)
	shuffle4(&ymm5, &ymm7, &ymm0, &ymm7)

	// level 4: offset = 8, step = 16
	// zeta indexes = 16, 17, 18, ..., 30, 31
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[40+8*off]), uint32(zetasMontgomeryAVX2[40+8*off+1]), uint32(zetasMontgomeryAVX2[40+8*off+2]), uint32(zetasMontgomeryAVX2[40+8*off+3]),
		uint32(zetasMontgomeryAVX2[40+8*off+4]), uint32(zetasMontgomeryAVX2[40+8*off+5]), uint32(zetasMontgomeryAVX2[40+8*off+6]), uint32(zetasMontgomeryAVX2[40+8*off+7]),
	})
	butterflyAVX2(&ymm3, &ymm4, &zetasYMM, &zetasYMM, &ymm3, &ymm4)
	butterflyAVX2(&ymm1, &ymm2, &zetasYMM, &zetasYMM, &ymm1, &ymm2)
	butterflyAVX2(&ymm8, &ymm0, &zetasYMM, &zetasYMM, &ymm8, &ymm0)
	butterflyAVX2(&ymm6, &ymm7, &zetasYMM, &zetasYMM, &ymm6, &ymm7)

	// Input word layout:
	//   ymm3 = [0 1 16 17 | 32 33 48 49]
	//   ymm1 = [2 3 18 19 | 34 35 50 51]
	//   ymm8 = [4 5 20 21 | 36 37 52 53]
	//   ymm6 = [6 7 22 23 | 38 39 54 55]
	//   ymm4 = [8 9 24 25 | 40 41 56 57]
	//   ymm2 = [10 11 26 27 | 42 43 58 59]
	//   ymm0 = [12 13 28 29 | 44 45 60 61]
	//   ymm7 = [14 15 30 31 | 46 47 62 63]
	// Required word layout:
	//   ymm5 = [0 8 16 24 | 32 40 48 56]
	//   ymm4 = [1 9 17 25 | 33 41 49 57]
	//   ymm3 = [2 10 18 26 | 34 42 50 58]
	//   ymm2 = [3 11 19 27 | 35 43 51 59]
	//   ymm1 = [4 12 20 28 | 36 44 52 60]
	//   ymm0 = [5 13 21 29 | 37 45 53 61]
	//   ymm8 = [6 14 22 30 | 38 46 54 62]
	//   ymm7 = [7 15 23 31 | 39 47 55 63]
	shuffle2(&ymm3, &ymm4, &ymm5, &ymm4)
	shuffle2(&ymm1, &ymm2, &ymm3, &ymm2)
	shuffle2(&ymm8, &ymm0, &ymm1, &ymm0)
	shuffle2(&ymm6, &ymm7, &ymm8, &ymm7)

	// level 5: offset = 4, step = 32
	// zeta indexes = 32, 33, 34, ..., 62, 63
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[72+8*off]), uint32(zetasMontgomeryAVX2[72+8*off+1]), uint32(zetasMontgomeryAVX2[72+8*off+2]), uint32(zetasMontgomeryAVX2[72+8*off+3]),
		uint32(zetasMontgomeryAVX2[72+8*off+4]), uint32(zetasMontgomeryAVX2[72+8*off+5]), uint32(zetasMontgomeryAVX2[72+8*off+6]), uint32(zetasMontgomeryAVX2[72+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm5, &ymm1, &zetasYMM, &zetasHYMM, &ymm5, &ymm1)
	butterflyAVX2(&ymm4, &ymm0, &zetasYMM, &zetasHYMM, &ymm4, &ymm0)
	butterflyAVX2(&ymm3, &ymm8, &zetasYMM, &zetasHYMM, &ymm3, &ymm8)
	butterflyAVX2(&ymm2, &ymm7, &zetasYMM, &zetasHYMM, &ymm2, &ymm7)

	// level 6: offset = 2, step = 64
	// zeta indexes = 64, 65, 66, ..., 62, 127
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[104+8*off]), uint32(zetasMontgomeryAVX2[104+8*off+1]), uint32(zetasMontgomeryAVX2[104+8*off+2]), uint32(zetasMontgomeryAVX2[104+8*off+3]),
		uint32(zetasMontgomeryAVX2[104+8*off+4]), uint32(zetasMontgomeryAVX2[104+8*off+5]), uint32(zetasMontgomeryAVX2[104+8*off+6]), uint32(zetasMontgomeryAVX2[104+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm5, &ymm3, &zetasYMM, &zetasHYMM, &ymm5, &ymm3)
	butterflyAVX2(&ymm4, &ymm2, &zetasYMM, &zetasHYMM, &ymm4, &ymm2)

	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[136+8*off]), uint32(zetasMontgomeryAVX2[136+8*off+1]), uint32(zetasMontgomeryAVX2[136+8*off+2]), uint32(zetasMontgomeryAVX2[136+8*off+3]),
		uint32(zetasMontgomeryAVX2[136+8*off+4]), uint32(zetasMontgomeryAVX2[136+8*off+5]), uint32(zetasMontgomeryAVX2[136+8*off+6]), uint32(zetasMontgomeryAVX2[136+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm1, &ymm8, &zetasYMM, &zetasHYMM, &ymm1, &ymm8)
	butterflyAVX2(&ymm0, &ymm7, &zetasYMM, &zetasHYMM, &ymm0, &ymm7)

	// level 7: offset = 1, step = 128
	// zeta indexes = 128, 129, 130, ..., 254, 255
	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[168+8*off]), uint32(zetasMontgomeryAVX2[168+8*off+1]), uint32(zetasMontgomeryAVX2[168+8*off+2]), uint32(zetasMontgomeryAVX2[168+8*off+3]),
		uint32(zetasMontgomeryAVX2[168+8*off+4]), uint32(zetasMontgomeryAVX2[168+8*off+5]), uint32(zetasMontgomeryAVX2[168+8*off+6]), uint32(zetasMontgomeryAVX2[168+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm5, &ymm4, &zetasYMM, &zetasHYMM, &ymm5, &ymm4)

	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[168+32+8*off]), uint32(zetasMontgomeryAVX2[168+32+8*off+1]), uint32(zetasMontgomeryAVX2[168+32+8*off+2]), uint32(zetasMontgomeryAVX2[168+32+8*off+3]),
		uint32(zetasMontgomeryAVX2[168+32+8*off+4]), uint32(zetasMontgomeryAVX2[168+32+8*off+5]), uint32(zetasMontgomeryAVX2[168+32+8*off+6]), uint32(zetasMontgomeryAVX2[168+32+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm3, &ymm2, &zetasYMM, &zetasHYMM, &ymm3, &ymm2)

	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[168+64+8*off]), uint32(zetasMontgomeryAVX2[168+64+8*off+1]), uint32(zetasMontgomeryAVX2[168+64+8*off+2]), uint32(zetasMontgomeryAVX2[168+64+8*off+3]),
		uint32(zetasMontgomeryAVX2[168+64+8*off+4]), uint32(zetasMontgomeryAVX2[168+64+8*off+5]), uint32(zetasMontgomeryAVX2[168+64+8*off+6]), uint32(zetasMontgomeryAVX2[168+64+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm1, &ymm0, &zetasYMM, &zetasHYMM, &ymm1, &ymm0)

	avx2.VMOVDQU_Luint32(&zetasYMM, []uint32{
		uint32(zetasMontgomeryAVX2[168+96+8*off]), uint32(zetasMontgomeryAVX2[168+96+8*off+1]), uint32(zetasMontgomeryAVX2[168+96+8*off+2]), uint32(zetasMontgomeryAVX2[168+96+8*off+3]),
		uint32(zetasMontgomeryAVX2[168+96+8*off+4]), uint32(zetasMontgomeryAVX2[168+96+8*off+5]), uint32(zetasMontgomeryAVX2[168+96+8*off+6]), uint32(zetasMontgomeryAVX2[168+96+8*off+7]),
	})
	avx2.VPSRLQ(&zetasHYMM, &zetasYMM, 32)
	butterflyAVX2(&ymm8, &ymm7, &zetasYMM, &zetasHYMM, &ymm8, &ymm7)

	// Input word layout:
	//   ymm5 = [0 8 16 24 | 32 40 48 56]
	//   ymm4 = [1 9 17 25 | 33 41 49 57]
	//   ymm3 = [2 10 18 26 | 34 42 50 58]
	//   ymm2 = [3 11 19 27 | 35 43 51 59]
	//   ymm1 = [4 12 20 28 | 36 44 52 60]
	//   ymm0 = [5 13 21 29 | 37 45 53 61]
	//   ymm8 = [6 14 22 30 | 38 46 54 62]
	//   ymm7 = [7 15 23 31 | 39 47 55 63]
	// Required word layout:
	//   ymm0 = [0 1 2 3 | 4 5 6 7]
	//   ymm1 = [8 9 10 11 | 12 13 14 15]
	//   ymm2 = [16 17 18 19 | 20 21 22 23]
	//   ymm3 = [24 25 26 27 | 28 29 30 31]
	//   ymm4 = [32 33 34 35 | 36 37 38 39]
	//   ymm5 = [40 41 42 43 | 44 45 46 47]
	//   ymm6 = [48 49 50 51 | 52 53 54 55]
	//   ymm7 = [56 57 58 59 | 60 61 62 63]

	// matrix transpose
	shuffle8(&ymm5, &ymm1, &ymm6, &ymm1)
	shuffle8(&ymm4, &ymm0, &ymm5, &ymm0)
	shuffle8(&ymm3, &ymm8, &ymm4, &ymm8)
	shuffle8(&ymm2, &ymm7, &ymm3, &ymm7)

	shuffle4(&ymm6, &ymm4, &ymm2, &ymm4)
	shuffle4(&ymm1, &ymm8, &ymm6, &ymm8)
	shuffle4(&ymm5, &ymm3, &ymm1, &ymm3)
	shuffle4(&ymm0, &ymm7, &ymm5, &ymm7)

	shuffle2(&ymm2, &ymm1, &ymm0, &ymm1)
	shuffle2(&ymm4, &ymm3, &ymm2, &ymm3)
	shuffle2(&ymm6, &ymm5, &ymm4, &ymm5)
	shuffle2(&ymm8, &ymm7, &ymm6, &ymm7)

	avx2.VMOVEDQU_Suint32(fUint32[off*64+0:], &ymm0)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+1*8:], &ymm1)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+2*8:], &ymm2)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+3*8:], &ymm3)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+4*8:], &ymm4)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+5*8:], &ymm5)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+6*8:], &ymm6)
	avx2.VMOVEDQU_Suint32(fUint32[off*64+7*8:], &ymm7)

	for i, v := range fUint32 {
		f[i] = fieldElement(v)
	}
}

func nttAVX2(f ringElement) nttElement {
	for off := 0; off < n/64; off++ {
		nttLevel0to1AVX2(f[:], off)
	}

	for off := 0; off < n/64; off++ {
		nttLevel2to7AVX2(f[:], off)
	}

	return nttElement(f)
}
