pragma solidity ^0.5.4;

import "../cryptography/AltBn128.sol";

library AltBn128Stub {
    function sign(uint256 secretKey) public view returns(uint256, uint256) {
        return AltBn128.sign(secretKey);
    }
}
