import { ethers, waffle } from "hardhat"
import { expect } from "chai"

import type { Signer } from "ethers"
import { randomBeaconDeployment } from "./fixtures"
import type { RandomBeaconStub } from "../typechain"

describe("RandomBeacon - Parameters", () => {
  let governance: Signer
  let thirdParty: Signer
  let randomBeacon: RandomBeaconStub

  // prettier-ignore
  before(async () => {
    [governance, thirdParty] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(randomBeaconDeployment)
    randomBeacon = contracts.randomBeacon as RandomBeaconStub
  })

  describe("updateRelayEntryParameters", () => {
    const relayRequestFee = 100
    const relayEntrySubmissionEligibilityDelay = 200
    const relayEntryHardTimeout = 300
    const callbackGasLimit = 400

    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateRelayEntryParameters(
              relayRequestFee,
              relayEntrySubmissionEligibilityDelay,
              relayEntryHardTimeout,
              callbackGasLimit
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx
      beforeEach(async () => {
        tx = await randomBeacon
          .connect(governance)
          .updateRelayEntryParameters(
            relayRequestFee,
            relayEntrySubmissionEligibilityDelay,
            relayEntryHardTimeout,
            callbackGasLimit
          )
      })

      it("should update the relay request fee", async () => {
        expect(await randomBeacon.relayRequestFee()).to.be.equal(
          relayRequestFee
        )
      })

      it("should update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(relayEntrySubmissionEligibilityDelay)
      })

      it("should update the relay entry hard timeout", async () => {
        expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(
          relayEntryHardTimeout
        )
      })

      it("should update the callback gas limit", async () => {
        expect(await randomBeacon.callbackGasLimit()).to.be.equal(
          callbackGasLimit
        )
      })

      it("should emit the RelayEntryParametersUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "RelayEntryParametersUpdated")
          .withArgs(
            relayRequestFee,
            relayEntrySubmissionEligibilityDelay,
            relayEntryHardTimeout,
            callbackGasLimit
          )
      })
    })
  })

  describe("updateAuthorizationParameters", () => {
    const minimumAuthorization = 4200000
    const authorizationDecreaseDelay = 86400

    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateAuthorizationParameters(
              minimumAuthorization,
              authorizationDecreaseDelay
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeacon
          .connect(governance)
          .updateAuthorizationParameters(
            minimumAuthorization,
            authorizationDecreaseDelay
          )
      })

      it("should update the group creation frequency", async () => {
        expect(await randomBeacon.minimumAuthorization()).to.be.equal(
          minimumAuthorization
        )
      })

      it("should update the authorization decrease delay", async () => {
        expect(await randomBeacon.authorizationDecreaseDelay()).to.be.equal(
          authorizationDecreaseDelay
        )
      })

      it("should emit the AuthorizationParametersUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "AuthorizationParametersUpdated")
          .withArgs(minimumAuthorization, authorizationDecreaseDelay)
      })
    })
  })

  describe("updateGroupCreationParameters", () => {
    const groupCreationFrequency = 100
    const groupLifetime = 200

    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateGroupCreationParameters(
              groupCreationFrequency,
              groupLifetime
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx
      beforeEach(async () => {
        tx = await randomBeacon
          .connect(governance)
          .updateGroupCreationParameters(groupCreationFrequency, groupLifetime)
      })

      it("should update the group creation frequency", async () => {
        expect(await randomBeacon.groupCreationFrequency()).to.be.equal(
          groupCreationFrequency
        )
      })

      it("should update the group lifetime", async () => {
        expect(await randomBeacon.groupLifetime()).to.be.equal(groupLifetime)
      })

      it("should emit the GroupCreationParametersUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "GroupCreationParametersUpdated")
          .withArgs(groupCreationFrequency, groupLifetime)
      })
    })
  })

  describe("updateDkgParameters", () => {
    const dkgResultChallengePeriodLength = 300
    const dkgResultSubmissionEligibilityDelay = 400

    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateDkgParameters(
              dkgResultChallengePeriodLength,
              dkgResultSubmissionEligibilityDelay
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx
      beforeEach(async () => {
        tx = await randomBeacon
          .connect(governance)
          .updateDkgParameters(
            dkgResultChallengePeriodLength,
            dkgResultSubmissionEligibilityDelay
          )
      })

      it("should update the DKG result challenge period length", async () => {
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          dkgResultChallengePeriodLength
        )
      })

      it("should update the DKG result submission eligibility delay", async () => {
        expect(
          await randomBeacon.dkgResultSubmissionEligibilityDelay()
        ).to.be.equal(dkgResultSubmissionEligibilityDelay)
      })

      it("should emit the DkgParametersUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "DkgParametersUpdated")
          .withArgs(
            dkgResultChallengePeriodLength,
            dkgResultSubmissionEligibilityDelay
          )
      })
    })
  })

  describe("updateRewardParameters", () => {
    const dkgResultSubmissionReward = 100
    const sortitionPoolUnlockingReward = 200
    const ineligibleOperatorNotifierReward = 300
    const sortitionPoolRewardsBanDuration = 400
    const relayEntryTimeoutNotificationRewardMultiplier = 10
    const unauthorizedSigningNotificationRewardMultiplier = 10
    const dkgMaliciousResultNotificationRewardMultiplier = 20

    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateRewardParameters(
              dkgResultSubmissionReward,
              sortitionPoolUnlockingReward,
              ineligibleOperatorNotifierReward,
              sortitionPoolRewardsBanDuration,
              relayEntryTimeoutNotificationRewardMultiplier,
              unauthorizedSigningNotificationRewardMultiplier,
              dkgMaliciousResultNotificationRewardMultiplier
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx
      beforeEach(async () => {
        tx = await randomBeacon
          .connect(governance)
          .updateRewardParameters(
            dkgResultSubmissionReward,
            sortitionPoolUnlockingReward,
            ineligibleOperatorNotifierReward,
            sortitionPoolRewardsBanDuration,
            relayEntryTimeoutNotificationRewardMultiplier,
            unauthorizedSigningNotificationRewardMultiplier,
            dkgMaliciousResultNotificationRewardMultiplier
          )
      })

      it("should update the DKG result submission reward", async () => {
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(
          dkgResultSubmissionReward
        )
      })

      it("should update the sortition pool unlocking reward", async () => {
        expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(
          sortitionPoolUnlockingReward
        )
      })

      it("should update the ineligible operator notifier reward", async () => {
        expect(
          await randomBeacon.ineligibleOperatorNotifierReward()
        ).to.be.equal(ineligibleOperatorNotifierReward)
      })

      it("should update the sortition pool rewards ban duration", async () => {
        expect(
          await randomBeacon.sortitionPoolRewardsBanDuration()
        ).to.be.equal(sortitionPoolRewardsBanDuration)
      })

      it("should update the relay entry timeout notification reward multiplier", async () => {
        expect(
          await randomBeacon.relayEntryTimeoutNotificationRewardMultiplier()
        ).to.be.equal(relayEntryTimeoutNotificationRewardMultiplier)
      })

      it("should update the DKG malicious result notification reward multiplier", async () => {
        expect(
          await randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        ).to.be.equal(dkgMaliciousResultNotificationRewardMultiplier)
      })

      it("should emit the RewardParametersUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "RewardParametersUpdated")
          .withArgs(
            dkgResultSubmissionReward,
            sortitionPoolUnlockingReward,
            ineligibleOperatorNotifierReward,
            sortitionPoolRewardsBanDuration,
            relayEntryTimeoutNotificationRewardMultiplier,
            unauthorizedSigningNotificationRewardMultiplier,
            dkgMaliciousResultNotificationRewardMultiplier
          )
      })
    })
  })

  describe("updateSlashingParameters", () => {
    const relayEntrySubmissionFailureSlashingAmount = 100
    const maliciousDkgResultSlashingAmount = 200
    const unauthorizedSigningSlashingAmount = 150

    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateSlashingParameters(
              relayEntrySubmissionFailureSlashingAmount,
              maliciousDkgResultSlashingAmount,
              unauthorizedSigningSlashingAmount
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx
      beforeEach(async () => {
        tx = await randomBeacon
          .connect(governance)
          .updateSlashingParameters(
            relayEntrySubmissionFailureSlashingAmount,
            maliciousDkgResultSlashingAmount,
            unauthorizedSigningSlashingAmount
          )
      })

      it("should update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(relayEntrySubmissionFailureSlashingAmount)
      })

      it("should update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeacon.maliciousDkgResultSlashingAmount()
        ).to.be.equal(maliciousDkgResultSlashingAmount)
      })

      it("should update the unauthorized signing slashing amount", async () => {
        expect(
          await randomBeacon.unauthorizedSigningSlashingAmount()
        ).to.be.equal(unauthorizedSigningSlashingAmount)
      })

      it("should emit the SlashingParametersUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeacon, "SlashingParametersUpdated")
          .withArgs(
            relayEntrySubmissionFailureSlashingAmount,
            maliciousDkgResultSlashingAmount,
            unauthorizedSigningSlashingAmount
          )
      })
    })
  })
})
