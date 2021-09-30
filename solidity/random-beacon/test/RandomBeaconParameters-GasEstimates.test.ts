import { ethers } from "hardhat"
import { BigNumber } from "ethers"
import { expect } from "chai"

import type {
  RandomBeaconParameters,
  RandomBeaconParametersV2,
} from "../typechain"

describe("RandomBeacon Parameters Gas Estimates", () => {
  let parametersV1: RandomBeaconParameters
  let parametersV2: RandomBeaconParametersV2

  beforeEach(async () => {
    const GovernableParameters = await ethers.getContractFactory(
      "GovernableParameters"
    )
    const governableParameters = await GovernableParameters.deploy()
    await governableParameters.deployed()

    const RandomBeaconParametersV1 = await ethers.getContractFactory(
      "RandomBeaconParameters",
      {
        libraries: {
          GovernableParameters: governableParameters.address,
        },
      }
    )
    parametersV1 = await RandomBeaconParametersV1.deploy()

    await parametersV1.deployed()

    const RandomBeaconParametersV2 = await ethers.getContractFactory(
      "RandomBeaconParametersV2"
    )
    parametersV2 = await RandomBeaconParametersV2.deploy()
    await parametersV2.deployed()
  })

  const parameter = "relayRequestFee"
  describe("parameters V1", async () => {
    const previousEstimatedGasV1 = 23474
    const previousEstimatedGasV2 = 25201

    let lowestEstimatedGas: BigNumber
    let estimatedGasV1: BigNumber
    let estimatedGasV2: BigNumber

    beforeEach(async () => {
      estimatedGasV1 = await parametersV1.estimateGas.relayRequestFee()
      estimatedGasV2 = await parametersV2.estimateGas.getParameter(parameter)

      console.log("estimated gas for V1", estimatedGasV1.toNumber())
      console.log("estimated gas for V2", estimatedGasV2.toNumber())

      lowestEstimatedGas = estimatedGasV1.lt(estimatedGasV2)
        ? estimatedGasV1
        : estimatedGasV2
    })

    it("V1 should match previous estimate", async () => {
      expect(
        estimatedGasV1,
        "gas estimate for V1 is not lower than previous estimate"
      ).to.be.eq(previousEstimatedGasV1)
    })

    it("V2 should match previous estimate", async () => {
      expect(
        estimatedGasV2,
        "gas estimate for V2 is not lower than previous estimate"
      ).to.be.eq(previousEstimatedGasV2)
    })

    it("V1 should be lower than lowest estimate", async () => {
      expect(
        estimatedGasV1,
        "gas estimate for V1 is not lower than estimate for V2"
      ).to.be.lte(lowestEstimatedGas)
    })

    it("V2 should be lower than lowest estimate", async () => {
      expect(
        estimatedGasV2,
        "gas estimate for V2 is not lower than estimate for V1"
      ).to.be.lte(lowestEstimatedGas)
    })
  })
})
