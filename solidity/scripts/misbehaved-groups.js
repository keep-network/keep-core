/**
 * Important: the KeepRandomBeaconOperatorRewardsStub
 * and the TokenStakingStub contract should be deployed to the network!
 * 
 * This script:
 *  - emits fake events (UnauthorizedSigningReported, RelayEntryTimeoutReported)
 *  */ 

// The KeepRandomBeaconOperatorRewardsStub contract should be deployed to the network
const Operator = artifacts.require("../stubs/KeepRandomBeaconOperatorRewardsStub.sol");
const crypto = require("crypto")

function getAccounts() {
  return new Promise((resolve, reject) => {
      web3.eth.getAccounts((error, accounts) => {
        resolve(accounts);
      });
  });
}

async function registerNewGroup (keepRandomBeaconOerator, accounts) {
  const groupPubKey = crypto.randomBytes(32);
  await keepRandomBeaconOerator.registerNewGroup(groupPubKey)
  await keepRandomBeaconOerator.setGroupMembers(groupPubKey, [accounts[1], accounts[2], accounts[1]]);
  console.log('created group', await keepRandomBeaconOerator.getGroupPublicKey(i));
}
  
module.exports = async function() {
    try {
        const accounts = await getAccounts();
        const keepRandomBeaconOerator = await Operator.deployed();
        const owner = accounts[0];
        await registerNewGroup(keepRandomBeaconOerator, accounts);

        await keepRandomBeaconOerator.reportUnauthorizedSigning(0,  Buffer.from('abc', 'hex'), { from: owner })
        await keepRandomBeaconOerator.reportRelayEntryTimeout({ from: owner })
      } catch(erorr) {
        console.log('Unexpected error', erorr);
        process.exit(1)
    }
    process.exit();
}
