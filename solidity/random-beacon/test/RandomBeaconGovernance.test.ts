import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import { randomBeaconDeployment } from "./fixtures"

import type { ContractTransaction, Signer } from "ethers"
import type { RandomBeacon, RandomBeaconGovernance } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("RandomBeaconGovernance", () => {
  let governance: Signer
  let thirdParty: Signer
  let randomBeacon: RandomBeacon
  let randomBeaconGovernance: RandomBeaconGovernance

  const governanceDelay = 604800 // 1 week

  const initialRelayRequestFee = 100000
  const initialRelayEntrySoftTimeout = 10
  const initialRelayEntryHardTimeout = 100
  const initialCallbackGasLimit = 900000
  const initialGroupCreationFrequency = 4
  const initialGroupLifeTime = 60 * 60 * 24 * 7
  const initialDkgResultChallengePeriodLength = 60
  const initialDkgResultSubmissionTimeout = 200
  const initialDkgSubmitterPrecedencePeriodLength = 180
  const initialDkgResultSubmissionReward = 500000
  const initialSortitionPoolUnlockingReward = 5000
  const initialIneligibleOperatorNotifierReward = 6000
  const initialRelayEntrySubmissionFailureSlashingAmount = 1000
  const initialMaliciousDkgResultSlashingAmount = 1000000000
  const initialUnauthorizedSigningSlashingAmount = 1000000000
  const initialSortitionPoolRewardsBanDuration = 1209600
  const initialRelayEntryTimeoutNotificationRewardMultiplier = 5
  const initialUnauthorizedSignatureNotificationRewardMultiplier = 5
  const initialMinimumAuthorization = 1000000
  const initialAuthorizationDecreaseDelay = 86400
  const initialDkgMaliciousResultNotificationRewardMultiplier = 5

  // prettier-ignore
  before(async () => {
    [governance, thirdParty] = await ethers.getSigners()

    const contracts = await waffle.loadFixture(randomBeaconDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon

    await randomBeacon
      .connect(governance)
      .updateRelayEntryParameters(
        initialRelayRequestFee,
        initialRelayEntrySoftTimeout,
        initialRelayEntryHardTimeout,
        initialCallbackGasLimit
      )
    await randomBeacon
      .connect(governance)
      .updateGroupCreationParameters(
        initialGroupCreationFrequency,
        initialGroupLifeTime
      )
    await randomBeacon
      .connect(governance)
      .updateDkgParameters(
        initialDkgResultChallengePeriodLength,
        initialDkgResultSubmissionTimeout,
        initialDkgSubmitterPrecedencePeriodLength
      )
    await randomBeacon
      .connect(governance)
      .updateRewardParameters(
        initialDkgResultSubmissionReward,
        initialSortitionPoolUnlockingReward,
        initialIneligibleOperatorNotifierReward,
        initialSortitionPoolRewardsBanDuration,
        initialRelayEntryTimeoutNotificationRewardMultiplier,
        initialUnauthorizedSignatureNotificationRewardMultiplier,
        initialDkgMaliciousResultNotificationRewardMultiplier
      )
    await randomBeacon
      .connect(governance)
      .updateSlashingParameters(
        initialRelayEntrySubmissionFailureSlashingAmount,
        initialMaliciousDkgResultSlashingAmount,
        initialUnauthorizedSigningSlashingAmount
      )

    await randomBeacon
      .connect(governance)
      .updateAuthorizationParameters(
        initialMinimumAuthorization,
        initialAuthorizationDecreaseDelay
      )

    const RandomBeaconGovernance = await ethers.getContractFactory(
      "RandomBeaconGovernance"
    )
    randomBeaconGovernance = await RandomBeaconGovernance.deploy(
      randomBeacon.address, governanceDelay
    )
    await randomBeaconGovernance.deployed()
    await randomBeacon.transferOwnership(randomBeaconGovernance.address)
  })

  describe("beginGovernanceDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginGovernanceDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginGovernanceDelayUpdate(1337)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the governance delay", async () => {
        expect(await randomBeaconGovernance.governanceDelay()).to.be.equal(
          governanceDelay
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingGovernanceDelayUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit GovernanceDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "GovernanceDelayUpdateStarted")
          .withArgs(1337, blockTimestamp)
      })
    })
  })

  describe("finalizeGovernanceDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeGovernanceDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGovernanceDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      before(async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGovernanceDelayUpdate(7331)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGovernanceDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await randomBeaconGovernance
            .connect(governance)
            .beginGovernanceDelayUpdate(7331)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeGovernanceDelayUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the governance delay", async () => {
          expect(await randomBeaconGovernance.governanceDelay()).to.be.equal(
            7331
          )
        })

        it("should emit GovernanceDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "GovernanceDelayUpdated")
            .withArgs(7331)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingGovernanceDelayUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginRelayRequestFeeUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayRequestFeeUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the relay request fee", async () => {
        expect(await randomBeacon.relayRequestFee()).to.be.equal(
          initialRelayRequestFee
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayRequestFeeUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the RelayRequestFeeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "RelayRequestFeeUpdateStarted")
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayRequestFeeUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginRelayRequestFeeUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the relay request fee", async () => {
          expect(await randomBeacon.relayRequestFee()).to.be.equal(123)
        })

        it("should emit RelayRequestFeeUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "RelayRequestFeeUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingRelayRequestFeeUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginRelayEntrySoftTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginRelayEntrySoftTimeoutUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginRelayEntrySoftTimeoutUpdate(0)
        ).to.be.revertedWith("Relay entry soft timeout must be > 0")
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySoftTimeoutUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySoftTimeoutUpdate(2)

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySoftTimeoutUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the relay entry soft timeout", async () => {
        expect(await randomBeacon.relayEntrySoftTimeout()).to.be.equal(
          initialRelayEntrySoftTimeout
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntrySoftTimeoutUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the RelayEntrySoftTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "RelayEntrySoftTimeoutUpdateStarted")
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntrySoftTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeRelayEntrySoftTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySoftTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySoftTimeoutUpdate(1)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySoftTimeoutUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntrySoftTimeoutUpdate(1)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySoftTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the relay entry soft timeout", async () => {
          expect(await randomBeacon.relayEntrySoftTimeout()).to.be.equal(1)
        })

        it("should emit RelayEntrySoftTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "RelayEntrySoftTimeoutUpdated")
            .withArgs(1)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingRelayEntrySoftTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginRelayEntryHardTimeoutUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the relay entry hard timeout", async () => {
        expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(
          initialRelayEntryHardTimeout
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntryHardTimeoutUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the RelayEntryHardTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "RelayEntryHardTimeoutUpdateStarted")
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntryHardTimeoutUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the relay entry hard timeout", async () => {
          expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(123)
        })

        it("should emit RelayEntryHardTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "RelayEntryHardTimeoutUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingRelayEntryHardTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginCallbackGasLimitUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginCallbackGasLimitUpdate(0)
        ).to.be.revertedWith("Callback gas limit must be > 0 and <= 1000000")
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(2)

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the update value is more than one million", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginCallbackGasLimitUpdate(1000001)
        ).to.be.revertedWith("Callback gas limit must be > 0 and <= 1000000")
      })
    })

    context("when the update value is one million", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(1000000)

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the callback gas limit", async () => {
        expect(await randomBeacon.callbackGasLimit()).to.be.equal(
          initialCallbackGasLimit
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingCallbackGasLimitUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the CallbackGasLimitUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "CallbackGasLimitUpdateStarted")
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginCallbackGasLimitUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the callback gas limit", async () => {
          expect(await randomBeacon.callbackGasLimit()).to.be.equal(123)
        })

        it("should emit CallbackGasLimitUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "CallbackGasLimitUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingCallbackGasLimitUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginGroupCreationFrequencyUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginGroupCreationFrequencyUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginGroupCreationFrequencyUpdate(0)
        ).to.be.revertedWith("Group creation frequency must be > 0")
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(2)

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the group creation frequency timeout", async () => {
        expect(await randomBeacon.groupCreationFrequency()).to.be.equal(
          initialGroupCreationFrequency
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingGroupCreationFrequencyUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the GroupCreationFrequencyUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "GroupCreationFrequencyUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeGroupCreationFrequencyUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginGroupCreationFrequencyUpdate(1)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the group creation frequency", async () => {
          expect(await randomBeacon.groupCreationFrequency()).to.be.equal(1)
        })

        it("should emit GroupCreationFrequencyUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "GroupCreationFrequencyUpdated")
            .withArgs(1)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingGroupCreationFrequencyUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than one day", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginGroupLifetimeUpdate(23 * 60 * 60 - 1) // 24 hours - 1sec
        ).to.be.revertedWith("Group lifetime must be >= 1 day and <= 2 weeks")
      })
    })

    context("when the update value is one day", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(24 * 60 * 60) // 24 hours

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the update value is more than 2 weeks", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginGroupLifetimeUpdate(14 * 24 * 60 * 60 + 1) // 14 days + 1 sec
        ).to.be.revertedWith("Group lifetime must be >= 1 day and <= 2 weeks")
      })
    })

    context("when the update value is 2 weeks", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(14 * 24 * 60 * 60) // 14 days

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the group lifetime", async () => {
        expect(await randomBeacon.groupLifetime()).to.be.equal(
          initialGroupLifeTime
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingGroupLifetimeUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the GroupLifetimeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "GroupLifetimeUpdateStarted")
          .withArgs(2 * 24 * 60 * 60, blockTimestamp) // 2 days
      })
    })
  })

  describe("finalizeGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the group lifetime", async () => {
          expect(await randomBeacon.groupLifetime()).to.be.equal(
            2 * 24 * 60 * 60
          )
        })

        it("should emit GroupLifetimeUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "GroupLifetimeUpdated")
            .withArgs(2 * 24 * 60 * 60) // 2 days
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingGroupLifetimeUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgResultChallengePeriodLengthUpdate(11)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than 10", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(9)
        ).to.be.revertedWith("DKG result challenge period length must be >= 10")
      })
    })

    context("when the update value is at least 10", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(10)
        await randomBeaconGovernance
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

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          initialDkgResultChallengePeriodLength
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the DkgResultChallengePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
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

          await randomBeaconGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(11)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result challenge period length", async () => {
          expect(
            await randomBeacon.dkgResultChallengePeriodLength()
          ).to.be.equal(11)
        })

        it("should emit DkgResultChallengePeriodLengthUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "DkgResultChallengePeriodLengthUpdated"
            )
            .withArgs(11)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultSubmissionTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionTimeoutUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(0)
        ).to.be.revertedWith("DKG result submission timeout must be > 0")
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(2)

        // works, did not revert

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG result submission timeout", async () => {
        expect(await randomBeacon.dkgResultSubmissionTimeout()).to.be.equal(
          initialDkgResultSubmissionTimeout
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultSubmissionTimeoutUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the DkgResultSubmissionTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(1)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")

        await restoreSnapshot()
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        const newValue = 234
        let tx

        before(async () => {
          await createSnapshot()

          await randomBeaconGovernance
            .connect(governance)
            .beginDkgResultSubmissionTimeoutUpdate(newValue)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionTimeoutUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result submission timeout", async () => {
          expect(await randomBeacon.dkgResultSubmissionTimeout()).to.be.equal(
            newValue
          )
        })

        it("should emit DkgResultSubmissionTimeoutUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "DkgResultSubmissionTimeoutUpdated"
            )
            .withArgs(newValue)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingDkgResultSubmissionTimeoutUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgSubmitterPrecedencePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
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
          randomBeaconGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
        ).not.to.be.reverted

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(2)
        ).not.to.be.reverted

        await restoreSnapshot()
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG submitter precedence period length", async () => {
        expect(
          await randomBeacon.dkgSubmitterPrecedencePeriodLength()
        ).to.be.equal(initialDkgSubmitterPrecedencePeriodLength)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgSubmitterPrecedencePeriodLengthUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the DkgSubmitterPrecedencePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
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

          await randomBeaconGovernance
            .connect(governance)
            .beginDkgSubmitterPrecedencePeriodLengthUpdate(1)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG submitter precedence period length", async () => {
          expect(
            await randomBeacon.dkgSubmitterPrecedencePeriodLength()
          ).to.be.equal(1)
        })

        it("should emit DkgSubmitterPrecedencePeriodLengthUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "DkgSubmitterPrecedencePeriodLengthUpdated"
            )
            .withArgs(1)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingDkgSubmitterPrecedencePeriodLengthUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgResultSubmissionRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionRewardUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the dkg result submission reward", async () => {
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(
          initialDkgResultSubmissionReward
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultSubmissionRewardUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the DkgResultSubmissionRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "DkgResultSubmissionRewardUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultSubmissionRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginDkgResultSubmissionRewardUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the dkg result submission reward", async () => {
          expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(
            123
          )
        })

        it("should emit DkgResultSubmissionRewardUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "DkgResultSubmissionRewardUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingDkgResultSubmissionRewardUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginSortitionPoolUnlockingRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginSortitionPoolUnlockingRewardUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the sortition pool unlocking reward", async () => {
        expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(
          initialSortitionPoolUnlockingReward
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingSortitionPoolUnlockingRewardUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the SortitionPoolUnlockingRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "SortitionPoolUnlockingRewardUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeSortitionPoolUnlockingRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginSortitionPoolUnlockingRewardUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the sortition pool unlocking reward", async () => {
          expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(
            123
          )
        })

        it("should emit SortitionPoolUnlockingRewardUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "SortitionPoolUnlockingRewardUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingSortitionPoolUnlockingRewardUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginIneligibleOperatorNotifierRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginIneligibleOperatorNotifierRewardUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginIneligibleOperatorNotifierRewardUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the ineligible operator notifier reward", async () => {
        expect(
          await randomBeacon.ineligibleOperatorNotifierReward()
        ).to.be.equal(initialIneligibleOperatorNotifierReward)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingIneligibleOperatorNotifierRewardUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the IneligibleOperatorNotifierRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "IneligibleOperatorNotifierRewardUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeIneligibleOperatorNotifierRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeIneligibleOperatorNotifierRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeIneligibleOperatorNotifierRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginIneligibleOperatorNotifierRewardUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeIneligibleOperatorNotifierRewardUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginIneligibleOperatorNotifierRewardUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeIneligibleOperatorNotifierRewardUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the ineligible operator notifier reward", async () => {
          expect(
            await randomBeacon.ineligibleOperatorNotifierReward()
          ).to.be.equal(123)
        })

        it("should emit IneligibleOperatorNotifierRewardUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "IneligibleOperatorNotifierRewardUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingIneligibleOperatorNotifierRewardUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginRelayEntrySubmissionFailureSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(initialRelayEntrySubmissionFailureSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the RelayEntrySubmissionFailureSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "RelayEntrySubmissionFailureSlashingAmountUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntrySubmissionFailureSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the relay entry submission failure slashing amount", async () => {
          expect(
            await randomBeacon.relayEntrySubmissionFailureSlashingAmount()
          ).to.be.equal(123)
        })

        it("should emit RelayEntrySubmissionFailureSlashingAmountUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "RelayEntrySubmissionFailureSlashingAmountUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginUnauthorizedSigningSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginUnauthorizedSigningSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginUnauthorizedSigningSlashingAmountUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the unauthorized signing slashing amount", async () => {
        expect(
          await randomBeacon.unauthorizedSigningSlashingAmount()
        ).to.be.equal(initialUnauthorizedSigningSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingUnauthorizedSigningSlashingAmountUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the UnauthorizedSigningSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "UnauthorizedSigningSlashingAmountUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeUnauthorizedSigningSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeUnauthorizedSigningSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeUnauthorizedSigningSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginUnauthorizedSigningSlashingAmountUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeUnauthorizedSigningSlashingAmountUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginUnauthorizedSigningSlashingAmountUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeUnauthorizedSigningSlashingAmountUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the unauthorized signing slashing amount", async () => {
          expect(
            await randomBeacon.unauthorizedSigningSlashingAmount()
          ).to.be.equal(123)
        })

        it("should emit UnauthorizedSigningSlashingAmountUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "UnauthorizedSigningSlashingAmountUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingUnauthorizedSigningSlashingAmountUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeacon.maliciousDkgResultSlashingAmount()
        ).to.be.equal(initialMaliciousDkgResultSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the MaliciousDkgResultSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
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

          await randomBeaconGovernance
            .connect(governance)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the malicious DKG result slashing amount", async () => {
          expect(
            await randomBeacon.maliciousDkgResultSlashingAmount()
          ).to.be.equal(123)
        })

        it("should emit MaliciousDkgResultSlashingAmountUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "MaliciousDkgResultSlashingAmountUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginSortitionPoolRewardsBanDurationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginSortitionPoolRewardsBanDurationUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolRewardsBanDurationUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the sortition pool rewards ban duration", async () => {
        expect(
          await randomBeacon.sortitionPoolRewardsBanDuration()
        ).to.be.equal(initialSortitionPoolRewardsBanDuration)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingSortitionPoolRewardsBanDurationUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the SortitionPoolRewardsBanDurationUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "SortitionPoolRewardsBanDurationUpdateStarted"
          )
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeSortitionPoolRewardsBanDurationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolRewardsBanDurationUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginSortitionPoolRewardsBanDurationUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolRewardsBanDurationUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the sortition pool rewards ban duration", async () => {
          expect(
            await randomBeacon.sortitionPoolRewardsBanDuration()
          ).to.be.equal(123)
        })

        it("should emit SortitionPoolRewardsBanDurationUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "SortitionPoolRewardsBanDurationUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingSortitionPoolRewardsBanDurationUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginUnauthorizedSigningNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called with value >100", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(101)
        ).to.be.revertedWith("Maximum value is 100")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the unauthorized signing notification reward multiplier", async () => {
        expect(
          await randomBeacon.unauthorizedSigningNotificationRewardMultiplier()
        ).to.be.equal(initialUnauthorizedSignatureNotificationRewardMultiplier)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingUnauthorizedSigningNotificationRewardMultiplierUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the UnauthorizedSigningNotificationRewardMultiplierUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "UnauthorizedSigningNotificationRewardMultiplierUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(100)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(100)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the unauthorized signing notification reward multiplier", async () => {
          expect(
            await randomBeacon.unauthorizedSigningNotificationRewardMultiplier()
          ).to.be.equal(100)
        })

        it("should emit UnauthorizedSigningNotificationRewardMultiplierUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "UnauthorizedSigningNotificationRewardMultiplierUpdated"
            )
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingUnauthorizedSigningNotificationRewardMultiplierUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginRelayEntryTimeoutNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called with value >100", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(101)
        ).to.be.revertedWith("Maximum value is 100")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the relay entry timeout notification reward multiplier", async () => {
        expect(
          await randomBeacon.relayEntryTimeoutNotificationRewardMultiplier()
        ).to.be.equal(initialRelayEntryTimeoutNotificationRewardMultiplier)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntryTimeoutNotificationRewardMultiplierUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the RelayEntryTimeoutNotificationRewardMultiplierUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "RelayEntryTimeoutNotificationRewardMultiplierUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(100)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(100)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the relay entry timeout notification reward multiplier", async () => {
          expect(
            await randomBeacon.relayEntryTimeoutNotificationRewardMultiplier()
          ).to.be.equal(100)
        })

        it("should emit RelayEntryTimeoutNotificationRewardMultiplierUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "RelayEntryTimeoutNotificationRewardMultiplierUpdated"
            )
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingRelayEntryTimeoutNotificationRewardMultiplierUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginMinimumAuthorizationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginMinimumAuthorizationUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the minimum authorization amount", async () => {
        expect(await randomBeacon.minimumAuthorization()).to.be.equal(
          initialMinimumAuthorization
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingMimimumAuthorizationUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the MinimumAuthorizationUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernance, "MinimumAuthorizationUpdateStarted")
          .withArgs(123, blockTimestamp)
      })
    })
  })

  describe("finalizeMinimumAuthorizationUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(123)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
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

          await randomBeaconGovernance
            .connect(governance)
            .beginMinimumAuthorizationUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeMinimumAuthorizationUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the minimum authorization amount", async () => {
          expect(await randomBeacon.minimumAuthorization()).to.be.equal(123)
        })

        it("should emit MinimumAuthorizationUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeaconGovernance, "MinimumAuthorizationUpdated")
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingMimimumAuthorizationUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginAuthorizationDecreaseDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginAuthorizationDecreaseDelayUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(123)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the authorization decrease delay", async () => {
        expect(await randomBeacon.authorizationDecreaseDelay()).to.be.equal(
          initialAuthorizationDecreaseDelay
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingAuthorizationDecreaseDelayUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the AuthorizationDecreaseDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(123)
        await helpers.time.increaseTime(governanceDelay - 60) // -1min
        await expect(
          randomBeaconGovernance
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

          await randomBeaconGovernance
            .connect(governance)
            .beginAuthorizationDecreaseDelayUpdate(123)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeAuthorizationDecreaseDelayUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the authorization decrease delay", async () => {
          expect(await randomBeacon.authorizationDecreaseDelay()).to.be.equal(
            123
          )
        })

        it("should emit AuthorizationDecreaseDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "AuthorizationDecreaseDelayUpdated"
            )
            .withArgs(123)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingAuthorizationDecreaseDelayUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })

  describe("beginDkgMaliciousResultNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(100)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when called with value >100", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(101)
        ).to.be.revertedWith("Maximum value is 100")
      })
    })

    context("when the caller is the owner and value is correct", () => {
      let tx

      before(async () => {
        await createSnapshot()

        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(100)
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should not update the DKG malicious result notification reward multiplier", async () => {
        expect(
          await randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        ).to.be.equal(initialDkgMaliciousResultNotificationRewardMultiplier)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgMaliciousResultNotificationRewardMultiplierUpdateTime()
        ).to.be.equal(governanceDelay)
      })

      it("should emit the DkgMaliciousResultNotificationRewardMultiplierUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "DkgMaliciousResultNotificationRewardMultiplierUpdateStarted"
          )
          .withArgs(100, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(100)

        await helpers.time.increaseTime(governanceDelay - 60) // -1min

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
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

          await randomBeaconGovernance
            .connect(governance)
            .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(100)

          await helpers.time.increaseTime(governanceDelay)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG malicious result notification reward multiplier", async () => {
          expect(
            await randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
          ).to.be.equal(100)
        })

        it("should emit DkgMaliciousResultNotificationRewardMultiplierUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "DkgMaliciousResultNotificationRewardMultiplierUpdated"
            )
            .withArgs(100)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingDkgMaliciousResultNotificationRewardMultiplierUpdateTime()
          ).to.be.revertedWith("Change not initiated")
        })
      }
    )
  })
})
