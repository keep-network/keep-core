import { waffle, helpers } from "hardhat"
import { expect } from "chai"
import { formatBytes32String } from "ethers/lib/utils"

import { walletRegistryFixture } from "./fixtures"
import { createNewWallet } from "./utils/wallets"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Wallet Creation", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry
  let walletOwner: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner } = await waffle.loadFixture(
      walletRegistryFixture
    ))
  })

  describe("getWallet", async () => {
    context("with wallet not registered", async () => {
      it("should revert", async () => {
        await expect(
          await walletRegistry.isWalletRegistered(
            formatBytes32String("NON EXISTING")
          )
        ).to.be.false
      })
    })

    context("with wallet registered", async () => {
      let walletID: string

      before("create a wallet", async () => {
        await createSnapshot()
        ;({ walletID } = await createNewWallet(walletRegistry, walletOwner))
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(await walletRegistry.isWalletRegistered(walletID)).to.be
          .true
      })
    })
  })
})
