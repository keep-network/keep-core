/* eslint-disable @typescript-eslint/no-unused-expressions */
import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import { randomBeaconDeployment } from "./fixtures"

import type { ContractTransaction, Signer } from "ethers"
import type { RandomBeaconStub } from "../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("RandomBeacon - Parameters", () => {
  let governance: Signer
  let thirdParty: Signer
  let thirdPartyContract: SignerWithAddress
  let randomBeacon: RandomBeaconStub

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;[governance, thirdParty, thirdPartyContract] = await ethers.getSigners()

    const contracts = await waffle.loadFixture(randomBeaconDeployment)
    randomBeacon = contracts.randomBeacon as RandomBeaconStub
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
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await randomBeacon
          .connect(governance)
          .updateRelayEntryParameters(
            newRelayEntrySoftTimeout,
            newRelayEntryHardTimeout,
            newCallbackGasLimit
          )
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

      it("should emit the RelayEntryParametersUpdated event", async () => {
        await expect(tx)
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
    const minimumAuthorization = 4200000
    const authorizationDecreaseDelay = 86400

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateAuthorizationParameters(
              minimumAuthorization,
              authorizationDecreaseDelay
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await randomBeacon
          .connect(governance)
          .updateAuthorizationParameters(
            minimumAuthorization,
            authorizationDecreaseDelay
          )
      })

      after(async () => {
        await restoreSnapshot()
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

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateGroupCreationParameters(
              groupCreationFrequency,
              groupLifetime
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await randomBeacon
          .connect(governance)
          .updateGroupCreationParameters(groupCreationFrequency, groupLifetime)
      })

      after(async () => {
        await restoreSnapshot()
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
    const dkgResultSubmissionTimeout = 400
    const dkgSubmitterPrecedencePeriodLength = 200

    context("when the caller is not the governance", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateDkgParameters(
              dkgResultChallengePeriodLength,
              dkgResultSubmissionTimeout,
              dkgSubmitterPrecedencePeriodLength
            )
        ).to.be.revertedWith("Caller is not the governance")
      })
    })

    context("when the caller is the governance", () => {
      context("when values are valid", () => {
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          tx = await randomBeacon
            .connect(governance)
            .updateDkgParameters(
              dkgResultChallengePeriodLength,
              dkgResultSubmissionTimeout,
              dkgSubmitterPrecedencePeriodLength
            )
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should update the DKG result challenge period length", async () => {
          expect(
            await randomBeacon.dkgResultChallengePeriodLength()
          ).to.be.equal(dkgResultChallengePeriodLength)
        })

        it("should update the DKG result submission timeout", async () => {
          expect(await randomBeacon.dkgResultSubmissionTimeout()).to.be.equal(
            dkgResultSubmissionTimeout
          )
        })

        it("should update the DKG submitter precedence period", async () => {
          expect(
            await randomBeacon.dkgSubmitterPrecedencePeriodLength()
          ).to.be.equal(dkgSubmitterPrecedencePeriodLength)
        })

        it("should emit the DkgParametersUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "DkgParametersUpdated")
            .withArgs(
              dkgResultChallengePeriodLength,
              dkgResultSubmissionTimeout,
              dkgSubmitterPrecedencePeriodLength
            )
        })
      })

      context("when values are invalid", () => {
        context(
          "when precedence period length is equal submission timeout",
          () => {
            it("should revert", async () => {
              const invalidDkgSubmitterPrecedencePeriodLength =
                dkgResultSubmissionTimeout

              await expect(
                randomBeacon
                  .connect(governance)
                  .updateDkgParameters(
                    dkgResultChallengePeriodLength,
                    dkgResultSubmissionTimeout,
                    invalidDkgSubmitterPrecedencePeriodLength
                  )
              ).to.be.revertedWith(
                "New value should be less than result submission timeout"
              )
            })
          }
        )

        context(
          "when precedence period length is greater than submission timeout",
          () => {
            it("should revert", async () => {
              const invalidDkgSubmitterPrecedencePeriodLength =
                dkgResultSubmissionTimeout + 1

              await expect(
                randomBeacon
                  .connect(governance)
                  .updateDkgParameters(
                    dkgResultChallengePeriodLength,
                    dkgResultSubmissionTimeout,
                    invalidDkgSubmitterPrecedencePeriodLength
                  )
              ).to.be.revertedWith(
                "New value should be less than result submission timeout"
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
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await randomBeacon
          .connect(governance)
          .updateRewardParameters(
            newSortitionPoolRewardsBanDuration,
            newRelayEntryTimeoutNotificationRewardMultiplier,
            newUnauthorizedSigningNotificationRewardMultiplier,
            newDkgMaliciousResultNotificationRewardMultiplier
          )
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

      it("should emit the RewardParametersUpdated event", async () => {
        await expect(tx)
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
      let tx: ContractTransaction

      before(async () => {
        await createSnapshot()

        tx = await randomBeacon
          .connect(governance)
          .updateSlashingParameters(
            newRelayEntrySubmissionFailureSlashingAmount,
            newMaliciousDkgResultSlashingAmount,
            newUnauthorizedSigningSlashingAmount
          )
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

      it("should emit the SlashingParametersUpdated event", async () => {
        await expect(tx)
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

          tx = await randomBeacon
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

          await randomBeacon
            .connect(governance)
            .setRequesterAuthorization(thirdPartyContract.address, true)

          tx = await randomBeacon
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
