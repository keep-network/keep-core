import { ethers, waffle } from "hardhat"
import { expect } from "chai"

import { constants, randomBeaconDeployment } from "./helpers/fixtures"

import type { Signer } from "ethers"
import type { RandomBeacon, RandomBeacon__factory } from "../typechain"

describe("RandomBeacon", () => {
  context("when contracts not deployed", () => {
    describe("constructor", async function () {
      let RandomBeacon: RandomBeacon__factory

      beforeEach(async function () {
        const DKG = await ethers.getContractFactory("DKG")
        const dkg = await DKG.deploy()

        const Groups = await ethers.getContractFactory("Groups")
        const groups = await Groups.deploy()

        RandomBeacon = await ethers.getContractFactory("RandomBeacon", {
          libraries: {
            DKG: dkg.address,
            Groups: groups.address,
          },
        })
      })

      it("sets default governable properties", async () => {
        const randomBeacon = await RandomBeacon.deploy(
          constants.groupSize,
          constants.signatureThreshold,
          constants.timeDKG
        )

        expect(await randomBeacon.relayRequestFee()).to.be.equal(0)
        expect(
          await randomBeacon.relayEntrySubmissionEligibilityDelay()
        ).to.be.equal(10)
        expect(await randomBeacon.relayEntryHardTimeout()).to.be.equal(5760)
        expect(await randomBeacon.callbackGasLimit()).to.be.equal(200000)
        expect(await randomBeacon.groupCreationFrequency()).to.be.equal(10)
        expect(await randomBeacon.groupLifetime()).to.be.equal(1209600)
        expect(await randomBeacon.dkgResultChallengePeriodLength()).to.be.equal(
          1440
        )
        expect(
          await randomBeacon.dkgResultSubmissionEligibilityDelay()
        ).to.be.equal(10)
        expect(await randomBeacon.dkgResultSubmissionReward()).to.be.equal(0)
        expect(await randomBeacon.sortitionPoolUnlockingReward()).to.be.equal(0)
        expect(
          await randomBeacon.relayEntrySubmissionFailureSlashingAmount()
        ).to.be.equal(ethers.BigNumber.from(10).pow(18).mul(1000))
        expect(
          await randomBeacon.maliciousDkgResultSlashingAmount()
        ).to.be.equal(ethers.BigNumber.from(10).pow(18).mul(50000))
      })

      it("sets provided properties", async () => {
        const randomBeacon = await RandomBeacon.deploy(10, 5, 30)

        expect(await randomBeacon.GROUP_SIZE()).to.be.equal(10)
        expect(await randomBeacon.SIGNATURE_THRESHOLD()).to.be.equal(5)
        expect(await randomBeacon.TIME_DKG()).to.be.equal(30)
      })

      it("reverts on groupSize equal zero", async () => {
        await expect(
          RandomBeacon.deploy(
            0,
            constants.signatureThreshold,
            constants.timeDKG
          )
        ).to.be.revertedWith("groupSize has to be greater than zero")
      })

      it("reverts on signatureThreshold equal zero", async () => {
        await expect(
          RandomBeacon.deploy(constants.groupSize, 0, constants.timeDKG)
        ).to.be.revertedWith("signatureThreshold has to be greater than zero")
      })

      it("reverts on timeDkg equal zero", async () => {
        await expect(
          RandomBeacon.deploy(
            constants.groupSize,
            constants.signatureThreshold,
            0
          )
        ).to.be.revertedWith("timeDkg has to be greater than zero")
      })

      it("reverts on signatureThreshold greater than groupSize", async () => {
        await expect(
          RandomBeacon.deploy(
            constants.groupSize,
            constants.groupSize + 1,
            constants.timeDKG
          )
        ).to.be.revertedWith(
          "signatureThreshold has to be less or equal groupSize"
        )
      })

      it("completes on signatureThreshold equal groupSize", async () => {
        const randomBeacon = await RandomBeacon.deploy(
          constants.groupSize,
          constants.groupSize,
          constants.timeDKG
        )

        expect(await randomBeacon.SIGNATURE_THRESHOLD()).to.be.equal(
          constants.groupSize
        )
      })
    })
  })

  context("when contracts deployed for test", () => {
    let governance: Signer
    let thirdParty: Signer
    let randomBeacon: RandomBeacon

    before(async function () {
      ;[governance, thirdParty] = await ethers.getSigners()
    })

    beforeEach("load test fixture", async function () {
      const contracts = await waffle.loadFixture(randomBeaconDeployment)

      randomBeacon = contracts.randomBeacon as RandomBeacon
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

    describe("updateGroupCreationParameters", () => {
      const groupCreationFrequency = 100
      const groupLifetime = 200
      const dkgResultChallengePeriodLength = 300
      const dkgResultSubmissionEligibilityDelay = 400

      context("when the caller is not the owner", () => {
        it("should revert", async () => {
          await expect(
            randomBeacon
              .connect(thirdParty)
              .updateGroupCreationParameters(
                groupCreationFrequency,
                groupLifetime,
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
            .updateGroupCreationParameters(
              groupCreationFrequency,
              groupLifetime,
              dkgResultChallengePeriodLength,
              dkgResultSubmissionEligibilityDelay
            )
        })

        it("should update the group creation frequency", async () => {
          expect(await randomBeacon.groupCreationFrequency()).to.be.equal(
            groupCreationFrequency
          )
        })

        it("should update the group lifetime", async () => {
          expect(await randomBeacon.groupLifetime()).to.be.equal(groupLifetime)
        })

        it("should update the DKG result challenge period length", async () => {
          expect(
            await randomBeacon.dkgResultChallengePeriodLength()
          ).to.be.equal(dkgResultChallengePeriodLength)
        })

        it("should update the DKG result submission eligibility delay", async () => {
          expect(
            await randomBeacon.dkgResultSubmissionEligibilityDelay()
          ).to.be.equal(dkgResultSubmissionEligibilityDelay)
        })

        it("should emit the GroupCreationParametersUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "GroupCreationParametersUpdated")
            .withArgs(
              groupCreationFrequency,
              groupLifetime,
              dkgResultChallengePeriodLength,
              dkgResultSubmissionEligibilityDelay
            )
        })
      })
    })

    describe("updateRewardParameters", () => {
      const dkgResultSubmissionReward = 100
      const sortitionPoolUnlockingReward = 200

      context("when the caller is not the owner", () => {
        it("should revert", async () => {
          await expect(
            randomBeacon
              .connect(thirdParty)
              .updateRewardParameters(
                dkgResultSubmissionReward,
                sortitionPoolUnlockingReward
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
              sortitionPoolUnlockingReward
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

        it("should emit the RewardParametersUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "RewardParametersUpdated")
            .withArgs(dkgResultSubmissionReward, sortitionPoolUnlockingReward)
        })
      })
    })

    describe("updateSlashingParameters", () => {
      const relayEntrySubmissionFailureSlashingAmount = 100
      const maliciousDkgResultSlashingAmount = 200

      context("when the caller is not the owner", () => {
        it("should revert", async () => {
          await expect(
            randomBeacon
              .connect(thirdParty)
              .updateSlashingParameters(
                relayEntrySubmissionFailureSlashingAmount,
                maliciousDkgResultSlashingAmount
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
              maliciousDkgResultSlashingAmount
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

        it("should emit the SlashingParametersUpdated event", async () => {
          await expect(tx)
            .to.emit(randomBeacon, "SlashingParametersUpdated")
            .withArgs(
              relayEntrySubmissionFailureSlashingAmount,
              maliciousDkgResultSlashingAmount
            )
        })
      })
    })
  })
})
