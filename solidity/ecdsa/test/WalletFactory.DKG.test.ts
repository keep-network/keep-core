import {
  deployments,
  ethers,
  waffle,
  helpers,
  getUnnamedAccounts,
} from "hardhat"
import { expect } from "chai"

import type { BigNumber, ContractTransaction, Signer } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  SortitionPool,
  Wallet,
  WalletFactory,
  WalletFactoryStub,
} from "../typechain"
import { constants, dkgState, params } from "./fixtures"
import ecdsaData from "./data/ecdsa"
import {
  calculateDkgSeed,
  DkgResult,
  noMisbehaved,
  signAndSubmitArbitraryDkgResult,
  signAndSubmitCorrectDkgResult,
  signAndSubmitUnrecoverableDkgResult,
} from "./utils/dkg"
import { registerOperators } from "./utils/operators"
import type { Operator } from "./utils/operators"
import { selectGroup, hashUint32Array } from "./utils/groups"
import { firstEligibleIndex, shiftEligibleIndex } from "./utils/submission"

const { mineBlocks, mineBlocksTo } = helpers.time
const { createSnapshot, restoreSnapshot } = helpers.snapshot

const { keccak256 } = ethers.utils

const fixture = async () => {
  await deployments.fixture(["WalletFactory"])

  const walletFactory: WalletFactoryStub & WalletFactory =
    await ethers.getContract("WalletFactory")
  const sortitionPool: SortitionPool = await ethers.getContract("SortitionPool")

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")
  const walletManager: SignerWithAddress = await ethers.getNamedSigner(
    "walletManager"
  )

  const thirdParty = await ethers.getSigner((await getUnnamedAccounts())[0])

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const operators = await registerOperators(
    walletFactory,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  return {
    walletFactory,
    sortitionPool,
    deployer,
    walletManager,
    thirdParty,
    operators,
  }
}

describe("WalletFactory", () => {
  const dkgTimeout: number =
    constants.offchainDkgTime +
    constants.groupSize * params.dkgResultSubmissionEligibilityDelay
  const groupPublicKey: string = ethers.utils.hexValue(ecdsaData.groupPubKey)
  const groupPublicKeyHash: string = ethers.utils.keccak256(groupPublicKey)
  const firstEligibleSubmitterIndex: number = firstEligibleIndex(
    keccak256(ecdsaData.groupPubKey)
  )

  let walletFactory: WalletFactoryStub & WalletFactory
  let sortitionPool: SortitionPool

  let deployer: SignerWithAddress
  let walletManager: SignerWithAddress
  let thirdParty: SignerWithAddress
  let operators: Operator[]

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({
      walletFactory,
      sortitionPool,
      deployer,
      walletManager,
      thirdParty,
      operators,
    } = await waffle.loadFixture(fixture))
  })

  describe("requestNewWallet", async () => {
    context("when called by a deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletFactory.connect(deployer).requestNewWallet()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletFactory.connect(thirdParty).requestNewWallet()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a wallet manager", async () => {
      context("with initial contract state", async () => {
        let tx: ContractTransaction
        let dkgSeed: BigNumber

        before("start wallet creation", async () => {
          await createSnapshot()
          ;({ tx, dkgSeed } = await requestNewWallet(
            walletFactory,
            walletManager
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should lock the sortition pool", async () => {
          await expect(await sortitionPool.isLocked()).to.be.true
        })

        it("should emit DkgStateLocked event", async () => {
          await expect(tx).to.emit(walletFactory, "DkgStateLocked")
        })

        it("should emit DkgStarted event", async () => {
          await expect(tx)
            .to.emit(walletFactory, "DkgStarted")
            .withArgs(dkgSeed)
        })
      })

      context("with wallet creation in progress", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before("start wallet creation", async () => {
          await createSnapshot()
          ;({ startBlock, dkgSeed } = await requestNewWallet(
            walletFactory,
            walletManager
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with dkg result not submitted", async () => {
          it("should revert", async () => {
            await expect(
              walletFactory.connect(walletManager).requestNewWallet()
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
              walletFactory,
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
                walletFactory.connect(walletManager).requestNewWallet()
              ).to.be.revertedWith("Current state is not IDLE")
            })
          })

          context("with dkg result approved", async () => {
            before("approve dkg result", async () => {
              await createSnapshot()

              await mineBlocks(params.dkgResultChallengePeriodLength)

              await walletFactory.connect(submitter).approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should succeed", async () => {
              await expect(
                walletFactory.connect(walletManager).requestNewWallet()
              ).to.not.be.reverted
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
                walletFactory,
                groupPublicKey,
                // Mix operators to make the result malicious
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
                startBlock,
                noMisbehaved
              ))

              await walletFactory.challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert", async () => {
              await expect(
                walletFactory.connect(walletManager).requestNewWallet()
              ).to.be.revertedWith("Current state is not IDLE")
            })
          })
        })

        context("with dkg timeout notified", async () => {
          // TODO: Implement
          it("should succeed")
        })
      })
    })
  })

  describe("getWalletCreationState", async () => {
    context("with initial contract state", async () => {
      it("should return IDLE state", async () => {
        expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
          walletFactory,
          walletManager
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("at the start of off-chain dkg period", async () => {
        it("should return KEY_GENERATION state", async () => {
          expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
          expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
            expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
              expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
              walletFactory,
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
              expect(await walletFactory.getWalletCreationState()).to.be.equal(
                dkgState.CHALLENGE
              )
            })
          })

          context("when dkg result was approved", async () => {
            before("approve dkg result", async () => {
              await createSnapshot()

              await mineBlocks(params.dkgResultChallengePeriodLength)

              await walletFactory.connect(submitter).approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return IDLE state", async () => {
              expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
              walletFactory,
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
              await walletFactory.challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return AWAITING_RESULT state", async () => {
              expect(await walletFactory.getWalletCreationState()).to.be.equal(
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
        await expect(await walletFactory.hasDkgTimedOut()).to.be.false
      })
    })

    context("when wallet creation started", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletFactory,
          walletManager
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("within off-chain dkg period", async () => {
        it("should return false", async () => {
          await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
            await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
              await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
              await expect(await walletFactory.hasDkgTimedOut()).to.be.true
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
              walletFactory,
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
                await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
                await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
                await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
                await expect(await walletFactory.hasDkgTimedOut()).to.be.false
              })
            })
          })

          context("when dkg result was approved", async () => {
            before(async () => {
              await createSnapshot()

              await mineBlocksTo(
                resultSubmissionBlock + params.dkgResultChallengePeriodLength
              )

              await walletFactory.connect(submitter).approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return false", async () => {
              await expect(await walletFactory.hasDkgTimedOut()).to.be.false
            })
          })
        })

        context("when malicious dkg result was submitted", async () => {
          let dkgResult: DkgResult

          before(async () => {
            await createSnapshot()
            ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
              walletFactory,
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

              const tx = await walletFactory.challengeDkgResult(dkgResult)
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
                await expect(await walletFactory.hasDkgTimedOut()).to.be.false
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
                await expect(await walletFactory.hasDkgTimedOut()).to.be.true
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
            walletFactory,
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
          walletFactory,
          walletManager
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
                walletFactory,
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
                  walletFactory,
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
                await expect(tx)
                  .to.emit(walletFactory, "DkgResultSubmitted")
                  .withArgs(
                    dkgResultHash,
                    dkgSeed,
                    dkgResult.submitterMemberIndex,
                    dkgResult.groupPubKey,
                    dkgResult.misbehavedMembersIndices,
                    dkgResult.signatures,
                    dkgResult.signingMembersIndices,
                    dkgResult.members
                  )
              })

              it("should register a candidate wallet", async () => {
                const wallets = await walletFactory.getWallets()
                expect(wallets).to.be.lengthOf(1)

                const wallet: Wallet = await ethers.getContractAt(
                  "Wallet",
                  wallets[0]
                )

                expect(await wallet.publicKeyHash()).to.be.equal(
                  groupPublicKeyHash
                )
                expect(await wallet.activationBlockNumber()).to.be.equal(0)
                expect(await wallet.membersIdsHash()).to.be.equal(
                  hashUint32Array(dkgResult.members)
                )
              })

              it("should emit WalletCreated event", async () => {
                const wallet: Wallet = await ethers.getContractAt(
                  "Wallet",
                  (
                    await walletFactory.getWallets()
                  )[0]
                )

                await expect(tx)
                  .to.emit(walletFactory, "WalletCreated")
                  .withArgs(wallet.address, dkgResultHash)
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
                    walletFactory,
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
                walletFactory,
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
                  walletFactory,
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

                await walletFactory
                  .connect(submitter)
                  .approveDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should revert", async () => {
                await expect(
                  signAndSubmitCorrectDkgResult(
                    walletFactory,
                    groupPublicKey,
                    dkgSeed,
                    startBlock,
                    noMisbehaved
                  )
                ).to.be.revertedWith("Sortition pool unlocked")
              })
            })

            context("with dkg result challenged", async () => {
              let challengeBlockNumber: number

              before(async () => {
                await createSnapshot()

                const tx = await walletFactory.challengeDkgResult(dkgResult)
                challengeBlockNumber = tx.blockNumber
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should register a candidate wallet", async () => {
                await createSnapshot()

                const { dkgResult: dkgResult2 } =
                  await signAndSubmitCorrectDkgResult(
                    walletFactory,
                    groupPublicKey,
                    dkgSeed,
                    startBlock,
                    noMisbehaved
                  )

                const wallets = await walletFactory.getWallets()
                expect(wallets).to.be.lengthOf(1)

                const wallet: Wallet = await ethers.getContractAt(
                  "Wallet",
                  wallets[0]
                )

                expect(await wallet.publicKeyHash()).to.deep.equal(
                  groupPublicKeyHash
                )
                expect(await wallet.activationBlockNumber()).to.be.equal(0)
                expect(await wallet.membersIdsHash()).to.be.equal(
                  hashUint32Array(dkgResult2.members)
                )

                await restoreSnapshot()
              })

              it("should emit WalletCreated event", async () => {
                await createSnapshot()

                const { transaction: tx, dkgResultHash } =
                  await signAndSubmitCorrectDkgResult(
                    walletFactory,
                    groupPublicKey,
                    dkgSeed,
                    startBlock,
                    noMisbehaved
                  )

                const wallet: Wallet = await ethers.getContractAt(
                  "Wallet",
                  (
                    await walletFactory.getWallets()
                  )[0]
                )

                await expect(tx)
                  .to.emit(walletFactory, "WalletCreated")
                  .withArgs(wallet.address, dkgResultHash)

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
                    walletFactory,
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
                  await expect(tx)
                    .to.emit(walletFactory, "DkgResultSubmitted")
                    .withArgs(
                      dkgResultHash,
                      dkgSeed,
                      dkgResult.submitterMemberIndex,
                      dkgResult.groupPubKey,
                      dkgResult.misbehavedMembersIndices,
                      dkgResult.signatures,
                      dkgResult.signingMembersIndices,
                      dkgResult.members
                    )
                })

                it("should emit WalletCreated event", async () => {
                  const wallet: Wallet = await ethers.getContractAt(
                    "Wallet",
                    (
                      await walletFactory.getWallets()
                    )[0]
                  )

                  await expect(tx)
                    .to.emit(walletFactory, "WalletCreated")
                    .withArgs(wallet.address, dkgResultHash)
                })

                it("should correctly set a wallet members", async () => {
                  const wallet: Wallet = await ethers.getContractAt(
                    "Wallet",
                    (
                      await walletFactory.getWallets()
                    )[0]
                  )

                  // misbehavedIndices: [2, 9, 11, 30, 60, 64]
                  const expectedMembers = [...dkgResult.members]
                  expectedMembers.splice(1, 1) // index -1
                  expectedMembers.splice(7, 1) // index -2 (cause expectedMembers already shrinked)
                  expectedMembers.splice(8, 1) // index -3
                  expectedMembers.splice(26, 1) // index -4
                  expectedMembers.splice(55, 1) // index -5
                  expectedMembers.splice(58, 1) // index -6

                  expect(await wallet.membersIdsHash()).to.be.equal(
                    hashUint32Array(expectedMembers)
                  )
                })
              }
            )

            context(
              "when misbehaved members are not in ascending order",
              async () => {
                const misbehavedIndices = [2, 9, 30, 11, 60, 64]

                before(async () => {
                  await createSnapshot()

                  await mineBlocksTo(startBlock + constants.offchainDkgTime)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should revert", async () => {
                  await expect(
                    signAndSubmitCorrectDkgResult(
                      walletFactory,
                      groupPublicKey,
                      dkgSeed,
                      startBlock,
                      misbehavedIndices
                    )
                  ).to.be.revertedWith(
                    "Array accessed at an out-of-bounds or negative index"
                  )
                })
              }
            )
          })
        })
      })

      // TODO: Check challenge adjust start block calculation for eligibility
      // TODO: Check that challenges add up the delay

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
                walletFactory,
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
          walletFactory,
          groupPublicKey,
          dkgSeed,
          startBlock,
          noMisbehaved,
          submitterIndex
        )

        await expect(tx)
          .to.emit(walletFactory, "DkgResultSubmitted")
          .withArgs(
            dkgResultHash,
            dkgSeed,
            dkgResult.submitterMemberIndex,
            dkgResult.groupPubKey,
            dkgResult.misbehavedMembersIndices,
            dkgResult.signatures,
            dkgResult.signingMembersIndices,
            dkgResult.members
          )

        await restoreSnapshot()
      }

      async function assertSubmissionReverts(
        submitterIndex: number,
        message = "Submitter is not eligible"
      ): Promise<void> {
        await expect(
          signAndSubmitCorrectDkgResult(
            walletFactory,
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
    // Just to make `approveDkgResult` call possible.
    const stubDkgResult: DkgResult = {
      groupPubKey: ecdsaData.groupPubKey,
      members: [1, 2, 3, 4],
      misbehavedMembersIndices: [],
      signatures: "0x01020304",
      signingMembersIndices: [1, 2, 3, 4],
      submitterMemberIndex: 1,
    }

    context("with initial contract state", async () => {
      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          walletFactory.approveDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start new wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletFactory,
          walletManager
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          walletFactory.approveDkgResult(stubDkgResult)
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
              walletFactory.approveDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })
        })

        context("with dkg result submitted", async () => {
          let resultSubmissionBlock: number
          let dkgResultHash: string
          let dkgResult: DkgResult
          let submitter: SignerWithAddress
          let wallet: Wallet
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
              walletFactory,
              groupPublicKey,
              dkgSeed,
              startBlock,
              noMisbehaved,
              submitterIndex
            ))

            resultSubmissionBlock = tx.blockNumber

            const wallets = await walletFactory.getWallets()
            expect(wallets).to.be.lengthOf(1)

            wallet = await ethers.getContractAt("Wallet", wallets[0])
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
                walletFactory.connect(submitter).approveDkgResult(dkgResult)
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
                //   await walletFactory.dkgRewardsPool()
                // initialSubmitterBalance = await testToken.balanceOf(
                //   await submitter.getAddress()
                // )
                tx = await walletFactory
                  .connect(submitter)
                  .approveDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultApproved event", async () => {
                await expect(tx)
                  .to.emit(walletFactory, "DkgResultApproved")
                  .withArgs(dkgResultHash, await submitter.getAddress())
              })

              it("should clean dkg data", async () => {
                await assertDkgResultCleanData(walletFactory)
              })

              it("should activate a candidate wallet", async () => {
                expect(await wallet.activationBlockNumber()).to.be.equal(
                  tx.blockNumber
                )
              })

              // it("should reward the submitter with tokens from DKG rewards pool", async () => {
              //   const currentDkgRewardsPoolBalance =
              //     await walletFactory.dkgRewardsPool()
              //   expect(
              //     initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
              //   ).to.be.equal(params.dkgResultSubmissionReward)

              //   const currentSubmitterBalance: BigNumber =
              //     await testToken.balanceOf(await submitter.getAddress())
              //   expect(
              //     currentSubmitterBalance.sub(initialSubmitterBalance)
              //   ).to.be.equal(params.dkgResultSubmissionReward)
              // })

              it("should emit WalletActivated event", async () => {
                await expect(tx)
                  .to.emit(walletFactory, "WalletActivated")
                  .withArgs(wallet.address)
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
                    walletFactory
                      .connect(thirdParty)
                      .approveDkgResult(dkgResult)
                  ).to.be.revertedWith(
                    "Only the DKG result submitter can approve the result at this moment"
                  )
                })
              })

              context("when the third party is eligible", async () => {
                let tx: ContractTransaction
                let initialDkgRewardsPoolBalance: BigNumber
                let initApproverBalance: BigNumber

                before(async () => {
                  await createSnapshot()

                  await mineBlocks(params.dkgResultSubmissionEligibilityDelay)
                  // initialDkgRewardsPoolBalance =
                  //   await walletFactory.dkgRewardsPool()
                  // initApproverBalance = await testToken.balanceOf(
                  //   await thirdParty.getAddress()
                  // )
                  tx = await walletFactory
                    .connect(thirdParty)
                    .approveDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should succeed", async () => {
                  await expect(tx)
                    .to.emit(walletFactory, "WalletActivated")
                    .withArgs(wallet.address)
                })

                // it("should pay the reward to the third party", async () => {
                //   const currentDkgRewardsPoolBalance =
                //     await walletFactory.dkgRewardsPool()
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
          let wallet: Wallet

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
                walletFactory,
                groupPublicKey,
                // Mix operators to make the result malicious.
                mixOperators(await selectGroup(sortitionPool, dkgSeed)),
                startBlock,
                noMisbehaved,
                maliciousSubmitter
              )

            await walletFactory.challengeDkgResult(maliciousDkgResult)

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
              walletFactory,
              groupPublicKey,
              dkgSeed,
              startBlock,
              noMisbehaved,
              anotherSubmitterIndex
            ))

            resultSubmissionBlock = tx.blockNumber

            const wallets = await walletFactory.getWallets()
            expect(wallets).to.be.lengthOf(1)

            wallet = await ethers.getContractAt("Wallet", wallets[0])
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
                walletFactory
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

              // initialDkgRewardsPoolBalance =
              //   await walletFactory.dkgRewardsPool()

              // initialSubmitterBalance = await testToken.balanceOf(
              //   await anotherSubmitter.getAddress()
              // )

              tx = await walletFactory
                .connect(anotherSubmitter)
                .approveDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultApproved event", async () => {
              await expect(tx)
                .to.emit(walletFactory, "DkgResultApproved")
                .withArgs(dkgResultHash, await anotherSubmitter.getAddress())
            })

            it("should activate a candidate wallet", async () => {
              expect(await wallet.activationBlockNumber()).to.be.equal(
                tx.blockNumber
              )
            })

            // it("should reward the submitter with tokens from DKG rewards pool", async () => {
            //   const currentDkgRewardsPoolBalance =
            //     await walletFactory.dkgRewardsPool()
            //   expect(
            //     initialDkgRewardsPoolBalance.sub(currentDkgRewardsPoolBalance)
            //   ).to.be.equal(params.dkgResultSubmissionReward)

            //   const currentSubmitterBalance: BigNumber =
            //     await testToken.balanceOf(await anotherSubmitter.getAddress())
            //   expect(
            //     currentSubmitterBalance.sub(initialSubmitterBalance)
            //   ).to.be.equal(params.dkgResultSubmissionReward)
            // })

            it("should emit WalletActivated event", async () => {
              await expect(tx)
                .to.emit(walletFactory, "WalletActivated")
                .withArgs(wallet.address)
            })

            it("should unlock the sortition pool", async () => {
              await expect(await sortitionPool.isLocked()).to.be.false
            })
          })
        })
      })

      context("with max periods duration", async () => {
        let tx: ContractTransaction
        let wallet: Wallet

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout - 1)

          const { dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            walletFactory,
            groupPublicKey,
            dkgSeed,
            startBlock,
            noMisbehaved
          )

          const wallets = await walletFactory.getWallets()
          expect(wallets).to.be.lengthOf(1)

          wallet = await ethers.getContractAt("Wallet", wallets[0])

          await mineBlocks(params.dkgResultChallengePeriodLength)

          tx = await walletFactory
            .connect(submitter)
            .approveDkgResult(dkgResult)
        })

        after(async () => {
          await restoreSnapshot()
        })

        // Just an explicit assertion to make sure transaction passes correctly
        // for max periods duration.
        it("should succeed", async () => {
          await expect(tx)
            .to.emit(walletFactory, "WalletActivated")
            .withArgs(wallet.address)
        })

        it("should unlock the sortition pool", async () => {
          await expect(await sortitionPool.isLocked()).to.be.false
        })
      })

      context("with misbehaved operators", async () => {
        const misbehavedIndices = [2, 10, 64]
        let misbehavedIds
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await mineBlocksTo(startBlock + dkgTimeout - 1)

          const { dkgResult, members, submitter } =
            await signAndSubmitCorrectDkgResult(
              walletFactory,
              groupPublicKey,
              dkgSeed,
              startBlock,
              misbehavedIndices
            )

          misbehavedIds = misbehavedIndices.map((i) => members[i - 1])

          await mineBlocks(params.dkgResultChallengePeriodLength)
          tx = await walletFactory
            .connect(submitter)
            .approveDkgResult(dkgResult)
        })

        after(async () => {
          await restoreSnapshot()
        })

        // it("should ban misbehaved operators from sortition pool rewards", async () => {
        //   const now = await helpers.time.lastBlockTime()
        //   const expectedUntil = now + params.sortitionPoolRewardsBanDuration

        //   await expect(tx)
        //     .to.emit(sortitionPool, "IneligibleForRewards")
        //     .withArgs(misbehavedIds, expectedUntil)
        // })

        it("should clean dkg data", async () => {
          await assertDkgResultCleanData(walletFactory)
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

    //       dkgRewardsPoolBalance = await walletFactory.dkgRewardsPool()

    //       // Set the DKG result submission reward to twice the amount of test
    //       // tokens in the DKG rewards pool
    //       await walletFactoryGovernance.beginDkgResultSubmissionRewardUpdate(
    //         dkgRewardsPoolBalance.mul(2)
    //       )
    //       await helpers.time.increaseTime(12 * 60 * 60)
    //       await walletFactoryGovernance.finalizeDkgResultSubmissionRewardUpdate()

    //       const [genesisTx, genesisSeed] = await genesis(walletFactory)
    //       const startBlock: number = genesisTx.blockNumber
    //       await mineBlocksTo(startBlock + dkgTimeout - 1)

    //       let dkgResult: DkgResult
    //       ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
    //         walletFactory,
    //         groupPublicKey,
    //         genesisSeed,
    //         startBlock,
    //         noMisbehaved
    //       ))

    //       initApproverBalance = await testToken.balanceOf(
    //         await submitter.getAddress()
    //       )

    //       await mineBlocks(params.dkgResultChallengePeriodLength)
    //       tx = await walletFactory
    //         .connect(submitter)
    //         .approveDkgResult(dkgResult)
    //     })

    //     after(async () => {
    //       await restoreSnapshot()
    //     })

    //     it("should succeed", async () => {
    //       await expect(tx)
    //         .to.emit(walletFactory, "GroupActivated")
    //         .withArgs(0, groupPublicKey)
    //     })

    //     it("should pay the approver the whole DKG rewards pool balance", async () => {
    //       expect(await walletFactory.dkgRewardsPool()).to.be.equal(0)

    //       const currentApproverBalance = await testToken.balanceOf(
    //         await submitter.getAddress()
    //       )
    //       expect(currentApproverBalance.sub(initApproverBalance)).to.be.equal(
    //         dkgRewardsPoolBalance
    //       )
    //     })
    //   }
    // )
  })

  describe("challengeDkgResult", async () => {
    // Just to make `challengeDkgResult` call possible.
    const stubDkgResult: DkgResult = {
      groupPubKey: ecdsaData.groupPubKey,
      members: [1, 2, 3, 4],
      misbehavedMembersIndices: [],
      signatures: "0x01020304",
      signingMembersIndices: [1, 2, 3, 4],
      submitterMemberIndex: 1,
    }

    context("with initial contract state", async () => {
      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          walletFactory.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })
    })

    context("with group creation in progress", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("start new wallet creation", async () => {
        await createSnapshot()
        ;({ startBlock, dkgSeed } = await requestNewWallet(
          walletFactory,
          walletManager
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE' error", async () => {
        await expect(
          walletFactory.challengeDkgResult(stubDkgResult)
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
              walletFactory.challengeDkgResult(stubDkgResult)
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
              walletFactory,
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

                tx = await walletFactory
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(walletFactory, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "Invalid group members"
                  )
              })

              it("should remove a candidate wallet", async () => {
                const wallets = await walletFactory.getWallets()
                expect(wallets).to.be.lengthOf(0)
              })

              it("should emit CandidateGroupRemoved event", async () => {
                await expect(tx)
                  .to.emit(walletFactory, "CandidateGroupRemoved")
                  .withArgs(groupPublicKey)
              })

              it("should not unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.true
              })

              // it("should emit DkgMaliciousResultSlashed event", async () => {
              //   await expect(tx)
              //     .to.emit(walletFactory, "DkgMaliciousResultSlashed")
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

                tx = await walletFactory
                  .connect(thirdParty)
                  .challengeDkgResult(dkgResult)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should emit DkgResultChallenged event", async () => {
                await expect(tx)
                  .to.emit(walletFactory, "DkgResultChallenged")
                  .withArgs(
                    dkgResultHash,
                    await thirdParty.getAddress(),
                    "Invalid group members"
                  )
              })

              it("should remove a candidate wallet", async () => {
                const wallets = await walletFactory.getWallets()
                expect(wallets).to.be.lengthOf(0)
              })

              it("should emit CandidateGroupRemoved event", async () => {
                await expect(tx)
                  .to.emit(walletFactory, "CandidateGroupRemoved")
                  .withArgs(groupPublicKey)
              })

              it("should not unlock the sortition pool", async () => {
                await expect(await sortitionPool.isLocked()).to.be.true
              })

              // it("should emit DkgMaliciousResultSlashed event", async () => {
              //   await expect(tx)
              //     .to.emit(walletFactory, "DkgMaliciousResultSlashed")
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
                walletFactory.challengeDkgResult(dkgResult)
              ).to.be.revertedWith("Challenge period has already passed")
            })
          })
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
                  walletFactory,
                  groupPublicKey,
                  await selectGroup(sortitionPool, dkgSeed),
                  startBlock,
                  noMisbehaved
                ))

              tx = await walletFactory
                .connect(thirdParty)
                .challengeDkgResult(dkgResult)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgResultChallenged event", async () => {
              await expect(tx)
                .to.emit(walletFactory, "DkgResultChallenged")
                .withArgs(
                  dkgResultHash,
                  await thirdParty.getAddress(),
                  "validation reverted"
                )
            })

            it("should remove a candidate group", async () => {
              const wallets = await walletFactory.getWallets()
              expect(wallets).to.be.lengthOf(0)
            })

            it("should emit CandidateGroupRemoved event", async () => {
              await expect(tx)
                .to.emit(walletFactory, "CandidateGroupRemoved")
                .withArgs(groupPublicKey)
            })

            it("should not unlock the sortition pool", async () => {
              await expect(await sortitionPool.isLocked()).to.be.true
            })

            // it("should emit DkgMaliciousResultSlashed event", async () => {
            //   await expect(tx)
            //     .to.emit(walletFactory, "DkgMaliciousResultSlashed")
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
              walletFactory,
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
              walletFactory.challengeDkgResult(dkgResult)
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

      let dkgResult: DkgResult

      const { startBlock } = await requestNewWallet(
        walletFactory,
        walletManager
      )

      await mineBlocks(constants.offchainDkgTime)

      // Submit result 1 at the beginning of the submission period
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        walletFactory,
        groupPublicKey,
        operators,
        startBlock,
        noMisbehaved
      ))

      await expect(
        (
          await walletFactory.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 1 after submission"
      ).to.equal(0)

      // Challenge result 1 at the beginning of the challenge period
      await walletFactory.challengeDkgResult(dkgResult)
      // 1 block for dkg result submission tx +
      // 1 block for challenge tx
      let expectedSubmissionOffset = 2

      await expect(
        (
          await walletFactory.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 1 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 2 in the middle of the submission period
      let blocksToMine =
        (constants.groupSize * params.dkgResultSubmissionEligibilityDelay) / 2
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        walletFactory,
        groupPublicKey,
        operators,
        startBlock,
        noMisbehaved,
        shiftEligibleIndex(firstEligibleSubmitterIndex, constants.groupSize / 2)
      ))

      await expect(
        (
          await walletFactory.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 2 after submission"
      ).to.equal(expectedSubmissionOffset) // same as before

      expectedSubmissionOffset += blocksToMine

      // Challenge result 2 in the middle of the challenge period
      await mineBlocks(params.dkgResultChallengePeriodLength / 2)
      expectedSubmissionOffset += params.dkgResultChallengePeriodLength / 2
      await walletFactory.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await walletFactory.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 2 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 3 at the end of the submission period
      blocksToMine =
        constants.groupSize * params.dkgResultSubmissionEligibilityDelay - 1
      await mineBlocks(blocksToMine)
      ;({ dkgResult } = await signAndSubmitArbitraryDkgResult(
        walletFactory,
        groupPublicKey,
        operators,
        startBlock,
        noMisbehaved,
        shiftEligibleIndex(firstEligibleSubmitterIndex, constants.groupSize - 1)
      ))

      await expect(
        (
          await walletFactory.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after submission"
      ).to.equal(expectedSubmissionOffset) // same as before

      expectedSubmissionOffset += blocksToMine

      // Challenge result 3 at the end of the challenge period
      blocksToMine = params.dkgResultChallengePeriodLength - 1
      await mineBlocks(blocksToMine)
      expectedSubmissionOffset += blocksToMine

      await expect(
        walletFactory.callStatic.notifyDkgTimeout()
      ).to.be.revertedWith("DKG has not timed out")

      await walletFactory.challengeDkgResult(dkgResult)
      expectedSubmissionOffset += 2 // 1 block for dkg result submission tx + 1 block for challenge tx

      await expect(
        (
          await walletFactory.getDkgData()
        ).resultSubmissionStartBlockOffset,
        "invalid resultSubmissionStartBlockOffset for result 3 after challenge"
      ).to.equal(expectedSubmissionOffset)

      // Submit result 4 after the submission period
      blocksToMine =
        constants.groupSize * params.dkgResultSubmissionEligibilityDelay
      await mineBlocks(blocksToMine)
      await expect(
        signAndSubmitArbitraryDkgResult(
          walletFactory,
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

      await walletFactory.notifyDkgTimeout()

      await restoreSnapshot()
    })
  })

  // TODO: Add tests for notifyDkgTimeout
})

async function assertDkgResultCleanData(walletFactory: WalletFactoryStub) {
  const dkgData = await walletFactory.getDkgData()

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

async function requestNewWallet(
  walletFactory: WalletFactory,
  walletManager: SignerWithAddress
) {
  const tx: ContractTransaction = await walletFactory
    .connect(walletManager)
    .requestNewWallet()

  const startBlock: number = tx.blockNumber

  const dkgSeed: BigNumber = calculateDkgSeed(
    await walletFactory.relayEntry(),
    startBlock
  )

  return { tx, startBlock, dkgSeed }
}
