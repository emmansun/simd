import matrix_util
import aes_sbox
import zuc_sbox

'''
生成基于AESNI的ZUC SBOX参数M1/C1/M2/C2，这里C1为0.
'''
def gen_all_m1_c1_m2_c2():
    aes_result_list = aes_sbox.get_all_WZY()
    zuc_result_list = zuc_sbox.get_all_WZY()
    Aaes = aes_sbox.affine_transform_matrix()
    Aaes_inv = Aaes.Inverse()
    Azuc = zuc_sbox.affine_transform_matrix()
    Caes = matrix_util.matrix_from_col_byte(0x63)
    for i, v1 in enumerate(aes_result_list):
        Xaes = matrix_util.matrix_from_cols(aes_sbox.gen_X(v1[0], v1[1], v1[2], v1[3], v1[4], v1[5]))
        Xaes_inv = Xaes.Inverse()
        for j, v2 in enumerate(zuc_result_list):
            Xzuc = matrix_util.matrix_from_cols(zuc_sbox.gen_X(v2[0], v2[1], v2[2], v2[3], v2[4], v2[5]))
            Xzuc_inv = Xzuc.Inverse()
            M1 = Xaes * Xzuc_inv
            M2 = Azuc * Xzuc * Xaes_inv * Aaes_inv
            C2 = M2 * Caes
            print(f'M1=','', end='')
            matrix_util.print_matrix_rows(M1)
            print(f' C1=','', end='')
            print(hex(0x0))
            print(f'M2=','', end='')
            matrix_util.print_matrix_rows(M2)
            print(f' C2=','', end='')
            print(hex(0x55 ^ matrix_util.matrix_col_byte(C2.GetColumn(0))))
            print()

gen_all_m1_c1_m2_c2()
