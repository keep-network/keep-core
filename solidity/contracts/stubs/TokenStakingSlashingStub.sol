pragma solidity 0.5.17;

import "../TokenStaking.sol";
import "../TokenStakingEscrow.sol";
import "../TokenGrant.sol";
import "../KeepRegistry.sol";

contract TokenStakingSlashingStub is TokenStaking {
    constructor(
        ERC20Burnable _token,
        TokenGrant _tokenGrant,
        TokenStakingEscrow _escrow,
        KeepRegistry _registry,
        uint256 _initializationPeriod
    ) TokenStaking(_token, _tokenGrant, _escrow, _registry, _initializationPeriod) public {
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
