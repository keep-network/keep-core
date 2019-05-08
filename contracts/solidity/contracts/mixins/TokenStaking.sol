pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "./StakeDelegatable.sol";

/**
 * @title TokenStaking
 * @dev A token staking mixin contract for a specified standard ERC20 token.
 */
contract TokenStaking is StakeDelegatable {

    using SafeERC20 for ERC20;

    event ReceivedApproval(uint256 _value);
    event Staked(address indexed from, uint256 value);

    ERC20 public token;

    mapping(address => uint256) public stakedBalances;

    /**
     * @notice Receives approval of token transfer and stakes the approved ammount.
     * @dev Makes sure provided token contract is the same one linked to this contract.
     * @param _from The owner of the tokens who approved them to transfer.
     * @param _value Approved amount for the transfer and stake.
     * @param _token Token contract address.
     * @param _extraData Data for stake delegation. This byte array must have the
     * following values concatenated: Magpie address (20 bytes) where the rewards for participation
     * are sent and the operator's ECDSA (65 bytes) signature of the address of the stake owner.
     */
    function receiveApproval(address _from, uint256 _value, address _token, bytes memory _extraData) public {
        emit ReceivedApproval(_value);
        require(ERC20(_token) == token, "Token contract must be the same one linked to this contract.");
        require(_value <= token.balanceOf(_from), "Sender must have enough tokens.");

        (address magpie, address operator) = _extractDelegationData(_from, _extraData);
        _delegateStake(_from, _value, magpie, operator);

        // Transfer tokens to this contract.
        token.transferFrom(_from, address(this), _value);

        // Maintain a record of the stake amount by the sender.
        stakedBalances[operator] = stakedBalances[operator].add(_value);
        emit Staked(operator, _value);
    }

    /**
     * @notice Transfer unstaked tokens to the owner.
     * @param _operator Address of the operator.
     * @param _value Amount of tokens to transfer.
     */
    function _transferUnstakedTokens(address _operator, uint256 _value) internal {
        address owner = operatorToOwner[_operator];
        token.safeTransfer(owner, _value);
        stakedBalances[_operator] = stakedBalances[_operator].sub(_value);
    }
}
