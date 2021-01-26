const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const {
  expectRevert,
  expectEvent,
  time,
} = require("@openzeppelin/test-helpers")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")

const { toBN } = require("web3-utils")

const KeepTokenGeyser = contract.fromArtifact("KeepTokenGeyser")
const KeepToken = contract.fromArtifact("KeepToken")
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

  const stakerInitialBalance = toBN(80e3).mul(tokenDecimalMultiplier) // 80k KEEP with 18-decimal precision
  const stakeAmount = toBN(20e3).mul(tokenDecimalMultiplier) // 20k KEEP
  const rewardsAmount = toBN(800e3).mul(tokenDecimalMultiplier) // 800k KEEP

  let token

  let tokenGeyser

  before(async () => {
    token = await KeepToken.new({ from: contractOwner })

    tokenGeyser = await KeepTokenGeyser.new(
      token.address,
      token.address, // TODO: Use two different contracts
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
    await token.transfer(staker1, stakerInitialBalance, { from: contractOwner })
    await token.transfer(staker2, stakerInitialBalance, { from: contractOwner })
    await token.transfer(rewardDistribution, rewardsAmount, {
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
        const escrow = await BatchedPhasedEscrow.new(token.address, {
          from: contractOwner,
        })

        // Configure escrow beneficiary.
        const escrowBeneficiary = await KeepTokenGeyserRewardsEscrowBeneficiary.new(
          token.address,
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
        await token.approveAndCall(escrow.address, rewardsAmount, [], {
          from: contractOwner,
        })

        // Initiate withdraw.
        const initialEscrowBalance = await token.balanceOf.call(escrow.address)

        await escrow.batchedWithdraw(
          [escrowBeneficiary.address],
          [rewardsAmount],
          { from: contractOwner }
        )

        expect(
          await token.balanceOf.call(escrow.address),
          "invalid staker's token balance"
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
            token.address,
            []
          ),
          "Caller is not the reward distribution"
        )
      })
    })
  })

  describe("stake", async () => {
    it("should update balances", async () => {
      await token.approve(tokenGeyser.address, stakeAmount, {
        from: staker1,
      })
      await tokenGeyser.stake(stakeAmount, [], { from: staker1 })

      expect(
        await token.balanceOf.call(staker1),
        "invalid staker's token balance"
      ).to.eq.BN(stakerInitialBalance.sub(stakeAmount))

      expect(
        await tokenGeyser.totalStakedFor(staker1),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount)
    })
  })

  describe("stakeFor", async () => {
    it("should update balances", async () => {
      const initialContractOwnerBalance = await token.balanceOf.call(
        contractOwner
      )

      await token.approve(tokenGeyser.address, stakeAmount, {
        from: contractOwner,
      })
      await tokenGeyser.stakeFor(beneficiary, stakeAmount, [], {
        from: contractOwner,
      })

      expect(
        await token.balanceOf.call(staker1),
        "invalid staker's token balance"
      ).to.eq.BN(stakerInitialBalance)

      expect(
        await token.balanceOf.call(contractOwner),
        "invalid contract owner's token balance"
      ).to.eq.BN(initialContractOwnerBalance.sub(stakeAmount))

      expect(
        await tokenGeyser.totalStakedFor(contractOwner),
        "invalid staker's staked balance"
      ).to.eq.BN(0)

      expect(
        await tokenGeyser.totalStakedFor(beneficiary),
        "invalid staker's staked balance"
      ).to.eq.BN(stakeAmount)
    })
  })
})
