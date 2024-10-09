from pyfinite import genericmatrix
import sbox_field
import aes_sbox
import zuc_sbox
XOR = lambda x,y:x^y
AND = lambda x,y:x&y
DIV = lambda x,y:x

'''
生成基于AESNI的ZUC SBOX参数M1/C1/M2/C2，这里C1为0.
'''
def gen_all_m1_c1_m2_c2():
    aes_result_list = sbox_field.get_all_WZY(aes_sbox.aesf)
    zuc_result_list = sbox_field.get_all_WZY(zuc_sbox.zucf)
    Aaes = sbox_field.to_matrix(aes_sbox.AES_A)
    Aaes_inv = Aaes.Inverse()
    Azuc = sbox_field.to_matrix(zuc_sbox.ZUC_A)
    Caes = genericmatrix.GenericMatrix(size=(8, 1), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range(8):
        Caes.SetRow(i, [aes_sbox.AES_C[i]])
    Czuc = genericmatrix.GenericMatrix(size=(8, 1), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range(8):
        Czuc.SetRow(i, [zuc_sbox.ZUC_C[i]])
    for i, v1 in enumerate(aes_result_list):
        Xaes = sbox_field.to_matrix(sbox_field.gen_X(aes_sbox.aesf, v1[0], v1[1], v1[2], v1[3], v1[4], v1[5]))
        Xaes_inv = Xaes.Inverse()
        for j, v2 in enumerate(zuc_result_list):
            Xzuc = sbox_field.to_matrix(sbox_field.gen_X(zuc_sbox.zucf, v2[0], v2[1], v2[2], v2[3], v2[4], v2[5]))
            Xzuc_inv = Xzuc.Inverse()
            M1 = Xaes * Xzuc_inv
            M2 = Azuc * Xzuc * Xaes_inv * Aaes_inv
            C2 = M2 * Caes
            print(f'M1=','', end='')
            sbox_field.print_m(sbox_field.matrix_rows(M1))
            print(f' C1=','', end='')
            print(hex(0x0))
            print(f'M2=','', end='')
            sbox_field.print_m(sbox_field.matrix_rows(M2))
            print(f' C2=','', end='')
            print(hex(0x55 ^ sbox_field.matrix_col_byte(C2.GetColumn(0))))
            print()

gen_all_m1_c1_m2_c2()
