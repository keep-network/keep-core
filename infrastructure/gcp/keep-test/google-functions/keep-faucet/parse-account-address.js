// Attempts to read address to be funded from the http request to the faucet
const url = require('url')

exports.parseAccountAddress = (request, response) => {
  requestUrl = url.parse(request.url, true)
  if (requestUrl.query.account) {
    return (account = requestUrl.query.account)
  } else if (!/^(0x)?[0-9a-f]{40}$/i.test(account)) {
    // check if it has the basic requirements of an account
    // double thanks to the Ethereum folks for figuring this regex out already
    return response.send(
      'Improperly formatted account account, please try a valid one.',
    )
  } else {
    console.log('No account set in HTTP request.')
    return response.send(
      'No account account set, please set an account with ?account=<accountAddress>',
    )
  }
}
