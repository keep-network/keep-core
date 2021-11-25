/* eslint-disable @typescript-eslint/no-unused-expressions, no-await-in-loop */

import { ethers, helpers, waffle } from "hardhat"
import { expect } from "chai"
import { submissionDeployment } from "./fixtures"
import type { Submission } from "../typechain"
import { OperatorID } from "./utils/operators"
import blsData from "./data/bls"

const { mineBlocks } = helpers.time
const testSeed = blsData.groupSignatureUint256
const testGroupSize = 8
const testEligibilityDelay = 10

describe("Submission", () => {
  let submission: Submission

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(submissionDeployment)

    submission = contracts.submission as Submission
  })

  describe("isEligible", () => {
    let block: number

    it("should correctly manage the eligibility queue", async () => {
      block = (await ethers.provider.getBlock("latest")).number

      // At the beginning only member 8 is eligible because
      // (testSeed % groupSize) + 1 = 8.
      await assertMembersEligible(block, [8])
      await assertMembersNotEligible(block, [1, 2, 3, 4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1])
      await assertMembersNotEligible(block, [2, 3, 4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1, 2])
      await assertMembersNotEligible(block, [3, 4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1, 2, 3])
      await assertMembersNotEligible(block, [4, 5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1, 2, 3, 4])
      await assertMembersNotEligible(block, [5, 6, 7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1, 2, 3, 4, 5])
      await assertMembersNotEligible(block, [6, 7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1, 2, 3, 4, 5, 6])
      await assertMembersNotEligible(block, [7])

      await mineBlocks(10)

      await assertMembersEligible(block, [8, 1, 2, 3, 4, 5, 6, 7])
    })
  })

  describe("getInactiveMembers", () => {
    let groupMembers: OperatorID[]

    beforeEach(async () => {
      groupMembers = [1, 2, 3, 4, 5, 6, 7, 8]
    })

    context("when submitter index is the first eligible index", () => {
      it("should return empty inactive members list", async () => {
        const inactiveMembers = await submission.getInactiveMembers(
          5,
          5,
          groupMembers
        )

        await expect(inactiveMembers.length).to.be.equal(0)
      })
    })

    context("when submitter index is bigger than first eligible index", () => {
      it("should return a proper inactive members list", async () => {
        const inactiveMembers = await submission.getInactiveMembers(
          8,
          5,
          groupMembers
        )

        await expect(inactiveMembers.length).to.be.equal(3)
        await expect(inactiveMembers[0]).to.be.equal(groupMembers[4])
        await expect(inactiveMembers[1]).to.be.equal(groupMembers[5])
        await expect(inactiveMembers[2]).to.be.equal(groupMembers[6])
      })
    })

    context("when submitter index is smaller than first eligible index", () => {
      it("should return a proper inactive members list", async () => {
        const inactiveMembers = await submission.getInactiveMembers(
          3,
          5,
          groupMembers
        )

        await expect(inactiveMembers.length).to.be.equal(6)
        await expect(inactiveMembers[0]).to.be.equal(groupMembers[4])
        await expect(inactiveMembers[1]).to.be.equal(groupMembers[5])
        await expect(inactiveMembers[2]).to.be.equal(groupMembers[6])
        await expect(inactiveMembers[3]).to.be.equal(groupMembers[7])
        await expect(inactiveMembers[4]).to.be.equal(groupMembers[0])
        await expect(inactiveMembers[5]).to.be.equal(groupMembers[1])
      })
    })
  })

  async function assertMembersEligible(
    protocolStartBlock: number,
    checkedMembers: number[]
  ) {
    const [firstEligibleIndex, lastEligibleIndex] =
      await submission.getEligibilityRange(
        testSeed,
        protocolStartBlock,
        testEligibilityDelay,
        testGroupSize
      )

    for (let i = 0; i < checkedMembers.length; i++) {
      expect(
        await submission.isEligible(
          checkedMembers[i],
          firstEligibleIndex,
          lastEligibleIndex
        )
      ).to.be.true
    }
  }

  async function assertMembersNotEligible(
    protocolStartBlock: number,
    checkedMembers: number[]
  ) {
    const [firstEligibleIndex, lastEligibleIndex] =
      await submission.getEligibilityRange(
        testSeed,
        protocolStartBlock,
        testEligibilityDelay,
        testGroupSize
      )

    for (let i = 0; i < checkedMembers.length; i++) {
      expect(
        await submission.isEligible(
          checkedMembers[i],
          firstEligibleIndex,
          lastEligibleIndex
        )
      ).to.be.false
    }
  }
})
