pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./libraries/grant/UnlockingSchedule.sol";
import "./GrantStakingPolicy.sol";

/// @title EmployeeStakingPolicy
/// @dev A staking policy which allows the grantee
/// to always stake the defined minimum stake,
/// or the unlocked amount if greater.
contract EmployeeStakingPolicy is GrantStakingPolicy {
    using SafeMath for uint256;
    using UnlockingSchedule for uint256;
    uint256 minimumStake;

    constructor(uint256 _minimumStake) public {
        minimumStake = _minimumStake;
    }

    function getStakeableAmount (
        uint256 _now,
        uint256 grantedAmount,
        uint256 duration,
        uint256 start,
        uint256 cliff,
        uint256 withdrawn
    ) public view returns (uint256) {
        uint256 unlocked = _now.getUnlockedAmount(
            grantedAmount,
            duration,
            start,
            cliff
        );
        uint256 remainingInGrant = grantedAmount.sub(withdrawn);
        uint256 unlockedInGrant = unlocked.sub(withdrawn);

        // Less than minimum stake remaining
        //   -> may stake what is remaining in grant
        if (remainingInGrant < minimumStake) { return remainingInGrant; }
        // At least minimum stake remaining in grant,
        // but unlocked amount is less than the minimum stake
        //   -> may stake the minimum stake
        if (unlockedInGrant < minimumStake) { return minimumStake; }
        // More than minimum stake unlocked in grant
        //   -> may stake the unlocked amount
        return unlockedInGrant;
    }
}
