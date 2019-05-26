import exceptThrow from './helpers/expectThrow';
import {getContracts} from './helpers/initContracts';

contract('TestKeepRandomBeaconBackend', function(accounts) {

  let backend;

  beforeEach(async () => {
    let contracts = await getContracts(accounts);
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
