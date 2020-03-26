pragma solidity ^0.5.4;

import "../TokenStaking.sol";

contract TokenStakingStub is TokenStaking {
    constructor(
        address _tokenAddress,
        address _registry,
        uint256 _initializationPeriod,
        uint256 _undelegationPeriod
    ) TokenStaking(_tokenAddress, _registry, _initializationPeriod, _undelegationPeriod) public {
    }

    function setInitializationPeriod(uint256 _initializationPeriod) public {
        initializationPeriod = _initializationPeriod;
    }
}
