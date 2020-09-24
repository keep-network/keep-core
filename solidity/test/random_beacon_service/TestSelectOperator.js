const {expectRevert} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {initContracts} = require('../helpers/initContracts')
const assert = require('chai').assert
const {contract, accounts} = require("@openzeppelin/test-environment")
const OperatorContract = contract.fromArtifact('KeepRandomBeaconOperatorStub')
const GasPriceOracle = contract.fromArtifact("GasPriceOracle")

describe('TestKeepRandomBeaconService/SelectOperator', function() {

  let registry, stakingContract, serviceContract, operatorContract, operatorContract2, operatorContract3;

  before(async () => {
    const gasPriceOracle = await GasPriceOracle.new({from: accounts[0]})
    const contracts = await initContracts(
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      OperatorContract
    );

    registry = contracts.registry;
    stakingContract = contracts.stakingContract;
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;

    // Create and initialize additional operator contracts
    operatorContract2 = await OperatorContract.new(
      serviceContract.address, 
      stakingContract.address, 
      registry.address,
      gasPriceOracle.address,
      {from: accounts[0]}
    );
    operatorContract3 = await OperatorContract.new(
      serviceContract.address, 
      stakingContract.address, 
      registry.address,
      gasPriceOracle.address,
      {from: accounts[0]}
    );

    await operatorContract.registerNewGroup("0x0", {from: accounts[0]});
    await operatorContract2.registerNewGroup("0x0", {from: accounts[0]});
    await operatorContract2.registerNewGroup("0x0", {from: accounts[0]});
    await operatorContract3.registerNewGroup("0x0", {from: accounts[0]});
    await operatorContract3.registerNewGroup("0x0", {from: accounts[0]});
    await operatorContract3.registerNewGroup("0x0", {from: accounts[0]});
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("service contract owner should be able to remove and add operator contracts.", async function() {
    let result = await serviceContract.selectOperatorContract(0, {from: accounts[0]});
    assert.equal(result, operatorContract.address, "Operator contract added during initialization should present in the service contract.");

    await serviceContract.removeOperatorContract(operatorContract.address, {from: accounts[0]});
    await expectRevert(
      serviceContract.selectOperatorContract(0, {from: accounts[0]}),
      "Total number of groups must be greater than zero"
    ); // Should revert since no operator contract present.

    await registry.approveOperatorContract(operatorContract2.address, {from: accounts[0]});
    await serviceContract.addOperatorContract(operatorContract2.address, {from: accounts[0]});
    result = await serviceContract.selectOperatorContract(0, {from: accounts[0]});
    assert.equal(result, operatorContract2.address, "Operator contract should be added");

  });

  it("should select contract from operators list according to the amount of groups.", async function() {
    await registry.approveOperatorContract(operatorContract2.address, {from: accounts[0]});
    await registry.approveOperatorContract(operatorContract3.address, {from: accounts[0]});
    serviceContract.addOperatorContract(operatorContract2.address, {from: accounts[0]});
    serviceContract.addOperatorContract(operatorContract3.address, {from: accounts[0]});

    let selectionCounter = {};
    selectionCounter[operatorContract.address] = 0;
    selectionCounter[operatorContract2.address] = 0;
    selectionCounter[operatorContract3.address] = 0;

    // Total max weight = 6 (Operator1 - 1 group, Operator2 - 2 groups, Operator3 - 3 groups)
    for(let i = 0; i < 6; i++) {
      let address = await serviceContract.selectOperatorContract(i, {from: accounts[0]});
      selectionCounter[address] = selectionCounter[address] + 1;
    }

    assert.equal(
      selectionCounter[operatorContract.address],
      (await operatorContract.numberOfGroups()).toNumber(), "Contract selection counter should be equal to the number of groups."
    );

    assert.equal(
      selectionCounter[operatorContract2.address],
      (await operatorContract2.numberOfGroups()).toNumber(), "Contract selection counter should be equal to the number of groups."
    );

    assert.equal(
      selectionCounter[operatorContract3.address],
      (await operatorContract3.numberOfGroups()).toNumber(), "Contract selection counter should be equal to the number of groups."
    );

    await registry.disableOperatorContract(operatorContract.address, {from: accounts[0]});
    await registry.disableOperatorContract(operatorContract2.address, {from: accounts[0]});
    await registry.disableOperatorContract(operatorContract3.address, {from: accounts[0]});

    await expectRevert(
      serviceContract.selectOperatorContract(0, {from: accounts[0]}),
      "Total number of groups must be greater than zero."
    );
  });
});
