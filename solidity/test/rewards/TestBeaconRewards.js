const { initContracts } = require('../helpers/initContracts')
const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
const crypto = require("crypto")
const stakeDelegate = require('../helpers/stakeDelegate')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('BeaconRewards', () => {

    // accounts
    const owner = accounts[0]
    let operators
    let beneficiaries

    // contracts
    let rewardsContract
    let token, stakingContract, operatorContract
    const groupActiveTime = 5 // from KeepRandomBeaconOperatorBeaconRewardsStub
    const relayEntryTimeout = 10 // from KeepRandomBeaconOperatorBeaconRewardsStub

    // system parameters
    const tokenDecimalMultiplier = web3.utils.toBN(10).pow(web3.utils.toBN(18))
    // 1,000,000,000 - total KEEP supply
    //   200,000,000 - 20% of the total supply goes to staker rewards
    //    20,000,000 - 10% of staker rewards goes to the random beacon stakers
    //    19,800,000 - 99% of staker rewards goes to the random beacon stakers 
    //                 operating after 2020-09-24
    const totalBeaconRewards = web3.utils.toBN(19800000).mul(tokenDecimalMultiplier)

    const groupSize = 64

    const expectedKeepAllocations = [
        792000, 1520640, 1748736, 1888635, 2077498, 1765874,
        1500993, 1275844, 1084467, 921797, 783528, 665998,
        566099, 481184, 409006, 347655, 295507, 251181,
        213504, 181478, 154257, 131118, 111450, 94733
    ]

    before(async() => {
        let contracts = await initContracts(
            contract.fromArtifact('TokenStaking'),
            contract.fromArtifact('KeepRandomBeaconService'),
            contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
            contract.fromArtifact('KeepRandomBeaconOperatorBeaconRewardsStub')
        )

        token = contracts.token
        stakingContract = contracts.stakingContract
        operatorContract = contracts.operatorContract

        rewardsContract = await contract.fromArtifact('BeaconRewards').new(
            token.address,
            operatorContract.address,
            stakingContract.address,
            { from: owner }
        )

        await token.approveAndCall(
            rewardsContract.address,
            totalBeaconRewards,
            "0x0",
            { from: owner }
        )
        await rewardsContract.markAsFunded({from: owner})

        // create 64 operators and beneficiaries, delegate stake for them
        const minimumStake = await stakingContract.minimumStake()
        operators = []
        beneficiaries = []
        for (i = 0; i < groupSize; i++) {
            const operator = accounts[i]
            const beneficiary = accounts[groupSize + i]
            const authorizer = operator

            operators.push(operator)
            beneficiaries.push(beneficiary)
            await stakeDelegate(
                stakingContract, 
                token, 
                owner, 
                operator, 
                beneficiary, 
                authorizer, 
                minimumStake
            )
        }
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("interval allocations", async() => {
        it("should equal expected ones having two groups created per interval", async() => {
            for (let i = 0; i < 24; i++) {
                const startOf = await rewardsContract.startOf(i)

                // 2 groups created in an interval
                await registerNewGroup(startOf.addn(1))
                await registerNewGroup(startOf.addn(2))

                await timeJumpToEndOfInterval(i)
                await rewardsContract.allocateRewards(i)
                
                const allocated = await rewardsContract.getAllocatedRewards(i)
                const allocatedKeep = allocated.div(tokenDecimalMultiplier)

                assertKeepIntervalAllocations(allocatedKeep, expectedKeepAllocations[i])
            }
        })

        it("should equal expected ones having more than two groups created per interval", async () => {
            for (let i = 0; i < 24; i++) {
                const startOf = await rewardsContract.startOf(i)

                // 5 groups created in an interval
                for (let j = 0; j < 5; j++) {
                    await registerNewGroup(startOf.addn(j+1))
                }

                await timeJumpToEndOfInterval(i)
                await rewardsContract.allocateRewards(i)
                
                const allocated = await rewardsContract.getAllocatedRewards(i)
                const allocatedKeep = allocated.div(tokenDecimalMultiplier)

                assertKeepIntervalAllocations(allocatedKeep, expectedKeepAllocations[i])
            }
        })

        it("should equal expected ones having just one group created per interval", async () => {
            // 1st interval expected allocation: 19,800,000 * 4% = 792,000
            // 1st interval adjusted: 792,000 / 2 = 396,000
            // Remaining pool: 19,800,000 - 396,000 = 19,404,000
            // 2nd interval expected allocation: 19,404,000 * 8% = 1,552,320
            // 2nd interval adjusted: 776160
            // etc.
            const adjustedKeepAllocations = [
                396000, 776160, 931392, 1061787, 1247600, 1154030,
                1067477, 987417, 913360, 844858, 781494, 722882,
                668666, 618516, 572127, 529218, 489526, 452812,
                418851, 387437, 358379, 331501, 306638, 283640,
            ]

            for (let i = 0; i < 24; i++) {
                const startOf = await rewardsContract.startOf(i)
                // one group created in an interval
                await registerNewGroup(startOf.addn(1))

                await timeJumpToEndOfInterval(i)
                await rewardsContract.allocateRewards(i)

                const allocated = await rewardsContract.getAllocatedRewards(i)                                
                const allocatedKeep = allocated.div(tokenDecimalMultiplier)

                assertKeepIntervalAllocations(allocatedKeep, adjustedKeepAllocations[i])
            }
        })
    })

    describe("rewards withdrawal", async () => {
        it("should be possible for stale groups", async () => {
            const startOf = await rewardsContract.startOf(0)
            await registerNewGroup(startOf.addn(1))
            await expireAllGroups()

            const isEligible = await rewardsContract.eligibleForReward(0)
            expect(isEligible).to.be.true

            await timeJumpToEndOfInterval(0)
            await rewardsContract.receiveReward(0)
            // ok, no revert
        })

        it("should not be possible for non-stale groups", async () => {
            const startOf = await rewardsContract.startOf(0)
            await registerNewGroup(startOf.addn(1))

            const isEligible = await rewardsContract.eligibleForReward(0)
            expect(isEligible).to.be.false

            await timeJumpToEndOfInterval(0)
            await expectRevert(
                rewardsContract.receiveReward(0),
                "Keep is not closed"
            )
        })

        it("should not be possible for terminated groups", async () => {
            const startOf = await rewardsContract.startOf(0)
            await registerNewGroup(startOf.addn(1))
            await operatorContract.terminateGroup(0)
            
            const isEligible = await rewardsContract.eligibleForReward(0)
            expect(isEligible).to.be.false

            await timeJumpToEndOfInterval(0)
            await expectRevert(
                rewardsContract.receiveReward(0),
                "Keep is not closed"
            )
        })

        it("should not count terminated groups when distributing rewards", async () => {
            const startOf = await rewardsContract.startOf(0)
            await registerNewGroup(startOf.addn(1))
            await registerNewGroup(startOf.addn(2))
            await operatorContract.terminateGroup(1)

            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            await expireAllGroups()
            await rewardsContract.receiveReward(0)
            // two groups but one of them is terminated and does not count here 
            // each beneficiary receives 792,000 / 2 / 64 = 6,187.5 KEEP => ~6,188
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(6188))

            // the remaining 396,000 stays in unallocated rewards but the fact
            // it terminated needs to be reported to recalculate the unallocated
            // amount
            let unallocated = await rewardsContract.unallocatedRewards()
            let unallocatedInKeep = unallocated.div(tokenDecimalMultiplier)
            expect(unallocatedInKeep).to.eq.BN(19008000)

            await rewardsContract.reportTermination(1)
            unallocated = await rewardsContract.unallocatedRewards()
            unallocatedInKeep = unallocated.div(tokenDecimalMultiplier)
            expect(unallocatedInKeep).to.eq.BN(19404000)
        })

        it("should not count a batch of terminated groups when distributing rewards", async () => {
            const startOf = await rewardsContract.startOf(0)
            await registerNewGroup(startOf.addn(1))
            await registerNewGroup(startOf.addn(2))
            await registerNewGroup(startOf.addn(3))

            const terminatedGroups = [1, 2]
            await operatorContract.terminateGroup(terminatedGroups[0])
            await operatorContract.terminateGroup(terminatedGroups[1])

            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            await expireAllGroups()
            await rewardsContract.receiveReward(0)
            // three groups but two of them were terminated and do not count here 
            // each beneficiary receives 792,000 / 3 / 64 = 4,125 KEEP
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(4125))

            // the remaining unallocated rewards pool has 19,008,000 KEEP
            // the remaining 528,000 stays in unallocated rewards but the fact
            // two keeps were terminated needs to be reported to recalculate the 
            // unallocated amount
            // unallocated amount: 19,008,000 + 528,000 = 19,536,000
            await rewardsContract.methods['reportTerminations(uint256[])'](
                terminatedGroups
            )
            unallocated = await rewardsContract.unallocatedRewards()
            unallocatedInKeep = unallocated.div(tokenDecimalMultiplier)
            expect(unallocatedInKeep).to.eq.BN(19536000)
        })

        it("should correctly distribute rewards to beneficiaries", async () => {
            let startOf = await rewardsContract.startOf(0)
            // 2 groups in the first interval, 792,000 KEEP to distribute
            // between 64 beneficiaries.
            await registerNewGroup(startOf.addn(1))
            await registerNewGroup(startOf.addn(2))
            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            await expireAllGroups()
            await rewardsContract.receiveReward(0)
            // each beneficiary receives 792,000 / 2 / 64 = 6,187.5 KEEP => ~6,188
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(6188))
            await rewardsContract.receiveReward(1)
            // each beneficiary receives 792,000 / 2 / 64 = 6187.5 KEEP
            // they should have 6,187.5 + 6,187.5 = 12,375 KEEP now
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(12375))

            // 1 group in the second interval, 760,320 KEEP to distribute
            // between 64 beneficiaries
            startOf = await rewardsContract.startOf(1)
            await registerNewGroup(startOf.addn(1))
            await timeJumpToEndOfInterval(1)
            await rewardsContract.allocateRewards(1)

            await expireAllGroups()
            await rewardsContract.receiveReward(2)
            // each beneficiary receives 760,320 / 64 = 11,880 KEEP
            // they should have 12,375 + 11,880 = 23,760 KEEP now
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(24255))
        })

        it("should correctly distribute rewards in batch", async () => {
            let startOf = await rewardsContract.startOf(0)
            // 2 groups in the first interval, 792,000 KEEP to distribute
            // between 64 beneficiaries.
            await registerNewGroup(startOf.addn(1))
            await registerNewGroup(startOf.addn(2))

            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            await expireAllGroups()
            const groupsReceivingRewards = [0, 1]
            await rewardsContract.receiveRewards(groupsReceivingRewards)
            // each beneficiary receives 792,000 / 2 / 64 = 6187.5 KEEP
            // each beneficiary was in 2 groups
            // each should receive 6,187.5 * 2 = 12,375 KEEP
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(12375))
        })
    })

    async function timeJumpToEndOfInterval(intervalNumber) {
        const endOf = await rewardsContract.endOf(intervalNumber)
        const now = await time.latest()

        if (now.lt(endOf)) {
            await time.increaseTo(endOf.addn(1))
        }
    }

    async function registerNewGroup(creationTimestamp) {
        const groupPublicKey = crypto.randomBytes(128)
        await operatorContract.registerNewGroup(groupPublicKey, operators, creationTimestamp)
    }

    async function expireAllGroups() {
        const currentBlock = await time.latestBlock()

        const groupStalingTime = groupActiveTime + relayEntryTimeout + 1
        await time.advanceBlockTo(currentBlock.addn(groupStalingTime))
        await operatorContract.expireOldGroups()
    }

    async function assertKeepBalanceOfBeneficiaries(expectedBalance) {
        // Solidity is not very good when it comes to floating point precision,
        // we are allowing for ~1 KEEP difference margin between expected and
        // actual value.
        const precision = 1

        for (let i = 0; i < beneficiaries.length; i++) {
            const balance = await token.balanceOf(beneficiaries[i])
            const balanceInKeep = balance.div(tokenDecimalMultiplier)

            expect(balanceInKeep).to.gte.BN(expectedBalance.subn(precision))
            expect(balanceInKeep).to.lte.BN(expectedBalance.addn(precision))
        }
    }

    async function assertKeepIntervalAllocations(actual, expected) {
        const precision = 1
    
        expect(actual).to.gte.BN(expected - precision)
        expect(actual).to.lte.BN(expected + precision)
    }
})
