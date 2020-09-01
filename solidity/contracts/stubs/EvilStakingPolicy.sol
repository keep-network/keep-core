pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../GrantStakingPolicy.sol";

/// @title EvilStakingPolicy
/// @dev A staking policy which allows the grantee to stake
/// a million times more than the grant amount.
contract EvilStakingPolicy is GrantStakingPolicy {
    using SafeMath for uint256;

    function getStakeableAmount (
        uint256 _now,
        uint256 grantedAmount,
        uint256 duration,
        uint256 start,
        uint256 cliff,
        uint256 withdrawn
    ) public view returns (uint256) {
        return grantedAmount.mul(1000000);
    }
}
