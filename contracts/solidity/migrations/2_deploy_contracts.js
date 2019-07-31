const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./cryptography/AltBn128.sol");
const BLS = artifacts.require("./cryptography/BLS.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const KeepRandomBeaconOperatorStub = artifacts.require("./KeepRandomBeaconOperatorStub.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;

const genesisEntry = web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862');
const genesisSeed = web3.utils.toBN('27182818284590452353602874713526624977572470936999595749669676277240766303535');
const genesisGroupPubKey = '0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0';

module.exports = async function(deployer) {
  await deployer.deploy(ModUtils);
  await deployer.link(ModUtils, AltBn128);
  await deployer.deploy(AltBn128);
  await deployer.link(AltBn128, BLS);
  await deployer.deploy(BLS);
  await deployer.deploy(KeepToken);
  await deployer.deploy(TokenStaking, KeepToken.address, withdrawalDelay);
  await deployer.deploy(TokenGrant, KeepToken.address, TokenStaking.address);
  await deployer.link(BLS, KeepRandomBeaconOperator);
  await deployer.link(BLS, KeepRandomBeaconOperatorStub);
  deployer.deploy(KeepRandomBeaconOperator);
  await deployer.deploy(KeepRandomBeaconServiceImplV1);
  await deployer.deploy(KeepRandomBeaconService, KeepRandomBeaconServiceImplV1.address);

  const keepRandomBeaconService = await KeepRandomBeaconServiceImplV1.at(KeepRandomBeaconService.address);
  const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

  keepRandomBeaconOperator.initialize(
    KeepRandomBeaconService.address,
    genesisEntry, genesisSeed, genesisGroupPubKey
  );

  // TODO: replace with a secure authorization protocol (addressed in RFC 11).
  keepRandomBeaconOperator.authorizeStakingContract(TokenStaking.address);

  keepRandomBeaconService.initialize(
    minPayment,
    withdrawalDelay,
    keepRandomBeaconOperator.address
  );
};
