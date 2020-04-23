/*
This is a deployment script using buidler.

To run it on ropsten it is required to set `CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY`
environment variable as it's used in buidler config. The private key should match
the account defined in the config.

Local network:
  npx buidler run --network local scripts/buidler-deploy.js

Ropsten network:
  CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY="0x...." \
    npx buidler run --network ropsten scripts/buidler-deploy.js
*/

const bre = require("@nomiclabs/buidler")
const path = require("path")
const fs = require("fs")
const clc = require("cli-color")

const dataDir = "data"
const dataFilePath = path.resolve(path.join(dataDir, "contracts.json"))

const initializationPeriod = 43200
const undelegationPeriod = 7776000

async function main() {
  const network = bre.buidlerArguments.network
  console.log("deploying contracts to network:", clc.magenta(network))

  const Registry = bre.artifacts.require("Registry")
  const KeepToken = bre.artifacts.require("KeepToken")
  const TokenGrant = bre.artifacts.require("TokenGrant")
  const ManagedGrantFactory = bre.artifacts.require("ManagedGrantFactory")
  const TokenStaking = bre.artifacts.require("TokenStaking")
  const PermissiveStakingPolicy = bre.artifacts.require(
    "PermissiveStakingPolicy"
  )
  const GuaranteedMinimumStakingPolicy = bre.artifacts.require(
    "GuaranteedMinimumStakingPolicy"
  )

  // REGISTRY
  const registry = await Registry.new()
  storeAddresses(Registry.contractName, registry)

  // TOKEN RELATED
  const keepToken = await KeepToken.new()
  storeAddresses(KeepToken.contractName, keepToken)

  const tokenGrant = await TokenGrant.new(keepToken.address)
  storeAddresses(TokenGrant.contractName, tokenGrant)

  const managedTokenGrant = await ManagedGrantFactory.new(
    keepToken.address,
    tokenGrant.address
  )
  storeAddresses(ManagedGrantFactory.contractName, managedTokenGrant)

  // STAKING RELATED
  const tokenStaking = await TokenStaking.new(
    keepToken.address,
    registry.address,
    initializationPeriod,
    undelegationPeriod
  )
  storeAddresses(TokenStaking.contractName, tokenStaking)

  const permissiveStakingPolicy = await PermissiveStakingPolicy.new()
  storeAddresses(PermissiveStakingPolicy.contractName, permissiveStakingPolicy)

  const guaranteedMinimumStakingPolicy = await GuaranteedMinimumStakingPolicy.new(
    tokenStaking.address
  )
  storeAddresses(
    GuaranteedMinimumStakingPolicy.contractName,
    guaranteedMinimumStakingPolicy
  )

  console.log(clc.green("deployment completed!"))
  console.debug(clc.blackBright("data stored in: ", dataFilePath))

  function storeAddresses(contractName, contract) {
    console.log(`deployed ${contractName} at ${contract.address}`)

    if (!fs.existsSync(dataDir)) {
      fs.mkdirSync(dataDir)
    }

    let content
    if (fs.existsSync(dataFilePath)) {
      content = JSON.parse(fs.readFileSync(dataFilePath, "utf8"))
    } else {
      content = {}
    }

    if (!content[network]) {
      content[network] = {}
    }

    content[network][contractName] = {
      transactionHash: contract.transactionHash,
      address: contract.address,
    }

    fs.writeFileSync(dataFilePath, JSON.stringify(content, null, 2))
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(clc.red(error))
    process.exit(1)
  })
