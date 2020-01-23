// Attempts to read address to be funded from the http request to the faucet
const url = require("url")

exports.parseAccountAddress = (request, response) => {
  requestUrl = url.parse(request.url, true)
  if (requestUrl.query.address) {
    return (address = requestUrl.query.address)
  } else if (!/^(0x)?[0-9a-f]{40}$/i.test(address)) {
    // check if it has the basic requirements of an address
    // double thanks to the Ethereum folks for figuring this regex out already
    return response.send(
      "Improperly formatted account address, please try a valid one.",
    )
  } else {
    console.log("No address set in HTTP request.")
    return response.send(
      "No account address set, please set an address with ?address=<accountAddress>",
    )
  }
}
