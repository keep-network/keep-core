// Function deps
const url = require("url")
const Web3 = /** @type {import('web3').default} */ (
  /** @type {unknown } */ (require("web3"))
)
const HDWalletProvider = require("@truffle/hdwallet-provider")

/** @typedef { import('bn.js') } BN */

// Ethereum host info
const ethereumHost = process.env.ETHEREUM_HOST
const ethereumNetworkId = process.env.ETHEREUM_NETWORK_ID

// Keep contract owner info
const keepContractOwnerAddress = process.env.KEEP_CONTRACT_OWNER_ADDRESS
const keepContractOwnerPrivateKey = process.env.KEEP_CONTRACT_OWNER_PRIVATE_KEY
const keepContractOwnerProvider = new HDWalletProvider(
  `${keepContractOwnerPrivateKey}`,
  `${ethereumHost}`
)

// Contract artifacts
const tokenGrantJson = require("@keep-network/keep-core/artifacts/TokenGrant.json")
const keepTokenJson = require("@keep-network/keep-core/artifacts/KeepToken.json")
const permissiveStakingPolicyJson = require("@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json")

// We override transactionConfirmationBlocks and transactionBlockTimeout because they're
// 25 and 50 blocks respectively at default.  The result of this on small private testnets
// is long wait times for scripts to execute.
const web3Options = {
  defaultBlock: "latest",
  defaultGas: 4712388,
  transactionBlockTimeout: 25,
  transactionConfirmationBlocks: 3,
  transactionPollingTimeout: 480,
  handleRevert: true,
}

// Setup web3 provider.  We use the keepContractOwner since it needs to sign the approveAndCall transaction.
const web3 = new Web3(keepContractOwnerProvider)
web3.options = web3Options

// TokenGrant
const tokenGrantAbi = tokenGrantJson.abi
const tokenGrantAddress = tokenGrantJson.networks[ethereumNetworkId].address
const tokenGrant = new web3.eth.Contract(tokenGrantAbi, tokenGrantAddress)

// KeepToken
const keepTokenAbi = keepTokenJson.abi
const keepTokenAddress = keepTokenJson.networks[ethereumNetworkId].address
const keepToken = new web3.eth.Contract(keepTokenAbi, keepTokenAddress)

// PermissiveStakingPolicy
const permissiveStakingPolicyAddress =
  permissiveStakingPolicyJson.networks[ethereumNetworkId].address

const ethAccountRegExp = /^(0x)?[0-9a-f]{40}$/i
const tokenDecimalMultiplier = web3.utils.toBN(10).pow(web3.utils.toBN(18))

const baseGrantAmount = web3.utils.toBN(300e3) // 300k tokens
const grantAmount = baseGrantAmount.mul(tokenDecimalMultiplier) // to 18-decimal precision
const ERRORS = {
  "max-keep": {
    status: 400,
    message: `
      Token grant failed: your account has the maximum testnet KEEP allowed.\n
      You can manage your token grants at: https://dashboard.test.keep.network\n
      If you have questions, you can find us on Discord: https://discord.gg/jqxBU4m\n
    `,
  },
  "unexpected-revert": {
    status: 500,
    message: `
      Unexpected revert during token grant creation:\n
      {reason}
    `,
  },
  "unexpected-error": {
    status: 500,
    message: `
      Token grant failed, consider trying again.\n
      If problems persist find us on Discord: https://discord.gg/jqxBU4m .\n
    `,
  },
}
const SUCCESSES = {
  created: `
      Created token grant with {grantAmount} KEEP for account: {granteeAccount}\n
      You can follow the transaction at https://ropsten.etherscan.io/tx/{transactionHash}\n
      You can manage your token grants at: https://dashboard.test.keep.network .\n
      You can find us on Discord at: https://discord.gg/jqxBU4m .\n
  `,
}

const SECOND = 1
const SECONDS = SECOND
const MINUTE = 60 * SECONDS
const MINUTES = MINUTE
const HOUR = 60 * MINUTES
const HOURS = HOUR

exports.issueGrant = async (request, response) => {
  response.type("text/plain")
  try {
    const requestUrl = url.parse(request.url, true)
    const account = /** @type {string} */ (requestUrl.query.account)

    if (!account) {
      console.error("Unspecified account.")
      return response
        .status(400)
        .send(
          "No account address set, please set an account with ?account=<address>\n"
        )
    } else if (!ethAccountRegExp.test(account)) {
      console.error("Bad account address [", account, "].")
      return response
        .status(400)
        .send(
          "Improperly formatted account address, please correct and try again.\n"
        )
    } else {
      try {
        const content = await issueGrant(account, grantAmount)
        response.send(
          SUCCESSES[content.code].replace(
            /{(.*?)}/g,
            (_, property) => content && content[property]
          )
        )
      } catch (e) {
        if (e.payload && e.payload.code && ERRORS[e.payload.code]) {
          const { code, content } = e.payload
          const { status, message } = ERRORS[code] || {}
          console.error("Caught error issuing grant: ", e, "url: ", request.url)
          response
            .status(status)
            .send(
              (message || `unknown error with code ${code}`).replace(
                /{(.*?)}/g,
                (_, property) => content && content[property]
              )
            )
        } else {
          throw e
        }
      }
    }
  } catch (e) {
    const { status, message } = ERRORS["unexpected-error"]
    console.error("Caught unexpected error: ", e, "url: ", request.url)
    response.status(status).send(message)
  }
}

/**
 * @param {string} granteeAccount The account to issue the grant to, if the
 *        account currently has < the grant amount granted.
 * @param {BN} grantAmount The amount to grant.
 * @param {number} [currentNonce] The nonce to start with; if left off, resolved
 *        by calling `getTransactionCount` for the grant creator account.
 * @param {number} [gasPrice] If set, the gas price to use.
 */
async function issueGrant(granteeAccount, grantAmount, currentNonce, gasPrice) {
  console.log(`Fetching existing balance for account [${granteeAccount}]...`)
  const existingBalance = await existingGrantBalance(granteeAccount)
  console.log(
    `Existing balance for account [${granteeAccount}] is [${existingBalance}].`
  )

  if (existingBalance.gte(grantAmount)) {
    console.warn(
      `[${granteeAccount}] requested grant while at limit. Balance: [${existingBalance}].`
    )

    throw new PayloadError({ code: "max-keep" })
  } else {
    // Date.now is ms from epoch, start is seconds from epoch.
    const start = web3.utils.toBN(Math.floor(Date.now() / 1000))
    const cliff = web3.utils.toBN(48 * HOURS)
    // Unlock = cliff means everything unlocks at once.
    const unlockingDuration = cliff
    const revocable = true

    console.log(
      `Submitting grant for [${grantAmount}] to [${granteeAccount}]...`
    )
    const grantData = web3.eth.abi.encodeParameters(
      [
        "address",
        "address",
        "uint256",
        "uint256",
        "uint256",
        "bool",
        "address",
      ],
      [
        keepContractOwnerAddress,
        granteeAccount,
        unlockingDuration,
        start,
        cliff,
        revocable,
        permissiveStakingPolicyAddress,
      ]
    )

    const nonce =
      currentNonce ||
      (await web3.eth.getTransactionCount(keepContractOwnerAddress, "pending"))

    console.log(
      `Test submission for account ${granteeAccount} with nonce ${nonce}...`
    )
    // Try calling; if this throws, we'll have a proper error message thanks
    // to handleRevert above.
    try {
      await keepToken.methods
        .approveAndCall(tokenGrantAddress, grantAmount, grantData)
        .call({ from: keepContractOwnerAddress, nonce: nonce })
    } catch (e) {
      e.reason = e.reason || e.message
      throw new PayloadError({ code: "unexpected-revert", content: e })
    }

    return new Promise((resolve, reject) => {
      // If the call didn't revert, try submitting the transaction proper.
      console.log(
        `Submitting transaction for account ${granteeAccount} with nonce ${nonce}...`
      )
      keepToken.methods
        .approveAndCall(tokenGrantAddress, grantAmount, grantData)
        .send({ from: keepContractOwnerAddress, nonce: nonce })
        .on("transactionHash", (hash) => {
          console.log(
            `Submitted grant for [${grantAmount}] to [${granteeAccount}] ` +
              `with hash [${hash}]`,
            `and nonce [${nonce}]`
          )

          resolve({
            code: "created",
            transactionHash: hash,
            granteeAccount,
            grantAmount,
          })
        })
        .on("error", (error) => {
          if (
            // Confirmed transaction with this nonce, so bump it.
            (error.message && error.message == "nonce too low") ||
            // Pending transaction with this nonce, so bump it and queue up.
            (error.message &&
              error.message == "replacement transaction underpriced")
          ) {
            console.error(
              `Error with account grant transaction for ${granteeAccount}, ` +
                `nonce too low at [${nonce}], retry at [${nonce + 1}].`
            )
            console.log("Retrying transaction with higher nonce...")
            resolve(issueGrant(granteeAccount, grantAmount, nonce + 1))
          } else {
            reject(
              new PayloadError({ code: "unexpected-error", content: error })
            )
          }
        })
    })
  }
}

/**
 * @param {string} account The account whose grant balance to check.
 */
async function existingGrantBalance(account) {
  const grantBalanceString = await tokenGrant.methods
    .balanceOf(account)
    .call({}, "pending")

  return web3.utils.toBN(grantBalanceString)
}

class PayloadError extends Error {
  constructor(payload) {
    super(`Error with payload: ${JSON.stringify(payload)}`)
    this.payload = payload
  }
}
