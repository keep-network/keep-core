import { ethers, waffle } from "hardhat"
import { expect } from "chai"

import { Signer, Contract } from "ethers"
import { randomBeaconDeployment } from "./fixtures"

import type { RandomBeacon, SortitionPoolStub } from "../typechain"

describe("RandomBeacon - Pool", () => {
  let operator: Signer
  let randomBeacon: Contract
  let sortitionPoolStub: Contract

  // prettier-ignore
  before(async () => {
    [operator] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(randomBeaconDeployment)

    sortitionPoolStub = contracts.sortitionPoolStub as SortitionPoolStub
    randomBeacon = contracts.randomBeacon as RandomBeacon
  })

  describe("registerMemberCandidate", () => {
    context("when the operator is not registered yet", () => {
      beforeEach(async () => {
        await randomBeacon.connect(operator).registerMemberCandidate()
      })

      it("should register the operator", async () => {
        await expect(
          await sortitionPoolStub.operators(await operator.getAddress())
        ).to.be.true
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
        await expect(
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
          await expect(
            await randomBeacon.isOperatorEligible(await operator.getAddress())
          ).to.be.false
        })
      }
    )
  })
})
