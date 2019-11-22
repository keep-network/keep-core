import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import generateTickets from './helpers/generateTickets';

contract('KeepRandomBeaconOperator', function(accounts) {
  const groupSize = 64;
  const numberOfTickets = 200;
  const numbersToSlice = 15 // if more than 15 => Error: Number can only safely store up to 53 bits

  let operatorContract, tickets, ticketToAdd8Bytes;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorTicketsOrderingStub.sol')
    );

    operatorContract = contracts.operatorContract;
    operatorContract.setGroupSize(groupSize);

    tickets = generateTickets(await operatorContract.getGroupSelectionRelayEntry(), accounts[0], numberOfTickets);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("ticket packing estimation", () => {

    it("should estimate average gas usage for adding tickets when a ticket type is uint64", async () => {
      // 2^64 - 1 = 18,446,744,073,709,551,615
      // uint64 max number = 18446744073709551615
      
      let estimates = [], ticketToAdd8Bytes;

      console.log("tickets.length:: ", tickets.length)
      
      for (let i = 0; i < tickets.length; i++) {
        ticketToAdd8Bytes = new web3.utils.BN(tickets[i].value.toString().slice(0, numbersToSlice))
        // console.log(ticketToAdd8Bytes.toNumber())

        let estimate = await operatorContract.addTicket.estimateGas(ticketToAdd8Bytes)
        estimates.push(estimate);

        await operatorContract.addTicket(ticketToAdd8Bytes);
      }

      let estimatesSum = estimates.reduce((acc, val) => acc + val, 0);
      console.log('addTicket() average = ' + estimatesSum/tickets.length);
    });

    it("should estimate average gas usage for adding tickets when parameters are packed", async () => {
      // await operatorContract.submitTicket(
      //   tickets[0].value, 
      //   operator1, 
      //   1, 
      //   {from: operator1}
      // );
  
      // let submittedCount = await operatorContract.submittedTicketsCount();
      // assert.equal(1, submittedCount, "Ticket should be accepted");
    });

  });

});

