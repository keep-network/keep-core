const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconImplV1 = artifacts.require("./KeepRandomBeaconImplV1.sol");
const KeepRandomBeacon = artifacts.require("./KeepRandomBeacon.sol");
const KeepGroupImplV1 = artifacts.require("./KeepGroupImplV1.sol");
const KeepGroup = artifacts.require("./KeepGroup.sol");


module.exports = async function() {

  const stakingProxy = await StakingProxy.deployed();
  const tokenStaking = await TokenStaking.deployed();
  const tokenGrant = await TokenGrant.deployed();
  const keepRandomBeacon = await KeepRandomBeacon.deployed();
  const keepRandomBeaconImplV1 = await KeepRandomBeaconImplV1.at(keepRandomBeacon.address);
  const keepGroup = await KeepGroup.deployed();
  const keepGroupImplV1 = await KeepGroupImplV1.at(keepGroup.address);

  // Authorize staking contracts to work via proxy
  if (!await stakingProxy.isAuthorized(tokenStaking.address)) {
    stakingProxy.authorizeContract(tokenStaking.address);
  }
  if (!await stakingProxy.isAuthorized(tokenGrant.address)) {
    stakingProxy.authorizeContract(tokenGrant.address);
  }

  // Initialize Keep Random beacon contract implementation
  if (!await keepRandomBeaconImplV1.initialized()) {
    await keepRandomBeaconImplV1.initialize(stakingProxy.address, 100, 200, 0);
  }

  // Initialize Keep Group contract implementation
  if (!await keepGroupImplV1.initialized()) {
    await keepGroupImplV1.initialize(2, 3, keepRandomBeacon.address);
  }
};
