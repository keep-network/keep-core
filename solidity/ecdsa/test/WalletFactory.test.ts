import {
  deployments,
  ethers,
  waffle,
  helpers,
  getUnnamedAccounts,
} from "hardhat"
import { expect } from "chai"

import type { BigNumber, ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { SortitionPool, WalletFactory } from "../typechain"
import { constants, params } from "./fixtures"
import ecdsaData from "./data/ecdsa"
import {
  calculateDkgSeed,
  DkgResult,
  noMisbehaved,
  signAndSubmitCorrectDkgResult,
} from "./utils/dkg"
import { registerOperators } from "./utils/operators"
import type { Operator } from "./utils/operators"

const { mineBlocks } = helpers.time

const fixture = async () => {
  await deployments.fixture(["WalletFactory"])

  const walletFactory: WalletFactory = await ethers.getContract("WalletFactory")
  const sortitionPool: SortitionPool = await ethers.getContract("SortitionPool")

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const signers = await registerOperators(
    walletFactory,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  return { walletFactory, sortitionPool, deployer, signers }
}

describe("WalletFactory", () => {
  const groupPublicKey: string = ethers.utils.hexValue(ecdsaData.groupPubKey)

  let walletFactory: WalletFactory
  let sortitionPool: SortitionPool

  let deployer: SignerWithAddress
  let signers: Operator[]

  beforeEach("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletFactory, sortitionPool, deployer, signers } =
      await waffle.loadFixture(fixture))
  })

  describe("createNewWallet", async () => {
    context("with initial contract state", async () => {
      let tx: ContractTransaction
      let expectedSeed: BigNumber

      beforeEach(async () => {
        tx = await walletFactory.createNewWallet()

        expectedSeed = calculateDkgSeed(
          await walletFactory.relayEntry(),
          tx.blockNumber
        )
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
          .withArgs(expectedSeed)
      })
    })

    context("with contract creation in progress", async () => {
      let startBlock: number
      let dkgSeed: BigNumber

      beforeEach("run contract creation", async () => {
        const tx = await walletFactory.createNewWallet()
        startBlock = tx.blockNumber

        dkgSeed = calculateDkgSeed(await walletFactory.relayEntry(), startBlock)
      })

      context("with dkg result not submitted", async () => {
        it("should revert", async () => {
          await expect(walletFactory.createNewWallet()).to.be.revertedWith(
            "Current state is not IDLE"
          )
        })
      })

      context("with dkg result submitted", async () => {
        let dkgResult: DkgResult
        let submitter: SignerWithAddress

        beforeEach(async () => {
          await mineBlocks(constants.offchainDkgTime)
          ;({ dkgResult, submitter } = await signAndSubmitCorrectDkgResult(
            walletFactory,
            groupPublicKey,
            dkgSeed,
            startBlock,
            noMisbehaved
          ))
        })

        // TODO: Add test cases to cover results that are approved, challenged or
        // pending.

        context("with dkg result not approved", async () => {
          it("should revert with 'current state is not IDLE' error", async () => {
            await expect(walletFactory.createNewWallet()).to.be.revertedWith(
              "Current state is not IDLE"
            )
          })
        })

        context("with dkg result challenged", async () => {})

        context("with dkg result approved", async () => {
          beforeEach(async () => {
            await mineBlocks(params.dkgResultChallengePeriodLength)

            await walletFactory.connect(submitter).approveDkgResult(dkgResult)
          })

          it.only("should emit DkgStarted event", async () => {
            const tx = await walletFactory.createNewWallet()
            const expectedSeed = calculateDkgSeed(
              await walletFactory.relayEntry(),
              tx.blockNumber
            )

            await expect(tx)
              .to.emit(walletFactory, "DkgStarted")
              .withArgs(expectedSeed)
          })
        })

        context("with dkg timeout notified", async () => {})
      })
    })
  })
})
