import { waffle } from "hardhat"
import { expect } from "chai"

import { walletRegistryFixture } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"

describe("WalletRegistry - Wallet Creation", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry

  let deployer: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, deployer, thirdParty } = await waffle.loadFixture(
      walletRegistryFixture
    ))
  })

  describe("updateDkgParameters", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateDkgParameters(1, 2, 3)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateDkgParameters(1, 2, 3)
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

  describe("updateWalletOwner", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateWalletOwner(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .updateWalletOwner(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })
})
