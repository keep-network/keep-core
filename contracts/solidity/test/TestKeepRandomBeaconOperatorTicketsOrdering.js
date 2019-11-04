import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";

contract('KeepRandomBeaconOperator', function(accounts) {
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
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("ticket insertion", function() {
    let groupSize;

    beforeEach(async () => {
      groupSize = await operatorContract.groupSize() // should be 10
    });

    describe("tickets array size is at its max capacity", function() {
      it("should reject a new ticket when it is higher than an existing highest one.", async () =>{
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 6, 11, 8, 12, 100, 200, 300];

        await addTickets(ticketsToAdd)

        let expectedTail = 9;
        // smallest index points to itself
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 8, 2, 5, 3, 7];
        // await logTicketStatus(ticketsToAdd, expectedOrderedIndices)
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)

      });

      it("should replace the highest existing with new ticket which is somewhere in the middle value range", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 135]; // 190 -> out

        await addTickets(ticketsToAdd)

        let expectedTail = 2;
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 9, 1, 5, 6, 4];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

      it("should replace highest ticket (tail) and become a new highest one (also tail)", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 185]; // 190 -> out

        await addTickets(ticketsToAdd)

        let expectedTail = 5;
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 2, 1, 9, 6, 4];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

      it("should add a new smallest ticket and remove the highest", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 2]; // 190 -> out

        await addTickets(ticketsToAdd)

        let expectedTail = 5;
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 2, 1, 9, 6, 4];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });
    });


    describe("tickets array size is less than a group size", function() {
      it("should add all the tickets and keep track the order", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 6, 11];

        await addTickets(ticketsToAdd)

        let expectedTail = 7;
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2, 5];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

      it("should add all the tickets and track the order when a latest ticket is between smallest and biggest", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 11, 6];

        await addTickets(ticketsToAdd)

        let expectedTail = 6;
        let expectedOrderedIndices = [0, 0, 4, 7, 1, 3, 5, 2];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)

      });

      it("should add all the tickets and track the order when a last added ticket is the smallest", async () => {
        let ticketsToAdd = [2, 3, 5, 7, 4, 9, 11, 1];

        await addTickets(ticketsToAdd)

        let expectedTail = 6;
        let expectedOrderedIndices = [7, 0, 4, 2, 1, 3, 5, 7];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

    });

    async function addTickets(ticketsToAdd) {
      for (let i = 0; i < ticketsToAdd.length; i++) {
        await operatorContract.addTicket(ticketsToAdd[i]);
      }
    }

    async function assertTicketsProperties(expectedTail, expectedOrderedIndices) {
      let tickets = await operatorContract.getTickets();
      let correctNumberOfTickets = tickets.length < groupSize ? tickets.length : groupSize;
      assert.equal(tickets.length, correctNumberOfTickets, "array of tickets should be the size of: " + correctNumberOfTickets)

      let tail = await operatorContract.getTail()
      assert.equal(expectedTail, tail.toString(), "tail index should be equal to " + expectedTail)

      for (let i = 0; i < tickets.length; i++) {
        let prevByIndex = await operatorContract.getOrderedLinkedTicketIndex(i)
        assert.equal(expectedOrderedIndices[i] + '', prevByIndex.toString())
      }
    }

    // For debugging purposes.
    async function logTicketStatus(ticketsToAdd, expectedOrderedIndices) {
      console.log("--------------")
      console.log("added tickets: [" + ticketsToAdd.toString() + "]")

      console.log("--------------")
      console.log("number added tickets: ", ticketsToAdd.length)

      console.log("--------------")
      let tail = await operatorContract.getTail();
      console.log("tail index: ", tail.toString());

      console.log("--------------")
      let ticketMaxValue = await operatorContract.getTicketMaxValue();
      console.log("max value ticket[tail]: ", ticketMaxValue.toString());

      console.log("--------------")
      let tickets = await operatorContract.getTickets();
      console.log("tickets on chain[" + tickets.toString() + "]")

      console.log("--------------")
      for (let i = 0; i < tickets.length; i++) {
        console.log("prev_tickets[" + i + "] -> " + await operatorContract.getOrderedLinkedTicketIndex(i) + " | " + expectedOrderedIndices[i]);
      }
    }
  })

});
