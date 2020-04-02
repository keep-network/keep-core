pragma solidity ^0.5.4;

import "../TokenStaking.sol";

contract TokenStakingSlashingStub is TokenStaking {
    constructor(
        address _tokenAddress,
        address _registry,
        uint256 _initializationPeriod,
        uint256 _undelegationPeriod
    ) TokenStaking(_tokenAddress, _registry, _initializationPeriod, _undelegationPeriod) public {
    }

    function slash(uint256 amountToSlash, address[] memory misbehavedOperators) public {
        for (uint i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            emit TokensSlashed(operator, 1 ether);
        }
    }

    function seize(
        uint256 amountToSeize,
        uint256 rewardMultiplier,
        address tattletale,
        address[] memory misbehavedOperators
    ) public {
        for (uint i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            emit TokensSeized(operator, 1 ether);
        }
    }
}
