import { ethers } from "hardhat"
import { Signer, Contract } from "ethers"
import { expect } from "chai"
import { increaseTime } from "./helpers/contract-test-helpers"

describe("RandomBeaconGovernance", () => {
  let governance: Signer
  let thirdParty: Signer
  let randomBeacon: Contract
  let randomBeaconGovernance: Contract

  beforeEach(async () => {
    const signers = await ethers.getSigners()
    governance = signers[0]
    thirdParty = signers[1]

    const RandomBeacon = await ethers.getContractFactory("RandomBeacon")
    randomBeacon = await RandomBeacon.deploy()
    await randomBeacon.deployed()

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
        expect(await randomBeacon.relayRequestFee()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayRequestFeeUpdateTime()
        ).to.be.equal(24 * 60 * 60) // 24 hours
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

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayRequestFeeUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

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
    })
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
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultSubmissionRewardUpdateTime()
        ).to.be.equal(24 * 60 * 60) // 24 hours
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

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(123)

        await increaseTime(24 * 60 * 60)

        tx = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgResultSubmissionRewardUpdate()
      })

      it("should update the dkg result submission reward", async () => {
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(123)
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
    })
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
        expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingSortitionPoolUnlockingRewardUpdateTime()
        ).to.be.equal(24 * 60 * 60) // 24 hours
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

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

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
    })
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
        ).to.be.equal(0)
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

        await increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(123)

        await increaseTime(14 * 24 * 60 * 60) // 2 weeks

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
    })
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
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.equal(24 * 60 * 60) // 24 hours
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

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(123)

        await increaseTime(24 * 60 * 60) // 24 hours

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
    })
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

    context("when the update value is less then required", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
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
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          0
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.equal(24 * 60 * 60) // 24 hours
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

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(11)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgResultChallengePeriodLengthUpdate()
      })

      it("should update the DKG result challenge period length", async () => {
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          11
        )
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
    })
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
          "Relay entry submission eligibility delay must be greater than 0 blocks"
        )
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
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingRelayEntrySubmissionEligibilityDelayUpdateTime()
        ).to.be.equal(24 * 60 * 60) // 24 hours
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

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(1)

        await increaseTime(24 * 60 * 60) // 24 hours

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
    })
  })

  describe("beginDkgSubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(thirdParty)
            .beginDkgSubmissionEligibilityDelayUpdate(1)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update value is less than zero", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
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
        tx = await randomBeaconGovernance
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(1)
      })

      it("should not update the DKG submission eligibility delay", async () => {
        expect(await randomBeacon.dkgSubmissionEligibilityDelay()).to.be.equal(
          0
        )
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingDkgSubmissionEligibilityDelayUpdateTime()
        ).to.be.equal(24 * 60 * 60)
      })

      it("should emit the DkgSubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
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
          randomBeaconGovernance
            .connect(thirdParty)
            .finalizeDkgSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(1)

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeDkgSubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginDkgSubmissionEligibilityDelayUpdate(1)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconGovernance
          .connect(governance)
          .finalizeDkgSubmissionEligibilityDelayUpdate()
      })

      it("should update the DKG submission eligibility delay", async () => {
        expect(await randomBeacon.dkgSubmissionEligibilityDelay()).to.be.equal(
          1
        )
      })

      it("should emit DkgSubmissionEligibilityDelayUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconGovernance,
            "DkgSubmissionEligibilityDelayUpdated"
          )
          .withArgs(1)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernance.getRemainingDkgSubmissionEligibilityDelayUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
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
        expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(0)
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

        await increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(123)

        await increaseTime(14 * 24 * 60 * 60) // 2 weeks

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
    })
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
        ).to.be.revertedWith(
          "Group creation frequency must be grater than zero"
        )
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
        expect(await randomBeacon.groupCreationFrequency()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernance.getRemainingGroupCreationFrequencyUpdateTime()
        ).to.be.equal(14 * 24 * 60 * 60) // 2 weeks
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

        await increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(1)

        await increaseTime(14 * 24 * 60 * 60) // 2 weeks

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
    })
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
            .beginGroupLifetimeUpdate(23 * 60 * 60) // 23 hours
        ).to.be.revertedWith("Group lifetime must be >= 1 day and <= 2 weeks")
      })
    })

    context("when the update value is more than 2 weeks", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginGroupLifetimeUpdate(15 * 24 * 60 * 60) // 15 days
        ).to.be.revertedWith("Group lifetime must be >= 1 day and <= 2 weeks")
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
        expect(await randomBeacon.groupLifetime()).to.be.equal(0)
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

        await increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginGroupLifetimeUpdate(2 * 24 * 60 * 60) // 2 days

        await increaseTime(14 * 24 * 60 * 60) // 2 weeks

        tx = await randomBeaconGovernance
          .connect(governance)
          .finalizeGroupLifetimeUpdate()
      })

      it("should update the group lifetime", async () => {
        expect(await randomBeacon.groupLifetime()).to.be.equal(2 * 24 * 60 * 60)
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
    })
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
        ).to.be.revertedWith("Callback gas limit must be > 0 and < 1000000")
      })
    })

    context("when the update value is million", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernance
            .connect(governance)
            .beginCallbackGasLimitUpdate(1000000)
        ).to.be.revertedWith("Callback gas limit must be > 0 and < 1000000")
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
        expect(await randomBeacon.callbackGasLimit()).to.be.equal(0)
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

        await increaseTime(13 * 24 * 60 * 60) // 13 days

        await expect(
          randomBeaconGovernance
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernance
          .connect(governance)
          .beginCallbackGasLimitUpdate(123)

        await increaseTime(14 * 24 * 60 * 60) // 2 weeks

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
    })
  })
})
