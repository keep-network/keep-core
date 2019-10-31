import { duration } from './increaseTime';
const BLS = artifacts.require('./cryptography/BLS.sol');
const GroupSelection = artifacts.require('./GroupSelection.sol');

async function initContracts(KeepToken, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator, KeepRandomBeaconOperatorGroups) {

  let token, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract, groupContract;

  let priceFeedEstimate = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)), // (20 Gwei = 20 * 10^9 wei)
    fluctuationMargin = 50, // 50%
    dkgContributionMargin = 10, // 10%
    withdrawalDelay = 1;

  // Initialize Keep token contract
  token = await KeepToken.new();

  // Initialize staking contract
  stakingContract = await TokenStaking.new(token.address, duration.days(30));

  // Initialize Keep Random Beacon service contract
  serviceContractImplV1 = await KeepRandomBeaconServiceImplV1.new();
  serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address);
  serviceContract = await KeepRandomBeaconServiceImplV1.at(serviceContractProxy.address)

  // Initialize Keep Random Beacon operator contract
  const bls = await BLS.new();
  await KeepRandomBeaconOperator.link("BLS", bls.address);
  const groupSelection = await GroupSelection.new();
  await KeepRandomBeaconOperator.link("GroupSelection", groupSelection.address);
  groupContract = await KeepRandomBeaconOperatorGroups.new();
  operatorContract = await KeepRandomBeaconOperator.new(serviceContractProxy.address, stakingContract.address, groupContract.address);
  await groupContract.setOperatorContract(operatorContract.address);

  await serviceContract.initialize(priceFeedEstimate, fluctuationMargin, dkgContributionMargin, withdrawalDelay, operatorContract.address);

  let dkgGasEstimate = await operatorContract.dkgGasEstimate();
  let gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(web3.utils.toBN(fluctuationMargin)).div(web3.utils.toBN(100)));

  // Genesis should include payment to cover DKG cost to create first group
  await operatorContract.genesis({value: dkgGasEstimate.mul(gasPriceWithFluctuationMargin)});

  return {
    token: token,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract,
    groupContract: groupContract
  };
};

module.exports.initContracts = initContracts;
