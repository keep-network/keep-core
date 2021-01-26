pragma solidity 0.5.17;

import "./KeepTokenGeyser.sol";
import "../KeepToken.sol";

// TODO: Add docs
contract KeepVault is KeepTokenGeyser {
    /**
     * @param _keepToken KEEP token contract address. It is a token that is accepted
     * as user's stake and that will be distributed as rewards.
     * @param _maxUnlockSchedules Max number of unlock stages, to guard against hitting gas limit.
     * @param _startBonus Starting time bonus, BONUS_DECIMALS fixed point.
     *                    e.g. 25% means user gets 25% of max distribution tokens.
     * @param _bonusPeriodSec Length of time for bonus to increase linearly to max.
     * @param _initialSharesPerToken Number of shares to mint per staking token on first stake.
     * @param _durationSec TODO: Update docs
     */
    constructor(
        KeepToken _keepToken,
        uint256 _maxUnlockSchedules,
        uint256 _startBonus,
        uint256 _bonusPeriodSec,
        uint256 _initialSharesPerToken,
        uint256 _durationSec
    )
        public
        KeepTokenGeyser(
            _keepToken,
            _keepToken,
            _maxUnlockSchedules,
            _startBonus,
            _bonusPeriodSec,
            _initialSharesPerToken,
            _durationSec
        )
    {}
}
