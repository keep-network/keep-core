const {initContracts} = require("./helpers/initContracts")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {expectRevert, expectEvent} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = require("@openzeppelin/test-helpers/src/constants")
const time = require("@openzeppelin/test-helpers/src/time")
const crypto = require("crypto")

const KeepToken = contract.fromArtifact("KeepToken")
const Escrow = contract.fromArtifact("Escrow")
const PhasedEscrow = contract.fromArtifact("PhasedEscrow")
const BatchedPhasedEscrow = contract.fromArtifact("BatchedPhasedEscrow")

const TestSimpleBeneficiary = contract.fromArtifact("TestSimpleBeneficiary")
const StakingPoolRewardsEscrowBeneficiary = contract.fromArtifact(
  "StakingPoolRewardsEscrowBeneficiary"
)
const StakerRewardsBeneficiary = contract.fromArtifact(
  "StakerRewardsBeneficiary"
)
const BeaconBackportRewardsEscrowBeneficiary = contract.fromArtifact(
  "BeaconBackportRewardsEscrowBeneficiary"
)
const BeaconRewardsEscrowBeneficiary = contract.fromArtifact(
  "BeaconRewardsEscrowBeneficiary"
)

const BeaconRewards = contract.fromArtifact("BeaconRewards")
const BeaconBackportRewards = contract.fromArtifact("BeaconBackportRewards")
const TestCurveRewards = contract.fromArtifact("TestCurveRewards")
const TestSimpleStakerRewards = contract.fromArtifact("TestSimpleStakerRewards")

const chai = require("chai")
chai.use(require("bn-chai")(web3.utils.BN))
const expect = chai.expect

describe("PhasedEscrow", () => {
  const owner = accounts[0]
  const updatedOwner = accounts[1]

  let beneficiary
  let updatedBeneficiary
  let rewardsBeneficiary

  let token
  let phasedEscrow

  before(async () => {
    token = await KeepToken.new({from: owner})
    phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
    beneficiary = await TestSimpleBeneficiary.new()
    updatedBeneficiary = await TestSimpleBeneficiary.new()
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("receiveApproval", async () => {
    it("fails for an unknown token", async () => {
      // It is another KeepToken contract deployment, not the one PhasedEscrow
      // has been created with.
      const unknownToken = await KeepToken.new({from: owner})
      const amountApproved = web3.utils.toBN(9991)

      await expectRevert(
        unknownToken.approveAndCall(
          phasedEscrow.address,
          amountApproved,
          "0x0",
          {from: owner}
        ),
        "Unsupported token"
      )
    })

    it("transfers all approved tokens", async () => {
      const amountApproved = web3.utils.toBN(9993)
      await token.approveAndCall(phasedEscrow.address, amountApproved, "0x0", {
        from: owner,
      })

      const actualBalance = await token.balanceOf(phasedEscrow.address)
      expect(actualBalance).to.eq.BN(amountApproved)
    })
  })

  describe("setBeneficiary", async () => {
    it("can be called by owner", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      // ok, no revert
    })

    it("can be called by updated owner", async () => {
      await phasedEscrow.transferOwnership(updatedOwner, {from: owner})

      await expectRevert(
        phasedEscrow.setBeneficiary(beneficiary.address, {from: owner}),
        "Ownable: caller is not the owner"
      )
      await phasedEscrow.setBeneficiary(beneficiary.address, {
        from: updatedOwner,
      })
      // ok, no revert
    })

    it("can not be called by non-owner", async () => {
      await expectRevert(
        phasedEscrow.setBeneficiary(beneficiary.address, {from: updatedOwner}),
        "Ownable: caller is not the owner"
      )
    })

    it("sets beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})

      expect(await phasedEscrow.beneficiary()).to.equal(
        beneficiary.address,
        "Unexpected beneficiary"
      )
    })

    it("emits an event", async () => {
      const receipt = await phasedEscrow.setBeneficiary(beneficiary.address, {
        from: owner,
      })

      expectEvent(receipt, "BeneficiaryUpdated", {
        beneficiary: beneficiary.address,
      })
    })
  })

  describe("withdrawFromEscrow", async () => {
    const amount = web3.utils.toBN(12090)
    let escrow

    beforeEach(async () => {
      escrow = await Escrow.new(token.address, {from: owner})
      await token.transfer(escrow.address, amount, {from: owner})
    })

    it("pulls all funds from a non-phased Escrow when having no tokens", async () => {
      const balanceBefore = await token.balanceOf(phasedEscrow.address)
      expect(balanceBefore).to.eq.BN(0)

      await escrow.setBeneficiary(phasedEscrow.address, {from: owner})
      await phasedEscrow.withdrawFromEscrow(escrow.address)

      const balanceAfter = await token.balanceOf(phasedEscrow.address)
      expect(balanceAfter).to.eq.BN(amount)
    })

    it("pulls all funds from a non-phased Escrow when having some tokens", async () => {
      const initialFunds = web3.utils.toBN(999)
      await token.transfer(phasedEscrow.address, initialFunds, {from: owner})
      const balanceBefore = await token.balanceOf(phasedEscrow.address)
      expect(balanceBefore).to.eq.BN(initialFunds)

      await escrow.setBeneficiary(phasedEscrow.address, {from: owner})
      await phasedEscrow.withdrawFromEscrow(escrow.address)

      const balanceAfter = await token.balanceOf(phasedEscrow.address)
      expect(balanceAfter).to.eq.BN(initialFunds.add(amount))
    })
  })

  describe("withdraw", async () => {
    it("can not be called if beneficiary wasn't set", async () => {
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: owner}),
        "Beneficiary not assigned"
      )
    })

    it("can not be called by non-owner", async () => {
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: beneficiary.address}),
        "Ownable: caller is not the owner"
      )
    })

    it("can be called by owner", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await phasedEscrow.withdraw(100, {from: owner})
      // ok, no reverts
    })

    it("fails when escrow is empty", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: owner}),
        "Not enough tokens for withdrawal"
      )
    })

    it("withdraws specified tokens to updated beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      const amount = web3.utils.toBN(987654321)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.withdraw(100, {from: owner})

      await phasedEscrow.setBeneficiary(updatedBeneficiary.address, {
        from: owner,
      })
      await phasedEscrow.withdraw(100, {from: owner})

      expect(await token.balanceOf(beneficiary.address)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(updatedBeneficiary.address)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        987654321 - 200,
        "Unexpected amount withdrawn"
      )
    })

    it("withdraws specified tokens to beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      const amount = web3.utils.toBN(123456789)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.withdraw(100, {from: owner})

      expect(await token.balanceOf(beneficiary.address)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        123456789 - 100,
        "Unexpected amount withdrawn"
      )
    })

    it("emits an event", async () => {
      await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
      const amount = web3.utils.toBN(100)
      await token.transfer(phasedEscrow.address, amount.muln(2), {from: owner})

      const receipt = await phasedEscrow.withdraw(amount, {from: owner})

      await expectEvent(receipt, "TokensWithdrawn", {
        beneficiary: beneficiary.address,
        amount: amount,
      })
    })
  })

  describe("when withdrawing to a StakingPoolRewardsEscrowBeneficiary", () => {
    const baseBalance = 123456789
    const transferAmount = 100

    before(async () => {
      rewardsContract = await TestCurveRewards.new(token.address)
      rewardsBeneficiary = await StakingPoolRewardsEscrowBeneficiary.new(
        token.address,
        rewardsContract.address,
        {from: owner}
      )

      await rewardsBeneficiary.transferOwnership(phasedEscrow.address, {
        from: owner,
      })
      const amount = web3.utils.toBN(baseBalance)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.setBeneficiary(rewardsBeneficiary.address, {
        from: owner,
      })
    })

    assertRewards(baseBalance, transferAmount)

    it("emits a RewardAdded event from the rewards beneficiary", async () => {
      const receipt = resolveAllLogs(
        (await phasedEscrow.withdraw(transferAmount, {from: owner})).receipt,
        {rewardsContract}
      )

      expectEvent(receipt, "RewardAdded", {
        reward: web3.utils.toBN(transferAmount),
      })
    })
  })

  describe("when withdrawing to a BeaconBackportRewardsEscrowBeneficiary", () => {
    const baseBalance = 200000000
    const transferAmount = 200000
    let stakingContract
    let operatorContract

    before(async () => {
      const contracts = await initContracts(
        contract.fromArtifact("TokenStaking"),
        contract.fromArtifact("KeepRandomBeaconService"),
        contract.fromArtifact("KeepRandomBeaconServiceImplV1"),
        contract.fromArtifact("KeepRandomBeaconOperatorBeaconRewardsStub")
      )

      stakingContract = contracts.stakingContract
      operatorContract = contracts.operatorContract

      phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
      const amount = web3.utils.toBN(baseBalance)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      rewardsContract = await BeaconBackportRewards.new(
        token.address,
        operatorContract.address,
        stakingContract.address
      )

      rewardsBeneficiary = await BeaconBackportRewardsEscrowBeneficiary.new(
        token.address,
        rewardsContract.address,
        {from: owner}
      )
      await rewardsBeneficiary.transferOwnership(phasedEscrow.address, {
        from: owner,
      })

      await phasedEscrow.setBeneficiary(rewardsBeneficiary.address, {
        from: owner,
      })
    })

    assertRewards(baseBalance, transferAmount)
  })

  describe("when withdrawing to a BeaconRewardsEscrowBeneficiary", () => {
    const baseBalance = 200000000
    const transferAmount = 19800000
    let stakingContract
    let operatorContract

    before(async () => {
      const contracts = await initContracts(
        contract.fromArtifact("TokenStaking"),
        contract.fromArtifact("KeepRandomBeaconService"),
        contract.fromArtifact("KeepRandomBeaconServiceImplV1"),
        contract.fromArtifact("KeepRandomBeaconOperatorBeaconRewardsStub")
      )

      stakingContract = contracts.stakingContract
      operatorContract = contracts.operatorContract

      phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
      const amount = web3.utils.toBN(baseBalance)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      rewardsContract = await BeaconRewards.new(
        token.address,
        operatorContract.address,
        stakingContract.address
      )

      rewardsBeneficiary = await BeaconRewardsEscrowBeneficiary.new(
        token.address,
        rewardsContract.address,
        {from: owner}
      )
      await rewardsBeneficiary.transferOwnership(phasedEscrow.address, {
        from: owner,
      })

      await phasedEscrow.setBeneficiary(rewardsBeneficiary.address, {
        from: owner,
      })
    })

    assertRewards(baseBalance, transferAmount)
  })

  async function assertRewards(baseBalance, transferAmount) {
    it("withdraws specified tokens from escrow", async () => {
      await phasedEscrow.withdraw(transferAmount, {from: owner})

      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        baseBalance - transferAmount,
        "Unexpected amount withdrawn"
      )
    })

    it("transfers specified tokens to rewards contract", async () => {
      await phasedEscrow.withdraw(transferAmount, {from: owner})

      expect(await token.balanceOf(rewardsContract.address)).to.eq.BN(
        transferAmount,
        "Unexpected amount deposited"
      )
    })

    it("leaves no tokens in the rewards beneficiary", async () => {
      await phasedEscrow.withdraw(transferAmount, {from: owner})

      expect(await token.balanceOf(rewardsBeneficiary.address)).to.eq.BN(
        0,
        "Unexpected amount left in rewards beneficiary"
      )
    })

    it("emits a TokensWithdrawn event to the rewards beneficiary", async () => {
      const receipt = await phasedEscrow.withdraw(transferAmount, {
        from: owner,
      })

      expectEvent(receipt, "TokensWithdrawn", {
        beneficiary: rewardsBeneficiary.address,
        amount: web3.utils.toBN(transferAmount),
      })
    })
  }
})

describe("BatchedPhasedEscrow", () => {
  const owner = accounts[1]
  const drawee = accounts[2]
  const updatedOwner = accounts[3]
  const updatedDrawee = accounts[4]

  let token
  let batchedPhasedEscrow

  let beneficiary1
  let beneficiary2
  let beneficiary3

  before(async () => {
    token = await KeepToken.new({from: owner})
    batchedPhasedEscrow = await BatchedPhasedEscrow.new(token.address, {
      from: owner,
    })

    beneficiary1 = await TestSimpleBeneficiary.new({from: owner})
    beneficiary2 = await TestSimpleBeneficiary.new({from: owner})
    beneficiary3 = await TestSimpleBeneficiary.new({from: owner})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("can be funded from PhasedEscrow contract", async () => {
    // This test verifies if BatchedPhasedEscrow can be funded from PhasedEscrow.
    // To perform such operation an intermediary beneficiary contract is needed
    // that will automatically transfer funds received from PhasedEscrow to
    // BatchedPhasedEscrow.
    // The tokens are transferred in the following way:
    //   PhasedEscrow -> StakerRewardsBeneficiary -> BatchedPhasedEscrow

    // Deploy PhasedEscrow and do initial funding.
    const amount = 9000

    const phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
    await token.transfer(phasedEscrow.address, amount, {
      from: owner,
    })

    // Deploy intermediary beneficiary and transfer its' ownership to PhasedEscrow.
    const beneficiary = await StakerRewardsBeneficiary.new(
      token.address,
      batchedPhasedEscrow.address,
      {from: owner}
    )
    await beneficiary.transferOwnership(phasedEscrow.address, {
      from: owner,
    })

    // Withdraw funds from PhasedEscrow
    await phasedEscrow.setBeneficiary(beneficiary.address, {from: owner})
    await phasedEscrow.withdraw(amount, {from: owner})

    // Verify that funds got transferred to BatchedPhasedEscrow
    expect(await token.balanceOf(batchedPhasedEscrow.address)).to.eq.BN(
      amount,
      `Unexpected batched escrow balance`
    )
    expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
      0,
      `Unexpected phased escrow balance`
    )
  })

  describe("receiveApproval", async () => {
    it("fails for an unknown token", async () => {
      // It is another KeepToken contract deployment, not the one PhasedEscrow
      // has been created with.
      const unknownToken = await KeepToken.new({from: owner})
      const amountApproved = web3.utils.toBN(9991)

      await expectRevert(
        unknownToken.approveAndCall(
          batchedPhasedEscrow.address,
          amountApproved,
          "0x0",
          {from: owner}
        ),
        "Unsupported token"
      )
    })

    it("transfers all approved tokens", async () => {
      const amountApproved = web3.utils.toBN(9993)
      await token.approveAndCall(
        batchedPhasedEscrow.address,
        amountApproved,
        "0x0",
        {from: owner}
      )

      const actualBalance = await token.balanceOf(batchedPhasedEscrow.address)
      expect(actualBalance).to.eq.BN(amountApproved)
    })
  })

  describe("beneficiary approval", async () => {
    it("can be done by owner", async () => {
      await batchedPhasedEscrow.approveBeneficiary(beneficiary1.address, {
        from: owner,
      })
      // ok, no revert
    })

    it("can be done by updated owner", async () => {
      await batchedPhasedEscrow.transferOwnership(updatedOwner, {from: owner})

      await expectRevert(
        batchedPhasedEscrow.approveBeneficiary(beneficiary1.address, {
          from: owner,
        }),
        "Ownable: caller is not the owner"
      )
      await batchedPhasedEscrow.approveBeneficiary(beneficiary1.address, {
        from: updatedOwner,
      })
      // ok, no revert
    })

    it("can not be done by non-owner", async () => {
      await expectRevert(
        batchedPhasedEscrow.approveBeneficiary(beneficiary1.address, {
          from: drawee,
        }),
        "Ownable: caller is not the owner"
      )
    })

    it("can not be done on zero address", async () => {
      await expectRevert(
        batchedPhasedEscrow.approveBeneficiary(ZERO_ADDRESS, {from: owner}),
        "Beneficiary can not be zero address"
      )
    })

    it("emits an event", async () => {
      const receipt = await batchedPhasedEscrow.approveBeneficiary(
        beneficiary1.address,
        {
          from: owner,
        }
      )

      expectEvent(receipt, "BeneficiaryApproved", {
        beneficiary: beneficiary1.address,
      })
    })

    it("maintains beneficiaries as non-approved by default", async () => {
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary1.address)
      ).to.be.false
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary2.address)
      ).to.be.false
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary3.address)
      ).to.be.false
    })

    it("approves a single beneficiary", async () => {
      await batchedPhasedEscrow.approveBeneficiary(beneficiary2.address, {
        from: owner,
      })

      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary1.address)
      ).to.be.false
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary2.address)
      ).to.be.true
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary3.address)
      ).to.be.false
    })

    it("approves multiple beneficiaries ", async () => {
      await batchedPhasedEscrow.approveBeneficiary(beneficiary1.address, {
        from: owner,
      })
      await batchedPhasedEscrow.approveBeneficiary(beneficiary2.address, {
        from: owner,
      })

      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary1.address)
      ).to.be.true
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary2.address)
      ).to.be.true
      expect(
        await batchedPhasedEscrow.isBeneficiaryApproved(beneficiary3.address)
      ).to.be.false
    })
  })

  describe("drawee role", async () => {
    it("is by default assigned to owner", async () => {
      expect(await batchedPhasedEscrow.drawee()).to.equal(owner)
    })

    it("can be transferred by owner", async () => {
      await batchedPhasedEscrow.setDrawee(updatedDrawee, {from: owner})
      // ok, no revert
    })

    it("can be transferred by updated owner", async () => {
      await batchedPhasedEscrow.transferOwnership(updatedOwner, {from: owner})

      await expectRevert(
        batchedPhasedEscrow.setDrawee(updatedDrawee, {from: owner}),
        "Ownable: caller is not the owner"
      )
      await batchedPhasedEscrow.setDrawee(updatedDrawee, {
        from: updatedOwner,
      })
      // ok, no revert
    })

    it("can not be transferred by non-owner", async () => {
      await expectRevert(
        batchedPhasedEscrow.setDrawee(updatedDrawee, {from: drawee}),
        "Ownable: caller is not the owner"
      )
    })

    it("can be transferred to another account", async () => {
      let receipt = await batchedPhasedEscrow.setDrawee(drawee, {
        from: owner,
      })

      expect(await batchedPhasedEscrow.drawee()).to.equal(drawee)
      expectEvent(receipt, "DraweeRoleTransferred", {
        oldDrawee: owner,
        newDrawee: drawee,
      })

      receipt = await batchedPhasedEscrow.setDrawee(updatedDrawee, {
        from: owner,
      })

      expect(await batchedPhasedEscrow.drawee()).to.equal(updatedDrawee)
      expectEvent(receipt, "DraweeRoleTransferred", {
        oldDrawee: drawee,
        newDrawee: updatedDrawee,
      })
    })
  })

  describe("batchedWithdraw", async () => {
    let beneficiaries
    let amounts
    let escrowBalance

    beforeEach(async () => {
      beneficiaries = [
        beneficiary1.address,
        beneficiary2.address,
        beneficiary3.address,
      ]

      amounts = [100, 200, 300]
      escrowBalance = 600

      await batchedPhasedEscrow.approveBeneficiary(beneficiary1.address, {
        from: owner,
      })
      await batchedPhasedEscrow.approveBeneficiary(beneficiary2.address, {
        from: owner,
      })
      await batchedPhasedEscrow.approveBeneficiary(beneficiary3.address, {
        from: owner,
      })

      await batchedPhasedEscrow.setDrawee(drawee, {
        from: owner,
      })

      await token.transfer(batchedPhasedEscrow.address, escrowBalance, {
        from: owner,
      })
    })

    it("can be called by drawee", async () => {
      await batchedPhasedEscrow.batchedWithdraw(beneficiaries, amounts, {
        from: drawee,
      })
      // ok, no revert
    })

    it("can not be called by owner if not drawee", async () => {
      await expectRevert(
        batchedPhasedEscrow.batchedWithdraw(beneficiaries, amounts, {
          from: owner,
        }),
        "Caller is not the drawee"
      )
    })

    it("can not be called by non-drawee", async () => {
      await expectRevert(
        batchedPhasedEscrow.batchedWithdraw(beneficiaries, amounts, {
          from: updatedDrawee,
        }),
        "Caller is not the drawee"
      )
    })

    it("reverts when input arrays have different lengths", async () => {
      await expectRevert(
        batchedPhasedEscrow.batchedWithdraw(beneficiaries, [100, 200], {
          from: drawee,
        }),
        "Mismatched arrays length"
      )
    })

    it("reverts when beneficiary is not IBeneficiaryContract", async () => {
      await expectRevert.unspecified(
        batchedPhasedEscrow.batchedWithdraw(
          [beneficiary1.address, beneficiary2.address, owner],
          amounts,
          {
            from: owner,
          }
        )
      )
    })

    it("reverts when beneficiary was not approved", async () => {
      const anotherBeneficiary = await TestSimpleBeneficiary.new({from: owner})

      await expectRevert(
        batchedPhasedEscrow.batchedWithdraw(
          [beneficiary1.address, anotherBeneficiary.address],
          [100, 200],
          {
            from: drawee,
          }
        ),
        "Beneficiary was not approved"
      )
    })

    it("reverts when there are not enough funds in the escrow", async () => {
      await expectRevert.unspecified(
        batchedPhasedEscrow.batchedWithdraw(beneficiaries, [100, 200, 301], {
          from: drawee,
        })
      )
    })

    it("withdraws specified tokens to beneficiaries", async () => {
      await batchedPhasedEscrow.batchedWithdraw(beneficiaries, amounts, {
        from: drawee,
      })

      for (let i = 0; i < beneficiaries.length; i++) {
        expect(await token.balanceOf(beneficiaries[i])).to.eq.BN(
          amounts[i],
          `Unexpected amount withdrawn for beneficiary ${i}`
        )
      }

      expect(await token.balanceOf(batchedPhasedEscrow.address)).to.eq.BN(
        0,
        `Unexpected escrow balance`
      )
    })
  })
})

describe("StakingPoolRewardsEscrowBeneficiary", () => {
  const owner = accounts[0]
  const thirdParty = accounts[1]

  const transferAmount = 1000

  let token
  let rewardsContract
  let rewardsBeneficiary

  before(async () => {
    token = await KeepToken.new({from: owner})
    rewardsContract = await TestCurveRewards.new(token.address)
    rewardsBeneficiary = await StakingPoolRewardsEscrowBeneficiary.new(
      token.address,
      rewardsContract.address,
      {from: owner}
    )

    const amount = web3.utils.toBN(transferAmount)
    await token.transfer(rewardsBeneficiary.address, amount, {from: owner})
  })

  describe("__escrowSentTokens", async () => {
    it("can be called by the owner", async () => {
      await rewardsBeneficiary.__escrowSentTokens(transferAmount, {
        from: owner,
      })
      // ok, no revert
    })

    it("can not be called by the non-owner", async () => {
      await expectRevert(
        rewardsBeneficiary.__escrowSentTokens(transferAmount, {
          from: thirdParty,
        }),
        "Ownable: caller is not the owner"
      )
    })
  })
})

describe("StakerRewardsBeneficiary", () => {
  const owner = accounts[0]
  const thirdParty = accounts[1]

  const transferAmount = 1000

  let token
  let rewardsContract
  let rewardsBeneficiary

  before(async () => {
    token = await KeepToken.new({from: owner})
    rewardsContract = await TestSimpleStakerRewards.new(token.address)
    rewardsBeneficiary = await StakerRewardsBeneficiary.new(
      token.address,
      rewardsContract.address,
      {from: owner}
    )

    const amount = web3.utils.toBN(transferAmount)
    await token.transfer(rewardsBeneficiary.address, amount, {from: owner})
  })

  describe("__escrowSentTokens", async () => {
    it("can be called by the owner", async () => {
      await rewardsBeneficiary.__escrowSentTokens(transferAmount, {
        from: owner,
      })
      // ok, no revert
    })

    it("can not be called by the non-owner", async () => {
      await expectRevert(
        rewardsBeneficiary.__escrowSentTokens(transferAmount, {
          from: thirdParty,
        }),
        "Ownable: caller is not the owner"
      )
    })
  })
})

describe("BeaconRewards to PhasedEscrow transfer", async () => {
  const owner = accounts[0]
  const operators = [accounts[1], accounts[2]]

  const tokenDecimalMultiplier = web3.utils.toBN(10).pow(web3.utils.toBN(18))
  const totalRewards = web3.utils.toBN(19800000).mul(tokenDecimalMultiplier)

  let token
  let operatorContract
  let phasedEscrow
  let rewardsContract

  before(async () => {
    const contracts = await initContracts(
      contract.fromArtifact("TokenStaking"),
      contract.fromArtifact("KeepRandomBeaconService"),
      contract.fromArtifact("KeepRandomBeaconServiceImplV1"),
      contract.fromArtifact("KeepRandomBeaconOperatorBeaconRewardsStub")
    )

    token = contracts.token
    operatorContract = contracts.operatorContract
    const stakingContract = contracts.stakingContract

    phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
    rewardsContract = await BeaconRewards.new(
      token.address,
      operatorContract.address,
      stakingContract.address,
      {from: owner}
    )

    await token.approveAndCall(rewardsContract.address, totalRewards, "0x0", {
      from: owner,
    })
    await rewardsContract.markAsFunded({from: owner})
  })

  it("moves all unallocated tokens to escrow", async () => {
    await rewardsContract.initiateRewardsUpgrade(phasedEscrow.address, {
      from: owner,
    })

    const now = await time.latest()

    await operatorContract.registerNewGroup(
      crypto.randomBytes(128),
      operators,
      now
    )

    const currentInterval = await rewardsContract.intervalOf(now)
    const currentIntervalEnd = await rewardsContract.endOf(currentInterval)
    await time.increaseTo(currentIntervalEnd.addn(1))

    await rewardsContract.finalizeRewardsUpgrade({from: owner})

    const escrowBalance = await token.balanceOf(phasedEscrow.address)
    const allocatedRewards = await rewardsContract.totalRewards()
    expect(escrowBalance).to.eq.BN(totalRewards.sub(allocatedRewards))
  })
})

// FIXME Move to a shared test utils library for all Keep projects.
/**
 * Uses the ABIs of all contracts in the `contractContainer` to resolve any
 * events they may have emitted into the given `receipt`'s logs. Typically
 * Truffle only resolves the events on the calling contract; this function
 * resolves all of the ones that can be resolved.
 *
 * @param {TruffleReceipt} receipt The receipt of a contract function call
 *        submission.
 * @param {ContractContainer} contractContainer An object that contains
 *        properties that are TruffleContracts. Not all properties in the
 *        container need be contracts, nor do all contracts need to have events
 *        in the receipt.
 *
 * @return {TruffleReceipt} The receipt, with its `logs` property updated to
 *         include all resolved logs.
 */
function resolveAllLogs(receipt, contractContainer) {
  const contracts = Object.entries(contractContainer)
    .map(([, value]) => value)
    .filter((_) => _.contract && _.address)

  const {resolved: resolvedLogs} = contracts.reduce(
    ({raw, resolved}, contract) => {
      const events = contract.contract._jsonInterface.filter(
        (_) => _.type === "event"
      )
      const contractLogs = raw.filter((_) => _.address == contract.address)

      const decoded = contractLogs.map((log) => {
        const event = events.find((_) => log.topics.includes(_.signature))
        const decoded = web3.eth.abi.decodeLog(
          event.inputs,
          log.data,
          log.topics.slice(1)
        )

        return Object.assign({}, log, {
          event: event.name,
          args: decoded,
        })
      })

      return {
        raw: raw.filter((_) => _.address != contract.address),
        resolved: resolved.concat(decoded),
      }
    },
    {raw: receipt.rawLogs, resolved: []}
  )

  return Object.assign({}, receipt, {
    logs: resolvedLogs,
  })
}
