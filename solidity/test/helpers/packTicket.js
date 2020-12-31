const { web3 } = require("@openzeppelin/test-environment")

function packTicket(ticketValueHex, index, operator) {
  const stakerValueBytes = web3.utils.hexToBytes(operator)

  const ticketBytes = web3.utils.hexToBytes(ticketValueHex)
  const ticketValue = ticketBytes.slice(0, 8) // Take the first 8 bytes of the ticket value

  const virtualStakerIndexPadded = web3.utils.padLeft(index, 8)
  const virtualStakerIndexBytes = web3.utils.hexToBytes(
    virtualStakerIndexPadded
  )

  return ticketValue.concat(stakerValueBytes).concat(virtualStakerIndexBytes)
}

module.exports = packTicket
