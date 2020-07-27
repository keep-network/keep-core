const stakeAndGenesis = require('../helpers/stakeAndGenesis')
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const blsData = require("../helpers/data.js")
const {initContracts} = require('../helpers/initContracts')
const {contract, web3, accounts} = require("@openzeppelin/test-environment")
const {expectRevert} = require("@openzeppelin/test-helpers")

const CallbackContract = contract.fromArtifact('CallbackContract');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect
const assert = chai.assert

// A set of integration tests for the beacon pricing mechanism related to
// callback reimbursement.
describe('KeepRandomBeacon/RelayRequestCallback', function() {
  let operatorContract, serviceContract, callbackContract;

  let customer = accounts[0];
  let operator = accounts[1]; // make sure these match the ones in stakeAndGenesis.js
  let beneficiary = accounts[4];

  let entryVerificationFee, dkgContributionFee, groupProfitFee;
  let baseCallbackGas;

  before(async () => {
    let contracts = await initContracts(
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorCallbackStub')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;
    callbackContract = await CallbackContract.new();

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

  it("should revert when callback gas exceeds gas limit", async () => {
    const callbackGas = 2000001;
    await expectRevert(
      serviceContract.entryFeeEstimate(callbackGas),
      "Callback gas exceeds 2000000 gas limit"
    );

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(2000000)
    await expectRevert(
      serviceContract.methods['requestRelayEntry(address,uint256)'](
        callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
      ),
      "Callback gas exceeds 2000000 gas limit"
    );
  });

  it("should produce entry when no callback was not provided", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: customer});
    await operatorContract.relayEntry(blsData.groupSignature, {from: operator});

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      blsData.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );
  });

  it("should reimburse submitter when no callback was provided", async () => {
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: customer});

    const beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {from: operator});
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    assert.isTrue(
      submitterReimbursement.eq(entryVerificationFee), 
      "Unexpected submitter reimbursement"
    );
  })

  it("should produce entry and execute callback if provided", async () => {
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
    );
    await operatorContract.relayEntry(blsData.groupSignature, {from: operator});

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      blsData.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );

    let result = web3.utils.toBN(await callbackContract.lastEntry());
    assert.isTrue(
      result.eq(blsData.groupSignatureNumber), 
      "Unexpected entry value passed to the callback"
    );
  })

  it("should reimburse submitter and customer for executing callback consuming less gas", async () => {
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);  
    callbackGas = callbackGas + 100000 // use 100k more gas than needed
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas) 
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
    );

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    // use the same gas price as the gas price ceiling
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('60', 'Gwei'));
    await operatorContract.relayEntry(blsData.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
    });

    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    let gasSurplus = customerSurplus.div(relayEntryTxGasPrice)
    await assertTotalSpending(
      entryFeeEstimate, 
      submitterReimbursement, 
      customerSurplus
    );

    // The additional 100k gas margin was never used so it should be returned
    // to the customer as surplus. Since gasleft() in Soldity inserts additional
    // opcodes, it's not possible to evaluate an exact gas spent value for
    // callback execution. Hence, we assume a healthy margin instead of an exact
    // value.
    expect(gasSurplus).to.gt.BN("98500")
    expect(gasSurplus).to.lte.BN("100000")
  })

  // gas price ceiling > tx.gasprice
  it("should reimburse submitter and customer for executing callback with lower tx.gasprice", async () => {
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use lower gas price when submitting entry
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('2', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
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

  // gas price ceiling == tx.gasprice
  it("should reimburse submitter and customer for executing callback with expected tx.gasprice", async () => {
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use the same gas price as the gas price ceiling
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('60', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
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

  // gas price ceiling < tx.gasprice
  it("should reimburse submitter and customer for executing callback with higher tx.gasprice", async () => {
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use higher price than the gas price ceiling
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('70', 'Gwei'));
    // higher tx.gasprice should not be used for reimbursement - maximum gas
    // price is the one from the gas price ceiling
    let gasPriceForReimbursement = web3.utils.toBN(web3.utils.toWei('60', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
      from: operator, 
      gasPrice: relayEntryTxGasPrice
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
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
    );
  
    await operatorContract.relayEntry(blsData.groupSignature, {from: operator});
  
    assert.equal((await operatorContract.getPastEvents())[1].event,
      'GroupSelectionStarted', "Should start group selection"
    );

    assert.equal((await operatorContract.getPastEvents())[1].args['newEntry'].toString(),
      blsData.groupSignatureNumber.toString(), "Should start group selection with new entry"
    );
    
    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      blsData.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );

    let result = web3.utils.toBN(await callbackContract.lastEntry());
    assert.isTrue(
      result.eq(blsData.groupSignatureNumber),
      "Unexpected entry value passed to the callback"
    );
  });

  // gas price ceiling == tx.gasprice
  it("should trigger new group creation, execute callback, reimburse submitter " + 
     "and customer with expected tx.gasprice", async () => {
    await fundDkgPool();

    // Request relay entry with a callback
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );

    // use the same gas price as the gas price ceiling
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('60', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
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

  // gas price ceiling > tx.gasprice
  it("should trigger new group selection, execute callback, reimburse submitter " +
     "and customer with lower tx.gasprice", async () => {
    await fundDkgPool();

    // Request relay entry with a callback
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );
  
    // use lower gas price when submitting entry
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('2', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
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
 
  // gas price ceiling < tx.gasprice
  it("should trigger new group selection, execute callback, reimburse submitter " +
    "and customer with higher tx.gasprice", async () => {
    await fundDkgPool();
  
    // Request relay entry with a callback
    let callbackGas = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {
        value: entryFeeEstimate,
        from: customer
      }
    );
  
    // use higher price than the gas price ceiling
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('70', 'Gwei'));
    // higher tx.gasprice should not be used for reimbursement - maximum gas
    // price is the one from the gas price ceiling
    let gasPriceForReimbursement = web3.utils.toBN(web3.utils.toWei('60', 'Gwei'));

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
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
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
    );
    
    await operatorContract.relayEntry(blsData.groupSignature, {from: operator});

    assert.equal((await operatorContract.getPastEvents())[1].event,
      'GroupSelectionStarted', "Should start group selection"
    );

    assert.equal((await operatorContract.getPastEvents())[1].args['newEntry'].toString(),
      blsData.groupSignatureNumber.toString(), "Should start group selection with new entry"
    );
  
    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      blsData.groupSignatureNumber.toString(), "Should emit event with the generated entry"
    );
  })

  it("should reimburse submitter when callback failed", async () => {
    let estimate = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    
    let callbackGas = estimate - 10; // Requestor provides wrong gas
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
    );
  
    // use the same gas price as the gas price ceiling
    let relayEntryTxGasPrice = web3.utils.toBN(web3.utils.toWei('60', 'Gwei'));
  
    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
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

  it("should return surplus to requestor when callback fails", async () => {
    let lastEntry = await callbackContract.lastEntry();
    let estimate = await callbackContract.__beaconCallback.estimateGas(blsData.groupSignature);
    let callbackGas = estimate - 25000; // Requestor provides wrong gas
    let entryFeeEstimate = await serviceContract.entryFeeEstimate(callbackGas)
    await serviceContract.methods['requestRelayEntry(address,uint256)'](
      callbackContract.address, callbackGas, {value: entryFeeEstimate, from: customer}
    );

    let customerStartBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryStartBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));
    await operatorContract.relayEntry(blsData.groupSignature, {
      from: operator
    });
    let customerEndBalance = web3.utils.toBN(await web3.eth.getBalance(customer));
    let beneficiaryEndBalance = web3.utils.toBN(await web3.eth.getBalance(beneficiary));

    let customerSurplus = customerEndBalance.sub(customerStartBalance);
    let submitterReimbursement = beneficiaryEndBalance.sub(beneficiaryStartBalance);

    await assertCallbackReimbursement(
      callbackGas, entryFeeEstimate, submitterReimbursement,
      customerSurplus, web3.utils.toBN(await web3.eth.getGasPrice())
    );

    const dkgSubmitterReimbursementFee = await operatorContract.dkgSubmitterReimbursementFee();
    const operatorContractBalance = await web3.eth.getBalance(operatorContract.address);

    assert.isTrue(
      web3.utils.toBN(operatorContractBalance).eq(
        dkgSubmitterReimbursementFee.add(groupProfitFee)
      ),
      "Unexpected operator contract balance"
    );

    assert.isTrue(
      lastEntry.eq(await callbackContract.lastEntry()),
      "Unexpected callback"
    );
  });

  // This function assets expected submitter reimbursement and customer surplus
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
    assertTotalSpending(entryFeeEstimate, submitterReimbursement, customerSurplus)

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

  // This function makes sure the total spending do not exceed the 
  // fee provided by the customer.
  // It makes sure beneficiary account balance change (submitter 
  // reward + callback  reimbursement) and customer account balance
  // change (surplus) add up to entry verification fee and callback fee as
  // calculated by the service contract. In other words, beacon does not spend
  // more than received and all required contributions (DKG fee + group profit 
  // fee) stay in the contract.
  async function assertTotalSpending(
    entryFeeEstimate, submitterReimbursement, customerSurplus
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
    let gasPriceCeiling = await operatorContract.gasPriceCeiling()
    let groupCreationFee = gasPriceCeiling.mul(groupSelectionGasEstimate)

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
    const groupCreationFee = await operatorContract.groupCreationFee();  
    await serviceContract.fundDkgFeePool({value: groupCreationFee});
  }
});
