const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {expectRevert, expectEvent} = require("@openzeppelin/test-helpers")

const KeepToken = contract.fromArtifact("KeepToken")
const PhasedEscrow = contract.fromArtifact("PhasedEscrow")

const chai = require("chai")
chai.use(require("bn-chai")(web3.utils.BN))
const expect = chai.expect

describe.only("PhasedEscrow", () => {
  const owner = accounts[0]
  const updatedOwner = accounts[1]

  const beneficiary = accounts[2]
  const updatedBeneficiary = accounts[3]

  let token
  let phasedEscrow

  before(async () => {
    token = await KeepToken.new({from: owner})
    phasedEscrow = await PhasedEscrow.new(token.address, {from: owner})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("setBeneficiary", async () => {
    it("can be called by owner", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})
      // ok, no revert
    })

    it("can be called by updated owner", async () => {
      await phasedEscrow.transferOwnership(updatedOwner, {from: owner})

      await expectRevert(
        phasedEscrow.setBeneficiary(beneficiary, {from: owner}),
        "Ownable: caller is not the owner"
      )
      await phasedEscrow.setBeneficiary(beneficiary, {from: updatedOwner})
      // ok, no revert
    })

    it("can not be called by non-owner", async () => {
      await expectRevert(
        phasedEscrow.setBeneficiary(beneficiary, {from: beneficiary}),
        "Ownable: caller is not the owner"
      )
    })

    it("sets beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})

      expect(await phasedEscrow.beneficiary()).to.equal(
        beneficiary,
        "Unexpected beneficiary"
      )
    })

    it("emits an event", async () => {
      const receipt = await phasedEscrow.setBeneficiary(beneficiary, {
        from: owner,
      })

      expectEvent(receipt, "BeneficiaryUpdated", {
        beneficiary: beneficiary,
      })
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
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: beneficiary}),
        "Ownable: caller is not the owner"
      )
    })

    it("can be called by owner", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})
      await token.transfer(phasedEscrow.address, 100, {from: owner})
      await phasedEscrow.withdraw(100, {from: owner})
      // ok, no reverts
    })

    it("fails when escrow is empty", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})
      await expectRevert(
        phasedEscrow.withdraw(100, {from: owner}),
        "Not enough tokens for withdrawal"
      )
    })

    it("withdraws specified tokens to beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})
      const amount = web3.utils.toBN(123456789)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.withdraw(100, {from: owner})

      expect(await token.balanceOf(beneficiary)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        123456789 - 100,
        "Unexpected amount withdrawn"
      )
    })

    it("withdraws specified tokens to updated beneficiary", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})

      const amount = web3.utils.toBN(987654321)
      await token.transfer(phasedEscrow.address, amount, {from: owner})

      await phasedEscrow.withdraw(100, {from: owner})

      await phasedEscrow.setBeneficiary(updatedBeneficiary, {from: owner})
      await phasedEscrow.withdraw(100, {from: owner})

      expect(await token.balanceOf(beneficiary)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(updatedBeneficiary)).to.eq.BN(
        100,
        "Unexpected amount withdrawn"
      )
      expect(await token.balanceOf(phasedEscrow.address)).to.eq.BN(
        987654321 - 200,
        "Unexpected amount withdrawn"
      )
    })

    it("emits an event", async () => {
      await phasedEscrow.setBeneficiary(beneficiary, {from: owner})
      const amount = web3.utils.toBN(100)
      await token.transfer(phasedEscrow.address, amount.muln(2), {from: owner})

      const receipt = await phasedEscrow.withdraw(amount, {from: owner})

      await expectEvent(receipt, "TokensWithdrawn", {
        beneficiary: beneficiary,
        amount: amount,
      })
    })
  })
})
