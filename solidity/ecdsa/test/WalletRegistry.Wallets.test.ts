import { helpers } from "hardhat"
import { expect } from "chai"
import { formatBytes32String } from "ethers/lib/utils"

import { walletRegistryFixture } from "./fixtures"
import { createNewWallet } from "./utils/wallets"

import type { IWalletOwner } from "../typechain/IWalletOwner"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"
import type { FakeContract } from "@defi-wonderland/smock"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Wallets", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry
  let walletOwner: FakeContract<IWalletOwner>

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
      let walletID: string

      before("create a wallet", async () => {
        await createSnapshot()
        ;({ walletID } = await createNewWallet(
          walletRegistry,
          walletOwner.wallet
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should return true", async () => {
        await expect(await walletRegistry.isWalletRegistered(walletID)).to.be
          .true
      })
    })
  })
})
