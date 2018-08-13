pragma solidity ^0.4.21;

library ModUtils {

    function modExp(uint256 a, uint256 exponent, uint256 p)
        public
        view returns(uint256 raised)
    {
        uint256[6] memory args = [32, 32, 32, a, exponent, p];
        uint256[1] memory output;
        /* solium-disable-next-line */
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
        view returns(uint256)
    {

        if (legendre(a, p) != 1) {
            return 0;
        }

        if (a == 0) {
            return 0;
        }

        if (p == 2) {
            return p;
        }

        if (p % 4 == 3) {
            return modExp(a, (p + 1) / 4, p);
        }

        uint256 s = p - 1;
        // log_2(256) = 8 => uint8
        uint8 e = 0;

        while (s % 2 == 0) {
            s = s / 2;
            e = e + 1;
        }

        // Note the smaller int- finding n with Legendre symbol or -1
        // should be quick
        uint32 n = 2;
        while (legendre(n, p) != -1) {
            n = n + 1;
        }

        uint256 x = modExp(a, (s + 1) / 2, p);
        uint256 b = modExp(a, s, p);
        uint256 g = modExp(n, s, p);
        uint8 r = e;
        uint256 gs = 0;
        uint8 m = 0;
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

            gs = modExp(g, uint256(2) ** (r - m - 1), p);
            g = (gs * gs) % p;
            x = (x * gs) % p;
            b = (b * g) % p;
            r = m;
        }
    }


    function legendre(uint256 a, uint256 p)
        public
        view returns(int8)
    {
        uint256 raised = modExp(a, (p - 1) / uint256(2), p);

        if (raised == 0 || raised == 1) {
            return int8(raised);
        }
        else if (raised == p - 1) {
            return -1;
        }

        require(false);
    }
}
