import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import { constants, dkgState, params, walletRegistryFixture } from "./fixtures"
import ecdsaData from "./data/ecdsa"
import {
  hashDKGMembers,
  noMisbehaved,
  signAndSubmitArbitraryDkgResult,
  signAndSubmitCorrectDkgResult,
  signAndSubmitUnrecoverableDkgResult,
  expectDkgResultSubmittedEvent,
} from "./utils/dkg"
import { selectGroup, hashUint32Array } from "./utils/groups"
import { firstEligibleIndex, shiftEligibleIndex } from "./utils/submission"
import { createNewWallet, getWalletID, requestNewWallet } from "./utils/wallets"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { BigNumber, ContractTransaction, Signer } from "ethers"
import type {
  SortitionPool,
  WalletRegistry,
  WalletRegistryStub,
  StakingStub,
} from "../typechain"
import type { DkgResult } from "./utils/dkg"
import type { Operator } from "./utils/operators"

const { to1e18 } = helpers.number
const { mineBlocks, mineBlocksTo } = helpers.time
const { createSnapshot, restoreSnapshot } = helpers.snapshot

const { keccak256 } = ethers.utils

describe("WalletRegistry - Wallet Creation", async () => {
  const dkgTimeout: number =
    constants.offchainDkgTime +
    constants.groupSize * params.dkgResultSubmissionEligibilityDelay
  const groupPublicKey: string = ethers.utils.hexValue(
    ecdsaData.group1.publicKey
  )
  const groupPublicKeyHash: string = ethers.utils.keccak256(groupPublicKey)
  const firstEligibleSubmitterIndex: number = firstEligibleIndex(
    keccak256(ecdsaData.group1.publicKey)
  )

  const stubDkgResult: DkgResult = {
    submitterMemberIndex: 1,
    groupPubKey: ecdsaData.group1.publicKey,
    signingMembersIndices: [1, 2, 3, 4],
    signatures: "0x01020304",
    misbehavedMembersIndices: [],
    members: [1, 2, 3, 4],
    membersHash: hashDKGMembers([1, 2, 3, 4]),
  }

  let walletRegistry: WalletRegistryStub & WalletRegistry
  let sortitionPool: SortitionPool
  let staking: StakingStub

  let deployer: SignerWithAddress
  let walletOwner: SignerWithAddress
  let thirdParty: SignerWithAddress

  let operators: Operator[]

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({
      walletRegistry,
      sortitionPool,
      walletOwner,
      deployer,
      thirdParty,
      operators,
      staking,
    } = await waffle.loadFixture(walletRegistryFixture))
  })

  describe("requestNewWallet", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).requestNewWallet()
        ).to.be.revertedWith("Caller is not the Wallet Owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).requestNewWallet()
        ).to.be.revertedWith("Caller is not the Wallet Owner")
      })
    })

    context("when called by the wallet owner", async () => {
      context("with initial contract state", async () => {
        let tx: ContractTransaction
        let dkgSeed: BigNumber

        before("start wallet creation", async () => {
          await createSnapshot()
          ;({ tx, dkgSeed } = await requestNewWallet(
            walletRegistry,
            walletOwner
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should lock the sortition pool", async () => {
          await expect(await sortitionPool.isLocked()).to.be.true
        })

        it("should emit DkgStateLocked event", async () => {
          await expect(tx).to.emit(walletRegistry, "DkgStateLocked")
        })

        it("should emit DkgStarted event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "DkgStarted")
            .withArgs(dkgSeed)
        })

        it("should not register new wallet", async () => {
          await expect(tx).not.to.emit(walletRegistry, "WalletCreated")
        })
      })

      context("with wallet creation in progress", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before("start wallet creation", async () => {
          await createSnapshot()
          ;({ startBlock, dkgSeed } = await requestNewWallet(
            walletRegistry,
            walletOwner
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with dkg result not submitted", async () => {
          it("should revert", async () => {
            await expect(
              requestNewWallet(walletRegistry, walletOwner)
            ).to.be.revertedWith("Current state is not IDLE")
          })
        })

        context("with valid dkg result submitted", async () => {
          let dkgResult: DkgResult
          let submitter: SignerWithAddress

          before("submit dkg result", async () => {
            await createSnapshot()

            await mineBlocks(constants.offchainDkgTime)
            ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
              startBlock,
              noMisbehaved
            ))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with dkg result not approved", async () => {
            it("should revert with 'current state is not IDLE' error", async () => {
              await expect(
                requestNewWallet(walletRegistry, walletOwner)
              ).to.be.revertedWith("Current state is not IDLE")
            })
          })

          context("with dkg result approved", async () => {
            before("approve dkg result", async () => {
              await createSnapshot()

              await mineBlocks(params.dkgResultChallengePeriodLength)

              await walletRegistry
                .connect(submitter)
                .approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should succeed", async () => {
              await expect(requestNewWallet(walletRegistry, walletOwner)).to.not
                .be.reverted
            })
          })
        })

        context("with invalid dkg result submitted", async () => {
          context("with dkg result challenged", async () => {
            let dkgResult: DkgResult

            before("submit and challenge dkg result", async () => {
              await createSnapshot()
              await mineBlocks(constants.offchainDkgTime)
              ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
                walletRegistry,
                groupPublicKey,
                // Mix operators to make the result malicious
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
                startBlock,
                noMisbehaved
              ))

              await walletRegistry.challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert", async () => {
              await expect(
                requestNewWallet(walletRegistry, walletOwner)
              ).to.be.revertedWith("Current state is not IDLE")
            })
          })
        })

        context("with dkg timeout notified", async () => {
          before("notify dkg timeout", async () => {
            await createSnapshot()

            await mineBlocks(dkgTimeout)

            await walletRegistry.notifyDkgTimeout()
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should succeed", async () => {
            await expect(requestNewWallet(walletRegistry, walletOwner)).not.to
              .be.reverted
          })
        })
      })
    })
  })

  describe("getWalletCreationState", async () => {
    context("with initial contract state", async () => {
      it("should return IDLE state", async () => {
        expect(await walletRegistry.getWalletCreationState()).to.be.equal(
          dkgState.IDLE
        )
      })
    })

    context("when wallet creation started", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletRegistry,
          walletOwner
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("at the start of off-chain dkg period", async () => {
        it("should return KEY_GENERATION state", async () => {
          expect(await walletRegistry.getWalletCreationState()).to.be.equal(
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
          expect(await walletRegistry.getWalletCreationState()).to.be.equal(
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
            expect(await walletRegistry.getWalletCreationState()).to.be.equal(
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
              expect(await walletRegistry.getWalletCreationState()).to.be.equal(
                dkgState.AWAITING_RESULT
              )
            })
          })
        })

        context("when dkg result was submitted", async () => {
          let dkgResult: DkgResult
          let submitter: SignerWithAddress

          before("submit dkg result", async () => {
            await createSnapshot()
            ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
              startBlock,
              noMisbehaved
            ))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when dkg result was not approved", async () => {
            it("should return CHALLENGE state", async () => {
              expect(await walletRegistry.getWalletCreationState()).to.be.equal(
                dkgState.CHALLENGE
              )
            })
          })

          context("when dkg result was approved", async () => {
            before("approve dkg result", async () => {
              await createSnapshot()

              await mineBlocks(params.dkgResultChallengePeriodLength)

              await walletRegistry
                .connect(submitter)
                .approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return IDLE state", async () => {
              expect(await walletRegistry.getWalletCreationState()).to.be.equal(
                dkgState.IDLE
              )
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DkgResult

          before("submit malicious dkg result", async () => {
            await createSnapshot()
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              walletRegistry,
              groupPublicKey,
              // Mix operators to make the result malicious
              mixOperators(await selectGroup(sortitionPool, dkgSeed)),
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
              await walletRegistry.challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return AWAITING_RESULT state", async () => {
              expect(await walletRegistry.getWalletCreationState()).to.be.equal(
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
        await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
      })
    })

    context("when wallet creation started", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletRegistry,
          walletOwner
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("within off-chain dkg period", async () => {
        it("should return false", async () => {
          await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
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
            await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
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
              await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
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
              await expect(await walletRegistry.hasDkgTimedOut()).to.be.true
            })
          })
        })

        context("when dkg result was submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResult: DkgResult
          let submitter: SignerWithAddress

          before(async () => {
            await createSnapshot()

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              submitter,
            } = await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
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
                await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
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
                await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
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
                await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
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
                await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
              })
            })
          })

          context("when dkg result was approved", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )

              await walletRegistry
                .connect(submitter)
                .approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return false", async () => {
              await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DkgResult

          before(async () => {
            await createSnapshot()
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              walletRegistry,
              groupPublicKey,
              // Mix operators to make the result malicious.
              mixOperators(await selectGroup(sortitionPool, dkgSeed)),
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

              const tx = await walletRegistry.challengeDkgResult(dkgResult)
              challengeBlockNumber = tx.blockNumber
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("at the end of dkg result submission period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  challengeBlockNumber +
                    constants.groupSize *
                      params.dkgResultSubmissionEligibilityDelay
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return false", async () => {
                await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
              })
            })

            context("after dkg result submission period", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(
                  challengeBlockNumber +
                    constants.groupSize *
                      params.dkgResultSubmissionEligibilityDelay +
                    1
                )
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return true", async () => {
                await expect(await walletRegistry.hasDkgTimedOut()).to.be.true
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
            walletRegistry,
            groupPublicKey,
            operators,
            1,
            noMisbehaved
          )
        ).to.be.revertedWith("Current state is not AWAITING_RESULT")
      })
    })

    context("with wallet creation in progress", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start new wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletRegistry,
          walletOwner
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with wallet creation not timed out", async () => {
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
                walletRegistry,
                groupPublicKey,
                dkgSeed,
                startBlock,
                noMisbehaved
              )
            ).to.be.revertedWith("Current state is not AWAITING_RESULT")
          })
        })

        context("with off-chain dkg time passed", async () => {
          context("with dkg result not submitted", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(startBlock + constants.offchainDkgTime)
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("with enough signatures on the result", async () => {
              let tx: ContractTransaction
              let dkgResult: DkgResult
              let dkgResultHash: string

              before(async () => {
                await createSnapshot()
                ;({
                  transaction: tx,
                  dkgResult,
                  dkgResultHash,
                } = await signAndSubmitArbitraryDkgResult(
                  walletRegistry,
                  groupPublicKey,
                  operators,
                  startBlock,
                  noMisbehaved,
                  firstEligibleSubmitterIndex,
                  constants.groupThreshold
                ))
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultSubmitted event", async () => {
                await expectDkgResultSubmittedEvent(tx, {
                  resultHash: dkgResultHash,
                  seed: dkgSeed,
                  result: dkgResult,
                })
              })

              it("should not register a new wallet", async () => {
                await expect(tx).not.to.emit(walletRegistry, "WalletCreated")
              })

              it("should not unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.true
              })
            })

            context("with not enough signatures on the result", async () => {
              it("should succeed", async () => {
                await createSnapshot()

                await expect(
                  signAndSubmitArbitraryDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    operators,
                    startBlock,
                    noMisbehaved,
                    firstEligibleSubmitterIndex,
                    constants.groupThreshold - 1
                  )
                ).to.not.be.reverted

                await restoreSnapshot()
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

              context("at the beginning of submission period", async () => {
                it("should succeed for the first submitter", async () => {
                  await assertSubmissionSucceeds(firstSubmitterIndex)
                })

                it("should revert for the second submitter", async () => {
                  await assertSubmissionReverts(secondSubmitterIndex)
                })

                it("should revert for the last submitter", async () => {
                  await assertSubmissionReverts(lastSubmitterIndex)
                })
              })

              context(
                "with first submitter eligibility delay period almost ended",
                async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      submissionStartBlockNumber +
                        params.dkgResultSubmissionEligibilityDelay -
                        2
                    )
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should succeed for the first submitter", async () => {
                    await assertSubmissionSucceeds(firstSubmitterIndex)
                  })

                  it("should revert for the second submitter", async () => {
                    await assertSubmissionReverts(secondSubmitterIndex)
                  })

                  it("should revert for the last submitter", async () => {
                    await assertSubmissionReverts(lastSubmitterIndex)
                  })
                }
              )

              context(
                "with first submitter eligibility delay period ended",
                async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      submissionStartBlockNumber +
                        params.dkgResultSubmissionEligibilityDelay -
                        1
                    )
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should succeed for the first submitter", async () => {
                    await assertSubmissionSucceeds(firstSubmitterIndex)
                  })

                  it("should succeed for the second submitter", async () => {
                    await assertSubmissionSucceeds(secondSubmitterIndex)
                  })

                  it("should revert for the third submitter", async () => {
                    await assertSubmissionReverts(thirdSubmitterIndex)
                  })

                  it("should revert for the last submitter", async () => {
                    await assertSubmissionReverts(lastSubmitterIndex)
                  })
                }
              )

              context(
                "with the last submitter eligibility delay period almost ended",
                async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      submissionStartBlockNumber +
                        constants.groupSize *
                          params.dkgResultSubmissionEligibilityDelay -
                        1
                    )
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should succeed for the first submitter", async () => {
                    await assertSubmissionSucceeds(firstSubmitterIndex)
                  })

                  it("should succeed for the last submitter", async () => {
                    await assertSubmissionSucceeds(lastSubmitterIndex)
                  })
                }
              )

              context(
                "with the last submitter eligibility delay period ended",
                async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      submissionStartBlockNumber +
                        constants.groupSize *
                          params.dkgResultSubmissionEligibilityDelay
                    )
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should revert for the first submitter", async () => {
                    await assertSubmissionReverts(
                      firstSubmitterIndex,
                      "DKG timeout already passed"
                    )
                  })

                  it("should revert for the last submitter", async () => {
                    await assertSubmissionReverts(
                      lastSubmitterIndex,
                      "DKG timeout already passed"
                    )
                  })
                }
              )
            })
          })

          context("with dkg result submitted", async () => {
            let dkgResult: DkgResult
            let submitter: SignerWithAddress
            let resultSubmissionBlock: number

            before(async () => {
              await createSnapshot()

              await mineBlocksTo(startBlock + constants.offchainDkgTime)

              let tx: ContractTransaction
              ;({
                transaction: tx,
                dkgResult,
                submitter,
              } = await signAndSubmitCorrectDkgResult(
                walletRegistry,
                groupPublicKey,
                dkgSeed,
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
                  walletRegistry,
                  groupPublicKey,
                  dkgSeed,
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

                await walletRegistry
                  .connect(submitter)
                  .approveDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should revert", async () => {
                await expect(
                  signAndSubmitCorrectDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    dkgSeed,
                    startBlock,
                    noMisbehaved
                  )
                ).to.be.revertedWith("Sortition pool unlocked")
              })
            })
          })

          context("with malicious dkg result submitted", async () => {
            let dkgResult: DkgResult

            before(async () => {
              await createSnapshot()

              await mineBlocksTo(startBlock + constants.offchainDkgTime)
              ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
                walletRegistry,
                groupPublicKey,
                // Mix operators to make the result malicious.
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
                startBlock,
                noMisbehaved
              ))
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("with dkg result challenged", async () => {
              let challengeBlockNumber: number

              before(async () => {
                await createSnapshot()

                const tx = await walletRegistry.challengeDkgResult(dkgResult)
                challengeBlockNumber = tx.blockNumber
              })

              after(async () => {
                await restoreSnapshot()
              })

              context("with a fresh dkg result", async () => {
                it("should emit DkgResultSubmitted event", async () => {
                  await createSnapshot()

                  const {
                    transaction: tx,
                    dkgResultHash: newDkgResultHash,
                    dkgResult: newDkgResult,
                  } = await signAndSubmitCorrectDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    dkgSeed,
                    startBlock,
                    noMisbehaved
                  )

                  await expectDkgResultSubmittedEvent(tx, {
                    resultHash: newDkgResultHash,
                    seed: dkgSeed,
                    result: newDkgResult,
                  })

                  await restoreSnapshot()
                })

                describe("submission eligibility verification", async () => {
                  let submissionStartBlockNumber: number

                  before(async () => {
                    await createSnapshot()

                    submissionStartBlockNumber = challengeBlockNumber
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  context("at the beginning of submission period", async () => {
                    it("should succeed for the first submitter", async () => {
                      await assertSubmissionSucceeds(firstSubmitterIndex)
                    })

                    it("should revert for the second submitter", async () => {
                      await assertSubmissionReverts(secondSubmitterIndex)
                    })

                    it("should revert for the last submitter", async () => {
                      await assertSubmissionReverts(lastSubmitterIndex)
                    })
                  })

                  context(
                    "with first submitter eligibility delay period almost ended",
                    async () => {
                      before(async () => {
                        await createSnapshot()

                        await mineBlocksTo(
                          submissionStartBlockNumber +
                            params.dkgResultSubmissionEligibilityDelay -
                            2
                        )
                      })

                      after(async () => {
                        await restoreSnapshot()
                      })

                      it("should succeed for the first submitter", async () => {
                        await assertSubmissionSucceeds(firstSubmitterIndex)
                      })

                      it("should revert for the second submitter", async () => {
                        await assertSubmissionReverts(secondSubmitterIndex)
                      })

                      it("should revert for the last submitter", async () => {
                        await assertSubmissionReverts(lastSubmitterIndex)
                      })
                    }
                  )

                  context(
                    "with first submitter eligibility delay period ended",
                    async () => {
                      before(async () => {
                        await createSnapshot()

                        await mineBlocksTo(
                          submissionStartBlockNumber +
                            params.dkgResultSubmissionEligibilityDelay -
                            1
                        )
                      })

                      after(async () => {
                        await restoreSnapshot()
                      })

                      it("should succeed for the first submitter", async () => {
                        await assertSubmissionSucceeds(firstSubmitterIndex)
                      })

                      it("should succeed for the second submitter", async () => {
                        await assertSubmissionSucceeds(secondSubmitterIndex)
                      })

                      it("should revert for the third submitter", async () => {
                        await assertSubmissionReverts(thirdSubmitterIndex)
                      })

                      it("should revert for the last submitter", async () => {
                        await assertSubmissionReverts(lastSubmitterIndex)
                      })
                    }
                  )

                  context(
                    "with the last submitter eligibility delay period almost ended",
                    async () => {
                      before(async () => {
                        await createSnapshot()

                        await mineBlocksTo(
                          submissionStartBlockNumber +
                            constants.groupSize *
                              params.dkgResultSubmissionEligibilityDelay -
                            1
                        )
                      })

                      after(async () => {
                        await restoreSnapshot()
                      })

                      it("should succeed for the first submitter", async () => {
                        await assertSubmissionSucceeds(firstSubmitterIndex)
                      })

                      it("should succeed for the last submitter", async () => {
                        await assertSubmissionSucceeds(lastSubmitterIndex)
                      })
                    }
                  )

                  context(
                    "with the last submitter eligibility delay period ended",
                    async () => {
                      before(async () => {
                        await createSnapshot()

                        await mineBlocksTo(
                          submissionStartBlockNumber +
                            constants.groupSize *
                              params.dkgResultSubmissionEligibilityDelay
                        )
                      })

                      after(async () => {
                        await restoreSnapshot()
                      })

                      it("should revert for the first submitter", async () => {
                        await assertSubmissionReverts(
                          firstSubmitterIndex,
                          "DKG timeout already passed"
                        )
                      })

                      it("should revert for the last submitter", async () => {
                        await assertSubmissionReverts(
                          lastSubmitterIndex,
                          "DKG timeout already passed"
                        )
                      })
                    }
                  )
                })
              })
            })
          })

          context("with misbehaved members", async () => {
            let tx: ContractTransaction
            let dkgResult: DkgResult
            let dkgResultHash: string

            context(
              "when misbehaved members are in ascending order",
              async () => {
                const misbehavedIndices = [2, 9, 11, 30, 60, 64]

                before(async () => {
                  await createSnapshot()

                  await mineBlocksTo(startBlock + constants.offchainDkgTime)
                  ;({
                    transaction: tx,
                    dkgResult,
                    dkgResultHash,
                  } = await signAndSubmitCorrectDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    dkgSeed,
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
                    seed: dkgSeed,
                    result: dkgResult,
                  })
                })
              }
            )
          })
        })
      })

      context("with wallet creation timed out", async () => {
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
                walletRegistry,
                groupPublicKey,
                dkgSeed,
                startBlock,
                noMisbehaved
              )
            ).to.be.revertedWith("DKG timeout already passed")
          })
        })
      })

      // Submission Eligibility Test Helpers
      const firstSubmitterIndex = firstEligibleSubmitterIndex
      const secondSubmitterIndex = shiftEligibleIndex(firstSubmitterIndex, 1)
      const thirdSubmitterIndex = shiftEligibleIndex(firstSubmitterIndex, 2)
      const lastSubmitterIndex = shiftEligibleIndex(
        firstSubmitterIndex,
        constants.groupSize - 1
      )

      async function assertSubmissionSucceeds(
        submitterIndex: number
      ): Promise<void> {
        await createSnapshot()

        const {
          transaction: tx,
          dkgResult,
          dkgResultHash,
        } = await signAndSubmitCorrectDkgResult(
          walletRegistry,
          groupPublicKey,
          dkgSeed,
          startBlock,
          noMisbehaved,
          submitterIndex
        )

        await expectDkgResultSubmittedEvent(tx, {
          resultHash: dkgResultHash,
          seed: dkgSeed,
          result: dkgResult,
        })

        await restoreSnapshot()
      }

      async function assertSubmissionReverts(
        submitterIndex: number,
        message = "Submitter is not eligible"
      ): Promise<void> {
        await expect(
          signAndSubmitCorrectDkgResult(
            walletRegistry,
            groupPublicKey,
            dkgSeed,
            startBlock,
            noMisbehaved,
            submitterIndex
          )
        ).to.be.revertedWith(message)
      }
    })
  })

  describe("approveDkgResult", async () => {
    context("with initial contract state", async () => {
      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          walletRegistry.approveDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start new wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletRegistry,
          walletOwner
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          walletRegistry.approveDkgResult(stubDkgResult)
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
              walletRegistry.approveDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })
        })

        context("with dkg result submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DkgResult
          let submitter: SignerWithAddress
          let walletID: string

          const submitterIndex = firstEligibleSubmitterIndex

          before(async () => {
            await createSnapshot()

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
              submitter,
            } = await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
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
                walletRegistry.connect(submitter).approveDkgResult(dkgResult)
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
              // let initialDkgRewardsPoolBalance: BigNumber
              // let initialSubmitterBalance: BigNumber

              before(async () => {
                await createSnapshot()

                // initialDkgRewardsPoolBalance =
                //   await walletRegistry.dkgRewardsPool()
                // initialSubmitterBalance = await testToken.balanceOf(
                //   await submitter.getAddress()
                // )
                tx = await walletRegistry
                  .connect(submitter)
                  .approveDkgResult(dkgResult)

                walletID = await getWalletID(tx)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultApproved event", async () => {
                await expect(tx)
                  .to.emit(walletRegistry, "DkgResultApproved")
                  .withArgs(dkgResultHash, await submitter.getAddress())
              })

              it("should clean dkg data", async () => {
                await assertDkgResultCleanData(walletRegistry)
              })

              it("should register a new wallet", async () => {
                await expect(walletID).to.be.equal(groupPublicKeyHash)

                const wallet = await walletRegistry.getWallet(walletID)

                await expect(wallet.membersIdsHash).to.be.equal(
                  hashUint32Array(dkgResult.members)
                )
              })

              // it("should reward the submitter with tokens from DKG rewards pool", async () => {
              //   const currentDkgRewardsPoolBalance =
              //     await walletRegistry.dkgRewardsPool()
              //   expect(
              //     initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
              //   ).to.be.equal(params.dkgResultSubmissionReward)

              //   const currentSubmitterBalance: BigNumber =
              //     await testToken.balanceOf(await submitter.getAddress())
              //   expect(
              //     currentSubmitterBalance.sub(initialSubmitterBalance)
              //   ).to.be.equal(params.dkgResultSubmissionReward)
              // })

              it("should emit WalletCreated event", async () => {
                await expect(tx)
                  .to.emit(walletRegistry, "WalletCreated")
                  .withArgs(walletID, dkgResultHash)
              })

              it("should unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.false
              })
            })

            context("when called by a third party", async () => {
              context("when the third party is not yet eligible", async () => {
                before(async () => {
                  await createSnapshot()

                  await mineBlocks(
                    params.dkgResultSubmissionEligibilityDelay - 1
                  )
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should revert", async () => {
                  await expect(
                    walletRegistry
                      .connect(thirdParty)
                      .approveDkgResult(dkgResult)
                  ).to.be.revertedWith(
                    "Only the DKG result submitter can approve the result at this moment"
                  )
                })
              })

              context("when the third party is eligible", async () => {
                let tx: ContractTransaction
                // let initialDkgRewardsPoolBalance: BigNumber
                // let initApproverBalance: BigNumber

                before(async () => {
                  await createSnapshot()

                  await mineBlocks(params.dkgResultSubmissionEligibilityDelay)
                  // initialDkgRewardsPoolBalance =
                  //   await walletRegistry.dkgRewardsPool()
                  // initApproverBalance = await testToken.balanceOf(
                  //   await thirdParty.getAddress()
                  // )
                  tx = await walletRegistry
                    .connect(thirdParty)
                    .approveDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should emit DkgResultApproved event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgResultApproved")
                    .withArgs(dkgResultHash, thirdParty.address)
                })

                // it("should pay the reward to the third party", async () => {
                //   const currentDkgRewardsPoolBalance =
                //     await walletRegistry.dkgRewardsPool()
                //   expect(
                //     initialDkgRewardsPoolBalance.sub(
                //       currentDkgRewardsPoolBalance
                //     )
                //   ).to.be.equal(params.dkgResultSubmissionReward)

                //   const currentApproverBalance = await testToken.balanceOf(
                //     await thirdParty.getAddress()
                //   )
                //   expect(
                //     currentApproverBalance.sub(initApproverBalance)
                //   ).to.be.equal(params.dkgResultSubmissionReward)
                // })
              })
            })
          })
        })

        context("when there was a challenged result before", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DkgResult

          // First result is malicious and submitter is also malicious
          const maliciousSubmitter = firstEligibleSubmitterIndex

          // Submit a second result by another submitter
          const submitterIndexShift = 5
          const anotherSubmitterIndex = shiftEligibleIndex(
            maliciousSubmitter,
            submitterIndexShift
          )
          let anotherSubmitter: Signer

          before(async () => {
            await createSnapshot()

            await mineBlocks(
              params.dkgResultSubmissionEligibilityDelay * submitterIndexShift
            )

            const { dkgResult: maliciousDkgResult } =
              await signAndSubmitArbitraryDkgResult(
                walletRegistry,
                groupPublicKey,
                // Mix operators to make the result malicious.
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
                startBlock,
                noMisbehaved,
                maliciousSubmitter
              )

            await walletRegistry.challengeDkgResult(maliciousDkgResult)

            await mineBlocks(
              params.dkgResultSubmissionEligibilityDelay * anotherSubmitterIndex
            )

            let tx: ContractTransaction
            ;({
              transaction: tx,
              dkgResult,
              dkgResultHash,
              submitter: anotherSubmitter,
            } = await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
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
                walletRegistry
                  .connect(anotherSubmitter)
                  .approveDkgResult(dkgResult)
              ).to.be.revertedWith("Challenge period has not passed yet")
            })
          })

          context("with challenge period passed", async () => {
            let tx: ContractTransaction
            let walletID: string
            let initialDkgRewardsPoolBalance: BigNumber
            let initialSubmitterBalance: BigNumber

            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )

              // initialDkgRewardsPoolBalance =
              //   await walletRegistry.dkgRewardsPool()

              // initialSubmitterBalance = await testToken.balanceOf(
              //   await anotherSubmitter.getAddress()
              // )

              tx = await walletRegistry
                .connect(anotherSubmitter)
                .approveDkgResult(dkgResult)

              walletID = await getWalletID(tx)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultApproved event", async () => {
              await expect(tx)
                .to.emit(walletRegistry, "DkgResultApproved")
                .withArgs(dkgResultHash, await anotherSubmitter.getAddress())
            })

            it("should register a new wallet", async () => {
              await expect(walletID).to.be.equal(groupPublicKeyHash)

              const wallet = await walletRegistry.getWallet(walletID)

              await expect(wallet.membersIdsHash).to.be.equal(
                hashUint32Array(dkgResult.members)
              )
            })

            // it("should reward the submitter with tokens from DKG rewards pool", async () => {
            //   const currentDkgRewardsPoolBalance =
            //     await walletRegistry.dkgRewardsPool()
            //   expect(
            //     initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
            //   ).to.be.equal(params.dkgResultSubmissionReward)

            //   const currentSubmitterBalance: BigNumber =
            //     await testToken.balanceOf(await anotherSubmitter.getAddress())
            //   expect(
            //     currentSubmitterBalance.sub(initialSubmitterBalance)
            //   ).to.be.equal(params.dkgResultSubmissionReward)
            // })

            it("should emit WalletCreated event", async () => {
              await expect(tx)
                .to.emit(walletRegistry, "WalletCreated")
                .withArgs(walletID, dkgResultHash)
            })

            it("should unlock the sortition pool", async () => {
              await expect(await sortitionPool.isLocked()).to.be.false
            })
          })
        })
      })

      context("with max periods duration", async () => {
        let tx: ContractTransaction
        let dkgResultHash: string
        let submitter: SignerWithAddress

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout - 1)

          let dkgResult: DkgResult
          ;({ dkgResult, dkgResultHash, submitter } =
            await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
              startBlock,
              noMisbehaved
            ))

          await mineBlocks(params.dkgResultChallengePeriodLength)

          tx = await walletRegistry
            .connect(submitter)
            .approveDkgResult(dkgResult)
        })

        after(async () => {
          await restoreSnapshot()
        })

        // Just an explicit assertion to make sure transaction passes correctly
        // for max periods duration.
        it("should emit DkgResultApproved event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "DkgResultApproved")
            .withArgs(dkgResultHash, submitter.address)
        })

        it("should unlock the sortition pool", async () => {
          await expect(await sortitionPool.isLocked()).to.be.false
        })
      })

      context("with misbehaved operators", async () => {
        const misbehavedIndices = [2, 9, 11, 30, 60, 64]
        let tx: ContractTransaction
        let walletID: string
        let dkgResult: DkgResult

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout - 1)

          let submitter
          ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            walletRegistry,
            groupPublicKey,
            dkgSeed,
            startBlock,
            misbehavedIndices
          ))

          await mineBlocks(params.dkgResultChallengePeriodLength)
          tx = await walletRegistry
            .connect(submitter)
            .approveDkgResult(dkgResult)

          walletID = await getWalletID(tx)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should register a new wallet", async () => {
          await expect(walletID).to.be.equal(groupPublicKeyHash)

          // misbehavedIndices: [2, 9, 11, 30, 60, 64]
          const expectedMembers = [...dkgResult.members]
          expectedMembers.splice(1, 1) // index -1
          expectedMembers.splice(7, 1) // index -2 (cause expectedMembers already shrinked)
          expectedMembers.splice(8, 1) // index -3
          expectedMembers.splice(26, 1) // index -4
          expectedMembers.splice(55, 1) // index -5
          expectedMembers.splice(58, 1) // index -6

          expect(
            (await walletRegistry.getWallet(walletID)).membersIdsHash
          ).to.be.equal(hashUint32Array(expectedMembers))
        })

        // it("should ban misbehaved operators from sortition pool rewards", async () => {
        //   const now = await helpers.time.lastBlockTime()
        //   const expectedUntil = now + params.sortitionPoolRewardsBanDuration

        //   await expect(tx)
        //     .to.emit(sortitionPool, "IneligibleForRewards")
        //     .withArgs(misbehavedIds, expectedUntil)
        // })

        it("should clean dkg data", async () => {
          await assertDkgResultCleanData(walletRegistry)
        })
      })

      // This case shouldn't happen in real life. When a result is submitted
      // with invalid order of misbehaved operators it should be challenged.
      context(
        "when misbehaved operators are not in ascending order",
        async () => {
          const misbehavedIndices = [2, 9, 30, 11, 60, 64]

          let dkgResult: DkgResult
          let submitter: SignerWithAddress

          before(async () => {
            await createSnapshot()

            await mineBlocksTo(startBlock + constants.offchainDkgTime)
            ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
              walletRegistry,
              groupPublicKey,
              dkgSeed,
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
              walletRegistry.connect(submitter).approveDkgResult(dkgResult)
            ).to.not.be.reverted
          })
        }
      )

      context("when misbehaved members contains duplicates", async () => {
        const misbehavedIndices = [2, 9, 9, 10]

        let dkgResult: DkgResult
        let submitter: SignerWithAddress

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + constants.offchainDkgTime)
          ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            walletRegistry,
            groupPublicKey,
            dkgSeed,
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
            walletRegistry.connect(submitter).approveDkgResult(dkgResult)
          ).to.not.be.reverted
        })
      })
    })

    // context(
    //   "when the balance of DKG rewards pool is smaller than the DKG submission reward",
    //   async () => {
    //     let dkgRewardsPoolBalance: BigNumber
    //     let tx: ContractTransaction
    //     let initApproverBalance: BigNumber
    //     let submitter: SignerWithAddress

    //     before(async () => {
    //       await createSnapshot()

    //       dkgRewardsPoolBalance = await walletRegistry.dkgRewardsPool()

    //       // Set the DKG result submission reward to twice the amount of test
    //       // tokens in the DKG rewards pool
    //       await walletRegistryGovernance.beginDkgResultSubmissionRewardUpdate(
    //         dkgRewardsPoolBalance.mul(2)
    //       )
    //       await helpers.time.increaseTime(12 * 60 * 60)
    //       await walletRegistryGovernance.finalizeDkgResultSubmissionRewardUpdate()

    //       const [genesisTx, genesisSeed] = await genesis(walletRegistry)
    //       const startBlock: number = genesisTx.blockNumber
    //       await mineBlocksTo(startBlock + dkgTimeout - 1)

    //       let dkgResult: DkgResult
    //       ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
    //         walletRegistry,
    //         groupPublicKey,
    //         genesisSeed,
    //         startBlock,
    //         noMisbehaved
    //       ))

    //       initApproverBalance = await testToken.balanceOf(
    //         await submitter.getAddress()
    //       )

    //       await mineBlocks(params.dkgResultChallengePeriodLength)
    //       tx = await walletRegistry
    //         .connect(submitter)
    //         .approveDkgResult(dkgResult)
    //     })

    //     after(async () => {
    //       await restoreSnapshot()
    //     })

    //     it("should succeed", async () => {
    //       await expect(tx)
    //         .to.emit(walletRegistry, "GroupActivated")
    //         .withArgs(0, groupPublicKey)
    //     })

    //     it("should pay the approver the whole DKG rewards pool balance", async () => {
    //       expect(await walletRegistry.dkgRewardsPool()).to.be.equal(0)

    //       const currentApproverBalance = await testToken.balanceOf(
    //         await submitter.getAddress()
    //       )
    //       expect(currentApproverBalance.sub(initApproverBalance)).to.be.equal(
    //         dkgRewardsPoolBalance
    //       )
    //     })
    //   }
    // )

    context("with wallet registered", async () => {
      const existingWalletPublicKey: string = ecdsaData.group1.publicKey
      let existingWalletID: string

      before("create a wallet", async () => {
        await createSnapshot()
        ;({ walletID: existingWalletID } = await createNewWallet(
          walletRegistry,
          walletOwner,
          existingWalletPublicKey
        ))

        await expect(existingWalletID).to.be.equal(
          keccak256(existingWalletPublicKey)
        )
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with new wallet requested", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before("start new wallet creation", async () => {
          await createSnapshot()
          ;({ startBlock, dkgSeed } = await requestNewWallet(
            walletRegistry,
            walletOwner
          ))

          await mineBlocks(constants.offchainDkgTime)
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with dkg result not submitted", async () => {
          it("should revert with 'Current state is not CHALLENGE' error", async () => {
            await expect(
              walletRegistry.approveDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })
        })

        context("with dkg result submitted", async () => {
          let dkgResultHash: string
          let dkgResult: DkgResult
          let submitter: SignerWithAddress
          let newWalletID: string

          const newResultPublicKey = ecdsaData.group2.publicKey
          const newResultPublicKeyHash = keccak256(newResultPublicKey)
          const newResultSubmitterIndex = firstEligibleIndex(
            newResultPublicKeyHash
          )

          before("submit dkg result", async () => {
            await createSnapshot()
            ;({ dkgResult, dkgResultHash, submitter } =
              await signAndSubmitCorrectDkgResult(
                walletRegistry,
                newResultPublicKey,
                dkgSeed,
                startBlock,
                noMisbehaved,
                newResultSubmitterIndex
              ))

            await mineBlocks(params.dkgResultChallengePeriodLength)
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("when called by a DKG result submitter", async () => {
            let tx: ContractTransaction

            before(async () => {
              await createSnapshot()

              tx = await walletRegistry
                .connect(submitter)
                .approveDkgResult(dkgResult)

              newWalletID = await getWalletID(tx)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultApproved event", async () => {
              await expect(tx)
                .to.emit(walletRegistry, "DkgResultApproved")
                .withArgs(dkgResultHash, await submitter.getAddress())
            })

            it("should clean dkg data", async () => {
              await assertDkgResultCleanData(walletRegistry)
            })

            it("should register a new wallet", async () => {
              await expect(newWalletID).to.be.equal(newResultPublicKeyHash)

              const wallet = await walletRegistry.getWallet(newWalletID)

              await expect(wallet.membersIdsHash).to.be.equal(
                hashUint32Array(dkgResult.members)
              )
            })

            it("should emit WalletCreated event", async () => {
              await expect(tx)
                .to.emit(walletRegistry, "WalletCreated")
                .withArgs(newWalletID, dkgResultHash)
            })

            it("should unlock the sortition pool", async () => {
              await expect(await sortitionPool.isLocked()).to.be.false
            })
          })
        })
      })
    })
  })

  describe("challengeDkgResult", async () => {
    context("with no wallets registered", async () => {
      it("should revert with 'Current state is not CHALLENGE'", async () => {
        await expect(
          walletRegistry.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })

      context("with group creation in progress", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before("start new wallet creation", async () => {
          await createSnapshot()
          ;({ startBlock, dkgSeed } = await requestNewWallet(
            walletRegistry,
            walletOwner
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert with 'Current state is not CHALLENGE'", async () => {
          await expect(
            walletRegistry.challengeDkgResult(stubDkgResult)
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
            it("should revert with 'Current state is not CHALLENGE'", async () => {
              await expect(
                walletRegistry.challengeDkgResult(stubDkgResult)
              ).to.be.revertedWith("Current state is not CHALLENGE")
            })
          })

          context("with malicious dkg result submitted", async () => {
            let resultSubmissionBlock: number
            let dkgResultHash: string
            let dkgResult: DkgResult
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
                walletRegistry,
                groupPublicKey,
                // Mix operators to make the result malicious.
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
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

                  tx = await walletRegistry
                    .connect(thirdParty)
                    .challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should emit DkgResultChallenged event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgResultChallenged")
                    .withArgs(
                      dkgResultHash,
                      await thirdParty.getAddress(),
                      "Invalid group members"
                    )
                })

                it("should not unlock the sortition pool", async () => {
                  await expect(await sortitionPool.isLocked()).to.be.true
                })

                it("should emit DkgMaliciousResultSlashed event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgMaliciousResultSlashed")
                    .withArgs(dkgResultHash, to1e18(50000), submitter.address)
                })

                it("should slash malicious result submitter", async () => {
                  await expect(tx)
                    .to.emit(staking, "Seized")
                    .withArgs(
                      to1e18(50000),
                      100,
                      await thirdParty.getAddress(),
                      [submitter.address]
                    )
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

                  tx = await walletRegistry
                    .connect(thirdParty)
                    .challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should emit DkgResultChallenged event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgResultChallenged")
                    .withArgs(
                      dkgResultHash,
                      await thirdParty.getAddress(),
                      "Invalid group members"
                    )
                })

                it("should not unlock the sortition pool", async () => {
                  await expect(await sortitionPool.isLocked()).to.be.true
                })

                it("should emit DkgMaliciousResultSlashed event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgMaliciousResultSlashed")
                    .withArgs(dkgResultHash, to1e18(50000), submitter.address)
                })

                it("should slash malicious result submitter", async () => {
                  await expect(tx)
                    .to.emit(staking, "Seized")
                    .withArgs(
                      to1e18(50000),
                      100,
                      await thirdParty.getAddress(),
                      [submitter.address]
                    )
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
                  walletRegistry.challengeDkgResult(dkgResult)
                ).to.be.revertedWith("Challenge period has already passed")
              })
            })

            context(
              "with challenged result not matching the submitted one",
              async () => {
                it("should revert with 'Result under challenge is different than the submitted one'", async () => {
                  const modifiedDkgResult: DkgResult = { ...dkgResult }
                  modifiedDkgResult.submitterMemberIndex =
                    dkgResult.submitterMemberIndex + 1

                  await expect(
                    walletRegistry.challengeDkgResult(modifiedDkgResult)
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
              let dkgResult: DkgResult
              let submitter: SignerWithAddress
              let tx: ContractTransaction

              before(async () => {
                await createSnapshot()
                ;({ dkgResult, dkgResultHash, submitter } =
                  await signAndSubmitUnrecoverableDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    await selectGroup(sortitionPool, dkgSeed),
                    startBlock,
                    noMisbehaved
                  ))

                tx = await walletRegistry
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(walletRegistry, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "validation reverted"
                  )
              })

              it("should not unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.true
              })

              it("should emit DkgMaliciousResultSlashed event", async () => {
                await expect(tx)
                  .to.emit(walletRegistry, "DkgMaliciousResultSlashed")
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

          context(
            "when misbehaved operators are not in ascending order",
            async () => {
              const misbehavedIndices = [2, 9, 30, 11, 60, 64]

              let tx: ContractTransaction
              let dkgResult: DkgResult
              let dkgResultHash: string

              before(async () => {
                await createSnapshot()
                ;({ dkgResult, dkgResultHash } =
                  await signAndSubmitCorrectDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    dkgSeed,
                    startBlock,
                    misbehavedIndices
                  ))

                tx = await walletRegistry
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(walletRegistry, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "Corrupted misbehaved members indices"
                  )
              })

              it("should not unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.true
              })
            }
          )

          context("when misbehaved members contains duplicates", async () => {
            const misbehavedIndices = [2, 9, 30, 30, 60, 64]

            let tx: ContractTransaction
            let dkgResult: DkgResult
            let dkgResultHash: string

            before(async () => {
              await createSnapshot()
              ;({ dkgResult, dkgResultHash } =
                await signAndSubmitCorrectDkgResult(
                  walletRegistry,
                  groupPublicKey,
                  dkgSeed,
                  startBlock,
                  misbehavedIndices
                ))

              tx = await walletRegistry
                .connect(thirdParty)
                .challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultChallenged event", async () => {
              await expect(tx)
                .to.emit(walletRegistry, "DkgResultChallenged")
                .withArgs(
                  dkgResultHash,
                  await thirdParty.getAddress(),
                  "Corrupted misbehaved members indices"
                )
            })

            it("should not unlock the sortition pool", async () => {
              await expect(await sortitionPool.isLocked()).to.be.true
            })
          })

          context("with correct dkg result submitted", async () => {
            let dkgResult: DkgResult

            before(async () => {
              await createSnapshot()
              ;({ dkgResult } = await signAndSubmitCorrectDkgResult(
                walletRegistry,
                groupPublicKey,
                dkgSeed,
                startBlock,
                noMisbehaved
              ))
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'unjustified challenge' error", async () => {
              await expect(
                walletRegistry.challengeDkgResult(dkgResult)
              ).to.be.revertedWith("unjustified challenge")
            })
          })
        })
      })
    })

    context("with wallet registered", async () => {
      before("create a wallet", async () => {
        await createSnapshot()

        await createNewWallet(walletRegistry, walletOwner)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE", async () => {
        await expect(
          walletRegistry.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })

      context("with group creation in progress", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before("start new wallet creation", async () => {
          await createSnapshot()
          ;({ startBlock, dkgSeed } = await requestNewWallet(
            walletRegistry,
            walletOwner
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert with 'Current state is not CHALLENGE'", async () => {
          await expect(
            walletRegistry.challengeDkgResult(stubDkgResult)
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
            it("should revert with 'Current state is not CHALLENGE'", async () => {
              await expect(
                walletRegistry.challengeDkgResult(stubDkgResult)
              ).to.be.revertedWith("Current state is not CHALLENGE")
            })
          })

          context("with malicious dkg result submitted", async () => {
            let resultSubmissionBlock: number
            let dkgResultHash: string
            let dkgResult: DkgResult

            before(async () => {
              await createSnapshot()

              let tx: ContractTransaction
              ;({
                transaction: tx,
                dkgResult,
                dkgResultHash,
              } = await signAndSubmitArbitraryDkgResult(
                walletRegistry,
                groupPublicKey,
                // Mix operators to make the result malicious.
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
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

                  tx = await walletRegistry
                    .connect(thirdParty)
                    .challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should emit DkgResultChallenged event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgResultChallenged")
                    .withArgs(
                      dkgResultHash,
                      await thirdParty.getAddress(),
                      "Invalid group members"
                    )
                })

                it("should not unlock the sortition pool", async () => {
                  await expect(await sortitionPool.isLocked()).to.be.true
                })

                // it("should emit DkgMaliciousResultSlashed event", async () => {
                //   await expect(tx)
                //     .to.emit(walletRegistry, "DkgMaliciousResultSlashed")
                //     .withArgs(dkgResultHash, to1e18(50000), submitter.address)
                // })

                // it("should slash malicious result submitter", async () => {
                //   await expect(tx)
                //     .to.emit(staking, "Seized")
                //     .withArgs(to1e18(50000), 100, await thirdParty.getAddress(), [
                //       submitter.address,
                //     ])
                // })
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

                  tx = await walletRegistry
                    .connect(thirdParty)
                    .challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should emit DkgResultChallenged event", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgResultChallenged")
                    .withArgs(
                      dkgResultHash,
                      await thirdParty.getAddress(),
                      "Invalid group members"
                    )
                })

                it("should not unlock the sortition pool", async () => {
                  await expect(await sortitionPool.isLocked()).to.be.true
                })

                // it("should emit DkgMaliciousResultSlashed event", async () => {
                //   await expect(tx)
                //     .to.emit(walletRegistry, "DkgMaliciousResultSlashed")
                //     .withArgs(dkgResultHash, to1e18(50000), submitter.address)
                // })

                // it("should slash malicious result submitter", async () => {
                //   await expect(tx)
                //     .to.emit(staking, "Seized")
                //     .withArgs(to1e18(50000), 100, await thirdParty.getAddress(), [
                //       submitter.address,
                //     ])
                // })
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
                  walletRegistry.challengeDkgResult(dkgResult)
                ).to.be.revertedWith("Challenge period has already passed")
              })
            })

            context(
              "with challenged result not matching the submitted one",
              async () => {
                it("should revert with 'Result under challenge is different than the submitted one'", async () => {
                  const modifiedDkgResult: DkgResult = { ...dkgResult }
                  modifiedDkgResult.submitterMemberIndex =
                    dkgResult.submitterMemberIndex + 1

                  await expect(
                    walletRegistry.challengeDkgResult(modifiedDkgResult)
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
              let dkgResult: DkgResult

              let tx: ContractTransaction

              before(async () => {
                await createSnapshot()
                ;({ dkgResult, dkgResultHash } =
                  await signAndSubmitUnrecoverableDkgResult(
                    walletRegistry,
                    groupPublicKey,
                    await selectGroup(sortitionPool, dkgSeed),
                    startBlock,
                    noMisbehaved
                  ))

                tx = await walletRegistry
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(walletRegistry, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "validation reverted"
                  )
              })

              it("should not unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.true
              })

              // it("should emit DkgMaliciousResultSlashed event", async () => {
              //   await expect(tx)
              //     .to.emit(walletRegistry, "DkgMaliciousResultSlashed")
              //     .withArgs(dkgResultHash, to1e18(50000), submitter.address)
              // })

              // it("should slash malicious result submitter", async () => {
              //   await expect(tx)
              //     .to.emit(staking, "Seized")
              //     .withArgs(to1e18(50000), 100, await thirdParty.getAddress(), [
              //       submitter.address,
              //     ])
              // })
            }
          )

          context("with correct dkg result submitted", async () => {
            let dkgResult: DkgResult

            before(async () => {
              await createSnapshot()
              ;({ dkgResult } = await signAndSubmitCorrectDkgResult(
                walletRegistry,
                groupPublicKey,
                dkgSeed,
                startBlock,
                noMisbehaved
              ))
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'unjustified challenge' error", async () => {
              await expect(
                walletRegistry.challengeDkgResult(dkgResult)
              ).to.be.revertedWith("unjustified challenge")
            })
          })
        })
      })
    })

    // This test checks that dkg timeout is adjusted in case of result challenges
    // to include the offset blocks that were mined until the invalid result
    // was challenged.
    it("should enforce submission start offset", async () => {
      await createSnapshot()

      let dkgResult: DkgResult

      const { startBlock } = await requestNewWallet(walletRegistry, walletOwner)

      await mineBlocks(constants.offchainDkgTime)

      // Submit result 1 at the beginning of the submission period
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        walletRegistry,
        groupPublicKey,
        operators,
        startBlock,
        noMisbehaved
      ))

      await expect(
        (
          await walletRegistry.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 1 after submission"
      ).to.equal(0)

      // Challenge result 1 at the beginning of the challenge period
      await walletRegistry.challengeDkgResult(dkgResult)
      // 1 block for dkg result submission tx +
      // 1 block for challenge tx
      let expectedSubmissionOffset = 2

      await expect(
        (
          await walletRegistry.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 1 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 2 in the middle of the submission period
      let blocksToMine =
        (constants.groupSize * params.dkgResultSubmissionEligibilityDelay) / 2
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        walletRegistry,
        groupPublicKey,
        operators,
        startBlock,
        noMisbehaved,
        shiftEligibleIndex(firstEligibleSubmitterIndex, constants.groupSize / 2)
      ))

      await expect(
        (
          await walletRegistry.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 2 after submission"
      ).to.equal(expectedSubmissionOffset) // same as before

      expectedSubmissionOffset += blocksToMine

      // Challenge result 2 in the middle of the challenge period
      await mineBlocks(params.dkgResultChallengePeriodLength / 2)
      expectedSubmissionOffset += params.dkgResultChallengePeriodLength / 2
      await walletRegistry.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await walletRegistry.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 2 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 3 at the end of the submission period
      blocksToMine =
        constants.groupSize * params.dkgResultSubmissionEligibilityDelay - 1
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        walletRegistry,
        groupPublicKey,
        operators,
        startBlock,
        noMisbehaved,
        shiftEligibleIndex(firstEligibleSubmitterIndex, constants.groupSize - 1)
      ))

      await expect(
        (
          await walletRegistry.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after submission"
      ).to.equal(expectedSubmissionOffset) // same as before

      expectedSubmissionOffset += blocksToMine

      // Challenge result 3 at the end of the challenge period
      blocksToMine = params.dkgResultChallengePeriodLength - 1
      await mineBlocks(blocksToMine)
      expectedSubmissionOffset += blocksToMine

      await expect(
        walletRegistry.callStatic.notifyDkgTimeout()
      ).to.be.revertedWith("DKG has not timed out")

      await walletRegistry.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await walletRegistry.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 4 after the submission period
      blocksToMine =
        constants.groupSize * params.dkgResultSubmissionEligibilityDelay
      await mineBlocks(blocksToMine)
      await expect(
        signAndSubmitArbitraryDkgResult(
          walletRegistry,
          groupPublicKey,
          operators,
          startBlock,
          noMisbehaved,
          shiftEligibleIndex(
            firstEligibleSubmitterIndex,
            constants.groupSize - 1
          )
        )
      ).to.be.revertedWith("DKG timeout already passed")

      await walletRegistry.notifyDkgTimeout()

      await restoreSnapshot()
    })
  })

  describe("notifyDkgTimeout", async () => {
    context("with initial contract state", async () => {
      it("should revert with 'DKG has not timed out' error", async () => {
        await expect(walletRegistry.notifyDkgTimeout()).to.be.revertedWith(
          "DKG has not timed out"
        )
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number

      before("start new wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock } = await requestNewWallet(walletRegistry, walletOwner))
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
            await expect(walletRegistry.notifyDkgTimeout()).to.be.revertedWith(
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
            await expect(walletRegistry.notifyDkgTimeout()).to.be.revertedWith(
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
            await expect(walletRegistry.notifyDkgTimeout()).to.be.revertedWith(
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

            // initialDkgRewardsPoolBalance = await walletRegistry.dkgRewardsPool()

            // initialNotifierBalance = await testToken.balanceOf(
            //   await thirdParty.getAddress()
            // )
            tx = await walletRegistry.connect(thirdParty).notifyDkgTimeout()
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should emit DkgTimedOut event", async () => {
            await expect(tx).to.emit(walletRegistry, "DkgTimedOut")
          })

          it("should clean dkg data", async () => {
            await assertDkgResultCleanData(walletRegistry)
          })

          // it("should reward the notifier with tokens from DKG rewards pool", async () => {
          //   const currentDkgRewardsPoolBalance =
          //     await walletRegistry.dkgRewardsPool()
          //   expect(
          //     initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
          //   ).to.be.equal(params.sortitionPoolUnlockingReward)

          //   const currentNotifierBalance: BigNumber = await testToken.balanceOf(
          //     await thirdParty.getAddress()
          //   )
          //   expect(
          //     currentNotifierBalance.sub(initialNotifierBalance)
          //   ).to.be.equal(params.sortitionPoolUnlockingReward)
          // })

          it("should unlock the sortition pool", async () => {
            await expect(await sortitionPool.isLocked()).to.be.false
          })
        })
      })
    })
  })

  describe("updateDkgParameters", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateDkgParameters(1, 2)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateDkgParameters(1, 2)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })

  describe("updateRewardParameters", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateRewardParameters(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateRewardParameters(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })

  describe("updateSlashingParameters", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateSlashingParameters(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateSlashingParameters(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })
})

async function assertDkgResultCleanData(walletRegistry: WalletRegistryStub) {
  const dkgData = await walletRegistry.getDkgData()

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

function mixOperators(operators: Operator[]): Operator[] {
  return operators
    .map((v) => ({ v, sort: Math.random() }))
    .sort((a, b) => a.sort - b.sort)
    .map(({ v }) => v)
}
