/* eslint-disable no-await-in-loop */
import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import type { BigNumber, ContractTransaction, Signer } from "ethers"
import type { Address } from "hardhat-deploy/types"
import blsData from "./data/bls"
import { constants, params, testDeployment } from "./fixtures"

import type { RandomBeacon } from "../typechain"

const { mineBlocks, mineBlocksTo } = helpers.time

const dkgState = {
  IDLE: 0,
  KEY_GENERATION: 1,
  AWAITING_RESULT: 2,
  CHALLENGE: 3,
}

describe("RandomBeacon", () => {
  const dkgTimeout: number =
    constants.offchainDkgTime +
    constants.groupSize * params.dkgResultSubmissionEligibilityDelay

  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let thirdParty: Signer
  let signers: DkgGroupSigners

  let randomBeacon: RandomBeacon

  before(async () => {
    thirdParty = await ethers.getSigner((await getUnnamedAccounts())[1])

    // Accounts offset provided to getDkgGroupSigners have to include number of
    // unnamed accounts that were already used.
    signers = await getDkgGroupSigners(constants.groupSize, 1)
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
  })

  describe("getDkgParameters", async () => {
    it("returns values", async () => {
      const result = await randomBeacon.getDkgParameters()

      expect(result[0]).to.be.equal(params.dkgResultChallengePeriodLength)
      expect(result[1]).to.be.equal(params.dkgResultSubmissionEligibilityDelay)
    })
  })

  describe("genesis", async () => {
    it("can be invoked by third party", async () => {
      await randomBeacon.connect(thirdParty).genesis()
    })

    context("with initial contract state", async () => {
      let tx: ContractTransaction
      let expectedSeed: BigNumber

      beforeEach("run genesis", async () => {
        // eslint-disable-next-line @typescript-eslint/no-extra-semi
        ;[tx, expectedSeed] = await genesis()
      })

      it("emits DkgStarted event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "DkgStarted")
          .withArgs(expectedSeed)
      })
    })

    context("with genesis in progress", async () => {
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()
        startBlock = genesisTx.blockNumber
      })

      context("with dkg result not submitted", async () => {
        it("reverts with 'current state is not IDLE' error", async () => {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "current state is not IDLE"
          )
        })
      })

      context("with dkg result submitted", async () => {
        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)
          await signAndSubmitDkgResult(signers, startBlock)
        })

        // TODO: Add test cases to cover results that are approved, challenged or
        // pending.

        context("with dkg result not approved", async () => {
          it("reverts with 'current state is not IDLE' error", async () => {
            await expect(randomBeacon.genesis()).to.be.revertedWith(
              "current state is not IDLE"
            )
          })
        })
      })
    })
  })

  describe("getGroupCreationState", async () => {
    context("with initial contract state", async () => {
      it("returns IDLE state", async () => {
        expect(await randomBeacon.getGroupCreationState()).to.be.equal(
          dkgState.IDLE
        )
      })
    })

    context("when genesis dkg started", async () => {
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()
        startBlock = genesisTx.blockNumber
      })

      context("at the start of off-chain dkg period", async () => {
        it("returns KEY_GENERATION state", async () => {
          expect(await randomBeacon.getGroupCreationState()).to.be.equal(
            dkgState.KEY_GENERATION
          )
        })
      })

      context("at the end of off-chain dkg period", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        it("returns KEY_GENERATION state", async () => {
          expect(await randomBeacon.getGroupCreationState()).to.be.equal(
            dkgState.KEY_GENERATION
          )
        })
      })

      context("after off-chain dkg period", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime + 1)
        })

        context("when dkg result was not submitted", async () => {
          it("returns AWAITING_RESULT state", async () => {
            expect(await randomBeacon.getGroupCreationState()).to.be.equal(
              dkgState.AWAITING_RESULT
            )
          })

          context("after the dkg timeout period", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + dkgTimeout + 1)
            })

            it("returns AWAITING_RESULT state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.AWAITING_RESULT
              )
            })
          })
        })

        context("when dkg result was submitted", async () => {
          beforeEach(async () => {
            await signAndSubmitDkgResult(signers, startBlock)
          })

          context("when dkg result was not approved", async () => {
            it("returns CHALLENGE state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.CHALLENGE
              )
            })
          })

          // TODO: Enable once approvals and challenges are implemented
          // context("when dkg result was approved", async function () {
          //   it("returns IDLE state", async function () {
          //     expect(await randomBeacon.getGroupCreationState()).to.be.equal(
          //       dkgState.IDLE
          //     )
          //   })
          // })

          // context("when dkg result was challenged", async function () {
          //   it("returns AWAITING_RESULT state", async function () {
          //     expect(await randomBeacon.getGroupCreationState()).to.be.equal(
          //       dkgState.AWAITING_RESULT
          //     )
          //   })
          // })
        })
      })
    })
  })

  describe("hasDkgTimedOut", async () => {
    context("with initial contract state", async () => {
      it("returns false", async () => {
        await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
      })
    })

    context("when genesis dkg started", async () => {
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()
        startBlock = genesisTx.blockNumber
      })

      context("within off-chain dkg period", async () => {
        it("returns false", async () => {
          await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
        })
      })

      context("after off-chain dkg period", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime + 1)
        })

        context("when dkg result was not submitted", async () => {
          it("returns false", async () => {
            await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
          })

          context("at the end of the dkg timeout period", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + dkgTimeout)
            })

            it("returns false", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
            })
          })

          context("after the dkg timeout period", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + dkgTimeout + 1)
            })

            it("returns true", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.true
            })
          })
        })

        context("when dkg result was submitted", async () => {
          beforeEach(async () => {
            await signAndSubmitDkgResult(signers, startBlock)
          })

          context("when dkg result was not approved", async () => {
            context("at the end of the dkg timeout period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout)
              })

              it("returns false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after the dkg timeout period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout + 1)
              })

              it("returns false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })
          })

          // TODO: Enable once approvals and challenges are implemented
          // context("when dkg result was approved", async function () {
          // it("returns false", async () => {
          //   await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
          // })
          // })
          // TODO: Add test cases to cover transition after challenge back to
          // the awaiting result state and counting timeout there.
          // context("when dkg result was challenged", async function () {
          // it("returns false", async () => {
          //   await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
          // })
          // })
        })
      })
    })
  })

  describe("submitDkgResult", async () => {
    // TODO: Add more tests to cover the DKG result verification function thoroughly.

    context("with initial contract state", async () => {
      it("reverts with 'current state is not AWAITING_RESULT' error", async () => {
        await expect(signAndSubmitDkgResult(signers, 1)).to.be.revertedWith(
          "current state is not AWAITING_RESULT"
        )
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis()

        startBlock = genesisTx.blockNumber
      })

      context("with group creation not timed out", async () => {
        context("with off-chain dkg time not passed", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime - 1)
          })

          it("reverts with 'current state is not AWAITING_RESULT' error", async () => {
            await expect(
              signAndSubmitDkgResult(signers, startBlock)
            ).to.revertedWith("current state is not AWAITING_RESULT")
          })
        })

        context("with off-chain dkg time passed", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime)
          })

          it("reverts with less than threshold signers", async () => {
            const filteredSigners = new Map(
              Array.from(signers).filter(
                ([index]) => index < constants.signatureThreshold
              )
            )

            await expect(
              signAndSubmitDkgResult(filteredSigners, startBlock)
            ).to.revertedWith("Too few signatures")
          })

          it("succeeds with threshold signers", async () => {
            const filteredSigners = new Map(
              Array.from(signers).filter(
                ([index]) => index <= constants.signatureThreshold
              )
            )

            const { transaction: tx, dkgResult } = await signAndSubmitDkgResult(
              filteredSigners,
              startBlock
            )

            await expect(tx)
              .to.emit(randomBeacon, "DkgResultSubmitted")
              .withArgs(dkgResult.groupPubKey, signers.get(1))
          })

          it("succeeds for the first submitter", async () => {
            const { transaction: tx, dkgResult } = await signAndSubmitDkgResult(
              signers,
              startBlock,
              1
            )
            await expect(tx)
              .to.emit(randomBeacon, "DkgResultSubmitted")
              .withArgs(dkgResult.groupPubKey, signers.get(1))
          })

          it("reverts for the second submitter", async () => {
            await expect(
              signAndSubmitDkgResult(signers, startBlock, 2)
            ).to.revertedWith("Submitter not eligible")
          })

          context(
            "with first submitter eligibility delay period almost ended",
            async () => {
              beforeEach(async () => {
                await mineBlocksTo(
                  startBlock +
                    constants.offchainDkgTime +
                    params.dkgResultSubmissionEligibilityDelay -
                    2
                )
              })

              it("succeeds for the first submitter", async () => {
                const { transaction: tx, dkgResult } =
                  await signAndSubmitDkgResult(signers, startBlock, 1)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(1))
              })

              it("reverts for the second submitter", async () => {
                await expect(
                  signAndSubmitDkgResult(signers, startBlock, 2)
                ).to.revertedWith("Submitter not eligible")
              })
            }
          )

          context(
            "with first submitter eligibility delay period ended",
            async () => {
              beforeEach(async () => {
                await mineBlocksTo(
                  startBlock +
                    constants.offchainDkgTime +
                    params.dkgResultSubmissionEligibilityDelay -
                    1
                )
              })

              it("succeeds for the first submitter", async () => {
                const { transaction: tx, dkgResult } =
                  await signAndSubmitDkgResult(signers, startBlock, 1)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(1))
              })

              it("succeeds for the second submitter", async () => {
                const { transaction: tx, dkgResult } =
                  await signAndSubmitDkgResult(signers, startBlock, 2)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(2))
              })

              it("reverts for the third submitter", async () => {
                await expect(
                  signAndSubmitDkgResult(signers, startBlock, 3)
                ).to.revertedWith("Submitter not eligible")
              })
            }
          )

          context(
            "with the last submitter eligibility delay period almost ended",
            async () => {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout - 1)
              })

              it("succeeds for the first submitter", async () => {
                const { transaction: tx, dkgResult } =
                  await signAndSubmitDkgResult(signers, startBlock, 1)

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(dkgResult.groupPubKey, signers.get(1))
              })

              it("succeeds for the last submitter", async () => {
                const { transaction: tx, dkgResult } =
                  await signAndSubmitDkgResult(
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

      context("with group creation timed out", async () => {
        beforeEach("increase time", async () => {
          await mineBlocksTo(startBlock + dkgTimeout)
        })

        context("with timeout not notified", async () => {
          it("reverts with dkg timeout already passed error", async () => {
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
    // eslint-disable-next-line @typescript-eslint/no-shadow
    signers: DkgGroupSigners,
    startBlock: number,
    submitterIndex = 1
  ): Promise<{
    transaction: ContractTransaction

    dkgResult: DkgResult
  }> {
    const noMisbehaved = "0x"

    const { members, signingMemberIndices, signaturesBytes } =
      await signDkgResult(signers, groupPublicKey, noMisbehaved, startBlock)

    const dkgResult = {
      submitterMemberIndex: submitterIndex,
      groupPubKey: blsData.groupPubKey,
      misbehaved: noMisbehaved,
      signatures: signaturesBytes,
      signingMemberIndices,
      members,
    }

    const transaction = await randomBeacon
      .connect(await ethers.getSigner(signers.get(submitterIndex)))
      .submitDkgResult(dkgResult)

    return { transaction, dkgResult }
  }
})

type DkgGroupSigners = Map<number, Address>

async function getDkgGroupSigners(
  groupSize: number,
  startAccountsOffset: number
): Promise<DkgGroupSigners> {
  const signers = new Map<number, Address>()

  for (let i = 1; i <= groupSize; i++) {
    const signer = (await getUnnamedAccounts())[startAccountsOffset + i]

    await expect(
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

  // eslint-disable-next-line no-restricted-syntax
  for (const [memberIndex, signer] of signers) {
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
