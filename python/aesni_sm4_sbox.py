import matrix_util
import aes_sbox
import sm4_sbox

'''
生成基于AESNI的SM4 SBOX参数M1/C1/M2/C2
'''
def gen_all_m1_c1_m2_c2():
    aes_result_list = aes_sbox.get_all_WZY()
    sm4_result_list = sm4_sbox.get_all_WZY()
    Aaes = aes_sbox.affine_transform_matrix()
    Aaes_inv = Aaes.Inverse()
    Asm4 = sm4_sbox.affine_transform_matrix()
    Caes = matrix_util.matrix_from_col_byte(0x63)
    Csm4 = matrix_util.matrix_from_col_byte(0xd3)
    for i, v1 in enumerate(aes_result_list):
        Xaes = matrix_util.matrix_from_cols(aes_sbox.gen_X(v1[0], v1[1], v1[2], v1[3], v1[4], v1[5]))
        Xaes_inv = Xaes.Inverse()
        for j, v2 in enumerate(sm4_result_list):
            Xsm4 = matrix_util.matrix_from_cols(sm4_sbox.gen_X(v2[0], v2[1], v2[2], v2[3], v2[4], v2[5]))
            Xsm4_inv = Xsm4.Inverse()
            M1 = Xaes * Xsm4_inv * Asm4
            C1 = Xaes * (Xsm4_inv * Csm4)
            M2 = Asm4 * Xsm4 * Xaes_inv * Aaes_inv
            C2 = M2 * Caes
            print(f'M1=','', end='')
            matrix_util.print_matrix_rows(M1)

            print(f' C1=','', end='')
            print(hex(matrix_util.matrix_col_byte(C1.GetColumn(0))))

            print(f'M2=','', end='')
            matrix_util.print_matrix_rows(M2)

            print(f' C2=','', end='')
            print(hex(0xd3 ^ matrix_util.matrix_col_byte(C2.GetColumn(0))))
            print()

gen_all_m1_c1_m2_c2()
