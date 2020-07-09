/**
 * Important: the KeepRandomBeaconOperatorRewardsStub contract should be deployed to the network!
 * 
 * This script:
 *  - delegates stake to operators. For each operator uses the same beneficiary address,
 *  - adds 30 mock groups of which 19 makes stale,
 *  - emits fake withdrawal events.
 *  */ 

const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
// The KeepRandomBeaconOperatorRewardsStub contract should be deployed to the network
const KeepRandomBeaconOperator = artifacts.require("../stubs/KeepRandomBeaconOperatorRewardsStub.sol");
const KeepRandomBeaconServiceImpl = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const crypto = require("crypto")

function getAccounts() {
  return new Promise((resolve, reject) => {
      web3.eth.getAccounts((error, accounts) => {
        resolve(accounts);
      });
  });
}

function formatAmount(amount, decimals) {
  return web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals)))
}
  
module.exports = async function() {
  try {

    const accounts = await getAccounts();
    const token = await KeepToken.deployed();
    const tokenStaking = await TokenStaking.deployed();
    const contractService = await KeepRandomBeaconService.deployed();
    const keepRandomBeaconService = await KeepRandomBeaconServiceImpl.at(contractService.address);
    const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();
  
    const beneficiary = accounts[5];
    const requestor = accounts[5];
    const owner = accounts[0];
      
    for(let i = 0; i < 5; i++) {
      const operator = accounts[i]
      let authorizer = accounts[i]
  
      const delegation = '0x' + Buffer.concat([
        Buffer.from(beneficiary.substr(2), 'hex'),
        Buffer.from(operator.substr(2), 'hex'),
        Buffer.from(authorizer.substr(2), 'hex')
      ]).toString('hex');
  
      const staked = await token.approveAndCall(
        tokenStaking.address, 
        formatAmount(20000000, 18),
        delegation,
        {from: owner}
      ).catch((err) => {
        console.log(`could not stake KEEP tokens for ${operator}: ${err}`);
      });
  
      await tokenStaking.authorizeOperatorContract(operator, keepRandomBeaconOperator.address, { from: authorizer });
  
      if (staked) {
        console.log(`successfully staked KEEP tokens for account ${operator}`);
      }
    }
  
    await registerNewGroups(10);
    await registerNewGroups(10);
    
    const entryFeeEstimate = await keepRandomBeaconService.entryFeeEstimate(0)
    await keepRandomBeaconService.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})
  
    await registerNewGroups(10);
  
    // terminate groups
    await keepRandomBeaconOperator.reportUnauthorizedSigning(25, { from: owner })
    await keepRandomBeaconOperator.reportUnauthorizedSigning(26, { from: owner })

    const numberOfGroups = await keepRandomBeaconOperator.numberOfGroups()
    const firstActiveIndex = await keepRandomBeaconOperator.getFirstActiveGroupIndex()
  
    const allGroups = (await keepRandomBeaconOperator.getNumberOfCreatedGroups()).toNumber()
    for (let i = 0; i < allGroups; i++) {
      const groupPublicKey =  await keepRandomBeaconOperator.getGroupPublicKey(i);
      const isStale = await keepRandomBeaconOperator.isStaleGroup(groupPublicKey);
  
      console.log('group: ', groupPublicKey, 'isStale: ', isStale);
  }
  
    console.log('number of active groups:', numberOfGroups.toString());
    console.log('first active index:', firstActiveIndex.toString());

    async function registerNewGroups (numberOfGroups) {
      const groupReward = web3.utils.toWei('145000', 'Gwei');
      for (let i = 0; i < numberOfGroups; i++) {
        console.log('register group', i+1);
        const groupPublicKey = crypto.randomBytes(32);
        await keepRandomBeaconOperator.registerNewGroup(groupPublicKey, [accounts[1], accounts[2]])
        await keepRandomBeaconOperator.addGroupMemberReward(groupPublicKey, groupReward);
        console.log('created group', await keepRandomBeaconOperator.getGroupPublicKey(i));
      }
    }
  } catch(error) {
    console.log('unexpected error', error);
    process.exit(1);
  }

  process.exit();
}
