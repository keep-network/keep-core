/**
 * This script creates various types of grants.
 */

const KeepToken = artifacts.require("./KeepToken.sol");
const ManagedGrantFactory = artifacts.require("./ManagedGrantFactory.sol");
const ManagedGrant = artifacts.require("./ManagedGrant.sol");
const GuaranteedMinimumStakingPolicy = artifacts.require("./GuaranteedMinimumStakingPolicy.sol");

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

module.exports = async function () {
  try {
    const accounts = await getAccounts();
    const token = await KeepToken.deployed();
    const managedGrantFactoryContract = await ManagedGrantFactory.deployed();
    const policy = await GuaranteedMinimumStakingPolicy.deployed();

    const owner = accounts[0];
    const newGrantee = accounts[1];
    const grantee = owner;
    const duration = 1;
    const start =  (await web3.eth.getBlock("latest")).timestamp;
    const cliff = 1;
    const revocable = false;
    const policyAddress = policy.address;
    
    const managedGrantExtraData = web3.eth.abi.encodeParameters(
      ["address", "uint256", "uint256", "uint256", "bool", "address"],
      [
        grantee,
        duration,
        start,
        cliff,
        revocable,
        policyAddress,
      ]
    )
    // Create managed grant
    await token.approveAndCall(managedGrantFactoryContract.address, formatAmount(12300000, 18), managedGrantExtraData, { from: owner })
    // Get the address of managed grant contract from an event.
    const managedGrant1Event = (await managedGrantFactoryContract.getPastEvents())[0]
    const managedGrant1Address = managedGrant1Event.args['grantAddress'];
    const managedGrant1 = await ManagedGrant.at(managedGrant1Address);

    // Reeasign grantee
    await managedGrant1.requestGranteeReassignment(newGrantee, { from: owner })
    await managedGrant1.confirmGranteeReassignment(newGrantee, { from: owner })

    // Create a second managed grant
    await token.approveAndCall(managedGrantFactoryContract.address, formatAmount(456000000, 18), managedGrantExtraData, { from: owner })
  } catch (err) {
    console.error('unexpected error:', err)
    process.exit(1)
  }

  process.exit();
};
