const { initContracts } = require('./helpers/initContracts')
const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
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
        requestor = accounts[1],
        operator1 = accounts[2],
        operator2 = accounts[3],
        operator3 = accounts[4],
        beneficiary1 = accounts[5],
        beneficiary2 = accounts[6],
        beneficiary3 = accounts[7]

    before(async () => {
        let contracts = await initContracts(
            contract.fromArtifact('TokenStaking'),
            contract.fromArtifact('KeepRandomBeaconService'),
            contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
            contract.fromArtifact('KeepRandomBeaconOperatorBeaconRewardsStub')
        )
        const termLength = 10
        const totalRewards = 6000

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

        rewards = await contract.fromArtifact('BeaconBackportRewardsStub').new(
            termLength,
            token.address,
            initiationTime,
            operatorContract.address,
            stakingContract.address,
            2, // groups 0, 1 and 2 eligible
            [1]
        )
        await token.approveAndCall(
            rewards.address,
            totalRewards,
            "0x0",
            { from: owner }
        )

        await time.increase(15)

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

    function bytes32(byte) {
        return "0x" + ("00" * 31) + byte
    }

    it("should have 2 groups", async () => {
        let count = await rewards.getKeepCount();
        expect(count).to.eq.BN(2);
    })

    it("should have 2 as the last eligible group", async () => {
        let count = await rewards.lastEligibleGroup();
        expect(count).to.eq.BN(2);
    })

    it("should exclude group 1", async () => {
        assert.isTrue(await rewards.isExcluded(1), "group 1 not excluded")
    })

    it("should recognize groups 0 and 2", async () => {
        let recognized0 = await rewards.recognizedByFactory(0);
        let recognized1 = await rewards.recognizedByFactory(1);
        let recognized2 = await rewards.recognizedByFactory(2);
        let recognized3 = await rewards.recognizedByFactory(3);

        assert.isTrue(recognized0, "group 0 not recognized")
        assert.isFalse(recognized1, "group 1 falsely recognized")
        assert.isTrue(recognized2, "group 2 not recognized")
        assert.isFalse(recognized3, "group 3 falsely recognized")
    })

    it("should have groups 0 and 2 eligible", async () => {
        let eligible0 = await rewards.eligibleForReward(0);
        let eligible1 = await rewards.eligibleForReward(1);
        let eligible2 = await rewards.eligibleForReward(2);
        let eligible3 = await rewards.eligibleForReward(3);

        assert.isTrue(eligible0, "group 0 ineligible")
        assert.isFalse(eligible1, "group 1 eligible")
        assert.isTrue(eligible2, "group 2 ineligible")
        assert.isFalse(eligible3, "group 3 eligible")
    })

    it("should recognize no group as terminated", async () => {
        let terminated0 = await rewards.eligibleButTerminatedWithUint(0);
        let terminated1 = await rewards.eligibleButTerminatedWithUint(1);
        let terminated2 = await rewards.eligibleButTerminatedWithUint(2);

        assert.isFalse(terminated0, "group 0 falsely terminated")
        assert.isFalse(terminated1, "group 1 falsely terminated")
        assert.isFalse(terminated2, "group 2 falsely terminated")
    })

    it("should register 2 keeps in the first interval", async () => {
        let count = await rewards.findEndpoint(await rewards.endOf(0))
        expect(count).to.eq.BN(2)
    })

    it("should receive rewards for groups 0 and 2", async () => {
        await rewards.receiveReward(0);

        let balance1 = await token.balanceOf(beneficiary1);
        expect(balance1).to.eq.BN(1000)

        await rewards.receiveReward(2);
        let balance2 = await token.balanceOf(beneficiary2);
        expect(balance2).to.eq.BN(4000)
    })
})
