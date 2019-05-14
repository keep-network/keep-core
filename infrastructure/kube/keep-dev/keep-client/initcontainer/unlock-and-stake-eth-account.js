const fs = require('fs');
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.HttpProvider(process.env.ETHEREUM_HOSTNAME + ":" + process.env.ETHEREUM_HOST_PORT));

// Contract setup
// stakingProxy
const stakingProxyContractJsonFile = "/tmp/StakingProxy.json";
const stakingProxyContractParsed = JSON.parse(fs.readFileSync(stakingProxyContractJsonFile));
const stakingProxyContractAbi = stakingProxyContractParsed.abi;
const stakingProxyContractAddress = stakingProxyContractParsed.networks[process.env.ETH_NETWORK_ID].address;
const stakingProxyContract = new web3.eth.Contract(stakingProxyContractAbi, stakingProxyContractAddress);

// tokenStaking
const tokenStakingContractJsonFile = "/tmp/TokenStaking.json";
const tokenStakingContractParsed = JSON.parse(fs.readFileSync(tokenStakingContractJsonFile));
const tokenStakingContractAbi = tokenStakingContractParsed.abi;
const tokenStakingContractAddress = tokenStakingContractParsed.networks[process.env.ETH_NETWORK_ID].address;
const tokenStakingContract = new web3.eth.Contract(tokenStakingContractAbi, tokenStakingContractAddress);

// keepToken
const keepTokenContractJsonFile = "/tmp/KeepToken.json";
const keepTokenContractParsed = JSON.parse(fs.readFileSync(keepTokenContractJsonFile));
const keepTokenContractAbi = keepTokenContractParsed.abi;
const keepTokenContractAddress = keepTokenContractParsed.networks[process.env.ETH_NETWORK_ID].address;
const keepTokenContract = new web3.eth.Contract(keepTokenContractAbi, keepTokenContractAddress);

// Ethereum account that contracts are migrated against
const contract_owner = process.env.CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS

// Stake a target eth account
async function stakeEthAccount() {

  let magpie = process.env.CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS;
  let operator = process.env.KEEP_CLIENT_ETH_ACCOUNT_ADDRESS;
  let operator_eth_account_password = process.env.KEEP_CLIENT_ETH_ACCOUNT_PASSWORD;

  let signature = Buffer.from((await web3.eth.sign(web3.utils.soliditySha3(contract_owner), operator)).substr(2), 'hex');
  let delegation = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex'), signature]).toString('hex');

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
    console.log("Unlocking account: " + operator);
    await web3.eth.personal.unlockAccount(operator, operator_eth_account_password, 700);
    console.log("Account: " + operator + " unlocked!");
  }
  catch(error) {
    console.error(error);
  }
}

// \heimdall aliens numbers
function formatAmount(amount, decimals) {
  return '0x' + web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals))).toString('hex')
}

unlockEthAccount();
stakeEthAccount();
