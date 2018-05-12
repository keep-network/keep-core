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

    struct G1 {
        uint256 x;
        uint256 y;
    }

    /**
     * @dev Hash a byte array message, m, and map it deterministically to a
     * point on G1. Note that this approach was chosen for its simplicity /
     * lower gas cost on the EVM, rather than good distribution of points on
     * G1.
     */
    function g1HashToPoint(bytes m)
        public
        constant returns(G1)
    {
        bytes32 h = sha256(m);
        uint256 x = uint256(h) % p;
        uint256 y = 0;

        while (true) {
            y = (x ** 3 + 3).modSqrt(p);
            if (y > 0) {
                return G1(x, y);
            }
            x += 1;
        }
    }
}
