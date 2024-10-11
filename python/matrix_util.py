from pyfinite import genericmatrix

XOR = lambda x,y:x^y
AND = lambda x,y:x&y
DIV = lambda x,y:x

def matrix_from_cols(cols):
    m = genericmatrix.GenericMatrix(size=(8, 8), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range (8):
        k = 7 - i
        j = 1 << k
        m.SetRow(i, [(cols[0] & j) >> k, (cols[1] & j) >> k, (cols[2] & j) >> k, (cols[3] & j) >> k, (cols[4] & j) >> k, (cols[5] & j) >> k, (cols[6] & j) >> k, (cols[7] & j) >> k])    

    return m

def matrix_from_rows(rows):
    m = genericmatrix.GenericMatrix(size=(8, 8), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range (8):
        m.SetRow(i, [(rows[i] >> 7)&1, (rows[i] >> 6)&1, (rows[i] >> 5)&1, (rows[i] >> 4)&1, (rows[i] >> 3)&1, (rows[i] >> 2)&1, (rows[i] >> 1)&1, (rows[i] >> 0)&1])
    return m

def matrix_from_col_byte(col_byte):
    m = genericmatrix.GenericMatrix(size=(8, 1), zeroElement=0, identityElement=1, add=XOR, mul=AND, sub=XOR, div=DIV)
    for i in range (8):
        m.SetRow(i, [col_byte >> (7-i) & 1])  
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

def print_m(m):
    for i, s in enumerate(m):
        print(f'0x%02x'%s,',', end='')

def print_sbox(sbox):
    for i, s in enumerate(sbox):
        print(f'%02x'%s,',', end='')
        if (i+1) % 16 == 0:
            print()

def print_matrix_rows(m):
    print_m(matrix_rows(m))

def gen_X_inv(x):
    m = matrix_from_cols(x)
    m_inv = m.Inverse()
    return matrix_cols(m_inv)
