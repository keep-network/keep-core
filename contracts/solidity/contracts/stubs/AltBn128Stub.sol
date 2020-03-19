pragma solidity ^0.5.4;

import "../cryptography/AltBn128.sol";

library AltBn128Stub {
    function g1HashToPoint(bytes memory message) public view returns(bytes memory) {
        AltBn128.G1Point memory point = AltBn128.g1HashToPoint(message);
        return AltBn128.g1Marshal(point);
    }
}
