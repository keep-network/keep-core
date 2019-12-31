const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const Operator = artifacts.require("../stubs/KeepRandomBeaconOperatorRewardsStub.sol")

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
    let magpie = accounts[6];
  
    let owner = accounts[0];
    
    for(let i = 0; i < 5; i++) {
      let operator = accounts[i]
       // The address where the rewards for participation are sent.
  
      // The owner provides to the contract a magpie address and the operator address. 
      let delegation = '0x' + Buffer.concat([
        Buffer.from(magpie.substr(2), 'hex'),
        Buffer.from(operator.substr(2), 'hex')
      ]).toString('hex');
  
      staked = await token.approveAndCall(
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

    console.log('check magpie operators')
    const operators = await tokenStaking.operatorsOfMagpie(magpie)
    console.log('magpie operators', operators)

    const keepRandomBeaconOerator = await Operator.deployed().catch(err => console.log('keepRandomBeaconOerator deploy', err))

    let groupPubKey = '0x' + Buffer.concat([Buffer.from(magpie.substr(2), 'hex')]).toString('hex');
    console.log('register new group', groupPubKey)
    await keepRandomBeaconOerator.registerNewGroup(groupPubKey).catch(err => console.log('register new group error', err))
    console.log('registered group public key', await keepRandomBeaconOerator.getGroupPublicKey(0))

    const numberOfGroups = await keepRandomBeaconOerator.numberOfGroups()
    console.log('number of groups', numberOfGroups.toString())
    await keepRandomBeaconOerator.addGroupMemberReward(groupPubKey, '2000')
    await keepRandomBeaconOerator.addGroupMember(groupPubKey, accounts[1])
    await keepRandomBeaconOerator.addGroupMember(groupPubKey, accounts[1])

    console.log('group member rewards', (await keepRandomBeaconOerator.getGroupMemberRewards(groupPubKey)).toString())
    console.log('is stable', await keepRandomBeaconOerator.isStaleGroup(groupPubKey))

    const indices = await keepRandomBeaconOerator.getGroupMemberIndices(groupPubKey, accounts[1], { from: accounts[1] })

    indices.forEach(i => {
        console.log('group member indices', i.toString())
    })
}