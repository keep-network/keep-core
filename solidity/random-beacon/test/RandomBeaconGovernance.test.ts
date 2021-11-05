import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import type { Signer } from "ethers"
import { randomBeaconDeployment } from "./fixtures"

import type {
  RandomBeacon,
  RandomBeaconGovernance,
  SortitionPoolStub,
} from "../typechain"

const fixture = async () => {
  const SortitionPoolStub = await ethers.getContractFactory("SortitionPoolStub")
  const sortitionPoolStub: SortitionPoolStub = await SortitionPoolStub.deploy()
  return randomBeaconDeployment(sortitionPoolStub)
}

describe("RandomBeaconGovernance", () => {
  let governance: Signer
  let thirdParty: Signer
  let randomBeacon: RandomBeacon
  let randomBeaconGovernance: RandomBeaconGovernance

  const initialRelayRequestFee = 100000
  const initialRelayEntrySubmissionEligibilityDelay = 10
  const initialRelayEntryHardTimeout = 100
  const initialCallbackGasLimit = 900000
  const initialGroupCreationFrequency = 4
  const initialGroupLifeTime = 60 * 60 * 24 * 7
  const initialDkgResultChallengePeriodLength = 60
  const initialDkgResultSubmissionEligibilityDelay = 10
  const initialDkgResultSubmissionReward = 500000
  const initialSortitionPoolUnlockingReward = 5000
  const initialRelayEntrySubmissionFailureSlashingAmount = 1000
  const initialMaliciousDkgResultSlashingAmount = 1000000000

  // prettier-ignore
  before(async () => {
    [governance, thirdParty] = await ethers.getSigners()
  })

  beforeEach(async () => {
    const contracts = await waffle.loadFixture(fixture)

    randomBeacon = contracts.randomBeacon as RandomBeacon

    await randomBeacon
      .connect(governance)
      .updateRelayEntryParameters(
        initialRelayRequestFee,
        initialRelayEntrySubmissionEligibilityDelay,
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
        initialDkgResultSubmissionEligibilityDelay
      )
    await randomBeacon
      .connect(governance)
      .updateRewardParameters(
        initialDkgResultSubmissionReward,
        initialSortitionPoolUnlockingReward
      )
    await randomBeacon
      .connect(governance)
      .updateSlashingParameters(
        initialRelayEntrySubmissionFailureSlashingAmount,
        initialMaliciousDkgResultSlashingAmount
      )

    const RandomBeaconGovernance = await ethers.getContractFactory(
      "RandomBeaconGovernance"
    )
    randomBeaconGovernance = await RandomBeaconGovernance.deploy(
      randomBeacon.address
    )
    await randomBeaconGovernance.deployed()
    await randomBeacon.transferOwnership(randomBeaconGovernance.address)
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

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayRequestFeeUpdate(123)
      })

      it("should not update the relay request fee", async () => {
        expect(await randomBeacon.relayRequestFee()).to.be.equal(
          initialRelayRequestFee
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayRequestFeeUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
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
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayRequestFeeUpdate(123)

        await helpers.time.increaseTime(11 * 60 * 60 - 10) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginRelayRequestFeeUpdate(123)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
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

  describe("beginRelayEntrySubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(0)
        ).to.be.revertedWith(
          "Relay entry submission eligibility delay must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(2)

        // works, did not revert
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)
      })

      it("should not update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(initialRelayEntrySubmissionEligibilityDelay)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntrySubmissionEligibilityDelayUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the RelayEntrySubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "RelayEntrySubmissionEligibilityDelayUpdateStarted"
          )
          .withArgs(1, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntrySubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(1)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        })

        it("should update the relay entry submission eligibility delay", async () => {
          expect(
            await randomBeacon.relayEntrySubmissionEligibilityDelay()
          ).to.be.equal(1)
        })

        it("should emit RelayEntrySubmissionEligibilityDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "RelayEntrySubmissionEligibilityDelayUpdated"
            )
            .withArgs(1)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingRelayEntrySubmissionEligibilityDelayUpdateTime()
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

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(123)
      })

      it("should not update the relay entry hard timeout", async () => {
        expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(
          initialRelayEntryHardTimeout
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntryHardTimeoutUpdateTime()
        ).to.be.equal(14 * 24 * 60 * 60) // 2 weeks
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
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(123)

        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntryHardTimeoutUpdate(123)

          await helpers.time.increaseTime(14 * 24 * 60 * 60) // 2 weeks

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
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
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(2)

        // works, did not revert
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
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(1000000)

        // works, did not revert
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(123)
      })

      it("should not update the callback gas limit", async () => {
        expect(await randomBeacon.callbackGasLimit()).to.be.equal(
          initialCallbackGasLimit
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingCallbackGasLimitUpdateTime()
        ).to.be.equal(14 * 24 * 60 * 60) // 2 weeks
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
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(123)

        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginCallbackGasLimitUpdate(123)

          await helpers.time.increaseTime(14 * 24 * 60 * 60) // 2 weeks

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
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
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(2)

        // works, did not revert
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)
      })

      it("should not update the group creation frequency timeout", async () => {
        expect(await randomBeacon.groupCreationFrequency()).to.be.equal(
          initialGroupCreationFrequency
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingGroupCreationFrequencyUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
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
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginGroupCreationFrequencyUpdate(1)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
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
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(24 * 60 * 60) // 24 hours

        // works, did not revert
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
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(14 * 24 * 60 * 60) // 14 days

        // works, did not revert
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days
      })

      it("should not update the group lifetime", async () => {
        expect(await randomBeacon.groupLifetime()).to.be.equal(
          initialGroupLifeTime
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingGroupLifetimeUpdateTime()
        ).to.be.equal(14 * 24 * 60 * 60) // 2 weeks
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
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

          await helpers.time.increaseTime(14 * 24 * 60 * 60) // 2 weeks

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
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
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(10)
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        // works, did not revert
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          initialDkgResultChallengePeriodLength
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
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
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(11)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
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

  describe("beginDkgResultSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgResultSubmissionEligibilityDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginDkgResultSubmissionEligibilityDelayUpdate(0)
        ).to.be.revertedWith(
          "DKG result submission eligibility delay must be > 0"
        )
      })
    })

    context("when the update value is at least one", () => {
      it("should accept the value", async () => {
        randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(1)
        randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(2)

        // works, did not revert
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(1)
      })

      it("should not update the DKG result submission eligibility delay", async () => {
        expect(
          await randomBeacon.dkgResultSubmissionEligibilityDelay()
        ).to.be.equal(initialDkgResultSubmissionEligibilityDelay)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultSubmissionEligibilityDelayUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
      })

      it("should emit the DkgResultSubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionEligibilityDelayUpdate(1)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginDkgResultSubmissionEligibilityDelayUpdate(1)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionEligibilityDelayUpdate()
        })

        it("should update the DKG result submission eligibility delay", async () => {
          expect(
            await randomBeacon.dkgResultSubmissionEligibilityDelay()
          ).to.be.equal(1)
        })

        it("should emit DkgResultSubmissionEligibilityDelayUpdated event", async () => {
          await expect(tx)
            .to.emit(
              randomBeaconGovernance,
              "DkgResultSubmissionEligibilityDelayUpdated"
            )
            .withArgs(1)
        })

        it("should reset the governance delay timer", async () => {
          await expect(
            randomBeaconGovernance.getRemainingDkgResultSubmissionEligibilityDelayUpdateTime()
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

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(123)
      })

      it("should not update the dkg result submission reward", async () => {
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(
          initialDkgResultSubmissionReward
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultSubmissionRewardUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
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
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(123)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginDkgResultSubmissionRewardUpdate(123)

          await helpers.time.increaseTime(24 * 60 * 60)

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
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

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(123)
      })

      it("should not update the sortition pool unlocking reward", async () => {
        expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(
          initialSortitionPoolUnlockingReward
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingSortitionPoolUnlockingRewardUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
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
        await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(123)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginSortitionPoolUnlockingRewardUpdate(123)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
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

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)
      })

      it("should not update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(initialRelayEntrySubmissionFailureSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime()
        ).to.be.equal(14 * 24 * 60 * 60) // 2 weeks
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
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

        await helpers.time.increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

          await helpers.time.increaseTime(14 * 24 * 60 * 60) // 2 weeks

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
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

      beforeEach(async () => {
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)
      })

      it("should not update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeacon.maliciousDkgResultSlashingAmount()
        ).to.be.equal(initialMaliciousDkgResultSlashingAmount)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.equal(12 * 60 * 60) // 12 hours
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
        await randomBeaconGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await helpers.time.increaseTime(11 * 60 * 60) // 11 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context(
      "when the update process is initialized and governance delay passed",
      () => {
        let tx

        beforeEach(async () => {
          await randomBeaconGovernance
            .connect(governance)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)

          await helpers.time.increaseTime(12 * 60 * 60) // 12 hours

          tx = await randomBeaconGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
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
})
