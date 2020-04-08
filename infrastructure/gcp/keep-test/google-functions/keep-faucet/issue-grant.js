// Function deps
const Web3 = require('web3')
const HDWalletProvider = require('@truffle/hdwallet-provider')

// Ethereum host info
const ethereumHost = process.env.ETHEREUM_HOST
const ethereumNetworkId = process.env.ETHEREUM_NETWORK_ID

// Keep contract owner info
const keepContractOwnerAddress = process.env.KEEP_CONTRACT_OWNER_ADDRESS
const keepContractOwnerPrivateKey = process.env.KEEP_CONTRACT_OWNER_PRIVATE_KEY
const keepContractOwnerProvider = new HDWalletProvider(
  `${keepContractOwnerPrivateKey}`,
  `${ethereumHost}`,
)

// Contract artifacts
const tokenGrantJsonFile = require('./node_modules/@keep-network/keep-core/artifacts/TokenGrant.json')
const keepTokenJsonFile = require('./node_modules/@keep-network/keep-core/artifacts/KeepToken.json')

// Parse grantee account address
const { parseAccountAddress } = require('./parse-account-address.js')

// We override transactionConfirmationBlocks and transactionBlockTimeout because they're
// 25 and 50 blocks respectively at default.  The result of this on small private testnets
// is long wait times for scripts to execute.
const web3Options = {
  defaultBlock: 'latest',
  defaultGas: 4712388,
  transactionBlockTimeout: 25,
  transactionConfirmationBlocks: 3,
  transactionPollingTimeout: 480,
}

// Setup web3 provider.  We use the keepContractOwner since it needs to sign the approveAndCall transaction.
const web3 = new Web3(keepContractOwnerProvider, null, web3Options)

// TokenGrant
const tokenGrantAbi = tokenGrantJsonFile.abi
const tokenGrantAddress = tokenGrantJsonFile.networks[ethereumNetworkId].address
const tokenGrant = new web3.eth.Contract(tokenGrantAbi, tokenGrantAddress)

// KeepToken
const keepTokenAbi = keepTokenJsonFile.abi
const keepTokenAddress = keepTokenJsonFile.networks[ethereumNetworkId].address
const keepToken = new web3.eth.Contract(keepTokenAbi, keepTokenAddress)

exports.issueGrant = async (request, response) => {
  try {
    const granteeAccount = parseAccountAddress(request, response)
    const unlockingDuration = 0
    const start = Math.floor(Date.now() / 1000)
    const cliff = 0
    const revocable = true
    const tokens = 300000
    const grantBalance = await tokenGrant.methods
      .balanceOf(granteeAccount)
      .call()
    var grantAmount = formatAmount(tokens, 18)

    if (grantBalance.gte(grantAmount)) {
      console.log(
        `${granteeAccount} requested grant while at limit. Balance: ${grantBalance}`,
      )
      return response.send(`
        Token grant failed, your account has the maximum testnet KEEP allowed.
        You can manage your token grants at: https://dashboard.test.keep.network
        If you have questions find us on Discord: https://discord.gg/jqxBU4m\n`)
    } else {
      const grantData = Buffer.concat([
        Buffer.from(granteeAccount.substr(2), 'hex'),
        web3.utils.toBN(unlockingDuration).toBuffer('be', 32),
        web3.utils.toBN(start).toBuffer('be', 32),
        web3.utils.toBN(cliff).toBuffer('be', 32),
        Buffer.from(revocable ? '01' : '00', 'hex'),
      ])

      await keepToken.methods
        .approveAndCall(tokenGrant.address, grantAmount, grantData)
        .send({ from: keepContractOwnerAddress })

      console.log(
        `Created grant for ${web3.utils.toBN(
          grantAmount,
        )} to: ${granteeAccount}`,
      )
      response.send(`
        Created token grant with ${web3.utils.toBN(
          grantAmount,
        )} KEEP for account: ${granteeAccount}
        You can manage your token grants at: https://dashboard.test.keep.network
        You can find us on Discord at: https://discord.gg/jqxBU4m\n`)
    }
  } catch (error) {
    console.log(error)
    return response.send(`
        Token grant failed, try again.
        If problems persist find us on Discord: https://discord.gg/jqxBU4m\n`)
  }
}

function formatAmount(amount, decimals) {
  return (
    '0x' +
    web3.utils
      .toBN(amount)
      .mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals)))
      .toString('hex')
  )
}
