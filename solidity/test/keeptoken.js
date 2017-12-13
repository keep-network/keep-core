import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenVesting = artifacts.require('./TokenVesting.sol');

contract('KeepToken', function(accounts) {

  let token, vestingContract, stakingContract,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingContract  = await TokenStaking.new(token.address, duration.days(30));
    vestingContract  = await TokenVesting.new(token.address, duration.days(30));
  });

  it("should send tokens correctly", async function() {
    let amount = 1000000000;

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Send tokens
    await token.transfer(account_two, amount, {from: account_one});

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - amount, "Amount wasn't correctly taken from the sender");
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + amount, "Amount wasn't correctly sent to the receiver");

  });

  it("should stake and unstake tokens correctly", async function() {
    
    let stakingAmount = 10000000;

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);

    // Stake tokens
    await token.approve(stakingContract.address, stakingAmount, {from: account_one});
    await stakingContract.stake(stakingAmount);
    
    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_one_stake_balance = await stakingContract.stakeBalanceOf.call(account_one);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - stakingAmount, "Staking amount should be transfered from sender balance");
    assert.equal(account_one_stake_balance.toNumber(), stakingAmount, "Staking amount should be added to the sender staking balance");
    
    // Initiate unstake tokens
    let stakeWithdrawalId = await stakingContract.initiateUnstake(stakingAmount, {from: account_one}).then((result)=>{
      // Look for initiateUnstake event in transaction receipt and get stake withdrawal id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "InitiateUnstake") {
          return log.args.id.toNumber();
        }
      }
    })

    // should not be able to finish unstake
    await exceptThrow(stakingContract.finishUnstake(stakeWithdrawalId));
    
    // jump in time, full withdrawal delay
    await increaseTimeTo(latestTime()+duration.days(30));

    // should be able to finish unstake
    await stakingContract.finishUnstake(stakeWithdrawalId);

    // check balances
    account_one_ending_balance = await token.balanceOf.call(account_one);
    account_one_stake_balance = await stakingContract.stakeBalanceOf.call(account_one);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber(), "Staking amount should be transfered to sender balance");
    assert.equal(account_one_stake_balance.toNumber(), 0, "Staking amount should be removed from sender staking balance");

  });


  it("should vest tokens correctly", async function() {

    let vestingAmount = 1000000000;
    let vestingDuration = duration.days(30);
    let vestingStart = latestTime();
    let vestingCliff = duration.days(10);
    let vestingRevocable = true;
    
    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Vest tokens
    await token.approve(vestingContract.address, vestingAmount, {from: account_one});
    let vestingId = await vestingContract.vest(vestingAmount, account_two, vestingDuration, 
      vestingStart, vestingCliff, vestingRevocable, {from: account_one}).then((result)=>{
      // Look for NewVesting event in transaction receipt and get vesting id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "NewVesting") {
          return log.args.id.toNumber();
        }
      }
    })
    
    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    let account_two_vesting_balance = await vestingContract.vestingBalanceOf.call(account_two);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - vestingAmount, "Vesting amount should be transfered from sender balance");
    assert.equal(account_two_vesting_balance.toNumber(), vestingAmount, "Vesting amount should be added to the beneficiary vesting balance");
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber(), "Beneficiary main balance should stay unchanged");
    
    // jump in time, third vesting duration
    await increaseTimeTo(latestTime()+vestingDuration/3);
    await vestingContract.releaseVesting(vestingId)
    
    // should release some of vesting to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.isBelow(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + vestingAmount/2, 'Should release some of the vesting to the main balance')

    // jump in time, full vesting duration
    await increaseTimeTo(latestTime()+vestingDuration);
    await vestingContract.releaseVesting(vestingId);

    // should release full vesting amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + vestingAmount, "Should release full vesting amount to the main balance");

    account_two_vesting_balance = await vestingContract.vestingBalanceOf.call(account_two);
    assert.equal(account_two_vesting_balance.toNumber(), 0, "Vesting amount should become 0");
    
  });

  it("should stake vested tokens correctly", async function() {
    
        let vestingAmount = 1000000000;
        let vestingDuration = duration.days(60);
        let vestingStart = latestTime();
        let vestingCliff = duration.days(10);
        let vestingRevocable = true;
    
        // Vest tokens
        await token.approve(vestingContract.address, vestingAmount, {from: account_one});
        let vestingId = await vestingContract.vest(vestingAmount, account_two, vestingDuration, 
          vestingStart, vestingCliff, vestingRevocable, {from: account_one}).then((result)=>{
          // Look for NewVesting event in transaction receipt and get vesting id
          for (var i = 0; i < result.logs.length; i++) {
            var log = result.logs[i];
            if (log.event == "NewVesting") {
              return log.args.id.toNumber();
            }
          }
        })
        
        // should throw if stake vesting called by anyone except vesting beneficiary
        await exceptThrow(vestingContract.stakeVesting(vestingId));

        // stake vesting can be only called by vesting beneficiary
        await vestingContract.stakeVesting(vestingId, {from: account_two});
        let account_two_vesting_stake_balance = await vestingContract.vestingStakeBalanceOf.call(account_two);
        assert.equal(account_two_vesting_stake_balance.toNumber(), vestingAmount, "Should stake vesting amount");

        // should throw if initiate unstake called by anyone except vesting beneficiary
        await exceptThrow(vestingContract.initiateUnstakeVesting(vestingId));

        // Initiate unstake of vested tokens by vesting beneficiary
        let stakeWithdrawalId = await vestingContract.initiateUnstakeVesting(vestingId, {from: account_two}).then((result)=>{
          // Look for InitiateUnstakeVesting event in transaction receipt and get stake withdrawal id
          for (var i = 0; i < result.logs.length; i++) {
            var log = result.logs[i];
            if (log.event == "InitiateUnstakeVesting") {
              return log.args.id.toNumber();
            }
          }
        });

        // should not be able to finish unstake before withdrawal delay is over
        await exceptThrow(vestingContract.finishUnstakeVesting(stakeWithdrawalId));

        // should not be able to release vesting as its still locked for staking
        await exceptThrow(vestingContract.releaseVesting(vestingId));

        // jump in time over withdrawal delay
        await increaseTimeTo(latestTime()+duration.days(30));
        await vestingContract.finishUnstakeVesting(stakeWithdrawalId);
        account_two_vesting_stake_balance = await vestingContract.vestingStakeBalanceOf.call(account_two);
        assert.equal(account_two_vesting_stake_balance.toNumber(), 0, "Stake vesting amount should be 0");

        // should be able to release 'releasable' vesting amount as it's not locked for staking anymore
        await vestingContract.releaseVesting(vestingId);
        let account_two_ending_balance = await token.balanceOf.call(account_two);
        assert.isAbove(account_two_ending_balance.toNumber(), vestingAmount/2, "Should have some released vesting amount");

      });
});
