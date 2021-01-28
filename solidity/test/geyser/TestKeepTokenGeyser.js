const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const {
  expectRevert,
  expectEvent,
  time,
} = require("@openzeppelin/test-helpers")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")

const { toBN } = require("web3-utils")

const KeepToken = contract.fromArtifact("KeepToken")
const TestToken = contract.fromArtifact("TestToken")
const KeepTokenGeyser = contract.fromArtifact("KeepTokenGeyser")
const BatchedPhasedEscrow = contract.fromArtifact("BatchedPhasedEscrow")
const KeepTokenGeyserRewardsEscrowBeneficiary = contract.fromArtifact(
  "KeepTokenGeyserRewardsEscrowBeneficiary"
)

const BN = web3.utils.BN
const chai = require("chai")
chai.use(require("bn-chai")(BN))
const expect = chai.expect

describe("KeepTokenGeyser", async () => {
  const contractOwner = accounts[1]
  const rewardDistribution = accounts[2]
  const staker1 = accounts[3]
  const staker2 = accounts[4]
  const beneficiary = accounts[5]
  const thirdParty = accounts[6]

  const maxUnlockSchedules = toBN(12) // ????
  const startBonus = toBN(30) //toBN(20) // 30% // ????
  const bonusPeriodSec = time.duration.weeks(4)
  const initialSharesPerToken = toBN(1) // ????
  const durationSec = time.duration.weeks(4)

  const tokenDecimalMultiplier = toBN(10e18) // 18-decimal precision

  const stakerInitialBalance1 = toBN(40e3).mul(tokenDecimalMultiplier) // 40k KEEP
  const stakeAmount1 = toBN(16e3).mul(tokenDecimalMultiplier) // 16k KEEP

  const stakerInitialBalance2 = toBN(120e3).mul(tokenDecimalMultiplier) // 120k KEEP
  const stakeAmount2 = toBN(112e3).mul(tokenDecimalMultiplier) // 112k KEEP

  const rewardsAmount = toBN(100e3).mul(tokenDecimalMultiplier) // 100k KEEP

  let stakeToken
  let keepToken
  let tokenGeyser

  before(async () => {
    stakeToken = await TestToken.new({ from: contractOwner })
    keepToken = await KeepToken.new({ from: contractOwner })

    tokenGeyser = await KeepTokenGeyser.new(
      stakeToken.address,
      keepToken.address,
      maxUnlockSchedules,
      startBonus,
      bonusPeriodSec,
      initialSharesPerToken,
      durationSec,
      {
        from: contractOwner,
      }
    )

    await tokenGeyser.setRewardDistribution(rewardDistribution, {
      from: contractOwner,
    })

    // Fund accounts with tokens.
    await stakeToken.mint(staker1, stakerInitialBalance1)
    await stakeToken.mint(staker2, stakerInitialBalance2)
    await keepToken.transfer(rewardDistribution, rewardsAmount, {
      from: contractOwner,
    })
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("lockTokens", async () => {
    describe("called via receiveApproval", async () => {
      it("should update balances", async () => {
        // Deploy Escrow.
        const escrow = await BatchedPhasedEscrow.new(keepToken.address, {
          from: contractOwner,
        })

        // Configure escrow beneficiary.
        const escrowBeneficiary = await KeepTokenGeyserRewardsEscrowBeneficiary.new(
          keepToken.address,
          tokenGeyser.address,
          {
            from: contractOwner,
          }
        )

        await escrowBeneficiary.transferOwnership(escrow.address, {
          from: contractOwner,
        })

        await escrow.approveBeneficiary(escrowBeneficiary.address, {
          from: contractOwner,
        })

        tokenGeyser.setRewardDistribution(escrowBeneficiary.address, {
          from: contractOwner,
        })

        // Transfer tokens to Escrow.
        await keepToken.approveAndCall(escrow.address, rewardsAmount, [], {
          from: contractOwner,
        })

        // Initiate withdraw.
        const initialEscrowBalance = await keepToken.balanceOf.call(
          escrow.address
        )

        await escrow.batchedWithdraw(
          [escrowBeneficiary.address],
          [rewardsAmount],
          { from: contractOwner }
        )

        expect(
          await keepToken.balanceOf.call(escrow.address),
          "invalid escrow's token balance"
        ).to.eq.BN(initialEscrowBalance.sub(rewardsAmount))

        expect(
          await tokenGeyser.totalStakedFor(escrowBeneficiary.address),
          "invalid reward distribution's staked balance"
        ).to.eq.BN(0)

        expect(
          await tokenGeyser.totalLocked(),
          "invalid total locked balance"
        ).to.eq.BN(rewardsAmount)
      })

      it("reverts when called by non-reward distribution", async () => {
        await expectRevert(
          tokenGeyser.receiveApproval(
            rewardDistribution,
            rewardsAmount,
            keepToken.address,
            []
          ),
          "Caller is not the reward distribution"
        )
      })

      it("reverts for not supported distribution token", async () => {
        await expectRevert(
          tokenGeyser.receiveApproval(
            rewardDistribution,
            rewardsAmount,
            thirdParty,
            []
          ),
          "Token is not supported distribution token"
        )
      })
    })
  })

  describe("setDurationSec", async () => {
    const newDurationSec = toBN(1230987)

    it("updates durationSec", async () => {
      await tokenGeyser.setDurationSec(newDurationSec, { from: contractOwner })

      expect(await tokenGeyser.durationSec.call()).to.eq.BN(newDurationSec)
    })

    it("reverts when called by non-owner", async () => {
      await expectRevert(
        tokenGeyser.setDurationSec(newDurationSec),
        "Ownable: caller is not the owner"
      )
    })

    it("emits event", async () => {
      const oldDurationSec = await tokenGeyser.durationSec.call()

      const receipt = await tokenGeyser.setDurationSec(newDurationSec, {
        from: contractOwner,
      })

      expectEvent(receipt, "DurationSecUpdated", {
        oldDurationSec: oldDurationSec,
        newDurationSec: newDurationSec,
      })
    })
  })

  describe("stake", async () => {
    it("should update balances", async () => {
      await stake(staker1, stakeAmount1)

      expect(
        await stakeToken.balanceOf.call(staker1),
        "invalid staker's token balance"
      ).to.eq.BN(stakerInitialBalance1.sub(stakeAmount1))

      expect(
        await tokenGeyser.totalStakedFor(staker1),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount1)
    })

    it("should emit event", async () => {
      const receipt = await stake(staker1, stakeAmount1)

      expectEvent(receipt, "Staked", {
        user: staker1,
        amount: stakeAmount1,
        total: stakeAmount1,
      })
    })

    it("allows stake top-ups", async () => {
      const topUpAmount = toBN(8e3).mul(tokenDecimalMultiplier) // 8k KEEP

      await stake(staker1, stakeAmount1)

      expect(
        await tokenGeyser.totalStakedFor(staker1),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount1)

      const receipt = await stake(staker1, topUpAmount)

      expect(
        await tokenGeyser.totalStakedFor(staker1),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount1.add(topUpAmount))

      expectEvent(receipt, "Staked", {
        user: staker1,
        amount: topUpAmount,
        total: stakeAmount1.add(topUpAmount),
      })
    })
  })

  describe("stakeFor", async () => {
    it("should update balances", async () => {
      await stakeToken.mint(contractOwner, stakeAmount1)

      await stakeToken.approve(tokenGeyser.address, stakeAmount1, {
        from: contractOwner,
      })
      await tokenGeyser.stakeFor(beneficiary, stakeAmount1, [], {
        from: contractOwner,
      })

      expect(
        await stakeToken.balanceOf.call(staker1),
        "invalid staker's token balance"
      ).to.eq.BN(stakerInitialBalance1)

      expect(
        await stakeToken.balanceOf.call(contractOwner),
        "invalid contract owner's token balance"
      ).to.eq.BN(0)

      expect(
        await tokenGeyser.totalStakedFor(contractOwner),
        "invalid staker's staked balance"
      ).to.eq.BN(0)

      expect(
        await tokenGeyser.totalStakedFor(beneficiary),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount1)
    })
  })

  describe("unstake", async () => {
    it("should calculate rewards for two stakers", async () => {
      const expectedRewards1 = toBN(125e2).mul(tokenDecimalMultiplier) // (16k / (16k + 112k)) * 100k = 12.5k KEEP
      const expectedRewards2 = toBN(875e2).mul(tokenDecimalMultiplier) // (112k / (16k + 112k)) * 100k = 87.5k KEEP

      await stake(staker1, stakeAmount1)
      await stake(staker2, stakeAmount2)

      const lockTimestamp1 = await lockTokens(rewardsAmount, durationSec)

      // End first interval.
      await time.increaseTo(lockTimestamp1.add(durationSec))

      await checkRewards(
        staker1,
        stakeAmount1,
        expectedRewards1,
        "invalid calculated staker's 1 rewards"
      )
      await checkRewards(
        staker2,
        stakeAmount2,
        expectedRewards2,
        "invalid calculated staker's 2 rewards"
      )
    })

    it("should calculate rewards in bonus period", async () => {
      // Here we estimate rewards taking into account bonus period.
      // With an assumption that bonus starts at 30% and goes to 100% over a bonus
      // period, we will check rewards in the middle of the bonus period.
      // In the middle of the bonus period rewards factor would be at 65%, hence
      // this is calculation of the expected rewards:
      //  staker 1: [(16k / (16k + 112k)) * 100k] * 50% * 65% = 4062.5k KEEP
      //  staker 2: [(112k / (16k + 112k)) * 100k] * 50% * 65% = 28437.5k KEEP
      const expectedRewards1 = toBN(40625).mul(tokenDecimalMultiplier).divn(10)
      const expectedRewards2 = toBN(284375).mul(tokenDecimalMultiplier).divn(10)

      await stake(staker1, stakeAmount1)
      await stake(staker2, stakeAmount2)

      const initTimestamp = await lockTokens(rewardsAmount, durationSec)

      // Pass the time to the middle of the bonus period.
      const passedPeriod = bonusPeriodSec.divn(2)
      await time.increaseTo(initTimestamp.add(passedPeriod))

      await checkRewards(
        staker1,
        stakeAmount1,
        expectedRewards1,
        `invalid calculated rewards for staker 1`
      )
      await checkRewards(
        staker2,
        stakeAmount2,
        expectedRewards2,
        `invalid calculated rewards for staker 2`
      )
    })

    it("should withdraw stake and rewards", async () => {
      const expectedRewards1 = toBN(125e2).mul(tokenDecimalMultiplier) // (16k / (16k + 112k)) * 100k = 12.5k KEEP
      const expectedRewards2 = toBN(875e2).mul(tokenDecimalMultiplier) // (112k / (16k + 112k)) * 100k = 87.5k KEEP

      await stake(staker1, stakeAmount1)
      await stake(staker2, stakeAmount2)

      const lockTimestamp1 = await lockTokens(rewardsAmount, durationSec)

      // End first interval.
      await time.increaseTo(lockTimestamp1.add(durationSec))

      await tokenGeyser.unstake(stakeAmount1, [], {
        from: staker1,
      })
      await tokenGeyser.unstake(stakeAmount2, [], {
        from: staker2,
      })

      // Validate stakers' stake token balances.
      expect(
        await stakeToken.balanceOf.call(staker1),
        "invalid staker's 1 token balance"
      ).to.eq.BN(stakerInitialBalance1)
      expect(
        await stakeToken.balanceOf.call(staker2),
        "invalid staker's 2 token balance"
      ).to.eq.BN(stakerInitialBalance2)

      // Validate stakers' distribution token balances.
      expect(
        await keepToken.balanceOf.call(staker1),
        "invalid staker's 1 rewards token balance"
      ).to.eq.BN(expectedRewards1)
      expect(
        await keepToken.balanceOf.call(staker2),
        "invalid staker's 2 rewards token balance"
      ).to.eq.BN(expectedRewards2)
    })

    async function checkRewards(staker, stakeAmount, expectedRewards, message) {
      const actualRewards = await tokenGeyser.unstakeQuery.call(stakeAmount, {
        from: staker,
      })

      expectCloseTo(actualRewards, expectedRewards, message)
    }
  })

  async function lockTokens(amount, durationSec) {
    await keepToken.approve(tokenGeyser.address, amount, {
      from: rewardDistribution,
    })

    const { receipt } = await tokenGeyser.lockTokens(amount, durationSec, {
      from: rewardDistribution,
    })

    const timestamp = toBN(
      (await web3.eth.getBlock(receipt.blockNumber)).timestamp
    )

    return timestamp
  }

  async function stake(staker, amount) {
    await stakeToken.approve(tokenGeyser.address, amount, { from: staker })

    return await tokenGeyser.stake(amount, [], { from: staker })
  }

  function expectCloseTo(actual, expected, message) {
    actualBN = toBN(actual)
    expectedBN = toBN(expected)

    const delta = actualBN.muln(1).divn(100) // approx. 1%

    if (
      actualBN.lt(expectedBN.sub(delta)) ||
      actualBN.gt(expectedBN.add(delta))
    ) {
      expect.fail(
        `${message}\nexpected : ${expectedBN.toString()}\nactual   : ${actualBN.toString()}`
      )
    }
  }
})
