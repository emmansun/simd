from pyfinite import ffield
import matrix_util
import sbox_field

def sm4_f():
    gen = 0b111110101
    return ffield.FField(8, gen, useLUT=0)    

sm4f = sm4_f()

SM4_A = [0b11100101, 0b11110010, 0b01111001, 0b10111100, 0b01011110, 0b00101111, 0b10010111, 0b11001011]
SM4_C = [1, 1, 0, 1, 0, 0, 1, 1]

def gen_X(W, W_2, Z, Z_4, Y, Y_16):
    return sbox_field.gen_X(sm4f, W, W_2, Z, Z_4, Y, Y_16)

def get_all_WZY():
    return sbox_field.get_all_WZY(sm4f)

def affine_transform_matrix():
    return matrix_util.matrix_from_cols(SM4_A)

def SM4_SBOX(X):
    X_inv = matrix_util.gen_X_inv(X)
    sbox = []
    for i in range(256):
        t = sbox_field.G256_new_basis(i, SM4_A)
        t ^= 0xd3
        t = sbox_field.G256_new_basis(t, X_inv)
        t = sbox_field.G256_inv(t)
        t = sbox_field.G256_new_basis(t, X)
        t = sbox_field.G256_new_basis(t, SM4_A)
        sbox.append(t ^ 0xd3)
    return sbox

def print_all_sm4_sbox():
    result_list = get_all_WZY()
    for i, v in enumerate(result_list):
        X = gen_X(v[0], v[1], v[2], v[3], v[4], v[5])
        matrix_util.print_sbox(SM4_SBOX(X))
        print()
