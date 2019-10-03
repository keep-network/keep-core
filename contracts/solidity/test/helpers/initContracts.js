import { duration } from './increaseTime';
const BLS = artifacts.require('./cryptography/BLS.sol');

async function initContracts(KeepToken, TokenStaking, KeepRandomBeaconService,
  KeepRandomBeaconServiceImplV1, KeepRandomBeaconOperator, KeepRandomBeaconOperatorGroups) {

  let token, stakingContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    operatorContract, groupContract;

  let dkgContributionMargin = web3.utils.toBN(10).mul(web3.utils.toBN(10**18)), // Fraction in % of the estimated cost of DKG that is included in relay request payment. Must include 18 decimal points.
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
  groupContract = await KeepRandomBeaconOperatorGroups.new();
  operatorContract = await KeepRandomBeaconOperator.new(serviceContractProxy.address, stakingContract.address, groupContract.address);
  await groupContract.setOperatorContract(operatorContract.address);

  await serviceContract.initialize(dkgContributionMargin, withdrawalDelay, operatorContract.address);

  // Add initial funds to the fee pool to trigger group creation without waiting for DKG fee accumulation
  let dkgGasEstimate = await operatorContract.dkgGasEstimate();
  let priceFeedEstimate = await serviceContract.priceFeedEstimate();
  await serviceContract.fundDkgFeePool({value: dkgGasEstimate.mul(priceFeedEstimate)});

  // Genesis should include payment to cover DKG cost to create first group
  await operatorContract.genesis({value: dkgGasEstimate.mul(web3.utils.toBN(20*1.5).mul(web3.utils.toBN(10**9)))});

  return {
    token: token,
    stakingContract: stakingContract,
    serviceContract: serviceContract,
    operatorContract: operatorContract,
    groupContract: groupContract
  };
};

module.exports.initContracts = initContracts;
