/**
 * Important: the KeepRandomBeaconOperatorRewardsStub contract should be deployed to the network!
 * 
 * This script:
 *  - delegates stake to operators. For each operator uses the same beneficiary(magpie) address,
 *  - adds 20 mock groups of which 9 makes stale,
 *  - emits fake withdrawal events.
 *  */ 

const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
// The KeepRandomBeaconOperatorRewardsStub contract should be deployed to the network
const Operator = artifacts.require("../stubs/KeepRandomBeaconOperatorRewardsStub.sol");
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

function mineBlocks(blocks) {
  for (let i = 0; i < blocks; i++) {
    web3.currentProvider.send({
      jsonrpc: "2.0",
      method: "evm_mine",
      id: 12345
    }, function(err, _) {
      if (err) console.log("Error mining a block. " + i, err)
    });
  }
}

function formatAmount(amount, decimals) {
  return web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals)))
}
  
module.exports = async function() {
  const accounts = await getAccounts();
  const token = await KeepToken.deployed();
  const tokenStaking = await TokenStaking.deployed();
  const contractService = await KeepRandomBeaconService.deployed();
  const keepRandomBeaconService = await KeepRandomBeaconServiceImpl.at(contractService.address);
  const keepRandomBeaconOerator = await Operator.deployed();

  const magpie = accounts[5];
  const requestor = accounts[5];
  const owner = accounts[0];
    
  for(let i = 0; i < 5; i++) {
    const operator = accounts[i]
    const delegation = '0x' + Buffer.concat([
      Buffer.from(magpie.substr(2), 'hex'),
      Buffer.from(operator.substr(2), 'hex')
    ]).toString('hex');

    const staked = await token.approveAndCall(
      tokenStaking.address, 
      formatAmount(20000000, 18),
      delegation,
      {from: owner}
    ).catch((err) => {
      console.log(`could not stake KEEP tokens for ${operator}: ${err}`);
    });

    if (staked) {
      console.log(`successfully staked KEEP tokens for account ${operator}`);
    }
  }

  await registerNewGroups(10);

  mineBlocks(10);
  const entryFeeEstimate = await keepRandomBeaconService.entryFeeEstimate(0)
  await keepRandomBeaconService.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor})

  await registerNewGroups(10);

  const numberOfGroups = await keepRandomBeaconOerator.numberOfGroups()
  const firstActiveIndex = await keepRandomBeaconOerator.getFirstActiveGroupIndex()

  for (let i = 0; i < firstActiveIndex; i++) {
    const groupPubKey =  await keepRandomBeaconOerator.getGroupPublicKey(i);
    const isStale = await keepRandomBeaconOerator.isStaleGroup(groupPubKey);

    console.log('group: ', groupPubKey, 'isStale: ', isStale);
  }

  console.log('number of groups:', numberOfGroups.toString());
  console.log('first active index:', firstActiveIndex.toString());

  await keepRandomBeaconOerator.emitRewardsWithdrawnEvent(accounts[1], 1)
  await keepRandomBeaconOerator.emitRewardsWithdrawnEvent(accounts[1], 3)
  await keepRandomBeaconOerator.emitRewardsWithdrawnEvent(accounts[1], 1)
  await keepRandomBeaconOerator.emitRewardsWithdrawnEvent(accounts[1], 5)
  await keepRandomBeaconOerator.emitRewardsWithdrawnEvent(accounts[1], 6)  

  async function registerNewGroups (numberOfGroups) {
    const groupReward = web3.utils.toWei('14500', 'Gwei');
    for (let i = 0; i < numberOfGroups; i++) {
      console.log('register group', i+1);
      const groupPubKey = crypto.randomBytes(32);
      await keepRandomBeaconOerator.registerNewGroup(groupPubKey)
      await keepRandomBeaconOerator.addGroupMemberReward(groupPubKey, groupReward);
      await keepRandomBeaconOerator.setGroupMembers(groupPubKey, [accounts[1], accounts[1], accounts[1]]);
      console.log('created group', await keepRandomBeaconOerator.getGroupPublicKey(i));
    }
  }
}
