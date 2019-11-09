import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import {bls} from './helpers/data';
import mineBlocks from './helpers/mineBlocks';

contract('KeepRandomBeaconOperator', function(accounts) {
    let serviceContract;
    let operatorContract;

    before(async () => {
        let contracts = await initContracts(
          artifacts.require('./KeepToken.sol'),
          artifacts.require('./TokenStaking.sol'),
          artifacts.require('./KeepRandomBeaconService.sol'),
          artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
          artifacts.require('./stubs/KeepRandomBeaconOperatorPricingStub.sol'),
          artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
        );
    
        serviceContract = contracts.serviceContract;
        operatorContract = contracts.operatorContract;

        await operatorContract.registerNewGroup(bls.groupPubKey);
    });

    beforeEach(async () => {
        await createSnapshot()
    });
    
    afterEach(async () => {
      await restoreSnapshot()
    });

    it("should correctly evaluate delay factor right after the request", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        let delayFactor = await operatorContract.delayFactor();        

        let expectedDelayFactor = web3.utils.toBN(10000000000000000);
        assert.isTrue(expectedDelayFactor.eq(delayFactor));
    });

    it("should correctly evaluate delay factor at the first submission block", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        mineBlocks(await operatorContract.relayEntryGenerationTime());

        let delayFactor = await operatorContract.delayFactor();

        let expectedDelayFactor = web3.utils.toBN(10000000000000000);
        assert.isTrue(expectedDelayFactor.eq(delayFactor));
    });

    it("should correctly evaluate delay factor at the second submission block", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        mineBlocks((await operatorContract.relayEntryGenerationTime()).addn(1));

        let delayFactor = await operatorContract.delayFactor();

        let expectedDelayFactor = web3.utils.toBN('9896104600694443');
        assert.isTrue(expectedDelayFactor.eq(delayFactor));
    });

    it("should correctly evaluate delay factor in the last block before timeout", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        mineBlocks((await operatorContract.relayEntryTimeout()).subn(1));

        let delayFactor = await operatorContract.delayFactor();        

        let expectedDelayFactor = web3.utils.toBN('271267361111');
        assert.isTrue(expectedDelayFactor.eq(delayFactor));        
    });
});
