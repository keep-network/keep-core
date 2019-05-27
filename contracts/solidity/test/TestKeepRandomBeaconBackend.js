import exceptThrow from './helpers/expectThrow';
import {initContracts} from './helpers/initContracts';

contract('TestKeepRandomBeaconBackend', function(accounts) {

  let backend;

  before(async () => {
    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconFrontendProxy.sol'),
      artifacts.require('./KeepRandomBeaconFrontendImplV1.sol'),
      artifacts.require('./KeepRandomBeaconBackend.sol')
    );
    backend = contracts.backend;
  });

  it("should fail to update minimum stake by non owner", async function() {
    await exceptThrow(backend.setMinimumStake(123, {from: accounts[1]}));
  });

  it("should be able to update minimum stake by the owner", async function() {
    await backend.setMinimumStake(123);
    let newMinStake = await backend.minimumStake();
    assert.equal(newMinStake, 123, "Should be able to get updated minimum stake.");
  });

  it("should be able to check if the implementation contract was initialized", async function() {
    assert.isTrue(await backend.initialized(), "Implementation contract should be initialized.");
  });
});
