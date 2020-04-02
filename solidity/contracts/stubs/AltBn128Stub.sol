pragma solidity ^0.5.4;

// it is just a stub, not a live deployment;
// we are fine with experimental feature
pragma experimental ABIEncoderV2;

import "../cryptography/AltBn128.sol";

contract AltBn128Stub {

    function publicG1Unmarshal(bytes memory m) public pure returns(AltBn128.G1Point memory) {
        return AltBn128.g1Unmarshal(m);
    }

    function publicG2Unmarshal(bytes memory m) public pure returns(AltBn128.G2Point memory) {
        return AltBn128.g2Unmarshal(m);
    }

    function publicG2Decompress(bytes memory m) public pure returns(AltBn128.G2Point memory) {
        return AltBn128.g2Decompress(m);
    }
}