const fs = require('fs');
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.HttpProvider(process.env.ETHEREUM_HOSTNAME + ":" + process.env.ETHEREUM_HOST_PORT));

// Set compiled contract json
const stakingProxyContractJsonFile = "/Users/sthompson22/projects/keep-core/contracts/solidity/build/contracts/StakingProxy.json";
const tokenStakingContractJsonFile = "/Users/sthompson22/projects/keep-core/contracts/solidity/build/contracts/TokenStaking.json";
const keepTokenContractJsonFile = "/Users/sthompson22/projects/keep-core/contracts/solidity/build/contracts/KeepToken.json";

// We need to do some parse-fu to get at the values we need
const stakingProxyContractParsed = JSON.parse(fs.readFileSync(stakingProxyContractJsonFile));
const tokenStakingContractParsed = JSON.parse(fs.readFileSync(tokenStakingContractJsonFile));
const keepTokenContractParsed = JSON.parse(fs.readFileSync(keepTokenContractJsonFile));

// .abi is used to set contract functions
const stakingProxyContractAbi = stakingProxyContractParsed.abi;
const tokenStakingContractAbi = tokenStakingContractParsed.abi;
const keepTokenContractAbi = keepTokenContractParsed.abi;

// Set the current contract address for the chosen network
const stakingProxyContractAddress = stakingProxyContractParsed.networks[process.env.ETH_NETWORK_ID].address;
const tokenStakingContractAddress = tokenStakingContractParsed.networks[process.env.ETH_NETWORK_ID].address;
const keepTokenContractAddress = keepTokenContractParsed.networks[process.env.ETH_NETWORK_ID].address;

// Set contracts to be used for staking an Ethereum account
const stakingProxyContract = new web3.eth.Contract(stakingProxyContractAbi, stakingProxyContractAddress);
const tokenStakingContract = new web3.eth.Contract(tokenStakingContractAbi, tokenStakingContractAddress);
const keepTokenContract = new web3.eth.Contract(keepTokenContractAbi, keepTokenContractAddress);

// Ethereum account that contracts were migrated against
const contract_owner = process.env.CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS

// \heimdall aliens numbers
function formatAmount(amount, decimals) {
  return '0x' + web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals))).toString('hex')
}

// Stake a target eth account
async function stakeEthAccount() {
  let magpie = process.env.CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS;
  let operator = process.env.KEEP_CLIENT_ETH_ACCOUNT_ADDRESS;

  let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(contract_owner), operator)).substr(2), 'hex');
  let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

  console.log(signature)
  console.log(delegation)

  try{
    if (!await stakingProxyContract.methods.isAuthorized(tokenStakingContract.address).send({from: contract_owner}).then((receipt) => {
        console.log("isAuthorized transaction receipt:")
        console.log(receipt)
    })) {
      await stakingProxyContract.methods.authorizeContract(tokenStakingContract.address).send({from: contract_owner}).then((receipt) => {
        console.log("authorizeContract transaction receipt:")
        console.log(receipt)
      })
    }
    console.log("stakingProxy/tokenStaking contracts authorized!")
  }
  catch(error) {
    console.error(error);
  }

  try {
    await keepTokenContract.methods.approveAndCall(
      tokenStakingContract.address,
      formatAmount(1000000, 18),
      delegation).send({from: contract_owner, gas: 4712388}).then((receipt) => {
        console.log("Account " + operator + " staked!");
        console.log(receipt);
      });
  }
  catch(error) {
    console.error(error);
  }
};

async function unlockEthAccount() {

  try {
    console.log("Unlocking account: " + contract_owner);
    await web3.eth.personal.unlockAccount(contract_owner, process.env.KEEP_CLIENT_ETH_ACCOUNT_PASSPHRASE, 700);
    console.log("Account: " + contract_owner + " unlocked!");
  }
  catch(error) {
    console.error(error);
  }
}

//unlockEthAccount();
//stakeEthAccount();
