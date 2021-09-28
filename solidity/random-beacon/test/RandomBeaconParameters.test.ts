import { ethers } from "hardhat"
import { Signer, Contract } from "ethers"
import { expect } from "chai"
import { increaseTime } from "./helpers/contract-test-helpers"

describe("RandomBeaconParameters", () => {
  const UPDATED_VALUE = 123
  const GOVERNANCE_DELAY = 7 * 24 * 60 * 60 // 7 days in seconds

  let governance: Signer
  let thirdParty: Signer
  let randomBeaconParameters: Contract

  beforeEach(async () => {
    const signers = await ethers.getSigners()
    governance = signers[0]
    thirdParty = signers[1]

    const GovernableParameters = await ethers.getContractFactory(
      "GovernableParameters"
    )
    const governableParameters = await GovernableParameters.deploy()
    await governableParameters.deployed()

    const RandomBeaconParameters = await ethers.getContractFactory(
      "RandomBeaconParameters",
      {
        libraries: {
          GovernableParameters: governableParameters.address,
        },
      }
    )
    randomBeaconParameters = await RandomBeaconParameters.deploy()
    await randomBeaconParameters.deployed()
  })

  describe("beginRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayRequestFeeUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayRequestFeeUpdate(UPDATED_VALUE)
      })

      it("should not update the relay request fee", async () => {
        expect(await randomBeaconParameters.relayRequestFee()).to.be.equal(0)
      })

      it("should emit the RelayRequestFeeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayRequestFeeUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayRequestFeeUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayRequestFeeUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayRequestFeeUpdate()
      })

      it("should update the relay request fee", async () => {
        expect(await randomBeaconParameters.relayRequestFee()).to.be.equal(
          UPDATED_VALUE
        )
      })

      it("should emit RelayRequestFeeUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayRequestFeeUpdated")
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginRelayEntrySubmissionFailureSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)
      })

      it("should not update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeaconParameters.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(0)
      })

      it("should emit the RelayEntrySubmissionFailureSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "RelayEntrySubmissionFailureSlashingAmountUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntrySubmissionFailureSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
      })

      it("should update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeaconParameters.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit RelayEntrySubmissionFailureSlashingAmountUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "RelayEntrySubmissionFailureSlashingAmountUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginRelayEntrySubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)
      })

      it("should not update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeaconParameters.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(0)
      })

      it("should emit the RelayEntrySubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "RelayEntrySubmissionEligibilityDelayUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntrySubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
      })

      it("should update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeaconParameters.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit RelayEntrySubmissionEligibilityDelayUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "RelayEntrySubmissionEligibilityDelayUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)
      })

      it("should not update the relay entry hard timeout", async () => {
        expect(
          await randomBeaconParameters.relayEntryHardTimeout()
        ).to.be.equal(0)
      })

      it("should emit the RelayEntryHardTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayEntryHardTimeoutUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayEntryHardTimeoutUpdate()
      })

      it("should update the relay entry hard timeout", async () => {
        expect(
          await randomBeaconParameters.relayEntryHardTimeout()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit RelayEntryHardTimeoutUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayEntryHardTimeoutUpdated")
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginDkgResultSubmissionRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)
      })

      it("should not update the dkg result submission reward", async () => {
        expect(
          await randomBeaconParameters.dkgResultSubmissionReward()
        ).to.be.equal(0)
      })

      it("should emit the DkgResultSubmissionRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgResultSubmissionRewardUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultSubmissionRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeDkgResultSubmissionRewardUpdate()
      })

      it("should update the dkg result submission reward", async () => {
        expect(
          await randomBeaconParameters.dkgResultSubmissionReward()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit DkgResultSubmissionRewardUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "DkgResultSubmissionRewardUpdated")
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)
      })

      it("should not update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeaconParameters.maliciousDkgResultSlashingAmount()
        ).to.be.equal(0)
      })

      it("should emit the MaliciousDkgResultSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "MaliciousDkgResultSlashingAmountUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeMaliciousDkgResultSlashingAmountUpdate()
      })

      it("should update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeaconParameters.maliciousDkgResultSlashingAmount()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit MaliciousDkgResultSlashingAmountUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "MaliciousDkgResultSlashingAmountUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginDkgSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginDkgSubmissionEligibilityDelayUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(UPDATED_VALUE)
      })

      it("should not update the DKG submission eligibility delay", async () => {
        expect(
          await randomBeaconParameters.dkgSubmissionEligibilityDelay()
        ).to.be.equal(0)
      })

      it("should emit the DkgSubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgSubmissionEligibilityDelayUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeDkgSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeDkgSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeDkgSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeDkgSubmissionEligibilityDelayUpdate()
      })

      it("should update the DKG submission eligibility delay", async () => {
        expect(
          await randomBeaconParameters.dkgSubmissionEligibilityDelay()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit DkgSubmissionEligibilityDelayUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgSubmissionEligibilityDelayUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(
          await randomBeaconParameters.dkgResultChallengePeriodLength()
        ).to.be.equal(0)
      })

      it("should emit the DkgResultChallengePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgResultChallengePeriodLengthUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeDkgResultChallengePeriodLengthUpdate()
      })

      it("should update the DKG result challenge period length", async () => {
        expect(
          await randomBeaconParameters.dkgResultChallengePeriodLength()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit DkgResultChallengePeriodLengthUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgResultChallengePeriodLengthUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginSortitionPoolUnlockingRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)
      })

      it("should not update the sortition pool unlocking reward", async () => {
        expect(
          await randomBeaconParameters.sortitionPoolUnlockingReward()
        ).to.be.equal(0)
      })

      it("should emit the SortitionPoolUnlockingRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "SortitionPoolUnlockingRewardUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeSortitionPoolUnlockingRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeSortitionPoolUnlockingRewardUpdate()
      })

      it("should update the sortition pool unlocking reward", async () => {
        expect(
          await randomBeaconParameters.sortitionPoolUnlockingReward()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit SortitionPoolUnlockingRewardUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "SortitionPoolUnlockingRewardUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginGroupCreationFrequencyUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)
      })

      it("should not update the group creation frequency timeout", async () => {
        expect(
          await randomBeaconParameters.groupCreationFrequency()
        ).to.be.equal(0)
      })

      it("should emit the GroupCreationFrequencyUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "GroupCreationFrequencyUpdateStarted"
          )
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeGroupCreationFrequencyUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeGroupCreationFrequencyUpdate()
      })

      it("should update the group creation frequency", async () => {
        expect(
          await randomBeaconParameters.groupCreationFrequency()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit GroupCreationFrequencyUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "GroupCreationFrequencyUpdated")
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginGroupLifetimeUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginGroupLifetimeUpdate(UPDATED_VALUE)
      })

      it("should not update the group lifetime", async () => {
        expect(await randomBeaconParameters.groupLifetime()).to.be.equal(0)
      })

      it("should emit the GroupLifetimeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "GroupLifetimeUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginGroupLifetimeUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginGroupLifetimeUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeGroupLifetimeUpdate()
      })

      it("should update the group lifetime", async () => {
        expect(await randomBeaconParameters.groupLifetime()).to.be.equal(
          UPDATED_VALUE
        )
      })

      it("should emit GroupLifetimeUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "GroupLifetimeUpdated")
          .withArgs(UPDATED_VALUE)
      })
    })
  })

  describe("beginCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginCallbackGasLimitUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginCallbackGasLimitUpdate(UPDATED_VALUE)
      })

      it("should not update the callback gas limit", async () => {
        expect(await randomBeaconParameters.callbackGasLimit()).to.be.equal(0)
      })

      it("should emit the CallbackGasLimitUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "CallbackGasLimitUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginCallbackGasLimitUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconParameters
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginCallbackGasLimitUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeCallbackGasLimitUpdate()
      })

      it("should update the callback gas limit", async () => {
        expect(await randomBeaconParameters.callbackGasLimit()).to.be.equal(
          UPDATED_VALUE
        )
      })

      it("should emit CallbackGasLimitUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "CallbackGasLimitUpdated")
          .withArgs(UPDATED_VALUE)
      })
    })
  })
})
