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

    function getP() public view returns (uint256) {
        return p;
    }

    /**
     * @dev Gets generator of G1 group.
     */
    function g1() public view returns (uint256[2]) {
        return [uint256(1), uint256(2)];
    }

    /**
     * @dev Gets generator of G2 group.
     */
    function g2() public view returns (uint256[4]) {
        return [
            11559732032986387107991004021392285783925812861821192530917403151452391805634,
            10857046999023057135944570762232829481370756359578518086990519993285655852781,
            4082367875863433681332203403145435568316851327593401208105741076214120093531,
            8495653923123431417604973247489272438418190587263600148770280649306958101930
        ];
    }

    function yFromX(uint256 x)
        private
        view returns(uint256)
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
        view returns(uint256, uint256)
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
     * @dev Calculates whether provided y coordinate is even or odd number.
     * @return 0x01 byte if y is an even number and 0x00 if it's odd.
     */
    function ySign(uint256 y) private returns (byte) {
        return bytes32(y)[31] & byte(1);
    }

    /**
     * @dev Compress a point on G1 to a single uint256 for serialization.
     */
    function g1Compress(uint256 x, uint256 y)
        public
        view returns(bytes32)
    {
        bytes32 m = bytes32(x);

        byte leadM = m[0] | ySign(y) << 7;
        bytes32 mask = 0xff << 31*8;
        m = (m & ~mask) | (leadM >> 0);

        return m;
    }

    /**
     * @dev Decompress a point on G1 from a single uint256.
     */
    function g1Decompress(bytes32 m)
        public
        view returns(uint256, uint256)
    {
        bytes32 mX = bytes32(0);
        byte leadX = m[0] & byte(127);
        bytes32 mask = 0xff << 31*8;
        mX = (m & ~mask) | (leadX >> 0);

        uint256 x = uint256(mX);
        uint256 y = yFromX(x);

        if (ySign(y) != (m[0] & byte(128)) >> 7) {
            y = p - y;
        }

        return (x, y);

    }

    /**
     * @dev Wrap the point addition pre-compile introduced in Byzantium. Return
     * the sum of two points on G1. Revert if the provided points aren't on the
     * curve.
     */
    function add(uint256[2] a, uint256[2] b) public view returns (uint256, uint256) {
        uint256[4] memory arg;
        arg[0] = a[0];
        arg[1] = a[1];
        arg[2] = b[0];
        arg[3] = b[1];
        uint256[2] memory c;
        /* solium-disable-next-line */
        assembly {
            // 0x60 is the ECADD precompile address
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
    function scalarMultiply(uint256[2] p_1, uint256 scalar) public view returns (uint256, uint256) {
        uint256[3] memory arg;
        arg[0] = p_1[0];
        arg[1] = p_1[1];
        arg[2] = scalar;
        uint256[2] memory p_2;
        /* solium-disable-next-line */
        assembly {
            // 0x70 is the ECMUL precompile address
            if iszero(call(not(0), 0x07, 0, arg, 0x60, p_2, 0x40)) {
                revert(0, 0)
            }
        }
        return (p_2[0], p_2[1]);
    }

    /**
     * @dev Wrap the bn256Pairing pre-compile introduced in Byzantium. Return
     * the result of a pairing check of 4 pairs (G1 p1, G2 p2, G1 p3, G2 p4)
     */
    function pairing(uint256[2] p1, uint256[4] p2, uint256[2] p3, uint256[4] p4) public view returns (bool) {
        uint256[12] memory arg = [
            p1[0], p1[1], p2[0], p2[1], p2[2], p2[3], p3[0], p3[1], p4[0], p4[1], p4[2], p4[3]
        ];
        uint[1] memory c;
        /* solium-disable-next-line */
        assembly {
            // call(gasLimit, to, value, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(call(not(0), 0x08, 0, arg, 0x180, c, 0x20)) {
                revert(0, 0)
            }
        }
        return c[0] != 0;
    }
}
