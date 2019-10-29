import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import {bls} from './helpers/data';

import stakeAndGenesis from './helpers/stakeAndGenesis';

describe('Keep random beacon pricing', function(accounts) {

    const groupSize = 20;

    let serviceContract;
    let operatorContract
    let dkgPayment;

    before(async () => {
        let contracts = await initContracts(
          artifacts.require('./KeepToken.sol'),
          artifacts.require('./TokenStaking.sol'),
          artifacts.require('./KeepRandomBeaconService.sol'),
          artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
          artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol'),
          artifacts.require('./KeepRandomBeaconOperatorGroups.sol')
        );
        
        serviceContract = contracts.serviceContract;
        operatorContract = contracts.operatorContract;
    
        await operatorContract.setGroupSize(groupSize);

        await stakeAndGenesis(accounts, contracts);    

        let dkgGasEstimateCost = await operatorContract.dkgGasEstimate();
        let fluctuationMargin = await operatorContract.fluctuationMargin();
        let priceFeedEstimate = await serviceContract.priceFeedEstimate();
        let gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)));
        dkgPayment = dkgGasEstimateCost.mul(gasPriceWithFluctuationMargin);
    });
    
    beforeEach(async () => {
        await createSnapshot()
    });
    
    afterEach(async () => {
        await restoreSnapshot()
    });

    it("should not trigger new group selection when there are not sufficient " +
       "funds in the DKG fee pool", async function() { 
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let insufficientPoolFunds = web3.utils.toBN(dkgPayment)
          .sub(web3.utils.toBN(contractBalance))
          .sub(web3.utils.toBN(1));
        
        await serviceContract.fundDkgFeePool({value: insufficientPoolFunds});

        await operatorContract.relayEntry(bls.nextGroupSignature);
        
        assert.isFalse(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should not start"
        );
    });

    it("should trigger new group selection when there are sufficient funds in the " +
       "DKG fee pool", async function() {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let sufficientPoolFunds = web3.utils.toBN(dkgPayment)
          .sub(web3.utils.toBN(contractBalance));
        
        await serviceContract.fundDkgFeePool({value: sufficientPoolFunds});

        await operatorContract.relayEntry(bls.nextGroupSignature);
        
        assert.isTrue(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should be started"
        );
    });

    it("should not trigger group selection while one is in progress", async function() {
        await serviceContract.fundDkgFeePool({value: 3 * dkgPayment});
  
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});
        await operatorContract.relayEntry(bls.nextGroupSignature);

        assert.isTrue(
          await operatorContract.isGroupSelectionInProgress(),
          "new group selection should be started"
        );
  
        let startBlock = await operatorContract.getTicketSubmissionStartBlock();

        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});
        await operatorContract.relayEntry(bls.nextNextGroupSignature);

        assert.isTrue(
            await operatorContract.isGroupSelectionInProgress(),
            "previous group selection should continue"
        );

        assert.isTrue(
            startBlock.eq(await operatorContract.getTicketSubmissionStartBlock()),
            "new group selection should not be started"
        );
    });
});