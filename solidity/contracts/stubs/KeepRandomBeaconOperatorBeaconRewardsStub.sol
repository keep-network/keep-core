pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorBeaconRewardsStub is KeepRandomBeaconOperator {
    using Groups for Groups.Group;

    uint256 constant GROUP_INDEX_FLAG = 1 << 255;

    constructor(
        address _serviceContract,
        address _stakingContract,
        address _registryContract,
        address _gasPriceOracle
    ) KeepRandomBeaconOperator(
        _serviceContract,
        _stakingContract,
        _registryContract,
        _gasPriceOracle
    ) public {
        groupSize = 3;
        groups.groupActiveTime = 5;
        groups.relayEntryTimeout = 10;
    }

    function registerNewGroup(bytes memory groupPublicKey, address[] memory members, uint256 creationTimestamp) public {
        groups.groupIndices[groupPublicKey] = (groups.groups.length ^ GROUP_INDEX_FLAG);
        groups.groups.push(Groups.Group(groupPublicKey, block.number, false, uint248(creationTimestamp)));
    
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }

    function terminateGroup(uint256 groupIndex) public {
        groups.terminateGroup(groupIndex);
    }

    function expireOldGroups() public {
        groups.expireOldGroups();
    }
}
