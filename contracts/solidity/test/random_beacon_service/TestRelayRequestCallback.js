import {bls} from '../helpers/data';
import {initContracts} from '../helpers/initContracts';
import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';
import stakeAndGenesis from '../helpers/stakeAndGenesis';

const CallbackContract = artifacts.require('./examples/CallbackContract.sol');

// A set of integration tests for the beacon pricing mechanism related to
// callback reimbursement.
contract('KeepRandomBeacon/RelayRequestCallback', function(accounts) {

  const groupSize = 3;
  const groupThreshold = 2;

  let operatorContract, serviceContract, callbackContract;

  let customer = accounts[0];
  let operator = accounts[1]; // make sure these match the ones in stakeAndGenesis.js
  let beneficiary = accounts[4];

  let entryVerificationFee, dkgContributionFee, groupProfitFee;
  let baseCallbackGas;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./stubs/TokenStakingStub.sol'),
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


    let feeBreakdown = await serviceContract.entryFeeBreakdown();
    entryVerificationFee = web3.utils.toBN(feeBreakdown.entryVerificationFee);
    dkgContributionFee = web3.utils.toBN(feeBreakdown.dkgContributionFee);
    groupProfitFee = web3.utils.toBN(feeBreakdown.groupProfitFee);

    baseCallbackGas = web3.utils.toBN(await serviceContract.baseCallbackGas())
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

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

    const beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {from: operator});
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    assert.isTrue(
      submitterReimbursement.eq(entryVerificationFee), 
      "Unexpected submitter reimbursement"
    );
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

  // price feed estimate > tx.gasprice
  it("should reimburse submitter and customer for executing callback with lower tx.gasprice", async () => {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use lower gas price when submitting entry
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('2', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, relayEntryTxGasPrice
    );
  });

  // price feed estimate == tx.gasprice
  it("should reimburse submitter and customer for executing callback with expected tx.gasprice", async () => {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use the same gas price as price feed
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('20', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, relayEntryTxGasPrice
    );
  });

  // price feed estimate < tx.gasprice
  it("should reimburse submitter and customer for executing callback with higher tx.gasprice", async () => {
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use higher price than the price feed
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('30', 'Gwei'));
    // higher tx.gasprice should not be used for reimbursement - maximum gas
    // price is the one from price feed
    let gasPriceForReimbursement = web3.utils.toBN(web3.utils.toWei('20', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice// tx.gasprice is the same as price feed estimate
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, gasPriceForReimbursement
    );
  })

  it("should trigger new group creation and execute callback if provided", async () => {
    await fundDkgPool();

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

  // price feed estimate == tx.gasprice
  it("should trigger new group creation, execute callback, reimburse submitter " + 
     "and customer with expected tx.gasprice", async () => {
    await fundDkgPool();

    // Request relay entry with a callback
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {
        value: entryFeeEstimate
      }
    );

    // use the same gas price as price feed
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('20', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackAndDKGReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, relayEntryTxGasPrice
    );
  })

  // price feed estimate > tx.gasprice
  it("should trigger new group selection, execute callback, reimburse submitter " +
     "and customer with lower tx.gasprice", async () => {
    await fundDkgPool();

    // Request relay entry with a callback
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {
        value: entryFeeEstimate
      }
    );
  
    // use lower gas price when submitting entry
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('2', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackAndDKGReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, relayEntryTxGasPrice
    );
  })
 
  // price feed estimate < tx.gasprice
  it("should trigger new group selection, execute callback, reimburse submitter " +
    "and customer with higher tx.gasprice", async () => {
    await fundDkgPool();
  
    // Request relay entry with a callback
    let callbackGas = await callbackContract.callback.estimateGas(bls.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {
        value: entryFeeEstimate
      }
    );
  
    // use higher price than the price feed
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('30', 'Gwei'));
    // higher tx.gasprice should not be used for reimbursement - maximum gas
    // price is the one from price feed
    let gasPriceForReimbursement = web3.utils.toBN(web3.utils.toWei('20', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackAndDKGReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, gasPriceForReimbursement
    );
  });

  it("should trigger new group creation when callback failed", async () => {
    await fundDkgPool();

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

  it("should reimburse submitter when callback failed", async () => {
    let estimate = await callbackContract.callback.estimateGas(bls.groupSignature);
    
    let callbackGas = estimate - 10; // Requestor provides wrong gas
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,string,uint256)'](
      callbackContract.address, "callback(uint256)", callbackGas, {value: entryFeeEstimate}
    );
  
    // use the same gas price as price feed
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('20', 'Gwei'));
  
    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(bls.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
  
    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);
  
    await assertCallbackReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement, 
      customerSurplus, relayEntryTxGasPrice
    );
  });

  // This function assets expexced submitter reimbursement and customer surplus
  // in a situation when callback is executed and no new group creation is
  // triggered.
  // 
  // Function performs two checks:
  //
  // 1. It makes sure beneficiary account balance change (submitter 
  // reward + callback  reimbursement) and customer account balance
  // change (surplus) add up to entry verification fee and callback fee as
  // calculated by the service contract. In other words, beacon does not spend
  // more than received and all required contributions (DKG fee + group profit 
  // fee) stay in the contract.
  //
  // 2. It calculates the expected gas cost of executing the callback and assert
  // it is the same as beneficiary account balance change without
  // entry verification fee.
  async function assertCallbackReimbursement(
    callbackGas, entryFeeEstimate, submitterReimbursement, 
    customerSurplus, txGasPrice
  ) {
    let totalSpentByBeacon = submitterReimbursement.add(customerSurplus);

    // calculate part the fee used for entry verification and callback
    let entryVerificationAndCallbackFee = web3.utils.toBN(entryFeeEstimate)
      .sub(dkgContributionFee)
      .sub(groupProfitFee)

    // the sum of ether paid to beneficiary and customer should equal
    // entry verification and callback fee passed to the beacon 
    assert.isTrue(
      entryVerificationAndCallbackFee.eq(totalSpentByBeacon), 
      "Beacon spent more than allowed"
    );  

    // expected and actual gas reimbursement should be _almost_ the same;
    // see the function doc for explanation about additional EVM opcodes
    // around 'call'
    let expectedCallbackReimbursementGas = baseCallbackGas.addn(callbackGas);
    let actualCallbackReimbursementGas = submitterReimbursement
      .sub(entryVerificationFee)
      .div(txGasPrice);

    assert.equal(      
      actualCallbackReimbursementGas.toNumber(), 
      expectedCallbackReimbursementGas.toNumber(),
      "Unexpected callback reimbursement"
    )
  }

  // This function asserts expected submitter reimbursement and customer surplus
  // in a situation when callback is executed and a new group creation is
  // triggered.
  //
  // Function performs two checks:
  //
  // 1. It makes sure beneficiary account balance change (submitter 
  // reward + callback  reimbursement + start group creation reimbursement) and 
  // customer account balance change (surplus) add up to entry verification fee,
  // start group creation fee  and callback fee as calculated by the service 
  // contract. In other words, beacon does not spend more than received and
  // all required contributions (DKG fee + group profit fee) stay in the
  // contract.
  //
  // 2. It calculates the expected gas cost of executing the callback and assert
  // it is the same as beneficiary account balance change without
  // entry verification fee and without group creation cost.
  async function assertCallbackAndDKGReimbursement(
    callbackGas, entryFeeEstimate, submitterReimbursement, 
    customerSurplus, txGasPrice
  ) {
    let totalSpentByBeacon = submitterReimbursement.add(customerSurplus);

    // calculate part the fee used for entry verification, group creation,
    // and callback
    let entryVerificationAndCallbackFee = web3.utils.toBN(entryFeeEstimate)
      .sub(dkgContributionFee)
      .sub(groupProfitFee)

    let groupSelectionGasEstimate = web3.utils.toBN(
      await operatorContract.groupSelectionGasEstimate()
    )
    let priceFeedEstimate = web3.utils.toBN(
      await operatorContract.priceFeedEstimate()
    )
    let groupCreationFee = priceFeedEstimate.mul(groupSelectionGasEstimate.add(
      groupSelectionGasEstimate.muln(50).divn(100) // with fluctuation margin
    ))

    // the sum of ether paid to beneficiary and customer should equal
    // entry verification, group creation, and callback fee passed to the beacon 
    let expectedTotalSpent = entryVerificationAndCallbackFee.add(groupCreationFee);
    assert.isTrue(
      expectedTotalSpent.eq(totalSpentByBeacon), 
      "Beacon spent more than allowed"
    );  

    // expected and actual gas reimbursement should be _almost_ the same;
    // see the function doc for explanation about additional EVM opcodes
    // around 'call'
    let expectedCallbackReimbursementGas = baseCallbackGas.addn(callbackGas);
    let actualCallbackReimbursementGas = submitterReimbursement
      .sub(entryVerificationFee)
      .sub(groupCreationFee)
      .div(txGasPrice);

    assert.equal(      
      actualCallbackReimbursementGas.toNumber(), 
      expectedCallbackReimbursementGas.toNumber(),
      "Unexpected callback reimbursement"
    )
  }

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
});
