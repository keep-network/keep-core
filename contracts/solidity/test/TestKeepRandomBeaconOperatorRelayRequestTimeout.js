import expectThrowWithMessage from './helpers/expectThrowWithMessage';
import {bls} from './helpers/data';
import {initContracts} from './helpers/initContracts';
import mineBlocks from './helpers/mineBlocks';

 contract('TestKeepRandomBeaconOperatorRelayRequestTimeout', function(accounts) {
  let serviceContract, operatorContract;
  const blocksForward = 20;

   describe("RelayRequestTimeout", function() {

     beforeEach(async () => {

       let contracts = await initContracts(
        accounts,
        artifacts.require('./KeepToken.sol'),
        artifacts.require('./StakingProxy.sol'),
        artifacts.require('./TokenStaking.sol'),
        artifacts.require('./KeepRandomBeaconService.sol'),
        artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
        artifacts.require('./KeepRandomBeaconOperatorStub.sol')
      );

      operatorContract = contracts.operatorContract;
      serviceContract = contracts.serviceContract;

       // Using stub method to add first group to help testing.
      await operatorContract.registerNewGroup(bls.groupPubKey);
    });

     it("should not throw an error when sigining is in progress and the block number > relay entry timeout", async function() {
      await serviceContract.requestRelayEntry(bls.seed, {value: 10});
      mineBlocks(blocksForward)
      await serviceContract.requestRelayEntry(bls.seed, {value: 10});

       assert.equal((await operatorContract.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on operator contract.");
    })

     it("should throw an error when sigining is in progress and the block number <= relay entry timeout", async function() {
      await serviceContract.requestRelayEntry(bls.seed, {value: 10});

       await expectThrowWithMessage(serviceContract.requestRelayEntry(bls.seed, {value: 10}), 'Relay entry request is in progress.');
    })

     it("should not throw an error when sigining is not in progress and the block number > relay entry timeout", async function() {
      mineBlocks(blocksForward)
      await serviceContract.requestRelayEntry(bls.seed, {value: 10});

       assert.equal((await operatorContract.getPastEvents())[0].event, 'RelayEntryRequested', "RelayEntryRequested event should occur on operator contract.");
    })
  })

 });