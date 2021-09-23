import { ethers } from "hardhat"
import { Signer, Contract } from "ethers"
import { expect } from "chai"
import { increaseTime } from "./helpers/contract-test-helpers"

describe("RandomBeaconGovernable", () => {
  const UPDATED_VALUE = 123
  const GOVERNANCE_DELAY = 7 * 24 * 60 * 60 // 7 days in seconds

  let governance: Signer
  let thirdParty: Signer
  let randomBeaconGovernable: Contract

  beforeEach(async () => {
    const signers = await ethers.getSigners()
    governance = signers[0]
    thirdParty = signers[1]

    const RandomBeaconGovernable = await ethers.getContractFactory(
      "RandomBeaconGovernable"
    )
    randomBeaconGovernable = await RandomBeaconGovernable.deploy()
    await randomBeaconGovernable.deployed()
  })

  describe("beginRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginRelayRequestFeeUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginRelayRequestFeeUpdate(UPDATED_VALUE)
      })

      it("should not update the relay request fee", async () => {
        expect(await randomBeaconGovernable.relayRequestFee()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingRelayRequestFeeUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the RelayRequestFeeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernable, "RelayRequestFeeUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayRequestFeeUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayRequestFeeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayRequestFeeUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeRelayRequestFeeUpdate()
      })

      it("should update the relay request fee", async () => {
        expect(await randomBeaconGovernable.relayRequestFee()).to.be.equal(
          UPDATED_VALUE
        )
      })

      it("should emit RelayRequestFeeUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconGovernable, "RelayRequestFeeUpdated")
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingRelayRequestFeeUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginDkgResultSubmissionRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)
      })

      it("should not update the dkg result submission reward", async () => {
        expect(
          await randomBeaconGovernable.dkgResultSubmissionReward()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingDkgResultSubmissionRewardUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the DkgResultSubmissionRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeDkgResultSubmissionRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginDkgResultSubmissionRewardUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeDkgResultSubmissionRewardUpdate()
      })

      it("should update the dkg result submission reward", async () => {
        expect(
          await randomBeaconGovernable.dkgResultSubmissionReward()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit DkgResultSubmissionRewardUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconGovernable, "DkgResultSubmissionRewardUpdated")
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingDkgResultSubmissionRewardUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginSortitionPoolUnlockingRewardUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)
      })

      it("should not update the sortition pool unlocking reward", async () => {
        expect(
          await randomBeaconGovernable.sortitionPoolUnlockingReward()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingSortitionPoolUnlockingRewardUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the SortitionPoolUnlockingRewardUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeSortitionPoolUnlockingRewardUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginSortitionPoolUnlockingRewardUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeSortitionPoolUnlockingRewardUpdate()
      })

      it("should update the sortition pool unlocking reward", async () => {
        expect(
          await randomBeaconGovernable.sortitionPoolUnlockingReward()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit SortitionPoolUnlockingRewardUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
            "SortitionPoolUnlockingRewardUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingSortitionPoolUnlockingRewardUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginRelayEntrySubmissionFailureSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)
      })

      it("should not update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeaconGovernable.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the RelayEntrySubmissionFailureSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntrySubmissionFailureSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
      })

      it("should update the relay entry submission failure slashing amount", async () => {
        expect(
          await randomBeaconGovernable.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit RelayEntrySubmissionFailureSlashingAmountUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
            "RelayEntrySubmissionFailureSlashingAmountUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginMaliciousDkgResultSlashingAmountUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)
      })

      it("should not update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeaconGovernable.maliciousDkgResultSlashingAmount()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the MaliciousDkgResultSlashingAmountUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeMaliciousDkgResultSlashingAmountUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginMaliciousDkgResultSlashingAmountUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeMaliciousDkgResultSlashingAmountUpdate()
      })

      it("should update the malicious DKG result slashing amount", async () => {
        expect(
          await randomBeaconGovernable.maliciousDkgResultSlashingAmount()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit MaliciousDkgResultSlashingAmountUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
            "MaliciousDkgResultSlashingAmountUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginDkgResultChallengePeriodLengthUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)
      })

      it("should not update the DKG result challenge period length", async () => {
        expect(
          await randomBeaconGovernable.dkgResultChallengePeriodLength()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the DkgResultChallengePeriodLengthUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeDkgResultChallengePeriodLengthUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginDkgResultChallengePeriodLengthUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeDkgResultChallengePeriodLengthUpdate()
      })

      it("should update the DKG result challenge period length", async () => {
        expect(
          await randomBeaconGovernable.dkgResultChallengePeriodLength()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit DkgResultChallengePeriodLengthUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
            "DkgResultChallengePeriodLengthUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingDkgResultChallengePeriodLengthUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginRelayEntrySubmissionEligibilityDelayUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)
      })

      it("should not update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeaconGovernable.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingRelayEntrySubmissionEligibilityDelayUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the RelayEntrySubmissionEligibilityDelayUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntrySubmissionEligibilityDelayUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeRelayEntrySubmissionEligibilityDelayUpdate()
      })

      it("should update the relay entry submission eligibility delay", async () => {
        expect(
          await randomBeaconGovernable.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit RelayEntrySubmissionEligibilityDelayUpdated event", async () => {
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
            "RelayEntrySubmissionEligibilityDelayUpdated"
          )
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingRelayEntrySubmissionEligibilityDelayUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)
      })

      it("should not update the relay entry hard timeout", async () => {
        expect(
          await randomBeaconGovernable.relayEntryHardTimeout()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingRelayEntryHardTimeoutUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the RelayEntryHardTimeoutUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernable, "RelayEntryHardTimeoutUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeRelayEntryHardTimeoutUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeRelayEntryHardTimeoutUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginRelayEntryHardTimeoutUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeRelayEntryHardTimeoutUpdate()
      })

      it("should update the relay entry hard timeout", async () => {
        expect(
          await randomBeaconGovernable.relayEntryHardTimeout()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit RelayEntryHardTimeoutUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconGovernable, "RelayEntryHardTimeoutUpdated")
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingRelayEntryHardTimeoutUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginGroupCreationFrequencyUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)
      })

      it("should not update the group creation frequency timeout", async () => {
        expect(
          await randomBeaconGovernable.groupCreationFrequency()
        ).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingGroupCreationFrequencyUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the GroupCreationFrequencyUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(
            randomBeaconGovernable,
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
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeGroupCreationFrequencyUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginGroupCreationFrequencyUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeGroupCreationFrequencyUpdate()
      })

      it("should update the group creation frequency", async () => {
        expect(
          await randomBeaconGovernable.groupCreationFrequency()
        ).to.be.equal(UPDATED_VALUE)
      })

      it("should emit GroupCreationFrequencyUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconGovernable, "GroupCreationFrequencyUpdated")
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingGroupCreationFrequencyUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginGroupLifetimeUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginGroupLifetimeUpdate(UPDATED_VALUE)
      })

      it("should not update the group lifetime", async () => {
        expect(await randomBeaconGovernable.groupLifetime()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingGroupLifetimeUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the GroupLifetimeUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernable, "GroupLifetimeUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeGroupLifetimeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginGroupLifetimeUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeGroupLifetimeUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginGroupLifetimeUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeGroupLifetimeUpdate()
      })

      it("should update the group lifetime", async () => {
        expect(await randomBeaconGovernable.groupLifetime()).to.be.equal(
          UPDATED_VALUE
        )
      })

      it("should emit GroupLifetimeUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconGovernable, "GroupLifetimeUpdated")
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingGroupLifetimeUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })

  describe("beginCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .beginCallbackGasLimitUpdate(UPDATED_VALUE)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconGovernable
          .connect(governance)
          .beginCallbackGasLimitUpdate(UPDATED_VALUE)
      })

      it("should not update the callback gas limit", async () => {
        expect(await randomBeaconGovernable.callbackGasLimit()).to.be.equal(0)
      })

      it("should start the governance delay timer", async () => {
        expect(
          await randomBeaconGovernable.getRemainingCallbackGasLimitUpdateTime()
        ).to.be.equal(GOVERNANCE_DELAY)
      })

      it("should emit the CallbackGasLimitUpdateStarted event", async () => {
        const blockTimestamp = (await ethers.provider.getBlock(tx.blockNumber))
          .timestamp
        await expect(tx)
          .to.emit(randomBeaconGovernable, "CallbackGasLimitUpdateStarted")
          .withArgs(UPDATED_VALUE, blockTimestamp)
      })
    })
  })

  describe("finalizeCallbackGasLimitUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(thirdParty)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Change not initiated")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginCallbackGasLimitUpdate(UPDATED_VALUE)

        await increaseTime(6 * 24 * 60 * 60) // 6 days

        await expect(
          randomBeaconGovernable
            .connect(governance)
            .finalizeCallbackGasLimitUpdate()
        ).to.be.revertedWith("Governance delay has not elapsed")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconGovernable
          .connect(governance)
          .beginCallbackGasLimitUpdate(UPDATED_VALUE)

        await increaseTime(GOVERNANCE_DELAY)

        tx = await randomBeaconGovernable
          .connect(governance)
          .finalizeCallbackGasLimitUpdate()
      })

      it("should update the callback gas limit", async () => {
        expect(await randomBeaconGovernable.callbackGasLimit()).to.be.equal(
          UPDATED_VALUE
        )
      })

      it("should emit CallbackGasLimitUpdated event", async () => {
        await expect(tx)
          .to.emit(randomBeaconGovernable, "CallbackGasLimitUpdated")
          .withArgs(UPDATED_VALUE)
      })

      it("should reset the governance delay timer", async () => {
        await expect(
          randomBeaconGovernable.getRemainingCallbackGasLimitUpdateTime()
        ).to.be.revertedWith("Change not initiated")
      })
    })
  })
})
