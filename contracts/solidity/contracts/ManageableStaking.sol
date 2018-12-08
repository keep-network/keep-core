pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";


/**
 * @title Manageable Staking
 * @dev A generic contract that provides authorized staking balance modifications.
 */
contract ManageableStaking {
    using SafeMath for uint256;

    mapping(address => uint256) public balances;
    mapping(address => mapping (address => bool)) public authorization;

    /**
     * @notice Authorize an address to transfer your staked tokens.
     * @param receiver Receiver of your staked tokens.
     */
    function authorize(address receiver) public {
        require(receiver != address(0), "Receiver address can't be zero.");
        require(receiver != msg.sender, "You can not authorize your own address.");
        authorization[msg.sender][receiver] = true;
    }

    /**
     * @notice Deauthorize an address from transferring your staked tokens.
     * @param receiver Receiver of your staked tokens.
     */
    function deauthorize(address receiver) public {
        delete authorization[msg.sender][receiver];
    }

    /**
     * @notice Authorized transfer of staked tokens.
     * @dev Transfer that can be executed only by a certain address
     * previously authorized by a staker.
     * @param staker Staker who authorized transfer.
     * @param amount Amount of the staked tokens.
     */
    function authorizedTransfer(address staker, uint256 amount) public {
        require(authorization[staker][msg.sender], "Receiver address must be authorized by the staker.");
        require(amount <= balances[staker], "Staker must have enough tokens to transfer.");
        balances[staker] = balances[staker].sub(amount);
        balances[msg.sender] = balances[msg.sender].add(amount);
    }
}
