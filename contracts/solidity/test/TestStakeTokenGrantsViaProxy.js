import { duration } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const Staking = artifacts.require('./Staking.sol');

contract('TestStakeTokenGrantsViaProxy', function(accounts) {

  let token, tokenGrant, stakingProxy, stakingContract,
    account_one = accounts[0],
    account_two = accounts[1],
    beneficiary = accounts[2],
    magpie = accounts[3],
    operator = accounts[4];

  beforeEach(async () => {
    token = await KeepToken.new();
    tokenGrant = await TokenGrant.new(token.address);
    stakingProxy = await StakingProxy.new();
    stakingContract = await Staking.new(token.address, tokenGrant.address, stakingProxy.address, duration.days(30));
  });

  it("should stake and unstake granted tokens via staking proxy contract", async function() {

    const proxyStakedEvent = stakingProxy.Staked();
    const proxyUnstakedEvent = stakingProxy.Unstaked();

    let amount = web3.utils.toBN(300);
    let vestingDuration = duration.days(60);
    let start = await latestTime();
    let cliff = duration.days(10);
    let revocable = true;

    // Grant tokens
    await token.approve(tokenGrant.address, amount, {from: account_one});
    let id = (await tokenGrant.grant(amount, beneficiary, vestingDuration,
      start, cliff, true, {from: account_one})).logs[0].args.id

    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(beneficiary), operator)).substr(2), 'hex');
    let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

    // Stake should fail since tokenGrant was not added to the staking proxy
    await exceptThrow(tokenGrant.approveAndCall(stakingContract.address, id, amount, delegation, {from: beneficiary}));

    // Non-owner of stakingProxy should not be able to authorize a token grant contract
    await exceptThrow(stakingProxy.authorizeContract(tokenGrant.address, {from: account_two}));

    // Owner of stakingProxy should be able to authorize a token grant contract
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one});
    await stakingProxy.authorizeContract(tokenGrant.address, {from: account_one})
    assert.equal(await stakingProxy.isAuthorized(tokenGrant.address), true, "StakingProxy owner should be able to authorize a token grant contract.");

    // Stake token grant using approveAndCall pattern
    let newGrantId = await tokenGrant.approveAndCall(stakingContract.address, id, amount, delegation, {from: beneficiary}).then((tx)=>{
      // Look for CreatedTokenGrant event in transaction receipt and get new grant id
      for (var i = 0; i < tx.logs.length; i++) {
        var log = tx.logs[i];
        if (log.event == "CreatedTokenGrant") {
          return log.args.id.toNumber();
        }
      }
    })

    assert.equal((await stakingProxy.getPastEvents())[0].event, 'Staked', "Staked event on the proxy contract should occur.");

    // Initiate unstake of granted tokens by grant beneficiary
    await stakingContract.initiateUnstake(id, operator, {from: beneficiary});
    assert.equal((await stakingProxy.getPastEvents())[0].event, 'Unstaked', "Unstaked event on the proxy contract should occur.");

    // Owner of stakingProxy should be able to deauthorize a token grant contract
    await stakingProxy.deauthorizeContract(tokenGrant.address, {from: account_one})
    assert.equal(await stakingProxy.isAuthorized(tokenGrant.address), false, "StakingProxy owner should be able to deauthorize a token grant contract.");

  });
});
