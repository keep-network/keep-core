// Attempts to read address to be funded from the http request to the faucet
const url = require('url')

exports.parseAccountAddress = (request, response) => {
  const requestUrl = url.parse(request.url, true)
  const account = requestUrl.query.account
  if (requestUrl.query.account) {
    return account
  } else if (!/^(0x)?[0-9a-f]{40}$/i.test(account)) {
    // check if it has the basic requirements of an account
    // double thanks to the Ethereum folks for figuring this regex out already
    return response.send(
      'Improperly formatted account account address, please try a valid one.\n',
    )
  } else {
    console.log('No account set in HTTP request.')
    return response.send(
      'No account account address set, please set an account with ?account=<accountAddress>\n',
    )
  }
}
