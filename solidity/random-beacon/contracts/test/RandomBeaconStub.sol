pragma solidity ^0.8.6;

import "../RandomBeacon.sol";
import "../libraries/DKG.sol";
import "../libraries/Callback.sol";
import "../libraries/Groups.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract RandomBeaconStub is RandomBeacon {
    constructor(
        ISortitionPool _sortitionPool,
        IERC20 _tToken,
        IStaking _staking
    ) RandomBeacon(_sortitionPool, _tToken, _staking) {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }

    function getCallbackData() external view returns (Callback.Data memory) {
        return callback;
    }

    function roughlyAddGroup(
        bytes calldata groupPubKey,
        uint32[] calldata members
    ) external {
        bytes32 groupPubKeyHash = keccak256(groupPubKey);

        Groups.Group memory group;
        group.groupPubKey = groupPubKey;
        group.members = members;
        /* solhint-disable-next-line not-rely-on-time */
        group.activationTimestamp = block.timestamp;

        groups.groupsData[groupPubKeyHash] = group;
        groups.groupsRegistry.push(groupPubKeyHash);
        groups.activeGroupsCount++;
    }

    function publicPunishOperators(
        address[] memory operators,
        uint256 punishmentDuration
    ) external {
        punishOperators(operators, punishmentDuration);
    }

    function hasGasDeposit(address operator) external view returns (bool) {
        return gasStation.gasDeposits[operator][0] != 0;
    }
}
