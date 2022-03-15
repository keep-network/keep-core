import { expect } from "chai"

import { walletRegistryFixture } from "./fixtures"

import type { IWalletOwner } from "../typechain/IWalletOwner"
import type { FakeContract } from "@defi-wonderland/smock"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"

describe("WalletRegistry - Parameters", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry

  let deployer: SignerWithAddress
  let walletOwner: FakeContract<IWalletOwner>
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
            .connect(walletOwner.wallet)
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
          walletRegistry
            .connect(walletOwner.wallet)
            .updateDkgParameters(1, 2, 3, 4)
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
          walletRegistry.connect(walletOwner.wallet).updateRewardParameters(1)
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
          walletRegistry.connect(walletOwner.wallet).updateSlashingParameters(1)
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
            .connect(walletOwner.wallet)
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

  describe("updateDkgResultSubmissionGas", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateDkgResultSubmissionGas(4200)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(walletOwner.wallet)
            .updateDkgResultSubmissionGas(4200)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateDkgResultSubmissionGas(4200)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })
  })

  describe("updateDkgResultApprovalGas", async () => {
    context("when called by the deployer", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(deployer).updateDkgResultApprovalGas(4200)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by the wallet owner", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry
            .connect(walletOwner.wallet)
            .updateDkgResultApprovalGas(4200)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).updateDkgResultApprovalGas(4200)
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
            .connect(walletOwner.wallet)
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
