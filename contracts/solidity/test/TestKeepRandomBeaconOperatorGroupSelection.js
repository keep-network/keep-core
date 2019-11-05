import expectThrow from './helpers/expectThrow';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import stakeDelegate from './helpers/stakeDelegate';
import {initContracts} from './helpers/initContracts';
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";

contract('KeepRandomBeaconOperator', function(accounts) {
  let operatorContract,
  owner = accounts[0], 
  magpie = accounts[1],
  operator1 = accounts[2], tickets1,
  operator2 = accounts[3], tickets2,
  operator3 = accounts[4], tickets3;

  const minimumStake = web3.utils.toBN(200000);
  const operator1StakingWeight = 2000;
  const operator2StakingWeight = 2000;
  const operator3StakingWeight = 3000;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorGroupSelectionStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
    );
    
    let token = contracts.token;
    let stakingContract = contracts.stakingContract;

    operatorContract = contracts.operatorContract;

    await operatorContract.setMinimumStake(minimumStake)

    await stakeDelegate(
      stakingContract, token, owner, operator1, magpie, 
      minimumStake.mul(web3.utils.toBN(operator1StakingWeight))
    );
    await stakeDelegate(
      stakingContract, token, owner, operator2, magpie, 
      minimumStake.mul(web3.utils.toBN(operator2StakingWeight))
    );
    await stakeDelegate(
      stakingContract, token, owner, operator3, magpie, 
      minimumStake.mul(web3.utils.toBN(operator3StakingWeight))
    );

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
    await operatorContract.submitTicket(
      tickets1[0].value, 
      operator1, 
      1, 
      {from: operator1}
    );

    let submittedCount = await operatorContract.submittedTicketsCount();
    assert.equal(1, submittedCount, "Ticket should be accepted");
  });

  it("should accept valid ticket with maximum virtual staker index", async () => {
    await operatorContract.submitTicket(
      tickets1[tickets1.length - 1].value,
      operator1,
      tickets1.length,
      {from: operator1}
    );

    let submittedCount = await operatorContract.submittedTicketsCount();
    assert.equal(1, submittedCount, "Ticket should be accepted");
  });

  it("should reject ticket with too high virtual staker index", async () => {
    await expectThrowWithMessage(
      operatorContract.submitTicket(
        tickets1[tickets1.length - 1].value,
        operator1,
        tickets1.length + 1,
        {from: operator1}
      ),
      "Invalid ticket"
    );
  });

  it("should reject ticket with invalid value", async() => {
    await expectThrowWithMessage(
      operatorContract.submitTicket(
        1337,
        operator1,
        1,
        {from: operator1}
      ),
      "Invalid ticket"
    );
  });

  it("should reject ticket with not matching operator", async() => {
    await expectThrowWithMessage(
      operatorContract.submitTicket(
        tickets1[0].value, 
        operator1, 
        1, 
        {from: operator2}
      ),
      "Invalid ticket"
    )
  });

  it("should reject ticket with not matching virtual staker index", async() => {
    await expectThrowWithMessage(
      operatorContract.submitTicket(
        tickets1[0].value, 
        operator1, 
        2, 
        {from: operator1}
      ),
      "Invalid ticket"
    )
  });

  it("should reject duplicate ticket", async () => {
    await operatorContract.submitTicket(
      tickets1[0].value, 
      operator1, 
      1, 
      {from: operator1}
    );

    await expectThrowWithMessage(
      operatorContract.submitTicket(
        tickets1[0].value, 
        operator1, 
        1, 
        {from: operator1}
      ),
      "Duplicate ticket"
    );
  })

  it("should trim selected participants to the group size", async () => {
    let groupSize = await operatorContract.groupSize();
  
    for (let i = 0; i < groupSize*2; i++) {
      await operatorContract.submitTicket(
        tickets1[i].value, 
        operator1, 
        tickets1[i].virtualStakerIndex, 
        {from: operator1}
      );
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

    await operatorContract.submitTicket(
      tickets1[0].value, 
      operator1, 
      tickets1[0].virtualStakerIndex, 
      {from: operator1}
    );
    await operatorContract.submitTicket(
      tickets2[0].value, 
      operator2, 
      tickets2[0].virtualStakerIndex, 
      {from: operator2}
    );
    await operatorContract.submitTicket(
      tickets3[0].value, 
      operator3, 
      tickets3[0].virtualStakerIndex, 
      {from: operator3}
    );

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
});
