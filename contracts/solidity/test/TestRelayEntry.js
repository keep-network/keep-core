import exceptThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';

contract('TestRelayEntry', function(accounts) {
  let frontend, operatorContract;

  before(async () => {

    let contracts = await initContracts(
      accounts,
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./StakingProxy.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconFrontend.sol'),
      artifacts.require('./KeepRandomBeaconFrontendImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperatorStub.sol')
    );
  
    operatorContract = contracts.operatorContract;
    frontend = contracts.frontend;
    // operatorContract.authorizeFrontendContract(frontend.address);

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    await frontend.requestRelayEntry(bls.seed, {value: 10});
  });

  it("should not be able to submit invalid relay entry", async function() {
    let requestID = 2;

    // Invalid signature
    let groupSignature = web3.utils.toBN('0x0fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8');

    await exceptThrow(operatorContract.relayEntry(requestID, groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed));
  });

  it("should be able to submit valid relay entry", async function() {
    let requestID = 2;

    await operatorContract.relayEntry(requestID, bls.groupSignature, bls.groupPubKey, bls.previousEntry, bls.seed);

    assert.equal((await frontend.getPastEvents())[0].args['requestResponse'].toString(),
      bls.groupSignature.toString(), "Should emit event with successfully submitted groupSignature."
    );

  });

});
