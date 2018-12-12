pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./StakingManager.sol";


/**
 * @title Manageable Staking
 * @dev A generic contract that provides authorized staking balance modifications.
 */
contract ManageableStaking is Verifier {
    using SafeMath for uint256;

    mapping(address => uint256) public balances;
    mapping(address => mapping (address => bool)) public authorization;

    /**
     * @notice Authorize an address to be able to transfer your staked tokens.
     * @param manager Manager of your staked tokens.
     * @param signature Manager address signed by you (ECDSA sig, i.e. using web3.eth.sign()).
     */
    function authorizeManager(address manager, bytes signature) public {
        require(manager != address(0), "Manager address can't be zero.");
        require(manager != msg.sender, "You can not authorize your own address.");
        StakingManager m = StakingManager(manager);
        authorization[msg.sender][manager] = true;
        m.receiveAuthorization(msg.sender, signature);
    }

    /**
     * @notice Deauthorize an address from being able to transfer your staked tokens.
     * @param manager Manager of your staked tokens.
     */
    function deauthorizeManager(address manager) public {
        delete authorization[msg.sender][manager];
    }

    /**
     * @notice Authorized transfer of staked tokens.
     * @dev Transfer that can be executed only by a certain address
     * previously authorized by a staker.
     * @param staker Staker who authorized transfer.
     * @param amount Amount of the staked tokens.
     */
    function authorizedTransfer(address staker, uint256 amount) public {
        require(isManagerAuthorizedFor(staker, msg.sender), "Manager address must be authorized by the staker.");
        require(amount <= balances[staker], "Staker must have enough tokens to transfer.");
        balances[staker] = balances[staker].sub(amount);
        balances[msg.sender] = balances[msg.sender].add(amount);
    }

    function isManagerAuthorizedFor(address staker, address manager) public view returns (bool) {
        return authorization[staker][manager];
    }

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _staker The address to query the balance of.
     * @return An uint256 representing the amount owned by the passed address.
     */
    function stakeBalanceOf(address _staker) public view returns (uint256 balance) {
        return balances[_staker];
    }
}
