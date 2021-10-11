import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { blsData } from "./helpers/data"
import { groupStateEnum } from "./helpers/enums"
import { constants, testDeployment } from "./helpers/fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { ContractTransaction } from "ethers"
import type { RandomBeacon } from "../typechain"

const { mineBlocks } = helpers.time

describe("RandomBeacon contract", function () {
  const dkgTimeout: number =
    constants.timeDKG +
    constants.groupSize * constants.dkgResultSubmissionEligibilityDelay +
    constants.dkgResultChallengePeriodLength

  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let signer1: SignerWithAddress
  let signer2: SignerWithAddress
  let signer3: SignerWithAddress

  let randomBeacon: RandomBeacon

  before(async function () {
    ;[signer1, signer2, signer3] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async function () {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
  })

  describe("genesis function call", async function () {
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

    context("with genesis started", async function () {
      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
      })

      context("with genesis in progress", async function () {
        it("reverts with dkg is currently in progress error", async function () {
          // TODO: It should return a dedicated error from RandomBeacon contract.
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "dkg is currently in progress"
          )
        })

        // TODO: add more tests when more scenarios are covered.
      })

      context("with genesis completed", async function () {
        beforeEach(async () => {
          await mineBlocks(constants.timeDKG)

          const [, approveDkgResultFunc] = await signAndSubmitDkgResult()

          await mineBlocks(
            (await randomBeacon.dkg()).dkgResultChallengePeriodLength.toNumber()
          )

          await approveDkgResultFunc()
        })

        it("reverts with not awaiting genesis error", async function () {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "not awaiting genesis"
          )
        })

        // TODO: Add tests
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
        let approveDkgResultFunc: Function

        beforeEach(async () => {
          await mineBlocks(constants.timeDKG)
          ;[, approveDkgResultFunc] = await signAndSubmitDkgResult()
        })

        context("when genesis dkg result was not approved", async function () {
          it("returns true", async function () {
            expect(await randomBeacon.isDkgInProgress()).to.be.true
          })
        })

        context("when genesis dkg result was approved", async function () {
          beforeEach(async () => {
            await mineBlocks(constants.dkgResultChallengePeriodLength)

            await approveDkgResultFunc()
          })

          it("returns false", async function () {
            expect(await randomBeacon.isDkgInProgress()).to.be.false
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
      let approveDkgResultFunc: Function

      beforeEach(async () => {
        await randomBeacon.genesis()
        await mineBlocks(constants.timeDKG)
        ;[, approveDkgResultFunc] = await signAndSubmitDkgResult()
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
            constants.groupSize *
              constants.dkgResultSubmissionEligibilityDelay +
              constants.dkgResultChallengePeriodLength
          )
        })

        context("when genesis dkg result was not approved", async function () {
          it("suceeds", async function () {
            await randomBeacon.notifyDkgTimeout()
          })
        })

        context("when genesis dkg result was approved", async function () {
          beforeEach(async () => {
            await approveDkgResultFunc()
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
        await expect(signAndSubmitDkgResult()).to.be.revertedWith(
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
        beforeEach("increase time", async () => {
          await mineBlocks(constants.timeDKG)
        })

        it("succeeds", async function () {
          await signAndSubmitDkgResult()
        })

        // TODO: Move to result approval
        // it("cleans up dkg data", async function () {
        //   await signAndSubmitDkgResult()

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
          beforeEach(async () => {
            await signAndSubmitDkgResult()
          })

          it("sets group state", async function () {
            const groupData = await randomBeacon.callStatic.getGroup(
              groupPublicKey
            )
          })
        })

        context("with timeout notified", async function () {
          beforeEach("notify dkg timeout", async () => {
            await randomBeacon.notifyDkgTimeout()
          })

          it("reverts with dkg is currently not in progress error", async function () {
            await expect(signAndSubmitDkgResult()).to.be.revertedWith(
              "dkg is currently not in progress"
            )
          })
        })
      })
    })
  })

  describe("approveDkgResult function call", async function () {
    context("with dkg result submitted", async function () {
      let approveDkgResultFunc: Function

      beforeEach(async () => {
        await randomBeacon.genesis()
        await mineBlocks(dkgTimeout)
        ;[, approveDkgResultFunc] = await signAndSubmitDkgResult()
      })

      context("with challenge period not passed", async function () {
        beforeEach(async () => {
          await mineBlocks(constants.dkgResultChallengePeriodLength - 2)
        })

        it("reverts with challenge period not passed error", async function () {
          await expect(approveDkgResultFunc()).revertedWith(
            "Challenge period has not passed yet"
          )
        })
      })

      context("with challenge period passed", async function () {
        beforeEach(async () => {
          await mineBlocks(constants.dkgResultChallengePeriodLength)

          await approveDkgResultFunc()
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

  async function signAndSubmitDkgResult(): Promise<
    [ContractTransaction, Function]
  > {
    const noMisbehaved = "0x"

    const signers = new Map<number, SignerWithAddress>([
      [1, signer1],
      [2, signer2],
      [3, signer3],
    ])

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

    const transaction = await randomBeacon.submitDkgResult(dkgResult)
    let eventFilter = randomBeacon.filters.DkgResultSubmitted()
    let events = await randomBeacon.queryFilter(eventFilter)
    expect(events).to.be.lengthOf(1)

    const resultIndex = events[0].args.index

    const approveDkgResultFunc = async () => {
      return randomBeacon.approveDkgResult(resultIndex, dkgResult)
    }

    return [transaction, approveDkgResultFunc]
  }
})

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

async function signDkgResult(
  signers: Map<number, SignerWithAddress>,
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
    members.push(await signer.getAddress())

    signingMembersIndexes.push(memberIndex)

    const signature = await signer.signMessage(
      ethers.utils.arrayify(resultHash)
    )

    signatures.push(signature)
  }

  const signaturesBytes: string = ethers.utils.hexConcat(signatures)

  return { members, signingMembersIndexes, signaturesBytes }
}
