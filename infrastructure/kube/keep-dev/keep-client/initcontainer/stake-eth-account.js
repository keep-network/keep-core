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
const stakingProxyContract = new web3.eth.Contract(stakingProxyContractAbi, "0x9F77E21dB16ef8218640Be6E62472B675366377f");
const tokenStakingContract = new web3.eth.Contract(tokenStakingContractAbi, "0xeFB960Ca430F982c0b6650C3544fD8244A117407");
const keepTokenContract = new web3.eth.Contract(keepTokenContractAbi, "0x67a6c635b967fDBDB313eF6C043117b6780f978E");

const owner = "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB"

// \heimdall aliens numbers
function formatAmount(amount, decimals) {
  return '0x' + web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals))).toString('hex')
}

// Stake a target eth account
async function stakeEthAccount() {
  let magpie = "0x0F0977c4161a371B5E5eE6a8F43Eb798cD1Ae1DB";
  let operator = "0xa924d3a62b2d515235e5de5d903c405cba7f0e86";

  let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(owner), operator)).substr(2), 'hex');
  let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

  console.log(signature)
  console.log(delegation)

  try{
    if (!await stakingProxyContract.methods.isAuthorized(tokenStakingContract.address).send({from: owner}).then((receipt) => {
        console.log("isAuthorized Transaction Receipt:")
        console.log(receipt)
    })) {
      await stakingProxyContract.methods.authorizeContract(tokenStakingContract.address).send({from: owner}).then((receipt) => {
        console.log("authorizeContract Transaction Receipt:")
        console.log(receipt)
      })
    }
    console.log("stakingProxy/tokenStaking Contracts Authorized!")
  }
  catch(error) {
    console.error(error);
  }

  try {
    await keepTokenContract.methods.approveAndCall(
      tokenStakingContract.address,
      formatAmount(1000000, 18),
      delegation).send({from: owner, gas: 4712388}).then((receipt) => {
        console.log("Account" + operator + "staked!");
        console.log(receipt);
      });
  }
  catch(error) {
    console.error(error);
  }
};

async function unlockEthAccount() {

  try {
    console.log("Unlocking account: " + owner);
    await web3.eth.personal.unlockAccount(owner, "doughnut_armenian_parallel_firework_backbite_employer_singlet", 700);
    console.log("Account: " + owner + " unlocked!");
  }
  catch(error) {
    console.error(error);
  }
}

unlockEthAccount();
stakeEthAccount();