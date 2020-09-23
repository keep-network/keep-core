const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")
const testRealData = require("./helpers/rewardsRealData.js")

const KeepToken = contract.fromArtifact('KeepToken')
const RewardsStub = contract.fromArtifact('RewardsStub');

describe.only('Beacon 2% Rewards', () => {

    let token, testValues, rewards
    const funder = accounts[5]

        
    async function createKeeps(timestamps) {
        rewards = await RewardsStub.new(
            token.address,
            testValues.minimumIntervalKeeps,
            testValues.initiationTime,
            testValues.intervalWeights,
            timestamps
        )
        await fund(testValues.totalRewards)
    }

    async function fund(amount) {
        await token.approveAndCall(
            rewards.address,
            amount,
            "0x0",
            { from: funder }
        )
    }
    
    before(async () => {
        testValues = await testRealData.testData()
        token = await KeepToken.new({ from: funder })
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })


    describe("allocate rewards for intervals", async () => {
        it("allocates the reward for each interval", async () => {
            let timestamps = testValues.rewardTimestamps
            let expectedAllocations = testValues.actualAllocations

            await createKeeps(timestamps)

            let sumOfActualAllocation = 0;
            let sumOfExpectedAllocation = 0;
            console.log("")
            console.log("Interval: Actual: Expected: ")
            for (let i = 0; i < expectedAllocations.length; i++) {
                await rewards.allocateRewards(i)
                let allocation = await rewards.getAllocatedRewards(i)
                sumOfActualAllocation += allocation.toNumber()
                sumOfExpectedAllocation += expectedAllocations[i]

                console.log(`${i} : ${allocation.toNumber()} : ${expectedAllocations[i]}`)
            }

            console.log("")
            console.log("Total actual allocations: ", sumOfActualAllocation)
            console.log("Total expected allocations: ", sumOfExpectedAllocation)
            console.log("Diff: ", sumOfExpectedAllocation - sumOfActualAllocation)


        })
    })
})
