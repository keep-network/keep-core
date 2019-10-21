import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
import mineBlocks from './helpers/mineBlocks';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";

contract('TestKeepRandomBeaconOperatorRelayEntryTimeout', function(accounts) {
  let operatorContract, serviceContract, fee;
  const blocksForward = 20;
  const requestCounter = 0;

  describe("RelayEntryTimeout", function() {

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
      serviceContract = contracts.serviceContract;

      // Using stub method to add first group to help testing.
      await operatorContract.registerNewGroup(bls.groupPubKey);
      // Passing a sender's authorization. accounts[0] is a msg.sender on blockchain
      await operatorContract.addServiceContract(accounts[0])
      fee = await serviceContract.entryFeeEstimate(0);
    });

    beforeEach(async () => {
      await createSnapshot()
    });

    afterEach(async () => {
      await restoreSnapshot()
    });

    it("should not throw an error when sigining is in progress and the block number > relay entry timeout", async function() {
      await operatorContract.sign(requestCounter, bls.seed, bls.previousEntry, {value: fee});
      mineBlocks(blocksForward)
      await operatorContract.sign(requestCounter, bls.seed, bls.previousEntry, {value: fee});

      assert.equal((await operatorContract.getPastEvents())[0].event, 'SignatureRequested', "SignatureRequested event should occur on operator contract.");
    })

    it("should throw an error when signing is in progress and the block number <= relay entry timeout", async function() {
      await operatorContract.sign(requestCounter, bls.seed, bls.previousEntry, {value: fee});

      await expectThrowWithMessage(operatorContract.sign(requestCounter, bls.seed, bls.previousEntry, {value: fee}), 'Relay entry is in progress.');
    })

    it("should not throw an error when sigining is not in progress and the block number > relay entry timeout", async function() {
      await operatorContract.sign(requestCounter, bls.seed, bls.previousEntry, {value: fee});

      assert.equal((await operatorContract.getPastEvents())[0].event, 'SignatureRequested', "SignatureRequested event should occur on operator contract.");
    })

  })
});