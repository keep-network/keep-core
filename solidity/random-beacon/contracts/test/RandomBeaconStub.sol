pragma solidity ^0.8.6;

import "../RandomBeacon.sol";
import "../libraries/DKG.sol";
import "../libraries/Callback.sol";
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

    function hasGasDeposit(address operator) external view returns (bool) {
        return gasStation.gasDeposits[operator][0] != 0;
    }

    function publicPunishOperators(
        address[] memory operators,
        uint256 punishmentDuration
    ) external {
        punishOperators(operators, punishmentDuration);
    }
}
