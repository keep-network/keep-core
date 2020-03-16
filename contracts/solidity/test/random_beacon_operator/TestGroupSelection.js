import expectThrow from '../helpers/expectThrow';
import mineBlocks from '../helpers/mineBlocks';
import packTicket from '../helpers/packTicket';
import generateTickets from '../helpers/generateTickets';
import stakeDelegate from '../helpers/stakeDelegate';
import {initContracts} from '../helpers/initContracts';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage';
import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';

contract('KeepRandomBeaconOperator/GroupSelection', function(accounts) {
  let operatorContract,
  owner = accounts[0], 
  magpie = accounts[1],
  operator1 = accounts[2], tickets1,
  operator2 = accounts[3], tickets2,
  operator3 = accounts[4], tickets3,
  authorizer = owner;

  const minimumStake = web3.utils.toBN(200000);
  const operator1StakingWeight = 2000;
  const operator2StakingWeight = 2000;
  const operator3StakingWeight = 3000;
  const groupSize = 3;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorGroupSelectionStub.sol')
    );
    
    let token = contracts.token;
    let stakingContract = contracts.stakingContract;

    operatorContract = contracts.operatorContract;

    await stakingContract.setMinimumStake(minimumStake, owner)
    await operatorContract.setGroupSize(groupSize)

    await stakeDelegate(
      stakingContract, token, owner, operator1, magpie, authorizer,
      minimumStake.mul(web3.utils.toBN(operator1StakingWeight))
    );
    await stakeDelegate(
      stakingContract, token, owner, operator2, magpie, authorizer,
      minimumStake.mul(web3.utils.toBN(operator2StakingWeight))
    );
    await stakeDelegate(
      stakingContract, token, owner, operator3, magpie, authorizer,
      minimumStake.mul(web3.utils.toBN(operator3StakingWeight))
    );

    await stakingContract.authorizeOperatorContract(operator1, operatorContract.address, {from: authorizer})
    await stakingContract.authorizeOperatorContract(operator2, operatorContract.address, {from: authorizer})
    await stakingContract.authorizeOperatorContract(operator3, operatorContract.address, {from: authorizer})

    tickets1 = generateTickets(
      await operatorContract.getGroupSelectionRelayEntry(), 
      operator1, 
      operator1StakingWeight
    );
    tickets2 = generateTickets(
      await operatorContract.getGroupSelectionRelayEntry(), 
      operator2, 
      operator2StakingWeight
    );
    tickets3 = generateTickets(
      await operatorContract.getGroupSelectionRelayEntry(), 
      operator3, 
      operator3StakingWeight
    );
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should fail to get selected participants before submission period is over", async () => {
    await expectThrow(operatorContract.selectedParticipants());
  });

  it("should accept valid ticket with minimum virtual staker index", async () => {
    let ticket = packTicket(tickets1[0].valueHex, 1, operator1);
    await operatorContract.submitTicket(ticket, {from: operator1});

    let submittedCount = await operatorContract.submittedTicketsCount();
    assert.equal(1, submittedCount, "Ticket should be accepted");
  });

  it("should accept valid ticket with maximum virtual staker index", async () => {
    let ticket = packTicket(tickets1[tickets1.length - 1].valueHex, tickets1.length, operator1);
    await operatorContract.submitTicket(ticket,{from: operator1});

    let submittedCount = await operatorContract.submittedTicketsCount();
    assert.equal(1, submittedCount, "Ticket should be accepted");
  });

  it("should reject ticket with too high virtual staker index", async () => {
    let ticket = packTicket(tickets1[tickets1.length - 1].valueHex, tickets1.length + 1, operator1);
    await expectThrowWithMessage(
      operatorContract.submitTicket(ticket, {from: operator1}),
      "Invalid ticket"
    );
  });

  it("should reject ticket with invalid value", async() => {
    let ticket = packTicket('0x1337', 1, operator1);
    await expectThrowWithMessage(
      operatorContract.submitTicket(ticket, {from: operator1}),
      "Invalid ticket"
    );
  });

  it("should reject ticket with not matching operator", async() => {
    let ticket = packTicket(tickets1[0].valueHex, 1, operator1);
    await expectThrowWithMessage(
      operatorContract.submitTicket(ticket, {from: operator2}),
      "Invalid ticket"
    )
  });

  it("should reject ticket with not matching virtual staker index", async() => {
    let ticket = packTicket(tickets1[0].valueHex, 2, operator1);
    await expectThrowWithMessage(
      operatorContract.submitTicket(ticket, {from: operator1}),
      "Invalid ticket"
    )
  });

  it("should reject duplicate ticket", async () => {
    let ticket = packTicket(tickets1[0].valueHex, 1, operator1);
    await operatorContract.submitTicket(ticket, {from: operator1});

    await expectThrowWithMessage(
      operatorContract.submitTicket(ticket, {from: operator1}),
      "Duplicate ticket"
    );
  });

  it("should trim selected participants to the group size", async () => {
    let groupSize = await operatorContract.groupSize();
    let ticket;
  
    for (let i = 0; i < groupSize*2; i++) {
      ticket = packTicket(tickets1[i].valueHex, tickets1[i].virtualStakerIndex, operator1);
      await operatorContract.submitTicket(ticket, {from: operator1});
    }
  
    mineBlocks(await operatorContract.ticketSubmissionTimeout());

    let selectedParticipants = await operatorContract.selectedParticipants();
    assert.equal(
      selectedParticipants.length, 
      groupSize, 
      "Selected participants list should be trimmed to groupSize length"
    );
  });

  it("should select participants by tickets in ascending order", async function() {
    let tickets = [
      {value: tickets1[0].value, operator: operator1},
      {value: tickets2[0].value, operator: operator2},
      {value: tickets3[0].value, operator: operator3}
    ];

    // Sort tickets in ascending order
    tickets = tickets.sort(function(a, b){return a.value-b.value});

    let ticket1 = packTicket(tickets1[0].valueHex, tickets1[0].virtualStakerIndex, operator1);
    await operatorContract.submitTicket(ticket1, {from: operator1});

    let ticket2 = packTicket(tickets2[0].valueHex, tickets2[0].virtualStakerIndex, operator2);
    await operatorContract.submitTicket(ticket2, {from: operator2});

    let ticket3 = packTicket(tickets3[0].valueHex, tickets3[0].virtualStakerIndex, operator3);
    await operatorContract.submitTicket(ticket3, {from: operator3});

    mineBlocks(await operatorContract.ticketSubmissionTimeout());

    let selectedParticipants = await operatorContract.selectedParticipants();
    assert.equal(
      selectedParticipants[0], 
      tickets[0].operator, 
      "Unexpected operator selected at position 0"
    );
    assert.equal(
      selectedParticipants[1], 
      tickets[1].operator, 
      "Unexpected operator selected at position 1"
    );
    assert.equal(
      selectedParticipants[2], 
      tickets[2].operator, 
      "Unexpected operator selected at position 2"
    );
  });

  it("should properly override previous group selection data", async function() {
    // Simulate previous data existence: operator 2 has submitted two tickets and operator 3 has submitted one ticket
    await operatorContract.submitTicket(
        packTicket(tickets2[10].valueHex, tickets2[10].virtualStakerIndex, operator2),
        {from: operator2}
    );
    await operatorContract.submitTicket(
        packTicket(tickets3[10].valueHex, tickets3[10].virtualStakerIndex, operator3),
        {from: operator3}
    );
    await operatorContract.submitTicket(
        packTicket(tickets2[11].valueHex, tickets2[11].virtualStakerIndex, operator2),
        {from: operator2}
    );

    mineBlocks(await operatorContract.ticketSubmissionTimeout());

    // Start new group selection
    const seed = await operatorContract.getGroupSelectionRelayEntry();
    await operatorContract.startGroupSelection(seed);

    let tickets = [
      {value: tickets1[0].value, operator: operator1},
      {value: tickets2[0].value, operator: operator2},
      {value: tickets3[0].value, operator: operator3}
    ];

    // Sort tickets in ascending order
    tickets = tickets.sort(function(a, b){return a.value-b.value});

    let ticket1 = packTicket(tickets1[0].valueHex, tickets1[0].virtualStakerIndex, operator1);
    await operatorContract.submitTicket(ticket1, {from: operator1});

    let ticket2 = packTicket(tickets2[0].valueHex, tickets2[0].virtualStakerIndex, operator2);
    await operatorContract.submitTicket(ticket2, {from: operator2});

    let ticket3 = packTicket(tickets3[0].valueHex, tickets3[0].virtualStakerIndex, operator3);
    await operatorContract.submitTicket(ticket3, {from: operator3});

    mineBlocks(await operatorContract.ticketSubmissionTimeout());

    let selectedParticipants = await operatorContract.selectedParticipants();
    assert.equal(
        selectedParticipants.length,
        3,
        "Unexpected number of selected participants"
    );
    assert.equal(
        selectedParticipants[0],
        tickets[0].operator,
        "Unexpected operator selected at position 0"
    );
    assert.equal(
        selectedParticipants[1],
        tickets[1].operator,
        "Unexpected operator selected at position 1"
    );
    assert.equal(
        selectedParticipants[2],
        tickets[2].operator,
        "Unexpected operator selected at position 2"
    );
  });
});
