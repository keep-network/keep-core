const { web3 } = require("@openzeppelin/test-environment")

function generateTickets(randomBeaconValue, stakerValue, stakerWeight) {
  const tickets = []
  for (let i = 1; i <= stakerWeight; i++) {
    const ticketValueHex = web3.utils.soliditySha3(
      { t: "uint", v: randomBeaconValue },
      { t: "uint", v: stakerValue },
      { t: "uint", v: i }
    )
    const ticketValue = web3.utils.toBN(ticketValueHex)
    const ticket = {
      valueHex: ticketValueHex,
      value: ticketValue,
      virtualStakerIndex: i,
    }
    tickets.push(ticket)
  }
  return tickets
}

module.exports = generateTickets
