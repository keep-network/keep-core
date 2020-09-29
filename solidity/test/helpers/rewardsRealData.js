const { time } = require("@openzeppelin/test-helpers")

async function testData() {
    const termLength = time.duration.days(30) //  30 days in sec
    const daysToFirstInterval = time.duration.days(30) // days left to first interval
    
    // Deployment today
    const now = await time.latest() // in sec.
    // Today + 30 days = first interval start
    const firstIntervalStart = now.add(daysToFirstInterval)
    
    // const firstIntervalStart = web3.utils.toBN(1000) // in sec.
    // const termLength = time.duration.seconds(100) //  30 days in sec

    const keepsInRewardIntervals = [
        3, 4, 3, 2, 3, 2,
        2, 3, 3, 3, 4, 2,
        3, 5, 6, 4, 5, 5,
        5, 6, 6, 6, 4, 3
    ]

    async function createRewardTimestampsInInterval() {
        const rewardTimestampsForIntervals = []
        const numberOfIntervals = keepsInRewardIntervals.length // should be 24

        for (let i = 0; i < numberOfIntervals; i++) {
            const rewardTimestampsForInterval = []
            // iterate over keeps in a given (i) interval
            for (let j = 0; j < keepsInRewardIntervals[i]; j++) {
                // rewardTimestamp = firstIntervalStart + (termLength * intervalNum) + numOfSec
                const rewardTimestamp = firstIntervalStart.add((termLength.muln(i))).addn(j)
                rewardTimestampsForInterval.push(rewardTimestamp)
            }

            rewardTimestampsForIntervals.push(rewardTimestampsForInterval)
        }

        return rewardTimestampsForIntervals.flat()
    }

    return {
        firstIntervalStart: firstIntervalStart,
        termLength: termLength,
        totalRewards: 20000000, // 2% KEEP of beacon subsidy pool
        minimumIntervalKeeps: 2, // TODO: define min keeps(groups) per interval

        rewardTimestamps: await createRewardTimestampsInInterval(),
        keepsInRewardIntervals: keepsInRewardIntervals,

        // percentage of unallocated rewards in 24 intervals (30days each)
        intervalWeights: [
            4, 8, 10, 12, 15, 15,
            15, 15, 15, 15, 15, 15,
            15, 15, 15, 15, 15, 15,
            15, 15, 15, 15, 15, 15
        ],
        expectedAllocations: [
            800000,
            1536000,
            1766400,
            1907712,
            2098483.20,
            1783710.72,
            1516154.11,
            1288731,
            1095421.35,
            931108.14,
            791441.92,
            672725.63,
            571816.79,
            486044.27,
            413137.63,
            351166.99,
            298491.94,
            253718.15,
            215660.42,
            183311.36,
            155814.66,
            132442.46,
            112576.09,
            95689.68,
        ],
        expectedMemberRewardsPerInterval: [
            12500,
            36500,
            64100,
            93908,
            126696.8,
            154567.28,
            178257.188,
            198393.6098,
            215509.5683,
            230058.1331,
            242424.4131,
            252935.7512,
            261870.3885,
            269464.8302,
            275920.1057,
            281407.0898,
            286071.0264,
            290035.3724,
            293405.0665,
            296269.3066,
            298703.9106,
            300773.324,
            302532.3254,
            304027.4766
        ],
        expectedAllocatedRewards: 19457758.5,
    }
};


module.exports.testData = testData
