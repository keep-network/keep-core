import { ethers } from "hardhat"
import { Signer, Contract } from "ethers"
import { expect } from "chai"

describe("RandomBeacon - Pool", () => {
  let governance: Signer
  let operator: Signer
  let randomBeacon: Contract
  let sortitionPoolStub: Contract

  beforeEach(async () => {
    const signers = await ethers.getSigners()
    governance = signers[0]
    operator = signers[1]

    const SortitionPoolStub = await ethers.getContractFactory(
      "SortitionPoolStub"
    )
    sortitionPoolStub = await SortitionPoolStub.deploy()
    await sortitionPoolStub.deployed()

    const RandomBeacon = await ethers.getContractFactory("RandomBeacon")
    randomBeacon = await RandomBeacon.deploy(sortitionPoolStub.address)
    await randomBeacon.deployed()
  })

  describe("registerMemberCandidate", () => {
    context("when the operator is not registered yet", () => {
      beforeEach(async () => {
        await randomBeacon.connect(operator).registerMemberCandidate()
      })

      it("should register the operator", async () => {
        expect(await sortitionPoolStub.operators(await operator.getAddress()))
          .to.be.true
      })
    })

    context("when the operator is already registered", () => {
      beforeEach(async () => {
        await randomBeacon.connect(operator).registerMemberCandidate()
      })

      it("should revert", async () => {
        await expect(
          randomBeacon.connect(operator).registerMemberCandidate()
        ).to.be.revertedWith("Operator is already registered")
      })
    })
  })

  describe("isOperatorEligible", () => {
    context("when the operator is eligible to join the sortition pool", () => {
      beforeEach(async () => {
        await sortitionPoolStub.setOperatorEligibility(
          await operator.getAddress(),
          true
        )
      })

      it("should return true", async () => {
        expect(
          await randomBeacon.isOperatorEligible(await operator.getAddress())
        ).to.be.true
      })
    })

    context(
      "when the operator is not eligible to join the sortition pool",
      () => {
        beforeEach(async () => {
          await sortitionPoolStub.setOperatorEligibility(
            await operator.getAddress(),
            false
          )
        })

        it("should return false", async () => {
          expect(
            await randomBeacon.isOperatorEligible(await operator.getAddress())
          ).to.be.false
        })
      }
    )
  })
})
