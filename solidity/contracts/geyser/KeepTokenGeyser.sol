pragma solidity 0.5.17;

import "./TokenGeyser.sol";
import "../KeepToken.sol";

// TODO: Rename to KeepVault?
// TODO: Add docs
contract KeepTokenGeyser is TokenGeyser {
    constructor(
        KeepToken stakingToken,
        KeepToken distributionToken,
        uint256 maxUnlockSchedules,
        uint256 startBonus_,
        uint256 bonusPeriodSec_,
        uint256 initialSharesPerToken
    )
        public
        TokenGeyser(
            stakingToken,
            distributionToken,
            maxUnlockSchedules,
            startBonus_,
            bonusPeriodSec_,
            initialSharesPerToken
        )
    {}

    // TODO: Add functions that may be required from the Token Dashboard perspective.
}
