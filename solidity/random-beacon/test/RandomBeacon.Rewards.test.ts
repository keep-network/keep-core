/* eslint-disable @typescript-eslint/no-extra-semi */
import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { to1e18 } from "@keep-network/hardhat-helpers/dist/src/number"

import { constants, testDeployment } from "./fixtures"
import { registerOperators } from "./utils/operators"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Operator } from "./utils/operators"
import type { RandomBeacon, T, SortitionPool, TokenStaking } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const fixture = async () => {
  const contracts = await testDeployment()

  // Accounts offset provided to slice getUnnamedSigners have to include number
  // of unnamed accounts that were already used.
  const signers = await registerOperators(
    contracts.randomBeacon as RandomBeacon,
    contracts.t as T,
    constants.groupSize,
    1
  )

  const randomBeacon = contracts.randomBeacon as RandomBeacon
  const staking = contracts.staking as TokenStaking
  const sortitionPool = contracts.sortitionPool as SortitionPool
  const t = contracts.t as T

  return {
    randomBeacon,
    staking,
    sortitionPool,
    t,
    signers,
  }
}

describe("RandomBeacon - Rewards", () => {
  let thirdParty: SignerWithAddress
  let signers: Operator[]

  let randomBeacon: RandomBeacon
  let staking: TokenStaking
  let sortitionPool: SortitionPool
  let t: T

  const rewardAmount = to1e18(100000)

  before(async () => {
    ;({ randomBeacon, staking, sortitionPool, t, signers } =
      await waffle.loadFixture(fixture))
    ;[thirdParty] = await ethers.getUnnamedSigners()
    const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

    // Allocate sortition pool rewards
    await t.connect(deployer).mint(deployer.address, rewardAmount)
    await t
      .connect(deployer)
      .approveAndCall(sortitionPool.address, rewardAmount, [])
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
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should withdraw rewards", async () => {
        const operator = signers[0].signer.address
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
})
