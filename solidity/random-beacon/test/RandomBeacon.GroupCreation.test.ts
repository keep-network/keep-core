/* eslint-disable @typescript-eslint/no-unused-expressions, @typescript-eslint/no-extra-semi */

import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import type { BigNumber, ContractTransaction, Signer } from "ethers"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import blsData from "./data/bls"
import { constants, dkgState, params, testDeployment } from "./fixtures"
import type {
  RandomBeacon,
  RandomBeaconGovernance,
  RandomBeaconStub,
  TestToken,
  SortitionPool,
  StakingStub,
} from "../typechain"
import {
  genesis,
  signAndSubmitCorrectDkgResult,
  signAndSubmitArbitraryDkgResult,
  DkgResult,
  noMisbehaved,
  getDkgResultSubmitterSigner,
} from "./utils/dkg"
import { registerOperators, Operator } from "./utils/operators"
import { selectGroup } from "./utils/groups"

const { mineBlocks, mineBlocksTo } = helpers.time
const { to1e18 } = helpers.number
const { keccak256 } = ethers.utils

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
  const randomBeacon = contracts.randomBeacon as RandomBeaconStub & RandomBeacon
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
    constants.offchainDkgTime +
    constants.groupSize * params.dkgResultSubmissionEligibilityDelay

  const dkgResultSubmissionReward = to1e18(5)
  const sortitionPoolUnlockingReward = to1e18(10)

  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let thirdParty: Signer
  let signers: Operator[]

  let randomBeaconGovernance: RandomBeaconGovernance
  let randomBeacon: RandomBeaconStub & RandomBeacon
  let sortitionPool: SortitionPool
  let staking: StakingStub
  let testToken: TestToken

  before(async () => {
    thirdParty = await ethers.getSigner((await getUnnamedAccounts())[0])
  })

  beforeEach("load test fixture", async () => {
    ;({
      randomBeaconGovernance,
      randomBeacon,
      sortitionPool,
      staking,
      testToken,
      signers,
    } = await waffle.loadFixture(fixture))

    await randomBeaconGovernance.beginDkgResultSubmissionRewardUpdate(
      dkgResultSubmissionReward
    )
    await randomBeaconGovernance.beginSortitionPoolUnlockingRewardUpdate(
      sortitionPoolUnlockingReward
    )
    await helpers.time.increaseTime(12 * 60 * 60)
    await randomBeaconGovernance.finalizeDkgResultSubmissionRewardUpdate()
    await randomBeaconGovernance.finalizeSortitionPoolUnlockingRewardUpdate()

    // Fund DKG rewards pool to make testing of rewards possible.
    const dkgRewardsPoolDonate = to1e18(100)
    await testToken.mint(await thirdParty.getAddress(), dkgRewardsPoolDonate)
    await testToken
      .connect(thirdParty)
      .approve(randomBeacon.address, dkgRewardsPoolDonate)
    await randomBeacon.fundDkgRewardsPool(
      await thirdParty.getAddress(),
      dkgRewardsPoolDonate
    )
  })

  describe("genesis", async () => {
    context("when called by a third party", async () => {
      it("should succeed", async () => {
        await randomBeacon.connect(thirdParty).genesis()
      })
    })

    context("with initial contract state", async () => {
      let tx: ContractTransaction
      let expectedSeed: BigNumber

      beforeEach("run genesis", async () => {
        ;[tx, expectedSeed] = await genesis(randomBeacon)
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

      beforeEach("run genesis", async () => {
        const [genesisTx, seed] = await genesis(randomBeacon)
        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      context("with dkg result not submitted", async () => {
        it("should revert with 'current state is not IDLE' error", async () => {
          await expect(randomBeacon.genesis()).to.be.revertedWith(
            "current state is not IDLE"
          )
        })
      })

      context("with dkg result submitted", async () => {
        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)
          await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved
          )
        })

        // TODO: Add test cases to cover results that are approved, challenged or
        // pending.

        context("with dkg result not approved", async () => {
          it("should revert with 'current state is not IDLE' error", async () => {
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
      it("should return IDLE state", async () => {
        expect(await randomBeacon.getGroupCreationState()).to.be.equal(
          dkgState.IDLE
        )
      })
    })

    context("when genesis dkg started", async () => {
      let startBlock: number
      let genesisSeed

      beforeEach("run genesis", async () => {
        const [genesisTx, seed] = await genesis(randomBeacon)
        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      context("at the start of off-chain dkg period", async () => {
        it("should return KEY_GENERATION state", async () => {
          expect(await randomBeacon.getGroupCreationState()).to.be.equal(
            dkgState.KEY_GENERATION
          )
        })
      })

      context("at the end of off-chain dkg period", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        it("should return KEY_GENERATION state", async () => {
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
          it("should return AWAITING_RESULT state", async () => {
            expect(await randomBeacon.getGroupCreationState()).to.be.equal(
              dkgState.AWAITING_RESULT
            )
          })

          context("after the dkg timeout period", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + dkgTimeout + 1)
            })

            it("should return AWAITING_RESULT state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.AWAITING_RESULT
              )
            })
          })
        })

        context("when dkg result was submitted", async () => {
          let dkgResult: DkgResult

          beforeEach(async () => {
            ;({ dkgResult } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            ))
          })

          context("when dkg result was not approved", async () => {
            it("should return CHALLENGE state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.CHALLENGE
              )
            })
          })

          context("when dkg result was approved", async () => {
            beforeEach(async () => {
              await mineBlocks(params.dkgResultChallengePeriodLength)

              const submitter = await getDkgResultSubmitterSigner(
                randomBeacon,
                dkgResult
              )

              await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
            })

            it("should return IDLE state", async () => {
              expect(await randomBeacon.getGroupCreationState()).to.be.equal(
                dkgState.IDLE
              )
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DkgResult

          beforeEach(async () => {
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            ))
          })

          context("when dkg result was challenged", async () => {
            beforeEach(async () => {
              await randomBeacon.challengeDkgResult(dkgResult)
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

      beforeEach("run genesis", async () => {
        const [genesisTx, seed] = await genesis(randomBeacon)
        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      context("within off-chain dkg period", async () => {
        it("should return false", async () => {
          await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
        })
      })

      context("after off-chain dkg period", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime + 1)
        })

        context("when dkg result was not submitted", async () => {
          it("should return false", async () => {
            await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
          })

          context("at the end of the dkg timeout period", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + dkgTimeout)
            })

            it("should return false", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
            })
          })

          context("after the dkg timeout period", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + dkgTimeout + 1)
            })

            it("should return true", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.true
            })
          })
        })

        context("when dkg result was submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResult: DkgResult

          beforeEach(async () => {
            let tx: ContractTransaction
            ;({ transaction: tx, dkgResult } =
              await signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              ))

            resultSubmissionBlock = tx.blockNumber
          })

          context("when dkg result was not approved", async () => {
            context("at the end of the dkg timeout period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout)
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after the dkg timeout period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout + 1)
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("at the end of the challenge period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(
                  resultSubmissionBlock + params.dkgResultChallengePeriodLength
                )
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after the challenge period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(
                  resultSubmissionBlock +
                    params.dkgResultChallengePeriodLength +
                    1
                )
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })
          })

          context("when dkg result was approved", async () => {
            beforeEach(async () => {
              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )

              const submitter = await getDkgResultSubmitterSigner(
                randomBeacon,
                dkgResult
              )

              await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
            })

            it("should return false", async () => {
              await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DkgResult

          beforeEach(async () => {
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious.
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            ))
          })

          context("when dkg result was challenged", async () => {
            let challengeBlockNumber: number

            beforeEach(async () => {
              const tx = await randomBeacon.challengeDkgResult(dkgResult)
              challengeBlockNumber = tx.blockNumber
            })

            context("at the end of dkg result submission period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(
                  challengeBlockNumber +
                    constants.groupSize *
                      params.dkgResultSubmissionEligibilityDelay
                )
              })

              it("should return false", async () => {
                await expect(await randomBeacon.hasDkgTimedOut()).to.be.false
              })
            })

            context("after dkg result submission period", async () => {
              beforeEach(async () => {
                await mineBlocksTo(
                  challengeBlockNumber +
                    constants.groupSize *
                      params.dkgResultSubmissionEligibilityDelay +
                    1
                )
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
    // TODO: Add more tests to cover the DKG result verification function thoroughly.
    // TODO: Add tests to cover misbehaved members

    context("with initial contract state", async () => {
      it("should revert with 'current state is not AWAITING_RESULT' error", async () => {
        await expect(
          signAndSubmitArbitraryDkgResult(
            randomBeacon,
            groupPublicKey,
            signers,
            1,
            noMisbehaved
          )
        ).to.be.revertedWith("current state is not AWAITING_RESULT")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let genesisSeed

      beforeEach("run genesis", async () => {
        const [genesisTx, seed] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      context("with group creation not timed out", async () => {
        context("with off-chain dkg time not passed", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime - 1)
          })

          it("should revert with 'current state is not AWAITING_RESULT' error", async () => {
            await expect(
              signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )
            ).to.be.revertedWith("current state is not AWAITING_RESULT")
          })
        })

        context("with off-chain dkg time passed", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime)
          })

          context("with enough signatures on the result", async () => {
            let tx: ContractTransaction
            let dkgResult: DkgResult
            let dkgResultHash: string

            beforeEach(async () => {
              const filteredSigners = signers.slice(
                0,
                constants.signatureThreshold
              )

              ;({
                transaction: tx,
                dkgResult,
                dkgResultHash,
              } = await signAndSubmitArbitraryDkgResult(
                randomBeacon,
                groupPublicKey,
                filteredSigners,
                startBlock,
                noMisbehaved,
                1,
                33
              ))
            })

            it("should succeed", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgResultSubmitted")
                .withArgs(
                  dkgResultHash,
                  genesisSeed,
                  dkgResult.submitterMemberIndex,
                  dkgResult.groupPubKey,
                  dkgResult.misbehavedMembersIndices,
                  dkgResult.signatures,
                  dkgResult.signingMembersIndices,
                  dkgResult.members
                )
            })

            it("should register a candidate group", async () => {
              const groupsRegistry = await randomBeacon.getGroupsRegistry()

              expect(groupsRegistry).to.be.lengthOf(1)
              expect(groupsRegistry[0]).to.deep.equal(keccak256(groupPublicKey))

              const storedGroup = await randomBeacon["getGroup(bytes)"](
                groupPublicKey
              )

              expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
              expect(storedGroup.activationBlockNumber).to.be.equal(0)
              expect(storedGroup.members).to.be.deep.equal(dkgResult.members)
            })
          })

          it("should succeed for the first submitter", async () => {
            const submitterIndex = 1

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
            await expect(tx)
              .to.emit(randomBeacon, "DkgResultSubmitted")
              .withArgs(
                dkgResultHash,
                genesisSeed,
                dkgResult.submitterMemberIndex,
                dkgResult.groupPubKey,
                dkgResult.misbehavedMembersIndices,
                dkgResult.signatures,
                dkgResult.signingMembersIndices,
                dkgResult.members
              )
          })

          it("should revert for the second submitter", async () => {
            await expect(
              signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved,
                2
              )
            ).to.be.revertedWith("Submitter not eligible")
          })

          it("should register a candidate group", async () => {
            const { dkgResult } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            )

            const groupsRegistry = await randomBeacon.getGroupsRegistry()

            expect(groupsRegistry).to.be.lengthOf(1)
            expect(groupsRegistry[0]).to.deep.equal(keccak256(groupPublicKey))

            const storedGroup = await randomBeacon["getGroup(bytes)"](
              groupPublicKey
            )

            expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
            expect(storedGroup.activationBlockNumber).to.be.equal(0)
            expect(storedGroup.members).to.be.deep.equal(dkgResult.members)
          })

          it("should emit CandidateGroupRegistered event", async () => {
            const { transaction: tx } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            )

            await expect(tx)
              .to.emit(randomBeacon, "CandidateGroupRegistered")
              .withArgs(groupPublicKey)
          })

          it("should not unlock the sortition pool", async () => {
            await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            )

            expect(await sortitionPool.isLocked()).to.be.true
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

              it("should succeed for the first submitter", async () => {
                const submitterIndex = 1

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

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(
                    dkgResultHash,
                    genesisSeed,
                    dkgResult.submitterMemberIndex,
                    dkgResult.groupPubKey,
                    dkgResult.misbehavedMembersIndices,
                    dkgResult.signatures,
                    dkgResult.signingMembersIndices,
                    dkgResult.members
                  )
              })

              it("should revert for the second submitter", async () => {
                await expect(
                  signAndSubmitCorrectDkgResult(
                    randomBeacon,
                    groupPublicKey,
                    genesisSeed,
                    startBlock,
                    noMisbehaved,
                    2
                  )
                ).to.be.revertedWith("Submitter not eligible")
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

              it("should succeed for the first submitter", async () => {
                const submitterIndex = 1

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

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(
                    dkgResultHash,
                    genesisSeed,
                    dkgResult.submitterMemberIndex,
                    dkgResult.groupPubKey,
                    dkgResult.misbehavedMembersIndices,
                    dkgResult.signatures,
                    dkgResult.signingMembersIndices,
                    dkgResult.members
                  )
              })

              it("should succeed for the second submitter", async () => {
                const submitterIndex = 2

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

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(
                    dkgResultHash,
                    genesisSeed,
                    dkgResult.submitterMemberIndex,
                    dkgResult.groupPubKey,
                    dkgResult.misbehavedMembersIndices,
                    dkgResult.signatures,
                    dkgResult.signingMembersIndices,
                    dkgResult.members
                  )
              })

              it("should revert for the third submitter", async () => {
                await expect(
                  signAndSubmitCorrectDkgResult(
                    randomBeacon,
                    groupPublicKey,
                    genesisSeed,
                    startBlock,
                    noMisbehaved,
                    3
                  )
                ).to.be.revertedWith("Submitter not eligible")
              })
            }
          )

          context(
            "with the last submitter eligibility delay period almost ended",
            async () => {
              beforeEach(async () => {
                await mineBlocksTo(startBlock + dkgTimeout - 1)
              })

              it("should succeed for the first submitter", async () => {
                const submitterIndex = 2

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

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(
                    dkgResultHash,
                    genesisSeed,
                    dkgResult.submitterMemberIndex,
                    dkgResult.groupPubKey,
                    dkgResult.misbehavedMembersIndices,
                    dkgResult.signatures,
                    dkgResult.signingMembersIndices,
                    dkgResult.members
                  )
              })

              it("should succeed for the last submitter", async () => {
                const submitterIndex = constants.groupSize

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

                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultSubmitted")
                  .withArgs(
                    dkgResultHash,
                    genesisSeed,
                    dkgResult.submitterMemberIndex,
                    dkgResult.groupPubKey,
                    dkgResult.misbehavedMembersIndices,
                    dkgResult.signatures,
                    dkgResult.signingMembersIndices,
                    dkgResult.members
                  )
              })
            }
          )

          context("with dkg result approved", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + constants.offchainDkgTime)

              await signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )
            })

            it("should revert 'current state is not AWAITING_RESULT' error", async () => {
              await expect(
                signAndSubmitCorrectDkgResult(
                  randomBeacon,
                  groupPublicKey,
                  genesisSeed,
                  startBlock,
                  noMisbehaved
                )
              ).to.be.revertedWith("current state is not AWAITING_RESULT")
            })
          })

          context("with dkg result challenged", async () => {
            beforeEach(async () => {
              await mineBlocksTo(startBlock + constants.offchainDkgTime)

              const { dkgResult } = await signAndSubmitArbitraryDkgResult(
                randomBeacon,
                groupPublicKey,
                // Mix signers to make the result malicious.
                mixSigners(await selectGroup(sortitionPool, genesisSeed)),
                startBlock,
                noMisbehaved
              )

              await randomBeacon.challengeDkgResult(dkgResult)
            })

            it("should allow first member to submit", async () => {
              const submitterIndex = 1

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

              await expect(tx)
                .to.emit(randomBeacon, "DkgResultSubmitted")
                .withArgs(
                  dkgResultHash,
                  genesisSeed,
                  dkgResult.submitterMemberIndex,
                  dkgResult.groupPubKey,
                  dkgResult.misbehavedMembersIndices,
                  dkgResult.signatures,
                  dkgResult.signingMembersIndices,
                  dkgResult.members
                )
            })

            it("should register a candidate group", async () => {
              const { dkgResult } = await signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )

              const groupsRegistry = await randomBeacon.getGroupsRegistry()

              expect(groupsRegistry).to.be.lengthOf(1)
              expect(groupsRegistry[0]).to.deep.equal(keccak256(groupPublicKey))

              const storedGroup = await randomBeacon["getGroup(bytes)"](
                groupPublicKey
              )

              expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
              expect(storedGroup.activationBlockNumber).to.be.equal(0)
              expect(storedGroup.members).to.be.deep.equal(dkgResult.members)
            })

            it("should emit CandidateGroupRegistered event", async () => {
              const { transaction: tx } = await signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )

              await expect(tx)
                .to.emit(randomBeacon, "CandidateGroupRegistered")
                .withArgs(groupPublicKey)
            })
          })

          context("with misbehaved members", async () => {
            let tx: ContractTransaction
            let dkgResult: DkgResult
            let dkgResultHash: string

            const misbehavedIndices = [2, 9, 11, 30, 60, 64]

            beforeEach(async () => {
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

            it("should succeed with misbehaved members", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgResultSubmitted")
                .withArgs(
                  dkgResultHash,
                  genesisSeed,
                  dkgResult.submitterMemberIndex,
                  dkgResult.groupPubKey,
                  dkgResult.misbehavedMembersIndices,
                  dkgResult.signatures,
                  dkgResult.signingMembersIndices,
                  dkgResult.members
                )
            })
          })
        })
      })

      // TODO: Check challenge adjust start block calculation for eligibility
      // TODO: Check that challenges add up the delay

      context("with group creation timed out", async () => {
        beforeEach("increase time", async () => {
          await mineBlocksTo(startBlock + dkgTimeout)
        })

        context("with timeout not notified", async () => {
          it("should revert with dkg timeout already passed error", async () => {
            await expect(
              signAndSubmitCorrectDkgResult(
                randomBeacon,
                groupPublicKey,
                genesisSeed,
                startBlock,
                noMisbehaved
              )
            ).to.be.revertedWith("dkg timeout already passed")
          })
        })
      })
    })
  })

  describe("approveDkgResult", async () => {
    // Just to make `approveDkgResult` call possible.
    const stubDkgResult: DkgResult = {
      groupPubKey: blsData.groupPubKey,
      members: [1, 2, 3, 4],
      misbehavedMembersIndices: [],
      signatures: "0x01020304",
      signingMembersIndices: [1, 2, 3, 4],
      submitterMemberIndex: 1,
    }

    context("with initial contract state", async () => {
      it("should revert with 'current state is not CHALLENGE' error", async () => {
        const submitter = await getDkgResultSubmitterSigner(
          randomBeacon,
          stubDkgResult
        )

        await expect(
          randomBeacon.connect(submitter).approveDkgResult(stubDkgResult)
        ).to.be.revertedWith("current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let genesisSeed: BigNumber

      beforeEach("run genesis", async () => {
        const [genesisTx, seed] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      it("should revert with 'current state is not CHALLENGE' error", async () => {
        const submitter = await getDkgResultSubmitterSigner(
          randomBeacon,
          stubDkgResult
        )

        await expect(
          randomBeacon.connect(submitter).approveDkgResult(stubDkgResult)
        ).to.be.revertedWith("current state is not CHALLENGE")
      })

      context("with off-chain dkg time passed", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        context("with dkg result not submitted", async () => {
          it("should revert with 'current state is not CHALLENGE' error", async () => {
            const submitter = await getDkgResultSubmitterSigner(
              randomBeacon,
              stubDkgResult
            )

            await expect(
              randomBeacon.connect(submitter).approveDkgResult(stubDkgResult)
            ).to.be.revertedWith("current state is not CHALLENGE")
          })
        })

        context("with dkg result submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DkgResult
          let submitter: SignerWithAddress
          const submitterIndex = 1

          beforeEach(async () => {
            let tx: ContractTransaction
            ;({
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
            ))

            submitter = await getDkgResultSubmitterSigner(
              randomBeacon,
              dkgResult
            )

            resultSubmissionBlock = tx.blockNumber
          })

          context("with challenge period not passed", async () => {
            beforeEach(async () => {
              await mineBlocksTo(
                resultSubmissionBlock +
                  params.dkgResultChallengePeriodLength -
                  1
              )
            })

            it("should revert with 'challenge period has not passed yet' error", async () => {
              await expect(
                randomBeacon.connect(submitter).approveDkgResult(dkgResult)
              ).to.be.revertedWith("challenge period has not passed yet")
            })
          })

          context("with challenge period passed", async () => {
            beforeEach(async () => {
              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )
            })

            context("when called by a DKG result submitter", async () => {
              let tx: ContractTransaction
              let initialDkgRewardsPoolBalance: BigNumber
              let initialSubmitterBalance: BigNumber

              beforeEach(async () => {
                initialDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                initialSubmitterBalance = await testToken.balanceOf(
                  await submitter.getAddress()
                )
                tx = await randomBeacon
                  .connect(submitter)
                  .approveDkgResult(dkgResult)
              })

              it("should emit DkgResultApproved event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "DkgResultApproved")
                  .withArgs(dkgResultHash, await submitter.getAddress())
              })

              it("should clean dkg data", async () => {
                await assertDkgResultCleanData(randomBeacon)
              })

              it("should activate a candidate group", async () => {
                const storedGroup = await randomBeacon["getGroup(bytes)"](
                  groupPublicKey
                )

                expect(storedGroup.activationBlockNumber).to.be.equal(
                  tx.blockNumber
                )
              })

              it("should reward the submitter with tokens from DKG rewards pool", async () => {
                const currentDkgRewardsPoolBalance =
                  await randomBeacon.dkgRewardsPool()
                expect(
                  initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
                ).to.be.equal(dkgResultSubmissionReward)

                const currentSubmitterBalance: BigNumber =
                  await testToken.balanceOf(await submitter.getAddress())
                expect(
                  currentSubmitterBalance.sub(initialSubmitterBalance)
                ).to.be.equal(dkgResultSubmissionReward)
              })

              it("should emit GroupActivated event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "GroupActivated")
                  .withArgs(0, groupPublicKey)
              })

              it("should unlock the sortition pool", async () => {
                expect(await sortitionPool.isLocked()).to.be.false
              })
            })

            context("when called by a third party", async () => {
              context("when the third party is not yet eligible", async () => {
                beforeEach(async () => {
                  await mineBlocks(
                    params.relayEntrySubmissionEligibilityDelay - 1
                  )
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
                let tx: ContractTransaction
                let initialDkgRewardsPoolBalance: BigNumber
                let initApproverBalance: BigNumber

                beforeEach(async () => {
                  await mineBlocks(params.relayEntrySubmissionEligibilityDelay)
                  initialDkgRewardsPoolBalance =
                    await randomBeacon.dkgRewardsPool()
                  initApproverBalance = await testToken.balanceOf(
                    await thirdParty.getAddress()
                  )
                  tx = await randomBeacon
                    .connect(thirdParty)
                    .approveDkgResult(dkgResult)
                })

                it("should succeed", async () => {
                  await expect(tx)
                    .to.emit(randomBeacon, "GroupActivated")
                    .withArgs(0, groupPublicKey)
                })

                it("should pay the reward to the third party", async () => {
                  const currentDkgRewardsPoolBalance =
                    await randomBeacon.dkgRewardsPool()
                  expect(
                    initialDkgRewardsPoolBalance.sub(
                      currentDkgRewardsPoolBalance
                    )
                  ).to.be.equal(dkgResultSubmissionReward)

                  const currentApproverBalance = await testToken.balanceOf(
                    await thirdParty.getAddress()
                  )
                  expect(
                    currentApproverBalance.sub(initApproverBalance)
                  ).to.be.equal(dkgResultSubmissionReward)
                })
              })
            })
          })
        })

        context("when there was a challenged result before", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DkgResult

          // First result is malicious and submitter is also malicious
          const maliciousSubmitter = 1

          // Submit a second result by another submitter
          const anotherSubmitterIndex = 5
          let anotherSubmitter: Signer

          beforeEach(async () => {
            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
            } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious.
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved,
              maliciousSubmitter
            ))

            await randomBeacon.challengeDkgResult(dkgResult)

            await mineBlocks(
              params.dkgResultSubmissionEligibilityDelay * anotherSubmitterIndex
            )
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
            } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved,
              anotherSubmitterIndex
            ))

            anotherSubmitter = await ethers.getSigner(
              await sortitionPool.getIDOperator(
                dkgResult.members[anotherSubmitterIndex - 1]
              )
            )
            resultSubmissionBlock = tx.blockNumber
          })

          context("with challenge period not passed", async () => {
            beforeEach(async () => {
              await mineBlocksTo(
                resultSubmissionBlock +
                  params.dkgResultChallengePeriodLength -
                  1
              )
            })

            it("should revert with 'challenge period has not passed yet' error", async () => {
              await expect(
                randomBeacon
                  .connect(anotherSubmitter)
                  .approveDkgResult(dkgResult)
              ).to.be.revertedWith("challenge period has not passed yet")
            })
          })

          context("with challenge period passed", async () => {
            let tx: ContractTransaction
            let initialDkgRewardsPoolBalance: BigNumber
            let initialSubmitterBalance: BigNumber

            beforeEach(async () => {
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

            it("should emit DkgResultApproved event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "DkgResultApproved")
                .withArgs(dkgResultHash, await anotherSubmitter.getAddress())
            })

            it("should activate a candidate group", async () => {
              const storedGroup = await randomBeacon["getGroup(bytes)"](
                groupPublicKey
              )

              expect(storedGroup.activationBlockNumber).to.be.equal(
                tx.blockNumber
              )
            })

            it("should reward the submitter with tokens from DKG rewards pool", async () => {
              const currentDkgRewardsPoolBalance =
                await randomBeacon.dkgRewardsPool()
              expect(
                initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
              ).to.be.equal(dkgResultSubmissionReward)

              const currentSubmitterBalance: BigNumber =
                await testToken.balanceOf(await anotherSubmitter.getAddress())
              expect(
                currentSubmitterBalance.sub(initialSubmitterBalance)
              ).to.be.equal(dkgResultSubmissionReward)
            })

            it("should emit GroupActivated event", async () => {
              await expect(tx)
                .to.emit(randomBeacon, "GroupActivated")
                .withArgs(0, groupPublicKey)
            })

            it("should unlock the sortition pool", async () => {
              expect(await sortitionPool.isLocked()).to.be.false
            })
          })
        })
      })

      context("with max periods duration", async () => {
        let tx: ContractTransaction

        beforeEach(async () => {
          await mineBlocksTo(startBlock + dkgTimeout - 1)

          const { dkgResult } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved
          )

          const submitter = await getDkgResultSubmitterSigner(
            randomBeacon,
            dkgResult
          )

          await mineBlocks(params.dkgResultChallengePeriodLength)

          tx = await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
        })

        // Just an explicit assertion to make sure transaction passes correctly
        // for max periods duration.
        it("should succeed", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "GroupActivated")
            .withArgs(0, groupPublicKey)
        })

        it("should unlock the sortition pool", async () => {
          expect(await sortitionPool.isLocked()).to.be.false
        })
      })

      context("with misbehaved operators", async () => {
        const misbehavedIndices = [2, 10, 64]
        let misbehavedIds
        let tx: ContractTransaction

        beforeEach(async () => {
          await mineBlocksTo(startBlock + dkgTimeout - 1)

          const { dkgResult, members } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            misbehavedIndices
          )

          misbehavedIds = misbehavedIndices.map((i) => members[i - 1])

          const submitter = await getDkgResultSubmitterSigner(
            randomBeacon,
            dkgResult
          )

          await mineBlocks(params.dkgResultChallengePeriodLength)
          tx = await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
        })

        it("should ban misbehaved operators from sortition pool rewards", async () => {
          const now = await helpers.time.lastBlockTime()
          const expectedUntil = now + 1209600 // 2 weeks

          await expect(tx)
            .to.emit(sortitionPool, "IneligibleForRewards")
            .withArgs(misbehavedIds, expectedUntil)
        })

        it("should clean dkg data", async () => {
          await assertDkgResultCleanData(randomBeacon)
        })
      })
    })

    context(
      "when the balance of DKG rewards pool is smaller than the DKG submission reward",
      async () => {
        let dkgRewardsPoolBalance: BigNumber
        let tx: ContractTransaction
        let initApproverBalance: BigNumber
        let submitter: SignerWithAddress

        beforeEach(async () => {
          dkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()

          // Set the DKG result submission reward to twice the amount of test
          // tokens in the DKG rewards pool
          await randomBeaconGovernance.beginDkgResultSubmissionRewardUpdate(
            dkgRewardsPoolBalance.mul(2)
          )
          await helpers.time.increaseTime(12 * 60 * 60)
          await randomBeaconGovernance.finalizeDkgResultSubmissionRewardUpdate()

          const [genesisTx, genesisSeed] = await genesis(randomBeacon)
          const startBlock: number = genesisTx.blockNumber
          await mineBlocksTo(startBlock + dkgTimeout - 1)

          const { dkgResult } = await signAndSubmitCorrectDkgResult(
            randomBeacon,
            groupPublicKey,
            genesisSeed,
            startBlock,
            noMisbehaved
          )

          submitter = await getDkgResultSubmitterSigner(randomBeacon, dkgResult)

          initApproverBalance = await testToken.balanceOf(
            await submitter.getAddress()
          )

          await mineBlocks(params.dkgResultChallengePeriodLength)
          tx = await randomBeacon.connect(submitter).approveDkgResult(dkgResult)
        })

        it("should succeed", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "GroupActivated")
            .withArgs(0, groupPublicKey)
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
      it("should revert with 'dkg has not timed out' error", async () => {
        await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
          "dkg has not timed out"
        )
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number

      beforeEach("run genesis", async () => {
        const [genesisTx] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
      })

      context("with dkg not timed out", async () => {
        context("with off-chain dkg time not passed", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime - 1)
          })

          it("should revert with 'dkg has not timed out' error", async () => {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "dkg has not timed out"
            )
          })
        })

        context("with off-chain dkg time passed", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + constants.offchainDkgTime)
          })

          it("should revert with 'dkg has not timed out' error", async () => {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "dkg has not timed out"
            )
          })
        })

        context("with result submission period almost ended", async () => {
          beforeEach(async () => {
            await mineBlocksTo(startBlock + dkgTimeout - 1)
          })

          it("should revert with 'dkg has not timed out' error", async () => {
            await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
              "dkg has not timed out"
            )
          })
        })
      })

      context("with dkg timed out", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + dkgTimeout)
        })

        context("called by a third party", async () => {
          let tx: ContractTransaction
          let initialDkgRewardsPoolBalance: BigNumber
          let initialNotifierBalance: BigNumber

          beforeEach(async () => {
            initialDkgRewardsPoolBalance = await randomBeacon.dkgRewardsPool()

            initialNotifierBalance = await testToken.balanceOf(
              await thirdParty.getAddress()
            )
            tx = await randomBeacon.connect(thirdParty).notifyDkgTimeout()
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
            ).to.be.equal(sortitionPoolUnlockingReward)

            const currentNotifierBalance: BigNumber = await testToken.balanceOf(
              await thirdParty.getAddress()
            )
            expect(
              currentNotifierBalance.sub(initialNotifierBalance)
            ).to.be.equal(sortitionPoolUnlockingReward)
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
    const stubDkgResult: DkgResult = {
      groupPubKey: blsData.groupPubKey,
      members: [1, 2, 3, 4],
      misbehavedMembersIndices: [],
      signatures: "0x01020304",
      signingMembersIndices: [1, 2, 3, 4],
      submitterMemberIndex: 1,
    }

    context("with initial contract state", async () => {
      it("should revert with 'current state is not CHALLENGE' error", async () => {
        await expect(
          randomBeacon.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let genesisSeed

      beforeEach("run genesis", async () => {
        const [genesisTx, seed] = await genesis(randomBeacon)

        startBlock = genesisTx.blockNumber
        genesisSeed = seed
      })

      it("should revert with 'current state is not CHALLENGE' error", async () => {
        await expect(
          randomBeacon.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("current state is not CHALLENGE")
      })

      context("with off-chain dkg time passed", async () => {
        beforeEach(async () => {
          await mineBlocksTo(startBlock + constants.offchainDkgTime)
        })

        context("with dkg result not submitted", async () => {
          it("should revert with 'current state is not CHALLENGE' error", async () => {
            await expect(
              randomBeacon.challengeDkgResult(stubDkgResult)
            ).to.be.revertedWith("current state is not CHALLENGE")
          })
        })

        context("with malicious dkg result submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DkgResult
          let submitter: SignerWithAddress

          beforeEach(async () => {
            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
            } = await signAndSubmitArbitraryDkgResult(
              randomBeacon,
              groupPublicKey,
              // Mix signers to make the result malicious.
              mixSigners(await selectGroup(sortitionPool, genesisSeed)),
              startBlock,
              noMisbehaved
            ))

            submitter = await getDkgResultSubmitterSigner(
              randomBeacon,
              dkgResult
            )

            resultSubmissionBlock = tx.blockNumber
          })

          context("at the beginning of challenge period", async () => {
            context("called by a third party", async () => {
              let tx: ContractTransaction
              beforeEach(async () => {
                tx = await randomBeacon
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
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

              it("should remove a candidate group", async () => {
                const groupsRegistry = await randomBeacon.getGroupsRegistry()

                expect(groupsRegistry).to.be.lengthOf(0)
              })

              it("should emit CandidateGroupRemoved event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "CandidateGroupRemoved")
                  .withArgs(groupPublicKey)
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
            beforeEach(async () => {
              await mineBlocksTo(
                resultSubmissionBlock +
                  params.dkgResultChallengePeriodLength -
                  1
              )
            })

            context("called by a third party", async () => {
              let tx: ContractTransaction
              beforeEach(async () => {
                tx = await randomBeacon
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
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

              it("should remove a candidate group", async () => {
                const groupsRegistry = await randomBeacon.getGroupsRegistry()

                expect(groupsRegistry).to.be.lengthOf(0)
              })

              it("should emit CandidateGroupRemoved event", async () => {
                await expect(tx)
                  .to.emit(randomBeacon, "CandidateGroupRemoved")
                  .withArgs(groupPublicKey)
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
            beforeEach(async () => {
              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )
            })

            it("should revert with 'challenge period has already passed' error", async () => {
              await expect(
                randomBeacon.challengeDkgResult(dkgResult)
              ).to.be.revertedWith("challenge period has already passed")
            })
          })
        })

        context("with correct dkg result submitted", async () => {
          let dkgResult: DkgResult

          beforeEach(async () => {
            ;({ dkgResult } = await signAndSubmitCorrectDkgResult(
              randomBeacon,
              groupPublicKey,
              genesisSeed,
              startBlock,
              noMisbehaved
            ))
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
      let dkgResult: DkgResult

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
      let blocksToMine =
        (constants.groupSize * params.dkgResultSubmissionEligibilityDelay) / 2
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock,
        noMisbehaved,
        constants.groupSize / 2
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
      blocksToMine =
        constants.groupSize * params.dkgResultSubmissionEligibilityDelay - 1
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock,
        noMisbehaved,
        constants.groupSize
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
      ).to.be.revertedWith("dkg has not timed out")

      await randomBeacon.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await randomBeacon.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 4 after the submission period
      blocksToMine =
        constants.groupSize * params.dkgResultSubmissionEligibilityDelay
      await mineBlocks(blocksToMine)
      await expect(
        signAndSubmitArbitraryDkgResult(
          randomBeacon,
          groupPublicKey,
          signers,
          startBlock,
          noMisbehaved,
          constants.groupSize
        )
      ).to.be.revertedWith("dkg timeout already passed")

      await randomBeacon.notifyDkgTimeout()
    })
  })

  describe("fundDkgRewardsPool", () => {
    const amount = to1e18(1000)

    let previousDkgRewardsPoolBalance: BigNumber
    let previousRandomBeaconBalance: BigNumber

    beforeEach(async () => {
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
})

async function assertDkgResultCleanData(randomBeacon: RandomBeaconStub) {
  const dkgData = await randomBeacon.getDkgData()

  expect(
    dkgData.parameters.resultChallengePeriodLength,
    "unexpected resultChallengePeriodLength"
  ).to.eq(params.dkgResultChallengePeriodLength)

  expect(
    dkgData.parameters.resultSubmissionEligibilityDelay,
    "unexpected resultSubmissionEligibilityDelay"
  ).to.eq(params.dkgResultSubmissionEligibilityDelay)

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
