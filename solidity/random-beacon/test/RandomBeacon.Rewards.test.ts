/* eslint-disable @typescript-eslint/no-extra-semi */
import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { to1e18 } from "@keep-network/hardhat-helpers/dist/src/number"

import { constants, testDeployment } from "./fixtures"
import { registerOperators } from "./utils/operators"
import { createGroup } from "./utils/groups"
import { signOperatorInactivityClaim } from "./utils/inacvitity"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Operator } from "./utils/operators"
import type {
  RandomBeacon,
  T,
  SortitionPool,
  TokenStaking,
  RandomBeaconGovernance,
} from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const fixture = async () => {
  const contracts = await testDeployment()

  // Accounts offset provided to slice getUnnamedSigners have to include number
  // of unnamed accounts that were already used.
  const operators = await registerOperators(
    contracts.randomBeacon as RandomBeacon,
    contracts.t as T,
    constants.groupSize,
    1
  )

  const randomBeacon = contracts.randomBeacon as RandomBeacon
  const randomBeaconGovernance =
    contracts.randomBeaconGovernance as RandomBeaconGovernance
  const staking = contracts.staking as TokenStaking
  const sortitionPool = contracts.sortitionPool as SortitionPool
  const t = contracts.t as T

  return {
    randomBeacon,
    randomBeaconGovernance,
    staking,
    sortitionPool,
    t,
    operators,
  }
}

describe("RandomBeacon - Rewards", () => {
  let deployer: SignerWithAddress
  let governance: SignerWithAddress
  let thirdParty: SignerWithAddress
  let operators: Operator[]

  let randomBeacon: RandomBeacon
  let randomBeaconGovernance: RandomBeaconGovernance
  let staking: TokenStaking
  let sortitionPool: SortitionPool
  let t: T

  const rewardAmount = to1e18(100000)

  before(async () => {
    ;({
      randomBeacon,
      randomBeaconGovernance,
      staking,
      sortitionPool,
      t,
      operators,
    } = await waffle.loadFixture(fixture))
    ;[thirdParty] = await ethers.getUnnamedSigners()

    deployer = await ethers.getNamedSigner("deployer")
    governance = await ethers.getNamedSigner("governance")
  })

  describe("withdrawRewards", () => {
    context("when called for an unknown operator", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon.withdrawRewards(thirdParty.address)
        ).to.be.revertedWith("Unknown operator")
      })
    })

    context("when called for a known operator", () => {
      before(async () => {
        await createSnapshot()

        // Allocate sortition pool rewards
        await t.connect(deployer).mint(deployer.address, rewardAmount)
        await t
          .connect(deployer)
          .approveAndCall(sortitionPool.address, rewardAmount, [])
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should withdraw rewards", async () => {
        const operator = operators[0].signer.address
        const stakingProvider = await randomBeacon.operatorToStakingProvider(
          operator
        )
        const { beneficiary } = await staking.rolesOf(stakingProvider)

        expect(await t.balanceOf(beneficiary)).to.equal(0)
        await randomBeacon.withdrawRewards(stakingProvider)
        expect(await t.balanceOf(beneficiary)).to.be.gt(0)
      })
    })
  })

  describe("withdrawIneligibleRewards", () => {
    context("when called not by the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .withdrawIneligibleRewards(thirdParty.address)
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when called by the governance", () => {
      before(async () => {
        await createSnapshot()

        await createGroup(randomBeacon, operators)
        const groupId = 0
        const group = await randomBeacon["getGroup(uint64)"](groupId)

        const inactiveMembersIndices = [1, 2, 3]

        const { signatures, signingMembersIndices } =
          await signOperatorInactivityClaim(
            operators,
            0,
            group.groupPubKey,
            inactiveMembersIndices,
            33
          )

        const claimSender = operators[0].signer
        await randomBeacon.connect(claimSender).notifyOperatorInactivity(
          {
            groupId,
            inactiveMembersIndices,
            signatures,
            signingMembersIndices,
          },
          0,
          operators.map((operator) => operator.id)
        )

        // Allocate sortition pool rewards
        await t.connect(deployer).mint(deployer.address, rewardAmount)
        await t
          .connect(deployer)
          .approveAndCall(sortitionPool.address, rewardAmount, [])
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should withdraw ineligible rewards", async () => {
        // Withdraw rewards for ineligible operator. This action recalculates
        // the balance of "ineligible rewards" available for withdrawal from
        // the Sortition Pool
        const operator = operators[0].signer.address
        const stakingProvider = await randomBeacon.operatorToStakingProvider(
          operator
        )
        await randomBeacon.withdrawRewards(stakingProvider)

        expect(await t.balanceOf(thirdParty.address)).to.equal(0)
        await randomBeaconGovernance
          .connect(governance)
          .withdrawIneligibleRewards(thirdParty.address)
        expect(await t.balanceOf(thirdParty.address)).to.be.gt(0)
      })
    })
  })
})
