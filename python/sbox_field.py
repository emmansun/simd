from pyfinite import ffield
from pyfinite import genericmatrix

XOR = lambda x,y:x^y
AND = lambda x,y:x&y
DIV = lambda x,y:x

def field_pow2(x, F):
    return F.Multiply(x, x)

def field_pow3(x, F):
    return F.Multiply(x, field_pow2(x, F))

def field_pow4(x, F):
    return field_pow2(field_pow2(x, F), F)

def field_pow16(x, F):
    return field_pow4(field_pow4(x, F), F)    

def get_all_WZY(F):
    result_list = []
    for i in range(256):
        if field_pow2(i, F)^i^1 == 0:
            W=i
            W_2 = field_pow2(W, F)
            N = W_2
            for j in range(256):
                if field_pow2(j, F)^j^W_2 == 0:
                    Z = j
                    Z_4 = field_pow4(Z, F)
                    u = F.Multiply(field_pow2(N, F), Z)
                    for k in range(256):
                        if field_pow2(k, F)^k^u == 0:
                            Y = k
                            Y_16 = field_pow16(k, F)
                            result_list.append([W, W_2, Z, Z_4, Y, Y_16])
    return result_list

def gen_X(F, W, W_2, Z, Z_4, Y, Y_16):
    W_2_Z_4_Y_16 = F.Multiply(F.Multiply(W_2, Z_4), Y_16)
    W_Z_4_Y_16 = F.Multiply(F.Multiply(W, Z_4), Y_16)
    W_2_Z_Y_16 = F.Multiply(F.Multiply(W_2, Z), Y_16)
    W_Z_Y_16 = F.Multiply(F.Multiply(W, Z), Y_16)
    W_2_Z_4_Y = F.Multiply(F.Multiply(W_2, Z_4), Y)
    W_Z_4_Y = F.Multiply(F.Multiply(W, Z_4), Y)
    W_2_Z_Y = F.Multiply(F.Multiply(W_2, Z), Y)
    W_Z_Y = F.Multiply(F.Multiply(W, Z), Y)
    return [W_2_Z_4_Y_16, W_Z_4_Y_16, W_2_Z_Y_16, W_Z_Y_16, W_2_Z_4_Y, W_Z_4_Y, W_2_Z_Y, W_Z_Y]

def to_matrix(x):
    m = genericmatrix.GenericMatrix(size=(8,8), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    m.SetRow(0, [(x[0] & 0x80) >> 7, (x[1] & 0x80) >> 7, (x[2] & 0x80) >> 7, (x[3] & 0x80) >> 7, (x[4] & 0x80) >> 7, (x[5] & 0x80) >> 7, (x[6] & 0x80) >> 7, (x[7] & 0x80) >> 7]) 
    m.SetRow(1, [(x[0] & 0x40) >> 6, (x[1] & 0x40) >> 6, (x[2] & 0x40) >> 6, (x[3] & 0x40) >> 6, (x[4] & 0x40) >> 6, (x[5] & 0x40) >> 6, (x[6] & 0x40) >> 6, (x[7] & 0x40) >> 6]) 
    m.SetRow(2, [(x[0] & 0x20) >> 5, (x[1] & 0x20) >> 5, (x[2] & 0x20) >> 5, (x[3] & 0x20) >> 5, (x[4] & 0x20) >> 5, (x[5] & 0x20) >> 5, (x[6] & 0x20) >> 5, (x[7] & 0x20) >> 5]) 
    m.SetRow(3, [(x[0] & 0x10) >> 4, (x[1] & 0x10) >> 4, (x[2] & 0x10) >> 4, (x[3] & 0x10) >> 4, (x[4] & 0x10) >> 4, (x[5] & 0x10) >> 4, (x[6] & 0x10) >> 4, (x[7] & 0x10) >> 4]) 
    m.SetRow(4, [(x[0] & 0x08) >> 3, (x[1] & 0x08) >> 3, (x[2] & 0x08) >> 3, (x[3] & 0x08) >> 3, (x[4] & 0x08) >> 3, (x[5] & 0x08) >> 3, (x[6] & 0x08) >> 3, (x[7] & 0x08) >> 3]) 
    m.SetRow(5, [(x[0] & 0x04) >> 2, (x[1] & 0x04) >> 2, (x[2] & 0x04) >> 2, (x[3] & 0x04) >> 2, (x[4] & 0x04) >> 2, (x[5] & 0x04) >> 2, (x[6] & 0x04) >> 2, (x[7] & 0x04) >> 2]) 
    m.SetRow(6, [(x[0] & 0x02) >> 1, (x[1] & 0x02) >> 1, (x[2] & 0x02) >> 1, (x[3] & 0x02) >> 1, (x[4] & 0x02) >> 1, (x[5] & 0x02) >> 1, (x[6] & 0x02) >> 1, (x[7] & 0x02) >> 1]) 
    m.SetRow(7, [(x[0] & 0x01) >> 0, (x[1] & 0x01) >> 0, (x[2] & 0x01) >> 0, (x[3] & 0x01) >> 0, (x[4] & 0x01) >> 0, (x[5] & 0x01) >> 0, (x[6] & 0x01) >> 0, (x[7] & 0x01) >> 0]) 
    return m

def matrix_col_byte(c):
    return (c[0] << 7) ^ (c[1] << 6) ^ (c[2] << 5) ^ (c[3] << 4) ^ (c[4] << 3) ^ (c[5] << 2) ^ (c[6] << 1) ^ (c[7] << 0)

def matrix_row_byte(c):
    return (c[0] << 7) ^ (c[1] << 6) ^ (c[2] << 5) ^ (c[3] << 4) ^ (c[4] << 3) ^ (c[5] << 2) ^ (c[6] << 1) ^ (c[7] << 0)

def matrix_cols(m):
    x = []
    for i in range(8):
        c = m.GetColumn(i)
        x.append(matrix_col_byte(c))
    return x

def matrix_rows(m):
    x = []
    for i in range(8):
        r = m.GetRow(i)
        x.append(matrix_row_byte(r))
    return x

def gen_X_inv(x):
    m = to_matrix(x)
    m_inv = m.Inverse()
    return matrix_cols(m_inv)

def G4_mul(x, y):
    '''
    GF(2^2) multiply operator, normal basis is {W^2, W}
    '''
    a = (x & 0x02) >> 1
    b = x & 0x01
    c = (y & 0x02) >> 1
    d = y & 0x01
    e = (a ^ b) & (c ^ d)
    return (((a & c) ^ e) << 1) | ((b & d) ^ e)

def G4_mul_N(x):
    '''
    GF(2^2) multiply N, normal basis is {W^2, W}, N = W^2
    '''
    a = (x & 0x02) >> 1
    b = x & 0x01
    p = b
    q = a ^ b
    return (p << 1) | q

def G4_mul_N2(x):
    '''
    GF(2^2) multiply N^2, normal basis is {W^2, W}, N = W^2
    '''
    a = (x & 0x02) >> 1
    b = x & 0x01
    return ((a ^ b) << 1) | a

def G4_inv(x):
    '''
    GF(2^2) inverse opertor
    '''        
    a = (x & 0x02) >> 1
    b = x & 0x01
    return (b << 1) | a

def G16_mul(x, y):
    '''
    GF(2^4) multiply operator, normal basis is {Z^4, Z}
    '''
    a = (x & 0xc) >> 2
    b = x & 0x03
    c = (y & 0xc) >> 2
    d = y & 0x03
    e = G4_mul(a ^ b, c ^ d)
    e = G4_mul_N(e)
    p = G4_mul(a, c) ^ e
    q = G4_mul(b, d) ^ e
    return (p << 2) | q

def G16_sq_mul_u(x):
    '''
    GF(2^4) x^2 * u operator, u = N^2 Z, N = W^2
    '''    
    a = (x & 0xc) >> 2
    b = x & 0x03
    p = G4_inv(a ^ b)
    q = G4_mul_N2(G4_inv(b))
    return (p << 2) | q

def G16_inv(x):
    '''
    GF(2^4) inverse opertor
    '''
    a = (x & 0xc) >> 2
    b = x & 0x03
    c = G4_mul_N(G4_inv(a ^ b))
    d = G4_mul(a, b)
    e = G4_inv(c ^ d)
    p = G4_mul(e, b)
    q = G4_mul(e, a)
    return (p << 2) | q

def G256_inv(x):
    '''
    GF(2^8) inverse opertor
    '''
    a = (x & 0xf0) >> 4
    b = x & 0x0f
    c = G16_sq_mul_u(a ^ b)
    d = G16_mul(a, b)
    e = G16_inv(c ^ d)
    p = G16_mul(e, b)
    q = G16_mul(e, a)
    return (p << 4) | q

def G256_new_basis(x, b):
    '''
    x presentation under new basis b
    '''
    y = 0
    for i in range(8):
        if x & (1<<((7-i))):
            y ^= b[i]
    return y

def print_m(m):
    for i, s in enumerate(m):
        print(f'0x%02x'%s,',', end='')  

def print_sbox(sbox):
    for i, s in enumerate(sbox):
        print(f'%02x'%s,',', end='')    
        if (i+1) % 16 == 0:
            print()
