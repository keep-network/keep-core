const { initContracts } = require('./helpers/initContracts')
const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")
const { time } = require("@openzeppelin/test-helpers")
const crypto = require("crypto")
const stakeDelegate = require('./helpers/stakeDelegate')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect
const assert = chai.assert

describe('BeaconBackportRewards', () => {

    let token, stakingContract, operatorContract, serviceContract, rewards,
        groupSize, minimumStake,
        group0, group1, group2, group3,
        owner = accounts[0],
        operator1 = accounts[2],
        operator2 = accounts[3],
        operator3 = accounts[4],
        beneficiary1 = accounts[5],
        beneficiary2 = accounts[6],
        beneficiary3 = accounts[7],
        excessRecipient = accounts[8]

    before(async () => {
        let contracts = await initContracts(
            contract.fromArtifact('TokenStaking'),
            contract.fromArtifact('KeepRandomBeaconService'),
            contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
            contract.fromArtifact('KeepRandomBeaconOperatorBeaconRewardsStub')
        )
        const termLength = 100
        const totalRewards = 9000

        token = contracts.token
        stakingContract = contracts.stakingContract
        operatorContract = contracts.operatorContract
        serviceContract = contracts.serviceContract

        groupSize = await operatorContract.groupSize()
        minimumStake = await stakingContract.minimumStake()

        await stakeDelegate(stakingContract, token, owner, operator1, beneficiary1, operator1, minimumStake)
        await stakeDelegate(stakingContract, token, owner, operator2, beneficiary2, operator2, minimumStake)
        await stakeDelegate(stakingContract, token, owner, operator3, beneficiary3, operator3, minimumStake)

        group0 = crypto.randomBytes(128)
        group1 = crypto.randomBytes(128)
        group2 = crypto.randomBytes(128)
        group3 = crypto.randomBytes(128)

        await operatorContract.registerNewGroup(group0, [operator1, operator2, operator2])
        await operatorContract.registerNewGroup(group1, [operator1, operator2, operator2])
        await operatorContract.registerNewGroup(group2, [operator1, operator2, operator2])
        await operatorContract.registerNewGroup(group3, [operator1, operator2, operator2])

        const initiationTime = await time.latest()

        await time.increase(250)

        rewards = await contract.fromArtifact('BeaconBackportRewardsStub').new(
            token.address,
            initiationTime,
            [50, 100],
            operatorContract.address,
            stakingContract.address,
            [2, 3], // groups 0~2 in first interval, 3 in second
            [1],
            excessRecipient
        )
        await token.approveAndCall(
            rewards.address,
            totalRewards,
            "0x0",
            { from: owner }
        )

        // make all groups expire
        let blockN = await time.latestBlock()
        await time.advanceBlockTo(blockN.addn(15))
        await operatorContract.expireOldGroups()
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    it("should have 4 groups", async () => {
        let count = await rewards.getKeepCount();
        expect(count).to.eq.BN(4);
    })

    it("should have 3 as the last eligible group", async () => {
        let count = await rewards.lastEligibleGroup();
        expect(count).to.eq.BN(3);
    })

    it("should have 1 as the last interval", async () => {
        expect(await rewards.lastInterval()).to.eq.BN(1);
    })

    it("should exclude group 1", async () => {
        assert.isTrue(await rewards.isExcluded(1), "group 1 not excluded")
    })

    it("should recognize groups 0~3", async () => {
        let recognized0 = await rewards.recognizedByFactory(0);
        let recognized1 = await rewards.recognizedByFactory(1);
        let recognized2 = await rewards.recognizedByFactory(2);
        let recognized3 = await rewards.recognizedByFactory(3);
        let recognized4 = await rewards.recognizedByFactory(4);

        assert.isTrue(recognized0, "group 0 not recognized")
        assert.isTrue(recognized1, "group 1 not recognized")
        assert.isTrue(recognized2, "group 2 not recognized")
        assert.isTrue(recognized3, "group 3 not recognized")
        assert.isFalse(recognized4, "group 4 falsely recognized")
    })

    it("should have groups 0, 2 and 3 eligible", async () => {
        let eligible0 = await rewards.eligibleForReward(0);
        let eligible1 = await rewards.eligibleForReward(1);
        let eligible2 = await rewards.eligibleForReward(2);
        let eligible3 = await rewards.eligibleForReward(3);

        assert.isTrue(eligible0, "group 0 ineligible")
        assert.isFalse(eligible1, "group 1 eligible")
        assert.isTrue(eligible2, "group 2 ineligible")
        assert.isTrue(eligible3, "group 3 ineligible")
    })

    it("should recognize group 1 as terminated", async () => {
        let terminated0 = await rewards.eligibleButTerminatedWithUint(0);
        let terminated1 = await rewards.eligibleButTerminatedWithUint(1);
        let terminated2 = await rewards.eligibleButTerminatedWithUint(2);
        let terminated3 = await rewards.eligibleButTerminatedWithUint(3);

        assert.isFalse(terminated0, "group 0 falsely terminated")
        assert.isTrue(terminated1, "group 1 not terminated")
        assert.isFalse(terminated2, "group 2 falsely terminated")
        assert.isFalse(terminated3, "group 3 falsely terminated")
    })

    it("should register 3 groups in the first interval", async () => {
        let count = await rewards.findEndpoint(await rewards.endOf(0))
        expect(count).to.eq.BN(3)
    })

    it("should register 1 group in the second interval", async () => {
        let count = await rewards.findEndpoint(await rewards.endOf(1))
        expect(count).to.eq.BN(4)
    })

    it("should receive rewards for groups 0, 2 and 3", async () => {
        // 4500 allocated to the first interval:
        //   1500 allocated to group 0
        //   1500 not allocated to group 1 because it's terminated
        //   1500 allocated to group 2
        // 4500 allocated to the second interval:
        //   4500 allocated to group 3

        // 1500 allocated to group 0
        await rewards.receiveReward(0);
        expect(await token.balanceOf(beneficiary1)).to.eq.BN(500)
        expect(await token.balanceOf(beneficiary2)).to.eq.BN(1000)

        // 1500 allocated to group 2
        await rewards.receiveReward(2);
        expect(await token.balanceOf(beneficiary1)).to.eq.BN(1000) // 500+500
        expect(await token.balanceOf(beneficiary2)).to.eq.BN(2000) // 1000+1000

        // 4500 allocated to group 3
        await rewards.receiveReward(3);
        expect(await token.balanceOf(beneficiary1)).to.eq.BN(2500) // 1000+1500
        expect(await token.balanceOf(beneficiary2)).to.eq.BN(5000) // 2000+3000
    })

    it("should withdraw excess rewards from terminated group 1", async () => {
        await rewards.allocateRewards(1);
        await rewards.reportTermination(1);
        await rewards.withdrawExcess();

        let balance = await token.balanceOf(excessRecipient);
        expect(balance).to.eq.BN(1500)
    })
})
