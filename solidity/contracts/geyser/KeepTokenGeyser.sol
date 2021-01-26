pragma solidity 0.5.17;

import "./TokenGeyser.sol";
import "../KeepToken.sol";
import {IStakerRewards, StakerRewardsBeneficiary} from "../PhasedEscrow.sol";

import "openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";

// TODO: Rename to KeepVault?
// TODO: Add docs
contract KeepTokenGeyser is TokenGeyser, IStakerRewards {
    event DurationSecUpdated(uint256 oldDurationSec, uint256 newDurationSec);

    uint256 public durationSec;

    /**
     * @param _stakingToken The token users deposit as a stake.
     * @param _distributionToken The token users are rewarded in and receive it as they unstake.
     * @param _maxUnlockSchedules Max number of unlock stages, to guard against hitting gas limit.
     * @param _startBonus Starting time bonus, BONUS_DECIMALS fixed point.
     *                    e.g. 25% means user gets 25% of max distribution tokens.
     * @param _bonusPeriodSec Length of time for bonus to increase linearly to max.
     * @param _initialSharesPerToken Number of shares to mint per staking token on first stake.
     * @param _durationSec Length of time to linear unlock the rewards tokens.
     */
    constructor(
        IERC20 _stakingToken,
        KeepToken _distributionToken,
        uint256 _maxUnlockSchedules,
        uint256 _startBonus,
        uint256 _bonusPeriodSec,
        uint256 _initialSharesPerToken,
        uint256 _durationSec
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
    {
        durationSec = _durationSec;
    }

    /// @notice Receives an approval of token transfer and locks the approved
    /// amount for rewards.
    /// @dev Requires the provided token contract to be the same as the distribution
    /// token supported by this contract.
    /// @param _from The owner of the tokens who approved them to stake.
    /// @param _value Approved amount of tokens for the transfer.
    /// @param _token Token contract address.
    /// @param _extraData Ignored.
    function receiveApproval(
        address _from,
        uint256 _value,
        address _token,
        bytes calldata _extraData
    ) external {
        require(
            KeepToken(_token) == getDistributionToken(),
            "Token is not supported distribution token"
        );

        lockTokens(_value, durationSec);
    }

    function setDurationSec(uint256 _newDurationSec) external onlyOwner {
        require(
            _newDurationSec > 0,
            "New duration has to be greater than zero"
        );

        emit DurationSecUpdated(durationSec, _newDurationSec);

        durationSec = _newDurationSec;
    }

    // TODO: Add functions that may be required from the Token Dashboard perspective.
}

/// @title KeepTokenGeyserRewardsEscrowBeneficiary
/// @notice Intermediate contract used to transfer tokens from PhasedEscrow to a
/// designated KeepTokenGeyser contract.
contract KeepTokenGeyserRewardsEscrowBeneficiary is StakerRewardsBeneficiary {
    constructor(IERC20 _token, IStakerRewards _stakerRewards)
        public
        StakerRewardsBeneficiary(_token, _stakerRewards)
    {}
}
