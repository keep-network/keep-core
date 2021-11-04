import { ethers, waffle } from "hardhat"
import { expect } from "chai"

import { Signer, Contract } from "ethers"
import { randomBeaconDeployment } from "./fixtures"

import type { RandomBeacon, SortitionPool, StakingStub } from "../typechain"

const fixture = async () => randomBeaconDeployment(undefined)

describe("RandomBeacon - Pool", () => {
  let operator: Signer
  let randomBeacon: Contract
  let sortitionPool: Contract
  let stakingStub: Contract

  // prettier-ignore
  before(async () => {
    [operator] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(fixture)
    sortitionPool = contracts.sortitionPool as SortitionPool
    stakingStub = contracts.stakingStub as StakingStub
    randomBeacon = contracts.randomBeacon as RandomBeacon
  })

  describe("registerMemberCandidate", () => {
    const minimumStake = 2000
    beforeEach(async () => {
      await stakingStub.setStake(operator.getAddress(), minimumStake)
    })

    context("when the operator is not registered yet", () => {
      beforeEach(async () => {
        await randomBeacon.connect(operator).registerMemberCandidate()
      })

      it("should register the operator", async () => {
        await expect(
          await sortitionPool.isOperatorInPool(await operator.getAddress())
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
    const minimumStake = 2000
    context("when the operator is eligible to join the sortition pool", () => {
      beforeEach(async () => {
        await stakingStub.setStake(operator.getAddress(), minimumStake)
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
          await stakingStub.setStake(operator.getAddress(), minimumStake - 1)
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
