import expectThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';

contract('TestKeepRandomBeaconOperatorRelayEntry', function(accounts) {
  let serviceContract, operatorContract, groupContract;

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol'),
      artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
    );

    operatorContract = contracts.operatorContract;
    groupContract = contracts.groupContract;
    serviceContract = contracts.serviceContract;

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    operatorContract.setGroupSize(3);
    let group = await groupContract.getGroupPublicKey(0);
    await operatorContract.addGroupMember(group, accounts[0]);
    await operatorContract.addGroupMember(group, accounts[1]);
    await operatorContract.addGroupMember(group, accounts[2]);

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});
  });

  it("should keep relay entry submission at reasonable price", async () => {
    let gasEstimate = await operatorContract.relayEntry.estimateGas(bls.nextGroupSignature);

    // Make sure no change will make the verification more expensive than it is 
    // now or that even if it happens, it will be a conscious change.
    assert.isBelow(gasEstimate, 462415, "Relay entry submission is too expensive")
  });

  it("should not be able to submit invalid relay entry", async function() {
    // Invalid signature
    let groupSignature = web3.utils.toBN('0x0fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8');

    await expectThrow(operatorContract.relayEntry(groupSignature));
  });

  it("should be able to submit valid relay entry", async function() {
    await operatorContract.relayEntry(bls.nextGroupSignature);

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.nextGroupSignature.toString(), "Should emit event with successfully submitted groupSignature."
    );
  });
});
