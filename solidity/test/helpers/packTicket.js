const {web3} = require("@openzeppelin/test-environment")

function packTicket(ticketValueHex, index, operator) {
    let stakerValueBytes = web3.utils.hexToBytes(operator);

    let ticketBytes = web3.utils.hexToBytes(ticketValueHex)
    let ticketValue = ticketBytes.slice(0, 8) // Take the first 8 bytes of the ticket value

    let virtualStakerIndexPadded = web3.utils.padLeft(index, 8)
    let virtualStakerIndexBytes = web3.utils.hexToBytes(virtualStakerIndexPadded)

    return ticketValue.concat(stakerValueBytes).concat(virtualStakerIndexBytes)
}

module.exports = packTicket