package arm64

func EIA16Bytes(data []byte, keys []uint32) uint32 {
	var (
		XTMP1             = Vector128{}
		XTMP2             = Vector128{}
		XTMP3             = Vector128{}
		XTMP4             = Vector128{}
		XTMP5             = Vector128{}
		XTMP6             = Vector128{}
		XDATA             = Vector128{}
		XDIGEST           = Vector128{}
		KS_L              = Vector128{}
		KS_M1             = Vector128{}
		BIT_REV_TAB_L     = Vector128{}
		BIT_REV_TAB_H     = Vector128{}
		BIT_REV_AND_TAB   = Vector128{}
		SHUF_MASK_DW0_DW1 = Vector128{}
		SHUF_MASK_DW2_DW3 = Vector128{}
	)

	VLD1_2D([]uint64{0x0e060a020c040800, 0x0f070b030d050901}, &BIT_REV_TAB_L)
	VLD1_2D([]uint64{0xe060a020c0408000, 0xf070b030d0509010}, &BIT_REV_TAB_H)
	VDUP_BYTE(0x0f, &BIT_REV_AND_TAB)
	VLD1_2D([]uint64{0xffffffff03020100, 0xffffffff07060504}, &SHUF_MASK_DW0_DW1)
	VLD1_2D([]uint64{0xffffffff0b0a0908, 0xffffffff0f0e0d0c}, &SHUF_MASK_DW2_DW3)

	// load data
	VLD1_16B(data, &XDATA)
	VAND(&XDATA, &BIT_REV_AND_TAB, &XTMP3)
	VUSHR_S(4, &XDATA, &XTMP1)
	VAND(&XTMP1, &BIT_REV_AND_TAB, &XTMP1)

	VTBL_B(&XTMP3, []*Vector128{&BIT_REV_TAB_H}, &XTMP3)
	VTBL_B(&XTMP1, []*Vector128{&BIT_REV_TAB_L}, &XTMP1)
	VEOR(&XTMP3, &XTMP1, &XTMP3) // XTMP3 - bit reverse data bytes

	// ZUC authentication part, 4x32 data bits
	// setup KS
	VLD1_4S(keys, &XTMP1)
	VLD1_4S(keys[4:], &XTMP2)
	VDUP_S(XTMP1.Uint32s()[1], &KS_L)
	VMOV_S(&XTMP1, &KS_L, 0, 1)
	VMOV_S(&XTMP1, &KS_L, 2, 2) // KS bits [63:32 31:0 95:64 63:32]
	VDUP_S(XTMP1.Uint32s()[3], &KS_M1)
	VMOV_S(&XTMP1, &KS_M1, 2, 1)
	VMOV_S(&XTMP2, &KS_M1, 0, 2)

	// setup data
	VTBL_B(&SHUF_MASK_DW0_DW1, []*Vector128{&XTMP3}, &XTMP1) // XTMP1 - Data bits [31:0 0s 63:32 0s]
	VTBL_B(&SHUF_MASK_DW2_DW3, []*Vector128{&XTMP3}, &XTMP2) // XTMP2 - Data bits [95:64 0s 127:96 0s]

	// clmul
	// xor the results from 4 32-bit words together
	// Calculate lower 32 bits of tag
	VPMULL(&KS_L, &XTMP1, &XTMP3)
	VPMULL2(&KS_L, &XTMP1, &XTMP4)
	VPMULL(&KS_M1, &XTMP2, &XTMP5)
	VPMULL2(&KS_M1, &XTMP2, &XTMP6)

	VEOR(&XTMP3, &XTMP4, &XTMP3)
	VEOR(&XTMP5, &XTMP6, &XTMP5)
	VEOR(&XTMP3, &XTMP5, &XDIGEST)

	return XDIGEST.Uint32s()[1]
}
