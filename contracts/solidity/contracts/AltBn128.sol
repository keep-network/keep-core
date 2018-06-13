pragma solidity ^0.4.21;
pragma experimental ABIEncoderV2;

import "./utils/ModUtils.sol";

/**
 * @title Operations on alt_bn128
 * @dev Implementations of common elliptic curve operations on Ethereum's
 * (poorly named) alt_bn128 curve. Whenever possible, use post-Byzantium
 * pre-compiled contracts to offset gas costs. Note that these pre-compiles
 * might not be available on all (eg private) chains.
 */
library AltBn128 {

    using ModUtils for uint256;

    uint256 constant p = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    function getP() public constant returns (uint256) {
        return p;
    }

    function yFromX(uint256 x)
        private
        constant returns(uint256)
    {
        return ((x.modExp(3, p) + 3) % p).modSqrt(p);
    }

    /**
     * @dev Hash a byte array message, m, and map it deterministically to a
     * point on G1. Note that this approach was chosen for its simplicity /
     * lower gas cost on the EVM, rather than good distribution of points on
     * G1.
     */
    function g1HashToPoint(bytes m)
        public
        constant returns(uint256, uint256)
    {
        bytes32 h = sha256(m);
        uint256 x = uint256(h) % p;
        uint256 y;

        while (true) {
            y = yFromX(x);
            if (y > 0) {
                return (x, y);
            }
            x += 1;
        }
    }

    /**
     * @dev Compress a point on G1 to a single uint256 for serialization.
     */
    function g1Compress(uint256 x, uint256 y)
        public
        constant returns(bytes32)
    {
        bytes32 m = bytes32(x);

        byte leadM = m[0] | ((bytes32(y)[31] & byte(1)) << 7);

        assembly {
            mstore(add(m, 1), leadM)
        }

        return m;
    }

    /**
     * @dev Decompress a point on G1 from a single uint256.
     */
    function g1Decompress(bytes32 m)
        public
        constant returns(uint256, uint256)
    {
        byte ySign = (m[0] ^ byte(0x10000000)) >> 7;
        bytes32 mX = bytes32(0);
        byte leadX = mX[0] & byte(0x01111111);

        assembly {
            mstore(add(mX, 32), m)
            mstore(add(mX, 1), leadX)
        }

        uint256 x = uint256(mX);
        uint256 y = yFromX(x);

        if (ySign != (bytes32(y)[0] ^ byte(0x10000000)) >> 7) {
            y = p - y;
        }

        return (x, y);

    }

    /**
     * @dev Wrap the point addition pre-compile introduced in Byzantium. Return
     * the sum of two points on G1. Revert if the provided points aren't on the
     * curve.
     */
    function add(uint256[2] a, uint256[2] b) public constant returns (uint256, uint256) {
        uint256[4] memory arg;
        arg[0] = a[0];
        arg[1] = a[1];
        arg[2] = b[0];
        arg[3] = b[1];
        uint256[2] memory c;
        assembly {
            if iszero(call(not(0), 0x06, 0, arg, 0x80, c, 0x40)) {
                revert(0, 0)
            }
        }
        return (c[0], c[1]);
    }

    /**
     * @dev Wrap the scalar point multiplication pre-compile introduced in
     * Byzantium. The result of a point from G1 multiplied by a scalar should
     * match the point added to itself the same number of times. Revert if the
     * provided point isn't on the curve.
     */
    function scalarMultiply(uint256[2] p_1, uint256 scalar) public constant returns (uint256, uint256) {
        uint256[3] memory arg;
        arg[0] = p_1[0];
        arg[1] = p_1[1];
        arg[2] = scalar;
        uint256[2] memory p_2;
        assembly {
            if iszero(call(not(0), 0x07, 0, arg, 0x60, p_2, 0x40)) {
                revert(0, 0)
            }
        }
        return (p_2[0], p_2[1]);
    }
}
