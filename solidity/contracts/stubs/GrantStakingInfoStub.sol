pragma solidity 0.5.17;

import "../libraries/staking/GrantStakingInfo.sol";

contract GrantStakingInfoStub {
    using GrantStakingInfo for GrantStakingInfo.Storage;
    GrantStakingInfo.Storage grantStakingInfo;

    function hasGrantDelegated(address operator) public view returns (bool) {
        return grantStakingInfo.hasGrantDelegated(operator);
    }

    function setGrantForOperator(address operator, uint256 grantId) public {
        grantStakingInfo.setGrantForOperator(operator, grantId);
    }

    function getGrantForOperator(address operator) public view returns (uint256) {
        return grantStakingInfo.getGrantForOperator(operator);
    }
}