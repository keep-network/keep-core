pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/StandardToken.sol';

/**
 * @title VestableToken
 * @dev Adds vesting functionality to the token. 
 * Vesting balance released to the beneficiary gradually like a
 * typical vesting scheme, with a cliff and vesting period. Optionally revocable by the
 * owner.
 */
contract VestableToken is StandardToken {
  using SafeMath for uint256;

  event NewVesting(uint256 id);
  event VestingReleased(uint256 amount);

  struct Vesting {
    address owner;
    address beneficiary;
    bool locked;
    bool revoked;
    bool revocable;
    uint256 amount;
    uint256 duration;
    uint256 start;
    uint256 cliff;
    uint256 released;
  }

  uint256 public numVestings;

  // Vesting balances
  // Sum of all vested tokens to a beneficiary
  mapping(address => uint256) public vestingBalances;

  // Vestings
  mapping(uint256 => Vesting) public vestings;
  
  
  function VestableToken() {
  }

  /**
  * @dev Gets the vesting balance of the specified address.
  * @param _owner The address to query the vesting balance of.
  * @return An uint256 representing the amount owned by the passed address.
  */
  function vestingBalanceOf(address _owner) public constant returns (uint256 balance) {
    return vestingBalances[_owner];
  }

  /**
   * @dev Creates a vesting schedule with its balance released to the
   * beneficiary gradually in a linear fashion until start + duration. By then all
   * of the balance will have vested.
   * @param _amount to be vested
   * @param _beneficiary address to whom vested tokens are transferred
   * @param _cliff duration in seconds of the cliff in which tokens will begin to vest
   * @param _duration duration in seconds of the period in which the tokens will vest
   * @param _revocable whether the vesting is revocable or not
   */
  function vest(uint256 _amount, address _beneficiary, uint256 _duration, uint256 _start, uint256 _cliff, bool _revocable) public returns (uint256) {
    require(_beneficiary != address(0));
    require(_cliff <= _duration);
    require(_amount <= balances[msg.sender]);
    
    // Create new vesting schedule
    uint256 id = numVestings++;
    vestings[id] = Vesting(msg.sender, _beneficiary, false, false, _revocable, _amount, _duration, _start, _start.add(_cliff), 0);

    // Transfer amount from sender to the beneficiary vesting balance
    balances[msg.sender] = balances[msg.sender].sub(_amount);
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
    require(!vestings[_id].revoked);
    uint256 unreleased = releasableVestedAmount(_id);
    require(unreleased > 0);

    // Update released amount
    vestings[_id].released = vestings[_id].released.add(unreleased);

    // Transfer tokens to beneficiary balance
    vestingBalances[vestings[_id].beneficiary] = vestingBalances[vestings[_id].beneficiary].sub(unreleased);
    balances[vestings[_id].beneficiary] = balances[vestings[_id].beneficiary].add(unreleased);

    VestingReleased(unreleased);
  }
  
  /**
   * @dev Calculates the amount that has already vested, 
   * inlcuding amount that could be already withdrawn by the beneficiary
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
}
