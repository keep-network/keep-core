import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');

contract('TestTokenGrantStake', function(accounts) {

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
    await increaseTimeTo(await latestTime()+duration.days(30));
    await grantContract.finishUnstake(stakeWithdrawalId);
    account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.isZero(), true, "Stake grant amount should be 0");
    assert.equal(await grantContract.operatorsOf.call(account_two), 0, "Operator should be released after finishing unstake");

    // should be able to release 'releasable' granted amount as it's not locked for staking anymore
    await grantContract.release(id);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.gte(amount.div(web3.utils.toBN(2))), true, "Should have some released grant amount");

  });

  it("should be able to delegate stake to the same operator after finishing unstaking", async function() {
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
    await increaseTimeTo(await latestTime()+duration.days(30));
    await grantContract.finishUnstake(stakeWithdrawalId);
    account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.isZero(), true, "Stake grant amount should be 0");
    assert.equal(await grantContract.operatorsOf.call(account_two), 0, "Operator should be released after finishing unstake");

    // should be able to release 'releasable' granted amount as it's not locked for staking anymore
    await grantContract.release(id);
    let account_two_ending_balance = await token.balanceOf.call(account_two);
    assert.equal(account_two_ending_balance.gte(amount.div(web3.utils.toBN(2))), true, "Should have some released grant amount");

    // Grant tokens
    await token.approve(grantContract.address, amount, {from: account_one});
    id = await grantContract.grant(amount, account_two, vestingDuration,
      start, cliff, revocable, {from: account_one}).then((result)=>{
      // Look for CreatedTokenGrant event in transaction receipt and get grant id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "CreatedTokenGrant") {
          return log.args.id.toNumber();
        }
      }
    })

    signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(account_two), account_two_operator)).substr(2), 'hex');
    delegation = Buffer.concat([Buffer.from(account_two_magpie.substr(2), 'hex'), signature]);

    // should throw if stake granted tokens called by anyone except grant beneficiary
    await exceptThrow(grantContract.stake(id, delegation));

    // stake granted tokens can be only called by grant beneficiary
    await grantContract.stake(id, delegation, {from: account_two});
    account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.eq(amount), true, "Should stake grant amount");
  });

  it("should not be able to use TokenGrant.initiateUnstake() to undelegate/unstake not owned stake", async function() {
    let amount = web3.utils.toBN(10000000);
    let amount2 = web3.utils.toBN(1000000);
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

    // stake granted tokens can be only called by grant beneficiary
    await grantContract.stake(id, delegation, {from: account_two});
    let account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    assert.equal(account_two_operator_stake_balance.eq(amount), true, "Should stake grant amount");

    // Grant tokens
    await token.approve(grantContract.address, amount2, {from: account_one});
    let id2 = await grantContract.grant(amount2, account_one, vestingDuration,
      start, cliff, revocable, {from: account_one}).then((result)=>{
      // Look for CreatedTokenGrant event in transaction receipt and get grant id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "CreatedTokenGrant") {
          return log.args.id.toNumber();
        }
      }
    })

    let signature2 = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(account_one), account_one_operator)).substr(2), 'hex');
    let delegation2 = Buffer.concat([Buffer.from(account_one_magpie.substr(2), 'hex'), signature2]);

    // stake granted tokens can be only called by grant beneficiary
    await grantContract.stake(id2, delegation2, {from: account_one});
    let account_one_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_one_operator);
    assert.equal(account_one_operator_stake_balance.eq(amount2), true, "Should stake grant amount");

    await increaseTimeTo(await latestTime()+duration.days(30));

    // Try to initiate unstake of granted tokens by not by grant beneficiary
    await exceptThrow(grantContract.initiateUnstake(id2, {from: account_two}));
  });
});
