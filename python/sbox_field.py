import matrix_util

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
