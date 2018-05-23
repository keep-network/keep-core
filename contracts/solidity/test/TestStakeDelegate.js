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

    // Send tokens to account two
    await token.transfer(account_two, 200, {from: account_one});

    // Stake tokens as account two
    await token.approveAndCall(stakingContract.address, 200, "", {from: account_two});
  });

  it("should not be able to delegate stake to an 'operator' address that is a staker", async function() {
    await exceptThrow(stakingContract.delegateStakeTo(account_two));
  });

  it("should be able to delegate stake to an 'operator' address to represent your stake balance", async function() {
    await stakingContract.delegateStakeTo(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Operator account should represent delegator's stake balance.");
  });

  it("should not be able to delegate stake to an 'operator' address that is already in use", async function() {
    await stakingContract.delegateStakeTo(account_three, {from: account_one});
    await exceptThrow(stakingContract.delegateStakeTo(account_three, {from: account_two}));
  });

  it("should be able to update 'operator' address", async function() {
    await stakingContract.delegateStakeTo(account_three, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Operator account should represent delegator's stake balance.");

    await stakingContract.delegateStakeTo(account_four, {from: account_one});
    assert.equal(await stakingProxy.balanceOf(account_three), 0, "Previous operator account should stop representing delegator's stake balance.");
    assert.equal(await stakingProxy.balanceOf(account_four), 200, "Updated operator account should represent delegator's stake balance.");
  });

  it("should be able to remove delegated operator address that represents your stake balance", async function() {
    await stakingContract.delegateStakeTo(account_three);
    assert.equal(await stakingProxy.balanceOf(account_three), 200, "Operator account should represent delegator's stake balance.");
    await stakingContract.removeDelegate();
    assert.equal(await stakingProxy.balanceOf(account_three), 0, "Operator account should stop representing delegator's stake balance.");
  });
});
