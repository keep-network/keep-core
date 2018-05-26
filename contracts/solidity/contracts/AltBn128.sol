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
        uint256 y = 0;

        while (true) {
            y = (x ** 3 + 3).modSqrt(p);
            if (y > 0) {
                return (x, y);
            }
            x += 1;
        }
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
