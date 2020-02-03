import {bls} from './helpers/data';
import mineBlocks from './helpers/mineBlocks';
import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";

import stakeAndGenesis from './helpers/stakeAndGenesis';

const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('TestKeepRandomBeaconServiceRelayRequestCallback', function(accounts) {

  const groupSize = 3;
  const groupThreshold = 2;

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

    await operatorContract.setGroupSize(groupSize);
    await operatorContract.setGroupThreshold(groupThreshold);

    await stakeAndGenesis(accounts, contracts);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
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
    // Fund DKG pool
    const groupCreationGasEstimate = await operatorContract.groupCreationGasEstimate();
    const fluctuationMargin = await operatorContract.fluctuationMargin();
    const priceFeedEstimate = await serviceContract.priceFeedEstimate();
    const gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)));
    await serviceContract.fundDkgFeePool({value: groupCreationGasEstimate.mul(gasPriceWithFluctuationMargin)});

    // Make sure DKG is possible
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "Group selection should be possible");

    // Request relay entry with a callback
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate});

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
     // Fund DKG pool
     const groupCreationGasEstimate = await operatorContract.groupCreationGasEstimate();
     const fluctuationMargin = await operatorContract.fluctuationMargin();
     const priceFeedEstimate = await serviceContract.priceFeedEstimate();
     const gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)));
     await serviceContract.fundDkgFeePool({value: groupCreationGasEstimate.mul(gasPriceWithFluctuationMargin)});
 
     // Make sure DKG is possible
     assert.isTrue(await operatorContract.isGroupSelectionPossible(), "Group selection should be possible");
 
     // Request relay entry with a callback using wrong gas estimate
     let callbackGas = 1; // wrong gas estimate
     let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
     await serviceContract.methods['requestRelayEntry(address,string,uint256)'](callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate});
 
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
  });
});
