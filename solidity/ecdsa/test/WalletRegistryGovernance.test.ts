import { ethers, helpers } from "hardhat"
import { expect } from "chai"

import { walletRegistryFixture, params } from "./fixtures"

import type { ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryGovernance } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot
const { to1e18 } = helpers.number

describe("WalletRegistryGovernance", async () => {
  let governance: SignerWithAddress
  let walletRegistry: WalletRegistry
  let walletRegistryGovernance: WalletRegistryGovernance
  let thirdParty: SignerWithAddress
  let walletOwner: SignerWithAddress

  const initialMinimumAuthorization = to1e18(100000)
  const initialAuthorizationDecreaseDelay = 5260000 // 2 months
  const initialMaliciousDkgResultSlashingAmount = to1e18(50000)
  const initialMaliciousDkgResultNotificationRewardMultiplier = 100

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({
      walletRegistry,
      walletRegistryGovernance,
      governance,
      thirdParty,
      walletOwner,
    } = await walletRegistryFixture())
  })

  describe("upgradeRandomBeacon", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .upgradeRandomBeacon(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      context("when new address is zero", () => {
        it("should revert when a new random beacon address is zero", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .upgradeRandomBeacon(ethers.constants.AddressZero)
          ).to.be.revertedWith("New random beacon address cannot be zero")
        })
      })

      context("when new address is not zero", () => {
        before(async () => {
          await createSnapshot()

          tx = await walletRegistryGovernance
            .connect(governance)
            .upgradeRandomBeacon(thirdParty.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the random beacon", async () => {
          expect(await walletRegistry.randomBeacon()).to.be.equal(
            thirdParty.address
          )
        })

        it("should emit RandomBeaconUpgraded event", async () => {
          await expect(tx)
            .to.emit(walletRegistry, "RandomBeaconUpgraded")
            .withArgs(thirdParty.address)
        })
      })
    })
  })

  describe("beginWalletOwnerUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginWalletOwnerUpdate(thirdParty.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginWalletOwnerUpdate(thirdParty.address)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update when a new owner is zero address", async () => {
        it("should revert", async () => {
          await expect(
            walletRegistryGovernance
              .connect(governance)
              .beginWalletOwnerUpdate(ethers.constants.AddressZero)
          ).to.be.revertedWith("New wallet owner address cannot be zero")
        })
      })

      it("should not update the wallet owner", async () => {
        expect(await walletRegistry.walletOwner()).to.be.equal(
          walletOwner.address
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingWalletOwnerUpdateTime()
        ).to.be.equal(14 * 24 * 60 * 60) // 2 weeks
      })

      it("should emit the WalletOwnerUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(walletRegistryGovernance, "WalletOwnerUpdateStarted")
          .withArgs(thirdParty.address, blockTimestamp)
      })
    })
  })

  describe("finalizeWalletOwnerUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeWalletOwnerUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeWalletOwnerUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginWalletOwnerUpdate(thirdParty.address)

        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeWalletOwnerUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginWalletOwnerUpdate(thirdParty.address)

          await helpers.time.increaseTime(14 * 24 * 60 * 60) // 2 weeks

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeWalletOwnerUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the wallet owner", async () => {
          expect(await walletRegistry.walletOwner()).to.be.equal(
            thirdParty.address
          )
        })

        it("should emit WalletOwnerUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "WalletOwnerUpdated")
            .withArgs(thirdParty.address)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingWalletOwnerUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMinimumAuthorizationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginMinimumAuthorizationUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the minimum authorization amount", async () => {
        expect(await walletRegistry.minimumAuthorization()).to.be.equal(
          initialMinimumAuthorization
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingMimimumAuthorizationUpdateTime()
        ).to.be.equal(24 * 14 * 60 * 60) // 2 weeks
      })

      it("should emit the MinimumAuthorizationUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "MinimumAuthorizationUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeMinimumAuthorizationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(123)

        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginMinimumAuthorizationUpdate(123)

          await helpers.time.increaseTime(24 * 14 * 60 * 60) // 2 weeks

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the minimum authorization amount", async () => {
          expect(await walletRegistry.minimumAuthorization()).to.be.equal(123)
        })

        it("should emit MinimumAuthorizationUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "MinimumAuthorizationUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingMimimumAuthorizationUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginAuthorizationDecreaseDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginAuthorizationDecreaseDelayUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the authorization decrease delay", async () => {
        expect(await walletRegistry.authorizationDecreaseDelay()).to.be.equal(
          initialAuthorizationDecreaseDelay
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingAuthorizationDecreaseDelayUpdateTime()
        ).to.be.equal(24 * 14 * 60 * 60) // 2 weeks
      })

      it("should emit the AuthorizationDecreaseDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "AuthorizationDecreaseDelayUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeAuthorizationDecreaseDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(123)
        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginAuthorizationDecreaseDelayUpdate(123)

          await helpers.time.increaseTime(24 * 14 * 60 * 60) // 2 weeks

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the authorization decrease delay", async () => {
          expect(await walletRegistry.authorizationDecreaseDelay()).to.be.equal(
            123
          )
        })

        it("should emit AuthorizationDecreaseDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "AuthorizationDecreaseDelayUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingAuthorizationDecreaseDelayUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the malicious DKG result slashing amount", async () => {
        expect(
          await walletRegistry.maliciousDkgResultSlashingAmount()
        ).to.be.equal(initialMaliciousDkgResultSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the MaliciousDkgResultSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "MaliciousDkgResultSlashingAmountUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the malicious DKG result slashing amount", async () => {
          expect(
            await walletRegistry.maliciousDkgResultSlashingAmount()
          ).to.be.equal(123)
        })

        it("should emit MaliciousDkgResultSlashingAmountUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "MaliciousDkgResultSlashingAmountUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMaliciousDkgResultNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called with value >100", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(101)
        ).to.be.revertedWith("Maximum value is 100")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG malicious result notification reward multiplier", async () => {
        expect(
          await walletRegistry.maliciousDkgResultNotificationRewardMultiplier()
        ).to.be.equal(initialMaliciousDkgResultNotificationRewardMultiplier)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingMaliciousDkgResultNotificationRewardMultiplierUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the MaliciousDkgResultNotificationRewardMultiplierUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "MaliciousDkgResultNotificationRewardMultiplierUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginMaliciousDkgResultNotificationRewardMultiplierUpdate(100)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG malicious result notification reward multiplier", async () => {
          expect(
            await walletRegistry.maliciousDkgResultNotificationRewardMultiplier()
          ).to.be.equal(100)
        })

        it("should emit MaliciousDkgResultNotificationRewardMultiplierUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "MaliciousDkgResultNotificationRewardMultiplierUpdated"
            )
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingMaliciousDkgResultNotificationRewardMultiplierUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgSeedTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgSeedTimeoutUpdate(11)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is equal 0", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(0)
        ).to.be.revertedWith("DKG seed timeout must be > 0")
      })
    })

    context("when the update value is at least 1", () => {
      before(async () => {
        await createSnapshot()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(1)
        ).not.to.be.reverted
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(11)
        ).not.to.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgSeedTimeoutUpdate(11)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG seed timeout", async () => {
        expect((await walletRegistry.dkgParameters()).seedTimeout).to.be.equal(
          params.dkgSeedTimeout
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgSeedTimeoutUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the DkgSeedTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(walletRegistryGovernance, "DkgSeedTimeoutUpdateStarted")
          .withArgs(11, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgSeedTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgSeedTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSeedTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgSeedTimeoutUpdate(11)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSeedTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgSeedTimeoutUpdate(11)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSeedTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG seed timeout", async () => {
          expect(
            (await walletRegistry.dkgParameters()).seedTimeout
          ).to.be.equal(11)
        })

        it("should emit DkgSeedTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(walletRegistryGovernance, "DkgSeedTimeoutUpdated")
            .withArgs(11)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgSeedTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgResultChallengePeriodLengthUpdate(11)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than 10", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(9)
        ).to.be.revertedWith("DKG result challenge period length must be >= 10")
      })
    })

    context("when the update value is at least 10", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(10)
        ).to.not.be.reverted

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(11)
        ).to.not.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(
          (await walletRegistry.dkgParameters()).resultChallengePeriodLength
        ).to.be.equal(params.dkgResultChallengePeriodLength)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the DkgResultChallengePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgResultChallengePeriodLengthUpdateStarted"
          )
          .withArgs(11, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(11)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result challenge period length", async () => {
          expect(
            (await walletRegistry.dkgParameters()).resultChallengePeriodLength
          ).to.be.equal(11)
        })

        it("should emit DkgResultChallengePeriodLengthUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgResultChallengePeriodLengthUpdated"
            )
            .withArgs(11)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultSubmissionTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionTimeoutUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(0)
        ).to.be.revertedWith(
          "DKG result submission eligibility delay must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(1)
        ).to.not.be.reverted
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(2)
        ).to.not.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result submission eligibility delay", async () => {
        expect(
          (await walletRegistry.dkgParameters()).resultSubmissionTimeout
        ).to.be.equal(params.dkgResultSubmissionTimeout)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgResultSubmissionTimeoutUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the DkgResultSubmissionTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgResultSubmissionTimeoutUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultSubmissionTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(10)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result submission eligibility delay", async () => {
          expect(
            (await walletRegistry.dkgParameters()).resultSubmissionTimeout
          ).to.be.equal(10)
        })

        it("should emit DkgResultSubmissionTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgResultSubmissionTimeoutUpdated"
            )
            .withArgs(10)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgResultSubmissionTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgSubmitterPrecedencePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(0)
        ).to.be.revertedWith(
          "DKG submitter precedence period length must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
        ).to.not.be.reverted
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(2)
        ).to.not.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG submitter precedence period length", async () => {
        expect(
          (await walletRegistry.dkgParameters()).submitterPrecedencePeriodLength
        ).to.be.equal(params.dkgSubmitterPrecedencePeriodLength)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingSubmitterPrecedencePeriodLengthUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the DkgSubmitterPrecedencePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgSubmitterPrecedencePeriodLengthUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgSubmitterPrecedencePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        before(async () => {
          await createSnapshot()

          await walletRegistryGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(10)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG submitter precedence period length", async () => {
          expect(
            (await walletRegistry.dkgParameters())
              .submitterPrecedencePeriodLength
          ).to.be.equal(10)
        })

        it("should emit DkgSubmitterPrecedencePeriodLengthUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgSubmitterPrecedencePeriodLengthUpdated"
            )
            .withArgs(10)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingSubmitterPrecedencePeriodLengthUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })
})
