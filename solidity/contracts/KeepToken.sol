pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/StandardToken.sol';

/**
 * @title KEEP Token
 * @dev Standard ERC20 token
 */
contract KeepToken is StandardToken {
  string public constant NAME = "KEEP Token";
  string public constant SYMBOL = "KEEP";
  uint256 public constant DECIMALS = 18;
  uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(DECIMALS));

  /**
   * @dev Gives msg.sender all of existing tokens.
   */
  function KeepToken() {
    totalSupply = INITIAL_SUPPLY;
    balances[msg.sender] = INITIAL_SUPPLY;
  }

  // /**
  //  * @notice Stake vesting.
  //  * Stakable vested amount is the amount of vested tokens minus what user already released from the vesting
  //  * @param _id Vesting ID
  //  */
  // function stakeVesting(uint256 _id) public {

  //   // Vesting must be unlocked and not revoked
  //   require(!vestings[_id].locked);
  //   require(!vestings[_id].revoked);
  
  //   // Make sure decision to unstake is up to the beneficiary of the vesting
  //   require(vestings[_id].beneficiary == msg.sender);
  //   // Calculate available amount. Amount of vested tokens minus what user already released
  //   uint256 available = vestings[_id].amount.sub(vestings[_id].released);
  //   require(available > 0);

  //   // Lock vesting from releasing it's balance
  //   vestings[_id].locked = true;
  
  //   // Transfer tokens to beneficiary's vesting stake balance
  //   vestingStakeBalances[vestings[_id].beneficiary] = vestingStakeBalances[vestings[_id].beneficiary].add(available);
  // }

  // /**
  //  * @notice Initiate unstake of the vesting.
  //  * @param _id Vesting ID
  //  */
  // function initiateUnstakeVesting(uint256 _id) public {

  //   // Vesting must be locked and not revoked
  //   require(vestings[_id].locked);
  //   require(!vestings[_id].revoked);

  //   // Make sure decision to unstake is up to the beneficiary of the vesting
  //   require(msg.sender == vestings[_id].beneficiary);
    
  //   // Vesting withdrawal start shouldn't be set
  //   require(vestingStakeWithdrawalStart[_id] == 0);

  //   // Set vesting stake withdrawal start
  //   vestingStakeWithdrawalStart[_id] = now;
  // }

  // /**
  // * @dev Finish unstake of the vesting
  // * @param _id Vesting ID
  // */
  // function finishUnstakeVesting(uint256 _id) public {

  //   // Vesting withdrawal start must be set
  //   require(vestingStakeWithdrawalStart[_id] > 0);

  //   // Vesting must be locked and not revoked
  //   require(vestings[_id].locked);
  //   require(!vestings[_id].revoked);

  //   // Vesting withdrawal delay should be over
  //   require(now >= vestingStakeWithdrawalStart[_id].add(withdrawalDelay));

  //   // Calculate vesting amount that was staked
  //   uint256 available = vestings[_id].amount.sub(vestings[_id].released);
  //   require(available > 0);

  //   // Remove tokens from vesting stake balance
  //   vestingStakeBalances[vestings[_id].beneficiary] = vestingStakeBalances[vestings[_id].beneficiary].sub(available);
    
  //   // Unlock vesting
  //   vestings[_id].locked = false;

  //   // Unset vesting withdrawal start
  //   vestingStakeWithdrawalStart[_id] = 0;

  // }
}
