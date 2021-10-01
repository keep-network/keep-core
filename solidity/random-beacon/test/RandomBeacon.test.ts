import { ethers } from "hardhat"
import { Signer, Contract } from "ethers"
import { expect } from "chai"

describe("RandomBeacon", () => {
  let governance: Signer
  let thirdParty: Signer
  let randomBeacon: Contract

  beforeEach(async () => {
    const signers = await ethers.getSigners()
    governance = signers[0]
    thirdParty = signers[1]

    const RandomBeacon = await ethers.getContractFactory("RandomBeacon")
    randomBeacon = await RandomBeacon.deploy()
    await randomBeacon.deployed()
  })

  describe("updateRelayEntryParameters", () => {
    const RELAY_REQUEST_FEE = 100
    const RELAY_ENTRY_SUBMISSION_ELIGIBILITY_DELAY = 200
    const RELAY_ENTRY_HARD_TIMEOUT = 300
    const CALLBACK_GAS_LIMIT = 400
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateRelayEntryParameters(
              RELAY_REQUEST_FEE,
              RELAY_ENTRY_SUBMISSION_ELIGIBILITY_DELAY,
              RELAY_ENTRY_HARD_TIMEOUT,
              CALLBACK_GAS_LIMIT
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      beforeEach(async () => {
        await randomBeacon
          .connect(governance)
          .updateRelayEntryParameters(
            RELAY_REQUEST_FEE,
            RELAY_ENTRY_SUBMISSION_ELIGIBILITY_DELAY,
            RELAY_ENTRY_HARD_TIMEOUT,
            CALLBACK_GAS_LIMIT
          )
      })

      it("should update the relay request fee", async () => {
        expect(await randomBeacon.relayRequestFee()).to.be.equal(
          RELAY_REQUEST_FEE
        )
      })

      it("should update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(RELAY_ENTRY_SUBMISSION_ELIGIBILITY_DELAY)
      })

      it("should update the relay entry hard timeout", async () => {
        expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(
          RELAY_ENTRY_HARD_TIMEOUT
        )
      })

      it("should update the callback gas limit", async () => {
        expect(await randomBeacon.callbackGasLimit()).to.be.equal(
          CALLBACK_GAS_LIMIT
        )
      })
    })
  })

  describe("updateGroupCreationParameters", () => {
    const GROUP_CREATION_FREQUENCY = 100
    const GROUP_LIFETIME = 200
    const DKG_RESULT_CHALLENGE_PERIOD_LENGTH = 300
    const DKG_SUBMISSION_ELIGIBILITY_DELAY = 400
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateGroupCreationParameters(
              GROUP_CREATION_FREQUENCY,
              GROUP_LIFETIME,
              DKG_RESULT_CHALLENGE_PERIOD_LENGTH,
              DKG_SUBMISSION_ELIGIBILITY_DELAY
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      beforeEach(async () => {
        await randomBeacon
          .connect(governance)
          .updateGroupCreationParameters(
            GROUP_CREATION_FREQUENCY,
            GROUP_LIFETIME,
            DKG_RESULT_CHALLENGE_PERIOD_LENGTH,
            DKG_SUBMISSION_ELIGIBILITY_DELAY
          )
      })

      it("should update the group creation frequency", async () => {
        expect(await randomBeacon.groupCreationFrequency()).to.be.equal(
          GROUP_CREATION_FREQUENCY
        )
      })

      it("should update the group lifetime", async () => {
        expect(await randomBeacon.groupLifetime()).to.be.equal(GROUP_LIFETIME)
      })

      it("should update the DKG result challenge period length", async () => {
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          DKG_RESULT_CHALLENGE_PERIOD_LENGTH
        )
      })

      it("should update the DKG submission eligibility delay", async () => {
        expect(await randomBeacon.dkgSubmissionEligibilityDelay()).to.be.equal(
          DKG_SUBMISSION_ELIGIBILITY_DELAY
        )
      })
    })
  })

  describe("updateRewardParameters", () => {
    const DKG_RESULT_SUBMISSION_REWARD = 100
    const SORTITION_POOL_UNLOCKING_REWARD = 200
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateRewardParameters(
              DKG_RESULT_SUBMISSION_REWARD,
              SORTITION_POOL_UNLOCKING_REWARD
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      beforeEach(async () => {
        await randomBeacon
          .connect(governance)
          .updateRewardParameters(
            DKG_RESULT_SUBMISSION_REWARD,
            SORTITION_POOL_UNLOCKING_REWARD
          )
      })

      it("should update the DKG result submission reward", async () => {
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(
          DKG_RESULT_SUBMISSION_REWARD
        )
      })

      it("should update the sortition pool unlocking reward", async () => {
        expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(
          SORTITION_POOL_UNLOCKING_REWARD
        )
      })
    })
  })

  describe("updateSlashingParameters", () => {
    const RELAY_ENTRY_SUBMISSION_FAILURE_SLASHING_AMOUNT = 100
    const MALICIOUS_DKG_RESULT_SLASHING_AMOUNT = 200
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeacon
            .connect(thirdParty)
            .updateSlashingParameters(
              RELAY_ENTRY_SUBMISSION_FAILURE_SLASHING_AMOUNT,
              MALICIOUS_DKG_RESULT_SLASHING_AMOUNT
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      beforeEach(async () => {
        await randomBeacon
          .connect(governance)
          .updateSlashingParameters(
            RELAY_ENTRY_SUBMISSION_FAILURE_SLASHING_AMOUNT,
            MALICIOUS_DKG_RESULT_SLASHING_AMOUNT
          )
      })

      it("should update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeacon.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(RELAY_ENTRY_SUBMISSION_FAILURE_SLASHING_AMOUNT)
      })

      it("should update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeacon.maliciousDkgResultSlashingAmount()
        ).to.be.equal(MALICIOUS_DKG_RESULT_SLASHING_AMOUNT)
      })
    })
  })
})
