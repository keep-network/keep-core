import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');

contract('KeepToken', function(accounts) {

  let token, grantContract, stakingContract, stakingProxy,
    account_one = accounts[0],
    account_one_operator = accounts[1],
    account_one_magpie = accounts[2],
    account_two = accounts[3],
    account_two_operator = accounts[4],
    account_two_magpie = accounts[5];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    grantContract = await TokenGrant.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address);
    await stakingProxy.authorizeContract(grantContract.address);
  });

  it("should send tokens correctly", async function() {
    let amount = web3.utils.toBN(1000000000);

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_two);

    // Send tokens
    await token.transfer(account_two, amount, {from: account_one});

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_two_ending_balance = await token.balanceOf.call(account_two);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount wasn't correctly taken from the sender");
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance.add(amount)), true, "Amount wasn't correctly sent to the receiver");

  });

  it("should stake and unstake tokens correctly", async function() {

    let stakingAmount = web3.utils.toBN(10000000);

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);

    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(account_one), account_one_operator)).substr(2), 'hex');
    let data = Buffer.concat([Buffer.from(account_one_magpie.substr(2), 'hex'), signature]);

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, '0x' + data.toString('hex'), {from: account_one});

    // Ending balances
    let account_one_ending_balance = await token.balanceOf.call(account_one);
    let account_one_operator_stake_balance = await stakingContract.stakeBalanceOf.call(account_one_operator);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(stakingAmount)), true, "Staking amount should be transfered from sender balance");
    assert.equal(account_one_operator_stake_balance.eq(stakingAmount), true, "Staking amount should be added to the sender staking balance");

    // Initiate unstake tokens as token owner
    let stakeWithdrawalId = await stakingContract.initiateUnstake(stakingAmount/2, account_one_operator, {from: account_one}).then((result)=>{
      // Look for initiateUnstake event in transaction receipt and get stake withdrawal id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "InitiatedUnstake") {
          return log.args.id.toNumber();
        }
      }
    })

    // Initiate unstake tokens as operator
    let stakeWithdrawalId2 = await stakingContract.initiateUnstake(stakingAmount/2, account_one_operator, {from: account_one_operator}).then((result)=>{
      // Look for initiateUnstake event in transaction receipt and get stake withdrawal id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "InitiatedUnstake") {
          return log.args.id.toNumber();
        }
      }
    })

    let withdrawals = await stakingContract.getWithdrawals(account_one);
    assert.equal(withdrawals.length, 2, "Withdrawal records must present for the staker");

    // should not be able to finish unstake
    await exceptThrow(stakingContract.finishUnstake(stakeWithdrawalId));

    // jump in time, full withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));

    // should be able to finish unstake
    await stakingContract.finishUnstake(stakeWithdrawalId);
    await stakingContract.finishUnstake(stakeWithdrawalId2);

    withdrawals = await stakingContract.getWithdrawals(account_one);
    assert.equal(withdrawals.length, 0, "Withdrawal record must be cleared for the staker");

    // check balances
    account_one_ending_balance = await token.balanceOf.call(account_one);
    account_one_operator_stake_balance = await stakingContract.stakeBalanceOf.call(account_one_operator);

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance), true, "Staking amount should be transfered to sender balance");
    assert.equal(account_one_operator_stake_balance.isZero(), true, "Staking amount should be removed from sender staking balance");

  });


  it("should grant tokens correctly", async function() {

    let amount = web3.utils.toBN(1000000000);
    let vestingDuration = duration.days(30);
    let start = await latestTime();
    let cliff = duration.days(10);
    let revocable = true;

    // Starting balances
    let account_one_starting_balance = await token.balanceOf.call(account_one);
    let account_two_starting_balance = await token.balanceOf.call(account_one_operator);

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

    assert.equal(account_one_ending_balance.eq(account_one_starting_balance.sub(amount)), true, "Amount should be transfered from sender balance");
    assert.equal(account_two_grant_balance.eq(amount), true, "Amount should be added to the beneficiary grant balance");
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance), true, "Beneficiary main balance should stay unchanged");

    // Should not be able to release token grant (0 unreleased amount)
    await exceptThrow(grantContract.release(id))

    // jump in time, third vesting duration
    await increaseTimeTo(await latestTime()+vestingDuration/3);

    // Should be able to release token grant unreleased amount
    await grantContract.release(id)

    // should release some of grant to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.lte(account_two_starting_balance.add(amount.div(web3.utils.toBN(2)))), true, 'Should release some of the grant to the main balance')

    // jump in time, full vesting duration
    await increaseTimeTo(await latestTime()+vestingDuration);
    await grantContract.release(id);

    // should release full grant amount to the main balance
    account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.eq(account_two_starting_balance.add(amount)), true, "Should release full grant amount to the main balance");

    account_two_grant_balance = await grantContract.balanceOf.call(account_two);
    assert.equal(account_two_grant_balance, 0, "Grant amount should become 0");

  });

  it("should stake granted tokens correctly", async function() {
    let amount = web3.utils.toBN(1000000000);
    let vestingDuration = duration.days(60);
    let start = await latestTime();
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

    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(account_two), account_two_operator)).substr(2), 'hex');
    let delegation = Buffer.concat([Buffer.from(account_two_magpie.substr(2), 'hex'), signature]);

    // should throw if stake granted tokens called by anyone except grant beneficiary
    await exceptThrow(grantContract.stake(id, delegation));

    // stake granted tokens can be only called by grant beneficiary
    await grantContract.stake(id, delegation, {from: account_two});
    let account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.eq(amount), true, "Should stake grant amount");

    // should throw if initiate unstake called by anyone except grant beneficiary
    await exceptThrow(grantContract.initiateUnstake(id, account_two_operator));

    // Initiate unstake of granted tokens by grant beneficiary
    let stakeWithdrawalId = await grantContract.initiateUnstake(id, account_two_operator, {from: account_two}).then((result)=>{
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
    await increaseTimeTo(await latestTime()+duration.days(30));
    await grantContract.finishUnstake(stakeWithdrawalId);
    account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.isZero(), true, "Stake grant amount should be 0");

    // should be able to release 'releasable' granted amount as it's not locked for staking anymore
    await grantContract.release(id);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.gte(amount.div(web3.utils.toBN(2))), true, "Should have some released grant amount");

  });
});
