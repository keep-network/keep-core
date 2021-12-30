import { deployments, ethers } from "hardhat"
import { expect } from "chai"

import type { BigNumber, ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/dist/src/signers"
import type { SortitionPool, WalletFactory } from "../typechain"
import { calculateDkgSeed } from "./utils/walletFactory"

describe("WalletFactory", () => {
  let walletFactory: WalletFactory
  let sortitionPool: SortitionPool

  let deployer: SignerWithAddress

  beforeEach("load test fixture", async () => {
    await deployments.fixture(["WalletFactory"])

    walletFactory = await ethers.getContract("WalletFactory")
    sortitionPool = await ethers.getContract("SortitionPool")

    deployer = await ethers.getNamedSigner("deployer")
  })

  describe("createNewWallet", async () => {
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
})
