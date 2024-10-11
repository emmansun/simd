import matrix_util
import aes_sbox
import zuc_sbox

'''
生成基于GFNI的ZUC SBOX参数M1/M2
'''
def gen_all_m1_c1_m2():
    aes_result_list = aes_sbox.get_all_WZY()
    zuc_result_list = zuc_sbox.get_all_WZY()
    Azuc = zuc_sbox.affine_transform_matrix()
    for i, v1 in enumerate(aes_result_list):
        Xaes = matrix_util.matrix_from_cols(aes_sbox.gen_X(v1[0], v1[1], v1[2], v1[3], v1[4], v1[5]))
        Xaes_inv = Xaes.Inverse()
        for j, v2 in enumerate(zuc_result_list):
            Xzuc = matrix_util.matrix_from_cols(zuc_sbox.gen_X(v2[0], v2[1], v2[2], v2[3], v2[4], v2[5]))
            Xzuc_inv = Xzuc.Inverse()
            M1 = Xaes * Xzuc_inv
            M2 = Azuc * Xzuc * Xaes_inv
            print(f'M1=','', end='')
            matrix_util.print_matrix_rows(M1)
            print()
            print(f'M2=','', end='')
            matrix_util.print_matrix_rows(M2)
            print()
            print()

gen_all_m1_c1_m2()
