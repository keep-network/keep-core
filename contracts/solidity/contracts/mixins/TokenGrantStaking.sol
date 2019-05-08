pragma solidity ^0.5.4;

import "../utils/UintArrayUtils.sol";
import "../TokenGrant.sol";
import "./StakeDelegatable.sol";

/**
 * @title TokenGrantStaking
 * @dev A token grant staking mixin contract for a specified token grant contract.
 */
contract TokenGrantStaking is StakeDelegatable {

    using UintArrayUtils for uint256[];

    event ReceivedGrantApproval(uint256 _value);
    event Staked(address indexed from, uint256 value);

    TokenGrant public tokenGrant;

    mapping(address => uint256) public stakedGrantsBalances;
    mapping(address => uint256[]) public stakedGrants;

    /**
     * @notice Receives approval of token transfer and stakes the approved ammount.
     * @dev Makes sure provided token contract is the same one linked to this contract.
     * @param _id The owner of the tokens who approved them to transfer.
     * @param _value Approved amount for the transfer and stake.
     * @param _tokenGrant Token grant contract address.
     * @param _extraData Data for stake delegation. This byte array must have the
     * following values concatenated: Magpie address (20 bytes) where the rewards for participation
     * are sent and the operator's ECDSA (65 bytes) signature of the address of the stake owner.
     */
    function receiveGrantApproval(uint256 _id, uint256 _value, address _tokenGrant, bytes memory _extraData) public {
        emit ReceivedGrantApproval(_value);
        require(TokenGrant(_tokenGrant) == tokenGrant, "Token grant contract must be the same one linked to this contract.");
        require(_value <= tokenGrant.transferableBalance(_id), "Token grant doesn't have enough amount available to transfer.");

        address _from = tokenGrant.grantBeneficiary(_id);
        (address magpie, address operator) = _extractDelegationData(_from, _extraData);
        _delegateStake(_from, _value, magpie, operator);

        // Transfer token grant.
        uint256 newGrantId = tokenGrant.transferFrom(_id, address(this), _value);

        // Maintain a record of the stake amount by the sender.
        stakedGrantsBalances[operator] = stakedGrantsBalances[operator].add(_value);
        stakedGrants[operator].push(newGrantId);
        emit Staked(operator, _value);
        if (address(stakingProxy) != address(0)) {
            stakingProxy.emitStakedEvent(operator, _value);
        }
    }

    function _transferUnstakedTokenGrants(address _operator, uint256 _value) internal returns (uint256){

        address owner = operatorToOwner[_operator];
        uint256 remaining = _value;

        for (uint i = 0; i < stakedGrants[_operator].length; i++) {

            uint256 grantId = stakedGrants[_operator][i];
            uint256 transferable = tokenGrant.transferableBalance(grantId);

            if (remaining >= transferable) {

                tokenGrant.transfer(grantId, owner, transferable);

                stakedGrants[_operator].removeValue(grantId);
                stakedGrantsBalances[_operator] = stakedGrantsBalances[_operator].sub(transferable);

                remaining = transferable.sub(remaining);
            } else {
                tokenGrant.transfer(grantId, owner, transferable.sub(remaining));
                stakedGrantsBalances[_operator] = stakedGrantsBalances[_operator].sub(transferable).sub(remaining);
                remaining = 0;
            }
        }

        return remaining;
    }
}
