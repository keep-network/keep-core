pragma solidity 0.5.17;

import "../TokenStaking.sol";
import "../TokenStakingEscrow.sol";
import "../TokenGrant.sol";
import "../KeepRegistry.sol";

/// @dev TokenStakingStub keeps the same minimum stake value of 100k KEEP for
/// all the time. This stub is used by tests for which we want to maintain the
/// same minimum stake value. MinimumStakeSchedule uses the time of deploying KEEP
/// token as the starting point of the minimum stake schedule.
/// Use this stub to keep it at the same, predictable level. Going back in time
/// in tests is not possible.
contract TokenStakingStub is TokenStaking {
    constructor(
        ERC20Burnable _token,
        TokenGrant _tokenGrant,
        TokenStakingEscrow _escrow,
        KeepRegistry _registry,
        uint256 _initializationPeriod
    ) TokenStaking(_token, _tokenGrant, _escrow, _registry, _initializationPeriod) public {
    }

    function minimumStake() public view returns (uint256) {
        return 100000 * 1e18;
    }
}
