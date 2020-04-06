const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const assert = require('chai').assert
const {contract} = require("@openzeppelin/test-environment")
const GroupSelectionStub = contract.fromArtifact('GroupSelectionStub');

describe('KeepRandomBeaconOperator/TicketsOrdering', function() {
  const groupSize = 10;
  
  let groupSelectionStub;

  before(async () => {
    groupSelectionStub = await GroupSelectionStub.new(groupSize);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("ticket insertion", () => {

    describe("tickets array size is at its max capacity", () => {

      it("should reject a new ticket when it is higher than the current highest one", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 6, 11, 8, 12, 100, 200, 300];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [1, 3, 5, 7, 4, 9, 6, 11, 8, 12]; // 100, 200, 300 -> out
        // indices          = [0, 1, 2, 3, 4, 5, 6, 7,  8,  9]
        // sorted tickets   = [1, 3, 4, 5, 6, 7, 8, 9, 11, 12]
        // sorted indices   = [0, 1, 4, 2, 6, 3, 8, 5,  7,  9]
        // 0->0    1->0    2->4    3->6    4->1
        // 5->8    6->2    7->5    8->3    9->7
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 8, 2, 5, 3, 7];
        let expectedTail = 9;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

      it("should replace the highest ticket with a new ticket which is somewhere in the middle value range", async () => {
        let ticketsToAdd = [5986, 6782, 5161, 7009, 8086, 1035, 5294, 9826, 6475, 9520, 4293];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [5986, 6782, 5161, 7009, 8086, 1035, 5294, 4293, 6475, 9520]; // 9826 -> out
        // indices          = [ 0  ,  1  ,  2  ,  3  ,  4  ,  5  ,  6  ,  7  ,  8  ,  9  ]
        // sorted tickets   = [1035, 4293, 5161, 5294, 5986, 6475, 6782, 7009, 8086, 9520]
        // sorted indices   = [ 5  ,  7  ,  2  ,  6  ,  0  ,  8  ,  1  ,  3  ,  4  ,  9  ]
        // 0->6    1->8    2->7    3->1    4->3
        // 5->5    6->2    7->5    8->0    9->4
        let expectedOrderedIndices = [6, 8, 7, 1, 3, 5, 2, 5, 0, 4];
        let expectedTail = 9;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

      it("should replace highest ticket (tail) and become a new highest (also tail)", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 185];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [151, 42, 175, 7, 128, 185, 74, 143, 88, 130]; // 190 -> out
        // indices          = [ 0,  1,  2,  3,  4,   5,   6,   7,   8,   9 ]
        // sorted tickets   = [ 7, 42, 74, 88, 128, 130, 143, 151, 175, 185]
        // sorted indices   = [ 3,  1,  6,  8,  4,   9,   7,   0,   2,   5 ]
        // 0->7    1->3    2->0    3->3    4->8
        // 5->2    6->1    7->9    8->6    9->4
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 2, 1, 9, 6, 4];
        let expectedTail = 5;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

      it("should add a new smallest ticket and remove the highest", async () => {
        let ticketsToAdd = [5986, 6782, 5161, 7009, 8086, 1035, 5294, 9826, 6475, 9520, 4293, 998];
        
        await addTickets(ticketsToAdd)
        
        let expectedTickets = [5986, 6782, 5161, 7009, 8086, 1035, 5294, 4293, 6475, 998]; // 9826 & 9520 -> out
        // indices          = [ 0  ,  1  ,  2  ,  3  ,  4  ,  5  ,  6  ,  7  ,  8  ,  9 ]
        // sorted tickets   = [998, 1035, 4293, 5161, 5294, 5986, 6475, 6782, 7009, 8086]
        // sorted indices   = [ 9  ,  5  ,  7 ,   2 ,  6  ,   0  ,  8  ,  1  ,  3  ,  4 ]
        // 0->6    1->8    2->7    3->1    4->3
        // 5->9    6->2    7->5    8->0    9->9
        let expectedOrderedIndices = [6, 8, 7, 1, 3, 9, 2, 5, 0, 9];
        let expectedTail = 4;

        await assertTickets(expectedTail, expectedOrderedIndices, expectedTickets)
      });

    });

    describe("tickets array size is less than a group size", () => {

      it("should add all the tickets and keep track the order when the latest ticket is the highest one", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 6, 11];

        await addTickets(ticketsToAdd)

        // expected tickets = [1, 3, 5, 7, 4, 9, 6, 11]
        // indices          = [0, 1, 2, 3, 4, 5, 6,  7]
        // sorted tickets   = [1, 3, 4, 5, 6, 7, 9, 11]
        // sorted indices   = [0, 1, 4, 2, 6, 3, 5,  7]
        // 0->0    1->0    2->4    3->6
        // 4->1    5->3    6->2    7->5
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2, 5];
        let expectedTail = 7;

        await assertTickets(expectedTail, expectedOrderedIndices, ticketsToAdd)
      });

      it("should add all the tickets and track the order when the latest ticket is somewhere in the middle value range", async () => {
        let ticketsToAdd = [1, 3, 5, 7, 4, 9, 11, 6];

        await addTickets(ticketsToAdd)

        // expected tickets = [1, 3, 5, 7, 4, 9, 11, 6]
        // indices          = [0, 1, 2, 3, 4, 5, 6,  7]
        // sorted tickets   = [1, 3, 4, 5, 6, 7, 9, 11]
        // sorted indices   = [0, 1, 4, 2, 7, 3, 5,  6]
        // 0->0    1->0    2->4    3->7
        // 4->1    5->3    6->5    7->2
        let expectedOrderedIndices = [0, 0, 4, 7, 1, 3, 5, 2];
        let expectedTail = 6;

        await assertTickets(expectedTail, expectedOrderedIndices, ticketsToAdd)
      });

      it("should add all the tickets and track the order when the last added ticket is the smallest", async () => {
        let ticketsToAdd = [151, 42, 175, 7, 128, 190, 74, 4];

        await addTickets(ticketsToAdd)

        // expected tickets = [151, 42,  175,  7, 128, 190, 74,   4 ]
        // indices          = [ 0,   1,   2,   3,  4,   5,   6,   7 ]
        // sorted tickets   = [ 4,   7,  42,  74, 128, 151, 175, 190]
        // sorted indices   = [ 7,   3,   1,   6,  4,   0,   2,   5 ]
        // 0->4    1->3    2->0    3->7
        // 4->6    5->2    6->1    7->7
        let expectedOrderedIndices = [4, 3, 0, 7, 6, 2, 1, 7];
        let expectedTail = 5;

        await assertTickets(expectedTail, expectedOrderedIndices, ticketsToAdd);
      });

    });

    async function addTickets(ticketsToAdd) {
      for (let i = 0; i < ticketsToAdd.length; i++) {
        await groupSelectionStub.addTicket(ticketsToAdd[i]);
      }
    };

    async function assertTickets(expectedTail, expectedLinkedTicketIndices, expectedTickets) {
      // Assert tickets size
      let tickets = await groupSelectionStub.getTickets();
      assert.isAtMost(
        tickets.length,
        groupSize, 
        "array of tickets cannot have more elements than the group size"
      );

      // Assert ticket values
      let actualTickets = [];
      for (let i = 0; i < tickets.length; i++) {
        actualTickets.push(tickets[i])
      }
      assert.sameOrderedMembers(
        actualTickets.map(bn => bn.toNumber()),
        expectedTickets,
        "unexpected ticket values"
      )

      // Assert tail
      let tail = await groupSelectionStub.getTail()
      assert.equal(tail.toString(), expectedTail, "unexpected tail index")

      // Assert the order of the tickets[] indices
      let actualLinkedTicketIndices = [];
      for (let i = 0; i < tickets.length; i++) {
        let actualIndex = await groupSelectionStub.getPreviousTicketIndex(i)
        actualLinkedTicketIndices.push(actualIndex)
      }
      assert.sameOrderedMembers(
        actualLinkedTicketIndices.map(bn => bn.toNumber()),
        expectedLinkedTicketIndices,
        'unexpected order of tickets'
      );
    };

  });

});
