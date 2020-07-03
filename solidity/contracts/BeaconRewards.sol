pragma solidity ^0.5.17;

import "./Rewards.sol";
import "./KeepRandomBeaconOperator.sol";
import "./TokenStaking.sol";

contract BeaconRewards is Rewards {
    KeepRandomBeaconOperator operatorContract;
    TokenStaking tokenStaking;

    constructor (
        uint256 _termLength,
        address _token,
        uint256 _minimumKeepsPerInterval,
        uint256 _firstIntervalStart,
        uint256[] memory _intervalWeights,
        address _operatorContract,
        address _stakingContract
    ) public Rewards(
        _termLength,
        _token,
        _minimumKeepsPerInterval,
        _firstIntervalStart,
        _intervalWeights
    ) {
        operatorContract = KeepRandomBeaconOperator(_operatorContract);
        tokenStaking = TokenStaking(_stakingContract);
    }

    function _getKeepCount() internal view returns (uint256) {
        return operatorContract.getNumberOfCreatedGroups();
    }

    function _getKeepAtIndex(uint256 i) internal view returns (bytes32) {
        return bytes32(i);
    }

    function _getCreationTime(bytes32 groupIndexBytes) internal view returns (uint256) {
        return operatorContract.getGroupRegistrationTime(uint256(groupIndexBytes));
    }

    function _isClosed(bytes32 groupIndexBytes) internal view returns (bool) {
        if (_isTerminated(groupIndexBytes)) { return false; }
        bytes memory groupPubkey = operatorContract.getGroupPublicKey(
            uint256(groupIndexBytes)
        );
        return operatorContract.isStaleGroup(groupPubkey);
    }

    function _isTerminated(bytes32 groupIndexBytes) internal view returns (bool) {
        return operatorContract.isGroupTerminated(uint256(groupIndexBytes));
    }

    function _recognizedByFactory(bytes32 groupIndexBytes) internal view returns (bool) {
        return _getKeepCount() > uint256(groupIndexBytes);
    }

    function _distributeReward(bytes32 groupIndexBytes, uint256 _value) internal {
        bytes memory groupPubkey = operatorContract.getGroupPublicKey(
            uint256(groupIndexBytes)
        );
        address[] memory members = operatorContract.getGroupMembers(groupPubkey);

        uint256 memberCount = members.length;
        uint256 dividend = _value.div(memberCount);

        require(dividend > 0, "Dividend value must be non-zero");

        for (uint16 i = 0; i < memberCount - 1; i++) {
            token.safeTransfer(
                tokenStaking.beneficiaryOf(members[i]),
                dividend
            );
        }

        // Transfer of dividend for the last member. Remainder might be equal to
        // zero in case of even distribution or some small number.
        uint256 remainder = _value.mod(memberCount);
        token.safeTransfer(
            tokenStaking.beneficiaryOf(members[memberCount - 1]),
            dividend.add(remainder)
        );
    }
}
