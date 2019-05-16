pragma solidity ^0.5.4;

import "./KeepGroupImplV1.sol";

contract KeepGroupImplV1Stub is KeepGroupImplV1 {
    
    /**
     * @dev Adds a new group based on groupPublicKey.
     * @param groupPublicKey is the identifier of the newly created group.
     */
    function registerNewGroup(bytes memory groupPublicKey) public {
        _groups.push(Group(groupPublicKey, block.number));
        address[] memory members = orderedParticipants();
        if (members.length > 0) {
            for (uint i = 0; i < _groupSize; i++) {
                _groupMembers[groupPublicKey].push(members[i]);
            }
        }
    }

    /**
     * @dev Gets the group registration block height.
     * @param groupIndex is the index of the queried group.
     */
    function getGroupRegistrationBlockHeight(uint256 groupIndex) public view returns(uint256) {
        return _groups[groupIndex].registrationBlockHeight;
    }

    /**
     * @dev Gets the public key of the group registered under the given index.
     * @param groupIndex is the index of the queried group.
     */
    function getGroupPublicKey(uint256 groupIndex) public view returns(bytes memory) {
        return _groups[groupIndex].groupPubKey;
    }

    /**
     * @dev Gets the value of expired offset.
     */
    function getExpiredOffset() public view returns(uint256) {
        return _expiredOffset;
    }
}