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

    // p is a prime over which we form a basic field
    // Taken from go-ethereum/crypto/bn256/cloudflare/constants.go
    uint256 constant p = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    function getP() public pure returns (uint256) {
        return p;
    }

    /**
     * @dev Gets generator of G1 group.
     * Taken from go-ethereum/crypto/bn256/cloudflare/curve.go
     */
    function g1() public pure returns (uint256[2]) {
        return [uint256(1), uint256(2)];
    }

    /**
     * @dev Gets generator of G2 group.
     * Taken from go-ethereum/crypto/bn256/cloudflare/twist.go
     */
    function g2() public pure returns (uint256[4]) {
        return [
            11559732032986387107991004021392285783925812861821192530917403151452391805634,
            10857046999023057135944570762232829481370756359578518086990519993285655852781,
            4082367875863433681332203403145435568316851327593401208105741076214120093531,
            8495653923123431417604973247489272438418190587263600148770280649306958101930
        ];
    }

    /**
     * @dev Gets twist curve B constant.
     * Taken from go-ethereum/crypto/bn256/cloudflare/twist.go
     */
    function twistB() public pure returns (uint256[2]) {
        return [
            19485874751759354771024239261021720505790618469301721065564631296452457478373,
            266929791119991161246907387137283842545076965332900288569378510910307636690
        ];
    }

    /**
     * @dev Gets root of the point where x and y are equal.
     */
    function hexRoot() public pure returns (uint256[2]) {
        return [
            21573744529824266246521972077326577680729363968861965890554801909984373949499,
            16854739155576650954933913186877292401521110422362946064090026408937773542853
        ];
    }

    /**
     * @dev yFromX computes a Y value for a point based on an X value. This
     * computation is simply evaluating the curve equation for Y on a
     * given X, and allows a point on the curve to be represented by just
     * an X value + a sign bit.
     */
    function yFromX(uint256 x)
        private
        view returns(uint256)
    {
        return ((x.modExp(3, p) + 3) % p).modSqrt(p);
    }

    /**
     * @dev gfP2YFromX computes a Y value for a gfP2 point based on an X value.
     * This computation is simply evaluating the curve equation for Y on a
     * given X, and allows a point on the curve to be represented by just
     * an X value + a sign bit.
     */
    function gfP2YFromX(uint256[2] _x)
        private
        pure returns(uint256[2] y)
    {
        uint256[2] memory x = gfP2Add(gfP2Pow(_x, 3), twistB());

        // Using formula y = x ^ (p^2 + 15) // 32) from 
        // https://github.com/ethereum/beacon_chain/blob/master/beacon_chain/utils/bls.py
        // (p^2 + 15) // 32) results into a big 512bit value, so breaking it to two uint256 as (a * a + b)
        uint256 a = 3869331240733915743250440106392954448556483137451914450067252501901456824595;
        uint256 b = 146360017852723390495514512480590656176144969185739259173561346299185050597;
  
        y = gfP2Multiply(gfP2Pow(gfP2Pow(x, a), a), gfP2Pow(x, b));

        // Multiply y by hexRoot constant to find correct y.
        while (!x2y(x, y)) {
            y = gfP2Multiply(y, hexRoot());
        }
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
     * @dev Calculates whether the provided number is even or odd.
     * @return 0x01 if y is an even number and 0x00 if it's odd.
     */
    function parity(uint256 value) private pure returns (byte) {
        return bytes32(value)[31] & byte(1);
    }

    /**
     * @dev Compress a point on G1 to a single uint256 for serialization.
     */
    function g1Compress(uint256 x, uint256 y)
        public
        pure returns(bytes32)
    {
        bytes32 m = bytes32(x);

        byte leadM = m[0] | parity(y) << 7;
        bytes32 mask = 0xff << 31*8;
        m = (m & ~mask) | (leadM >> 0);

        return m;
    }

    /**
     * @dev Compress a point on G2 to a pair of uint256 for serialization.
     */
    function g2Compress(uint256[2] x, uint256[2] y)
        public
        pure returns(bytes)
    {
        bytes32 m = bytes32(x[0]);

        byte leadM = m[0] | parity(y[0]) << 7;
        bytes32 mask = 0xff << 31*8;
        m = (m & ~mask) | (leadM >> 0);

        return abi.encodePacked(m, bytes32(x[1]));
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

        if (parity(y) != (m[0] & byte(128)) >> 7) {
            y = p - y;
        }

        return (x, y);
    }

    /**
     * @dev Decompress a point on G2 from a pair of uint256.
     */
    function g2Decompress(bytes m)
        public
        view returns(uint256[2], uint256[2])
    {
        bytes32 x1;
        bytes32 x2;
        uint256 temp;

        // Extract two bytes32 from bytes array
        /* solium-disable-next-line */
        assembly {
            temp := add(m, 32)
            x1 := mload(temp)
            temp := add(m, 64)
            x2 := mload(temp)
        }

        bytes32 mX = bytes32(0);
        byte leadX = x1[0] & byte(127);
        bytes32 mask = 0xff << 31*8;
        mX = (x1 & ~mask) | (leadX >> 0);

        uint256[2] memory x = [uint256(x2), uint256(mX)];
        uint256[2] memory y = gfP2YFromX(x);
        y = [y[1], y[0]];
        x = [x[1], x[0]];

        if (parity(y[0]) != (m[0] & byte(128)) >> 7) {
            y[0] = p - y[0];
            y[1] = p - y[1];
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
     * @dev Return the sum of two gfP2 points.
     */
    function gfP2Add(uint256[2] a, uint256[2] b) internal pure returns(uint256[2]) {
        return (
            [addmod(a[0], b[0], p),
            addmod(a[1], b[1], p)]
        );
    }

    /**
     * @dev Return multiplication of two gfP2 points.
     */
    function gfP2Multiply(uint256[2] a, uint256[2] b) internal pure returns(uint256[2]) {
        return (
            [addmod(mulmod(a[0], b[0], p), p - mulmod(a[1], b[1], p), p),
            addmod(mulmod(a[0], b[1], p), mulmod(a[1], b[0], p), p)]
        );
    }

    /**
     * @dev Return gfP2 element to the power of the provided exponent.
     */
    function gfP2Pow(uint256[2] _a, uint256 _exp) internal pure returns(uint256[2] result) {
        uint256 exp = _exp;
        uint256[2] memory a;
        result[0] = 1;
        result[1] = 0;
        a[0] = _a[0];
        a[1] = _a[1];

        // Reduce exp dividing by 2 gradually to 0 while computing final
        // result only when exp is an odd number.
        while (exp > 0) {
            if (parity(exp) == 1) {
                result = gfP2Multiply(result, a);
            }

            exp = exp / 2;
            a = gfP2Multiply(a, a);
        }
    }

    /**
     * @dev Return true if y^2 equals x.
     */
    function x2y(uint256[2] x, uint256[2] y) internal pure returns(bool) {
       
        uint256[2] memory y2;
        y2 = gfP2Pow(y, 2);

        return (y2[0] == x[0] && y2[1] == x[1]);
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
     * @dev Wrap the pairing check pre-compile introduced in Byzantium. Return
     * the result of a pairing check of 2 pairs (G1 p1, G2 p2) (G1 p3, G2 p4)
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
