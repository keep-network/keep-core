import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import { walletRegistryFixture, params } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryGovernance } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot
const { to1e18 } = helpers.number

describe("WalletRegistryGovernance", async () => {
  let governance: SignerWithAddress
  let walletRegistry: WalletRegistry
  let walletRegistryGovernance: WalletRegistryGovernance
  let thirdParty: SignerWithAddress

  const initialMaliciousDkgResultSlashingAmount = to1e18(50000)
  const initialMaliciousDkgResultNotificationRewardMultiplier = 100

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletRegistryGovernance, governance, thirdParty } =
      await waffle.loadFixture(walletRegistryFixture))
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

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(10)
        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        // works, did not revert

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

  describe("beginDkgResultSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionEligibilityDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .beginDkgResultSubmissionEligibilityDelayUpdate(0)
        ).to.be.revertedWith(
          "DKG result submission eligibility delay must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(1)
        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(2)

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result submission eligibility delay", async () => {
        expect(
          (await walletRegistry.dkgParameters())
            .resultSubmissionEligibilityDelay
        ).to.be.equal(params.dkgResultSubmissionEligibilityDelay)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await walletRegistryGovernance.getRemainingDkgResultSubmissionEligibilityDelayUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the DkgResultSubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            walletRegistryGovernance,
            "DkgResultSubmissionEligibilityDelayUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await walletRegistryGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(1)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
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
            .beginDkgResultSubmissionEligibilityDelayUpdate(1)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await walletRegistryGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result submission eligibility delay", async () => {
          expect(
            (await walletRegistry.dkgParameters())
              .resultSubmissionEligibilityDelay
          ).to.be.equal(1)
        })

        it("should emit DkgResultSubmissionEligibilityDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(
              walletRegistryGovernance,
              "DkgResultSubmissionEligibilityDelayUpdated"
            )
            .withArgs(1)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            walletRegistryGovernance.getRemainingDkgResultSubmissionEligibilityDelayUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })
})
