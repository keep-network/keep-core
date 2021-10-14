import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import { blsData } from "./data/bls"
import { constants, params, testDeployment } from "./fixtures"

import type { BigNumber, ContractTransaction, Signer } from "ethers"
import type { RandomBeacon } from "../typechain"
import type { Address } from "hardhat-deploy/types"

const { mineBlocks, mineBlocksTo } = helpers.time

describe("RandomBeacon", function () {
  const dkgTimeout: number =
    constants.offchainDkgTime +
    constants.groupSize * params.dkgResultSubmissionEligibilityDelay

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
      let expectedSeed: BigNumber

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
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()
        startBlock = genesisTx.blockNumber
      })

      context("with dkg result not submitted", async function () {
        it("reverts with dkg is currently in progress error", async function () {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "dkg is currently in progress"
          )
        })
      })

      context("with dkg result submitted", async function () {
        // TODO: Add test cases to cover results that are approved, challenged or
        // pending.

        context("with dkg result not approved", async function () {
          beforeEach(async () => {
            await mineBlocks(constants.offchainDkgTime)

            await signAndSubmitDkgResult(signers, startBlock)
          })

          it("reverts with dkg is currently in progress error", async function () {
            await expect(randomBeacon.genesis()).to.be.revertedWith(
              "dkg is currently in progress"
            )
          })
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
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()
        startBlock = genesisTx.blockNumber
      })

      it("returns true", async function () {
        expect(await randomBeacon.isDkgInProgress()).to.be.true
      })

      context("when genesis dkg result was submitted", async function () {
        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)
          await signAndSubmitDkgResult(signers, startBlock)
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
        await expect(signAndSubmitDkgResult(signers, 1)).to.be.revertedWith(
          "dkg is currently not in progress"
        )
      })
    })

    context("with group creation in progress", async function () {
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()

        startBlock = genesisTx.blockNumber
      })

      context("with group creation not timed out", async function () {
        context("with off-chain dkg time not passed", async function () {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime - 2)
          })

          it("reverts with submitter not eligible error", async function () {
            await expect(
              signAndSubmitDkgResult(signers, startBlock)
            ).to.revertedWith("Submitter not eligible")
          })
        })

        context("with off-chain dkg time passed", async function () {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime - 1)
          })

          it("reverts with less than threshold signers", async function () {
            const filteredSigners = new Map(
              Array.from(signers).filter(([index]) => {
                return index < constants.signatureThreshold
              })
            )

            await expect(
              signAndSubmitDkgResult(filteredSigners, startBlock)
            ).to.revertedWith("Too few signatures")
          })

          it("succeeds with threshold signers", async function () {
            const filteredSigners = new Map(
              Array.from(signers).filter(([index]) => {
                return index <= constants.signatureThreshold
              })
            )

            const { transaction: tx, dkgResult } = await signAndSubmitDkgResult(
              filteredSigners,
              startBlock
            )

            await expect(tx)
              .to.emit(randomBeacon, "DkgResultSubmitted")
              .withArgs(dkgResult.groupPubKey, signers.get(1))
          })

          it("succeeds for the first submitter", async function () {
            const { transaction: tx, dkgResult } = await signAndSubmitDkgResult(
              signers,
              startBlock,
              1
            )
            await expect(tx)
              .to.emit(randomBeacon, "DkgResultSubmitted")
              .withArgs(dkgResult.groupPubKey, signers.get(1))
          })

          it("reverts for the second submitter", async function () {
            await expect(
              signAndSubmitDkgResult(signers, startBlock, 2)
            ).to.revertedWith("Submitter not eligible")
          })

          context(
            "with first submitter eligibility delay period almost ended",
            async function () {
              beforeEach(async () => {
                await mineBlocksTo(
                  startBlock +
                    constants.offchainDkgTime +
                    params.dkgResultSubmissionEligibilityDelay -
                    2
                )
              })

              it("succeeds for the first submitter", async function () {
                const {
                  transaction: tx,
                  dkgResult
                } = await signAndSubmitDkgResult(signers, startBlock, 1)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(1))
              })

              it("reverts for the second submitter", async function () {
                await expect(
                  signAndSubmitDkgResult(signers, startBlock, 2)
                ).to.revertedWith("Submitter not eligible")
              })
            }
          )

          context(
            "with first submitter eligibility delay period ended",
            async function () {
              beforeEach(async () => {
                await mineBlocksTo(
                  startBlock +
                    constants.offchainDkgTime +
                    params.dkgResultSubmissionEligibilityDelay -
                    1
                )
              })

              it("succeeds for the first submitter", async function () {
                const {
                  transaction: tx,
                  dkgResult
                } = await signAndSubmitDkgResult(signers, startBlock, 1)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(1))
              })

              it("succeeds for the second submitter", async function () {
                const {
                  transaction: tx,
                  dkgResult
                } = await signAndSubmitDkgResult(signers, startBlock, 2)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(2))
              })

              it("reverts for the third submitter", async function () {
                await expect(
                  signAndSubmitDkgResult(signers, startBlock, 3)
                ).to.revertedWith("Submitter not eligible")
              })
            }
          )

          context(
            "with the last submitter eligibility delay period almost ended",
            async function () {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout - 1)
              })

              it(`succeeds for the first submitter`, async function () {
                const {
                  transaction: tx,
                  dkgResult
                } = await signAndSubmitDkgResult(signers, startBlock, 1)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(1))
              })

              it(`succeeds for the last submitter`, async function () {
                const {
                  transaction: tx,
                  dkgResult
                } = await signAndSubmitDkgResult(
                  signers,
                  startBlock,
                  constants.groupSize
                )

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(
                    dkgResult.groupPubKey,
                    signers.get(constants.groupSize)
                  )
              })
            }
          )
        })
      })

      context("with group creation timed out", async function () {
        beforeEach("increase time", async () => {
          await mineBlocksTo(startBlock + dkgTimeout)
        })

        context("with timeout not notified", async function () {
          it("reverts with dkg timeout already passed error", async function () {
            await expect(
              signAndSubmitDkgResult(signers, startBlock)
            ).to.revertedWith("dkg timeout already passed")
          })
        })
      })
    })
  })

  async function genesis(): Promise<[ContractTransaction, BigNumber]> {
    const tx = await randomBeacon.genesis()

    const expectedSeed = ethers.BigNumber.from(
      ethers.utils.keccak256(
        ethers.utils.solidityPack(
          ["uint256", "uint256"],
          [await randomBeacon.genesisSeed(), tx.blockNumber]
        )
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
    signers: DkgGroupSigners,
    startBlock: number,
    submitterIndex: number = 1
  ): Promise<{
    transaction: ContractTransaction

    dkgResult: DkgResult
  }> {
    const noMisbehaved = "0x"

    const {
      members,
      signingMemberIndices,
      signaturesBytes
    } = await signDkgResult(signers, groupPublicKey, noMisbehaved, startBlock)

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
  misbehaved: string,
  startBlock: number
) {
  const resultHash = ethers.utils.solidityKeccak256(
    ["bytes", "bytes", "uint256"],
    [groupPublicKey, misbehaved, startBlock]
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

  const signaturesBytes: string = ethers.utils.hexConcat(signatures)

  return { members, signingMemberIndices, signaturesBytes }
}
