import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServiceRelayRequestCallback', function(accounts) {

  let operatorContract, serviceContract, callbackContract;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
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
    let group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.addGroupMember(group, accounts[0]);
  });

  it("should produce entry if callback contract was not provided", async function() {
    let minimumPayment = await serviceContract.minimumPayment()
    await serviceContract.requestRelayEntry(bls.seed, {value: minimumPayment});
    await operatorContract.relayEntry(bls.nextGroupSignature);

    let result = await serviceContract.previousEntry();
    assert.isTrue(result.eq(bls.nextGroupSignature), "Value should be updated on beacon contract.");
  });

  it("should successfully call method on a callback contract", async function() {
    let minimumPayment = await serviceContract.minimumPayment()
    await serviceContract.methods['requestRelayEntry(uint256,address,string)'](bls.seed, callbackContract.address, "callback(uint256)", {value: minimumPayment});

    await operatorContract.relayEntry(bls.nextNextGroupSignature);

    let result = await callbackContract.lastEntry();
    assert.isTrue(result.eq(bls.nextNextGroupSignature), "Value updated by the callback should be the same as relay entry.");
  });

});
