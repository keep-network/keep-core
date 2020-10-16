pragma solidity 0.5.17;

import "../KeepRandomBeaconOperator.sol";

contract KeepRandomBeaconOperatorBeaconRewardsStub is KeepRandomBeaconOperator {
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
        groups.addGroup(groupPublicKey);
        groups.groups[groups.groups.length - 1].registrationTime = uint248(creationTimestamp);
    
        groups.setGroupMembers(groupPublicKey, members, hex"");
    }

    function terminateGroup(uint256 groupIndex) public {
        groups.terminateGroup(groupIndex);
    }

    function expireOldGroups() public {
        groups.expireOldGroups();
    }
}
