import exceptThrow from './helpers/expectThrow';
import {initContracts} from './helpers/initContracts';
import {bls} from './helpers/data';
const OperatorContract = artifacts.require('./KeepRandomBeaconOperatorStub.sol')

contract('TestKeepRandomBeaconServiceSelectOperator', function(accounts) {

  let config, stakingProxy, serviceContract, operatorContract, operatorContract2, operatorContract3;

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      OperatorContract
    );

    config = contracts.config;
    stakingProxy = contracts.stakingProxy;
    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;

    // Create and initialize additional operator contracts
    operatorContract2 = await OperatorContract.new();
    operatorContract2.initialize(
      stakingProxy.address, serviceContract.address, config.minimumStake, config.groupThreshold,
      config.groupSize, config.timeoutInitial, config.timeoutSubmission, config.timeoutChallenge, config.timeDKG, config.resultPublicationBlockStep,
      config.activeGroupsThreshold, config.groupActiveTime, config.relayRequestTimeout,
      bls.groupSignature, bls.groupPubKey
    );

    operatorContract3 = await OperatorContract.new();
    operatorContract3.initialize(
      stakingProxy.address, serviceContract.address, config.minimumStake, config.groupThreshold,
      config.groupSize, config.timeoutInitial, config.timeoutSubmission, config.timeoutChallenge, config.timeDKG, config.resultPublicationBlockStep,
      config.activeGroupsThreshold, config.groupActiveTime, config.relayRequestTimeout,
      bls.groupSignature, bls.groupPubKey
    );

  });

  it("service contract owner should be able to remove and add operator contracts.", async function() {

    let result = await serviceContract.selectOperatorContract();
    assert.equal(result, operatorContract.address, "Operator contract added during initialization should present in the service contract.");

    await serviceContract.removeOperatorContract(operatorContract.address);
    await exceptThrow(serviceContract.selectOperatorContract()); // Should revert since no operator contract present.

    await serviceContract.addOperatorContract(operatorContract.address);
    assert.equal(result, operatorContract.address, "Operator contract should be added");

  });

});
