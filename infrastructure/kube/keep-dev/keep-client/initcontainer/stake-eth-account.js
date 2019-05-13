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

const stakingProxyContractAbi = stakingProxyContractParsed.abi;
const tokenStakingContractAbi = tokenStakingContractParsed.abi;
const keepTokenContractAbi = keepTokenContractParsed.abi;

// Set contracts
const stakingProxyContract = new web3.eth.Contract(stakingProxyContractAbi, "0xD4A23BF413c0C11084C6D25BCA1Afb2305781E80");
const tokenStakingContract = new web3.eth.Contract(tokenStakingContractAbi, "0x765E2963955b98E1789972277136D9D735c022e9");
const keepTokenContract = new web3.eth.Contract(keepTokenContractAbi, "0xdc3C6306a34005d3Eba0777E558715Fc2e21C5ba");

// \heimdall aliens numbers
function formatAmount(amount, decimals) {
  return '0x' + web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals))).toString('hex')
}

// Stake a target eth account
async function stakeEthAccount() {
  let owner = "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB";
  let magpie = "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB";
  let operator = "0xA86c468475EF9C2ce851Ea4125424672C3F7e0C8";

  let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
  let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

  try{
    if (!await stakingProxyContract.methods.isAuthorized(tokenStakingContract.address).send({from: owner})) {
      stakingProxyContract.methods.authorizeContract(tokenStakingContract.address).send({from: owner})
    }
  }
  catch(error) {
    console.error(error);
  }

  try {
    await keepTokenContract.methods.approveAndCall(
      tokenStakingContract.address,
      formatAmount(500, 18),
      delegation).send({from: owner, gas: 4712388})
  }
  catch(error) {
    console.error(error);
  }
};

stakeEthAccount();