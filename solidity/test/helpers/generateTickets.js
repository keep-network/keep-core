const {web3} = require("@openzeppelin/test-environment")

function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
    let tickets = [];
    for (let i = 1; i <= stakerWeight; i++) {
      let ticketValueHex = web3.utils.soliditySha3({t: 'uint', v: randomBeaconValue}, {t: 'uint', v: stakerValue}, {t: 'uint', v: i})
      let ticketValue = web3.utils.toBN(ticketValueHex);
      let ticket = {
        valueHex: ticketValueHex,
        value: ticketValue,
        virtualStakerIndex: i
      }
      tickets.push(ticket);
    }
    return tickets
  }

  module.exports = generateTickets