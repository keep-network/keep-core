import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
import expectThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');

contract('KeepTokenBoundless', function(accounts) {

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

  it("should we be able to use TokenStaking.finishUnstake() to block competing operator/owner address?", async function() {

    let stakingAmount = web3.utils.toBN(100);

    // Send tokens
    await token.transfer(account_two, stakingAmount, {from: account_one});

    // Starting balances
    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(account_one), account_one_operator)).substr(2), 'hex');
    let data = Buffer.concat([Buffer.from(account_one_magpie.substr(2), 'hex'), signature]);
    
    let signature2 = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(account_two), account_two_operator)).substr(2), 'hex');
    let data2 = Buffer.concat([Buffer.from(account_two_magpie.substr(2), 'hex'), signature2]);

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, '0x' + data.toString('hex'), {from: account_one});
    await token.approveAndCall(stakingContract.address, stakingAmount, '0x' + data2.toString('hex'), {from: account_two});
    
    // Initiate unstake tokens as token owner
    let stakeWithdrawalId = await stakingContract.initiateUnstake(stakingAmount, account_one_operator, {from: account_one}).then((result)=>{
      // Look for initiateUnstake event in transaction receipt and get stake withdrawal id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "InitiatedUnstake") {
          return log.args.id.toNumber();
        }
      }
    })

    // // jump in time, full withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));

    // // should be able to finish unstake
    await stakingContract.finishUnstake(stakeWithdrawalId, account_two_operator, {from: account_two}); // <- not an operator for this stake, this will block account_two
    await exceptThrow(stakingContract.initiateUnstake(stakingAmount, account_two_operator, {from: account_two})); // <- account_two is blocked, cannot initiate unstake (or do any other operations), cause account_two_operator address was already released
  });

  it("should we be able to use TokenGrant.initiateUnstake() to unstake stake not our delegation?", async function() {
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
    
    // Initiate unstake of granted tokens by grant beneficiary
    let stakeWithdrawalId = await grantContract.initiateUnstake(id2, account_two_operator, {from: account_two}).then((result)=>{ // <- id2 is for account_one!!!
      // Look for InitiatedTokenGrantUnstake event in transaction receipt and get stake withdrawal id
      for (var i = 0; i < result.logs.length; i++) {
        var log = result.logs[i];
        if (log.event == "InitiatedTokenGrantUnstake") {
          return log.args.id.toNumber();
        }
      }
    });

    // jump in time over withdrawal delay
    await increaseTimeTo(await latestTime()+duration.days(30));
    await grantContract.finishUnstake(stakeWithdrawalId, account_two_operator, {from: account_two_operator}); // <- finishes unstake of not owned delegation succesfully 
    //await grantContract.finishUnstake(stakeWithdrawalId, {from: account_two_operator}); // <- finishes unstake of not owned delegation succesfully
    account_two_operator_stake_balance = await grantContract.stakeBalanceOf.call(account_two_operator);
    console.log("After finishUnstake of account_one delegation acount_two_operator_stake_balance =", Number(account_two_operator_stake_balance));
    assert.equal(account_two_operator_stake_balance.isZero(), false, "Should stake grant amount be equal to amount - amount2?");
    assert.equal(await grantContract.operatorsOf.call(account_two), 0, "Operator should be released after finishing unstake");

    await expectThrow(grantContract.initiateUnstake(id2, account_one_operator, {from: account_one})); // <-  account_one is blocked
  });
});
