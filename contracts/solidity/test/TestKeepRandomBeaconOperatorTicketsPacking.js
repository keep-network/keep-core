import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import generateTickets from './helpers/generateTickets';

contract('KeepRandomBeaconOperator', function(accounts) {
  const groupSize = 64;
  const numberOfTickets = 200;
  const numbersToSlice = 18 

  let operatorContract, tickets, ticketToAdd8Bytes;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorTicketsPackingStub.sol')
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

    it("should estimate average gas usage for addTickets() when a ticket type is uint64", async () => {
      // 2^64 - 1 = 18,446,744,073,709,551,615
      // uint64 max number = 18446744073709551615
      
      let estimates = [], estimate, ticketToAdd8Bytes;

      console.log("tickets.length:: ", tickets.length)
      
      for (let i = 0; i < tickets.length; i++) {
        ticketToAdd8Bytes = new web3.utils.BN(tickets[i].value.toString().slice(0, numbersToSlice))
        // console.log(ticketToAdd8Bytes.toNumber())

        estimate = await operatorContract.addTicket.estimateGas(ticketToAdd8Bytes)
        estimates.push(estimate);

        await operatorContract.addTicket(ticketToAdd8Bytes);
      }

      let estimatesSum = estimates.reduce((acc, val) => acc + val, 0);
      console.log('addTicket() average = ' + estimatesSum/tickets.length);
    });

    it("should estimate average gas usage for submitTicket() when parameters are packed", async () => {
      let ticket, ticketHex, ticketBytes, estimate,
      virtualStakerHex, virtualStakerIndexBytes, ticketBytesCombined,
      estimates = [];

      // staker value - 20bytes
      let stakerValueBytes = web3.utils.hexToBytes(accounts[0]) //staker address
      // console.log("stakerValueBytes.length", stakerValueBytes.length)
      // console.log("stakerValueBytes: ", stakerValueBytes)

      for (let i = 0; i < tickets.length; i++) {
        // ticket value - 8bytes
        ticket = new web3.utils.BN(tickets[i].value.toString().slice(0, numbersToSlice))
        ticketHex = web3.utils.toHex(ticket)
        ticketBytes = web3.utils.hexToBytes(ticketHex);
        // console.log("ticketBytes.length", ticketBytes.length)
        // console.log("ticketBytes: ", ticketBytes)
        
        // add virtual staker index - 4bytes
        virtualStakerHex = web3.utils.padLeft(i, 8)
        virtualStakerIndexBytes = web3.utils.hexToBytes(virtualStakerHex)
        // console.log("virtualStakerBytes.length", virtualStakerIndexBytes.length)
        // console.log("virtualStakerIndexBytes: ", virtualStakerIndexBytes)
        
        ticketBytesCombined = ticketBytes.concat(stakerValueBytes).concat(virtualStakerIndexBytes)
        // console.log("ticketBytesCombined.length", ticketBytesCombined.length)
        // console.log("ticketBytesCombined: ", ticketBytesCombined)
        
        estimate = await operatorContract.submitTicket.estimateGas(ticketBytesCombined)
        estimates.push(estimate);
  
        await operatorContract.submitTicket(
          ticketBytesCombined,
          {from: accounts[0]}
        );
      }

      let estimatesSum = estimates.reduce((acc, val) => acc + val, 0);
      console.log('submitTicket() average = ' + estimatesSum/tickets.length);
  
    });

  });

});

