export default function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
    let tickets = [];
    for (let i = 1; i <= stakerWeight; i++) {
      let ticketValue = web3.utils.toBN(
        web3.utils.soliditySha3({t: 'uint', v: randomBeaconValue}, {t: 'uint', v: stakerValue}, {t: 'uint', v: i})
      );
      let ticket = {
        value: ticketValue,
        virtualStakerIndex: i
      }
      tickets.push(ticket);
    }
    return tickets
  }