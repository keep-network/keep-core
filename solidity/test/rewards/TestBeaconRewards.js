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
    const totalBeaconRewards = web3.utils.toBN(20000000).mul(tokenDecimalMultiplier)

    const groupSize = 64

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
            stakingContract.address
        )

        await token.approveAndCall(
            rewardsContract.address,
            totalBeaconRewards,
            "0x0",
            { from: owner }
        )

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
            const expectedKeepAllocations = [                
                800000,  1536000, 1766400, 1907712, 2098483, 1783710,
                1516154, 1288730, 1095421, 931108,  791441,  672725,
                571816,  486044,  413137,  351166,  298491,  253718,
                215660,  183311,  155814,  132442,  112576,  95689
            ]

            for (let i = 0; i < 24; i++) {
                // 2 groups created in an interval
                await registerNewGroup()
                await registerNewGroup()

                await timeJumpToEndOfInterval(i)
                await rewardsContract.allocateRewards(i)
                
                const allocated = await rewardsContract.getAllocatedRewards(i)                                
                const allocatedKeep = allocated.div(tokenDecimalMultiplier)
                
                expect(allocatedKeep).to.eq.BN(expectedKeepAllocations[i])
            }
        })

        it("should equal expected ones having more than two groups created per interval", async () => {
            const expectedKeepAllocations = [                
                800000,  1536000, 1766400, 1907712, 2098483, 1783710,
                1516154, 1288730, 1095421, 931108,  791441,  672725,
                571816,  486044,  413137,  351166,  298491,  253718,
                215660,  183311,  155814,  132442,  112576,  95689
            ]

            for (let i = 0; i < 24; i++) {
                // 5 groups created in an interval
                for (let j = 0; j < 5; j++) {
                    await registerNewGroup()
                }

                await timeJumpToEndOfInterval(i)
                await rewardsContract.allocateRewards(i)
                
                const allocated = await rewardsContract.getAllocatedRewards(i)                                
                const allocatedKeep = allocated.div(tokenDecimalMultiplier)
                
                expect(allocatedKeep).to.eq.BN(expectedKeepAllocations[i])
            }
        })

        it("should equal expected ones having just one group created per interval", async () => {
            const expectedKeepAllocations = [
                400000,  784000, 940800, 1072512, 1260201, 1165686,
                1078259, 997390, 922586, 853392,  789387,  730183,
                675419,  624763, 577906, 534563,  494470,  457385,
                423081,  391350, 361999, 334849,  309735,  286505
            ]

            for (let i = 0; i < 24; i++) {
                // one group created in an interval
                await registerNewGroup()

                await timeJumpToEndOfInterval(i)
                await rewardsContract.allocateRewards(i)

                const allocated = await rewardsContract.getAllocatedRewards(i)                                
                const allocatedKeep = allocated.div(tokenDecimalMultiplier)
                
                expect(allocatedKeep).to.eq.BN(expectedKeepAllocations[i])
            }
        })
    })

    describe("rewards withdrawal", async () => {
        it("should be possible for stale groups", async () => {
            await registerNewGroup()
            await expireAllGroups()     

            const isEligible = await rewardsContract.eligibleForReward(0)
            expect(isEligible).to.be.true

            await timeJumpToEndOfInterval(0)
            await rewardsContract.receiveReward(0)
            // ok, no revert
        })

        it("should not be possible for non-stale groups", async () => {
            await registerNewGroup()

            const isEligible = await rewardsContract.eligibleForReward(0)
            expect(isEligible).to.be.false

            await timeJumpToEndOfInterval(0)
            await expectRevert(
                rewardsContract.receiveReward(0),
                "Keep is not closed"
            )
        })

        it("should not be possible for terminated groups", async () => {
            await registerNewGroup()
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
            await registerNewGroup()
            await registerNewGroup()
            await operatorContract.terminateGroup(1)

            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            await expireAllGroups()
            await rewardsContract.receiveReward(0)
            // two groups but one of them is terminated and does not count here 
            // each beneficiary receives 800000 / 2 / 64 = 6250 KEEP
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(6250))

            // the remaining 400000 stays in unallocated rewards but the fact
            // it terminated needs to be reported to recalculate the unallocated
            // amount
            let unallocated = await rewardsContract.unallocatedRewards()
            let unallocatedInKeep = unallocated.div(tokenDecimalMultiplier)
            expect(unallocatedInKeep).to.eq.BN(19200000)

            await rewardsContract.reportTermination(1)
            unallocated = await rewardsContract.unallocatedRewards()
            unallocatedInKeep = unallocated.div(tokenDecimalMultiplier)
            expect(unallocatedInKeep).to.eq.BN(19600000)
        })

        it("should correctly distribute rewards to beneficiaries", async () => {
            // 2 groups in the first interval, 800000 KEEP to distribute
            // between 64 beneficiaries.
            await registerNewGroup()
            await registerNewGroup()
            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            await expireAllGroups()
            await rewardsContract.receiveReward(0)
            // each beneficiary receives 800000 / 2 / 64 = 6250 KEEP
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(6250))
            await rewardsContract.receiveReward(1)
            // each beneficiary receives 800000 / 2 / 64 = 6250 KEEP
            // they should have 6250 + 6250 = 12500 KEEP now
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(12500))

            // 1 group in the second interval, 768000 KEEP to distribute
            // between 64 beneficiaries
            await registerNewGroup()
            await timeJumpToEndOfInterval(1)
            await rewardsContract.allocateRewards(1)

            await expireAllGroups()
            await rewardsContract.receiveReward(2)
            // each beneficiary receives 768000 / 64 = 12000 KEEP;
            // they should have 12250 + 12000 = 24500 KEEP now
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(24500))
        })
    })

    async function timeJumpToEndOfInterval(intervalNumber) {
        const endOf = await rewardsContract.endOf(intervalNumber)
        const now = await time.latest()

        if (now.lt(endOf)) {
            await time.increaseTo(endOf.addn(1))
        }
    }

    async function registerNewGroup() {
        const groupPublicKey = crypto.randomBytes(128)
        await operatorContract.registerNewGroup(groupPublicKey, operators)
    }

    async function expireAllGroups() {        
        const currentBlock = await time.latestBlock()

        const groupStalingTime = groupActiveTime + relayEntryTimeout + 1
        await time.advanceBlockTo(currentBlock.addn(groupStalingTime))
        await operatorContract.expireOldGroups()
    }

    async function assertKeepBalanceOfBeneficiaries(expectedBalance) {
        for (let i = 0; i < beneficiaries.length; i++) {
            const balance = await token.balanceOf(beneficiaries[i])
            const balanceInKeep = balance.div(tokenDecimalMultiplier)
            expect(balanceInKeep).to.eq.BN(expectedBalance)
        }
    }
})
