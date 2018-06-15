import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');


contract('TestStakeDelegate', function(accounts) {

  let token, stakingProxy, stakingContract,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2],
    account_four = accounts[3];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})

    // Stake tokens as account one
    await token.approveAndCall(stakingContract.address, 200, "", {from: account_one});

    // Send tokens to the accounts
    await token.transfer(account_two, 200, {from: account_one});
    await token.transfer(account_three, 500, {from: account_one});
  
    // Stake tokens as account two
    await token.approveAndCall(stakingContract.address, 200, "", {from: account_two});
  });

  it("should not be able to delegate stake to a delegate address that is a staker", async function() {
    await stakingContract.requestDelegateFor(account_one, {from: account_two});
    await exceptThrow(stakingContract.approveDelegateAt(account_two));
  });

  it("should be able to delegate stake to a delegate address to represent your stake balance", async function() {
    await stakingContract.requestDelegateFor(account_one, {from: account_three});
    await stakingContract.approveDelegateAt(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Delegate account should represent delegator's stake balance.");
  });

  it("should not be able to delegate stake to a delegate address that is not approved", async function() {
    await stakingContract.requestDelegateFor(account_two, {from: account_three});
    await stakingContract.approveDelegateAt(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 0, "Delegate account should be zero since there were no handshake with delegator.");
  });

  it("should be able to update delegate address if new delegate request exist", async function() {
    await stakingContract.requestDelegateFor(account_one, {from: account_three});
    await stakingContract.approveDelegateAt(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Delegate account should represent delegator's stake balance.");

    await stakingContract.approveDelegateAt(account_four, {from: account_one});
    await stakingContract.requestDelegateFor(account_one, {from: account_four});
    assert.equal(await stakingProxy.balanceOf(account_three), 0, "Previous delegate account should stop representing delegator's stake balance.");
    assert.equal(await stakingProxy.balanceOf(account_four), 200, "Updated delegate account should represent delegator's stake balance.");
  });

  it("should be able to remove delegated address that represents your stake balance", async function() {
    await stakingContract.requestDelegateFor(account_one, {from: account_three});
    await stakingContract.approveDelegateAt(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Delegate account should represent delegator's stake balance.");
    await stakingContract.removeDelegate();
    assert.equal(await stakingProxy.balanceOf(account_three), 0, "Delegate account should stop representing delegator's stake balance.");
  });

  it("should be able to change stake and get delegate to reflect updated balance", async function() {
    await stakingContract.requestDelegateFor(account_one, {from: account_three});
    await stakingContract.approveDelegateAt(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Delegate account should represent delegator's stake balance.");

    // Stake more tokens
    await token.approveAndCall(stakingContract.address, 100, "", {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 300, "Delegate account should reflect delegator's updated stake balance.");

    // Unstake everything
    await stakingContract.initiateUnstake(300);
    assert.equal(await stakingProxy.balanceOf(account_three), 0, "Delegate account should reflect delegator's updated stake balance.");
  });

  it("should revert if delegate tries to stake", async function() {
    await stakingContract.requestDelegateFor(account_one, {from: account_three});
    await stakingContract.approveDelegateAt(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Delegate account should represent delegator's stake balance.");

    // Stake tokens as account three
    await exceptThrow(token.approveAndCall(stakingContract.address, 500, "", {from: account_three}));
  });
});
