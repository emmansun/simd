from pyfinite import ffield
import sbox_field

def sm4_f():
    gen = 0b111110101
    return ffield.FField(8, gen, useLUT=0)    

sm4f = sm4_f()

SM4_A = [0b11100101, 0b11110010, 0b01111001, 0b10111100, 0b01011110, 0b00101111, 0b10010111, 0b11001011]
SM4_C = [1, 1, 0, 1, 0, 0, 1, 1]

def SM4_SBOX(X, X_inv):
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
    result_list = sbox_field.get_all_WZY(sm4f)
    for i, v in enumerate(result_list):
        X = sbox_field.gen_X(sm4f, v[0], v[1], v[2], v[3], v[4], v[5])
        X_inv = sbox_field.gen_X_inv(X)
        sbox_field.print_sbox(SM4_SBOX(X, X_inv))
        print()
