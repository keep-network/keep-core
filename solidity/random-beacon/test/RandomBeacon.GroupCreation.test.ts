/* eslint-disable @typescript-eslint/no-unused-expressions, @typescript-eslint/no-extra-semi */

import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"

import blsData from "./data/bls"
import { constants, dkgState, params, testDeployment } from "./fixtures"
import {
  genesis,
  signAndSubmitCorrectDkgResult,
  signAndSubmitArbitraryDkgResult,
  noMisbehaved,
  signAndSubmitUnrecoverableDkgResult,
  hashDKGMembers,
  expectDkgResultSubmittedEvent,
} from "./utils/dkg"
import { registerOperators } from "./utils/operators"
import { selectGroup, hashUint32Array } from "./utils/groups"
import { fakeTokenStaking } from "./mocks/staking"

import type { Operator } from "./utils/operators"
import type { BeaconDkg as DKG } from "../typechain/RandomBeaconStub"
import type { FakeContract } from "@defi-wonderland/smock"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  RandomBeacon,
  RandomBeaconGovernance,
  TestToken,
  SortitionPool,
  StakingStub,
  IRandomBeaconStaking,
} from "../typechain"
import type { BigNumber, ContractTransaction, Signer } from "ethers"

const { mineBlocks, mineBlocksTo } = helpers.time
const { to1e18 } = helpers.number
const { keccak256 } = ethers.utils
const { createSnapshot, restoreSnapshot } = helpers.snapshot

// FIXME: As a workaround for a bug https://github.com/dethcrypto/TypeChain/issues/601
// we declare a new type instead of using `RandomBeaconStub & RandomBeacon` intersection.
type RandomBeaconTest = RandomBeacon & {
  getDkgData: () => Promise<DKG.DataStructOutput>
}

const fixture = async () => {
  const contracts = await testDeployment()

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const signers = await registerOperators(
    contracts.randomBeacon as RandomBeacon,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  const randomBeaconGovernance =
    contracts.randomBeaconGovernance as RandomBeaconGovernance
  const randomBeacon = contracts.randomBeacon as RandomBeaconTest
  const sortitionPool = contracts.sortitionPool as SortitionPool
  const staking = contracts.stakingStub as StakingStub
  const testToken = contracts.testToken as TestToken

  return {
    randomBeaconGovernance,
    randomBeacon,
    sortitionPool,
    staking,
    testToken,
    signers,
  }
}

// Test suite covering group creation in RandomBeacon contract.
// It covers DKG and Groups libraries usage in the process of group creation.
describe("RandomBeacon - Group Creation", () => {
  const dkgTimeout: number =
    constants.offchainDkgTime + params.dkgResultSubmissionTimeout
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let thirdParty: Signer
  let signers: Operator[]

  let randomBeaconGovernance: RandomBeaconGovernance
  let randomBeacon: RandomBeaconTest
  let sortitionPool: SortitionPool
  let staking: StakingStub
  let testToken: TestToken

  before(async () => {
    thirdParty = await ethers.getSigner((await getUnnamedAccounts())[0])
    let randomBeaconStub: RandomBeaconTest
    ;({
      randomBeaconGovernance,
      randomBeacon: randomBeaconStub,
      sortitionPool,
      staking,
      testToken,
      signers,
    } = await waffle.loadFixture(fixture))

    randomBeacon = randomBeaconStub

    // Fund DKG rewards pool to make testing of rewards possible.
    await fundDkgRewardsPool(to1e18(100))
  })

  describe("genesis", async () => {
    context("when called by a third party", async () => {
      let tx: Promise<ContractTransaction>

      before("run genesis", async () => {
        await createSnapshot()

        tx = randomBeacon.connect(thirdParty).genesis()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should succeed", async () => {
        await expect(tx).to.not.be.reverted
      })
    })

    context("with initial contract state", async () => {
      let tx: ContractTransaction
      let expectedSeed: BigNumber

      before("run genesis", async () => {
        await createSnapshot()
        ;[tx, expectedSeed] = await genesis(randomBeacon)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should lock the sortition pool", async () => {
        expect(await sortitionPool.isLocked()).to.be.true
      })

      it("should emit DkgStateLocked event", async () => {
        await expect(tx).to.emit(randomBeacon, "DkgStateLocked")
      })

      it("should emit DkgStarted event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "DkgStarted")
          .withArgs(expectedSeed)
      })
    })

    context("with genesis in progress", async () => {
      let startBlock: number
      let genesisSeed: BigNumber

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx, seed] = await genesis(randomBeacon)
        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with dkg result not submitted", async () => {
        it("should revert with 'Current state is not IDLE' error", async () => {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "Current state is not IDLE"
          )
        })
      })

      context("with valid dkg result submitted", async () => {
        let dkgResult: DKG.ResultStruct
        let submitter: SignerWithAddress

        before(async () => {
          await createSnapshot()

          await mineBlocks(constants.offchainDkgTime)
          ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with dkg result not approved", async () => {
          it("should revert with 'Current state is not IDLE' error", async () => {
            await expect(randomBeacon.genesis()).to.be.revertedWith(
              "Current state is not IDLE"
            )
          })
        })

        context("with dkg result approved", async () => {
          before("approve dkg result", async () => {
            await createSnapshot()

            await mineBlocks(params.dkgResultChallengePeriodLength)

            await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should succeed", async () => {
            await expect(randomBeacon.genesis()).to.be.revertedWith(
              "Not awaiting genesis"
            )
          })
        })
      })

      context("with invalid dkg result submitted", async () => {
        context("with dkg result challenged", async () => {
          before("submit and challenge dkg result", async () => {
            await createSnapshot()

            await mineBlocks(constants.offchainDkgTime)
            const { dkgResult } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix operators to make the result malicious
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            )

            await randomBeacon.challengeDkgResult(dkgResult)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert", async () => {
            await expect(randomBeacon.genesis()).to.be.revertedWith(
              "Current state is not IDLE"
            )
          })
        })
      })

      context("with dkg timeout notified", async () => {
        before("notify dkg timeout", async () => {
          await createSnapshot()

          await mineBlocks(dkgTimeout)

          await randomBeacon.notifyDkgTimeout()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should succeed", async () => {
          await expect(randomBeacon.genesis()).to.not.be.reverted
        })
      })
    })
  })

  describe("getGroupCreationState", async () => {
    context("with initial contract state", async () => {
      it("should return IDLE state", async () => {
        expect(await randomBeacon.getGroupCreationState()).to.be.equal(
          dkgState.IDLE
        )
      })
    })

    context("when genesis dkg started", async () => {
      let startBlock: number
      let genesisSeed: BigNumber

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx, seed] = await genesis(randomBeacon)
        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("at the start of off-chain dkg period", async () => {
        it("should return KEY_GENERATION state", async () => {
          expect(await randomBeacon.getGroupCreationState()).to.be.equal(
            dkgState.KEY_GENERATION
          )
        })
      })

      context("at the end of off-chain dkg period", async () => {
        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should return KEY_GENERATION state", async () => {
          expect(await randomBeacon.getGroupCreationState()).to.be.equal(
            dkgState.KEY_GENERATION
          )
        })
      })

      context("after off-chain dkg period", async () => {
        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + constants.offchainDkgTime + 1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when dkg result was not submitted", async () => {
          it("should return AWAITING_RESULT state", async () => {
            expect(await randomBeacon.getGroupCreationState()).to.be.equal(
              dkgState.AWAITING_RESULT
            )
          })

          context("after the dkg timeout period", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(startBlock + dkgTimeout + 1)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return AWAITING_RESULT state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.AWAITING_RESULT
              )
            })
          })
        })

        context("when dkg result was submitted", async () => {
          let dkgResult: DKG.ResultStruct
          let submitter: SignerWithAddress

          before("submit dkg result", async () => {
            await createSnapshot()
            ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            ))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when dkg result was not approved", async () => {
            it("should return CHALLENGE state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.CHALLENGE
              )
            })
          })

          context("when dkg result was approved", async () => {
            before("approve dkg result", async () => {
              await createSnapshot()

              await mineBlocks(params.dkgResultChallengePeriodLength)

              await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return IDLE state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.IDLE
              )
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DKG.ResultStruct

          before("submit malicious dkg result", async () => {
            await createSnapshot()
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            ))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when dkg result was challenged", async () => {
            before("challenge dkg result", async () => {
              await createSnapshot()

              await randomBeacon.challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return AWAITING_RESULT state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.AWAITING_RESULT
              )
            })
          })
        })
      })
    })
  })

  describe("hasDkgTimedOut", async () => {
    context("with initial contract state", async () => {
      it("should return false", async () => {
        await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
      })
    })

    context("when genesis dkg started", async () => {
      let startBlock: number
      let genesisSeed

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx, seed] = await genesis(randomBeacon)
        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("within off-chain dkg period", async () => {
        it("should return false", async () => {
          await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
        })
      })

      context("after off-chain dkg period", async () => {
        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + constants.offchainDkgTime + 1)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("when dkg result was not submitted", async () => {
          it("should return false", async () => {
            await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
          })

          context("at the end of the dkg timeout period", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(startBlock + dkgTimeout)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return false", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
            })
          })

          context("after the dkg timeout period", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(startBlock + dkgTimeout + 1)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return true", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.true
            })
          })
        })

        context("when dkg result was submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResult: DKG.ResultStruct
          let submitter: SignerWithAddress

          before(async () => {
            await createSnapshot()

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              submitter,
            } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            ))

            resultSubmissionBlock = tx.blockNumber
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when dkg result was not approved", async () => {
            context("at the end of the dkg timeout period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(startBlock + dkgTimeout)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after the dkg timeout period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(startBlock + dkgTimeout + 1)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("at the end of the challenge period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  resultSubmissionBlock + params.dkgResultChallengePeriodLength
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after the challenge period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  resultSubmissionBlock +
                    params.dkgResultChallengePeriodLength +
                    1
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })
          })

          context("when dkg result was approved", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )

              await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return false", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DKG.ResultStruct

          before(async () => {
            await createSnapshot()
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious.
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            ))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when dkg result was challenged", async () => {
            let challengeBlockNumber: number

            before(async () => {
              await createSnapshot()

              const tx = await randomBeacon.challengeDkgResult(dkgResult)
              challengeBlockNumber = tx.blockNumber
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("at the end of dkg result submission period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  challengeBlockNumber + params.dkgResultSubmissionTimeout
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after dkg result submission period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  challengeBlockNumber + params.dkgResultSubmissionTimeout + 1
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return true", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.true
              })
            })
          })
        })
      })
    })
  })

  describe("submitDkgResult", async () => {
    context("with initial contract state", async () => {
      it("should revert with 'Current state is not AWAITING_RESULT' error", async () => {
        await expect(
          signAndSubmitArbitraryDkgResult(
            randomBeacon,
            groupPublicKey,
            signers,
            1,
            noMisbehaved
          )
        ).to.be.revertedWith("Current state is not AWAITING_RESULT")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let genesisSeed: BigNumber

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx, seed] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with group creation not timed out", async () => {
        context("with off-chain dkg time not passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(startBlock + constants.offchainDkgTime - 1)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'Current state is not AWAITING_RESULT' error", async () => {
            await expect(
              signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )
            ).to.be.revertedWith("Current state is not AWAITING_RESULT")
          })
        })

        context("with off-chain dkg time passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(startBlock + constants.offchainDkgTime)
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with enough signatures on the result", async () => {
            let tx: ContractTransaction
            let dkgResult: DKG.ResultStruct
            let dkgResultHash: string

            before(async () => {
              await createSnapshot()
              ;({
                transaction: tx,
                dkgResult,
                dkgResultHash,
              } = await signAndSubmitArbitraryDkgResult(
                randomBeacon,
                groupPublicKey,
                signers,
                startBlock,
                noMisbehaved,
                1,
                undefined,
                constants.groupThreshold
              ))
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultSubmitted event", async () => {
              await expectDkgResultSubmittedEvent(tx, {
                resultHash: dkgResultHash,
                seed: genesisSeed,
                result: dkgResult,
              })
            })

            it("should not register a group", async () => {
              const groupsRegistry = await randomBeacon.getGroupsRegistry()

              expect(groupsRegistry).to.be.lengthOf(0)
            })

            it("should not emit GroupRegistered event", async () => {
              await expect(tx).not.to.emit(randomBeacon, "GroupRegistered")
            })

            it("should not unlock the sortition pool", async () => {
              expect(await sortitionPool.isLocked()).to.be.true
            })
          })

          context("with not enough signatures on the result", async () => {
            let tx: ContractTransaction
            let dkgResult: DKG.ResultStruct
            let dkgResultHash: string

            before(async () => {
              await createSnapshot()
              ;({
                transaction: tx,
                dkgResult,
                dkgResultHash,
              } = await signAndSubmitArbitraryDkgResult(
                randomBeacon,
                groupPublicKey,
                signers,
                startBlock,
                noMisbehaved,
                1,
                undefined,
                constants.groupThreshold - 1
              ))
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultSubmitted event", async () => {
              await expectDkgResultSubmittedEvent(tx, {
                resultHash: dkgResultHash,
                seed: genesisSeed,
                result: dkgResult,
              })
            })

            it("should not register a group", async () => {
              const groupsRegistry = await randomBeacon.getGroupsRegistry()

              expect(groupsRegistry).to.be.lengthOf(0)
            })

            it("should not emit GroupRegistered event", async () => {
              await expect(tx).not.to.emit(randomBeacon, "GroupRegistered")
            })

            it("should not unlock the sortition pool", async () => {
              expect(await sortitionPool.isLocked()).to.be.true
            })
          })

          describe("submission eligibility verification", async () => {
            let submissionStartBlockNumber: number

            before(async () => {
              await createSnapshot()

              submissionStartBlockNumber =
                startBlock + constants.offchainDkgTime
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("at the beginning of the submission period", async () => {
              it("should succeed for the first member", async () => {
                await assertSubmissionSucceeds(1)
              })

              it("should succeed for the second member", async () => {
                await assertSubmissionSucceeds(2)
              })

              it("should succeed for the last member", async () => {
                await assertSubmissionSucceeds(constants.groupSize - 1)
              })
            })

            context("at the end of the submission period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  submissionStartBlockNumber +
                    params.dkgResultSubmissionTimeout -
                    1
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should succeed for the first member", async () => {
                await assertSubmissionSucceeds(1)
              })

              it("should succeed for the second member", async () => {
                await assertSubmissionSucceeds(2)
              })

              it("should succeed for the last member", async () => {
                await assertSubmissionSucceeds(constants.groupSize - 1)
              })
            })

            context("after the submission period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  submissionStartBlockNumber + params.dkgResultSubmissionTimeout
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should revert for the first member", async () => {
                await assertSubmissionReverts(1)
              })

              it("should revert for the second member", async () => {
                await assertSubmissionReverts(2)
              })

              it("should revert for the last member", async () => {
                await assertSubmissionReverts(constants.groupSize - 1)
              })
            })
          })

          context("with dkg result submitted", async () => {
            let dkgResult: DKG.ResultStruct
            let submitter: SignerWithAddress
            let resultSubmissionBlock: number

            before(async () => {
              await createSnapshot()

              let tx: ContractTransaction
              ;({
                transaction: tx,
                dkgResult,
                submitter,
              } = await signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              ))

              resultSubmissionBlock = tx.blockNumber
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert 'Current state is not AWAITING_RESULT' error", async () => {
              await expect(
                signAndSubmitCorrectDkgResult(
                  randomBeacon,
                  groupPublicKey,
                  genesisSeed,
                  startBlock,
                  noMisbehaved
                )
              ).to.be.revertedWith("Current state is not AWAITING_RESULT")
            })

            context("with dkg result approved", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  resultSubmissionBlock + params.dkgResultChallengePeriodLength
                )

                await randomBeacon
                  .connect(submitter)
                  .approveDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should revert", async () => {
                await expect(
                  signAndSubmitCorrectDkgResult(
                    randomBeacon,
                    groupPublicKey,
                    genesisSeed,
                    startBlock,
                    noMisbehaved
                  )
                ).to.be.revertedWith("Sortition pool unlocked")
              })
            })
          })

          context("with dkg result challenged", async () => {
            let challengeBlockNumber: number

            before(async () => {
              await createSnapshot()
              await mineBlocksTo(startBlock + constants.offchainDkgTime)
              const { dkgResult } = await signAndSubmitArbitraryDkgResult(
                randomBeacon,
                groupPublicKey,
                // Mix signers to make the result malicious.
                mixSigners(await selectGroup(sortitionPool, genesisSeed)),
                startBlock,
                noMisbehaved
              )

              const tx = await randomBeacon.challengeDkgResult(dkgResult)
              challengeBlockNumber = tx.blockNumber
            })

            after(async () => {
              await restoreSnapshot()
            })

            describe("group registration", async () => {
              let tx: ContractTransaction

              before(async () => {
                await createSnapshot()
                ;({ transaction: tx } = await signAndSubmitCorrectDkgResult(
                  randomBeacon,
                  groupPublicKey,
                  genesisSeed,
                  startBlock,
                  noMisbehaved
                ))
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should not register a group", async () => {
                const groupsRegistry = await randomBeacon.getGroupsRegistry()

                expect(groupsRegistry).to.be.lengthOf(0)
              })

              it("should not emit GroupRegistered event", async () => {
                await expect(tx).not.to.emit(randomBeacon, "GroupRegistered")
              })
            })

            describe("submission eligibility verification", async () => {
              let submissionStartBlockNumber: number

              beforeEach(async () => {
                await createSnapshot()

                submissionStartBlockNumber = challengeBlockNumber
              })

              afterEach(async () => {
                await restoreSnapshot()
              })

              context("at the beginning of the submission period", async () => {
                it("should succeed for the first member", async () => {
                  await assertSubmissionSucceeds(1)
                })

                it("should succeed for the second member", async () => {
                  await assertSubmissionSucceeds(2)
                })

                it("should succeed for the last member", async () => {
                  await assertSubmissionSucceeds(constants.groupSize - 1)
                })
              })

              context("at the end of the submission period", async () => {
                before(async () => {
                  await createSnapshot()

                  await mineBlocksTo(
                    submissionStartBlockNumber +
                      params.dkgResultSubmissionTimeout -
                      1
                  )
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should succeed for the first member", async () => {
                  await assertSubmissionSucceeds(1)
                })

                it("should succeed for the second member", async () => {
                  await assertSubmissionSucceeds(2)
                })

                it("should succeed for the last member", async () => {
                  await assertSubmissionSucceeds(constants.groupSize - 1)
                })
              })

              context("after the submission period", async () => {
                before(async () => {
                  await createSnapshot()

                  await mineBlocksTo(
                    submissionStartBlockNumber +
                      params.dkgResultSubmissionTimeout
                  )
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should revert for the first member", async () => {
                  await assertSubmissionReverts(1)
                })

                it("should revert for the second member", async () => {
                  await assertSubmissionReverts(2)
                })

                it("should revert for the last member", async () => {
                  await assertSubmissionReverts(constants.groupSize - 1)
                })
              })
            })

            context("with misbehaved members", async () => {
              let tx: ContractTransaction
              let dkgResult: DKG.ResultStruct
              let dkgResultHash: string

              context(
                "when misbehaved members are in ascending order",
                async () => {
                  const misbehavedIndices = [2, 9, 11, 30, 60, 64]

                  before(async () => {
                    await createSnapshot()
                    ;({
                      transaction: tx,
                      dkgResult,
                      dkgResultHash,
                    } = await signAndSubmitCorrectDkgResult(
                      randomBeacon,
                      groupPublicKey,
                      genesisSeed,
                      startBlock,
                      misbehavedIndices
                    ))
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should emit DkgResultSubmitted", async () => {
                    await expectDkgResultSubmittedEvent(tx, {
                      resultHash: dkgResultHash,
                      seed: genesisSeed,
                      result: dkgResult,
                    })
                  })
                }
              )
            })
          })
        })
      })

      context("with group creation timed out", async () => {
        before("increase time", async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with timeout not notified", async () => {
          it("should revert with DKG timeout already passed error", async () => {
            await expect(
              signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )
            ).to.be.revertedWith("DKG timeout already passed")
          })
        })
      })

      // Submission Test Helpers
      async function assertSubmissionSucceeds(
        submitterIndex: number
      ): Promise<void> {
        await createSnapshot()

        const {
          transaction: tx,
          dkgResult,
          dkgResultHash,
        } = await signAndSubmitCorrectDkgResult(
          randomBeacon,
          groupPublicKey,
          genesisSeed,
          startBlock,
          noMisbehaved,
          submitterIndex
        )

        await expectDkgResultSubmittedEvent(tx, {
          resultHash: dkgResultHash,
          seed: genesisSeed,
          result: dkgResult,
        })

        await restoreSnapshot()
      }

      async function assertSubmissionReverts(
        submitterIndex: number,
        message = "DKG timeout already passed"
      ): Promise<void> {
        await expect(
          signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved,
            submitterIndex
          )
        ).to.be.revertedWith(message)
      }
    })
  })

  describe("approveDkgResult", async () => {
    // Just to make `approveDkgResult` call possible.
    const stubDkgResult: DKG.ResultStruct = {
      groupPubKey: blsData.groupPubKey,
      members: [1, 2, 3, 4],
      misbehavedMembersIndices: [],
      signatures: "0x01020304",
      signingMembersIndices: [1, 2, 3, 4],
      submitterMemberIndex: 1,
      membersHash: hashDKGMembers([1, 2, 3, 4], []),
    }

    context("with initial contract state", async () => {
      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          randomBeacon.approveDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let genesisSeed: BigNumber

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx, seed] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          randomBeacon.approveDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })

      context("with off-chain dkg time passed", async () => {
        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with dkg result not submitted", async () => {
          it("should revert with 'Current state is not CHALLENGE' error", async () => {
            await expect(
              randomBeacon.approveDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })
        })

        context("with dkg result submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DKG.ResultStruct
          let submitter: SignerWithAddress

          const submitterIndex = 1

          before(async () => {
            await createSnapshot()

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
              submitter,
            } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved,
              submitterIndex
            ))

            resultSubmissionBlock = tx.blockNumber
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with challenge period not passed", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock +
                  params.dkgResultChallengePeriodLength -
                  1
              )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'Challenge period has not passed yet' error", async () => {
              await expect(
                randomBeacon.connect(submitter).approveDkgResult(dkgResult)
              ).to.be.revertedWith("Challenge period has not passed yet")
            })
          })

          context("with challenge period passed", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("when called by a DKG result submitter", async () => {
              let tx: ContractTransaction
              let initialDkgRewardsPoolBalance: BigNumber
              let initialSubmitterBalance: BigNumber

              before(async () => {
                await createSnapshot()

                initialDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                initialSubmitterBalance = await testToken.balanceOf(
                  await submitter.getAddress()
                )
                tx = await randomBeacon
                  .connect(submitter)
                  .approveDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultApproved event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultApproved")
                  .withArgs(dkgResultHash, await submitter.getAddress())
              })

              it("should clean dkg data", async () => {
                await assertDkgResultCleanData(randomBeacon)
              })

              it("should register a group", async () => {
                const groupsRegistry = await randomBeacon.getGroupsRegistry()

                expect(groupsRegistry).to.be.lengthOf(1)
                expect(groupsRegistry[0]).to.deep.equal(
                  keccak256(groupPublicKey)
                )

                const storedGroup = await randomBeacon["getGroup(bytes)"](
                  groupPublicKey
                )

                expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
                expect(storedGroup.registrationBlockNumber).to.be.equal(
                  tx.blockNumber
                )
                expect(storedGroup.membersHash).to.be.equal(
                  hashUint32Array(dkgResult.members)
                )
              })

              it("should emit GroupRegistered event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "GroupRegistered")
                  .withArgs(0, groupPublicKey)
              })

              it("should reward the submitter with tokens from DKG rewards pool", async () => {
                const currentDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                expect(
                  initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
                ).to.be.equal(params.dkgResultSubmissionReward)

                const currentSubmitterBalance: BigNumber =
                  await testToken.balanceOf(await submitter.getAddress())
                expect(
                  currentSubmitterBalance.sub(initialSubmitterBalance)
                ).to.be.equal(params.dkgResultSubmissionReward)
              })

              it("should unlock the sortition pool", async () => {
                expect(await sortitionPool.isLocked()).to.be.false
              })
            })

            context("when called by a third party", async () => {
              context("when the third party is not yet eligible", async () => {
                before(async () => {
                  await createSnapshot()

                  await mineBlocks(
                    params.dkgSubmitterPrecedencePeriodLength - 1
                  )
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should revert", async () => {
                  await expect(
                    randomBeacon.connect(thirdParty).approveDkgResult(dkgResult)
                  ).to.be.revertedWith(
                    "Only the DKG result submitter can approve the result at this moment"
                  )
                })
              })

              context("when the third party is eligible", async () => {
                let tx: Promise<ContractTransaction>
                let initialDkgRewardsPoolBalance: BigNumber
                let initApproverBalance: BigNumber

                before(async () => {
                  await createSnapshot()

                  await mineBlocks(params.dkgSubmitterPrecedencePeriodLength)
                  initialDkgRewardsPoolBalance =
                    await randomBeacon.dkgRewardsPool()
                  initApproverBalance = await testToken.balanceOf(
                    await thirdParty.getAddress()
                  )
                  tx = randomBeacon
                    .connect(thirdParty)
                    .approveDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should succeed", async () => {
                  await expect(tx).to.not.be.reverted
                })

                it("should pay the reward to the third party", async () => {
                  const currentDkgRewardsPoolBalance =
                    await randomBeacon.dkgRewardsPool()
                  expect(
                    initialDkgRewardsPoolBalance.sub(
                      currentDkgRewardsPoolBalance
                    )
                  ).to.be.equal(params.dkgResultSubmissionReward)

                  const currentApproverBalance = await testToken.balanceOf(
                    await thirdParty.getAddress()
                  )
                  expect(
                    currentApproverBalance.sub(initApproverBalance)
                  ).to.be.equal(params.dkgResultSubmissionReward)
                })
              })
            })
          })
        })

        context("when there was a challenged result before", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DKG.ResultStruct

          // First result is malicious and submitter is also malicious
          const maliciousSubmitter = 1

          // Submit a second result by another submitter
          const anotherSubmitterIndex = 6
          let anotherSubmitter: Signer

          before(async () => {
            await createSnapshot()

            const { dkgResult: maliciousDkgResult } =
              await signAndSubmitArbitraryDkgResult(
                randomBeacon,
                groupPublicKey,
                // Mix signers to make the result malicious.
                mixSigners(await selectGroup(sortitionPool, genesisSeed)),
                startBlock,
                noMisbehaved,
                maliciousSubmitter
              )

            await randomBeacon.challengeDkgResult(maliciousDkgResult)

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
              submitter: anotherSubmitter,
            } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved,
              anotherSubmitterIndex
            ))

            resultSubmissionBlock = tx.blockNumber
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with challenge period not passed", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock +
                  params.dkgResultChallengePeriodLength -
                  1
              )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'Challenge period has not passed yet' error", async () => {
              await expect(
                randomBeacon
                  .connect(anotherSubmitter)
                  .approveDkgResult(dkgResult)
              ).to.be.revertedWith("Challenge period has not passed yet")
            })
          })

          context("with challenge period passed", async () => {
            let tx: ContractTransaction
            let initialDkgRewardsPoolBalance: BigNumber
            let initialSubmitterBalance: BigNumber

            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )

              initialDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()

              initialSubmitterBalance = await testToken.balanceOf(
                await anotherSubmitter.getAddress()
              )

              tx = await randomBeacon
                .connect(anotherSubmitter)
                .approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultApproved event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgResultApproved")
                .withArgs(dkgResultHash, await anotherSubmitter.getAddress())
            })

            it("should register a group", async () => {
              const groupsRegistry = await randomBeacon.getGroupsRegistry()

              expect(groupsRegistry).to.be.lengthOf(1)
              expect(groupsRegistry[0]).to.deep.equal(keccak256(groupPublicKey))

              const storedGroup = await randomBeacon["getGroup(bytes)"](
                groupPublicKey
              )

              expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
              expect(storedGroup.registrationBlockNumber).to.be.equal(
                tx.blockNumber
              )
              expect(storedGroup.membersHash).to.be.equal(
                hashUint32Array(dkgResult.members)
              )
            })

            it("should emit GroupRegistered event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "GroupRegistered")
                .withArgs(0, groupPublicKey)
            })

            it("should reward the submitter with tokens from DKG rewards pool", async () => {
              const currentDkgRewardsPoolBalance =
                await randomBeacon.dkgRewardsPool()
              expect(
                initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
              ).to.be.equal(params.dkgResultSubmissionReward)

              const currentSubmitterBalance: BigNumber =
                await testToken.balanceOf(await anotherSubmitter.getAddress())
              expect(
                currentSubmitterBalance.sub(initialSubmitterBalance)
              ).to.be.equal(params.dkgResultSubmissionReward)
            })

            it("should unlock the sortition pool", async () => {
              expect(await sortitionPool.isLocked()).to.be.false
            })
          })
        })
      })

      context("with max periods duration", async () => {
        let tx: Promise<ContractTransaction>

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout - 1)

          const { dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved
          )

          await mineBlocks(params.dkgResultChallengePeriodLength)

          tx = randomBeacon.connect(submitter).approveDkgResult(dkgResult)
        })

        after(async () => {
          await restoreSnapshot()
        })

        // Just an explicit assertion to make sure transaction passes correctly
        // for max periods duration.
        it("should succeed", async () => {
          await expect(tx).to.not.be.reverted
        })

        it("should unlock the sortition pool", async () => {
          await expect(await sortitionPool.isLocked()).to.be.false
        })
      })

      context("with misbehaved operators", async () => {
        const misbehavedIndices = [2, 9, 11, 30, 60, 64]
        let misbehavedIds
        let tx: ContractTransaction
        let dkgResult: DKG.ResultStruct

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout - 1)

          let members: number[]
          let submitter: SignerWithAddress
          ;({ dkgResult, members, submitter } =
            await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              misbehavedIndices
            ))

          misbehavedIds = misbehavedIndices.map((i) => members[i - 1])

          await mineBlocks(params.dkgResultChallengePeriodLength)
          tx = await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should correctly set a group members hash", async () => {
          const storedGroup = await randomBeacon["getGroup(bytes)"](
            groupPublicKey
          )

          // misbehavedIndices: [2, 9, 11, 30, 60, 64]
          const expectedMembers = [...dkgResult.members]
          expectedMembers.splice(1, 1) // index -1
          expectedMembers.splice(7, 1) // index -2 (cause expectedMembers already shrinked)
          expectedMembers.splice(8, 1) // index -3
          expectedMembers.splice(26, 1) // index -4
          expectedMembers.splice(55, 1) // index -5
          expectedMembers.splice(58, 1) // index -6
          expect(storedGroup.membersHash).to.be.equal(
            hashUint32Array(expectedMembers)
          )
        })

        it("should ban misbehaved operators from sortition pool rewards", async () => {
          const now = await helpers.time.lastBlockTime()
          const expectedUntil = now + params.sortitionPoolRewardsBanDuration

          await expect(tx)
            .to.emit(sortitionPool, "IneligibleForRewards")
            .withArgs(misbehavedIds, expectedUntil)
        })

        it("should clean dkg data", async () => {
          await assertDkgResultCleanData(randomBeacon)
        })
      })

      // This case shouldn't happen in real life. When a result is submitted
      // with invalid order of misbehaved operators it should be challenged.
      context(
        "when misbehaved operators are not in ascending order",
        async () => {
          const misbehavedIndices = [2, 9, 30, 11, 60, 64]

          let dkgResult: DKG.ResultStruct
          let submitter: SignerWithAddress

          before(async () => {
            await createSnapshot()
            await mineBlocksTo(startBlock + constants.offchainDkgTime)
            ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              misbehavedIndices
            ))

            await mineBlocks(params.dkgResultChallengePeriodLength)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should succeed", async () => {
            await expect(
              randomBeacon.connect(submitter).approveDkgResult(dkgResult)
            ).to.not.be.reverted
          })
        }
      )

      context("when misbehaved members contains duplicates", async () => {
        const misbehavedIndices = [2, 9, 9, 10]

        let dkgResult: DKG.ResultStruct
        let submitter: SignerWithAddress

        before(async () => {
          await createSnapshot()
          await mineBlocksTo(startBlock + constants.offchainDkgTime)
          ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            misbehavedIndices
          ))

          await mineBlocks(params.dkgResultChallengePeriodLength)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should succeed", async () => {
          await expect(
            randomBeacon.connect(submitter).approveDkgResult(dkgResult)
          ).to.not.be.reverted
        })
      })
    })

    context(
      "when the balance of DKG rewards pool is smaller than the DKG submission reward",
      async () => {
        let dkgRewardsPoolBalance: BigNumber
        let tx: Promise<ContractTransaction>
        let initApproverBalance: BigNumber
        let submitter: SignerWithAddress

        before(async () => {
          await createSnapshot()

          dkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()

          // Set the DKG result submission reward to twice the amount of test
          // tokens in the DKG rewards pool
          await randomBeaconGovernance.beginDkgResultSubmissionRewardUpdate(
            dkgRewardsPoolBalance.mul(2)
          )
          await helpers.time.increaseTime(params.governanceDelay)
          await randomBeaconGovernance.finalizeDkgResultSubmissionRewardUpdate()

          const [genesisTx, genesisSeed] = await genesis(randomBeacon)
          const startBlock: number = genesisTx.blockNumber
          await mineBlocksTo(startBlock + dkgTimeout - 1)

          let dkgResult: DKG.ResultStruct
          ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved
          ))

          initApproverBalance = await testToken.balanceOf(
            await submitter.getAddress()
          )

          await mineBlocks(params.dkgResultChallengePeriodLength)
          tx = randomBeacon.connect(submitter).approveDkgResult(dkgResult)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should succeed", async () => {
          await expect(tx).to.not.be.reverted
        })

        it("should pay the approver the whole DKG rewards pool balance", async () => {
          expect(await randomBeacon.dkgRewardsPool()).to.be.equal(0)

          const currentApproverBalance = await testToken.balanceOf(
            await submitter.getAddress()
          )
          expect(currentApproverBalance.sub(initApproverBalance)).to.be.equal(
            dkgRewardsPoolBalance
          )
        })
      }
    )
  })

  describe("notifyDkgTimeout", async () => {
    context("with initial contract state", async () => {
      it("should revert with 'DKG has not timed out' error", async () => {
        await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
          "DKG has not timed out"
        )
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with dkg not timed out", async () => {
        context("with off-chain dkg time not passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(startBlock + constants.offchainDkgTime - 1)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'DKG has not timed out' error", async () => {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "DKG has not timed out"
            )
          })
        })

        context("with off-chain dkg time passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(startBlock + constants.offchainDkgTime)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'DKG has not timed out' error", async () => {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "DKG has not timed out"
            )
          })
        })

        context("with result submission period almost ended", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(startBlock + dkgTimeout - 1)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'DKG has not timed out' error", async () => {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "DKG has not timed out"
            )
          })
        })
      })

      context("with dkg timed out", async () => {
        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("called by a third party", async () => {
          let tx: ContractTransaction
          let initialDkgRewardsPoolBalance: BigNumber
          let initialNotifierBalance: BigNumber

          before(async () => {
            await createSnapshot()

            initialDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()

            initialNotifierBalance = await testToken.balanceOf(
              await thirdParty.getAddress()
            )
            tx = await randomBeacon.connect(thirdParty).notifyDkgTimeout()
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should emit DkgTimedOut event", async () => {
            await expect(tx).to.emit(randomBeacon, "DkgTimedOut")
          })

          it("should clean dkg data", async () => {
            await assertDkgResultCleanData(randomBeacon)
          })

          it("should reward the notifier with tokens from DKG rewards pool", async () => {
            const currentDkgRewardsPoolBalance =
              await randomBeacon.dkgRewardsPool()
            expect(
              initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
            ).to.be.equal(params.sortitionPoolUnlockingReward)

            const currentNotifierBalance: BigNumber = await testToken.balanceOf(
              await thirdParty.getAddress()
            )
            expect(
              currentNotifierBalance.sub(initialNotifierBalance)
            ).to.be.equal(params.sortitionPoolUnlockingReward)
          })

          it("should unlock the sortition pool", async () => {
            expect(await sortitionPool.isLocked()).to.be.false
          })
        })
      })
    })
  })

  describe("challengeDkgResult", async () => {
    // Just to make `challengeDkgResult` call possible.
    const stubDkgResult: DKG.ResultStruct = {
      groupPubKey: blsData.groupPubKey,
      members: [1, 2, 3, 4],
      misbehavedMembersIndices: [],
      signatures: "0x01020304",
      signingMembersIndices: [1, 2, 3, 4],
      submitterMemberIndex: 1,
      membersHash: hashDKGMembers([1, 2, 3, 4], []),
    }

    context("with initial contract state", async () => {
      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          randomBeacon.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let genesisSeed

      before("run genesis", async () => {
        await createSnapshot()

        const [genesisTx, seed] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          randomBeacon.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })

      context("with off-chain dkg time passed", async () => {
        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with dkg result not submitted", async () => {
          it("should revert with 'Current state is not CHALLENGE' error", async () => {
            await expect(
              randomBeacon.challengeDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })
        })

        context("with invalid members hash submitted", async () => {
          let dkgResultHash: string
          let dkgResult: DKG.ResultStruct
          let members: Operator[]
          let membersIds: number[]

          before(async () => {
            await createSnapshot()

            members = await selectGroup(sortitionPool, genesisSeed)
            membersIds = members.map((m) => m.id)
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when misbehaved members are present", async () => {
            it("should emit DkgResultChallenged event", async () => {
              const actualMisbehavedIndices = [1, 7, 14, 42]
              const invalidMisbehavedIndices = [3, 8, 15, 41]
              const invalidMembersHash = hashDKGMembers(
                membersIds,
                invalidMisbehavedIndices
              )

              ;({ dkgResult, dkgResultHash } =
                await signAndSubmitArbitraryDkgResult(
                  randomBeacon,
                  groupPublicKey,
                  members,
                  startBlock,
                  actualMisbehavedIndices,
                  undefined,
                  invalidMembersHash
                ))

              const tx = await randomBeacon
                .connect(thirdParty)
                .challengeDkgResult(dkgResult)

              await expect(tx)
                .to.emit(randomBeacon, "DkgResultChallenged")
                .withArgs(
                  dkgResultHash,
                  await thirdParty.getAddress(),
                  "Invalid members hash"
                )
            })
          })

          context("when misbehaved members are not present", async () => {
            it("should emit DkgResultChallenged event", async () => {
              const invalidMisbehavedIndices = [3]
              const invalidMembersHash = hashDKGMembers(
                membersIds,
                invalidMisbehavedIndices
              )

              ;({ dkgResult, dkgResultHash } =
                await signAndSubmitArbitraryDkgResult(
                  randomBeacon,
                  groupPublicKey,
                  members,
                  startBlock,
                  noMisbehaved,
                  undefined,
                  invalidMembersHash
                ))

              const tx = await randomBeacon
                .connect(thirdParty)
                .challengeDkgResult(dkgResult)

              await expect(tx)
                .to.emit(randomBeacon, "DkgResultChallenged")
                .withArgs(
                  dkgResultHash,
                  await thirdParty.getAddress(),
                  "Invalid members hash"
                )
            })
          })
        })

        context("with malicious dkg result submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DKG.ResultStruct
          let submitter: SignerWithAddress

          before(async () => {
            await createSnapshot()

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
              submitter,
            } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious.
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            ))

            resultSubmissionBlock = tx.blockNumber
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("at the beginning of challenge period", async () => {
            context("called by a third party", async () => {
              let tx: ContractTransaction
              before(async () => {
                await createSnapshot()

                tx = await randomBeacon
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "Invalid group members"
                  )
              })

              it("should not unlock the sortition pool", async () => {
                expect(await sortitionPool.isLocked()).to.be.true
              })

              it("should emit DkgMaliciousResultSlashed event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "DkgMaliciousResultSlashed")
                  .withArgs(dkgResultHash, to1e18(50000), submitter.address)
              })

              it("should slash malicious result submitter", async () => {
                await expect(tx)
                  .to.emit(staking, "Seized")
                  .withArgs(to1e18(50000), 100, await thirdParty.getAddress(), [
                    submitter.address,
                  ])
              })
            })
          })

          context("at the end of challenge period", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock +
                  params.dkgResultChallengePeriodLength -
                  1
              )
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("called by a third party", async () => {
              let tx: ContractTransaction
              before(async () => {
                await createSnapshot()

                tx = await randomBeacon
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "Invalid group members"
                  )
              })

              it("should not unlock the sortition pool", async () => {
                expect(await sortitionPool.isLocked()).to.be.true
              })

              it("should emit DkgMaliciousResultSlashed event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "DkgMaliciousResultSlashed")
                  .withArgs(dkgResultHash, to1e18(50000), submitter.address)
              })

              it("should slash malicious result submitter", async () => {
                await expect(tx)
                  .to.emit(staking, "Seized")
                  .withArgs(to1e18(50000), 100, await thirdParty.getAddress(), [
                    submitter.address,
                  ])
              })
            })
          })

          context("with challenge period passed", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'Challenge period has already passed' error", async () => {
              await expect(
                randomBeacon.challengeDkgResult(dkgResult)
              ).to.be.revertedWith("Challenge period has already passed")
            })
          })

          context("with token staking seize call failure", async () => {
            let tokenStakingFake: FakeContract<IRandomBeaconStaking>
            let tx: Promise<ContractTransaction>

            before(async () => {
              await createSnapshot()

              tokenStakingFake = await fakeTokenStaking(randomBeacon)
              tokenStakingFake.seize.reverts()

              tx = randomBeacon
                .connect(thirdParty)
                .challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()

              tokenStakingFake.seize.reset()
            })

            it("should succeed", async () => {
              await expect(tx).to.not.be.reverted
            })

            it("should emit DkgMaliciousResultSlashingFailed", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgMaliciousResultSlashingFailed")
                .withArgs(
                  dkgResultHash,
                  params.maliciousDkgResultSlashingAmount,
                  submitter.address
                )
            })
          })

          context(
            "with challenged result not matching the submitted one",
            async () => {
              it("should revert with 'Result under challenge is different than the submitted one'", async () => {
                const modifiedDkgResult: DKG.ResultStruct = { ...dkgResult }
                const modifiedMembersHash = hashUint32Array(
                  modifiedDkgResult.members.splice(42, 1)
                )
                modifiedDkgResult.membersHash = modifiedMembersHash

                await expect(
                  randomBeacon.challengeDkgResult(modifiedDkgResult)
                ).to.be.revertedWith(
                  "Result under challenge is different than the submitted one"
                )
              })
            }
          )
        })

        context(
          "with dkg result submitted with unrecoverable signatures",
          async () => {
            let dkgResultHash: string
            let dkgResult: DKG.ResultStruct
            let submitter: SignerWithAddress
            let tx: ContractTransaction

            before(async () => {
              await createSnapshot()
              ;({ dkgResult, dkgResultHash, submitter } =
                await signAndSubmitUnrecoverableDkgResult(
                  randomBeacon,
                  groupPublicKey,
                  await selectGroup(sortitionPool, genesisSeed),
                  startBlock,
                  noMisbehaved
                ))

              tx = await randomBeacon
                .connect(thirdParty)
                .challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultChallenged event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgResultChallenged")
                .withArgs(
                  dkgResultHash,
                  await thirdParty.getAddress(),
                  "validation reverted"
                )
            })

            it("should not unlock the sortition pool", async () => {
              expect(await sortitionPool.isLocked()).to.be.true
            })

            it("should emit DkgMaliciousResultSlashed event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgMaliciousResultSlashed")
                .withArgs(dkgResultHash, to1e18(50000), submitter.address)
            })

            it("should slash malicious result submitter", async () => {
              await expect(tx)
                .to.emit(staking, "Seized")
                .withArgs(to1e18(50000), 100, await thirdParty.getAddress(), [
                  submitter.address,
                ])
            })
          }
        )

        context("with correct dkg result submitted", async () => {
          let dkgResult: DKG.ResultStruct

          before(async () => {
            await createSnapshot()
            ;({ dkgResult } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            ))
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'unjustified challenge' error", async () => {
            await expect(
              randomBeacon.challengeDkgResult(dkgResult)
            ).to.be.revertedWith("unjustified challenge")
          })
        })
      })
    })

    // This test checks that dkg timeout is adjusted in case of result challenges
    // to include the offset blocks that were mined until the invalid result
    // was challenged.
    it("should enforce submission start offset", async () => {
      await createSnapshot()

      let dkgResult: DKG.ResultStruct

      const [genesisTx] = await genesis(randomBeacon)
      const startBlock = genesisTx.blockNumber

      await mineBlocks(constants.offchainDkgTime)

      // Submit result 1 at the beginning of the submission period
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock,
        noMisbehaved
      ))

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 1 after submission"
      ).to.equal(0)

      // Challenge result 1 at the beginning of the challenge period
      await randomBeacon.challengeDkgResult(dkgResult)
      // 1 block for dkg result submission tx +
      // 1 block for challenge tx
      let expectedSubmissionOffset = 2

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 1 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 2 in the middle of the submission period
      let blocksToMine = params.dkgResultSubmissionTimeout / 2
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock,
        noMisbehaved
      ))

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 2 after submission"
      ).to.equal(expectedSubmissionOffset) // same as before

      expectedSubmissionOffset += blocksToMine

      // Challenge result 2 in the middle of the challenge period
      await mineBlocks(params.dkgResultChallengePeriodLength / 2)
      expectedSubmissionOffset += params.dkgResultChallengePeriodLength / 2
      await randomBeacon.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 2 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 3 at the end of the submission period
      blocksToMine = params.dkgResultSubmissionTimeout - 1
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock,
        noMisbehaved
      ))

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after submission"
      ).to.equal(expectedSubmissionOffset) // same as before

      expectedSubmissionOffset += blocksToMine

      // Challenge result 3 at the end of the challenge period
      blocksToMine = params.dkgResultChallengePeriodLength - 1
      await mineBlocks(blocksToMine)
      expectedSubmissionOffset += blocksToMine

      await expect(
        randomBeacon.callStatic.notifyDkgTimeout()
      ).to.be.revertedWith("DKG has not timed out")

      await randomBeacon.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 4 after the submission period
      blocksToMine = params.dkgResultSubmissionTimeout
      await mineBlocks(blocksToMine)
      await expect(
        signAndSubmitArbitraryDkgResult(
          randomBeacon,
          groupPublicKey,
          signers,
          startBlock,
          noMisbehaved
        )
      ).to.be.revertedWith("DKG timeout already passed")

      await randomBeacon.notifyDkgTimeout()

      await restoreSnapshot()
    })
  })

  describe("fundDkgRewardsPool", () => {
    const amount = to1e18(1000)

    let previousDkgRewardsPoolBalance: BigNumber
    let previousRandomBeaconBalance: BigNumber

    before(async () => {
      await createSnapshot()

      previousDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()
      previousRandomBeaconBalance = await testToken.balanceOf(
        randomBeacon.address
      )

      await testToken.mint(await thirdParty.getAddress(), amount)
      await testToken.connect(thirdParty).approve(randomBeacon.address, amount)

      await randomBeacon.fundDkgRewardsPool(
        await thirdParty.getAddress(),
        amount
      )
    })

    after(async () => {
      await restoreSnapshot()
    })

    it("should increase the DKG rewards pool balance", async () => {
      const currentDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()
      expect(
        currentDkgRewardsPoolBalance.sub(previousDkgRewardsPoolBalance)
      ).to.be.equal(amount)
    })

    it("should transfer tokens to the random beacon contract", async () => {
      const currentRandomBeaconBalance = await testToken.balanceOf(
        randomBeacon.address
      )
      expect(
        currentRandomBeaconBalance.sub(previousRandomBeaconBalance)
      ).to.be.equal(amount)
    })
  })

  async function fundDkgRewardsPool(donateAmount: BigNumber) {
    await testToken.mint(await thirdParty.getAddress(), donateAmount)
    await testToken
      .connect(thirdParty)
      .approve(randomBeacon.address, donateAmount)

    await randomBeacon.fundDkgRewardsPool(
      await thirdParty.getAddress(),
      donateAmount
    )
  }
})

async function assertDkgResultCleanData(randomBeacon: {
  getDkgData: () => Promise<DKG.DataStructOutput>
}) {
  const dkgData: DKG.DataStructOutput = await randomBeacon.getDkgData()

  expect(
    dkgData.parameters.resultChallengePeriodLength,
    "unexpected resultChallengePeriodLength"
  ).to.eq(params.dkgResultChallengePeriodLength)

  expect(
    dkgData.parameters.resultSubmissionTimeout,
    "unexpected resultSubmissionTimeout"
  ).to.eq(params.dkgResultSubmissionTimeout)

  expect(
    dkgData.parameters.submitterPrecedencePeriodLength,
    "unexpected submitterPrecedencePeriodLength"
  ).to.eq(params.dkgSubmitterPrecedencePeriodLength)

  expect(dkgData.startBlock, "unexpected startBlock").to.eq(0)

  expect(
    dkgData.resultSubmissionStartBlockOffset,
    "unexpected resultSubmissionStartBlockOffset"
  ).to.eq(0)

  expect(dkgData.submittedResultHash, "unexpected submittedResultHash").to.eq(
    ethers.constants.HashZero
  )

  expect(dkgData.submittedResultBlock, "unexpected submittedResultBlock").to.eq(
    0
  )
}

function mixSigners(signers: Operator[]): Operator[] {
  return signers
    .map((v) => ({ v, sort: Math.random() }))
    .sort((a, b) => a.sort - b.sort)
    .map(({ v }) => v)
}
