import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
import exceptThrow from './helpers/expectThrow';
const KeepToken = artifacts.require('./KeepToken.sol');
const StakingProxy = artifacts.require('./StakingProxy.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const Proxy = artifacts.require('./KeepRandomBeacon.sol');
const KeepRandomBeaconImplV1 = artifacts.require('./KeepRandomBeaconImplV1.sol');

contract('TestKeepRandomBeaconDelegateStake', function(accounts) {

  let token, stakingProxy, stakingContract, implV1, proxy, implViaProxy,
    account_one = accounts[0],
    account_two = accounts[1],
    account_three = accounts[2],
    account_four = accounts[3];

  beforeEach(async () => {
    token = await KeepToken.new();
    stakingProxy = await StakingProxy.new();
    stakingContract = await TokenStaking.new(token.address, stakingProxy.address, duration.days(30));
    await stakingProxy.authorizeContract(stakingContract.address, {from: account_one})

    implV1 = await KeepRandomBeaconImplV1.new();
    proxy = await Proxy.new('v1', implV1.address);
    implViaProxy = await KeepRandomBeaconImplV1.at(proxy.address);
    await implViaProxy.initialize(stakingProxy.address, 100, 200, duration.days(30));

    // Stake tokens as account one
    await token.approveAndCall(stakingContract.address, 200, "", {from: account_one});

    // Send tokens to account two
    await token.transfer(account_two, 200, {from: account_one});

    // Stake tokens as account two
    await token.approveAndCall(stakingContract.address, 200, "", {from: account_two});
  });

  it("should not be able to delegate stake to an 'operator' address that is a staker", async function() {
    await exceptThrow(implViaProxy.delegateStakeTo(account_two));
  });

  it("should not be able to delegate stake to an 'operator' address that is already in use", async function() {
    await implViaProxy.delegateStakeTo(account_three);
    await exceptThrow(implViaProxy.delegateStakeTo(account_three, {from: account_two}));
  });

  it("should be able to delegate stake to an 'operator' address to represet your stake balance", async function() {
    await implViaProxy.delegateStakeTo(account_three);
    assert.equal(await implViaProxy.hasMinimumStake(account_three), true, "Operator account should represent delegator's stake balance.");
  });

  it("should be able to update 'operator' address", async function() {
    await implViaProxy.delegateStakeTo(account_three);
    assert.equal(await implViaProxy.hasMinimumStake(account_three), true, "Operator account should represent delegator's stake balance.");

    await implViaProxy.delegateStakeTo(account_four);
    assert.equal(await implViaProxy.hasMinimumStake(account_three), false, "Previous operator account should stop representing delegator's stake balance.");
    assert.equal(await implViaProxy.hasMinimumStake(account_four), true, "Updated operator account should represent delegator's stake balance.");
  });

  it("should be able to remove delagated operator address that represent your stake balance", async function() {
    await implViaProxy.delegateStakeTo(account_three);
    assert.equal(await implViaProxy.hasMinimumStake(account_three), true, "Operator account should represent delegator's stake balance.");
    await implViaProxy.removeDelegate();
    assert.equal(await implViaProxy.hasMinimumStake(account_three), false, "Operator account should stop representing delegator's stake balance.");
  });
});
