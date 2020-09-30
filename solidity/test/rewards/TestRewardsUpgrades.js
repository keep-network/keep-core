const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { expectRevert, time } = require("@openzeppelin/test-helpers")

const KeepToken = contract.fromArtifact('KeepToken')

const RewardsStub = contract.fromArtifact('RewardsStub');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect
const assert = chai.assert

describe('Rewards/Upgrading', () => {
    const notFunder = accounts[1]
    const funder = accounts[9]

    const termLength = 100

    const minimumIntervalKeeps = 2
    const initiationTime = 1000
    const timestamps = [
        // No keeps in interval 0
        1100, // 1 keep in interval 1
        1200, 1201 // 2 keeps in interval 2
    ]
    const intervalWeights = [40, 50, 100]
    const totalRewards = 1000000

    const expectedAllocationsWithoutUpgrade = [0, 250000, 750000]
    const expectedAllocationsWithUpgrade = [0, 150000, 300000]
    const expectedNewRewardsTokens = [400000, 550000, 550000]

    let rewards
    let token
    let newRewards

    async function fund(amount) {
        await token.approveAndCall(
            rewards.address,
            amount,
            "0x0",
            { from: funder }
        )
    }

    before(async () => {
        token = await KeepToken.new({ from: funder })
        rewards = await RewardsStub.new(
            token.address,
            minimumIntervalKeeps,
            initiationTime,
            intervalWeights,
            timestamps,
            termLength,
            { from: funder }
        )
        await fund(totalRewards)

        newRewards = await RewardsStub.new(
            token.address,
            minimumIntervalKeeps,
            initiationTime,
            intervalWeights,
            timestamps,
            termLength,
            { from: funder }
        )
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("setNewRewards", async () => {
        it("can be called by the owner", async () => {
            await rewards.setNewRewards(newRewards.address, { from: funder })
            let reportedNewRewards = await rewards.newRewards()
            expect(reportedNewRewards).to.equal(newRewards.address)
        })

        it("can't be called by other accounts", async () => {
            await expectRevert(
                rewards.setNewRewards(newRewards.address, { from: notFunder }),
                "Ownable: caller is not the owner."
            )
        })
    })

    describe("allocateRewards", async () => {
        it("allocates the reward for each interval without upgrade", async () => {
            let expectedAllocations = expectedAllocationsWithoutUpgrade
            for (let i = 0; i < expectedAllocations.length; i++) {
                await rewards.allocateRewards(i)
                let allocation = await rewards.getAllocatedRewards(i)
                expect(allocation.toNumber()).to.equal(expectedAllocations[i])
            }
        })

        it("allocates the reward for each interval with upgrade", async () => {
            let expectedAllocations = expectedAllocationsWithUpgrade
            await rewards.setNewRewards(newRewards.address, { from: funder })
            for (let i = 0; i < expectedAllocations.length; i++) {
                await rewards.allocateRewards(i)
                let allocation = await rewards.getAllocatedRewards(i)
                expect(allocation.toNumber()).to.equal(expectedAllocations[i])
                let tokensInNewRewards = await newRewards.unallocatedRewards()
                expect(tokensInNewRewards.toNumber()).to.equal(expectedNewRewardsTokens[i])
            }
        })
    })

    describe("reportTermination", async () => {
        it("unallocates rewards allocated to terminated keeps when not upgraded", async () => {
            await rewards.setCloseTime(timestamps[1])
            await rewards.receiveReward(1, { from: notFunder }) // allocate rewards

            await rewards.terminate(2)
            let preUnallocated = await rewards.unallocatedRewards()
            await rewards.reportTermination(2)
            let postUnallocated = await rewards.unallocatedRewards()
            expect(postUnallocated.toNumber()).to.equal(
                preUnallocated.toNumber() + 375000
            )
        })

        it("transfers rewards allocated to terminated keeps when upgraded", async () => {
            await rewards.setCloseTime(timestamps[1])
            await rewards.receiveReward(1, { from: notFunder }) // allocate rewards

            await rewards.setNewRewards(newRewards.address, { from: funder })
            await rewards.terminate(2)
            let preUnallocated = await rewards.unallocatedRewards()
            await rewards.reportTermination(2)
            let postUnallocated = await rewards.unallocatedRewards()
            expect(postUnallocated.toNumber()).to.equal(
                preUnallocated.toNumber()
            )
            let inNewRewards = await newRewards.unallocatedRewards()
            expect(inNewRewards.toNumber()).to.equal(375000)
        })
    })
})
