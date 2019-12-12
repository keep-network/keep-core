import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import {bls} from './helpers/data';

import stakeAndGenesis from './helpers/stakeAndGenesis';

contract('KeepRandomBeaconService', (accounts) => {

    const groupSize = 20;
    const groupThreshold = 11;

    let serviceContract;
    let operatorContract
    let dkgPayment;

    before(async () => {
        let contracts = await initContracts(
          artifacts.require('./KeepToken.sol'),
          artifacts.require('./TokenStaking.sol'),
          artifacts.require('./KeepRandomBeaconService.sol'),
          artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
          artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
        );
        
        serviceContract = contracts.serviceContract;
        operatorContract = contracts.operatorContract;
    
        await operatorContract.setGroupSize(groupSize);
        await operatorContract.setGroupThreshold(groupThreshold);

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
       "funds in the DKG fee pool", async () => { 
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let insufficientPoolFunds = web3.utils.toBN(dkgPayment)
          .sub(web3.utils.toBN(contractBalance))
          .sub(web3.utils.toBN(1));
        
        await serviceContract.fundDkgFeePool({value: insufficientPoolFunds});

        await operatorContract.relayEntry(bls.groupSignature);
        
        assert.isFalse(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should not start"
        );
    });

    it("should trigger new group selection when there are sufficient funds in the " +
       "DKG fee pool", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let sufficientPoolFunds = web3.utils.toBN(dkgPayment)
          .sub(web3.utils.toBN(contractBalance));
        
        await serviceContract.fundDkgFeePool({value: sufficientPoolFunds});

        await operatorContract.relayEntry(bls.groupSignature);
        
        assert.isTrue(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should be started"
        );
    });

    it("should not trigger group selection while one is in progress", async () => {
        await serviceContract.fundDkgFeePool({value: 3 * dkgPayment});
  
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
        await operatorContract.relayEntry(bls.groupSignature);

        assert.isTrue(
          await operatorContract.isGroupSelectionInProgress(),
          "new group selection should be started"
        );
  
        let startBlock = await operatorContract.getTicketSubmissionStartBlock();

        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate});
        let contractBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));
        await operatorContract.relayEntry(bls.nextGroupSignature);

        assert.isTrue(
            web3.utils.toBN(await web3.eth.getBalance(serviceContract.address)).eq(contractBalance),
            "service contract balance should not change"
        )

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