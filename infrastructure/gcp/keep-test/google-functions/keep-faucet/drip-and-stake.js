const fs = require('fs');
const Web3 = require('web3');
const HDWalletProvider = require("@truffle/hdwallet-provider");

// Ethereum host info
const ethereumHost = process.env.ETHEREUM_HOST;
const ethereumNetworkId = process.env.ETHEREUM_NETWORK_ID;

// Keep contract owner info
const keepContractOwnerAddress = process.env.KEEP_CONTRACT_OWNER_ADDRESS;
const keepContractOwnerPrivateKey = process.env.KEEP_CONTRACT_OWNER_PRIVATE_KEY;
const keepContractOwnerProvider = new HDWalletProvider(`${keepContractOwnerPrivateKey}`, `${ethereumHost}`);

// Operator info
const { parseAccountAddress } = require('./parse-account-address.js')

// We override transactionConfirmationBlocks and transactionBlockTimeout because they're
// 25 and 50 blocks respectively at default.  The result of this on small private testnets
// is long wait times for scripts to execute.
const web3_options = {
    defaultBlock: 'latest',
    defaultGas: 4712388,
    transactionBlockTimeout: 25,
    transactionConfirmationBlocks: 3,
    transactionPollingTimeout: 480
};

// Setup web3 provider.  We use the keepContractOwner since it needs to sign the approveAndCall transaction.
const web3 = new Web3(keepContractOwnerProvider, null, web3_options);

// TokenStaking
const tokenStakingContractJsonFile = `./TokenStaking.json`;
const tokenStakingContractParsed = JSON.parse(fs.readFileSync(tokenStakingContractJsonFile));
const tokenStakingContractAbi = tokenStakingContractParsed.abi;
const tokenStakingContractAddress = tokenStakingContractParsed.networks[ethereumNetworkId].address;
const tokenStakingContract = new web3.eth.Contract(tokenStakingContractAbi, tokenStakingContractAddress);

// KeepToken
const keepTokenContractJsonFile = `./KeepToken.json`;
const keepTokenContractParsed = JSON.parse(fs.readFileSync(keepTokenContractJsonFile));
const keepTokenContractAbi = keepTokenContractParsed.abi;
const keepTokenContractAddress = keepTokenContractParsed.networks[ethereumNetworkId].address;
const keepTokenContract = new web3.eth.Contract(keepTokenContractAbi, keepTokenContractAddress);

exports.dripAndStake = async (request, response) => {

  try {
    let operatorAddress = parseAccountAddress(request, response)
    console.log(operatorAddress)
    let delegation = '0x' + Buffer.concat([
      Buffer.from(keepContractOwnerAddress.substr(2), 'hex'),
      Buffer.from(operatorAddress.substr(2), 'hex')
    ]).toString('hex');

    await keepTokenContract.methods.approveAndCall(
      tokenStakingContract.address,
      formatAmount(20000000, 18),
      delegation).send({from: keepContractOwnerAddress})

    console.log(`${operatorAddress} staked with 20000000 KEEP!`);
    response.send(
      `${operatorAddress} staked with 20000000 KEEP!`
    );

  } catch (error) {
    console.log(error);
    return response.send(
      'Staking failed, find an adult at Keep.'
    )
  }
};

function formatAmount(amount, decimals) {
  return '0x' + web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals))).toString('hex');
};
