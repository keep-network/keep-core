import mineBlocks from '../helpers/mineBlocks';
import { duration, increaseTimeTo } from '../helpers/increaseTime';
import latestTime from '../helpers/latestTime';
import expectThrowWithMessage from '../helpers/expectThrowWithMessage'
import grantTokens from '../helpers/grantTokens';
import { createSnapshot, restoreSnapshot } from '../helpers/snapshot'
import delegateStakeFromGrant from '../helpers/delegateStakeFromGrant'

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

const KeepToken = artifacts.require('./KeepToken.sol');
const TokenStaking = artifacts.require('./TokenStaking.sol');
const TokenGrant = artifacts.require('./TokenGrant.sol');
const Registry = artifacts.require("./Registry.sol");

contract('TokenGrant/Reassignment', function(accounts) {

    let tokenContract, registryContract, grantContract, stakingContract;

    const grantManager = accounts[0],
          grantee = accounts[1],
          newGrantee = accounts[2],
          thirdParty = accounts[3];

    let grantId;
    let grantStart;

    const grantAmount = web3.utils.toBN(1000000000),
          grantVestingDuration = duration.days(60),
          grantCliff = duration.days(10),
          grantRevocable = false;

    const initializationPeriod = 10;
    const undelegationPeriod = 30;

    before(async () => {
        tokenContract = await KeepToken.new({from: grantManager});
        registryContract = await Registry.new();
        stakingContract = await TokenStaking.new(
            tokenContract.address, 
            registryContract.address, 
            initializationPeriod, 
            undelegationPeriod
        );
        grantContract = await TokenGrant.new(tokenContract.address);

        await grantContract.authorizeStakingContract(
            stakingContract.address,
            {from: grantManager}
        );

        grantStart = await latestTime();

        // Grant tokens
        grantId = await grantTokens(
            grantContract, 
            tokenContract, 
            grantAmount, 
            grantManager, 
            grantee, 
            grantVestingDuration, 
            grantStart, 
            grantCliff, 
            grantRevocable
        );
    });

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("requestGranteeReassignment", async () => {
        it("should let the grantee request reassignment", async () => {
            await grantContract.requestGranteeReassignment(
                grantId,
                newGrantee,
                {from: grantee}
            );
            let requested = await grantContract.getPendingReassignment(grantId);
            expect(requested).to.equal(newGrantee);
        })

        it("should not let the manager request reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.requestGranteeReassignment(
                    grantId,
                    newGrantee,
                    {from: grantManager}
                ),
                "Only grantee may request reassignment"
            );
        })

        it("should not let the new grantee request reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.requestGranteeReassignment(
                    grantId,
                    newGrantee,
                    {from: newGrantee}
                ),
                "Only grantee may request reassignment"
            );
        })

        it("should not let a third party request reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.requestGranteeReassignment(
                    grantId,
                    newGrantee,
                    {from: thirdParty}
                ),
                "Only grantee may request reassignment"
            );
        })

        it("should require the new grantee to be different", async () => {
            await expectThrowWithMessage(
                grantContract.requestGranteeReassignment(
                    grantId,
                    grantee,
                    {from: grantee}
                ),
                "New grantee must be different"
            );
        })

        it("should require the new grantee to be not null", async () => {
            await expectThrowWithMessage(
                grantContract.requestGranteeReassignment(
                    grantId,
                    "0x0000000000000000000000000000000000000000",
                    {from: grantee}
                ),
                "Must specify new grantee address"
            );
        })
    })

    describe("confirmGranteeReassignment", async () => {
        beforeEach(async () => {
            await grantContract.requestGranteeReassignment(
                grantId,
                newGrantee,
                {from: grantee}
            );
        })

        it("should let the grant manager confirm reassignment", async () => {
            await grantContract.confirmGranteeReassignment(
                grantId,
                {from: grantManager}
            );
            let hasRequest = await grantContract.hasPendingReassignment(grantId);
            expect(hasRequest).to.equal(false);
            let changedGrant = await grantContract.getGrant(grantId);
            let changedGrantee = changedGrant[5];
            expect(changedGrantee).to.equal(newGrantee);
        })

        it("should not let the old grantee confirm reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.confirmGranteeReassignment(
                    grantId,
                    {from: grantee}
                ),
                "Only grant manager may confirm grantee reassignment"
            );
        })

        it("should not let the new grantee confirm reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.confirmGranteeReassignment(
                    grantId,
                    {from: newGrantee}
                ),
                "Only grant manager may confirm grantee reassignment"
            );
        })

        it("should not let a third party confirm reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.confirmGranteeReassignment(
                    grantId,
                    {from: thirdParty}
                ),
                "Only grant manager may confirm grantee reassignment"
            );
        })

        it("should require a request to be pending", async () => {
            let secondGrantId = await grantTokens(
                grantContract, 
                tokenContract, 
                grantAmount, 
                grantManager, 
                newGrantee, 
                grantVestingDuration, 
                grantStart, 
                grantCliff, 
                grantRevocable
            );
            await expectThrowWithMessage(
                grantContract.confirmGranteeReassignment(
                    secondGrantId,
                    {from: grantManager}
                ),
                "No reassignment requested"
            );
        })
    })

    describe("refuseGranteeReassignment", async () => {
        beforeEach(async () => {
            await grantContract.requestGranteeReassignment(
                grantId,
                newGrantee,
                {from: grantee}
            );
        })

        it("should let the grant manager refuse reassignment", async () => {
            await grantContract.refuseGranteeReassignment(
                grantId,
                {from: grantManager}
            );
            let hasRequest = await grantContract.hasPendingReassignment(grantId);
            expect(hasRequest).to.be.false;
            let changedGrant = await grantContract.getGrant(grantId);
            let unchangedGrantee = changedGrant[5];
            expect(unchangedGrantee).to.equal(grantee);
        })

        it("should not let the old grantee refuse reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.refuseGranteeReassignment(
                    grantId,
                    {from: grantee}
                ),
                "Only grant manager may refuse grantee reassignment"
            );
        })

        it("should not let the new grantee refuse reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.refuseGranteeReassignment(
                    grantId,
                    {from: newGrantee}
                ),
                "Only grant manager may refuse grantee reassignment"
            );
        })

        it("should not let a third party refuse reassignment", async () => {
            await expectThrowWithMessage(
                grantContract.refuseGranteeReassignment(
                    grantId,
                    {from: thirdParty}
                ),
                "Only grant manager may refuse grantee reassignment"
            );
        })

        it("should require a request to be pending", async () => {
            let secondGrantId = await grantTokens(
                grantContract, 
                tokenContract, 
                grantAmount, 
                grantManager, 
                newGrantee, 
                grantVestingDuration, 
                grantStart, 
                grantCliff, 
                grantRevocable
            );
            await expectThrowWithMessage(
                grantContract.refuseGranteeReassignment(
                    secondGrantId,
                    {from: grantManager}
                ),
                "No reassignment requested"
            );
        })
    })
});
