import { ethers } from "hardhat"
import { Signer, Contract, BigNumber } from "ethers"
import { expect } from "chai"
import { increaseTime } from "./helpers/contract-test-helpers"

import type { RandomBeaconParametersV2 } from "../typechain"

describe("RandomBeaconParametersV2", () => {
  let governance: Signer
  let thirdParty: Signer
  let randomBeaconParameters: RandomBeaconParametersV2

  beforeEach(async () => {
    const signers = await ethers.getSigners()
    governance = signers[0]
    thirdParty = signers[1]

    const RandomBeaconParametersV2 = await ethers.getContractFactory(
      "RandomBeaconParametersV2"
    )
    randomBeaconParameters = await RandomBeaconParametersV2.deploy()
    await randomBeaconParameters.deployed()
  })

  const parameter = "relayRequestFee"
  const newValue = 123
  const emptyValue = ethers.utils.formatBytes32String("")

  describe("beginRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters
            .connect(thirdParty)
            .beginUpdate(parameter, newValue)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      let tx

      beforeEach(async () => {
        tx = await randomBeaconParameters
          .connect(governance)
          .beginUpdate(parameter, newValue)
      })

      it("should not update the relay request fee", async () => {
        const value = await randomBeaconParameters.getParameter(parameter)
        console.log(parameter, value)
        expect(value).to.be.equal(emptyValue)
      })

      it("should emit the RelayRequestFeeUpdateStarted event", async () => {
        await expect(tx)
          .to.emit(randomBeaconParameters, "UpdateStarted")
          .withArgs(parameter, newValue)
      })
    })
  })

  describe("finalizeRelayRequestFeeUpdate", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters.connect(thirdParty).finalizeUpdate(parameter)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the update process is not initialized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconParameters.connect(governance).finalizeUpdate(parameter)
        ).to.be.revertedWith("GovernanceUpdateNotInitiated()")
      })
    })

    context("when the governance delay has not passed", () => {
      it("should revert", async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginUpdate(parameter, newValue)

        await increaseTime(23 * 60 * 60) // 23 hours

        await expect(
          randomBeaconParameters.connect(governance).finalizeUpdate(parameter)
        ).to.be.revertedWith("GovernanceDelayNotElapsed()")
      })
    })

    context("when the update process is initialized", () => {
      let tx

      beforeEach(async () => {
        await randomBeaconParameters
          .connect(governance)
          .beginUpdate(parameter, newValue)

        expect(
          await randomBeaconParameters.getParameterNewValue(parameter)
        ).to.be.equal(newValue)

        await increaseTime(24 * 60 * 60) // 24 hours

        tx = await randomBeaconParameters
          .connect(governance)
          .finalizeUpdate(parameter)
      })

      it("should update the relay request fee", async () => {
        expect(
          await randomBeaconParameters.getParameter(parameter)
        ).to.be.equal(newValue)
      })
    })
  })
})
