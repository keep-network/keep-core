const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const Operator = artifacts.require("../stubs/KeepRandomBeaconOperatorRewardsStub.sol")
const KeepRandomBeaconServiceImpl = artifacts.require("./KeepRandomBeaconServiceImplV1.sol")
const crypto = require("crypto")

function formatAmount(amount, decimals) {
    return web3.utils.toBN(amount).mul(web3.utils.toBN(10).pow(web3.utils.toBN(decimals)))
  }
  
  function getAccounts() {
    return new Promise((resolve, reject) => {
        web3.eth.getAccounts((error, accounts) => {
          resolve(accounts);
        });
    });
  };
  
  module.exports = async function() {
    const accounts = await getAccounts();
    const token = await KeepToken.deployed();
    const tokenStaking = await TokenStaking.deployed();
    const keepRandomBeaconService = await KeepRandomBeaconServiceImpl.deployed()
    const keepRandomBeaconOerator = await Operator.deployed()

    const magpie = accounts[5];
    const requestor = accounts[5]
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
        console.log(`successfully staked KEEP tokens for account ${operator}`)
      }
    }

    const groupReward = web3.utils.toWei('14500', 'Gwei')
    for (let i = 0; i < 100; i++) {
      console.log('register group', i+1)
      const groupPubKey = crypto.randomBytes(32)
      await keepRandomBeaconOerator.registerNewGroup(groupPubKey).catch(err => console.log('register new group error', err))
      await keepRandomBeaconOerator.addGroupMemberReward(groupPubKey, groupReward)
      await keepRandomBeaconOerator.addGroupMember(groupPubKey, accounts[1])
      await keepRandomBeaconOerator.addGroupMember(groupPubKey, accounts[1])
      await keepRandomBeaconOerator.addGroupMember(groupPubKey, accounts[1])
      console.log('created group', await keepRandomBeaconOerator.getGroupPublicKey(i))
    }

    mineBlocks(3196)
    const entryFeeEstimate = await keepRandomBeaconService.entryFeeEstimate(0).catch(error => console.log('error service entry fee', error))
    await keepRandomBeaconService.methods['requestRelayEntry()']({value: entryFeeEstimate, from: requestor}).catch(error => console.log('request relay entry entry service error', error))

    const numberOfGroups = await keepRandomBeaconOerator.numberOfGroups().catch(error => console.log('error num of groups', error))
    const firstActiveIndex = await keepRandomBeaconOerator.getFirstActiveGroupIndex().catch(error => console.log('error first active', error))

    console.log('number of groups:', numberOfGroups.toString())
    console.log('last active index:', firstActiveIndex.toString())
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