const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")

const KeepTokenGeyser = contract.fromArtifact("KeepTokenGeyser")
const KeepToken = contract.fromArtifact("KeepToken")

const BN = web3.utils.BN
const chai = require("chai")
chai.use(require("bn-chai")(BN))
const expect = chai.expect

describe("KeepTokenGeyser", () => {
  const contractOwner = accounts[0]
  const user = accounts[1]

  const maxUnlockSchedules = 12 // ????
  const startBonus = 5 // ????
  const bonusPeriodSec = 63113904 // ???? 24 months
  const initialSharesPerToken = 1 // ????

  const tokenDecimalMultiplier = web3.utils.toBN(10e18) // 18-decimal precision

  const userInitialBalance = web3.utils.toBN(80e3).mul(tokenDecimalMultiplier) // 80k KEEP with 18-decimal precision

  let token
  let tokenGeyser

  before(async () => {
    token = await KeepToken.new({ from: contractOwner })
    tokenGeyser = await KeepTokenGeyser.new(
      token.address,
      token.address,
      maxUnlockSchedules,
      startBonus,
      bonusPeriodSec,
      initialSharesPerToken,
      {
        from: contractOwner,
      }
    )

    await token.approve(user, userInitialBalance, { from: contractOwner })
    await token.transfer(user, userInitialBalance, { from: contractOwner })
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("stake", async () => {
    const stakeAmount = web3.utils.toBN(20e3).mul(tokenDecimalMultiplier)

    it("should update balances", async () => {
      await token.approve(tokenGeyser.address, stakeAmount, {
        from: user,
      })
      await tokenGeyser.stake(stakeAmount, [], { from: user })

      expect(
        await token.balanceOf.call(user),
        "invalid user's token balance"
      ).to.eq.BN(userInitialBalance.sub(stakeAmount))

      expect(
        await tokenGeyser.totalStakedFor(user),
        "invalid user's token geyser staked balance"
      ).to.eq.BN(stakeAmount)
    })
  })
})
