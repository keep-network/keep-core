const { time } = require("@openzeppelin/test-helpers")
const { web3 } = require("@openzeppelin/test-environment")

function testData() {
    const initiationTime = web3.utils.toBN(1000) // in sec.
    const termLength = time.duration.seconds(100) //  30 days in sec

    const keepsInRewardInterval = [
        3, 4, 3, 2, 3, 2,
        2, 3, 3, 3, 4, 2,
        3, 5, 6, 4, 5, 5,
        5, 6, 6, 6, 4, 3
    ]

    function createRewardTimestampsInInterval() {
        const rewardTimestampsForIntervals = []
        const numberOfIntervals = keepsInRewardInterval.length // should be 24

        for (let i = 0; i < numberOfIntervals; i++) {
            const rewardTimestampsForInterval = []
            // iterate over keeps in a given (i) interval
            for (let j = 0; j < keepsInRewardInterval[i]; j++) {
                // rewardTimestamp = initiationTime + (termLength * intervalNum) + numOfSec
                const rewardTimestamp = initiationTime.add((termLength.muln(i))).addn(j)
                rewardTimestampsForInterval.push(rewardTimestamp)
            }

            rewardTimestampsForIntervals.push(rewardTimestampsForInterval)
        }

        return rewardTimestampsForIntervals.flat()
    }

    return {
        initiationTime: initiationTime,
        termLength: termLength,
        totalRewards: 20000000, // 2% KEEP of beacon subsidy pool
        minimumIntervalKeeps: 2, // TODO: define min keeps(groups) per interval

        rewardTimestamps: createRewardTimestampsInInterval(),
        keepsInRewardIntervals: keepsInRewardInterval,

        // percentage of unallocated rewards in 24 intervals (30days each)
        intervalWeights: [
            4, 8, 10, 12, 15, 15,
            15, 15, 15, 15, 15, 15,
            15, 15, 15, 15, 15, 15,
            15, 15, 15, 15, 15, 15
        ],
        actualAllocations: [
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
        ]
    }
};


module.exports.testData = testData
