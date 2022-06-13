/* eslint-disable @typescript-eslint/no-unused-expressions */
import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import { params, randomBeaconDeployment } from "./fixtures"

import type { ContractTransaction, Signer } from "ethers"
import type { RandomBeaconStub, RandomBeaconGovernance } from "../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("RandomBeacon - Parameters", () => {
  const BLOCK_TIME = 15

  let governance: Signer
  let thirdParty: Signer
  let thirdPartyContract: SignerWithAddress
  let randomBeacon: RandomBeaconStub
  let randomBeaconGovernance: RandomBeaconGovernance

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;[thirdParty, thirdPartyContract] = await ethers.getSigners()
    ;({ governance } = await helpers.signers.getNamedSigners())

    const contracts = await waffle.loadFixture(randomBeaconDeployment)
    randomBeacon = contracts.randomBeacon as RandomBeaconStub
    randomBeaconGovernance =
      contracts.randomBeaconGovernance as RandomBeaconGovernance
  })

  describe("updateRelayEntryParameters", () => {
    const newRelayEntrySoftTimeout = 200
    const newRelayEntryHardTimeout = 300
    const newCallbackGasLimit = 400

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateRelayEntryParameters(
              newRelayEntrySoftTimeout,
              newRelayEntryHardTimeout,
              newCallbackGasLimit
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx1: ContractTransaction
      let tx2: ContractTransaction
      let tx3: ContractTransaction

      before(async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySoftTimeoutUpdate(newRelayEntrySoftTimeout)
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(newRelayEntryHardTimeout)
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(newCallbackGasLimit)

        await helpers.time.increaseTime(params.governanceDelay)

        tx1 = await randomBeaconGovernance
          .connect(governance)
          .finalizeRelayEntrySoftTimeoutUpdate()
        tx2 = await randomBeaconGovernance
          .connect(governance)
          .finalizeRelayEntryHardTimeoutUpdate()
        tx3 = await randomBeaconGovernance
          .connect(governance)
          .finalizeCallbackGasLimitUpdate()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should update the relay entry soft timeout", async () => {
        const { relayEntrySoftTimeout } =
          await randomBeacon.relayEntryParameters()
        expect(relayEntrySoftTimeout).to.be.equal(newRelayEntrySoftTimeout)
      })

      it("should update the relay entry hard timeout", async () => {
        const { relayEntryHardTimeout } =
          await randomBeacon.relayEntryParameters()
        expect(relayEntryHardTimeout).to.be.equal(newRelayEntryHardTimeout)
      })

      it("should update the callback gas limit", async () => {
        const { callbackGasLimit } = await randomBeacon.relayEntryParameters()
        expect(callbackGasLimit).to.be.equal(newCallbackGasLimit)
      })

      it("should emit the RelayEntryParametersUpdated event for soft timeout", async () => {
        await expect(tx1)
          .to.emit(randomBeacon, "RelayEntryParametersUpdated")
          .withArgs(
            newRelayEntrySoftTimeout,
            params.relayEntryHardTimeout,
            params.callbackGasLimit
          )
      })

      it("should emit the RelayEntryParametersUpdated event for hard timeout", async () => {
        await expect(tx2)
          .to.emit(randomBeacon, "RelayEntryParametersUpdated")
          .withArgs(
            newRelayEntrySoftTimeout,
            newRelayEntryHardTimeout,
            params.callbackGasLimit
          )
      })

      it("should emit the RelayEntryParametersUpdated event for gas callback", async () => {
        await expect(tx3)
          .to.emit(randomBeacon, "RelayEntryParametersUpdated")
          .withArgs(
            newRelayEntrySoftTimeout,
            newRelayEntryHardTimeout,
            newCallbackGasLimit
          )
      })
    })
  })

  describe("updateAuthorizationParameters", () => {
    const newMinimumAuthorization = 4200000
    const newAuthorizationDecreaseDelay = 86400
    const newAuthorizationDecreaseChangePeriod = 43200

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateAuthorizationParameters(
              newMinimumAuthorization,
              newAuthorizationDecreaseDelay,
              newAuthorizationDecreaseChangePeriod
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx1: ContractTransaction
      let tx2: ContractTransaction
      let tx3: ContractTransaction

      before(async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginMinimumAuthorizationUpdate(newMinimumAuthorization)
        await randomBeaconGovernance
          .connect(governance)
          .beginAuthorizationDecreaseDelayUpdate(newAuthorizationDecreaseDelay)
        await randomBeaconGovernance
          .connect(governance)
          .beginAuthorizationDecreaseChangePeriodUpdate(
            newAuthorizationDecreaseChangePeriod
          )

        await helpers.time.increaseTime(params.governanceDelay)

        tx1 = await randomBeaconGovernance
          .connect(governance)
          .finalizeMinimumAuthorizationUpdate()
        tx2 = await randomBeaconGovernance
          .connect(governance)
          .finalizeAuthorizationDecreaseDelayUpdate()
        tx3 = await randomBeaconGovernance
          .connect(governance)
          .finalizeAuthorizationDecreaseChangePeriodUpdate()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should update the group creation frequency", async () => {
        expect(await randomBeacon.minimumAuthorization()).to.be.equal(
          newMinimumAuthorization
        )
      })

      it("should update the authorization decrease delay", async () => {
        const { authorizationDecreaseDelay } =
          await randomBeacon.authorizationParameters()
        expect(authorizationDecreaseDelay).to.be.equal(
          authorizationDecreaseDelay
        )
      })

      it("should update the authorization decrease change period", async () => {
        const { authorizationDecreaseChangePeriod } =
          await randomBeacon.authorizationParameters()
        expect(authorizationDecreaseChangePeriod).to.be.equal(
          authorizationDecreaseChangePeriod
        )
      })

      it("should emit the AuthorizationParametersUpdated event for new minimum authorization", async () => {
        await expect(tx1)
          .to.emit(randomBeacon, "AuthorizationParametersUpdated")
          .withArgs(
            newMinimumAuthorization,
            params.authorizationDecreaseDelay,
            params.authorizationDecreaseChangePeriod
          )
      })

      it("should emit the AuthorizationParametersUpdated event for new authorization decrease delay", async () => {
        await expect(tx2)
          .to.emit(randomBeacon, "AuthorizationParametersUpdated")
          .withArgs(
            newMinimumAuthorization,
            newAuthorizationDecreaseDelay,
            params.authorizationDecreaseChangePeriod
          )
      })

      it("should emit the AuthorizationParametersUpdated event for new authorization decrease change period", async () => {
        await expect(tx3)
          .to.emit(randomBeacon, "AuthorizationParametersUpdated")
          .withArgs(
            newMinimumAuthorization,
            newAuthorizationDecreaseDelay,
            newAuthorizationDecreaseChangePeriod
          )
      })
    })
  })

  describe("updateGroupCreationParameters", () => {
    const newGroupCreationFrequency = 100
    const newGroupLifetime = (2 * 24 * 60 * 60) / BLOCK_TIME // 2days assuming 15s block time
    const newDkgResultChallengePeriodLength = 300
    const newDkgResultSubmissionTimeout = 400
    const newDkgSubmitterPrecedencePeriodLength = 200

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateGroupCreationParameters(
              newGroupCreationFrequency,
              newGroupLifetime,
              newDkgResultChallengePeriodLength,
              newDkgResultSubmissionTimeout,
              newDkgSubmitterPrecedencePeriodLength
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx1: ContractTransaction
      let tx2: ContractTransaction
      let tx3: ContractTransaction
      let tx4: ContractTransaction
      let tx5: ContractTransaction

      before(async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(newGroupCreationFrequency)
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(newGroupLifetime)
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(
            newDkgResultChallengePeriodLength
          )
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionTimeoutUpdate(newDkgResultSubmissionTimeout)
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgSubmitterPrecedencePeriodLengthUpdate(
            newDkgSubmitterPrecedencePeriodLength
          )

        await helpers.time.increaseTime(params.governanceDelay)

        tx1 = await randomBeaconGovernance
          .connect(governance)
          .finalizeGroupCreationFrequencyUpdate()
        tx2 = await randomBeaconGovernance
          .connect(governance)
          .finalizeGroupLifetimeUpdate()
        tx3 = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgResultChallengePeriodLengthUpdate()
        tx4 = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgResultSubmissionTimeoutUpdate()
        tx5 = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should update the group creation frequency", async () => {
        const { groupCreationFrequency } =
          await randomBeacon.groupCreationParameters()
        expect(groupCreationFrequency).to.be.equal(groupCreationFrequency)
      })

      it("should update the group lifetime", async () => {
        const { groupLifetime } = await randomBeacon.groupCreationParameters()
        expect(groupLifetime).to.be.equal(newGroupLifetime)
      })

      it("should update the DKG result challenge period length", async () => {
        const { dkgResultChallengePeriodLength } =
          await randomBeacon.groupCreationParameters()
        expect(dkgResultChallengePeriodLength).to.be.equal(
          newDkgResultChallengePeriodLength
        )
      })

      it("should update the DKG result submission timeout", async () => {
        const { dkgResultSubmissionTimeout } =
          await randomBeacon.groupCreationParameters()
        expect(dkgResultSubmissionTimeout).to.be.equal(
          newDkgResultSubmissionTimeout
        )
      })

      it("should update the DKG submitter precedence period", async () => {
        const { dkgSubmitterPrecedencePeriodLength } =
          await randomBeacon.groupCreationParameters()
        expect(dkgSubmitterPrecedencePeriodLength).to.be.equal(
          newDkgSubmitterPrecedencePeriodLength
        )
      })

      it("should emit the GroupCreationParametersUpdated event for new group creation frequency", async () => {
        await expect(tx1)
          .to.emit(randomBeacon, "GroupCreationParametersUpdated")
          .withArgs(
            newGroupCreationFrequency,
            params.groupLifeTime,
            params.dkgResultChallengePeriodLength,
            params.dkgResultSubmissionTimeout,
            params.dkgSubmitterPrecedencePeriodLength
          )
      })

      it("should emit the GroupCreationParametersUpdated event for new group life time", async () => {
        await expect(tx2)
          .to.emit(randomBeacon, "GroupCreationParametersUpdated")
          .withArgs(
            newGroupCreationFrequency,
            newGroupLifetime,
            params.dkgResultChallengePeriodLength,
            params.dkgResultSubmissionTimeout,
            params.dkgSubmitterPrecedencePeriodLength
          )
      })

      it("should emit the GroupCreationParametersUpdated event for new dkg result challenge period length", async () => {
        await expect(tx3)
          .to.emit(randomBeacon, "GroupCreationParametersUpdated")
          .withArgs(
            newGroupCreationFrequency,
            newGroupLifetime,
            newDkgResultChallengePeriodLength,
            params.dkgResultSubmissionTimeout,
            params.dkgSubmitterPrecedencePeriodLength
          )
      })

      it("should emit the GroupCreationParametersUpdated event for new dkg result submission timeout", async () => {
        await expect(tx4)
          .to.emit(randomBeacon, "GroupCreationParametersUpdated")
          .withArgs(
            newGroupCreationFrequency,
            newGroupLifetime,
            newDkgResultChallengePeriodLength,
            newDkgResultSubmissionTimeout,
            params.dkgSubmitterPrecedencePeriodLength
          )
      })

      it("should emit the GroupCreationParametersUpdated event for new dkg submitter precedence period length", async () => {
        await expect(tx5)
          .to.emit(randomBeacon, "GroupCreationParametersUpdated")
          .withArgs(
            newGroupCreationFrequency,
            newGroupLifetime,
            newDkgResultChallengePeriodLength,
            newDkgResultSubmissionTimeout,
            newDkgSubmitterPrecedencePeriodLength
          )
      })

      context("when values are invalid", () => {
        context(
          "when precedence period length is equal submission timeout",
          () => {
            it("should revert", async () => {
              const invalidDkgSubmitterPrecedencePeriodLength =
                newDkgResultSubmissionTimeout

              await randomBeaconGovernance
                .connect(governance)
                .beginDkgSubmitterPrecedencePeriodLengthUpdate(
                  invalidDkgSubmitterPrecedencePeriodLength
                )

              await helpers.time.increaseTime(params.governanceDelay)

              await expect(
                randomBeaconGovernance
                  .connect(governance)
                  .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
              ).to.be.revertedWith(
                "Submitter precedence period length should be less than the result submission timeout"
              )
            })
          }
        )

        context(
          "when precedence period length is greater than submission timeout",
          () => {
            it("should revert", async () => {
              const invalidDkgSubmitterPrecedencePeriodLength =
                newDkgResultSubmissionTimeout + 1

              await randomBeaconGovernance
                .connect(governance)
                .beginDkgSubmitterPrecedencePeriodLengthUpdate(
                  invalidDkgSubmitterPrecedencePeriodLength
                )

              await helpers.time.increaseTime(params.governanceDelay)

              await expect(
                randomBeaconGovernance
                  .connect(governance)
                  .finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
              ).to.be.revertedWith(
                "Submitter precedence period length should be less than the result submission timeout"
              )
            })
          }
        )
      })
    })
  })

  describe("updateRewardParameters", () => {
    const newSortitionPoolRewardsBanDuration = 400
    const newRelayEntryTimeoutNotificationRewardMultiplier = 10
    const newUnauthorizedSigningNotificationRewardMultiplier = 10
    const newDkgMaliciousResultNotificationRewardMultiplier = 20

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateRewardParameters(
              newSortitionPoolRewardsBanDuration,
              newRelayEntryTimeoutNotificationRewardMultiplier,
              newUnauthorizedSigningNotificationRewardMultiplier,
              newDkgMaliciousResultNotificationRewardMultiplier
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx1: ContractTransaction
      let tx2: ContractTransaction
      let tx3: ContractTransaction
      let tx4: ContractTransaction

      before(async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolRewardsBanDurationUpdate(
            newSortitionPoolRewardsBanDuration
          )
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(
            newRelayEntryTimeoutNotificationRewardMultiplier
          )
        await randomBeaconGovernance
          .connect(governance)
          .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(
            newUnauthorizedSigningNotificationRewardMultiplier
          )
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(
            newDkgMaliciousResultNotificationRewardMultiplier
          )

        await helpers.time.increaseTime(params.governanceDelay)

        tx1 = await randomBeaconGovernance
          .connect(governance)
          .finalizeSortitionPoolRewardsBanDurationUpdate()
        tx2 = await randomBeaconGovernance
          .connect(governance)
          .finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
        tx3 = await randomBeaconGovernance
          .connect(governance)
          .finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate()
        tx4 = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should update the sortition pool rewards ban duration", async () => {
        const { sortitionPoolRewardsBanDuration } =
          await randomBeacon.rewardParameters()
        expect(sortitionPoolRewardsBanDuration).to.be.equal(
          newSortitionPoolRewardsBanDuration
        )
      })

      it("should update the relay entry timeout notification reward multiplier", async () => {
        const { relayEntryTimeoutNotificationRewardMultiplier } =
          await randomBeacon.rewardParameters()
        expect(relayEntryTimeoutNotificationRewardMultiplier).to.be.equal(
          newRelayEntryTimeoutNotificationRewardMultiplier
        )
      })

      it("should update the unauthorized signing notification reward multiplier", async () => {
        const { unauthorizedSigningNotificationRewardMultiplier } =
          await randomBeacon.rewardParameters()
        expect(unauthorizedSigningNotificationRewardMultiplier).to.be.equal(
          newUnauthorizedSigningNotificationRewardMultiplier
        )
      })

      it("should update the DKG malicious result notification reward multiplier", async () => {
        const { dkgMaliciousResultNotificationRewardMultiplier } =
          await randomBeacon.rewardParameters()
        expect(dkgMaliciousResultNotificationRewardMultiplier).to.be.equal(
          newDkgMaliciousResultNotificationRewardMultiplier
        )
      })

      it("should emit the RewardParametersUpdated event for new sortition pool rewards ban duration", async () => {
        await expect(tx1)
          .to.emit(randomBeacon, "RewardParametersUpdated")
          .withArgs(
            newSortitionPoolRewardsBanDuration,
            params.relayEntryTimeoutNotificationRewardMultiplier,
            params.unauthorizedSigningNotificationRewardMultiplier,
            params.dkgMaliciousResultNotificationRewardMultiplier
          )
      })

      it("should emit the RewardParametersUpdated event for new relay entry timeout notification reward multiplier", async () => {
        await expect(tx2)
          .to.emit(randomBeacon, "RewardParametersUpdated")
          .withArgs(
            newSortitionPoolRewardsBanDuration,
            newRelayEntryTimeoutNotificationRewardMultiplier,
            params.unauthorizedSigningNotificationRewardMultiplier,
            params.dkgMaliciousResultNotificationRewardMultiplier
          )
      })

      it("should emit the RewardParametersUpdated event for new unauthorized signing notification reward multiplier", async () => {
        await expect(tx3)
          .to.emit(randomBeacon, "RewardParametersUpdated")
          .withArgs(
            newSortitionPoolRewardsBanDuration,
            newRelayEntryTimeoutNotificationRewardMultiplier,
            newUnauthorizedSigningNotificationRewardMultiplier,
            params.dkgMaliciousResultNotificationRewardMultiplier
          )
      })

      it("should emit the RewardParametersUpdated event for new dkg malicious result notification reward multiplier", async () => {
        await expect(tx4)
          .to.emit(randomBeacon, "RewardParametersUpdated")
          .withArgs(
            newSortitionPoolRewardsBanDuration,
            newRelayEntryTimeoutNotificationRewardMultiplier,
            newUnauthorizedSigningNotificationRewardMultiplier,
            newDkgMaliciousResultNotificationRewardMultiplier
          )
      })
    })
  })

  describe("updateSlashingParameters", () => {
    const newRelayEntrySubmissionFailureSlashingAmount = 100
    const newMaliciousDkgResultSlashingAmount = 200
    const newUnauthorizedSigningSlashingAmount = 150

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateSlashingParameters(
              newRelayEntrySubmissionFailureSlashingAmount,
              newMaliciousDkgResultSlashingAmount,
              newUnauthorizedSigningSlashingAmount
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx1: ContractTransaction
      let tx2: ContractTransaction
      let tx3: ContractTransaction

      before(async () => {
        await createSnapshot()

        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(
            newRelayEntrySubmissionFailureSlashingAmount
          )
        await randomBeaconGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(
            newMaliciousDkgResultSlashingAmount
          )
        await randomBeaconGovernance
          .connect(governance)
          .beginUnauthorizedSigningSlashingAmountUpdate(
            newUnauthorizedSigningSlashingAmount
          )

        await helpers.time.increaseTime(params.governanceDelay)

        tx1 = await randomBeaconGovernance
          .connect(governance)
          .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        tx2 = await randomBeaconGovernance
          .connect(governance)
          .finalizeMaliciousDkgResultSlashingAmountUpdate()
        tx3 = await randomBeaconGovernance
          .connect(governance)
          .finalizeUnauthorizedSigningSlashingAmountUpdate()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should update the relay entry submission failure slashing amount", async () => {
        const { relayEntrySubmissionFailureSlashingAmount } =
          await randomBeacon.slashingParameters()
        expect(relayEntrySubmissionFailureSlashingAmount).to.be.equal(
          newRelayEntrySubmissionFailureSlashingAmount
        )
      })

      it("should update the malicious DKG result slashing amount", async () => {
        const { maliciousDkgResultSlashingAmount } =
          await randomBeacon.slashingParameters()
        expect(maliciousDkgResultSlashingAmount).to.be.equal(
          newMaliciousDkgResultSlashingAmount
        )
      })

      it("should update the unauthorized signing slashing amount", async () => {
        const { unauthorizedSigningSlashingAmount } =
          await randomBeacon.slashingParameters()
        expect(unauthorizedSigningSlashingAmount).to.be.equal(
          newUnauthorizedSigningSlashingAmount
        )
      })

      it("should emit the SlashingParametersUpdated event for new relay entry submission failure slashing amount", async () => {
        await expect(tx1)
          .to.emit(randomBeacon, "SlashingParametersUpdated")
          .withArgs(
            newRelayEntrySubmissionFailureSlashingAmount,
            params.maliciousDkgResultSlashingAmount,
            params.unauthorizedSigningSlashingAmount
          )
      })

      it("should emit the SlashingParametersUpdated event new malicious dkg result slashing amount", async () => {
        await expect(tx2)
          .to.emit(randomBeacon, "SlashingParametersUpdated")
          .withArgs(
            newRelayEntrySubmissionFailureSlashingAmount,
            newMaliciousDkgResultSlashingAmount,
            params.unauthorizedSigningSlashingAmount
          )
      })

      it("should emit the SlashingParametersUpdated event new unauthorized signing slashing amount", async () => {
        await expect(tx3)
          .to.emit(randomBeacon, "SlashingParametersUpdated")
          .withArgs(
            newRelayEntrySubmissionFailureSlashingAmount,
            newMaliciousDkgResultSlashingAmount,
            newUnauthorizedSigningSlashingAmount
          )
      })
    })
  })

  describe("authorizedRequesters", () => {
    it("should be false by default", async () => {
      const isAuthorized = await randomBeacon.authorizedRequesters(
        thirdPartyContract.address
      )
      await expect(isAuthorized).to.be.false
    })
  })

  describe("setRequesterAuthorization", () => {
    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .setRequesterAuthorization(thirdPartyContract.address, true)
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      context("when authorizing a contract", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          tx = await randomBeaconGovernance
            .connect(governance)
            .setRequesterAuthorization(thirdPartyContract.address, true)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should set contract as authorized", async () => {
          const isAuthorized = await randomBeacon.authorizedRequesters(
            thirdPartyContract.address
          )
          expect(isAuthorized).to.be.true
        })

        it("should emit RequesterAuthorizationUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "RequesterAuthorizationUpdated")
            .withArgs(thirdPartyContract.address, true)
        })
      })

      context("when deauthorizing the contract", async () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          await randomBeaconGovernance
            .connect(governance)
            .setRequesterAuthorization(thirdPartyContract.address, true)

          tx = await randomBeaconGovernance
            .connect(governance)
            .setRequesterAuthorization(thirdPartyContract.address, false)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should set contract as not authorized", async () => {
          const isAuthorized = await randomBeacon.authorizedRequesters(
            thirdPartyContract.address
          )
          expect(isAuthorized).to.be.false
        })

        it("should emit RequesterAuthorizationUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "RequesterAuthorizationUpdated")
            .withArgs(thirdPartyContract.address, false)
        })
      })
    })
  })
})
