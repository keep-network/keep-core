import { Contract } from "ethers"
import { ethers } from "hardhat"

import type { RandomBeacon, RandomBeaconGovernance } from "../../typechain"

export const constants = {
  groupSize: 3,
  signatureThreshold: 2,
  timeDKG: 13,
  dkgSubmissionEligibilityDelay: 10,
}

interface DeployedContracts {
  [key: string]: Contract
}

export async function randomBeaconDeployment(): Promise<DeployedContracts> {
  const DKG = await ethers.getContractFactory("DKG")
  const dkg = await DKG.deploy()

  const RandomBeacon = await ethers.getContractFactory("RandomBeacon", {
    libraries: {
      DKG: dkg.address,
    },
  })

  const randomBeacon: RandomBeacon = await RandomBeacon.deploy(
    constants.groupSize,
    constants.signatureThreshold,
    constants.timeDKG
  )

  await randomBeacon.deployed()

  const contracts: DeployedContracts = { randomBeacon }

  return contracts
}

export async function testDeployment(): Promise<DeployedContracts> {
  const contracts = await randomBeaconDeployment()

  const RandomBeaconGovernance = await ethers.getContractFactory(
    "RandomBeaconGovernance"
  )
  const randomBeaconGovernance: RandomBeaconGovernance = await RandomBeaconGovernance.deploy(
    contracts.randomBeacon.address
  )
  await randomBeaconGovernance.deployed()
  await contracts.randomBeacon.transferOwnership(randomBeaconGovernance.address)

  const newContracts = { randomBeaconGovernance }

  return { ...contracts, ...newContracts }
}
