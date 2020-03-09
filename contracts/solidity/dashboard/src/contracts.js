import KeepToken from './contracts/KeepToken.json'
import TokenStaking from './contracts/TokenStaking.json'
import TokenGrant from './contracts/TokenGrant.json'
import KeepRandomBeaconOperator from './contracts/KeepRandomBeaconOperator.json'
import Registry from './contracts/Registry.json'

export async function getKeepToken(web3) {
  return getContract(web3, KeepToken)
}

export async function getTokenStaking(web3) {
  return getContract(web3, TokenStaking)
}

export async function getTokenGrant(web3) {
  return getContract(web3, TokenGrant)
}

export async function getKeepRandomBeaconOperator(web3) {
  return getContract(web3, KeepRandomBeaconOperator)
}

export async function getRegistry(web3) {
  return getContract(web3, Registry)
}

export async function getKeepTokenContractDeployerAddress(web3) {
  const deployTransactionHash = getTransactionHashOfContractDeploy(KeepToken)
  const transaction = await web3.eth.getTransaction(deployTransactionHash)

  return transaction.from
}

export async function getContracts(web3) {
  const contracts = await Promise.all([
    getKeepToken(web3),
    getTokenGrant(web3),
    getTokenStaking(web3),
    getKeepRandomBeaconOperator(web3),
    getRegistry(web3),
  ])

  return {
    token: contracts[0],
    grantContract: contracts[1],
    stakingContract: contracts[2],
    keepRandomBeaconOperatorContract: contracts[3],
    registryContract: contracts[4],
  }
}

async function getContract(web3, contract) {
  const address = getContractAddress(contract)
  const code = await web3.eth.getCode(address)

  checkCodeIsValid(code)
  return new web3.eth.Contract(contract.abi, address)
}

function checkCodeIsValid(code) {
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address')
}

function getTransactionHashOfContractDeploy({ networks }) {
  return networks[Object.keys(networks)[0]].transactionHash
}

function getContractAddress({ networks }) {
  return networks[Object.keys(networks)[0]].address
};

