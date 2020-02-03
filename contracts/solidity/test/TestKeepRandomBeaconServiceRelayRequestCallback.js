import {bls} from './helpers/data';
import mineBlocks from './helpers/mineBlocks';
import {initContracts} from './helpers/initContracts';
const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServiceRelayRequestCallback', function(accounts) {

  let operatorContract, serviceContract, callbackContract;

  beforeEach(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    operatorContract.setGroupSize(3);
    let group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.setGroupMembers(group, [accounts[0], accounts[1], accounts[2]]);
  });

  it("should produce entry if callback contract was not provided", async function() {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
    await operatorContract.relayEntry(bls.groupSignature);

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );
  });

  it("should successfully call method on a callback contract", async function() {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate});

    await operatorContract.relayEntry(bls.groupSignature);

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );

    let result = web3.utils.toBN(await callbackContract.lastEntry());
    assert.isTrue(
      result.eq(bls.groupSignatureNumber), 
      "Unexpected entry value passed to the callback"
    );
  });

  it("should successfully call method on a callback contract and trigger new group creation", async function() {
    mineBlocks(130); // Make sure dkgTimeout passed so relay entry can start group selection
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "GroupSelectionPossible");
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate});

    // Fund DKG pool
    const groupCreationGasEstimate = await operatorContract.groupCreationGasEstimate();
    const fluctuationMargin = await operatorContract.fluctuationMargin();
    const priceFeedEstimate = await serviceContract.priceFeedEstimate();
    const gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)));
    await serviceContract.fundDkgFeePool({value: groupCreationGasEstimate.mul(gasPriceWithFluctuationMargin)});

    await operatorContract.relayEntry(bls.groupSignature);

    assert.equal((await operatorContract.getPastEvents())[1].event,
      'GroupSelectionStarted', "Should start group selection"
    );

    assert.equal((await operatorContract.getPastEvents())[1].args['newEntry'].toString(),
      bls.groupSignatureNumber.toString(), "Should start group selection with new entry"
    );
    
    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );

    let result = web3.utils.toBN(await callbackContract.lastEntry());
    assert.isTrue(
      result.eq(bls.groupSignatureNumber),
      "Unexpected entry value passed to the callback"
    );
  });

  it("should submit relay entry and trigger new group creation with failed callback", async function() {
    mineBlocks(130); // Make sure dkgTimeout passed so relay entry can start group selection
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "GroupSelectionPossible");
    let callbackGas = 1; // Wrong gas estimation
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas) 
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate});

    // Fund DKG pool
    const groupCreationGasEstimate = await operatorContract.groupCreationGasEstimate();
    const fluctuationMargin = await operatorContract.fluctuationMargin();
    const priceFeedEstimate = await serviceContract.priceFeedEstimate();
    const gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)));
    await serviceContract.fundDkgFeePool({value: groupCreationGasEstimate.mul(gasPriceWithFluctuationMargin)});

    await operatorContract.relayEntry(bls.groupSignature);

    assert.equal((await operatorContract.getPastEvents())[1].event,
      'GroupSelectionStarted', "Should start group selection"
    );

    assert.equal((await operatorContract.getPastEvents())[1].args['newEntry'].toString(),
      bls.groupSignatureNumber.toString(), "Should start group selection with new entry"
    );

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );

    let result = web3.utils.toBN(await callbackContract.lastEntry());
      assert.isFalse(
      result.eq(bls.groupSignatureNumber),
      "Unexpected callback"
    );
  });
});
