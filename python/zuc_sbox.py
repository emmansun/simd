from pyfinite import ffield
import sbox_field

def zuc_f():
    gen = 0b110001011
    return ffield.FField(8, gen, useLUT=0)    

zucf = zuc_f()

ZUC_A = [0b01110111, 0b10111011, 0b11011101, 0b11101110, 0b11001011, 0b01101101, 0b00111110, 0b10010111]
ZUC_C = [0, 1, 0, 1, 0, 1, 0, 1]

def ZUC_SBOX(X, X_inv):
    sbox = []
    for i in range(256):
        t = sbox_field.G256_new_basis(i, X_inv)
        t = sbox_field.G256_inv(t)
        t = sbox_field.G256_new_basis(t, X)
        t = sbox_field.G256_new_basis(t, ZUC_A)
        sbox.append(t ^ 0x55)
    return sbox

def print_all_zuc_sbox():
    result_list = sbox_field.get_all_WZY(zucf)
    for i, v in enumerate(result_list):
        X = sbox_field.gen_X(zucf, v[0], v[1], v[2], v[3], v[4], v[5])
        X_inv = sbox_field.gen_X_inv(X)
        sbox_field.print_sbox(ZUC_SBOX(X, X_inv))
        print()
