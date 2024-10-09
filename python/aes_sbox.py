from pyfinite import ffield
import sbox_field

def aes_f():
    gen = 0b100011011
    return ffield.FField(8, gen, useLUT=0)

aesf = aes_f()

AES_A = [0b10001111, 0b11000111, 0b11100011, 0b11110001, 0b11111000, 0b01111100, 0b00111110, 0b00011111]
AES_C = [0, 1, 1, 0, 0, 0, 1, 1]

def AES_SBOX(X, X_inv):
    sbox = []
    for i in range(256):
        t = sbox_field.G256_new_basis(i, X_inv)
        t = sbox_field.G256_inv(t)
        t = sbox_field.G256_new_basis(t, X)
        t = sbox_field.G256_new_basis(t, AES_A)
        sbox.append(t ^ 0x63)
    return sbox

def print_all_aes_sbox():
    result_list = sbox_field.get_all_WZY(aesf)
    for i, v in enumerate(result_list):
        X = sbox_field.gen_X(aesf, v[0], v[1], v[2], v[3], v[4], v[5])
        X_inv = sbox_field.gen_X_inv(X)
        sbox_field.print_sbox(AES_SBOX(X, X_inv))
        print()
