pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

library MinimumStakeSchedule {
    using SafeMath for uint256;

    uint256 public constant schedule = 86400 * 365 * 2; // 2 years in seconds (seconds per day * days in a year * years)
    uint256 public constant steps = 10;
    uint256 public constant base = 10000 * 1e18;

    function current(uint256 scheduleStart) public view returns (uint256) {
        if (now < scheduleStart.add(schedule)) {
            uint256 currentStep = steps.mul(now.sub(scheduleStart)).div(schedule);
            return base.mul(steps.sub(currentStep));
        }
        return base;
    }
}