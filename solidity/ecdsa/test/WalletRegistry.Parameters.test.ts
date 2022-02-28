import { expect } from "chai"

import { walletRegistryFixture } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"

describe("WalletRegistry - Parameters", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry

  let deployer: SignerWithAddress
  let walletOwner: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner, deployer, thirdParty } =
      await walletRegistryFixture())
  })

  describe("updateAuthorizationParameters", () => {
    context("when called by the deployer", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateAuthorizationParameters(1, 2)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by the wallet owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(walletOwner)
            .updateAuthorizationParameters(1, 2)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateAuthorizationParameters(1, 2)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })

  describe("updateDkgParameters", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateDkgParameters(1, 2, 3, 4)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(walletOwner).updateDkgParameters(1, 2, 3, 4)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateDkgParameters(1, 2, 3, 4)
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

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(walletOwner).updateRewardParameters(1)
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

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(walletOwner).updateSlashingParameters(1)
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

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(walletOwner)
            .updateWalletOwner(thirdParty.address)
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

  describe("upgradeRandomBeacon", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(deployer)
            .upgradeRandomBeacon(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(walletOwner)
            .upgradeRandomBeacon(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(thirdParty)
            .upgradeRandomBeacon(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })
})
