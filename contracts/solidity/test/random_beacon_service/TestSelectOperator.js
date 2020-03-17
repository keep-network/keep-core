import expectThrow from '../helpers/expectThrow';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage';
import {initContracts} from '../helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "../helpers/snapshot";
const OperatorContract = artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')

contract('TestKeepRandomBeaconService/SelectOperator', function() {

  let registry, stakingContract, serviceContract, operatorContract, operatorContract2, operatorContract3;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      OperatorContract
    );

    registry = contracts.registry;
    stakingContract = contracts.stakingContract;
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;

    // Create and initialize additional operator contracts
    operatorContract2 = await OperatorContract.new(serviceContract.address, stakingContract.address);
    operatorContract3 = await OperatorContract.new(serviceContract.address, stakingContract.address);

    operatorContract.registerNewGroup("0x0");
    operatorContract2.registerNewGroup("0x0");
    operatorContract2.registerNewGroup("0x0");
    operatorContract3.registerNewGroup("0x0");
    operatorContract3.registerNewGroup("0x0");
    operatorContract3.registerNewGroup("0x0");
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("service contract owner should be able to remove and add operator contracts.", async function() {
    let result = await serviceContract.selectOperatorContract(0);
    assert.equal(result, operatorContract.address, "Operator contract added during initialization should present in the service contract.");

    await serviceContract.removeOperatorContract(operatorContract.address);
    await expectThrow(serviceContract.selectOperatorContract(0)); // Should revert since no operator contract present.

    await registry.approveOperatorContract(operatorContract2.address);
    await serviceContract.addOperatorContract(operatorContract2.address);
    result = await serviceContract.selectOperatorContract(0);
    assert.equal(result, operatorContract2.address, "Operator contract should be added");

  });

  it("should select contract from operators list according to the amount of groups.", async function() {
    await registry.approveOperatorContract(operatorContract2.address);
    await registry.approveOperatorContract(operatorContract3.address);
    serviceContract.addOperatorContract(operatorContract2.address);
    serviceContract.addOperatorContract(operatorContract3.address);

    let selectionCounter = {};
    selectionCounter[operatorContract.address] = 0;
    selectionCounter[operatorContract2.address] = 0;
    selectionCounter[operatorContract3.address] = 0;

    // Total max weight = 6 (Operator1 - 1 group, Operator2 - 2 groups, Operator3 - 3 groups)
    for(let i = 0; i < 6; i++) {
      let address = await serviceContract.selectOperatorContract(i);
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

    await registry.disableOperatorContract(operatorContract.address);
    await registry.disableOperatorContract(operatorContract2.address);
    await registry.disableOperatorContract(operatorContract3.address);

    await expectThrowWithMessage(
      serviceContract.selectOperatorContract(0),
      "Total number of groups must be greater than zero."
    );
  });

});
