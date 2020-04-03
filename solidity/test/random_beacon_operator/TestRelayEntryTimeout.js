import expectThrowWithMessage from '../helpers/expectThrowWithMessage';
import {bls} from '../helpers/data';
import {initContracts} from '../helpers/initContracts';
import mineBlocks from '../helpers/mineBlocks';
import {createSnapshot, restoreSnapshot} from '../helpers/snapshot';

contract("KeepRandomBeaconOperator/RelayEntryTimeout", function(accounts) {
  let operatorContract, serviceContract, fee;
  const blocksForward = 20;
  const requestCounter = 0;

  before(async() => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    ); 

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;

    await operatorContract.addServiceContract(accounts[0])  

    await operatorContract.registerNewGroup(bls.groupPubKey);
    await operatorContract.setGroupMembers(bls.groupPubKey, [accounts[0]]);

    fee = await serviceContract.entryFeeEstimate(0);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should not throw an error when entry is in progress and " +
     "block number > relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, bls.previousEntry, {value: fee}
    );

    mineBlocks(blocksForward)

    await operatorContract.sign(
      requestCounter, bls.previousEntry, {value: fee}
    );

    assert.equal(
      (await operatorContract.getPastEvents())[0].event, 
      "RelayEntryRequested", 
      "RelayEntryRequested event should occur on operator contract"
    );
  });

  it("should throw an error when entry is in progress and " + 
     "block number <= relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, bls.previousEntry, {value: fee}
    );

    await expectThrowWithMessage(
      operatorContract.sign(requestCounter, bls.previousEntry, {value: fee}), 
      "Beacon is busy"
    );
  });

  it("should not throw an error when entry is not in progress and " + 
     "block number > relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, bls.previousEntry, {value: fee}
      );

    assert.equal(
      (await operatorContract.getPastEvents())[0].event, 
      "RelayEntryRequested", 
      "RelayEntryRequested event should occur on operator contract."
    );
  });

  it("should not allow to submit relay entry after timeout", async () => {
    await operatorContract.sign(
      requestCounter, bls.previousEntry, {value: fee}
    );

    mineBlocks(await operatorContract.relayEntryTimeout());

    await expectThrowWithMessage(
      operatorContract.relayEntry(bls.groupSignature), 
      "Entry timed out"
    );
  });
});
