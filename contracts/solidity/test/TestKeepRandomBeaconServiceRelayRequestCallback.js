import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import stakeAndGenesis from './helpers/stakeAndGenesis';

const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

contract('KeepRandomBeacon', function(accounts) {

  const groupSize = 3;
  const groupThreshold = 2;

  let operatorContract, serviceContract, callbackContract,
    submitterExpectedReimbursement, submitterExpectedReimbursementGroupSelection, priceFeedEstimate,
    operator = accounts[1], // make sure these match the ones in stakeAndGenesis.js
    beneficiary = accounts[4];

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorCallbackStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

    await operatorContract.setGroupSize(groupSize);
    await operatorContract.setGroupThreshold(groupThreshold);

    await stakeAndGenesis(accounts, contracts);

    priceFeedEstimate = await operatorContract.priceFeedEstimate();
    const entryFeeBreakdown = await serviceContract.entryFeeBreakdown()
    const groupSelectionGasEstimate = await operatorContract.groupSelectionGasEstimate();
    const fluctuationMargin = await operatorContract.fluctuationMargin();
    const gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).divn(100));
    submitterExpectedReimbursement = entryFeeBreakdown.entryVerificationFee;
    submitterExpectedReimbursementGroupSelection = submitterExpectedReimbursement.add(
      groupSelectionGasEstimate.mul(gasPriceWithFluctuationMargin)
    )
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  // Sends to DKG fee pool on the service contract enough ether to start
  // a new group creation.
  async function fundDkgPool() {
    const groupCreationGasEstimate = await operatorContract.groupCreationGasEstimate();
    const fluctuationMargin = await operatorContract.fluctuationMargin();
    const priceFeedEstimate = await serviceContract.priceFeedEstimate();
    const gasPriceWithFluctuationMargin = priceFeedEstimate.add(
      priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100))
    );
    
    await serviceContract.fundDkgFeePool(
      {value: groupCreationGasEstimate.mul(gasPriceWithFluctuationMargin)}
    );
  }

  it("should produce entry when no callback was not provided", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );
  });

  it("should reimburse submitter when no callback was provided", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

    const startBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});
    let submitterReimbursement = web3.utils.toBN((await web3.eth.getBalance(beneficiary))).sub(startBalance);

    assert.isTrue(submitterReimbursement.eq(submitterExpectedReimbursement), "Unexpected submitter reimbursement");
  })

  it("should produce entry and execute callback if provided", async () => {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );

    let result = web3.utils.toBN(await callbackContract.lastEntry());
    assert.isTrue(
      result.eq(bls.groupSignatureNumber), 
      "Unexpected entry value passed to the callback"
    );
  })

  it("should reimburse submitter for executing a callback", async () => {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );

    const startBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});
    const endBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let submitterReimbursement = endBalance.sub(startBalance);

    let expectedReimbursmentGas = (await serviceContract.baseCallbackGas()).addn(callbackGas);
    submitterExpectedReimbursement = submitterExpectedReimbursement.add(expectedReimbursmentGas.mul(priceFeedEstimate));
    assert.isTrue(submitterReimbursement.eq(submitterExpectedReimbursement), "Unexpected submitter reimbursement");
  });

  it("should trigger new group creation and execute callback if provided", async () => {
    fundDkgPool();

    // Make sure DKG is possible
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "Group selection should be possible");

    // Request relay entry with a callback
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );
  
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});
  
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

  it("should trigger new group creation, execute callback, and reimburse submitter", async () => {
    fundDkgPool();

    // Make sure DKG is possible
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "Group selection should be possible");

    // Request relay entry with a callback
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );
  
    let startBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});
    let endBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let submitterReimbursement = endBalance.sub(startBalance);
  
    let expectedReimbursmentGas = (await serviceContract.baseCallbackGas()).addn(callbackGas);
    submitterExpectedReimbursementGroupSelection = submitterExpectedReimbursementGroupSelection.add(
      expectedReimbursmentGas.mul(priceFeedEstimate)
    );

    assert.isTrue(
      submitterReimbursement.eq(submitterExpectedReimbursementGroupSelection), 
      "Unexpected submitter reimbursement"
    );
  })

  it("should trigger new group creation when callback failed", async () => {
    fundDkgPool();

    // Make sure DKG is possible
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "Group selection should be possible");

    // Request relay entry with a callback using wrong gas estimate
    let realCallbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let callbackGas = 1; // Requestor provides wrong gas
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );
    
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});

    assert.equal((await operatorContract.getPastEvents())[1].event,
      'GroupSelectionStarted', "Should start group selection"
    );

    assert.equal((await operatorContract.getPastEvents())[1].args['newEntry'].toString(),
      bls.groupSignatureNumber.toString(), "Should start group selection with new entry"
    );
  
    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );
  })

  it("should trigger new group creation and reimburse submitter when callback failed", async () => {
    fundDkgPool();

    // Make sure DKG is possible
    assert.isTrue(await operatorContract.isGroupSelectionPossible(), "Group selection should be possible");

    // Request relay entry with a callback using wrong gas estimate
    let realCallbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let callbackGas = 1; // Requestor provides wrong gas
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );

    let startBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});
    let endBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let submitterReimbursement = endBalance.sub(startBalance);

    let expectedReimbursmentGas = (await serviceContract.baseCallbackGas()).addn(realCallbackGas);
    submitterExpectedReimbursementGroupSelection = submitterExpectedReimbursementGroupSelection.add(
      expectedReimbursmentGas.mul(priceFeedEstimate)
    );

    assert.isTrue(
      submitterReimbursement.eq(submitterExpectedReimbursementGroupSelection), 
      "Unexpected submitter reimbursement"
    );
  });
});
