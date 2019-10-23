import expectThrow from './helpers/expectThrow';
import mineBlocks from './helpers/mineBlocks';
import generateTickets from './helpers/generateTickets';
import stakeDelegate from './helpers/stakeDelegate';
import {initContracts} from './helpers/initContracts';
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";

contract('KeepRandomBeaconOperator', function(accounts) {
  let token, stakingContract, operatorContract,
  owner = accounts[0], magpie = accounts[1],
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
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
    );
    
    token = contracts.token;
    operatorContract = contracts.operatorContract;
    stakingContract = contracts.stakingContract;

    operatorContract.setMinimumStake(minimumStake)

    await stakeDelegate(stakingContract, token, owner, operator1, magpie, minimumStake.mul(web3.utils.toBN(operator1StakingWeight)))
    await stakeDelegate(stakingContract, token, owner, operator2, magpie, minimumStake.mul(web3.utils.toBN(operator2StakingWeight)))
    await stakeDelegate(stakingContract, token, owner, operator3, magpie, minimumStake.mul(web3.utils.toBN(operator3StakingWeight)))

    tickets1 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator1, operator1StakingWeight);
    tickets2 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator2, operator2StakingWeight);
    tickets3 = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), operator3, operator3StakingWeight);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should be able to get staking weight", async function() {
    assert.isTrue(web3.utils.toBN(2000).eq(await operatorContract.stakingWeight(operator1)), "Should have expected staking weight.");
    assert.isTrue(web3.utils.toBN(3000).eq(await operatorContract.stakingWeight(operator3)), "Should have expected staking weight.");
  });

  it("should fail to get selected tickets before submission period is over", async function() {
    await expectThrow(operatorContract.selectedTickets());
  });

  it("should fail to get selected participants before submission period is over", async function() {
    await expectThrow(operatorContract.selectedParticipants());
  });

  it("should be able to output submitted tickets in ascending ordered", async function() {
    let tickets = [];

    await operatorContract.submitTicket(tickets1[0].value, operator1, tickets1[0].virtualStakerIndex, {from: operator1});
    tickets.push(tickets1[0].value);

    await operatorContract.submitTicket(tickets2[0].value, operator2, tickets2[0].virtualStakerIndex, {from: operator2});
    tickets.push(tickets2[0].value);

    await operatorContract.submitTicket(tickets3[0].value, operator3, tickets3[0].virtualStakerIndex, {from: operator3});
    tickets.push(tickets3[0].value);

    tickets = tickets.sort(function(a, b){return a-b}); // Sort numbers in ascending order

    // Test tickets ordering
    let orderedTickets = await operatorContract.orderedTickets();
    assert.isTrue(orderedTickets[0].eq(tickets[0]), "Tickets should be in ascending order.");
    assert.isTrue(orderedTickets[1].eq(tickets[1]), "Tickets should be in ascending order.");
    assert.isTrue(orderedTickets[2].eq(tickets[2]), "Tickets should be in ascending order.");
  });

  it("should be able to verify a ticket", async function() {
    await operatorContract.submitTicket(tickets1[0].value, operator1, 1, {from: operator1});

    assert.isTrue(await operatorContract.isTicketValid(
      operator1, tickets1[0].value, operator1, tickets1[0].virtualStakerIndex
    ), "Should be able to verify a valid ticket.");

    let lastTicketIndex = tickets1.length - 1;
    let maxVirtualStakerIndexTicket = tickets1[lastTicketIndex].virtualStakerIndex;
    assert.isTrue(await operatorContract.isTicketValid(
      operator1, tickets1[lastTicketIndex].value, operator1, maxVirtualStakerIndexTicket
    ), "Should be able to verify a valid ticket with the maximum allowed staker index");

    let invalidVirtualStakerIndex = operator1StakingWeight + 1;
    assert.isFalse(await operatorContract.isTicketValid(
      operator1, tickets1[0].value, operator1, invalidVirtualStakerIndex
    ), "Should fail while verifying a submitted ticket due to invalid number of virtual stakers");
    
    assert.isFalse(await operatorContract.isTicketValid(
      operator1, 0, operator2, tickets1[0].virtualStakerIndex
    ), "Should fail while verifying a submitted ticket due to invalid ticket value");
    
    assert.isFalse(await operatorContract.isTicketValid(
      operator1, tickets1[0].value, operator2, tickets1[0].virtualStakerIndex
    ), "Should fail while verifying a submitted ticket due to invalid stake value");
      
    assert.isFalse(await operatorContract.isTicketValid(
      operator1, tickets1[0].value, operator1, 2
    ), "Should fail while verifying a submitted ticket due to invalid virtual staker index");
  });

  it("should revert the transaction when the ticket has been already submitted", async function() {
    await operatorContract.submitTicket(tickets1[0].value, operator1, 1, {from: operator1});

    await expectThrowWithMessage(
      operatorContract.submitTicket(tickets1[0].value, operator1, 1, {from: operator1}),
      "Ticket with the given value has already been submitted."
    );
  })

  it("should be able to get selected tickets and participants after submission period is over", async function() {
    let groupSize = await operatorContract.groupSize();
  
    for (let i = 0; i < groupSize*2; i++) {
      await operatorContract.submitTicket(tickets1[i].value, operator1, tickets1[i].virtualStakerIndex, {from: operator1});
    }
  
    mineBlocks(await operatorContract.ticketReactiveSubmissionTimeout());
    let selectedTickets = await operatorContract.selectedTickets();
    assert.equal(selectedTickets.length, groupSize, "Should be trimmed to groupSize length.");
  
    let selectedParticipants = await operatorContract.selectedParticipants();
    assert.equal(selectedParticipants.length, groupSize, "Should be trimmed to groupSize length.");
  });
});
