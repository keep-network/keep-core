import expectThrow from './helpers/expectThrow';
import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';

contract('KeepRandomBeaconOperator', (accounts) => {
  let serviceContract, operatorContract;

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    );

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;

    // Using stub method to add first group to help testing.
    await operatorContract.registerNewGroup(bls.groupPubKey);
    operatorContract.setGroupSize(3);
    let group = await operatorContract.getGroupPublicKey(0);
    await operatorContract.addGroupMember(group, accounts[0]);
    await operatorContract.addGroupMember(group, accounts[1]);
    await operatorContract.addGroupMember(group, accounts[2]); 

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
  });

  it("should keep relay entry submission at reasonable price", async () => {
    let gasEstimate = await operatorContract.relayEntry.estimateGas(bls.groupSignature);

    // Make sure no change will make the verification more expensive than it is 
    // now or that even if it happens, it will be a conscious decision.
    assert.isBelow(gasEstimate, 369544, "Relay entry submission is too expensive")
  });

  it("should not allow to submit invalid relay entry", async () => {
      // Invalid signature
      let groupSignature = "0x0fb34abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8";

      await expectThrow(operatorContract.relayEntry(groupSignature));
    });

  it("should allow to submit valid relay entry", async () => {
    await operatorContract.relayEntry(bls.groupSignature);

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with generated entry"
    );
  });
});
