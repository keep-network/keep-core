pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/BasicToken.sol';
import 'zeppelin-solidity/contracts/token/StandardToken.sol';
import 'zeppelin-solidity/contracts/token/SafeERC20.sol';
import 'zeppelin-solidity/contracts/math/SafeMath.sol';

/**
 * @title TokenVesting
 * @dev Token Vesting contract that vest and stake/unstake vested ERC20 tokens.  
 * Vesting balance released to the beneficiary gradually like a
 * typical vesting scheme, with a cliff and vesting period. Optionally revocable by the
 * owner.
 */
contract TokenVesting is BasicToken {
  using SafeMath for uint256;
  using SafeERC20 for StandardToken;

  event NewVesting(uint256 id);
  event VestingReleased(uint256 amount);
  event InitiateUnstakeVesting(uint256 id);
  event Revoked(uint256 id);

  // Token contract
  StandardToken public token;

  struct Vesting {
    address owner; // creator of vesting 
    address beneficiary; // address to which vested tokens are transferred
    bool locked; // whether the vesting is locked (i.e. for staking)
    bool revoked; // whether the vesting is revoked
    bool revocable; // whether creator of vesting can revoke it
    uint256 amount; // amount to be vested
    uint256 duration; // duration in seconds of the period in which the tokens will vest
    uint256 start; // timestamp at which vesting will start
    uint256 cliff; // duration in seconds of the cliff after which tokens will begin to vest
    uint256 released; // amount that was released to the beneficiary
  }

  uint256 public stakeWithdrawalDelay;
  uint256 public numVestings;

  // Vesting balances
  // Sum of all vested tokens to a beneficiary
  mapping(address => uint256) public vestingBalances;

  // Vestings
  mapping(uint256 => Vesting) public vestings;
  
  // Vestings stake balances
  mapping(address => uint256) public vestingStakeBalances;

  // Vesting stake withdrawals
  mapping(uint256 => uint256) public vestingStakeWithdrawalStart;

  function TokenVesting(StandardToken _token, uint256 _delay) {
    stakeWithdrawalDelay = _delay;
    token = _token;
  }

  /**
  * @dev Gets the vesting balance of the specified address.
  * @param _owner The address to query the vesting balance of.
  * @return An uint256 representing the vesting balance owned by the passed address.
  */
  function vestingBalanceOf(address _owner) public constant returns (uint256 balance) {
    return vestingBalances[_owner];
  }

  /**
  * @dev Gets the vesting stake balance of the specified address.
  * @param _owner The address to query the vesting balance of.
  * @return An uint256 representing the vesting stake balance owned by the passed address.
  */
  function vestingStakeBalanceOf(address _owner) public constant returns (uint256 balance) {
    return vestingStakeBalances[_owner];
  }

  /**
   * @dev Creates a vesting schedule with its balance released to the
   * beneficiary gradually in a linear fashion until start + duration. By then all
   * of the balance will have vested.
   * @param _amount to be vested
   * @param _beneficiary address to which vested tokens are transferred
   * @param _cliff duration in seconds of the cliff after which tokens will begin to vest
   * @param _duration duration in seconds of the period in which the tokens will vest
   * @param _start timestamp at which vesting will start
   * @param _revocable whether the vesting is revocable or not
   */
  function vest(uint256 _amount, address _beneficiary, uint256 _duration, uint256 _start, uint256 _cliff, bool _revocable) public returns (uint256) {
    require(_beneficiary != address(0));
    require(_cliff <= _duration);
    require(_amount <= token.balanceOf(msg.sender));
    
    // Create new vesting schedule
    uint256 id = numVestings++;
    vestings[id] = Vesting(msg.sender, _beneficiary, false, false, _revocable, _amount, _duration, _start, _start.add(_cliff), 0);

    // Transfer amount from sender to this vesting contract
    // Sender should approve the amount first by calling approve() on the token
    token.transferFrom(msg.sender, this, _amount);

    // Keep record of the vested amount 
    vestingBalances[_beneficiary] = vestingBalances[_beneficiary].add(_amount);
    NewVesting(id);
    return id;
  }

  /**
   * @notice Transfers vested tokens to beneficiary.
   * @param _id Vesting ID
   */
  function releaseVesting(uint256 _id) public {
    require(!vestings[_id].locked);
    uint256 unreleased = releasableVestedAmount(_id);
    require(unreleased > 0);

    // Update released amount
    vestings[_id].released = vestings[_id].released.add(unreleased);

    // Update beneficiary vesting balance
    vestingBalances[vestings[_id].beneficiary] = vestingBalances[vestings[_id].beneficiary].sub(unreleased);

    // Transfer tokens from this vesting contract balance to the beneficiary token balance
    token.safeTransfer(vestings[_id].beneficiary, unreleased);

    VestingReleased(unreleased);
  }
  
  /**
   * @dev Calculates the amount that has already vested, 
   * including any tokens that have already been withdrawn by the beneficiary 
   * as well as any tokens that are available to withdraw but have not yet been withdrawn
   * @param _id Vesting ID
   */
  function vestedAmount(uint256 _id) public constant returns (uint256) {
    uint256 balance = vestings[_id].amount;

    if (now < vestings[_id].cliff) {
      return 0; // Cliff period is not over
    } else if (now >= vestings[_id].start.add(vestings[_id].duration) || vestings[_id].revoked) {
      return balance; // Vesting period is finished.
    } else {
      return balance.mul(now.sub(vestings[_id].start)).div(vestings[_id].duration);
    }
  }

  /**
  * @dev Calculates the amount that has already vested but hasn't been released yet.
  * @param _id Vesting ID
  */
  function releasableVestedAmount(uint256 _id) public constant returns (uint256) {
    uint256 released = vestings[_id].released;
    return vestedAmount(_id).sub(released);
  }


  /**
   * @notice Stake vesting.
   * Stakable vested amount is the amount of vested tokens minus what user already released from the vesting
   * @param _id Vesting ID
   */
  function stakeVesting(uint256 _id) public {

    // Vesting must be unlocked and not revoked
    require(!vestings[_id].locked);
    require(!vestings[_id].revoked);
  
    // Make sure decision to stake is up to the beneficiary of the vesting
    require(vestings[_id].beneficiary == msg.sender);
    // Calculate available amount. Amount of vested tokens minus what user already released
    uint256 available = vestings[_id].amount.sub(vestings[_id].released);
    require(available > 0);

    // Lock vesting from releasing its balance
    vestings[_id].locked = true;
  
    // Transfer tokens to beneficiary's vesting stake balance
    vestingStakeBalances[vestings[_id].beneficiary] = vestingStakeBalances[vestings[_id].beneficiary].add(available);
  }

  /**
   * @notice Initiate unstake of the vesting.
   * @param _id Vesting ID
   */
  function initiateUnstakeVesting(uint256 _id) public {

    // Vesting must be locked and not revoked
    require(vestings[_id].locked);
    require(!vestings[_id].revoked);

    // Make sure decision to unstake is up to the beneficiary of the vesting
    require(msg.sender == vestings[_id].beneficiary);
    
    // Vesting withdrawal start shouldn't be set
    require(vestingStakeWithdrawalStart[_id] == 0);

    // Set vesting stake withdrawal start
    vestingStakeWithdrawalStart[_id] = now;

    InitiateUnstakeVesting(_id);
  }

  /**
  * @dev Finish unstake of the vesting
  * @param _id Vesting ID
  */
  function finishUnstakeVesting(uint256 _id) public {

    // Vesting withdrawal start must be set
    require(vestingStakeWithdrawalStart[_id] > 0);

    // Vesting must be locked and not revoked
    require(vestings[_id].locked);
    require(!vestings[_id].revoked);

    // Vesting withdrawal delay should be over
    require(now >= vestingStakeWithdrawalStart[_id].add(stakeWithdrawalDelay));

    // Calculate vesting amount that was staked
    uint256 available = vestings[_id].amount.sub(vestings[_id].released);
    require(available > 0);

    // Remove tokens from vesting stake balance
    vestingStakeBalances[vestings[_id].beneficiary] = vestingStakeBalances[vestings[_id].beneficiary].sub(available);
    
    // Unlock vesting
    vestings[_id].locked = false;

    // Unset vesting withdrawal start
    vestingStakeWithdrawalStart[_id] = 0;

  }


  /**
   * @notice Allows the owner to revoke the vesting. 
   * Tokens already vested (releasable amount) remain so beneficiary can still release them
   * the rest are returned to the vesting owner.
   * @param _id Vesting ID
   */
  function revoke(uint256 _id) public {

    // only vesting owner can revoke
    require(vestings[_id].owner == msg.sender);

    // vesting must be revocable in the first place
    require(vestings[_id].revocable);

    // vesting must not be already revoked
    require(!vestings[_id].revoked);

    // vesting must not be locked for staking
    require(!vestings[_id].locked);

    uint256 unreleased = releasableVestedAmount(_id);
    uint256 refund = vestings[_id].amount.sub(unreleased);
    vestings[_id].revoked = true;

    // Update beneficiary vesting balance
    vestingBalances[vestings[_id].beneficiary] = vestingBalances[vestings[_id].beneficiary].sub(refund);

    // Transfer tokens from this vesting contract balance to the owner of the vesting
    token.safeTransfer(vestings[_id].owner, refund);
    Revoked(_id);
  }
}
