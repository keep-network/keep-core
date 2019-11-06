import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";

contract('KeepRandomBeaconOperator', function() {
  const groupSize = 10;

  let operatorContract;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorTicketsOrderingStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
    );

    operatorContract = contracts.operatorContract;
    operatorContract.setGroupSize(groupSize);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("ticket insertion", function() {

    describe("tickets array size is at its max capacity", function() {

      it("should reject a new ticket when it is higher than the current highest one", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 6, 11, 8, 12, 100, 200, 300];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [1, 3, 5, 7, 4, 9, 6, 11, 8, 12]; // 100, 200, 300 -> out
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 8, 2, 5, 3, 7];
        let expectedTail = 9;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

      it("should replace the highest current with a new ticket which is somewhere in the middle value range", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 135];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [151, 42, 175, 7, 128, 135, 74, 143, 88, 130]; // 190 -> out
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 9, 1, 5, 6, 4];
        let expectedTail = 2;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

      it("should replace highest ticket (tail) and become a new highest one (also tail)", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 185];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [151, 42, 175, 7, 128, 185, 74, 143, 88, 130]; // 190 -> out
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 2, 1, 9, 6, 4];
        let expectedTail = 5;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

      it("should add a new smallest ticket and remove the highest", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 2];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [151, 42, 175, 7, 128, 2, 74, 143, 88, 130]; // 190 -> out
        let expectedOrderedIndices = [7, 3, 0, 5, 8, 5, 1, 9, 6, 4];
        let expectedTail = 2;
    
        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

    });

    describe("tickets array size is less than a group size", function() {

      it("should add all the tickets and keep track the order", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 6, 11];

        await addTickets(ticketsToAdd)

        let expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2, 5];
        let expectedTail = 7;

        await assertTickets(expectedTail, expectedOrderedIndices, ticketsToAdd)
      });

      it("should add all the tickets and track the order when a latest ticket is between smallest and biggest", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 11, 6];

        await addTickets(ticketsToAdd)

        let expectedOrderedIndices = [0, 0, 4, 7, 1, 3, 5, 2];
        let expectedTail = 6;

        await assertTickets(expectedTail, expectedOrderedIndices, ticketsToAdd)
      });

      it("should add all the tickets and track the order when a last added ticket is the smallest", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 4];

        await addTickets(ticketsToAdd)

        let expectedOrderedIndices = [4, 3, 0, 7, 6, 2, 1, 7];
        let expectedTail = 5;

        await assertTickets(expectedTail, expectedOrderedIndices, ticketsToAdd);
      });

    });

    async function addTickets(ticketsToAdd) {
      for (let i = 0; i < ticketsToAdd.length; i++) {
        await operatorContract.addTicket(ticketsToAdd[i]);
      }
    };

    async function assertTickets(expectedTail, expectedLinkedTicketIndices, expectedTickets) {
      // Assert tickets size
      let tickets = await operatorContract.getTickets();
      assert.isAtMost(
        tickets.length,
        groupSize, 
        "array of tickets cannot have more elements than the group size"
      );

      // Assert ticket values
      let actualTickets = [];
      for (let i = 0; i < tickets.length; i++) {
        actualTickets.push(Number(tickets[i]))
      }
      assert.sameOrderedMembers(actualTickets, expectedTickets, "array values should be the same")

      // Assert tail
      let tail = await operatorContract.getTail()
      assert.equal(expectedTail, tail.toString(), "unexpected tail index")

      // Assert the order of the tickets[] indices
      let actualLinkedTicketIndices = [];
      for (let i = 0; i < tickets.length; i++) {
        let actualIndex = await operatorContract.getPreviousTicketIndex(i)
        actualLinkedTicketIndices.push(Number(actualIndex))
      }
      assert.sameOrderedMembers(
        actualLinkedTicketIndices,
        expectedLinkedTicketIndices,
        'array values (indices) should be the same'
      );
    };

  });

});
