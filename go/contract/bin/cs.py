
# this is a nusance to get installed -use the .js version-

from ethereum import utils

def checksum_encode(addr): # Takes a 20-byte binary address as input
    o = ''
    v = utils.big_endian_to_int(utils.sha3(addr))
    for i, c in enumerate(addr.encode('hex')):
        if c in '0123456789':
            o += c
        else:
            o += c.upper() if (v & (2**(255 - i))) else c.lower()
    return '0x'+o


eip55 = checksum_encode(0x9fbda871d559710256a2502a2517b794b482db40)
print eip55

