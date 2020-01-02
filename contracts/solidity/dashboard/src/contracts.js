import KeepToken from './contracts/KeepToken.json'
import TokenStaking from './contracts/TokenStaking.json'
import TokenGrant from './contracts/TokenGrant.json'
import KeepRandomBeaconOperator from './contracts/KeepRandomBeaconOperator.json'

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

async function getContract(web3, contract) {
  const address = getContractAddress(contract)
  const code = await web3.eth.getCode(address)

  checkCodeIsValid(code)
  return new web3.eth.Contract(contract.abi, address)
}

function checkCodeIsValid(code) {
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address')
}

function getContractAddress({ networks }) {
  return networks[Object.keys(networks)[0]].address
};

