pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./ManageableStaking.sol";
import "./Verifier.sol";


/**
 * @title Staking Manager
 * @dev A generic contract that provides authorized staking balance modifications.
 */
contract StakingManager is Verifier {

    mapping (address => address) private _manageableContractFor;

    /**
     * @notice Receives authorization to modify staked balance of a staker.
     * @param staker The staker who authorized balance modifications for this contract.
     * @param signature Address of this contract signed by the staker.
     */
    function receiveAuthorization(address staker, bytes signature) public {

        // Expecting manageableStaking contract as a caller for this method.
        ManageableStaking stakingContract = ManageableStaking(msg.sender);

        require(stakingContract.stakeBalanceOf(staker) > 0, "Staker must have staked token.");
        require(_manageableContractFor[staker] == 0, "Manageable contract for the staker is already set.");
        require(isSigned(keccak256(address(this)), signature, staker), "This contract address must be signed by the staker.");

        _manageableContractFor[staker] = msg.sender;
    }

    function manageableContractFor(address staker) public view returns (address) {
        return _manageableContractFor[staker];
    }

    function _transfer(address staker, uint256 amount) internal {
        ManageableStaking stakingContract = ManageableStaking(_manageableContractFor[staker]);
        stakingContract.authorizedTransfer(staker, amount);
    }

}
