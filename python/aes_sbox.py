from pyfinite import ffield
import sbox_field
import matrix_util

def aes_f():
    gen = 0b100011011
    return ffield.FField(8, gen, useLUT=0)

aesf = aes_f()

AES_A = [0b10001111, 0b11000111, 0b11100011, 0b11110001, 0b11111000, 0b01111100, 0b00111110, 0b00011111]
AES_C = [0, 1, 1, 0, 0, 0, 1, 1]

def gen_X(W, W_2, Z, Z_4, Y, Y_16):
    return sbox_field.gen_X(aesf, W, W_2, Z, Z_4, Y, Y_16)

def get_all_WZY():
    return sbox_field.get_all_WZY(aesf)

def affine_transform_matrix():
    return matrix_util.matrix_from_cols(AES_A)

def AES_SBOX(X):
    X_inv = matrix_util.gen_X_inv(X)
    sbox = []
    for i in range(256):
        t = sbox_field.G256_new_basis(i, X_inv)
        t = sbox_field.G256_inv(t)
        t = sbox_field.G256_new_basis(t, X)
        t = sbox_field.G256_new_basis(t, AES_A)
        sbox.append(t ^ 0x63)
    return sbox

def print_all_aes_sbox():
    result_list = get_all_WZY()
    for i, v in enumerate(result_list):
        X = gen_X(v[0], v[1], v[2], v[3], v[4], v[5])
        matrix_util.print_sbox(AES_SBOX(X))
        print()
