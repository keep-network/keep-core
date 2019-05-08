import { duration } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const Staking = artifacts.require('./Staking.sol');

contract('TestStakeViaProxy', function(accounts) {

  let token, tokenGrant, stakingProxy, stakingContract,
    owner = accounts[0],
    operator = accounts[1];

  beforeEach(async () => {
    token = await KeepToken.new();
    tokenGrant = await TokenGrant.new(token.address);
  });

  it("should stake and unstake via staking proxy contract", async function() {
    stakingProxy = await StakingProxy.new();
    stakingContract = await Staking.new(token.address, tokenGrant.address, stakingProxy.address, duration.days(30));

    let stakingAmount = 10000000;

    let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
    let delegation = '0x' + Buffer.concat([Buffer.from(owner.substr(2), 'hex'), signature]).toString('hex');

    // Stake should fail since stakingContract was not added to the proxy
    await exceptThrow(token.approveAndCall(stakingContract.address, stakingAmount, delegation, {from: owner}));

    // Non-owner of stakingProxy should not be able to authorize a staking contract
    await exceptThrow(stakingProxy.authorizeContract(stakingContract.address, {from: operator}));

    // Owner of stakingProxy should be able to authorize a staking contract
    await stakingProxy.authorizeContract(stakingContract.address, {from: owner})
    assert.equal(await stakingProxy.isAuthorized(stakingContract.address), true, "StakingProxy owner should be able to authorize a staking contract.");

    // Stake tokens using approveAndCall pattern
    await token.approveAndCall(stakingContract.address, stakingAmount, delegation, {from: owner});
    assert.equal((await stakingProxy.getPastEvents())[0].event, 'Staked', "Staked event on the proxy contract should occur.");

    // Initiate unstake tokens
    await stakingContract.initiateUnstake(stakingAmount, operator, {from: owner});
    assert.equal((await stakingProxy.getPastEvents())[0].event, 'Unstaked', "Unstaked event on the proxy contract should occur.");

    // Owner of stakingProxy should be able to deauthorize a staking contract
    await stakingProxy.deauthorizeContract(stakingContract.address, {from: owner})
    assert.equal(await stakingProxy.isAuthorized(stakingContract.address), false, "StakingProxy owner should be able to deauthorize a staking contract.");

  });
});
