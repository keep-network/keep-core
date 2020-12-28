const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const {testValues} = require("./rewardsData.js")

const KeepToken = contract.fromArtifact("KeepToken")

const RewardsStub = contract.fromArtifact("RewardsStub")

const BN = web3.utils.BN
const chai = require("chai")
chai.use(require("bn-chai")(BN))
const expect = chai.expect
const assert = chai.assert

describe("Rewards", () => {
  const owner = accounts[0]
  const aliceBeneficiary = accounts[1]
  const funder = accounts[9]

  const termLength = 100

  let rewards
  let token

  async function createKeeps(timestamps) {
    rewards = await RewardsStub.new(
      token.address,
      testValues.minimumIntervalKeeps,
      testValues.initiationTime,
      testValues.intervalWeights,
      timestamps,
      termLength,
      {from: owner}
    )
    await fund(testValues.totalRewards)
    await rewards.markAsFunded({from: owner})
  }

  async function fund(amount) {
    await token.approveAndCall(rewards.address, amount, "0x0", {from: funder})
  }

  before(async () => {
    token = await KeepToken.new({from: funder})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("receiveApproval", async () => {
    it("funds the rewards correctly", async () => {
      await createKeeps([])
      const preRewards = await rewards.totalRewards()
      expect(preRewards.toNumber()).to.equal(testValues.totalRewards)

      await fund(testValues.totalRewards)
      const postRewards = await rewards.totalRewards()
      expect(postRewards.toNumber()).to.equal(testValues.totalRewards * 2)
    })

    it("collects tokens sent outside `approveAndCall`", async () => {
      await createKeeps([])
      await token.transfer(rewards.address, testValues.totalRewards, {
        from: funder,
      })

      const preRewards = await rewards.totalRewards()
      expect(preRewards.toNumber()).to.equal(testValues.totalRewards)

      await fund(0)
      const postRewards = await rewards.totalRewards()
      expect(postRewards.toNumber()).to.equal(testValues.totalRewards * 2)
    })
  })

  describe("markAsFunded", async () => {
    let newRewards

    beforeEach(async () => {
      newRewards = await RewardsStub.new(
        token.address,
        testValues.minimumIntervalKeeps,
        testValues.initiationTime,
        testValues.intervalWeights,
        [],
        termLength,
        {from: owner}
      )
    })

    it("can not be called by non-owner", async () => {
      await expectRevert(
        newRewards.markAsFunded({from: funder}),
        "Ownable: caller is not the owner"
      )
    })

    it("prevents from allocating rewards if not previously called", async () => {
      await expectRevert(
        newRewards.allocateRewards(0),
        "Contract has not been funded yet"
      )
    })
  })

  describe("eligibleForReward", async () => {
    it("returns true for happily closed keeps", async () => {
      await createKeeps([1000])
      await rewards.setCloseTime(1000)
      const eligible = await rewards.eligibleForReward(0)
      expect(eligible).to.equal(true)
    })

    it("returns false for terminated keeps", async () => {
      await createKeeps([1000])
      await rewards.terminate(0)
      const eligible = await rewards.eligibleForReward(0)
      expect(eligible).to.equal(false)
    })

    it("returns false for active keeps", async () => {
      await createKeeps([1000])
      await rewards.setCloseTime(999)
      const eligible = await rewards.eligibleForReward(0)
      expect(eligible).to.equal(false)
    })

    it("returns false for unrecognized keeps", async () => {
      await createKeeps([1000])
      const eligible = await rewards.eligibleForReward(1)
      expect(eligible).to.equal(false)
    })
  })

  describe("eligibleButTerminated", async () => {
    it("returns false for happily closed keeps", async () => {
      await createKeeps([1000])
      await rewards.setCloseTime(1000)
      const eligible = await rewards.eligibleButTerminatedWithUint(0)
      expect(eligible).to.equal(false)
    })

    it("returns true for terminated keeps", async () => {
      await createKeeps([1000])
      await rewards.terminate(0)
      const eligible = await rewards.eligibleButTerminatedWithUint(0)
      expect(eligible).to.equal(true)
    })

    it("returns false for active keeps", async () => {
      await createKeeps([1000])
      await rewards.setCloseTime(999)
      const eligible = await rewards.eligibleButTerminatedWithUint(0)
      expect(eligible).to.equal(false)
    })

    it("returns false for unrecognized keeps", async () => {
      await createKeeps([1000])
      const eligible = await rewards.eligibleButTerminatedWithUint(1)
      expect(eligible).to.equal(false)
    })
  })

  describe("intervalOf", async () => {
    it("returns the correct interval", async () => {
      await createKeeps([])
      const interval999 = await rewards.intervalOf(999)
      expect(interval999.toNumber()).to.equal(0)
      const interval1000 = await rewards.intervalOf(1000)
      expect(interval1000.toNumber()).to.equal(0)
      const interval1001 = await rewards.intervalOf(1001)
      expect(interval1001.toNumber()).to.equal(0)
      const interval1099 = await rewards.intervalOf(1099)
      expect(interval1099.toNumber()).to.equal(0)
      const interval1100 = await rewards.intervalOf(1100)
      expect(interval1100.toNumber()).to.equal(1)
      const interval1101 = await rewards.intervalOf(1101)
      expect(interval1101.toNumber()).to.equal(1)
      const interval1000000 = await rewards.intervalOf(1000000)
      expect(interval1000000.toNumber()).to.equal(9990)
    })
  })

  describe("startOf", async () => {
    it("returns the start of the interval", async () => {
      await createKeeps([])
      const start0 = await rewards.startOf(0)
      expect(start0.toNumber()).to.equal(1000)
      const start1 = await rewards.startOf(1)
      expect(start1.toNumber()).to.equal(1100)
      const start9990 = await rewards.startOf(9990)
      expect(start9990.toNumber()).to.equal(1000000)
    })
  })

  describe("endOf", async () => {
    it("returns the end of the interval", async () => {
      await createKeeps([])
      const end0 = await rewards.endOf(0)
      expect(end0.toNumber()).to.equal(1100)
      const end1 = await rewards.endOf(1)
      expect(end1.toNumber()).to.equal(1200)
      const end9990 = await rewards.endOf(9990)
      expect(end9990.toNumber()).to.equal(1000100)
    })
  })

  describe("findEndpoint", async () => {
    const increment = 1000

    it("returns 0 when no keeps have been created", async () => {
      await createKeeps([])
      const targetTimestamp = await time.latest()
      time.increase(increment)

      const index = await rewards.findEndpoint(targetTimestamp)
      expect(index.toNumber()).to.equal(0)
    })

    it("returns 0 when all current keeps were created after the interval", async () => {
      const timestamps = testValues.defaultTimestamps
      await createKeeps(timestamps)
      const targetTimestamp = 500
      const expectedIndex = 0
      const index = await rewards.findEndpoint(targetTimestamp)

      expect(index.toNumber()).to.equal(expectedIndex)
    })

    it("returns the first index outside the interval", async () => {
      const timestamps = testValues.defaultTimestamps
      await createKeeps(timestamps)
      for (let i = 0; i < timestamps.length; i++) {
        const expectedIndex = i
        const targetTimestamp = timestamps[i]
        const index = await rewards.findEndpoint(targetTimestamp)

        expect(index.toNumber()).to.equal(expectedIndex)
      }
    })

    it("returns the number of keeps when all current keeps were created in the interval", async () => {
      const timestamps = testValues.defaultTimestamps
      await createKeeps(timestamps)
      const targetTimestamp = 2000
      const expectedIndex = 16
      const index = await rewards.findEndpoint(targetTimestamp)

      expect(index.toNumber()).to.equal(expectedIndex)
    })

    it("returns the correct index when duplicates are present", async () => {
      const timestamps = [1001, 1001, 1002, 1002]
      await createKeeps(timestamps)
      const targetTimestamp = 1002
      const expectedIndex = 2
      const index = await rewards.findEndpoint(targetTimestamp)

      expect(index.toNumber()).to.equal(expectedIndex)
    })

    it("reverts if the endpoint is in the future", async () => {
      await createKeeps([])
      const recentTimestamp = await time.latest()
      const targetTimestamp = recentTimestamp + increment
      await expectRevert(
        rewards.findEndpoint(targetTimestamp),
        "interval hasn't ended yet"
      )
    })
  })

  describe("getEndpoint", async () => {
    it("returns the correct number of keeps for the interval", async () => {
      const timestamps = testValues.defaultTimestamps
      await createKeeps(timestamps)
      const keepCount = await rewards.getEndpoint.call(0)
      expect(keepCount.toNumber()).to.equal(timestamps.length)
    })

    it("returns 0 for intervals with no keeps", async () => {
      const timestamps = [1200, 1201]
      await createKeeps(timestamps)
      const keepCount = await rewards.getEndpoint.call(1)
      expect(keepCount.toNumber()).to.equal(0)
    })

    it("reverts if the interval hasn't ended", async () => {
      await createKeeps([])
      const recentTimestamp = await time.latest()
      const targetTimestamp = recentTimestamp + testValues.termLength
      const targetInterval = await rewards.intervalOf(targetTimestamp)
      await expectRevert(
        rewards.getEndpoint(targetInterval),
        "Interval hasn't ended yet"
      )
    })
  })

  describe("keepsInInterval", async () => {
    it("returns the correct number of keeps for the interval", async () => {
      const timestamps = testValues.rewardTimestamps
      const expectedCounts = testValues.keepsInRewardIntervals
      await createKeeps(timestamps)
      for (let i = 0; i < expectedCounts.length; i++) {
        const keepCount = await rewards.keepsInInterval.call(i)
        expect(keepCount.toNumber()).to.equal(expectedCounts[i])
      }
    })
  })

  describe("getIntervalWeight", async () => {
    it("returns the weight of a defined interval", async () => {
      await createKeeps([])
      const weight0 = await rewards.getIntervalWeight(0)
      expect(weight0.toNumber()).to.equal(20)
      const weight3 = await rewards.getIntervalWeight(3)
      expect(weight3.toNumber()).to.equal(50)
    })

    it("returns the last defined interval weight after the defined intervals", async () => {
      await createKeeps([])
      const weight4 = await rewards.getIntervalWeight(4)
      expect(weight4.toNumber()).to.equal(50)
    })
  })

  describe("getIntervalCount", async () => {
    it("returns the number of defined intervals", async () => {
      await createKeeps([])
      const intervalCount = await rewards.getIntervalCount()
      expect(intervalCount.toNumber()).to.equal(4)
    })
  })

  describe("baseAllocation", async () => {
    it("returns the maximum reward of a defined interval", async () => {
      await createKeeps([])
      const expectedAllocations = testValues.inVacuumBaseRewards
      for (let i = 0; i < expectedAllocations.length; i++) {
        const allocation = await rewards.baseAllocation(i)
        expect(allocation.toNumber()).to.equal(expectedAllocations[i])
      }
    })
  })

  describe("adjustedAllocation", async () => {
    it("returns the adjusted reward allocation of the interval", async () => {
      const timestamps = testValues.rewardTimestamps
      const expectedAllocations = testValues.inVacuumAdjustedRewards
      await createKeeps(timestamps)
      for (let i = 0; i < expectedAllocations.length; i++) {
        const allocation = await rewards.adjustedAllocation.call(i)
        expect(allocation.toNumber()).to.equal(expectedAllocations[i])
      }
    })
  })

  describe("rewardPerKeep", async () => {
    it("returns the per keep allocation of the interval", async () => {
      const timestamps = testValues.rewardTimestamps
      const expectedAllocations = testValues.inVacuumPerKeepRewards
      await createKeeps(timestamps)
      for (let i = 0; i < expectedAllocations.length; i++) {
        const allocation = await rewards.rewardPerKeep.call(i)
        expect(allocation.toNumber()).to.equal(expectedAllocations[i])
      }
    })
  })

  describe("allocateRewards", async () => {
    it("allocates the reward for each interval", async () => {
      const timestamps = testValues.rewardTimestamps
      const expectedAllocations = testValues.actualAllocations
      await createKeeps(timestamps)
      for (let i = 0; i < expectedAllocations.length; i++) {
        await rewards.allocateRewards(i)
        const allocation = await rewards.getAllocatedRewards(i)
        expect(allocation.toNumber()).to.equal(expectedAllocations[i])
      }
    })

    it("allocates the rewards recursively", async () => {
      const timestamps = testValues.rewardTimestamps
      const expectedAllocations = testValues.actualAllocations
      await createKeeps(timestamps)
      await rewards.allocateRewards(expectedAllocations.length - 1)
      for (let i = 0; i < expectedAllocations.length; i++) {
        const allocation = await rewards.getAllocatedRewards(i)
        expect(allocation.toNumber()).to.equal(expectedAllocations[i])
      }
    })
  })

  describe("isAllocated", async () => {
    it("returns false before allocation and true after allocation", async () => {
      const timestamps = testValues.rewardTimestamps
      const expectedAllocations = testValues.actualAllocations
      await createKeeps(timestamps)
      for (let i = 0; i < expectedAllocations.length; i++) {
        const preAllocated = await rewards.isAllocated(i)
        expect(preAllocated).to.equal(false)
        await rewards.allocateRewards(i)
        const postAllocated = await rewards.isAllocated(i)
        expect(postAllocated).to.equal(true)
      }
    })
  })

  describe("receiveReward", async () => {
    it("lets closed keeps claim the reward correctly", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)
      await rewards.setCloseTime(timestamps[0])
      await rewards.receiveReward(0, {from: aliceBeneficiary})
      const aliceBalance = await token.balanceOf(aliceBeneficiary)
      expect(aliceBalance.toNumber()).to.equal(66666)
    })

    it("doesn't let keeps claim rewards again", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)
      await rewards.setCloseTime(timestamps[0])
      await rewards.receiveReward(0, {from: aliceBeneficiary})
      await expectRevert(
        rewards.receiveReward(0, {from: aliceBeneficiary}),
        "Rewards already claimed"
      )
    })

    it("doesn't let active keeps claim the reward", async () => {
      await createKeeps(testValues.rewardTimestamps)
      await expectRevert(
        rewards.receiveReward(0, {from: aliceBeneficiary}),
        "Keep is not closed"
      )
    })

    it("doesn't let terminated keeps claim the reward", async () => {
      await createKeeps(testValues.rewardTimestamps)
      await rewards.terminate(0)
      await expectRevert(
        rewards.receiveReward(0, {from: aliceBeneficiary}),
        "Keep is not closed"
      )
    })

    it("doesn't let unrecognized keeps claim the reward", async () => {
      await createKeeps(testValues.rewardTimestamps)
      await expectRevert(
        rewards.receiveReward(testValues.rewardTimestamps.length),
        "Keep not recognized by factory"
      )
    })

    it("requires that the interval is over", async () => {
      const recentTimestamp = await time.latest()
      const targetTimestamp = recentTimestamp + 1000
      await createKeeps([targetTimestamp])
      await rewards.setCloseTime(targetTimestamp)
      await expectRevert(rewards.receiveReward(0), "Interval hasn't ended yet")
    })

    it("emits an event", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)
      await rewards.setCloseTime(timestamps[0])
      await rewards.receiveReward(0, {from: aliceBeneficiary})
      assert.equal(
        (await rewards.getPastEvents())[0].event,
        "RewardReceived",
        "Should emit event"
      )
    })
  })

  describe("receiveRewards", async () => {
    it("lets closed keeps claim the rewards from multiple keeps correctly", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)
      await rewards.setCloseTime(timestamps[2])
      const rewardsReceivingKeeps = [0, 1]

      await rewards.methods["receiveRewards(uint256[])"](
        rewardsReceivingKeeps,
        {from: aliceBeneficiary}
      )
      const aliceBalance = await token.balanceOf(aliceBeneficiary)
      // Beneficiary will receive 200,000 / 3 = 66,666 per keep
      // 66,666 * 2 = 133,332 KEEP rewards total for being in 2 closed keeps
      expect(aliceBalance.toNumber()).to.equal(133332)
    })
  })

  describe("reportTermination", async () => {
    it("unallocates rewards allocated to terminated keep", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)

      await rewards.setCloseTime(testValues.rewardTimestamps[0])
      await rewards.receiveReward(0, {from: aliceBeneficiary}) // allocate rewards

      await rewards.terminate(1)
      const preUnallocated = await rewards.unallocatedRewards()
      await rewards.reportTermination(1)
      const postUnallocated = await rewards.unallocatedRewards()
      expect(postUnallocated.toNumber()).to.equal(
        preUnallocated.toNumber() + 66666
      )
    })

    it("doesn't unallocate rewards twice for the same keep", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)
      await rewards.terminate(0)
      await rewards.reportTermination(0)
      await expectRevert(
        rewards.reportTermination(0),
        "Rewards already claimed"
      )
    })

    it("doesn't unallocate active keeps' rewards", async () => {
      await createKeeps(testValues.rewardTimestamps)
      await expectRevert(rewards.reportTermination(0), "Keep is not terminated")
    })

    it("doesn't unallocate closed keeps' rewards", async () => {
      await createKeeps(testValues.rewardTimestamps)
      await rewards.setCloseTime(testValues.rewardTimestamps[0])
      await expectRevert(rewards.reportTermination(0), "Keep is not terminated")
    })

    it("doesn't unallocate unrecognized keeps' rewards", async () => {
      await createKeeps(testValues.rewardTimestamps)
      await expectRevert(
        rewards.reportTermination(testValues.rewardTimestamps.length),
        "Keep not recognized by factory"
      )
    })

    it("requires that the interval is over", async () => {
      const recentTimestamp = await time.latest()
      const targetTimestamp = recentTimestamp + 1000
      await createKeeps([targetTimestamp])
      await rewards.terminate(0)
      await expectRevert(
        rewards.reportTermination(0),
        "Interval hasn't ended yet"
      )
    })
  })

  describe("reportTerminations", async () => {
    it("unallocates rewards allocated to terminated keeps in batch", async () => {
      const timestamps = testValues.rewardTimestamps
      await createKeeps(timestamps)

      await rewards.setCloseTime(testValues.rewardTimestamps[0])
      await rewards.allocateRewards(0)

      await rewards.terminate(1)
      await rewards.terminate(2)

      const preUnallocated = await rewards.unallocatedRewards()

      const terminatedIdentifiers = [1, 2]
      await rewards.methods["reportTerminations(uint256[])"](
        terminatedIdentifiers
      )
      const actual = await rewards.unallocatedRewards()
      // 200,000 KEEP were allocated for the first interval
      // 800,000 KEEP remaining in unallocated pool
      // 2 out of 3 keeps were terminated
      // 200,000 / 3 = 66,666 rewards per keep
      // 66,666 * 2 = 133,332 returned back to unallocated pool
      expect(actual.toNumber()).to.equal(preUnallocated.toNumber() + 133332)
    })
  })
})
