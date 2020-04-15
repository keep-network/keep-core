// Function deps
const url = require('url')
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
const tokenGrantJson = require('@keep-network/keep-core/artifacts/TokenGrant.json')
const keepTokenJson = require('@keep-network/keep-core/artifacts/KeepToken.json')
const permissiveStakingPolicyJson = require('@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json')

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
const tokenGrantAbi = tokenGrantJson.abi
const tokenGrantAddress = tokenGrantJson.networks[ethereumNetworkId].address
const tokenGrant = new web3.eth.Contract(tokenGrantAbi, tokenGrantAddress)
tokenGrant.options.handleRevert = true

// KeepToken
const keepTokenAbi = keepTokenJson.abi
const keepTokenAddress = keepTokenJson.networks[ethereumNetworkId].address
const keepToken = new web3.eth.Contract(keepTokenAbi, keepTokenAddress)
keepToken.options.handleRevert = true

// PermissiveStakingPolicy
const permissiveStakingPolicyAddress = permissiveStakingPolicyJson.networks[ethereumNetworkId].address

const ethAccountRegExp = /^(0x)?[0-9a-f]{40}$/i
const tokenDecimalMultiplier = web3.utils.toBN(10).pow(web3.utils.toBN(18))

exports.issueGrant = async (request, response) => {
  response.type('text/plain')
  try {
    const requestUrl = url.parse(request.url, true)
    const account = requestUrl.query.account

    if (! account) {
      console.error("Unspecified account.")
      return response.status(400).send(
        "No account address set, please set an account with ?accoun=<address>\n"
      )
    } else if (!ethAccountRegExp.test(account)) {
      console.error("Bad account address [", account, "].")
      return response.status(400).send(
        "Improperly formatted account address, please correct and try again.\n"
      )
    } else {
      const granteeAccount = account
      const unlockingDuration = web3.utils.toBN(0)
      const start = web3.utils.toBN(Math.floor(Date.now() / 1000))
      const cliff = web3.utils.toBN(0)
      const revocable = false
      const tokens = web3.utils.toBN(300000)
      console.log(`Fetching existing balance for account [${granteeAccount}]...`)
      const grantBalanceString = await tokenGrant.methods
        .balanceOf(granteeAccount)
        .call()
      const grantBalance = web3.utils.toBN(grantBalanceString)
      console.log(`Existing balance for account [${granteeAccount}] is [${grantBalance}].`)
      const grantAmount = tokens.mul(tokenDecimalMultiplier)

      if (grantBalance.gte(grantAmount)) {
        console.warn(
          `[${granteeAccount}] requested grant while at limit. Balance: [${grantBalance}].`,
        )
        return response.status(400).send(`
          Token grant failed: your account has the maximum testnet KEEP allowed.\n
          You can manage your token grants at: https://dashboard.test.keep.network\n
          If you have questions, you can find us on Discord: https://discord.gg/jqxBU4m\n`
        )
      } else {
        console.log(
          `Submitting grant for [${grantAmount}] to [${granteeAccount}]...`,
        )
        const grantData = Buffer.concat([
          Buffer.from(granteeAccount.substr(2), 'hex', 20),
          unlockingDuration.toBuffer('be', 32),
          start.toBuffer('be', 32),
          cliff.toBuffer('be', 32),
          Buffer.from(revocable ? '01' : '00', 'hex'),
          Buffer.from(permissiveStakingPolicyAddress.substr(2), 'hex', 20)
        ])

        console.log("Test submission...")
        // Try calling; if this throws, we'll have a proper error message thanks
        // to handleRevert above.
        await keepToken.methods
          .approveAndCall(tokenGrantAddress, grantAmount, grantData)
          .call({ from: keepContractOwnerAddress })

        console.log("Submitting transaction...")
        // If the call didn't revert, try submitting the transaction proper.
        keepToken.methods
          .approveAndCall(tokenGrantAddress, grantAmount, grantData)
          .send({ from: keepContractOwnerAddress })
          .on('transactionHash', (hash) => {
            console.log(
              `Submitted grant for [${grantAmount}] to [${granteeAccount}] ` +
                `with hash [${hash}].`,
            )
            response.send(`
              Created token grant with ${grantAmount} KEEP for account: ${granteeAccount}\n
              You can follow the transaction at https://ropsten.etherscan.io/tx/${hash}\n
              You can manage your token grants at: https://dashboard.test.keep.network .\n
              You can find us on Discord at: https://discord.gg/jqxBU4m .\n
            `)
          })
          .on('error', (error) => {
            console.error(
              `Error with account grant transaction: [${error}]; URL was [${request.url}].`
            )
            if (! response.headersSent) {
              response.status(500).send(`
                  Token grant failed, try again.\n
                  If problems persist find us on Discord: https://discord.gg/jqxBU4m .\n
              `)
            }
          })
      }
    }
  } catch (error) {
    console.error(
      `Error while requesting account grant: [${error}]; URL was [${request.url}].`
    )
    return response.status(500).send(`
        Token grant failed, try again.\n
        If problems persist find us on Discord: https://discord.gg/jqxBU4m .\n
    `)
  }
}
