const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const stakeAndGenesis = require('../helpers/stakeAndGenesis')
const blsData = require("../helpers/data.js")
const {initContracts} = require('../helpers/initContracts')
const assert = require('chai').assert
const {contract, accounts, web3} = require("@openzeppelin/test-environment")

describe('KeepRandomBeaconService/PricingDkg', () => {
    let serviceContract;
    let operatorContract
    let groupCreationFee;

    before(async () => {
        let contracts = await initContracts(
          contract.fromArtifact('TokenStaking'),
          contract.fromArtifact('KeepRandomBeaconService'),
          contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
          contract.fromArtifact('KeepRandomBeaconOperatorPricingDKGStub')
        );

        serviceContract = contracts.serviceContract;
        operatorContract = contracts.operatorContract;

        await stakeAndGenesis(accounts, contracts);    

        groupCreationFee = await operatorContract.groupCreationFee();
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
        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

        let contractBalance = await web3.eth.getBalance(serviceContract.address);

        let insufficientPoolFunds = web3.utils.toBN(groupCreationFee)
          .sub(web3.utils.toBN(contractBalance))
          .sub(web3.utils.toBN(1));
        
        await serviceContract.fundDkgFeePool({value: insufficientPoolFunds});

        await operatorContract.relayEntry(blsData.groupSignature);
        
        assert.isFalse(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should not start"
        );
    });

    it("should trigger new group selection when there are sufficient funds in the " +
       "DKG fee pool", async () => {
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0);
        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});

        let sufficientPoolFunds = web3.utils.toBN(groupCreationFee);
        
        await serviceContract.fundDkgFeePool({value: sufficientPoolFunds});

        await operatorContract.relayEntry(blsData.groupSignature);
        
        assert.isTrue(
            await operatorContract.isGroupSelectionInProgress(), 
            "new group selection should be started"
        );
    });

    it("should not trigger group selection while one is in progress", async () => {
        await serviceContract.fundDkgFeePool({value: 3 * groupCreationFee});
  
        let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});
        await operatorContract.relayEntry(blsData.groupSignature);

        assert.isTrue(
          await operatorContract.isGroupSelectionInProgress(),
          "new group selection should be started"
        );
  
        let startBlock = await operatorContract.getTicketSubmissionStartBlock();

        await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate, from: accounts[0]});
        let contractBalance = web3.utils.toBN(await web3.eth.getBalance(serviceContract.address));
        await operatorContract.relayEntry(blsData.nextGroupSignature);

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
