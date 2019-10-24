import {initContracts} from './helpers/initContracts';
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot";
import {bls} from './helpers/data';
import stakeDelegate from './helpers/stakeDelegate';
import runGenesisGroupSelection from './helpers/runGenesisGroupSelection';

contract('KeepRandomBeaconService', function(accounts) {

    const groupSize = 20;
    const minimumStake = web3.utils.toBN(200000);

    const operator1StakingWeight = 2000;
    const operator2StakingWeight = 2000;
    const operator3StakingWeight = 3000;
    
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
    
        operatorContract.setGroupSize(groupSize);
        operatorContract.setMinimumStake(minimumStake);

        let stakingContract = contracts.stakingContract;
        let token = contracts.token;
  
        let owner = accounts[0];
        let operator1 = accounts[1];
        let operator2 = accounts[2];
        let operator3 = accounts[3];

        await stakeDelegate(stakingContract, token, owner, operator1, operator1, minimumStake.mul(web3.utils.toBN(operator1StakingWeight)));
        await stakeDelegate(stakingContract, token, owner, operator2, operator2, minimumStake.mul(web3.utils.toBN(operator2StakingWeight)));
        await stakeDelegate(stakingContract, token, owner, operator3, operator3, minimumStake.mul(web3.utils.toBN(operator3StakingWeight)));

        await runGenesisGroupSelection(operatorContract, operator1, operator2, operator3);    

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

    it("should not trigger new group selection when there are not enough " +
       "funds in the DKG fee pool", async function() { 
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let notEnough = web3.utils.toBN(dkgPayment)
          .sub(web3.utils.toBN(contractBalance))
          .sub(web3.utils.toBN(1));
        
        await serviceContract.fundDkgFeePool({value: notEnough});

        await operatorContract.relayEntry(bls.nextGroupSignature);
        
        assert.isFalse(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should not start"
        );
    });

    it("should trigger new group selection when there are enough funds in the " +
       "DKG fee pool", async function() {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let enough = web3.utils.toBN(dkgPayment)
          .sub(web3.utils.toBN(contractBalance));
        
        await serviceContract.fundDkgFeePool({value: enough});

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