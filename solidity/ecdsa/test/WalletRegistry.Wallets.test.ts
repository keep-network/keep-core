import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { keccak256 } from "ethers/lib/utils"

import ecdsaData from "./data/ecdsa"
import { createNewWallet } from "./utils/wallets"
import { walletRegistryFixture } from "./fixtures"

import type { ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const { HashZero } = ethers.constants

describe("WalletRegistry - Wallets", async () => {
  let walletRegistry: WalletRegistry

  let walletOwner: SignerWithAddress

  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner, thirdParty } = await waffle.loadFixture(
      walletRegistryFixture
    ))
  })

  describe("requestSignature", async () => {
    const digest = ecdsaData.group1.digest1

    context("when wallets are registered", async () => {
      let walletID1: string
      let walletID2: string
      const notExistingWalletID: string = keccak256("0x01234567")

      before("create wallets", async () => {
        await createSnapshot()
        ;({ walletID: walletID1 } = await createNewWallet(
          walletRegistry,
          walletOwner,
          ecdsaData.group1.publicKey
        ))
        ;({ walletID: walletID2 } = await createNewWallet(
          walletRegistry,
          walletOwner,
          ecdsaData.group2.publicKey
        ))
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("with first wallet", async () => {
        context("called by a third party", async () => {
          it("should revert", async () => {
            await expect(
              walletRegistry
                .connect(thirdParty)
                .requestSignature(walletID1, digest)
            ).to.be.revertedWith("Caller is not the Wallet Owner")
          })
        })

        context("called by the owner", async () => {
          let tx: ContractTransaction

          before("create wallets", async () => {
            await createSnapshot()

            tx = await walletRegistry
              .connect(walletOwner)
              .requestSignature(walletID1, digest)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should emit SignatureRequested event", async () => {
            await expect(tx)
              .to.emit(walletRegistry, "SignatureRequested")
              .withArgs(walletID1, digest)
          })

          it("should set digest to sign on the wallet", async () => {
            const walletData = await walletRegistry.getWallet(walletID1)

            await expect(walletData.digestToSign).to.be.equal(digest)
          })

          it("should not set digest to sign for another wallet", async () => {
            const walletData = await walletRegistry.getWallet(walletID2)

            await expect(walletData.digestToSign).to.be.equal(HashZero)
          })
        })
      })

      context("with not existing wallet ID", async () => {
        it("should revert", async () => {
          await expect(
            walletRegistry
              .connect(walletOwner)
              .requestSignature(notExistingWalletID, digest)
          ).to.be.revertedWith(
            "Wallet with given public key hash doesn't exist"
          )
        })
      })
    })
  })

  describe("submitSignature", async () => {
    context("when wallets are registered", async () => {
      let walletID1: string

      before("create wallets", async () => {
        await createSnapshot()
        ;({ walletID: walletID1 } = await createNewWallet(
          walletRegistry,
          walletOwner,
          ecdsaData.group1.publicKey
        ))
        await createNewWallet(
          walletRegistry,
          walletOwner,
          ecdsaData.group2.publicKey
        )
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("called by a third party", async () => {
        context("when signature was not requested", async () => {})
        context("when signature was requested", async () => {
          let tx: ContractTransaction

          before("request signature", async () => {
            await createSnapshot()

            await walletRegistry
              .connect(walletOwner)
              .requestSignature(walletID1, ecdsaData.group1.digest1)

            tx = await walletRegistry
              .connect(thirdParty)
              .submitSignature(walletID1, ecdsaData.group1.signature1)
          })

          after(async () => {
            await restoreSnapshot()
          })

          it("should emit SignatureSubmitted event", async () => {
            await expect(tx)
              .to.emit(walletRegistry, "SignatureSubmitted")
              .withArgs(walletID1, ecdsaData.group1.digest1, [
                ecdsaData.group1.signature1.v,
                ecdsaData.group1.signature1.r,
                ecdsaData.group1.signature1.s,
              ])
          })

          it("should reset digest to sign on the wallet", async () => {
            const walletData = await walletRegistry.getWallet(walletID1)

            await expect(walletData.digestToSign).to.be.equal(HashZero)
          })

          it("should not reset digest to sign for another wallet")
        })
      })

      // TODO: Add more tests
    })
  })
})
