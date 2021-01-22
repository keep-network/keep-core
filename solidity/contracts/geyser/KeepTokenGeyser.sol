pragma solidity 0.5.17;

import "./TokenGeyser.sol";
import "../KeepToken.sol";

import "openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";

// TODO: Rename to KeepVault?
// TODO: Add docs
contract KeepTokenGeyser is TokenGeyser {
    /**
     * @param _stakingToken The token users deposit as a stake.
     * @param _distributionToken The token users are rewarded in and receive it as they unstake.
     * @param _maxUnlockSchedules Max number of unlock stages, to guard against hitting gas limit.
     * @param _startBonus Starting time bonus, BONUS_DECIMALS fixed point.
     *                    e.g. 25% means user gets 25% of max distribution tokens.
     * @param _bonusPeriodSec Length of time for bonus to increase linearly to max.
     * @param _initialSharesPerToken Number of shares to mint per staking token on first stake.
     */
    constructor(
        IERC20 _stakingToken,
        KeepToken _distributionToken,
        uint256 _maxUnlockSchedules,
        uint256 _startBonus,
        uint256 _bonusPeriodSec,
        uint256 _initialSharesPerToken
    )
        public
        TokenGeyser(
            _stakingToken,
            _distributionToken,
            _maxUnlockSchedules,
            _startBonus,
            _bonusPeriodSec,
            _initialSharesPerToken
        )
    {}

    // TODO: Add functions that may be required from the Token Dashboard perspective.
}
