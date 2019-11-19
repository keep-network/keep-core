import * as KeepToken from "./contracts/KeepToken.json"
import * as TokenStaking from "./contracts/TokenStaking.json"
import * as TokenGrant from "./contracts/TokenGrant.json"

export async function getKeepToken(web3) {
  const address = getContractAddress(KeepToken.default);

  const code = await web3.eth.getCode(address);
  
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');
  return new web3.eth.Contract(KeepToken.abi, address)
}

export async function getTokenStaking(web3) {
  const address = getContractAddress(TokenStaking.default);
  
  const code = await web3.eth.getCode(address);

  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');
  return new web3.eth.Contract(TokenStaking.abi, address)
}

export async function getTokenGrant(web3) {
  const address = getContractAddress(TokenGrant.default);
  
  const code = await web3.eth.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');
  return new web3.eth.Contract(TokenGrant.abi, address)
}

const getContractAddress = ({ networks }) => networks[Object.keys(networks)[0]].address;

