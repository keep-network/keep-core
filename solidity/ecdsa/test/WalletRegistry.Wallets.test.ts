import { waffle, helpers } from "hardhat"
import { expect } from "chai"
import { formatBytes32String } from "ethers/lib/utils"

import { walletRegistryFixture } from "./fixtures"
import { createNewWallet } from "./utils/wallets"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Wallets", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry
  let walletOwner: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner } = await walletRegistryFixture())
  })

  describe("isWalletRegistered", async () => {
    context("with wallet not registered", async () => {
      it("should return false", async () => {
        await expect(
          await walletRegistry.isWalletRegistered(
            formatBytes32String("NON EXISTING")
          )
        ).to.be.false
      })
    })

    context("with wallet registered", async () => {
      let publicKeyHash: string

      before("create a wallet", async () => {
        await createSnapshot()
        ;({ publicKeyHash } = await createNewWallet(
          walletRegistry,
          walletOwner
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should return true", async () => {
        await expect(await walletRegistry.isWalletRegistered(publicKeyHash)).to
          .be.true
      })
    })
  })
})
