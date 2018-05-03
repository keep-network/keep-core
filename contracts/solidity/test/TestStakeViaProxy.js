import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');

contract('TestStakeViaProxy', function(accounts) {

  let token, stakingProxy, stakingContract,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
  });

  it("should stake and unstake via staking proxy contract", async function() {
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));

    const proxyStakedEvent = stakingProxy.Staked();
    const proxyUnstakedEvent = stakingProxy.Unstaked();

    let stakingAmount = 10000000;

    // Stake should fail since stakingContract was not added to the proxy
    await exceptThrow(token.approveAndCall(stakingContract.address, stakingAmount, "", {from: account_one}));

    // Non-owner of stakingProxy should not be able to authorize a staking contract
    await exceptThrow(stakingProxy.authorizeContract(stakingContract.address, {from: account_two}));

    // Owner of stakingProxy should be able to authorize a staking contract
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})
    stakingProxy.isAuthorized(stakingContract.address).then(function(result){
      assert.equal(result, true, "StakingProxy owner should be able to authorize a staking contract.");
    })

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, "", {from: account_one});
    proxyStakedEvent.get(function(error, result){
      assert.equal(result[0].event, 'Staked', "Staked event on the proxy contract should occur.");
    });

    // Initiate unstake tokens
    let stakeWithdrawalId = await stakingContract.initiateUnstake(stakingAmount, {from: account_one});
    proxyUnstakedEvent.get(function(error, result){
      assert.equal(result[0].event, 'Unstaked', "Unstaked event on the proxy contract should occur.");
    });

    // Owner of stakingProxy should be able to deauthorize a staking contract
    await stakingProxy.deauthorizeContract(stakingContract.address, {from: account_one})
    stakingProxy.isAuthorized(stakingContract.address).then(function(result){
      assert.equal(result, false, "StakingProxy owner should be able to deauthorize a staking contract.");
    })

  });
});
