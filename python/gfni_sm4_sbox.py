from pyfinite import genericmatrix
import sbox_field
import aes_sbox
import sm4_sbox
XOR = lambda x,y:x^y
AND = lambda x,y:x&y
DIV = lambda x,y:x

'''
生成基于GFNI的SM4 SBOX参数M1/C1/M2
'''
def gen_all_m1_c1_m2():
    aes_result_list = sbox_field.get_all_WZY(aes_sbox.aesf)
    sm4_result_list = sbox_field.get_all_WZY(sm4_sbox.sm4f)
    Aaes = sbox_field.to_matrix(aes_sbox.AES_A)
    Aaes_inv = Aaes.Inverse()
    Asm4 = sbox_field.to_matrix(sm4_sbox.SM4_A)
    Caes = genericmatrix.GenericMatrix(size=(8, 1), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range(8):
        Caes.SetRow(i, [aes_sbox.AES_C[i]])
    Csm4 = genericmatrix.GenericMatrix(size=(8, 1), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range(8):
        Csm4.SetRow(i, [sm4_sbox.SM4_C[i]])
    for i, v1 in enumerate(aes_result_list):
        Xaes = sbox_field.to_matrix(sbox_field.gen_X(aes_sbox.aesf, v1[0], v1[1], v1[2], v1[3], v1[4], v1[5]))
        Xaes_inv = Xaes.Inverse()
        for j, v2 in enumerate(sm4_result_list):
            Xsm4 = sbox_field.to_matrix(sbox_field.gen_X(sm4_sbox.sm4f, v2[0], v2[1], v2[2], v2[3], v2[4], v2[5]))
            Xsm4_inv = Xsm4.Inverse()
            M1 = Xaes * Xsm4_inv * Asm4
            C1 = Xaes * (Xsm4_inv * Csm4)
            M2 = Asm4 * Xsm4 * Xaes_inv 
            print(f'M1=','', end='')
            sbox_field.print_m(sbox_field.matrix_rows(M1))
            print(f' C1=','', end='')
            print(hex(sbox_field.matrix_col_byte(C1.GetColumn(0))))
            print(f'M2=','', end='')
            sbox_field.print_m(sbox_field.matrix_rows(M2))
            print()            

gen_all_m1_c1_m2()
