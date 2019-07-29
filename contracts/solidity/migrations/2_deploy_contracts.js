const KeepToken = artifacts.require("./KeepToken.sol");
const ModUtils = artifacts.require("./utils/ModUtils.sol");
const AltBn128 = artifacts.require("./cryptography/AltBn128.sol");
const BLS = artifacts.require("./cryptography/BLS.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");
const KeepRandomBeaconService = artifacts.require("./KeepRandomBeaconService.sol");
const KeepRandomBeaconServiceImplV1 = artifacts.require("./KeepRandomBeaconServiceImplV1.sol");
const KeepRandomBeaconOperator = artifacts.require("./KeepRandomBeaconOperator.sol");
const KeepRandomBeaconOperatorStub = artifacts.require("./KeepRandomBeaconOperatorStub.sol");

const withdrawalDelay = 86400; // 1 day
const minPayment = 1;
const minStake = web3.utils.toBN(200000).mul(web3.utils.toBN(10**18));

const groupThreshold = 3;
const groupSize = 5;
const timeoutInitial = 4;
const timeoutSubmission = 4;
const timeoutChallenge = 4;
const resultPublicationBlockStep = 3;
const activeGroupsThreshold = 5;
const groupActiveTime = 300;
// Time in blocks it takes to execute relay entry signing.
// 1 state with state.MessagingStateDelayBlocks which is set to 1
// 1 state with state.MessagingStateActiveBlocks which is set to 3
const relayEntrySigningTime = 4

// Deadline in blocks for relay entry publication after the first 
// group member becomes eligible to submit the result.
// Deadline should not be shorter than the time it takes for the
// last group member to become eligible plus at least one block 
// to submit.
const relayEntryPublicationDeadline = 20

// The maximum time it may take for relay entry to appear on 
// chain after relay request has been published
const relayEntryTimeout = relayEntrySigningTime + relayEntryPublicationDeadline

// timeDKG - Timeout in blocks after DKG result is complete and ready to be published.
// 7 states with state.MessagingStateActiveBlocks which is set to 3
// 7 states with state.MessagingStateDelayBlocks which is set to 1
// the rest of the states use state.SilentStateDelayBlocks and
// state.SilentStateActiveBlocks which are both set to 0.
const timeDKG = 7*(3+1);

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
  await deployer.deploy(StakingProxy);
  await deployer.deploy(TokenStaking, KeepToken.address, StakingProxy.address, withdrawalDelay);
  await deployer.deploy(TokenGrant, KeepToken.address, StakingProxy.address, withdrawalDelay);
  await deployer.link(BLS, KeepRandomBeaconOperator);
  await deployer.link(BLS, KeepRandomBeaconOperatorStub);
  deployer.deploy(KeepRandomBeaconOperator);
  await deployer.deploy(KeepRandomBeaconServiceImplV1);
  await deployer.deploy(KeepRandomBeaconService, KeepRandomBeaconServiceImplV1.address);

  const keepRandomBeaconService = await KeepRandomBeaconServiceImplV1.at(KeepRandomBeaconService.address);
  const keepRandomBeaconOperator = await KeepRandomBeaconOperator.deployed();

  // Initialize contract genesis entry value and genesis group defined in Go client submitGenesisRelayEntry()
  keepRandomBeaconOperator.initialize(
    StakingProxy.address, KeepRandomBeaconService.address, minStake, groupThreshold, groupSize,
    timeoutInitial, timeoutSubmission, timeoutChallenge, timeDKG, resultPublicationBlockStep,
    activeGroupsThreshold, groupActiveTime, relayEntryTimeout,
    [genesisEntry, genesisSeed], genesisGroupPubKey
  );

  keepRandomBeaconService.initialize(
    minPayment,
    withdrawalDelay,
    keepRandomBeaconOperator.address
  );
};
