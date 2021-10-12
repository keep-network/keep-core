import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import { blsData } from "./helpers/data"
import { constants, params, testDeployment } from "./helpers/fixtures"

import type { BigNumber, ContractTransaction, Signer } from "ethers"
import type { RandomBeacon } from "../typechain"
import type { Address } from "hardhat-deploy/types"

const { mineBlocks } = helpers.time

describe("RandomBeacon contract", function () {
  const dkgTimeout: number =
    constants.timeDKG +
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
          .withArgs(
            expectedSeed,
            await randomBeacon.GROUP_SIZE(),
            await randomBeacon.dkgResultSubmissionEligibilityDelay()
          )
      })

      it("sets values", async function () {
        const dkgData = await randomBeacon.dkg()

        expect(dkgData.seed).to.eq(expectedSeed)
        expect(dkgData.groupSize).to.eq(await randomBeacon.GROUP_SIZE())
        expect(dkgData.signatureThreshold).to.eq(
          await randomBeacon.SIGNATURE_THRESHOLD()
        )
        expect(dkgData.dkgResultSubmissionEligibilityDelay).to.eq(
          await randomBeacon.dkgResultSubmissionEligibilityDelay()
        )
        expect(dkgData.dkgResultChallengePeriodLength).to.eq(
          await randomBeacon.dkgResultChallengePeriodLength()
        )
        expect(dkgData.timeDKG).to.eq(await randomBeacon.TIME_DKG())
        expect(dkgData.startBlock).to.be.eq(tx.blockNumber)
      })
    })

    context("with genesis in progress", async function () {
      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
      })

      context("with dkg result not submitted", async function () {
        it("reverts with dkg is currently in progress error", async function () {
          // TODO: It should return a dedicated error from RandomBeacon contract.
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "dkg is currently in progress"
          )
        })
      })

      context("with dkg result submitted", async function () {
        beforeEach(async () => {
          await mineBlocks(constants.timeDKG)

          await signAndSubmitDkgResult(signers)
        })

        it("reverts with dkg is currently in progress error", async function () {
          // TODO: It should return a dedicated error from RandomBeacon contract.
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "dkg is currently in progress"
          )
        })
      })
    })

    context("with genesis completed", async function () {
      context("with active group available", async function () {
        beforeEach(async () => {
          await randomBeacon.genesis()

          await mineBlocks(constants.timeDKG)

          const { resultIndex, dkgResult } = await signAndSubmitDkgResult(
            signers
          )

          await mineBlocks(
            (await randomBeacon.dkg()).dkgResultChallengePeriodLength.toNumber()
          )

          await randomBeacon.approveDkgResult(resultIndex, dkgResult)
        })

        it("reverts with not awaiting genesis error", async function () {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "not awaiting genesis"
          )
        })
      })

      // TODO: Add tests to cover scenartios of terminated, expired and pending groups.
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

      context("when dkg timeout was notified", async function () {
        beforeEach("notify dkg timeout", async () => {
          await mineBlocks(dkgTimeout)

          await randomBeacon.notifyDkgTimeout()
        })

        it("returns false", async function () {
          expect(await randomBeacon.isDkgInProgress()).to.be.false
        })
      })

      context("when genesis dkg result was submitted", async function () {
        let resultIndex: BigNumber
        let dkgResult: DkgResult

        beforeEach(async () => {
          await mineBlocks(constants.timeDKG)
          ;({ resultIndex, dkgResult } = await signAndSubmitDkgResult(signers))
        })

        context("when genesis dkg result was not approved", async function () {
          it("returns true", async function () {
            expect(await randomBeacon.isDkgInProgress()).to.be.true
          })
        })

        context("when genesis dkg result was approved", async function () {
          beforeEach(async () => {
            await mineBlocks(params.dkgResultChallengePeriodLength)

            await randomBeacon.approveDkgResult(resultIndex, dkgResult)
          })

          it("returns false", async function () {
            expect(await randomBeacon.isDkgInProgress()).to.be.false
          })
        })

        context("when genesis dkg result was challenged", async function () {
          beforeEach(async () => {
            await randomBeacon.challengeDkgResult(resultIndex, dkgResult)
          })

          it("returns true", async function () {
            expect(await randomBeacon.isDkgInProgress()).to.be.true
          })
        })
      })
    })
  })

  describe("notifyDkgTimeout function call", async function () {
    context("with group creation in progress", async function () {
      let expectedSeed: string

      beforeEach("run genesis", async () => {
        ;[, expectedSeed] = await genesis()
      })

      context("with group creation not timed out", async function () {
        it("reverts with timeout not passed yet error", async function () {
          await mineBlocks(dkgTimeout - 1)

          await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
            `timeout not passed yet`
          )
        })
      })

      context("with group creation timed out", async function () {
        let tx: ContractTransaction

        beforeEach("notify dkg timeout", async () => {
          await mineBlocks(dkgTimeout)
          tx = await randomBeacon.notifyDkgTimeout()
        })

        it("emits an event", async function () {
          await expect(tx)
            .to.emit(randomBeacon, "DkgTimedOut")
            .withArgs(expectedSeed)
        })

        it("cleans up dkg data", async function () {
          await assertDkgCleanData(randomBeacon)
        })
      })
    })

    context("when genesis dkg result was submitted", async function () {
      let resultIndex: BigNumber
      let dkgResult: DkgResult

      beforeEach(async () => {
        await randomBeacon.genesis()

        await mineBlocks(constants.timeDKG)
        ;({ resultIndex, dkgResult } = await signAndSubmitDkgResult(signers))
      })

      context("when challenge period not passed", async function () {
        context("when genesis dkg result was not approved", async function () {
          it("reverts with timeout not passed yet error", async function () {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              `timeout not passed yet`
            )
          })
        })
      })

      context("when challenge period passed", async function () {
        beforeEach(async () => {
          await mineBlocks(
            constants.groupSize * params.dkgResultSubmissionEligibilityDelay +
              params.dkgResultChallengePeriodLength
          )
        })

        context("when genesis dkg result was not approved", async function () {
          it("suceeds", async function () {
            await randomBeacon.notifyDkgTimeout()
          })
        })

        context("when genesis dkg result was approved", async function () {
          beforeEach(async () => {
            await randomBeacon.approveDkgResult(resultIndex, dkgResult)
          })

          it("reverts with dkg is currently not in progress error", async function () {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "dkg is currently not in progress"
            )
          })
        })
      })
    })

    context("with group creation not in progress", async function () {
      it("reverts with dkg is currently not in progress error", async function () {
        await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
          "dkg is currently not in progress"
        )
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
        let resultIndex: BigNumber
        let dkgResult: DkgResult

        beforeEach(async () => {
          await mineBlocks(constants.timeDKG)
          ;({
            transaction: tx,
            resultIndex,
            dkgResult,
          } = await signAndSubmitDkgResult(signers))
        })

        it("emits DkgStarted event", async function () {
          await expect(tx)
            .to.emit(randomBeacon, "DkgResultSubmitted")
            .withArgs(
              expectedSeed,
              resultIndex,
              dkgResult.submitterMemberIndex,
              dkgResult.groupPubKey,
              dkgResult.misbehaved,
              dkgResult.signatures,
              dkgResult.signingMembersIndexes,
              dkgResult.members
            )
        })

        it("registers a pending group", async function () {
          const group = await randomBeacon.getGroup(dkgResult.groupPubKey)

          expect(group.groupPubKey).to.be.equal(dkgResult.groupPubKey)
          expect(group.activationTimestamp).to.be.equal(0)
        })

        // TODO: Add test for DKG result validation

        // TODO: Move to result approval
        // it("cleans up dkg data", async function () {
        //   await signAndSubmitDkgResult(signers)

        //   // TODO: This test will be enhanced to clean up only after all DKG
        //   // results were submitted.
        //   await assertDkgCleanData(randomBeacon)
        // })
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

        context("with timeout notified", async function () {
          beforeEach("notify dkg timeout", async () => {
            await randomBeacon.notifyDkgTimeout()
          })

          it("reverts with dkg is currently not in progress error", async function () {
            await expect(signAndSubmitDkgResult(signers)).to.be.revertedWith(
              "dkg is currently not in progress"
            )
          })
        })
      })
    })
  })

  describe("approveDkgResult function call", async function () {
    context("with dkg result submitted", async function () {
      let resultIndex: BigNumber
      let dkgResult: DkgResult

      beforeEach(async () => {
        await randomBeacon.genesis()

        await mineBlocks(dkgTimeout)
        ;({ resultIndex, dkgResult } = await signAndSubmitDkgResult(signers))
      })

      context("with challenge period not passed", async function () {
        beforeEach(async () => {
          await mineBlocks(params.dkgResultChallengePeriodLength - 2)
        })

        it("reverts with challenge period not passed error", async function () {
          await expect(
            randomBeacon.approveDkgResult(resultIndex, dkgResult)
          ).revertedWith("Challenge period has not passed yet")
        })
      })

      context("with challenge period passed", async function () {
        beforeEach(async () => {
          await mineBlocks(params.dkgResultChallengePeriodLength)

          await randomBeacon.approveDkgResult(resultIndex, dkgResult)
        })

        it("cleans up dkg data", async function () {
          await assertDkgCleanData(randomBeacon)
        })
      })
    })
  })

  async function genesis(): Promise<[ContractTransaction, string]> {
    const tx = await randomBeacon.genesis()

    const expectedSeed = ethers.utils.keccak256(
      ethers.utils.solidityPack(
        ["uint256", "uint256"],
        [await randomBeacon.GENESIS_SEED(), tx.blockNumber]
      )
    )

    return [tx, expectedSeed]
  }

  interface DkgResult {
    submitterMemberIndex: number
    groupPubKey: string
    misbehaved: string
    signatures: string
    signingMembersIndexes: number[]
    members: string[]
  }

  async function signAndSubmitDkgResult(
    signers: DkgGroupSigners
  ): Promise<{
    transaction: ContractTransaction
    resultIndex: BigNumber
    dkgResult: DkgResult
  }> {
    const noMisbehaved = "0x"

    expect(signers.size, "unexpected signers map size").to.be.equal(
      constants.groupSize
    )

    const {
      members,
      signingMembersIndexes,
      signaturesBytes,
    } = await signDkgResult(signers, groupPublicKey, noMisbehaved)

    const submitterIndex = 1

    const dkgResult = {
      submitterMemberIndex: submitterIndex,
      groupPubKey: blsData.groupPubKey,
      misbehaved: noMisbehaved,
      signatures: signaturesBytes,
      signingMembersIndexes: signingMembersIndexes,
      members: members,
    }

    const transaction = await randomBeacon
      .connect(await ethers.getSigner(signers.get(submitterIndex)))
      .submitDkgResult(dkgResult)
    let eventFilter = randomBeacon.filters.DkgResultSubmitted()
    let events = await randomBeacon.queryFilter(eventFilter)
    expect(events).to.be.lengthOf(1)

    const resultIndex = events[0].args.index

    return { transaction, resultIndex, dkgResult }
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
  const signingMembersIndexes: number[] = []
  const signatures: string[] = []

  for (let [memberIndex, signer] of signers) {
    members.push(signer)

    signingMembersIndexes.push(memberIndex)

    const ethersSigner = await ethers.getSigner(signer)

    const signature = await ethersSigner.signMessage(
      ethers.utils.arrayify(resultHash)
    )

    signatures.push(signature)
  }

  expect(
    signingMembersIndexes.length,
    "unexpected signingMembersIndexes array size"
  ).to.be.equal(signers.size)

  expect(signatures.length, "unexpected signatures array size").to.be.equal(
    signers.size
  )

  const signaturesBytes: string = ethers.utils.hexConcat(signatures)

  return { members, signingMembersIndexes, signaturesBytes }
}

async function assertDkgCleanData(randomBeacon: RandomBeacon) {
  const dkgData = await randomBeacon.dkg()

  expect(dkgData.seed, "unexpected seed").to.eq(0)
  expect(dkgData.groupSize, "unexpected groupSize").to.eq(0)
  expect(dkgData.signatureThreshold, "unexpected signatureThreshold").to.eq(0)
  expect(dkgData.timeDKG, "unexpected timeDKG").to.eq(0)
  expect(
    dkgData.dkgResultSubmissionEligibilityDelay,
    "unexpected dkgResultSubmissionEligibilityDelay"
  ).to.eq(0)
  expect(
    dkgData.dkgResultChallengePeriodLength,
    "unexpected dkgResultChallengePeriodLength"
  ).to.eq(0)
  expect(dkgData.startBlock, "unexpected startBlock").to.eq(0)
}
