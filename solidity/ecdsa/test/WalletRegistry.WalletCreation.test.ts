import { ethers, helpers } from "hardhat"
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
  signDkgResult,
} from "./utils/dkg"
import { selectGroup, hashUint32Array } from "./utils/groups"
import { createNewWallet } from "./utils/wallets"
import { submitRelayEntry } from "./utils/randomBeacon"

import type { IWalletOwner } from "../typechain/IWalletOwner"
import type { BigNumber, ContractTransaction, Signer } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  SortitionPool,
  WalletRegistry,
  WalletRegistryStub,
  StakingStub,
} from "../typechain"
import type { DkgResult, DkgResultSubmittedEventArgs } from "./utils/dkg"
import type { Operator } from "./utils/operators"
import type { FakeContract } from "@defi-wonderland/smock"

const { to1e18 } = helpers.number
const { mineBlocks, mineBlocksTo } = helpers.time
const { createSnapshot, restoreSnapshot } = helpers.snapshot

const { keccak256 } = ethers.utils

describe("WalletRegistry - Wallet Creation", async () => {
  const dkgTimeout: number = params.dkgResultSubmissionTimeout
  const groupPublicKey: string = ethers.utils.hexValue(
    ecdsaData.group1.publicKey
  )
  const walletID: string = ethers.utils.keccak256(groupPublicKey)

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
  let walletOwner: FakeContract<IWalletOwner>

  let deployer: SignerWithAddress
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
    } = await walletRegistryFixture())
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

        before("start wallet creation", async () => {
          await createSnapshot()
          tx = await walletRegistry
            .connect(walletOwner.wallet)
            .requestNewWallet()
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

        it("should not emit DkgStarted event", async () => {
          await expect(tx).not.to.emit(walletRegistry, "DkgStarted")
        })

        it("should transition DKG to AWAITING_SEED state", async () => {
          await expect(
            await walletRegistry.getWalletCreationState()
          ).to.be.equal(dkgState.AWAITING_SEED)
        })

        it("should not set dkg details", async () => {
          const dkgData = await walletRegistry.getDkgData()

          await expect(dkgData.seed).to.be.equal(0)
          await expect(dkgData.startBlock).to.be.equal(0)
        })

        it("should not register new wallet", async () => {
          await expect(tx).not.to.emit(walletRegistry, "WalletCreated")
        })
      })

      context("with new wallet requested", async () => {
        before("request new wallet", async () => {
          await createSnapshot()
          await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with relay entry not submitted", async () => {
          it("should revert with 'Current state is not IDLE' error", async () => {
            await expect(
              walletRegistry.connect(walletOwner.wallet).requestNewWallet()
            ).to.be.revertedWith("Current state is not IDLE")
          })

          context("with relay entry submitted", async () => {
            let startBlock: number
            let dkgSeed: BigNumber

            before("submit relay entry", async () => {
              await createSnapshot()
              ;({ startBlock, dkgSeed } = await submitRelayEntry(
                walletRegistry
              ))
            })

            after(async () => {
              await restoreSnapshot()
            })

            context("with dkg result not submitted", async () => {
              it("should revert with 'Current state is not IDLE' error", async () => {
                await expect(
                  walletRegistry.connect(walletOwner.wallet).requestNewWallet()
                ).to.be.revertedWith("Current state is not IDLE")
              })
            })

            context("with valid dkg result submitted", async () => {
              let dkgResult: DkgResult
              let submitter: SignerWithAddress

              before("submit dkg result", async () => {
                await createSnapshot()
                ;({ dkgResult, submitter } =
                  await signAndSubmitCorrectDkgResult(
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
                    walletRegistry
                      .connect(walletOwner.wallet)
                      .requestNewWallet()
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
                  await expect(
                    walletRegistry
                      .connect(walletOwner.wallet)
                      .requestNewWallet()
                  ).to.not.be.reverted
                })
              })
            })

            context("with invalid dkg result submitted", async () => {
              let dkgResult: DkgResult

              before("submit invalid dkg result", async () => {
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

              context("with dkg result challenged", async () => {
                before("challenge dkg result", async () => {
                  await createSnapshot()

                  await walletRegistry.challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should revert", async () => {
                  await expect(
                    walletRegistry
                      .connect(walletOwner.wallet)
                      .requestNewWallet()
                  ).to.be.revertedWith("Current state is not IDLE")
                })
              })
            })

            context("with dkg timeout notified", async () => {
              before("notify dkg timeout", async () => {
                await createSnapshot()

                await mineBlocksTo(startBlock + dkgTimeout)

                await walletRegistry.notifyDkgTimeout()
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should succeed", async () => {
                await expect(
                  walletRegistry.connect(walletOwner.wallet).requestNewWallet()
                ).not.to.be.reverted
              })
            })
          })
        })
      })
    })
  })

  // Tests for `__beaconCallback` were implemented in `WalletRegistry.RandomBeacon.test.ts`

  describe("getWalletCreationState", async () => {
    context("with initial contract state", async () => {
      it("should return IDLE state", async () => {
        expect(await walletRegistry.getWalletCreationState()).to.be.equal(
          dkgState.IDLE
        )
      })
    })

    context("with new wallet requested", async () => {
      before("request new wallet", async () => {
        await createSnapshot()
        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
        it("should return AWAITING_SEED state", async () => {
          expect(await walletRegistry.getWalletCreationState()).to.be.equal(
            dkgState.AWAITING_SEED
          )
        })

        context("with relay entry submitted", async () => {
          let startBlock: number
          let dkgSeed: BigNumber

          before(async () => {
            await createSnapshot()
            ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with dkg result not submitted", async () => {
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
                expect(
                  await walletRegistry.getWalletCreationState()
                ).to.be.equal(dkgState.AWAITING_RESULT)
              })
            })

            context("with valid dkg result submitted", async () => {
              let dkgResult: DkgResult
              let submitter: SignerWithAddress

              before("submit dkg result", async () => {
                await createSnapshot()
                ;({ dkgResult, submitter } =
                  await signAndSubmitCorrectDkgResult(
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
                it("should return CHALLENGE state", async () => {
                  expect(
                    await walletRegistry.getWalletCreationState()
                  ).to.be.equal(dkgState.CHALLENGE)
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

                it("should return IDLE state", async () => {
                  expect(
                    await walletRegistry.getWalletCreationState()
                  ).to.be.equal(dkgState.IDLE)
                })
              })
            })

            context("with invalid dkg result submitted", async () => {
              let dkgResult: DkgResult

              before("submit invalid dkg result", async () => {
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

              context("with dkg result challenged", async () => {
                before("challenge dkg result", async () => {
                  await createSnapshot()
                  await walletRegistry.challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should return AWAITING_RESULT state", async () => {
                  expect(
                    await walletRegistry.getWalletCreationState()
                  ).to.be.equal(dkgState.AWAITING_RESULT)
                })
              })
            })

            context("with dkg timeout notified", async () => {
              before("notify dkg timeout", async () => {
                await createSnapshot()

                await mineBlocksTo(startBlock + dkgTimeout)

                await walletRegistry.notifyDkgTimeout()
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should return IDLE state", async () => {
                expect(
                  await walletRegistry.getWalletCreationState()
                ).to.be.equal(dkgState.IDLE)
              })
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

    context("with new wallet requested", async () => {
      let requestNewWalletStartBlock: number

      before(async () => {
        await createSnapshot()
        const tx = await walletRegistry
          .connect(walletOwner.wallet)
          .requestNewWallet()

        requestNewWalletStartBlock = tx.blockNumber
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
        it("should return false", async () => {
          await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
        })

        context("at the end of the dkg timeout period", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(requestNewWalletStartBlock + dkgTimeout)
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

            await mineBlocksTo(requestNewWalletStartBlock + dkgTimeout + 1)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should return false", async () => {
            await expect(await walletRegistry.hasDkgTimedOut()).to.be.false
          })
        })

        context("with relay entry submitted", async () => {
          let startBlock: number
          let dkgSeed: BigNumber

          before(async () => {
            await createSnapshot()
            ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with dkg result not submitted", async () => {
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

            context("with valid dkg result submitted", async () => {
              let resultSubmissionBlock: number
              let dkgResult: DkgResult
              let submitter: SignerWithAddress

              before("submit dkg result", async () => {
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

              context("with dkg result not approved", async () => {
                context("at the end of the dkg timeout period", async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(startBlock + dkgTimeout)
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should return false", async () => {
                    await expect(await walletRegistry.hasDkgTimedOut()).to.be
                      .false
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
                    await expect(await walletRegistry.hasDkgTimedOut()).to.be
                      .false
                  })
                })

                context("at the end of the challenge period", async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      resultSubmissionBlock +
                        params.dkgResultChallengePeriodLength
                    )
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should return false", async () => {
                    await expect(await walletRegistry.hasDkgTimedOut()).to.be
                      .false
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
                    await expect(await walletRegistry.hasDkgTimedOut()).to.be
                      .false
                  })
                })
              })

              context("when dkg result was approved", async () => {
                before(async () => {
                  await createSnapshot()

                  await mineBlocksTo(
                    resultSubmissionBlock +
                      params.dkgResultChallengePeriodLength
                  )

                  await walletRegistry
                    .connect(submitter)
                    .approveDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should return false", async () => {
                  await expect(await walletRegistry.hasDkgTimedOut()).to.be
                    .false
                })
              })
            })

            context("with invalid dkg result submitted", async () => {
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

                context(
                  "at the end of dkg result submission period",
                  async () => {
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
                      await expect(await walletRegistry.hasDkgTimedOut()).to.be
                        .false
                    })
                  }
                )

                context("after dkg result submission period", async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      challengeBlockNumber +
                        params.dkgResultSubmissionTimeout +
                        1
                    )
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should return true", async () => {
                    await expect(await walletRegistry.hasDkgTimedOut()).to.be
                      .true
                  })
                })
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

    context("with new wallet requested", async () => {
      before("request new wallet", async () => {
        await createSnapshot()
        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
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

        context("with relay entry submitted", async () => {
          let startBlock: number
          let dkgSeed: BigNumber

          before(async () => {
            await createSnapshot()
            ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with wallet creation not timed out", async () => {
            context("with dkg result not submitted", async () => {
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
                    1,
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
                before(async () => {
                  await createSnapshot()
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should succeed", async () => {
                  await expect(
                    signAndSubmitArbitraryDkgResult(
                      walletRegistry,
                      groupPublicKey,
                      operators,
                      startBlock,
                      noMisbehaved,
                      1,
                      constants.groupThreshold - 1
                    )
                  ).to.not.be.reverted
                })
              })

              context("with the submission period started", async () => {
                beforeEach(async () => {
                  await createSnapshot()
                })

                afterEach(async () => {
                  await restoreSnapshot()
                })

                context(
                  "at the beginning of the submission period",
                  async () => {
                    it("should succeed for the first member", async () => {
                      await assertSubmissionSucceeds(1)
                    })

                    it("should succeed for the second member", async () => {
                      await assertSubmissionSucceeds(2)
                    })

                    it("should succeed for the last member", async () => {
                      await assertSubmissionSucceeds(constants.groupSize - 1)
                    })
                  }
                )

                context("at the end of the submission period", async () => {
                  before(async () => {
                    await createSnapshot()

                    await mineBlocksTo(
                      startBlock + params.dkgResultSubmissionTimeout - 1
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
                      startBlock + params.dkgResultSubmissionTimeout
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
            })

            context("with dkg result submitted", async () => {
              let dkgResult: DkgResult
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
                    resultSubmissionBlock +
                      params.dkgResultChallengePeriodLength
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

            context("with invalid dkg result submitted", async () => {
              let dkgResult: DkgResult
              let dkgResultHash: string
              let submitter: SignerWithAddress

              before(async () => {
                await createSnapshot()
                ;({ dkgResult, dkgResultHash, submitter } =
                  await signAndSubmitArbitraryDkgResult(
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

                context("with the same dkg result", async () => {
                  before(async () => {
                    await createSnapshot()
                  })

                  after(async () => {
                    await restoreSnapshot()
                  })

                  it("should succeed", async () => {
                    const tx = await walletRegistry
                      .connect(submitter)
                      .submitDkgResult(dkgResult)

                    await expectDkgResultSubmittedEvent(tx, {
                      resultHash: dkgResultHash,
                      seed: dkgSeed,
                      result: dkgResult,
                    })
                  })
                })

                context("with a fresh dkg result", async () => {
                  context("", async () => {
                    let tx: ContractTransaction
                    let expectedEventArgs: DkgResultSubmittedEventArgs

                    before(async () => {
                      await createSnapshot()

                      const {
                        transaction,
                        dkgResultHash: newDkgResultHash,
                        dkgResult: newDkgResult,
                      } = await signAndSubmitCorrectDkgResult(
                        walletRegistry,
                        ecdsaData.group2.publicKey,
                        dkgSeed,
                        startBlock,
                        noMisbehaved
                      )

                      tx = transaction

                      expectedEventArgs = {
                        resultHash: newDkgResultHash,
                        seed: dkgSeed,
                        result: newDkgResult,
                      }
                    })

                    after(async () => {
                      await restoreSnapshot()
                    })

                    it("should emit DkgResultSubmitted event", async () => {
                      await expectDkgResultSubmittedEvent(tx, expectedEventArgs)
                    })
                  })

                  context("with the submission period started", async () => {
                    let submissionStartBlockNumber: number

                    before(async () => {
                      await createSnapshot()

                      submissionStartBlockNumber = challengeBlockNumber

                      await mineBlocksTo(submissionStartBlockNumber)
                    })

                    after(async () => {
                      await restoreSnapshot()
                    })

                    context(
                      "at the beginning of the submission period",
                      async () => {
                        it("should succeed for the first member", async () => {
                          await assertSubmissionSucceeds(1)
                        })

                        it("should succeed for the second member", async () => {
                          await assertSubmissionSucceeds(2)
                        })

                        it("should succeed for the last member", async () => {
                          await assertSubmissionSucceeds(
                            constants.groupSize - 1
                          )
                        })
                      }
                    )

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
                })
              })
            })

            context("with misbehaved members", async () => {
              let tx: ContractTransaction
              let dkgResult: DkgResult
              let dkgResultHash: string

              context(
                "with misbehaved members in ascending order",
                async () => {
                  const misbehavedIndices = [2, 9, 11, 30, 60, 64]

                  before(async () => {
                    await createSnapshot()
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
            message = "DKG timeout already passed"
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

    context("with new wallet requested", async () => {
      before("request new wallet", async () => {
        await createSnapshot()
        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
        it("should revert with 'Current state is not CHALLENGE' error", async () => {
          await expect(
            walletRegistry.approveDkgResult(stubDkgResult)
          ).to.be.revertedWith("Current state is not CHALLENGE")
        })

        context("with relay entry submitted", async () => {
          let startBlock: number
          let dkgSeed: BigNumber

          before("submit relay entry", async () => {
            await createSnapshot()
            ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'Current state is not CHALLENGE' error", async () => {
            await expect(
              walletRegistry.approveDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
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
                context(
                  "when the third party is not yet eligible",
                  async () => {
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
                        walletRegistry
                          .connect(thirdParty)
                          .approveDkgResult(dkgResult)
                      ).to.be.revertedWith(
                        "Only the DKG result submitter can approve the result at this moment"
                      )
                    })
                  }
                )

                context("when the third party is eligible", async () => {
                  let tx: ContractTransaction
                  // let initialDkgRewardsPoolBalance: BigNumber
                  // let initApproverBalance: BigNumber

                  before(async () => {
                    await createSnapshot()

                    await mineBlocks(params.dkgSubmitterPrecedencePeriodLength)
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
            const maliciousSubmitter = 1

            // Submit a second result by another submitter
            const anotherSubmitterIndex = 6
            let anotherSubmitter: Signer

            before(async () => {
              await createSnapshot()

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
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should register a new wallet", async () => {
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
                ;({ dkgResult, submitter } =
                  await signAndSubmitCorrectDkgResult(
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
          walletOwner.wallet,
          existingWalletPublicKey
        ))

        await expect(await walletRegistry.isWalletRegistered(existingWalletID))
          .to.be.true
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with new wallet creation started", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before(
          "request new wallet creation and submit relay entry",
          async () => {
            await createSnapshot()
            await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
            ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
          }
        )

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

          const newResultPublicKey = ecdsaData.group2.publicKey
          const newWalletID = keccak256(newResultPublicKey)
          const newResultSubmitterIndex = 1

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

      context("with new wallet requested", async () => {
        before("request new wallet", async () => {
          await createSnapshot()
          await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
        })

        after(async () => {
          await restoreSnapshot()
        })

        context("with relay entry not submitted", async () => {
          it("should revert with 'Current state is not CHALLENGE'", async () => {
            await expect(
              walletRegistry.challengeDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })

          context("with relay entry submitted", async () => {
            let startBlock: number
            let dkgSeed: BigNumber

            before("submit relay entry", async () => {
              await createSnapshot()
              ;({ startBlock, dkgSeed } = await submitRelayEntry(
                walletRegistry
              ))
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

            context("with invalid dkg result submitted", async () => {
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
                    resultSubmissionBlock +
                      params.dkgResultChallengePeriodLength
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
                    const modifiedDkgResult: DkgResult = {
                      ...dkgResult,
                    }
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

              context("with token staking seize call failure", async () => {
                const slashingAmount = params.minimumAuthorization.add(1)

                let tx: Promise<ContractTransaction>

                before(async () => {
                  await createSnapshot()

                  await walletRegistry.setMaliciousDkgResultSlashingAmount(
                    slashingAmount
                  )

                  tx = walletRegistry
                    .connect(thirdParty)
                    .challengeDkgResult(dkgResult)
                })

                after(async () => {
                  await restoreSnapshot()
                })

                it("should succeed", async () => {
                  await expect(tx).to.not.be.reverted
                })

                it("should emit DkgMaliciousResultSlashingFailed", async () => {
                  await expect(tx)
                    .to.emit(walletRegistry, "DkgMaliciousResultSlashingFailed")
                    .withArgs(dkgResultHash, slashingAmount, submitter.address)
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
                    .withArgs(
                      to1e18(50000),
                      100,
                      await thirdParty.getAddress(),
                      [submitter.address]
                    )
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
    })

    context("with wallet registered", async () => {
      before("create a wallet", async () => {
        await createSnapshot()

        await createNewWallet(walletRegistry, walletOwner.wallet)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert with 'Current state is not CHALLENGE", async () => {
        await expect(
          walletRegistry.challengeDkgResult(stubDkgResult)
        ).to.be.revertedWith("Current state is not CHALLENGE")
      })

      context("with new wallet creation started", async () => {
        let startBlock: number
        let dkgSeed: BigNumber

        before(
          "request new wallet creation and submit relay entry",
          async () => {
            await createSnapshot()
            await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
            ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
          }
        )

        after(async () => {
          await restoreSnapshot()
        })

        it("should revert with 'Current state is not CHALLENGE'", async () => {
          await expect(
            walletRegistry.challengeDkgResult(stubDkgResult)
          ).to.be.revertedWith("Current state is not CHALLENGE")
        })

        context("with dkg result not submitted", async () => {
          it("should revert with 'Current state is not CHALLENGE'", async () => {
            await expect(
              walletRegistry.challengeDkgResult(stubDkgResult)
            ).to.be.revertedWith("Current state is not CHALLENGE")
          })
        })

        context("with invalid dkg result submitted", async () => {
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

    // This test checks that dkg timeout is adjusted in case of result challenges
    // to include the offset blocks that were mined until the invalid result
    // was challenged.
    describe("submission start offset", async () => {
      before(async () => {
        await createSnapshot()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should enforce submission start offset", async () => {
        let dkgResult: DkgResult

        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
        const { startBlock } = await submitRelayEntry(walletRegistry)

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
        let blocksToMine = params.dkgResultSubmissionTimeout / 2
        await mineBlocks(blocksToMine)
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
        blocksToMine = params.dkgResultSubmissionTimeout - 1
        await mineBlocks(blocksToMine)
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
        blocksToMine = params.dkgResultSubmissionTimeout
        await mineBlocks(blocksToMine)
        await expect(
          signAndSubmitArbitraryDkgResult(
            walletRegistry,
            groupPublicKey,
            operators,
            startBlock,
            noMisbehaved
          )
        ).to.be.revertedWith("DKG timeout already passed")

        await walletRegistry.notifyDkgTimeout()
      })
    })
  })

  describe("isDkgResultValid", async () => {
    context("with group creation not in progress", async () => {
      it("should revert with 'DKG has not been started'", async () => {
        await expect(
          walletRegistry.isDkgResultValid(stubDkgResult)
        ).to.be.revertedWith("DKG has not been started")
      })
    })

    context("with new wallet creation in progress", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      before("request new wallet creation and submit relay entry", async () => {
        await createSnapshot()
        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()
        ;({ startBlock, dkgSeed } = await submitRelayEntry(walletRegistry))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with invalid result", async () => {
        it("should return false and an error message", async () => {
          const expectedValidationResult = [false, "Malformed signatures array"]

          const validationResult = await walletRegistry.isDkgResultValid(
            stubDkgResult
          )

          await expect(validationResult).to.be.deep.equal(
            expectedValidationResult
          )
        })
      })

      context("with valid result", async () => {
        let dkgResult: DkgResult

        before("start new wallet creation", async () => {
          await createSnapshot()
          ;({ dkgResult } = await signDkgResult(
            await selectGroup(sortitionPool, dkgSeed),
            groupPublicKey,
            noMisbehaved,
            startBlock
          ))
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should return true", async () => {
          const expectedValidationResult = [true, ""]

          const validationResult = await walletRegistry.isDkgResultValid(
            dkgResult
          )

          await expect(validationResult).to.be.deep.equal(
            expectedValidationResult
          )
        })
      })
    })
  })

  describe("hasSeedTimedOut", async () => {
    context("with initial contract state", async () => {
      it("should return false", async () => {
        await expect(await walletRegistry.hasSeedTimedOut()).to.be.false
      })
    })

    context("with new wallet requested", async () => {
      let requestNewWalletStartBlock: number

      before(async () => {
        await createSnapshot()
        const tx = await walletRegistry
          .connect(walletOwner.wallet)
          .requestNewWallet()

        requestNewWalletStartBlock = tx.blockNumber
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
        context("with seed timeout period not passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(
              requestNewWalletStartBlock + params.dkgSeedTimeout
            )
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should return false", async () => {
            await expect(await walletRegistry.hasSeedTimedOut()).to.be.false
          })
        })

        context("with relay entry submitted", async () => {
          before(async () => {
            await createSnapshot()
            await submitRelayEntry(walletRegistry)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should return false", async () => {
            await expect(await walletRegistry.hasSeedTimedOut()).to.be.false
          })
        })

        context("with seed timeout period passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(
              requestNewWalletStartBlock + params.dkgSeedTimeout + 1
            )
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should return true", async () => {
            await expect(await walletRegistry.hasSeedTimedOut()).to.be.true
          })

          context("with relay entry submitted", async () => {
            before(async () => {
              await createSnapshot()
              await submitRelayEntry(walletRegistry)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should return false", async () => {
              await expect(await walletRegistry.hasSeedTimedOut()).to.be.false
            })
          })
        })
      })
    })
  })

  describe("notifySeedTimeout", async () => {
    context("with initial contract state", async () => {
      it("should revert with 'DKG has not timed out' error", async () => {
        await expect(walletRegistry.notifySeedTimeout()).to.be.revertedWith(
          "Awaiting seed has not timed out"
        )
      })
    })

    context("with new wallet requested", async () => {
      let requestNewWalletStartBlock: number

      before(async () => {
        await createSnapshot()
        const tx = await walletRegistry
          .connect(walletOwner.wallet)
          .requestNewWallet()

        requestNewWalletStartBlock = tx.blockNumber
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
        context("with seed timeout period not passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(
              requestNewWalletStartBlock + params.dkgSeedTimeout - 1
            )
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should revert with 'Awaiting seed has not timed out' error", async () => {
            await expect(walletRegistry.notifySeedTimeout()).to.be.revertedWith(
              "Awaiting seed has not timed out"
            )
          })

          context("with relay entry submitted", async () => {
            before(async () => {
              await createSnapshot()
              await submitRelayEntry(walletRegistry)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'Awaiting seed has not timed out' error", async () => {
              await expect(
                walletRegistry.notifySeedTimeout()
              ).to.be.revertedWith("Awaiting seed has not timed out")
            })
          })
        })

        context("with seed timeout period passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(
              requestNewWalletStartBlock + params.dkgSeedTimeout
            )
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("called by a third party", async () => {
            let tx: ContractTransaction

            before(async () => {
              await createSnapshot()

              tx = await walletRegistry.connect(thirdParty).notifySeedTimeout()
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should emit DkgSeedTimedOut event", async () => {
              await expect(tx).to.emit(walletRegistry, "DkgSeedTimedOut")
            })

            it("should clean dkg data", async () => {
              await assertDkgResultCleanData(walletRegistry)
            })

            it("should unlock the sortition pool", async () => {
              await expect(await sortitionPool.isLocked()).to.be.false
            })
          })

          context("with relay entry submitted", async () => {
            before(async () => {
              await createSnapshot()
              await submitRelayEntry(walletRegistry)
            })

            after(async () => {
              await restoreSnapshot()
            })

            it("should revert with 'Awaiting seed has not timed out' error", async () => {
              await expect(
                walletRegistry.notifySeedTimeout()
              ).to.be.revertedWith("Awaiting seed has not timed out")
            })
          })
        })
      })
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

    context("with new wallet requested", async () => {
      let requestNewWalletStartBlock: number

      before(async () => {
        await createSnapshot()
        const tx = await walletRegistry
          .connect(walletOwner.wallet)
          .requestNewWallet()

        requestNewWalletStartBlock = tx.blockNumber
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with relay entry not submitted", async () => {
        context("with dkg timeout period passed", async () => {
          before(async () => {
            await createSnapshot()

            await mineBlocksTo(requestNewWalletStartBlock + dkgTimeout + 1)
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

        context("with relay entry submitted", async () => {
          let startBlock: number

          before(async () => {
            await createSnapshot()
            ;({ startBlock } = await submitRelayEntry(walletRegistry))
          })

          after(async () => {
            await restoreSnapshot()
          })

          context("with dkg not timed out", async () => {
            context("with result submission period almost ended", async () => {
              before(async () => {
                await createSnapshot()

                await mineBlocksTo(startBlock + dkgTimeout - 1)
              })

              after(async () => {
                await restoreSnapshot()
              })

              it("should revert with 'DKG has not timed out' error", async () => {
                await expect(
                  walletRegistry.notifyDkgTimeout()
                ).to.be.revertedWith("DKG has not timed out")
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

function mixOperators(operators: Operator[]): Operator[] {
  return operators
    .map((v) => ({ v, sort: Math.random() }))
    .sort((a, b) => a.sort - b.sort)
    .map(({ v }) => v)
}
