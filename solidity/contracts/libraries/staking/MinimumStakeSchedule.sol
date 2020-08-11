pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

/// @notice MinimumStakeSchedule defines the minimum stake parametrization and
/// schedule. It starts with a minimum stake of 100k KEEP. Over the following
/// 2 years, starting from the moment KEEP token has been deployed, the minimum
/// stake is lowered periodically using a uniform stepwise function, eventually
/// ending at 10k.
library MinimumStakeSchedule {
    using SafeMath for uint256;

    // Apr-28-2020 02:52:46 AM UTC when KEEP token has been deployed
    // TX:  0xea22d72bc7de4c82798df7194734024a1f2fd57b173d0e065864ff4e9d3dc014
    uint256 public constant scheduleStart = 1588042366;
    
    // 2 years in seconds (seconds per day * days in a year * years)
    uint256 public constant schedule = 86400 * 365 * 2;
    uint256 public constant steps = 10;
    uint256 public constant base = 10000 * 1e18;

    /// @notice Returns the current value of the minimum stake. The minimum
    /// stake is lowered periodically over the course of 2 years since the time
    /// KEEP token has been deployed and eventually ends at 10k KEEP.
    function current() public view returns (uint256) {
        return current(scheduleStart);
    }

    function current(uint256 scheduleStart) internal view returns (uint256) {
        if (now < scheduleStart.add(schedule)) {
            uint256 currentStep = steps.mul(now.sub(scheduleStart)).div(schedule);
            return base.mul(steps.sub(currentStep));
        }
        return base;
    }
}