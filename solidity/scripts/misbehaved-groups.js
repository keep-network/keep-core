/**
 * Important: the KeepRandomBeaconOperatorRewardsStub
 * and the TokenStakingSlashingStubs contract should be deployed to the network!
 * 
 * This script:
 *  - emits fake events (UnauthorizedSigningReported, RelayEntryTimeoutReported)
 *  */ 

// The KeepRandomBeaconOperatorRewardsStub contract should be deployed to the network
const KeepRandomBeaconOperator = artifacts.require("../stubs/KeepRandomBeaconOperatorRewardsStub.sol");
const crypto = require("crypto")

function getAccounts() {
  return new Promise((resolve, reject) => {
      web3.eth.getAccounts((error, accounts) => {
        resolve(accounts);
      });
  });
}

async function registerNewGroup (keepRandomBeaconOperator, accounts) {
  const groupPublicKey = crypto.randomBytes(32);
  await keepRandomBeaconOperator.registerNewGroup(groupPublicKey)
  await keepRandomBeaconOperator.setGroupMembers(groupPublicKey, [accounts[1], accounts[2], accounts[1]]);
}
  
module.exports = async function() {
    try {
        const accounts = await getAccounts();
        const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();
        const owner = accounts[0];
        await registerNewGroup(keepRandomBeaconOperator, accounts);

        await keepRandomBeaconOperator.reportUnauthorizedSigning(0,  Buffer.from('abc', 'hex'), { from: owner })
        await keepRandomBeaconOperator.reportRelayEntryTimeout({ from: owner })
      } catch(erorr) {
        console.log('Unexpected error', erorr);
        process.exit(1)
    }
    process.exit();
}
