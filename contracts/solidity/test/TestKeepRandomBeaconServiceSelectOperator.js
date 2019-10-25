import expectThrow from './helpers/expectThrow';
import {initContracts} from './helpers/initContracts';
import {bls} from './helpers/data';
const OperatorContract = artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
const GroupContract = artifacts.require('./stubs/KeepRandomBeaconOperatorGroups.sol')
const TicketContract = artifacts.require('./stubs/KeepRandomBeaconOperatorTickets.sol')

contract('TestKeepRandomBeaconServiceSelectOperator', function() {

  let stakingContract, serviceContract, operatorContract, ticketContract,
    groupContract2, operatorContract2,
    groupContract3, operatorContract3;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      OperatorContract,
      GroupContract,
      TicketContract
    );

    stakingContract = contracts.stakingContract;
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;
    ticketContract = contracts.ticketContract;

    // Create and initialize additional operator contracts
    groupContract2 = await GroupContract.new();
    operatorContract2 = await OperatorContract.new(serviceContract.address, stakingContract.address, groupContract2.address, ticketContract.address);
    await groupContract2.setOperatorContract(operatorContract2.address);
    groupContract3 = await GroupContract.new();
    operatorContract3 = await OperatorContract.new(serviceContract.address, stakingContract.address, groupContract3.address, ticketContract.address);
    await groupContract3.setOperatorContract(operatorContract3.address);

    operatorContract.registerNewGroup("0x0");
    operatorContract2.registerNewGroup("0x0");
    operatorContract2.registerNewGroup("0x0");
    operatorContract3.registerNewGroup("0x0");
    operatorContract3.registerNewGroup("0x0");
    operatorContract3.registerNewGroup("0x0");
  });

  it("service contract owner should be able to remove and add operator contracts.", async function() {
    let result = await serviceContract.selectOperatorContract(0);
    assert.equal(result, operatorContract.address, "Operator contract added during initialization should present in the service contract.");

    await serviceContract.removeOperatorContract(operatorContract.address);
    await expectThrow(serviceContract.selectOperatorContract(0)); // Should revert since no operator contract present.

    await serviceContract.addOperatorContract(operatorContract.address);
    result = await serviceContract.selectOperatorContract(0);
    assert.equal(result, operatorContract.address, "Operator contract should be added");

  });

  it("should select contract from operators list according to the amount of groups.", async function() {
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

  });

});
