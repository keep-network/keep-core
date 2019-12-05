pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";

/**
 @dev Interface of recipient contract for approveAndCall pattern.
*/
interface tokenRecipient {function receiveApproval(address _from, uint256 _value, address _token, address _magpie, address _operator) external;}


/**
 * @title KEEP Token
 * @dev Standard ERC20 token
 */
contract KeepToken is ERC20 {
    string public constant NAME = "KEEP Token";
    string public constant SYMBOL = "KEEP";
    uint8 public constant DECIMALS = 18; // The number of digits after the decimal place when displaying token values on-screen.
    uint256 public constant INITIAL_SUPPLY = 10**27; // 1 billion tokens, 18 decimal places.

    /**
     * @dev Gives msg.sender all of existing tokens.
     */
    constructor() public {
        _mint(msg.sender, INITIAL_SUPPLY);
    }

    /**
     * @notice Set allowance for other address and notify.
     * Allows `_spender` to spend no more than `_value` tokens
     * on your behalf and then ping the contract about it.
     * @param _spender The address authorized to spend.
     * @param _value The max amount they can spend.
     * @param _magpie Magpie address where the rewards for participation are sent.
     * @param _operator The address of a party authorized to operate a stake on behalf of a given owner.
     */
    function approveAndCall(address _spender, uint256 _value, address _magpie, address _operator) public returns (bool success) {
        tokenRecipient spender = tokenRecipient(_spender);
        if (approve(_spender, _value)) {
            spender.receiveApproval(msg.sender, _value, address(this), _magpie, _operator);
            return true;
        }
    }

}
