const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { expectRevert, expectEvent } = require("@openzeppelin/test-helpers")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")

const { toBN } = require("web3-utils")

const TestToken = contract.fromArtifact("TestToken")
const TokenGeyser = contract.fromArtifact("TokenGeyser")

const BN = web3.utils.BN
const chai = require("chai")
chai.use(require("bn-chai")(BN))
const expect = chai.expect

// We test only parts that we modified in the original `ampleforth/token-geyser`
// implementation.
describe("TokenGeyser", async () => {
  const contractOwner = accounts[1]
  const rewardDistribution = accounts[2]
  const thirdParty = accounts[6]

  const maxUnlockSchedules = toBN(12)
  const startBonus = toBN(100)
  const bonusPeriodSec = toBN(1)
  const initialSharesPerToken = toBN(1)
  const rewardsAmount = toBN(5000)

  let stakeToken
  let distributionToken
  let tokenGeyser

  before(async () => {
    stakeToken = await TestToken.new({ from: contractOwner })
    distributionToken = await TestToken.new({ from: contractOwner })

    tokenGeyser = await TokenGeyser.new(
      stakeToken.address,
      distributionToken.address,
      maxUnlockSchedules,
      startBonus,
      bonusPeriodSec,
      initialSharesPerToken,
      {
        from: contractOwner,
      }
    )

    await tokenGeyser.setRewardDistribution(rewardDistribution, {
      from: contractOwner,
    })

    // Fund accounts with tokens.
    await distributionToken.mint(rewardDistribution, rewardsAmount)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("constructor", async () => {
    it("sets reward distribution to owner by default", async () => {
      const tokenGeyser = await TokenGeyser.new(
        stakeToken.address,
        distributionToken.address,
        maxUnlockSchedules,
        startBonus,
        bonusPeriodSec,
        initialSharesPerToken,
        {
          from: contractOwner,
        }
      )

      expect(await tokenGeyser.rewardDistribution.call()).to.equal(
        contractOwner
      )
    })
  })

  describe("setRewardDistribution", async () => {
    const newRewardDistribution = thirdParty

    it("updates rewardDistribution", async () => {
      await tokenGeyser.setRewardDistribution(newRewardDistribution, {
        from: contractOwner,
      })

      expect(await tokenGeyser.rewardDistribution.call()).to.eq.BN(
        newRewardDistribution
      )
    })

    it("reverts when called by non-owner", async () => {
      await expectRevert(
        tokenGeyser.setRewardDistribution(newRewardDistribution),
        "Ownable: caller is not the owner"
      )
    })

    it("emits event", async () => {
      const receipt = await tokenGeyser.setRewardDistribution(
        newRewardDistribution,
        {
          from: contractOwner,
        }
      )

      expectEvent(receipt, "RewardDistributionRoleTransferred", {
        oldRewardDistribution: rewardDistribution,
        newRewardDistribution: newRewardDistribution,
      })
    })
  })

  describe("lockTokens", async () => {
    const durationSec = 1234567890

    it("succeeds when called by reward distribution", async () => {
      await distributionToken.approve(tokenGeyser.address, rewardsAmount, {
        from: rewardDistribution,
      })

      await tokenGeyser.lockTokens(rewardsAmount, durationSec, {
        from: rewardDistribution,
      })
    })

    it("reverts when called by non-reward distribution", async () => {
      await expectRevert(
        tokenGeyser.lockTokens(rewardsAmount, durationSec),
        "Caller is not the reward distribution"
      )
    })
  })
})
