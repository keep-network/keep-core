import exceptThrow from './helpers/expectThrow';
import {initContracts} from './helpers/initContracts';

contract('TestKeepRandomBeaconServiceSelectOperator', function(accounts) {

  let serviceContract, operatorContract;

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );

    serviceContract = contracts.serviceContract;
    operatorContract = contracts.operatorContract;

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
