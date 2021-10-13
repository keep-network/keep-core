import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import { blsData } from "./data/bls"
import { constants, params, testDeployment } from "./fixtures"

import type { ContractTransaction, Signer } from "ethers"
import type { RandomBeacon } from "../typechain"
import type { Address } from "hardhat-deploy/types"

const { mineBlocks } = helpers.time

describe("RandomBeacon contract", function () {
  const dkgTimeout: number =
    constants.offchainDkgTime +
    constants.groupSize * params.dkgResultSubmissionEligibilityDelay +
    params.dkgResultChallengePeriodLength

  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let thirdParty: Signer
  let signers: DkgGroupSigners

  let randomBeacon: RandomBeacon

  before(async function () {
    thirdParty = await ethers.getSigner((await getUnnamedAccounts())[1])

    // Accounts offset provided to getDkgGroupSigners have to include number of
    // unnamed accounts that were already used.
    signers = await getDkgGroupSigners(constants.groupSize, 1)
  })

  beforeEach("load test fixture", async function () {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
  })

  describe("genesis function call", async function () {
    it("can be invoked by third party", async function () {
      await randomBeacon.connect(thirdParty).genesis()
    })

    context("with initial contract state", async function () {
      let tx: ContractTransaction
      let expectedSeed: string

      beforeEach("run genesis", async () => {
        ;[tx, expectedSeed] = await genesis()
      })

      it("emits DkgStarted event", async function () {
        await expect(tx)
          .to.emit(randomBeacon, "DkgStarted")
          .withArgs(expectedSeed)
      })
    })

    context("with genesis in progress", async function () {
      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
      })

      context("with dkg result not submitted", async function () {
        it("reverts with dkg is currently in progress error", async function () {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "dkg is currently in progress"
          )
        })
      })

      context("with dkg result submitted", async function () {
        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)

          await signAndSubmitDkgResult(signers)
        })

        it("reverts with dkg is currently in progress error", async function () {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "dkg is currently in progress"
          )
        })
      })
    })
  })

  describe("isDkgInProgress function call", async function () {
    context("with initial contract state", async function () {
      it("returns false", async function () {
        expect(await randomBeacon.isDkgInProgress()).to.be.false
      })
    })

    context("when genesis dkg started", async function () {
      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
      })

      it("returns true", async function () {
        expect(await randomBeacon.isDkgInProgress()).to.be.true
      })

      context("when genesis dkg result was submitted", async function () {
        let dkgResult: DkgResult

        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)
          ;({ dkgResult } = await signAndSubmitDkgResult(signers))
        })

        context("when genesis dkg result was not approved", async function () {
          it("returns true", async function () {
            expect(await randomBeacon.isDkgInProgress()).to.be.true
          })
        })
      })
    })
  })

  describe("submitDkgResult function call", async function () {
    // TODO: Add more tests to cover the DKG result verification function thoroughly.

    context("with initial contract state", async function () {
      it("reverts with dkg is currently not in progress error", async function () {
        await expect(signAndSubmitDkgResult(signers)).to.be.revertedWith(
          "dkg is currently not in progress"
        )
      })
    })

    context("with group creation in progress", async function () {
      let expectedSeed: string

      beforeEach("run genesis", async () => {
        ;[, expectedSeed] = await genesis()
      })

      context("with group creation not timed out", async function () {
        let tx: ContractTransaction
        let dkgResult: DkgResult

        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)
          ;({
            transaction: tx,

            dkgResult
          } = await signAndSubmitDkgResult(signers))
        })

        it("emits DkgResultSubmitted event", async function () {
          await expect(tx)
            .to.emit(randomBeacon, "DkgResultSubmitted")
            .withArgs(dkgResult.groupPubKey, signers.get(1))
        })
      })

      context("with group creation timed out", async function () {
        beforeEach("increase time", async () => {
          await mineBlocks(dkgTimeout)
        })

        context("with timeout not notified", async function () {
          it("succeeds", async function () {
            await signAndSubmitDkgResult(signers)
          })
        })
      })
    })
  })

  async function genesis(): Promise<[ContractTransaction, string]> {
    const tx = await randomBeacon.genesis()

    const expectedSeed = ethers.utils.keccak256(
      ethers.utils.solidityPack(
        ["uint256", "uint256"],
        [await randomBeacon.genesisSeed(), tx.blockNumber]
      )
    )

    return [tx, expectedSeed]
  }

  interface DkgResult {
    submitterMemberIndex: number
    groupPubKey: string
    misbehaved: string
    signatures: string
    signingMemberIndices: number[]
    members: string[]
  }

  async function signAndSubmitDkgResult(
    signers: DkgGroupSigners
  ): Promise<{
    transaction: ContractTransaction

    dkgResult: DkgResult
  }> {
    const noMisbehaved = "0x"

    expect(signers.size, "unexpected signers map size").to.be.equal(
      constants.groupSize
    )

    const {
      members,
      signingMemberIndices,
      signaturesBytes
    } = await signDkgResult(signers, groupPublicKey, noMisbehaved)

    const submitterIndex = 1

    const dkgResult = {
      submitterMemberIndex: submitterIndex,
      groupPubKey: blsData.groupPubKey,
      misbehaved: noMisbehaved,
      signatures: signaturesBytes,
      signingMemberIndices: signingMemberIndices,
      members: members
    }

    const transaction = await randomBeacon
      .connect(await ethers.getSigner(signers.get(submitterIndex)))
      .submitDkgResult(dkgResult)

    return { transaction, dkgResult }
  }
})

interface DkgGroupSigners extends Map<number, Address> {}

async function getDkgGroupSigners(
  groupSize: number,
  startAccountsOffset: number
): Promise<DkgGroupSigners> {
  const signers = new Map<number, Address>()

  for (let i = 1; i <= groupSize; i++) {
    const signer = (await getUnnamedAccounts())[startAccountsOffset + i]

    expect(
      signer,
      `signer [${i}] is not defined; check hardhat network configuration`
    ).is.not.empty

    signers.set(i, signer)
  }

  return signers
}

async function signDkgResult(
  signers: DkgGroupSigners,
  groupPublicKey: string,
  misbehaved: string
) {
  const resultHash = ethers.utils.solidityKeccak256(
    ["bytes", "bytes"],
    [groupPublicKey, misbehaved]
  )

  const members: string[] = []
  const signingMemberIndices: number[] = []
  const signatures: string[] = []

  for (let [memberIndex, signer] of signers) {
    members.push(signer)

    signingMemberIndices.push(memberIndex)

    const ethersSigner = await ethers.getSigner(signer)

    const signature = await ethersSigner.signMessage(
      ethers.utils.arrayify(resultHash)
    )

    signatures.push(signature)
  }

  expect(
    signingMemberIndices.length,
    "unexpected signingMemberIndices array size"
  ).to.be.equal(signers.size)

  expect(signatures.length, "unexpected signatures array size").to.be.equal(
    signers.size
  )

  const signaturesBytes: string = ethers.utils.hexConcat(signatures)

  return { members, signingMemberIndices, signaturesBytes }
}
