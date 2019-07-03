import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServiceRelayRequestCallback', function(accounts) {

  let operatorContract, serviceContract, callbackContract;

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
  });

  it("should produce entry if callback contract was not provided", async function() {
    await serviceContract.requestRelayEntry(bls.seed, {value: 10});
    let signingId = (await operatorContract.getPastEvents())[0].args['signingId'].toNumber();
    await operatorContract.relayEntry(signingId, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    let result = await serviceContract.previousEntry();
    assert.isTrue(result.eq(bls.groupSignature), "Value should be updated on beacon contract.");
  });

  it("should successfully call method on a callback contract", async function() {
    await serviceContract.methods['requestRelayEntry(uint256,address,string)'](bls.seed, callbackContract.address, "callback(uint256)", {value: 10});
    let signingId = (await operatorContract.getPastEvents())[0].args['signingId'].toNumber();

    let result = await callbackContract.lastEntry();
    assert.isFalse(result.eq(bls.groupSignature), "Entry value on the callback contract should not be the same as next relay entry.");

    await operatorContract.relayEntry(signingId, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    result = await callbackContract.lastEntry();
    assert.isTrue(result.eq(bls.groupSignature), "Value updated by the callback should be the same as relay entry.");
  });

});
