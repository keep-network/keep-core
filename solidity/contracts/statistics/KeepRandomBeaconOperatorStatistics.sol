pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorStatistics {
    using SafeMath for uint256;

    KeepRandomBeaconOperator public operatorContract;

    constructor(
        address _operatorContract
    ) public {
        operatorContract = KeepRandomBeaconOperator(_operatorContract);
    }

    /**
     * @dev Gets all indices in the provided group for a member.
     */
    function getGroupMemberIndices(
        bytes memory groupPubKey,
        address member
    ) public view returns (uint256[] memory indices) {
        address[] memory members = operatorContract.getGroupMembers(groupPubKey);

        uint256 counter;
        for (uint i = 0; i < members.length; i++) {
            if (members[i] == member) {
                counter++;
            }
        }

        indices = new uint256[](counter);
        counter = 0;
        for (uint i = 0; i < members.length; i++) {
            if (members[i] == member) {
                indices[counter] = i;
                counter++;
            }
        }
    }

    function awaitingRewards(
        address operator,
        uint256 groupIndex
    ) public view returns (uint256 rewards) {
        if (operatorContract.hasWithdrawnRewards(operator, groupIndex)) {
            return 0;
        }
        bytes memory groupPubKey = operatorContract.getGroupPublicKey(groupIndex);
        uint256 memberRewards = operatorContract.getGroupMemberRewards(groupPubKey);

        uint256 memberCount = getGroupMemberIndices(groupPubKey, operator).length;

        return memberRewards.mul(memberCount);
    }

    function withdrawableRewards(
        address operator,
        uint256 groupIndex
    ) public view returns (uint256 rewards) {
        bytes memory groupPubKey = operatorContract.getGroupPublicKey(groupIndex);
        if (operatorContract.isStaleGroup(groupPubKey)) {
            return awaitingRewards(operator, groupIndex);
        } else {
            return 0;
        }
    }
}
