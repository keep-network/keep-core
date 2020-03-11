import expectThrow from '../helpers/expectThrow';
import {initContracts} from '../helpers/initContracts';

contract('TestKeepRandomBeaconOperator', function(accounts) {

  let operatorContract;

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperator.sol')
    );
    operatorContract = contracts.operatorContract;
  });

  it("should fail to update minimum stake by non owner", async function() {
    await expectThrow(operatorContract.setMinimumStake(123, {from: accounts[1]}));
  });

  it("should be able to update minimum stake by the owner", async function() {
    await operatorContract.setMinimumStake(123);
    let newMinStake = await operatorContract.minimumStake();
    assert.equal(newMinStake, 123, "Should be able to get updated minimum stake.");
  });
});
