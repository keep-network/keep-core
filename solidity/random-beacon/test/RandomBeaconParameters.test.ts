import { ethers } from "hardhat"
import { Signer, Contract } from "ethers"
import { expect } from "chai"
import { increaseTime } from "./helpers/contract-test-helpers"

describe("RandomBeaconParameters", () => {
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
            .beginRelayRequestFeeUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayRequestFeeUpdate(123)
      })

      it("should not update the relay request fee", async () => {
        expect(await randomBeaconParameters.relayRequestFee()).to.be.equal(0)
      })

      it("should emit the RelayRequestFeeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayRequestFeeUpdateStarted")
          .withArgs(123, blockTimestamp)
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
          .beginRelayRequestFeeUpdate(123)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginRelayRequestFeeUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayRequestFeeUpdate()
      })

      it("should update the relay request fee", async () => {
        expect(await randomBeaconParameters.relayRequestFee()).to.be.equal(123)
      })

      it("should emit RelayRequestFeeUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayRequestFeeUpdated")
          .withArgs(123)
      })
    })
  })

  describe("beginRelayEntrySubmissionFailureSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)
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
          .withArgs(123, blockTimestamp)
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
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

        await increaseTime(13 * 24 * 60 * 60) // 13 days

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
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

        await increaseTime(14 * 24 * 60 * 60) // 14 days

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
      })

      it("should update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeaconParameters.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(123)
      })

      it("should emit RelayEntrySubmissionFailureSlashingAmountUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "RelayEntrySubmissionFailureSlashingAmountUpdated"
          )
          .withArgs(123)
      })
    })
  })

  describe("beginRelayEntrySubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(0)
        ).to.be.revertedWith(
          "Relay entry submission eligibility delay must be greater than 0 blocks"
        )
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)
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
          .withArgs(1, blockTimestamp)
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
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
      })

      it("should update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeaconParameters.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(1)
      })

      it("should emit RelayEntrySubmissionEligibilityDelayUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "RelayEntrySubmissionEligibilityDelayUpdated"
          )
          .withArgs(1)
      })
    })
  })

  describe("beginRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginRelayEntryHardTimeoutUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(123)
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
          .withArgs(123, blockTimestamp)
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
          .beginRelayEntryHardTimeoutUpdate(123)

        await increaseTime(13 * 24 * 60 * 60) // 13 days

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
          .beginRelayEntryHardTimeoutUpdate(123)

        await increaseTime(14 * 24 * 60 * 60) // 14 days

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeRelayEntryHardTimeoutUpdate()
      })

      it("should update the relay entry hard timeout", async () => {
        expect(
          await randomBeaconParameters.relayEntryHardTimeout()
        ).to.be.equal(123)
      })

      it("should emit RelayEntryHardTimeoutUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "RelayEntryHardTimeoutUpdated")
          .withArgs(123)
      })
    })
  })

  describe("beginDkgResultSubmissionRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginDkgResultSubmissionRewardUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(123)
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
          .withArgs(123, blockTimestamp)
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
          .beginDkgResultSubmissionRewardUpdate(123)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginDkgResultSubmissionRewardUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeDkgResultSubmissionRewardUpdate()
      })

      it("should update the dkg result submission reward", async () => {
        expect(
          await randomBeaconParameters.dkgResultSubmissionReward()
        ).to.be.equal(123)
      })

      it("should emit DkgResultSubmissionRewardUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "DkgResultSubmissionRewardUpdated")
          .withArgs(123)
      })
    })
  })

  describe("beginMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginMaliciousDkgResultSlashingAmountUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)
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
          .withArgs(123, blockTimestamp)
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
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeMaliciousDkgResultSlashingAmountUpdate()
      })

      it("should update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeaconParameters.maliciousDkgResultSlashingAmount()
        ).to.be.equal(123)
      })

      it("should emit MaliciousDkgResultSlashingAmountUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "MaliciousDkgResultSlashingAmountUpdated"
          )
          .withArgs(123)
      })
    })
  })

  describe("beginDkgSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginDkgSubmissionEligibilityDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginDkgSubmissionEligibilityDelayUpdate(0)
        ).to.be.revertedWith(
          "DKG submission eligibility delay must be greater than 0 blocks"
        )
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(1)
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
          .withArgs(1, blockTimestamp)
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
          .beginDkgSubmissionEligibilityDelayUpdate(1)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginDkgSubmissionEligibilityDelayUpdate(1)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeDkgSubmissionEligibilityDelayUpdate()
      })

      it("should update the DKG submission eligibility delay", async () => {
        expect(
          await randomBeaconParameters.dkgSubmissionEligibilityDelay()
        ).to.be.equal(1)
      })

      it("should emit DkgSubmissionEligibilityDelayUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgSubmissionEligibilityDelayUpdated"
          )
          .withArgs(1)
      })
    })
  })

  describe("beginDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginDkgResultChallengePeriodLengthUpdate(11)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less then required", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginDkgResultChallengePeriodLengthUpdate(10)
        ).to.be.revertedWith(
          "DKG result challenge period length must be grater than 10 blocks"
        )
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)
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
          .withArgs(11, blockTimestamp)
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
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeDkgResultChallengePeriodLengthUpdate()
      })

      it("should update the DKG result challenge period length", async () => {
        expect(
          await randomBeaconParameters.dkgResultChallengePeriodLength()
        ).to.be.equal(11)
      })

      it("should emit DkgResultChallengePeriodLengthUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "DkgResultChallengePeriodLengthUpdated"
          )
          .withArgs(11)
      })
    })
  })

  describe("beginSortitionPoolUnlockingRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginSortitionPoolUnlockingRewardUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(123)
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
          .withArgs(123, blockTimestamp)
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
          .beginSortitionPoolUnlockingRewardUpdate(123)

        await increaseTime(23 * 60 * 60) // 23 hours

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
          .beginSortitionPoolUnlockingRewardUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeSortitionPoolUnlockingRewardUpdate()
      })

      it("should update the sortition pool unlocking reward", async () => {
        expect(
          await randomBeaconParameters.sortitionPoolUnlockingReward()
        ).to.be.equal(123)
      })

      it("should emit SortitionPoolUnlockingRewardUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconParameters,
            "SortitionPoolUnlockingRewardUpdated"
          )
          .withArgs(123)
      })
    })
  })

  describe("beginGroupCreationFrequencyUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginGroupCreationFrequencyUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginGroupCreationFrequencyUpdate(0)
        ).to.be.revertedWith(
          "Group creation frequency must be grater than zero"
        )
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)
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
          .withArgs(1, blockTimestamp)
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
          .beginGroupCreationFrequencyUpdate(1)

        await increaseTime(13 * 24 * 60 * 60) // 13 days

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
          .beginGroupCreationFrequencyUpdate(1)

        await increaseTime(14 * 24 * 60 * 60) // 14 days

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeGroupCreationFrequencyUpdate()
      })

      it("should update the group creation frequency", async () => {
        expect(
          await randomBeaconParameters.groupCreationFrequency()
        ).to.be.equal(1)
      })

      it("should emit GroupCreationFrequencyUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "GroupCreationFrequencyUpdated")
          .withArgs(1)
      })
    })
  })

  describe("beginGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than one day", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginGroupLifetimeUpdate(23 * 60 * 60) // 23 hours
        ).to.be.revertedWith("Group lifetime must be >= 1 day and <= 2 weeks")
      })
    })

    context("when the update value is more than 2 weeks", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginGroupLifetimeUpdate(15 * 24 * 60 * 60) // 15 days
        ).to.be.revertedWith("Group lifetime must be >= 1 day and <= 2 weeks")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days
      })

      it("should not update the group lifetime", async () => {
        expect(await randomBeaconParameters.groupLifetime()).to.be.equal(0)
      })

      it("should emit the GroupLifetimeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "GroupLifetimeUpdateStarted")
          .withArgs(2 * 24 * 60 * 60, blockTimestamp) // 2 days
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
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

        await increaseTime(13 * 24 * 60 * 60) // 13 days

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
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

        await increaseTime(14 * 24 * 60 * 60) // 14 days

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeGroupLifetimeUpdate()
      })

      it("should update the group lifetime", async () => {
        expect(await randomBeaconParameters.groupLifetime()).to.be.equal(
          2 * 24 * 60 * 60 // 2 days
        )
      })

      it("should emit GroupLifetimeUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "GroupLifetimeUpdated")
          .withArgs(2 * 24 * 60 * 60) // 2 days
      })
    })
  })

  describe("beginCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginCallbackGasLimitUpdate(123)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginCallbackGasLimitUpdate(0)
        ).to.be.revertedWith("Callback gas limit must be > 0 and < 1000000")
      })
    })

    context("when the update value is million", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(governance)
            .beginCallbackGasLimitUpdate(1000000)
        ).to.be.revertedWith("Callback gas limit must be > 0 and < 1000000")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginCallbackGasLimitUpdate(123)
      })

      it("should not update the callback gas limit", async () => {
        expect(await randomBeaconParameters.callbackGasLimit()).to.be.equal(0)
      })

      it("should emit the CallbackGasLimitUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconParameters, "CallbackGasLimitUpdateStarted")
          .withArgs(123, blockTimestamp)
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
          .beginCallbackGasLimitUpdate(123)

        await increaseTime(13 * 24 * 60 * 60) // 13 days

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
          .beginCallbackGasLimitUpdate(123)

        await increaseTime(14 * 24 * 60 * 60) // 14 days

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeCallbackGasLimitUpdate()
      })

      it("should update the callback gas limit", async () => {
        expect(await randomBeaconParameters.callbackGasLimit()).to.be.equal(123)
      })

      it("should emit CallbackGasLimitUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "CallbackGasLimitUpdated")
          .withArgs(123)
      })
    })
  })
})
