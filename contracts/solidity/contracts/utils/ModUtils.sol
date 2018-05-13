pragma solidity ^0.4.21;

library ModUtils {

    function modExp(uint256 a, uint256 exponent, uint256 p)
        public
        constant returns(uint256 raised)
    {
        uint256[6] memory args = [32, 32, 32, a, exponent, p];
        uint256[1] memory output;
        assembly {
            // 0x05 is the modular exponent contract address
            if iszero(call(not(0), 0x05, 0, args, 0xc0, output, 0x20)) {
                revert(0, 0)
            }
        }
        return output[0];
    }

    function modSqrt(uint256 a, uint256 p)
        public
        constant returns(uint256)
    {
        if (a == 0) {
            return 0;
        }
        else if (p == 2) {
            return p;
        }
        else if (p % 4 == 3) {
            return modExp(a, (p + 1) / 4, p);
        }
        else if (legendre(a, p) != 1) {
            return 0;
        }

        uint256 s = p - 1;
        uint256 e = 0;

        while (s % 2 == 0) {
            s = s / 2; // TODO check operator rounds like Python
            e = e + 1;
        }

        uint256 n = 2;
        while (legendre(n, p) != -1) {
            n = n + 1;
        }

		uint256 x = modExp(a, (s + 1) / 2, p);
		uint256 b = modExp(a, s, p);
		uint256 g = modExp(n, s, p);
		uint256 r = e;
        uint256 gs = 0;
        uint256 m = 0;
        uint256 t = b;

		while (true) {
			t = b;
			m = 0;

            for (m = 0; m < r; m++) {
                if (t == 1) {
                    break;
                }
				t = modExp(t, 2, p);
            }

            if (m == 0) {
                return x;
            }

			gs = modExp(g, 2 ** (r - m - 1), p);
			g = (gs * gs) % p;
			x = (x * gs) % p;
			b = (b * g) % p;
			r = m;
		}
    }


    function legendre(uint256 a, uint256 p)
        public
        constant returns(int)
    {
        uint256 raised = modExp(a, (p - 1) / uint256(2), p);

        if (raised == 0 || raised == 1) {
            return int(raised);
        }
        else if (raised == p - 1) {
            return -1;
        }

        require(false);
    }
}
