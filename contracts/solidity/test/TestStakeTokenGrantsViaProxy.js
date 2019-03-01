import { duration } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');

contract('TestStakeTokenGrantsViaProxy', function(accounts) {

  let token, stakingProxy, grantContract,
    account_one = accounts[0],
    account_two = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
  });

  it("should stake and unstake granted tokens via staking proxy contract", async function() {
    stakingProxy = await StakingProxy.new();
    grantContract = await TokenGrant.new(token.address, stakingProxy.address, duration.days(30));

    const proxyStakedEvent = stakingProxy.Staked();
    const proxyUnstakedEvent = stakingProxy.Unstaked();

    let amount = 1000000000;
    let vestingDuration = duration.days(60);
    let start = await latestTime();
    let cliff = duration.days(10);
    let revocable = true;

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

    // Stake should fail since grantContract was not added to the staking proxy
    await exceptThrow(grantContract.stake(id, {from: account_two}));

    // Non-owner of stakingProxy should not be able to authorize a token grant contract
    await exceptThrow(stakingProxy.authorizeContract(grantContract.address, {from: account_two}));

    // Owner of stakingProxy should be able to authorize a token grant contract
    await stakingProxy.authorizeContract(grantContract.address, {from: account_one})
    assert.equal(await stakingProxy.isAuthorized(grantContract.address), true, "StakingProxy owner should be able to authorize a token grant contract.");

    // Stake granted tokens
    await grantContract.stake(id, {from: account_two})
    assert.equal((await stakingProxy.getPastEvents())[0].event, 'Staked', "Staked event on the proxy contract should occur.");

    // Initiate unstake of granted tokens by grant beneficiary
    await grantContract.initiateUnstake(id, {from: account_two});
    assert.equal((await stakingProxy.getPastEvents())[0].event, 'Unstaked', "Unstaked event on the proxy contract should occur.");

    // Owner of stakingProxy should be able to deauthorize a token grant contract
    await stakingProxy.deauthorizeContract(grantContract.address, {from: account_one})
    assert.equal(await stakingProxy.isAuthorized(grantContract.address), false, "StakingProxy owner should be able to deauthorize a token grant contract.");

  });
});
