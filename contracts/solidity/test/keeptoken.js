import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');

contract('KeepToken', function(accounts) {

  let token, grantContract, stakingContract,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingContract  = await TokenStaking.new(token.address, duration.days(30));
    grantContract  = await TokenGrant.new(token.address, duration.days(30));
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

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, "", {from: account_one});

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_one_stake_balance = await stakingContract.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - stakingAmount, "Staking amount should be transfered from sender balance");
    assert.equal(account_one_stake_balance.toNumber(), stakingAmount, "Staking amount should be added to the sender staking balance");
    
    // Initiate unstake tokens
    let stakeWithdrawalId = await stakingContract.initiateUnstake(stakingAmount, {from: account_one}).then((result)=>{
      // Look for initiateUnstake event in transaction receipt and get stake withdrawal id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "InitiatedUnstake") {
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
    account_one_stake_balance = await stakingContract.balanceOf.call(account_one);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber(), "Staking amount should be transfered to sender balance");
    assert.equal(account_one_stake_balance.toNumber(), 0, "Staking amount should be removed from sender staking balance");

  });


  it("should grant tokens correctly", async function() {

    let amount = 1000000000;
    let vestingDuration = duration.days(30);
    let start = latestTime();
    let cliff = duration.days(10);
    let revocable = true;
    
    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Grant tokens
    await token.approve(grantContract.address, amount, {from: account_one});
    let id = await grantContract.grant(amount, account_two, vestingDuration, 
      start, cliff, revocable, {from: account_one}).then((result)=>{
      // Look for CreatedTokenGrant event in transaction receipt and get vesting id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "CreatedTokenGrant") {
          return log.args.id.toNumber();
        }
      }
    })
    
    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    let account_two_grant_balance = await grantContract.balanceOf.call(account_two);

    assert.equal(account_one_ending_balance.toNumber(), account_one_starting_balance.toNumber() - amount, "Amount should be transfered from sender balance");
    assert.equal(account_two_grant_balance.toNumber(), amount, "Amount should be added to the beneficiary grant balance");
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber(), "Beneficiary main balance should stay unchanged");
    
    // jump in time, third vesting duration
    await increaseTimeTo(latestTime()+vestingDuration/3);
    await grantContract.release(id)
    
    // should release some of grant to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.isBelow(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + amount/2, 'Should release some of the grant to the main balance')

    // jump in time, full vesting duration
    await increaseTimeTo(latestTime()+vestingDuration);
    await grantContract.release(id);

    // should release full grant amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.toNumber(), account_two_starting_balance.toNumber() + amount, "Should release full grant amount to the main balance");

    account_two_grant_balance = await grantContract.balanceOf.call(account_two);
    assert.equal(account_two_grant_balance.toNumber(), 0, "Grant amount should become 0");
    
  });

  it("should stake granted tokens correctly", async function() {
    
        let amount = 1000000000;
        let vestingDuration = duration.days(60);
        let start = latestTime();
        let cliff = duration.days(10);
        let revocable = true;
    
        // Grant tokens
        await token.approve(grantContract.address, amount, {from: account_one});
        let id = await grantContract.grant(amount, account_two, vestingDuration, 
          start, cliff, revocable, {from: account_one}).then((result)=>{
          // Look for CreatedTokenGrant event in transaction receipt and get grant id
          for (var i = 0; i < result.logs.length; i++) {
            var log = result.logs[i];
            if (log.event == "CreatedTokenGrant") {
              return log.args.id.toNumber();
            }
          }
        })
        
        // should throw if stake granted tokens called by anyone except grant beneficiary
        await exceptThrow(grantContract.stake(id));

        // stake granted tokens can be only called by grant beneficiary
        await grantContract.stake(id, {from: account_two});
        let account_two_grant_stake_balance = await grantContract.stakeBalanceOf.call(account_two);
        assert.equal(account_two_grant_stake_balance.toNumber(), amount, "Should stake grant amount");

        // should throw if initiate unstake called by anyone except grant beneficiary
        await exceptThrow(grantContract.initiateUnstake(id));

        // Initiate unstake of granted tokens by grant beneficiary
        let stakeWithdrawalId = await grantContract.initiateUnstake(id, {from: account_two}).then((result)=>{
          // Look for InitiatedTokenGrantUnstake event in transaction receipt and get stake withdrawal id
          for (var i = 0; i < result.logs.length; i++) {
            var log = result.logs[i];
            if (log.event == "InitiatedTokenGrantUnstake") {
              return log.args.id.toNumber();
            }
          }
        });

        // should not be able to finish unstake before withdrawal delay is over
        await exceptThrow(grantContract.finishUnstake(stakeWithdrawalId));

        // should not be able to release grant as its still locked for staking
        await exceptThrow(grantContract.release(id));

        // jump in time over withdrawal delay
        await increaseTimeTo(latestTime()+duration.days(30));
        await grantContract.finishUnstake(stakeWithdrawalId);
        account_two_grant_stake_balance = await grantContract.stakeBalanceOf.call(account_two);
        assert.equal(account_two_grant_stake_balance.toNumber(), 0, "Stake grant amount should be 0");

        // should be able to release 'releasable' granted amount as it's not locked for staking anymore
        await grantContract.release(id);
        let account_two_ending_balance = await token.balanceOf.call(account_two);
        assert.isAtLeast(account_two_ending_balance.toNumber(), amount/2, "Should have some released grant amount");

      });
});
