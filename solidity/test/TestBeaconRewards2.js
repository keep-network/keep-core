const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")
const testRealData = require("./helpers/rewardsRealData.js")
const { time } = require("@openzeppelin/test-helpers")
const { initContracts } = require('./helpers/initContracts')
const stakeDelegate = require('./helpers/stakeDelegate')
const crypto = require("crypto")

const BeaconRewardsStub = contract.fromArtifact('BeaconRewardsStub');

const BN = web3.utils.BN

describe('Beacon Only 2% Rewards', () => {
    const owner = accounts[0]
    const groupSize = 64 // total number of operators available for test

    let tokenContract, testValues, rewards, minimumStake, allOperators

    before(async () => {
        let contracts = await initContracts(
            contract.fromArtifact('TokenStaking'),
            contract.fromArtifact('KeepRandomBeaconService'),
            contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
            contract.fromArtifact('KeepRandomBeaconOperatorBeaconRewardsStub')
        )

        tokenContract = contracts.token
        stakingContract = contracts.stakingContract
        operatorContract = contracts.operatorContract
        
        testValues = await testRealData.testData()
        rewards = await BeaconRewardsStub.new(
            tokenContract.address,
            testValues.firstIntervalStart,
            testValues.intervalWeights,
            operatorContract.address,
            stakingContract.address,
        )

        console.log("'deployment' time: ", (await time.latest()).toString())
    
        await tokenContract.approveAndCall(
            rewards.address,
            testValues.totalRewards,
            "0x0",
            { from: owner }
        )

        minimumStake = await stakingContract.minimumStake()
        allOperators = []
        for (let i = 1; i <= groupSize; i++) {
            const operator = accounts[i]
            const beneficiary = operator
            const authorizer = operator
            await stakeDelegate(stakingContract, tokenContract, owner, operator, beneficiary, authorizer, minimumStake.muln(20))
            allOperators.push(operator)
        }
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe.only("receive rewards", async () => {
        it("allocates the reward for each interval", async () => {
            // pass "0" termLenght, so the inverval starts with 1
            await time.increaseTo(testValues.firstIntervalStart.add(time.duration.minutes(1)))

            console.log("firstIntervalStart: ", (await rewards.firstIntervalStart()).toString())
            console.log("termLength: ", (await rewards.termLength()).toString())

            // iterate over 24 months
            for (let i = 0; i < testValues.intervalWeights.length; i++) {
                // iterate over keeps in the given interval and register a group
                for (let keepNumber = 0; keepNumber < testValues.keepsInRewardIntervals[i]; keepNumber++) {
                    let group = crypto.randomBytes(128)
                    // let operators = []
                    // for (let op = 0; op < groupSize; op++) { // 64 operators in a given group
                    //     // let random = getRandomInt(groupSize) // total of 100 operators accounts available (some will be )
                    //     let operator = allOperators[random]
                    //     operators.push(operator)
                    // }
                    await operatorContract.registerNewGroup(group, allOperators)
                }
                console.log(`interval ${i}; groups creation time:  ${(await time.latest()).toString()}`)
                
                // go to the next interval
                await time.increase(testValues.termLength.add(time.duration.minutes(1)))
                const latestBlock = await time.latestBlock()
                await time.advanceBlockTo(latestBlock.addn(20))
            }
            await operatorContract.expireOldGroups()

            let totalKeepCount = await rewards.getKeepCount()
            console.log("Total keeps created: ", totalKeepCount.toString())

            let keepArrIndex = 0
            let accumulatedKeeps = testValues.keepsInRewardIntervals[keepArrIndex] - 1
            
            let totalRewardsBalance = new BN(0)
            for (let i = 0; i < totalKeepCount; i++) {
                await rewards.receiveReward(i)
                
                if (accumulatedKeeps == i) {
                    console.log("------------------")
                    for (let j = 0; j < allOperators.length; j++) {
                        const operatorBalance = await tokenContract.balanceOf(allOperators[j])
                        const diff = operatorBalance - testValues.expectedMemberRewardsPerInterval[keepArrIndex]
                        console.log(`interval ${i}: rewards accumulated for ${allOperators[j]} actual: ${operatorBalance}; expected: ${testValues.expectedMemberRewardsPerInterval[keepArrIndex]}; diff: ${diff}`)
                        if (i == totalKeepCount - 1) {
                            totalRewardsBalance = operatorBalance.add(totalRewardsBalance)
                        }
                    }
                    keepArrIndex++
                    accumulatedKeeps += testValues.keepsInRewardIntervals[keepArrIndex]
                }
            }

            const expectedAllocatedRewards = new BN(testValues.expectedAllocatedRewards) 
            const diffRewards = expectedAllocatedRewards.sub(totalRewardsBalance)
            console.log(`\nActual allocated rewards among members: ${totalRewardsBalance.toString()}; expected: ${expectedAllocatedRewards}; diff ${diffRewards.toString()}`)
        })
    })
})


function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
}
