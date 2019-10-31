contract('TestTickets', function() {

  let ticketsContract, groupSize;

  const Tickets = artifacts.require('./Tickets.sol');
  
  describe("ticket insertion", function() {

    beforeEach(async () => {
      ticketsContract = await Tickets.new();
      groupSize = await ticketsContract.groupSize() // should be 10
    });

    describe("tickets array size is at its max capacity", function() {
      it("should reject a new ticket when it is higher than an existing highest one.", async function() {
        let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 6, 11, 8, 12, 100, 200, 300];
        
        await submitTickets(ticketsToSubmit)

        let expectedTail = 9;
        // smallest index points to itself
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 8, 2, 5, 3, 7];
        // await logTicketStatus(ticketsToSubmit, expectedOrderedIndices)
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)

      });
      
      it("should replace the highest existing with new ticket which is somewhere in the middle value range", async function() {
        let ticketsToSubmit = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 135]; // 190 -> out

        await submitTickets(ticketsToSubmit)
        
        let expectedTail = 2;
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 9, 1, 5, 6, 4];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

      it("should replace highest ticket (tail) and become a new highest one (also tail)", async function() {
        let ticketsToSubmit = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 185]; // 190 -> out

        await submitTickets(ticketsToSubmit)
        
        let expectedTail = 5;
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 2, 1, 9, 6, 4];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

      it("should add a new smallest ticket and remove the highest", async function() {
        let ticketsToSubmit = [151, 42, 175, 7, 128, 190, 74, 143, 88, 130, 2]; // 190 -> out

        await submitTickets(ticketsToSubmit)
        
        let expectedTail = 5;
        let expectedOrderedIndices = [7, 3, 0, 3, 8, 2, 1, 9, 6, 4];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });
    });

    
    describe("tickets array size is less than a group size", function() {
      it("should add all the tickets and keep track the order", async function() {
        let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 6, 11];

        await submitTickets(ticketsToSubmit)
        
        let expectedTail = 7;
        let expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2, 5];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

      it("should add all the tickets and track the order when a latest ticket is between smallest and biggest", async function() {
        let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 11, 6];

        await submitTickets(ticketsToSubmit)
        
        let expectedTail = 6;
        let expectedOrderedIndices = [0, 0, 4, 7, 1, 3, 5, 2];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)

      });

      it("should add all the tickets and track the order when a last added ticket is the smallest", async function() {
        let ticketsToSubmit = [2, 3, 5, 7, 4, 9, 11, 1];

        await submitTickets(ticketsToSubmit)
        
        let expectedTail = 6;
        let expectedOrderedIndices = [7, 0, 4, 2, 1, 3, 5, 7];
        await assertTicketsProperties(expectedTail, expectedOrderedIndices)
      });

    });

    async function submitTickets(ticketsToSubmit) {
      for (let i = 0; i < ticketsToSubmit.length; i++) {
        await ticketsContract.submitTicket(ticketsToSubmit[i]);
      }
    }

    async function assertTicketsProperties(expectedTail, expectedOrderedIndices) {
      let tickets = await ticketsContract.getTickets();
      let correctNumberOfTickets = tickets.length < groupSize ? tickets.length : groupSize;
      assert.equal(tickets.length, correctNumberOfTickets, "array of tickets should be the size of: " + correctNumberOfTickets)
      
      let tail = await ticketsContract.getTail()
      assert.equal(expectedTail, tail.toString(), "tail index should be equal to " + expectedTail)

      for (let i = 0; i < tickets.length; i++) {
        let prevByIndex = await ticketsContract.getOrderedLinkedTicketIndices(i)
        assert.equal(expectedOrderedIndices[i] + '', prevByIndex.toString())
      }
    }

    // I'd leave it for debugging purposes.
    async function logTicketStatus(ticketsToSubmit, expectedPrev) {
      console.log("--------------")
      console.log("submitted tickets: [" + ticketsToSubmit.toString() + "]")
  
      console.log("--------------")
      console.log("number submitted tickets: ", ticketsToSubmit.length)
  
      console.log("--------------")
      let tail = await ticketsContract.getTail();
      console.log("tail index: ", tail.toString());
  
      console.log("--------------")
      let ticketMaxValue = await ticketsContract.getTicketMaxValue();
      console.log("max value ticket[tail]: ", ticketMaxValue.toString());
  
      console.log("--------------")
      let ordered = await ticketsContract.getOrdered();
      console.log("ordered[" + ordered.toString() + "]")

      console.log("--------------")
      let tickets = await ticketsContract.getTickets();
      console.log("tickets on chain[" + tickets.toString() + "]")
      
      console.log("--------------")
      for (let i = 0; i < tickets.length; i++) {
        console.log("prev_tickets[" + i + "] -> " + await ticketsContract.getPreviousTicketsByIndex(i) + " | " + expectedPrev[i]);
      }
    }

  }) 

});
