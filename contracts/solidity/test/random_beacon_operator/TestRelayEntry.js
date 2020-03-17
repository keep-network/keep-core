import expectThrow from '../helpers/expectThrow';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage';
import {bls} from '../helpers/data';
import {initContracts} from '../helpers/initContracts';
import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';

contract('KeepRandomBeaconOperator/RelayEntry', (accounts) => {
  let serviceContract, operatorContract;

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./stubs/TokenStakingStub.sol'),
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
    await operatorContract.setGroupMembers(group, [accounts[0], accounts[1], accounts[2]]);

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should keep relay entry submission at reasonable price", async () => {
    let gasEstimate = await operatorContract.relayEntry.estimateGas(bls.groupSignature);

    // Make sure no change will make the verification more expensive than it is
    // now or that even if it happens, it will be a conscious decision.
    assert.isBelow(gasEstimate, 378902, "Relay entry submission is too expensive")
  });

  it("should not allow to submit corrupted relay entry", async () => {
      // This is not a valid G1 point
      let groupSignature = "0x11134abfa2a9844a58776650e399bca3e08ab134e42595e03e3efc5a0472bcd8";

      await expectThrow(operatorContract.relayEntry(groupSignature));
  })

  it("should not allow to submit invalid relay entry", async () => {
      // Signature is a valid G1 point but it is not a signature over the
      // expected input.
      await expectThrowWithMessage(
        operatorContract.relayEntry(bls.nextGroupSignature),
        "Invalid signature"
      );
  });

  it("should allow to submit valid relay entry", async () => {
    await operatorContract.relayEntry(bls.groupSignature);

    assert.equal((await serviceContract.getPastEvents())[0].args['entry'].toString(),
      bls.groupSignatureNumber.toString(), "Should emit event with generated entry"
    );
  });

  it("should allow to submit only one entry", async () => {
    await operatorContract.relayEntry(bls.groupSignature);

    await expectThrowWithMessage(
      operatorContract.relayEntry(bls.groupSignature),
      "Entry was submitted"
    );
  });
});
