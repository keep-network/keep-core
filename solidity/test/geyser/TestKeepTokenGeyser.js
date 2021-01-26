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
  const startBonus = toBN(100) //toBN(20) // 20% // ????
  const bonusPeriodSec = toBN(1) // time.duration.days(4)
  const initialSharesPerToken = toBN(1) // ????
  const durationSec = time.duration.days(7)

  const tokenDecimalMultiplier = toBN(10e18) // 18-decimal precision

  const stakerInitialBalance1 = toBN(80e3).mul(tokenDecimalMultiplier) // 80k KEEP
  const stakeAmount1 = toBN(20e3).mul(tokenDecimalMultiplier) // 20k KEEP

  const stakeAmount2 = toBN(40e3).mul(tokenDecimalMultiplier) // 40k KEEP
  const stakerInitialBalance2 = toBN(50e3).mul(tokenDecimalMultiplier) // 50k KEEP

  const rewardsAmount = toBN(900e3).mul(tokenDecimalMultiplier) // 900k KEEP

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
      await stakeToken.approve(tokenGeyser.address, stakeAmount1, {
        from: staker1,
      })
      await tokenGeyser.stake(stakeAmount1, [], { from: staker1 })

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
      await stakeToken.approve(tokenGeyser.address, stakeAmount1, {
        from: staker1,
      })
      const receipt = await tokenGeyser.stake(stakeAmount1, [], {
        from: staker1,
      })

      expectEvent(receipt, "Staked", {
        user: staker1,
        amount: stakeAmount1,
        total: stakeAmount1,
      })
    })

    it("allows stake top-ups", async () => {
      await stakeToken.approve(
        tokenGeyser.address,
        stakeAmount1.add(stakeAmount2),
        {
          from: staker1,
        }
      )
      await tokenGeyser.stake(stakeAmount1, [], { from: staker1 })

      expect(
        await tokenGeyser.totalStakedFor(staker1),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount1)

      const receipt = await tokenGeyser.stake(stakeAmount2, [], {
        from: staker1,
      })

      expect(
        await tokenGeyser.totalStakedFor(staker1),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount1.add(stakeAmount2))

      expectEvent(receipt, "Staked", {
        user: staker1,
        amount: stakeAmount2,
        total: stakeAmount1.add(stakeAmount2),
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
    it("should calculate rewards", async () => {
      const stakeAmount1 = toBN(20e3).mul(tokenDecimalMultiplier) // 20k KEEP
      const stakeAmount2 = toBN(40e3).mul(tokenDecimalMultiplier) // 40k KEEP
      const rewardsAmount = toBN(900e3).mul(tokenDecimalMultiplier) // 900k KEEP

      const expectedRewards1 = toBN(300e3).mul(tokenDecimalMultiplier) // 300k KEEP
      const expectedRewards2 = toBN(600e3).mul(tokenDecimalMultiplier) // 600k KEEP

      await stakeTokenApprove(staker1, tokenGeyser.address, stakeAmount1)
      await stakeTokenApprove(staker2, tokenGeyser.address, stakeAmount2)
      await keepTokenApprove(
        rewardDistribution,
        tokenGeyser.address,
        rewardsAmount
      )

      await tokenGeyser.stake(stakeAmount1, [], { from: staker1 })
      await tokenGeyser.stake(stakeAmount2, [], { from: staker2 })

      const lockTokensTX = await tokenGeyser.lockTokens(
        rewardsAmount,
        durationSec,
        {
          from: rewardDistribution,
        }
      )

      const initTimestamp = new BN(
        (await web3.eth.getBlock(lockTokensTX.blockNumber)).timestamp
      )

      await time.increase(initTimestamp.add(durationSec))

      const rewards1 = await tokenGeyser.unstakeQuery.call(stakeAmount1, {
        from: staker1,
      })
      const rewards2 = await tokenGeyser.unstakeQuery.call(stakeAmount2, {
        from: staker2,
      })

      expectCloseTo(
        rewards1,
        expectedRewards1,
        "invalid calculated staker's 1 rewards"
      )
      expectCloseTo(
        rewards2,
        expectedRewards2,
        "invalid calculated staker's 2 rewards"
      )
    })

    it("should withdraw stake and rewards", async () => {
      const rewardsAmount = toBN(900e3).mul(tokenDecimalMultiplier) // 900k KEEP

      const expectedRewards1 = toBN(300e3).mul(tokenDecimalMultiplier) // (20k / (20k + 40k)) * 900k = 300k KEEP
      const expectedRewards2 = toBN(600e3).mul(tokenDecimalMultiplier) // (40k / (20k + 40k)) * 900k = 600k KEEP

      await stakeTokenApprove(staker1, tokenGeyser.address, stakeAmount1)
      await stakeTokenApprove(staker2, tokenGeyser.address, stakeAmount2)
      await keepTokenApprove(
        rewardDistribution,
        tokenGeyser.address,
        rewardsAmount
      )

      await tokenGeyser.stake(stakeAmount1, [], { from: staker1 })
      await tokenGeyser.stake(stakeAmount2, [], { from: staker2 })

      const lockTokensTX = await tokenGeyser.lockTokens(
        rewardsAmount,
        durationSec,
        {
          from: rewardDistribution,
        }
      )

      const initTimestamp = new BN(
        (await web3.eth.getBlock(lockTokensTX.blockNumber)).timestamp
      )

      await time.increase(initTimestamp.add(durationSec))

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
      expectCloseTo(
        await keepToken.balanceOf.call(staker1),
        expectedRewards1,
        "invalid staker's 1 rewards token balance"
      )
      expectCloseTo(
        await keepToken.balanceOf.call(staker2),
        expectedRewards2,
        "invalid staker's 2 rewards token balance"
      )
    })
  })

  async function keepTokenApprove(from, to, amount) {
    await keepToken.approve(to, amount, { from: from })
  }
  async function stakeTokenApprove(from, to, amount) {
    await stakeToken.approve(to, amount, { from: from })
  }

  function expectCloseTo(actual, expected, message) {
    const delta = tokenDecimalMultiplier

    if (actual.lt(expected.sub(delta)) || actual.gt(expected.add(delta))) {
      expect.fail(
        `${message}\nexpected : ${expected.toString()}\nactual   : ${actual.toString()}`
      )
    }
  }
})
