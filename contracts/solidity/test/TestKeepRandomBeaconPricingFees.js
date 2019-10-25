import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import {bls} from './helpers/data';
import mineBlocks from './helpers/mineBlocks';

describe('Keep random beacon pricing', function(accounts) {
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

    it("should correctly evaluate entry verification fee", async () => {
        await serviceContract.setPriceFeedEstimate(200);
        await operatorContract.setEntryVerificationGasEstimate(12);        

        let fees = await serviceContract.entryFeeBreakdown();
        let entryVerificationFee = fees.entryVerificationFee;

        let expectedEntryVerificationFee = 3600; // 200 * 12 * 150%
        assert.equal(expectedEntryVerificationFee, entryVerificationFee);
    });

    it("should correctly evaluate DKG contribution fee", async () => {
        await serviceContract.setPriceFeedEstimate(123);
        await operatorContract.setDkgGasEstimate(13);

        let fees = await serviceContract.entryFeeBreakdown();
        let dkgContributionFee = fees.dkgContributionFee;

        let expectedDkgContributionFee = 159; // 123 * 13 * 10%
        assert.equal(expectedDkgContributionFee, dkgContributionFee);
    });

    it("should correctly evaluate callback fee", async function() {
        await serviceContract.setPriceFeedEstimate(160);

        let callbackGas = 1091;

        let callbackFee = await serviceContract.callbackFee(callbackGas);
        
        let expectedCallbackFee = 261840; // 1091 * 160 * 150%
        assert.equal(expectedCallbackFee, callbackFee);
    });

    it("should correctly evaluate entry fee estimate", async () => {
        await serviceContract.setPriceFeedEstimate(200);
        await operatorContract.setEntryVerificationGasEstimate(12); 
        await operatorContract.setDkgGasEstimate(14); 
        await operatorContract.setGroupSize(13);
        await operatorContract.setGroupMemberBaseReward(3);

        let callbackGas = 7;

        let entryFeeEstimate = await serviceContract.entryFeeEstimate(
            callbackGas
        );

        // entry verification fee = 12 * 200 * 150% = 3600
        // dkg contribution fee = 14 * 200 * 10% = 280
        // group profit fee = 13 * 3 = 39
        // callback fee = 7 * 200 * 150% = 2100
        // entry fee = 3600 + 280 + 39 + 2100 = 6019
        let expectedEntryFeeEstimate = 6019;
        assert.equal(expectedEntryFeeEstimate, entryFeeEstimate)
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

        let expectedDelayFactor = web3.utils.toBN(8711111111111110);
        assert.isTrue(expectedDelayFactor.eq(delayFactor));
    });

    it("should correctly evaluate delay factor in the last block before timeout", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        mineBlocks((await operatorContract.relayEntryTimeout()).subn(1));

        let delayFactor = await operatorContract.delayFactor();        

        let expectedDelayFactor = web3.utils.toBN(44444444444444);
        assert.isTrue(expectedDelayFactor.eq(delayFactor));        
    });
});
