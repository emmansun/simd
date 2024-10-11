from pyfinite import ffield
import matrix_util
import sbox_field

def zuc_f():
    gen = 0b110001011
    return ffield.FField(8, gen, useLUT=0)    

zucf = zuc_f()

ZUC_A = [0b01110111, 0b10111011, 0b11011101, 0b11101110, 0b11001011, 0b01101101, 0b00111110, 0b10010111]
ZUC_C = [0, 1, 0, 1, 0, 1, 0, 1]

def gen_X(W, W_2, Z, Z_4, Y, Y_16):
    return sbox_field.gen_X(zucf, W, W_2, Z, Z_4, Y, Y_16)

def get_all_WZY():
    return sbox_field.get_all_WZY(zucf)

def affine_transform_matrix():
    return matrix_util.matrix_from_cols(ZUC_A)

def ZUC_SBOX(X):
    X_inv = matrix_util.gen_X_inv(X)
    sbox = []
    for i in range(256):
        t = sbox_field.G256_new_basis(i, X_inv)
        t = sbox_field.G256_inv(t)
        t = sbox_field.G256_new_basis(t, X)
        t = sbox_field.G256_new_basis(t, ZUC_A)
        sbox.append(t ^ 0x55)
    return sbox

def print_all_zuc_sbox():
    result_list = get_all_WZY()
    for i, v in enumerate(result_list):
        X = gen_X(v[0], v[1], v[2], v[3], v[4], v[5])
        matrix_util.print_sbox(ZUC_SBOX(X))
        print()
