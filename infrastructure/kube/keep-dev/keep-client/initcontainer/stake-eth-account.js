const fs = require('fs');
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.HttpProvider('http://eth-tx-node.default.svc.cluster.local:8545'));

// Set contract abi file
const stakingProxyContractJsonFile = "/Users/sthompson22/projects/keep-core/contracts/solidity/build/contracts/StakingProxy.json";
const tokenStakingContractJsonFile = "/Users/sthompson22/projects/keep-core/contracts/solidity/build/contracts/TokenStaking.json";
const keepTokenContractJsonFile = "/Users/sthompson22/projects/keep-core/contracts/solidity/build/contracts/KeepToken.json";

// Select abi
const stakingProxyContractParsed = JSON.parse(fs.readFileSync(stakingProxyContractJsonFile));
const tokenStakingContractParsed = JSON.parse(fs.readFileSync(tokenStakingContractJsonFile));
const keepTokenContractParsed = JSON.parse(fs.readFileSync(keepTokenContractJsonFile));

const stakingProxyContractAbi = stakingProxyContractParsed.abi
const tokenStakingContractAbi = tokenStakingContractParsed.abi
const keepTokenContractAbi = keepTokenContractParsed.abi

// Set contracts
const stakingProxyContract = new web3.eth.Contract([stakingProxyContractAbi], 0xd51b7aEC4d83B187A7810E22f8DfAcbd84136451);
const tokenStakingContract = new web3.eth.Contract([tokenStakingContractAbi], 0xAb23e60EE417940903c1a440c31E8FA29837cb43);
const keepTokenContract = new web3.eth.Contract([keepTokenContractAbi], 0x9AF9C7d3B2720cBd8ea7088c3733eA3e797Ad402);

// \heimdall aliens numbers
function formatAmount(amount, decimals) {
  return web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals)))
}

// Stake a target eth account
async function stakeEthAccount() {
  let owner = "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"
  let magpie = "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"
  let operator = "0xA86c468475EF9C2ce851Ea4125424672C3F7e0C8"

  let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
  let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

  if (!await stakingProxyContract.isAuthorized(tokenStakingContract.address)) {
    stakingProxyContract.authorizeContract(tokenStakingContract.address);
  }

  //await keepTokenContract.approveAndCall(
  //  tokenStakingContract.address,
  //  formatAmount(1000000, 18),
  //  delegation,
  //  {from: owner});
}

stakeEthAccount()