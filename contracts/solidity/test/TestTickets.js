import { AssertionError } from 'assert';

contract('TestTickets', function(accounts) {

  let ticketsContract;

  const Tickets = artifacts.require('./Tickets.sol');
  
  beforeEach(async () => {
    ticketsContract = await Tickets.new();
    await ticketsContract.cleanup()
  });


  // it("should be able to add tickets and read from the chain!!!", async function() {

  //   let numberOfTicketsToSubmit = 10;
  //   let randomNumbersArray = [];

  //   for (let i = 0; i < numberOfTicketsToSubmit; i++) {
  //      let randomNumber = Math.floor((Math.random() * 100) + 1); // 1 to 100
  //      if (!randomNumbersArray.includes(randomNumber)) {
  //        randomNumbersArray.push(randomNumber);
  //        await ticketsContract.submitTicket(randomNumber);
  //      }
  //   }
  //   let tail = await ticketsContract.getTail()
  //   console.log("ticketsContract.getTail()", tail.toString())

  //   let tickets = await ticketsContract.getTickets();
  //   console.log("tickets.lenght", tickets.length)
  //   for (let i = 0; i < tickets.length; i++) {
  //     console.log(tickets[i].toString())
  //   }
  // })

  describe("when tickets array size is less than a group size", function() {

    it("should be able to track tickets order when a last ticket is the biggest", async function() {
      let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 6, 11]; // expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2, 5] works
      // let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 6]; //expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2] works
      // let ticketsToSubmit = [1, 3, 5, 7, 4, 9]; //expectedOrderedIndices = [0, 0, 4, 2, 1, 3] works
      // let ticketsToSubmit = [1, 3, 5, 7, 4]; //expectedOrderedIndices = [0, 0, 4, 2, 1] works

      for (let i = 0; i < ticketsToSubmit.length; i++) {
          await ticketsContract.submitTicket(ticketsToSubmit[i]);
      }
      
      let tickets = await ticketsContract.getTickets();
      assert.equal(tickets.length, ticketsToSubmit.length, "array of tickets should be the size of: " + ticketsToSubmit.length)
      
      let tail = await ticketsContract.getTail()
      assert.equal(7, tail.toString(), "tail index should be equal to 7")
      
      let expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2, 5]
      // let expectedOrderedIndices = [0, 0, 4, 6, 1, 3, 2]
      // let expectedOrderedIndices = [0, 0, 4, 2, 1, 3]
      // let expectedOrderedIndices = [0, 0, 4, 2, 1]
      // await logTicketStatus(ticketsToSubmit, expectedOrderedIndices)

      for (let i = 0; i < ticketsToSubmit.length; i++) {
        let prevByIndex = await ticketsContract.getPreviousTicketsByIndex(i)
        assert.equal(expectedOrderedIndices[i] + '', prevByIndex.toString())
      }

    });

    it("should be able to track tickets order when a last ticket is between smallest and biggest", async function() {
      let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 11, 6]; // expectedOrderedIndices = [0, 0, 4, 7, 1, 3, 5, 2]
      // let ticketsToSubmit = [1, 3, 5, 7, 4, 9, 11]; // expectedOrderedIndices = [0, 0, 4, 2, 1, 3, 5] works

      for (let i = 0; i < ticketsToSubmit.length; i++) {
        await ticketsContract.submitTicket(ticketsToSubmit[i]);
      }

      let expectedOrderedIndices = [0, 0, 4, 7, 1, 3, 5, 2]
      // let expectedOrderedIndices = [0, 0, 4, 2, 1, 3, 5]
      // await logTicketStatus(ticketsToSubmit, expectedOrderedIndices)
      for (let i = 0; i < ticketsToSubmit.length; i++) {
        let prevByIndex = await ticketsContract.getPreviousTicketsByIndex(i)
        assert.equal(expectedOrderedIndices[i] + '', prevByIndex.toString())
      }
    });

    it("should be able to track order when a last ticket is the smallest", async function() {
      // let ticketsToSubmit = [2, 3, 5, 7, 4, 9, 11, 1]; //expectedOrderedIndices = [7, 0, 4, 2, 1, 3, 5, 0]
      let ticketsToSubmit = [2, 3, 5, 7, 4, 9, 11]; 
      
      for (let i = 0; i < ticketsToSubmit.length; i++) {
        await ticketsContract.submitTicket(ticketsToSubmit[i]);
      }

      // let expectedOrderedIndices = [7, 0, 4, 2, 1, 3, 5, 0]
      // // let expectedOrderedIndices = [0, 0, 4, 2, 1, 3, 5]
      // await logTicketStatus(ticketsToSubmit, expectedOrderedIndices)
      // for (let i = 0; i < ticketsToSubmit.length; i++) {
      //   let prevByIndex = await ticketsContract.getPreviousTicketsByIndex(i)
      //   assert.equal(expectedOrderedIndices[i] + '', prevByIndex.toString())
      // }
    });

    async function logTicketStatus(ticketsToSubmit, expectedPrev) {
      console.log("--------------")
      console.log("submitted tickets: [" + ticketsToSubmit.toString() + "]")

      console.log("--------------")
      console.log("submitted tickets length: ", ticketsToSubmit.length)

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
      for (let i = 0; i < ticketsToSubmit.length; i++) {
        console.log("prev_tickets[" + i + "] -> " + await ticketsContract.getPreviousTicketsByIndex(i) + " | " + expectedPrev[i]);
      }
      
      console.log("--------------")
      let j = await ticketsContract.getJIndex();
      console.log("j: ", j.toString())
    }

  });

  // it("should not be able to add a new ticket with higher value than existing highest", async function() {

    //   let numberOfTicketsToSubmit = 10;
    //   let randomNumbersArray = [];
    //   let maxRandomTicketValue = 100;

    //   for (let i = 0; i < numberOfTicketsToSubmit;) {
    //     let randomNumber = Math.floor((Math.random() * maxRandomTicketValue) + 1); // 1 to 100
    //     if (!randomNumbersArray.includes(randomNumber)) {
    //       randomNumbersArray.push(randomNumber);
    //       await ticketsContract.submitTicket(randomNumber);
    //       i++;
    //     }
    //   }

    //   await ticketsContract.submitTicket(maxRandomTicketValue + 1);

    //   let tickets = await ticketsContract.getTickets();
      
    //   assert.equal(tickets.length, numberOfTicketsToSubmit, "array of tickets should be the size of: " + numberOfTicketsToSubmit)
    // });


});
