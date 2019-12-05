pragma solidity ^0.5.4;

import "../utils/ModUtils.sol";

/**
 * @title Operations on alt_bn128
 * @dev Implementations of common elliptic curve operations on Ethereum's
 * (poorly named) alt_bn128 curve. Whenever possible, use post-Byzantium
 * pre-compiled contracts to offset gas costs. Note that these pre-compiles
 * might not be available on all (eg private) chains.
 */
library AltBn128 {

    using ModUtils for uint256;

    // G1Point implements a point in G1 group.
    struct G1Point {
        uint256 x;
        uint256 y;
    }

    // gfP2 implements a field of size p² as a quadratic extension of the base field.
    struct gfP2 {
        uint256 x;
        uint256 y;
    }

    // G2Point implements a point in G2 group.
    struct G2Point {
        gfP2 x;
        gfP2 y;
    }

    // p is a prime over which we form a basic field
    // Taken from go-ethereum/crypto/bn256/cloudflare/constants.go
    uint256 constant p = 21888242871839275222246405745257275088696311157297823662689037894645226208583;

    function getP() internal pure returns (uint256) {
        return p;
    }

    /**
     * @dev Gets generator of G1 group.
     * Taken from go-ethereum/crypto/bn256/cloudflare/curve.go
     */
    function g1() internal pure returns (G1Point memory) {
        return G1Point(uint256(1), uint256(2));
    }

    /**
     * @dev Gets generator of G2 group.
     * Taken from go-ethereum/crypto/bn256/cloudflare/twist.go
     */
    function g2() internal pure returns (G2Point memory) {
        return G2Point(
            gfP2(
                11559732032986387107991004021392285783925812861821192530917403151452391805634,
                10857046999023057135944570762232829481370756359578518086990519993285655852781
            ),
            gfP2(
                4082367875863433681332203403145435568316851327593401208105741076214120093531,
                8495653923123431417604973247489272438418190587263600148770280649306958101930
            )
        );
    }

    /**
     * @dev Gets twist curve B constant.
     * Taken from go-ethereum/crypto/bn256/cloudflare/twist.go
     */
    function twistB() private pure returns (gfP2 memory) {
        return gfP2(
            266929791119991161246907387137283842545076965332900288569378510910307636690,
            19485874751759354771024239261021720505790618469301721065564631296452457478373
        );
    }

    /**
     * @dev Gets root of the point where x and y are equal.
     */
    function hexRoot() private pure returns (gfP2 memory) {
        return gfP2(
            21573744529824266246521972077326577680729363968861965890554801909984373949499,
            16854739155576650954933913186877292401521110422362946064090026408937773542853
        );
    }

    /**
     * @dev g1YFromX computes a Y value for a G1 point based on an X value.
     * This computation is simply evaluating the curve equation for Y on a
     * given X, and allows a point on the curve to be represented by just
     * an X value + a sign bit.
     */
    function g1YFromX(uint256 x)
        internal
        view returns(uint256)
    {
        return ((x.modExp(3, p) + 3) % p).modSqrt(p);
    }

    /**
     * @dev g2YFromX computes a Y value for a G2 point based on an X value.
     * This computation is simply evaluating the curve equation for Y on a
     * given X, and allows a point on the curve to be represented by just
     * an X value + a sign bit.
     */
    function g2YFromX(gfP2 memory _x)
        internal
        pure returns(gfP2 memory y)
    {
        gfP2 memory x = gfP2Add(gfP2Pow(_x, 3), twistB());

        // Using formula y = x ^ (p^2 + 15) / 32 from
        // https://github.com/ethereum/beacon_chain/blob/master/beacon_chain/utils/bls.py
        // (p^2 + 15) / 32 results into a big 512bit value, so breaking it to two uint256 as (a * a + b)
        uint256 a = 3869331240733915743250440106392954448556483137451914450067252501901456824595;
        uint256 b = 146360017852723390495514512480590656176144969185739259173561346299185050597;
  
        y = gfP2Multiply(gfP2Pow(gfP2Pow(x, a), a), gfP2Pow(x, b));
        
        // Multiply y by hexRoot constant to find correct y.
        while (!g2X2y(x, y)) {
            y = gfP2Multiply(y, hexRoot());
        }
    }

    /**
     * @dev Hash a byte array message, m, and map it deterministically to a
     * point on G1. Note that this approach was chosen for its simplicity /
     * lower gas cost on the EVM, rather than good distribution of points on
     * G1.
     */
    function g1HashToPoint(bytes memory m)
        internal
        view returns(G1Point memory)
    {
        bytes32 h = sha256(m);
        uint256 x = uint256(h) % p;
        uint256 y;

        while (true) {
            y = g1YFromX(x);
            if (y > 0) {
                return G1Point(x, y);
            }
            x += 1;
        }
    }

    /**
     * @dev Calculates whether the provided number is even or odd.
     * @return 0x01 if y is an even number and 0x00 if it's odd.
     */
    function parity(uint256 value) private pure returns (byte) {
        return bytes32(value)[31] & 0x01;
    }

    /**
     * @dev Compress a point on G1 to a single uint256 for serialization.
     */
    function g1Compress(G1Point memory point)
        internal
        pure returns(bytes32)
    {
        bytes32 m = bytes32(point.x);

        byte leadM = m[0] | parity(point.y) << 7;
        uint256 mask = 0xff << 31*8;
        m = (m & ~bytes32(mask)) | (leadM >> 0);

        return m;
    }

    /**
     * @dev Compress a point on G2 to a pair of uint256 for serialization.
     */
    function g2Compress(G2Point memory point)
        internal
        pure returns(bytes memory)
    {
        bytes32 m = bytes32(point.x.x);

        byte leadM = m[0] | parity(point.y.x) << 7;
        uint256 mask = 0xff << 31*8;
        m = (m & ~bytes32(mask)) | (leadM >> 0);

        return abi.encodePacked(m, bytes32(point.x.y));
    }

    /**
     * @dev Decompress a point on G1 from a single uint256.
     */
    function g1Decompress(bytes32 m)
        internal
        view returns(G1Point memory)
    {
        bytes32 mX = bytes32(0);
        byte leadX = m[0] & 0x7f;
        uint256 mask = 0xff << 31*8;
        mX = (m & ~bytes32(mask)) | (leadX >> 0);

        uint256 x = uint256(mX);
        uint256 y = g1YFromX(x);

        if (parity(y) != (m[0] & 0x80) >> 7) {
            y = p - y;
        }

        require(isG1PointOnCurve(G1Point(x, y)), "Malformed bn256.G1 point.");

        return G1Point(x, y);
    }

    /**
     * @dev Unmarshals a point on G1 from bytes in an uncompressed form.
     */
    function g1Unmarshal(bytes memory m) internal pure returns(G1Point memory) {
        bytes32 x;
        bytes32 y;

        /* solium-disable-next-line */
        assembly {
            x := mload(add(m, 0x20))
            y := mload(add(m, 0x40))
        }

        return G1Point(uint256(x), uint256(y));
    }

    /**
     * @dev Unmarshals a point on G2 from bytes in an uncompressed form.
     */
    function g2Unmarshal(bytes memory m) internal pure returns(G2Point memory) {
        bytes32 xx;
        bytes32 xy;
        bytes32 yx;
        bytes32 yy;

        /* solium-disable-next-line */
        assembly {
            xx := mload(add(m, 0x20))
            xy := mload(add(m, 0x40))
            yx := mload(add(m, 0x60))
            yy := mload(add(m, 0x80))
        }

        return G2Point(gfP2(uint256(xx), uint256(xy)), gfP2(uint256(yx),uint256(yy)));
    }

    /**
     * @dev Decompress a point on G2 from a pair of uint256.
     */
    function g2Decompress(bytes memory m)
        internal
        pure returns(G2Point memory)
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
        byte leadX = x1[0] & 0x7f;
        uint256 mask = 0xff << 31*8;
        mX = (x1 & ~bytes32(mask)) | (leadX >> 0);

        gfP2 memory x = gfP2(uint256(mX), uint256(x2));
        gfP2 memory y = g2YFromX(x);

        if (parity(y.x) != (m[0] & 0x80) >> 7) {
            y.x = p - y.x;
            y.y = p - y.y;
        }

        require(isG2PointOnCurve(G2Point(x, y)), "Malformed bn256.G2 point.");
        return G2Point(x, y);
    }

    /**
     * @dev Wrap the point addition pre-compile introduced in Byzantium. Return
     * the sum of two points on G1. Revert if the provided points aren't on the
     * curve.
     */
    function g1Add(G1Point memory a, G1Point memory b) internal view returns (G1Point memory) {
        uint256[4] memory arg;
        arg[0] = a.x;
        arg[1] = a.y;
        arg[2] = b.x;
        arg[3] = b.y;
        uint256[2] memory c;

        /* solium-disable-next-line */
        assembly {
            // 0x60 is the ECADD precompile address
            if iszero(staticcall(not(0), 0x06, arg, 0x80, c, 0x40)) {
                revert(0, 0)
            }
        }

        return G1Point(c[0], c[1]);
    }

    /**
     * @dev Return the sum of two gfP2 field elements.
     */
    function gfP2Add(gfP2 memory a, gfP2 memory b) internal pure returns(gfP2 memory) {
        return gfP2(
            addmod(a.x, b.x, p),
            addmod(a.y, b.y, p)
        );
    }

    /**
     * @dev Return multiplication of two gfP2 field elements.
     */
    function gfP2Multiply(gfP2 memory a, gfP2 memory b) internal pure returns(gfP2 memory) {
        return gfP2(
            addmod(mulmod(a.x, b.y, p), mulmod(b.x, a.y, p), p),
            addmod(mulmod(a.y, b.y, p), p - mulmod(a.x, b.x, p), p)
        );
    }

    /**
     * @dev Return gfP2 element to the power of the provided exponent.
     */
    function gfP2Pow(gfP2 memory _a, uint256 _exp) internal pure returns(gfP2 memory result) {
        uint256 exp = _exp;
        gfP2 memory a;
        result.x = 0;
        result.y = 1;
        a.x = _a.x;
        a.y = _a.y;

        // Reduce exp dividing by 2 gradually to 0 while computing final
        // result only when exp is an odd number.
        while (exp > 0) {
            if (parity(exp) == 0x01) {
                result = gfP2Multiply(result, a);
            }

            exp = exp / 2;
            a = gfP2Multiply(a, a);
        }
    }

    /**
     * @dev Return true if G2 point's y^2 equals x.
     */
    function g2X2y(gfP2 memory x, gfP2 memory y) internal pure returns(bool) {
       
        gfP2 memory y2;
        y2 = gfP2Pow(y, 2);

        return (y2.x == x.x && y2.y == x.y);
    }

    /**
     * @dev Return true if G1 point is on the curve.
     */
    function isG1PointOnCurve(G1Point memory point) internal view returns (bool) {
        return point.y.modExp(2, p) == (point.x.modExp(3, p) + 3) % p;
    }

    /**
     * @dev Return true if G2 point is on the curve.
     */
    function isG2PointOnCurve(G2Point memory point) internal pure returns(bool) {

        gfP2 memory y2;
        gfP2 memory x3;

        y2 = gfP2Pow(point.y, 2);
        x3 = gfP2Add(gfP2Pow(point.x, 3), twistB());

        return (y2.x == x3.x && y2.y == x3.y);
    }

    /**
     * @dev Wrap the scalar point multiplication pre-compile introduced in
     * Byzantium. The result of a point from G1 multiplied by a scalar should
     * match the point added to itself the same number of times. Revert if the
     * provided point isn't on the curve.
     */
    function scalarMultiply(G1Point memory p_1, uint256 scalar) internal view returns (G1Point memory) {
        uint256[3] memory arg;
        arg[0] = p_1.x;
        arg[1] = p_1.y;
        arg[2] = scalar;
        uint256[2] memory p_2;
        /* solium-disable-next-line */
        assembly {
            // 0x70 is the ECMUL precompile address
            if iszero(staticcall(not(0), 0x07, arg, 0x60, p_2, 0x40)) {
                revert(0, 0)
            }
        }
        return G1Point(p_2[0], p_2[1]);
    }

    /**
     * @dev Wrap the pairing check pre-compile introduced in Byzantium. Return
     * the result of a pairing check of 2 pairs (G1 p1, G2 p2) (G1 p3, G2 p4)
     */
    function pairing(G1Point memory p1, G2Point memory p2, G1Point memory p3, G2Point memory p4) internal view returns (bool) {
        uint256[12] memory arg = [
            p1.x, p1.y, p2.x.x, p2.x.y, p2.y.x, p2.y.y, p3.x, p3.y, p4.x.x, p4.x.y, p4.y.x, p4.y.y
        ];
        uint[1] memory c;
        /* solium-disable-next-line */
        assembly {
            // call(gasLimit, to, value, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(not(0), 0x08, arg, 0x180, c, 0x20)) {
                revert(0, 0)
            }
        }
        return c[0] != 0;
    }

}
