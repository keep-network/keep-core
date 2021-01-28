/**
▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓▌        ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
  ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓

                           Trust math, not hardware.
*/

pragma solidity ^0.5.17;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/math/Math.sol";

import "./KeepToken.sol";
import "./RewardPayouts.sol";
import "./libraries/BinarySearch.sol";

/// @title KEEP Signer Subsidy Rewards
/// @notice A contract for distributing KEEP token rewards to keeps.
/// When a reward contract is created, the creator defines a reward schedule
/// consisting of one or more reward intervals and their interval weights,
/// the length of reward intervals, and the quota of how many keeps must be
/// created in an interval for the full reward for that interval to be paid out.
///
/// The amount of KEEP to be distributed is determined by funding the contract,
/// and additional KEEP can be added at any time.
/// The reward contract is funded with `approveAndCall` with no extra data,
/// but it also collects any KEEP mistakenly sent to it in any other way.
///
/// An interval is defined by the timestamps [startOf, endOf);
/// a keep created at the time `startOf(i)` belongs to interval `i`
/// and one created at `endOf(i)` belongs to `i+1`.
///
/// When an interval is over, it will be allocated a percentage of the remaining
/// unallocated rewards based on its weight, and adjusted by the number of keeps
/// created in the interval if the quota is not met.
///
/// The adjustment for not meeting the keep quota is a percentage that equals
/// the percentage of the quota that was met; if the number of keeps created is
/// 80% of the quota then 80% of the base reward will be allocated for the
/// interval.
///
/// Any unallocated rewards will stay in the unallocated rewards pool,
/// to be allocated for future intervals. Intervals past the initially defined
/// schedule have a weight of 100%, meaning that all remaining unallocated
/// rewards will be allocated to the interval.
///
/// Keeps of the appropriate type can receive rewards once the interval they
/// were created in is over, and the keep has closed happily.
/// There is no time limit to receiving rewards, nor is there need to wait for
/// all keeps from the interval to close.
/// Calling `receiveReward` automatically allocates the rewards for the interval
/// the specified keep was created in and all previous intervals.
///
/// If a keep is terminated, that fact can be reported to the reward contract.
/// Reporting a terminated keep returns its allocated reward to the pool of
/// unallocated rewards.
///
/// @dev A concrete implementation of the abstract rewards contract must specify
/// functions for accessing information about keeps and paying out rewards.
/// For the purpose of rewards, Random Beacon signing groups count as "keeps"
/// and the beacon operator contract acts as the "factory".
contract Rewards is RewardPayouts {
    constructor (
        address _token,
        uint256 _firstIntervalStart,
        uint256[] memory _intervalWeights,
        uint256 _termLength,
        uint256 _minimumKeepsPerInterval
    ) public RewardPayouts(
        _token,
        _firstIntervalStart,
        _intervalWeights,
        _termLength,
        _minimumKeepsPerInterval
    ) {
    }

    function _binarySearchGetValue(uint256 index) internal view returns (uint256) {
        return _getCreationTime(_getKeepAtIndex(index));
    }

    function _findEndpoint(uint256 intervalEndpoint) internal view returns (uint256) {
        require(
            intervalEndpoint <= block.timestamp,
            "interval hasn't ended yet"
        );
        uint256 keepCount = _getKeepCount();

        return BinarySearch.find(
            _binarySearchGetValue,
            keepCount,
            intervalEndpoint
        );
    }

    /// @notice Get the identifier of the keep at the given index,
    /// when all keeps created by the factory are ordered by creation time.
    /// @param index The index of the queried keep.
    /// @return The `bytes32` identifier of the keep at the given index.
    /// @dev Implementation is not required to check if a keep with the given
    /// index exists.
    function _getKeepAtIndex(uint256 index) internal view returns (bytes32);
}
