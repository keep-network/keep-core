const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { expectRevert, time } = require("@openzeppelin/test-helpers")

const KeepToken = contract.fromArtifact('KeepToken')
const RewardsStub = contract.fromArtifact("RewardsStub");
const NewRewardsStub = contract.fromArtifact("NewRewardsStub");

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe("Rewards/Upgrades", () => {
    const owner = accounts[0]
    const thirdParty = accounts[1]
    const beneficiary = accounts[3]

    const intervalWeights = [5, 10, 15, 20]
    const totalRewards = 1000000
    const minimumKeepsPerInterval = 1
    const termLength = 1000
    let timestamps
    let firstIntervalStart

    let token, rewards, newRewards

    before(async () => {
        token = await KeepToken.new({ from: owner })
        newRewards = await NewRewardsStub.new()

        firstIntervalStart = await time.latest()
        timestamps = [
            101, 150,    // 2 keep in interval 0
            1100,        // 1 keep in interval 1
            2200, 2201   // 2 keeps in interval 2
        ].map((t) => firstIntervalStart.addn(t).toNumber())

        rewards = await RewardsStub.new(
            token.address,
            minimumKeepsPerInterval,
            firstIntervalStart,
            intervalWeights,
            timestamps,
            termLength,
            {from: owner}
         )
         await token.approveAndCall(
            rewards.address,
            totalRewards,
            "0x0",
            {from: owner}
        )
        await rewards.markAsFunded({from: owner})
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("upgrades", async () => {
        it("can be initiated only by contract owner", async () => {            
            await expectRevert(
                rewards.initiateRewardsUpgrade(
                    newRewards.address,
                    {from: thirdParty}
                ),
                "Ownable: caller is not the owner"
            )
        })

        it("can be finalized only by contract owner", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )
            await expectRevert(
                rewards.finalizeRewardsUpgrade({from: thirdParty}),
                "Ownable: caller is not the owner."
            )            
        })

        it("cannot be finalized without initiating first", async () => {
            await expectRevert(
                rewards.finalizeRewardsUpgrade({from: owner}),
                "Upgrade not initiated"
            ) 
        })

        it("cannot be finalized before the current interval ends", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )
            await expectRevert(
                rewards.finalizeRewardsUpgrade({from: owner}),
                "Interval hasn't ended yet"
            )
        })

        it("cannot be finalized another time without initiating again", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )

            await time.increase(termLength + 1) // interval 0 ends
            await rewards.finalizeRewardsUpgrade({from: owner})
            await expectRevert(
                rewards.finalizeRewardsUpgrade({from: owner}),
                "Upgrade not initiated"
            ) 
        })

        it("should not change the current interval allocation", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )

            await time.increase(termLength + 1) // interval 0 ends
            await rewards.setCloseTime(timestamps[1])

            await rewards.finalizeRewardsUpgrade({from: owner})

            const allocation = await rewards.getAllocatedRewards(0)
            expect(allocation).to.eq.BN(50000) // 5% of 1 000 000
        })

        it("allocates all possible intervals", async () => {
            await time.increase(termLength + 1) // interval 0 ends

            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )

            await time.increase(termLength + 1) // interval 1 ends
            await rewards.setCloseTime(timestamps[2])

            await rewards.finalizeRewardsUpgrade({from: owner})

            const allocation0 = await rewards.getAllocatedRewards(0)
            const allocation1 = await rewards.getAllocatedRewards(1)
            
            expect(allocation0).to.eq.BN(50000) // 5% of 1000000
            expect(allocation1).to.eq.BN(95000) // 10% of (1000000 - 50000)
            await expectRevert(
                rewards.getAllocatedRewards(2),
                "Interval not allocated yet"
            )
        })

        it("moves all unallocated rewards to new contract", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )

            await time.increase(2 * termLength + 1)  
            
            await rewards.setCloseTime(timestamps[2])
            await rewards.finalizeRewardsUpgrade({from: owner})

            const newContractBalance = await token.balanceOf(newRewards.address)
            // interval 0 allocates 50000
            // interval 1 allocates 95000
            // 1000000 - 50000 - 95000 = 855000 should be transferred to the
            // new contract
            expect(newContractBalance).to.eq.BN(855000)
        })

        it("correctly updates reward balances", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )

            await time.increase(2 * termLength + 1)  
            
            await rewards.setCloseTime(timestamps[2])
            await rewards.finalizeRewardsUpgrade({from: owner})

            const totalRewards = await rewards.totalRewards()
            const unallocatedRewards = await rewards.unallocatedRewards()
            const paidOutRewards = await rewards.paidOutRewards()

            // interval 0 allocates 50000
            // interval 1 allocates 95000
            // 50000 + 95000 = 145000
            expect(totalRewards).to.eq.BN(145000)
            expect(unallocatedRewards).to.eq.BN(0)
            // nothing yet withdrawn but 855000 transferred to the new contract
            expect(paidOutRewards).to.eq.BN(855000)
        })

        it("let to withdraw outstanding rewards after migration", async () => {
            await rewards.initiateRewardsUpgrade(
                newRewards.address,
                {from: owner}
            )

            await time.increase(2 * termLength + 1)  
            
            await rewards.setCloseTime(timestamps[2])
            await rewards.finalizeRewardsUpgrade({from: owner})

            await rewards.receiveReward(0, { from: beneficiary })
            await rewards.receiveReward(1, { from: beneficiary })
            await rewards.receiveReward(2, { from: beneficiary })
            const beneficiaryBalance = await token.balanceOf(beneficiary)
            // interval 0 allocates 50000
            // interval 1 allocates 95000
            // 50000 + 95000 = 145000
            expect(beneficiaryBalance).to.eq.BN(145000)
        })
    })
})