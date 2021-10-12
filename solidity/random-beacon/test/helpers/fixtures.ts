import { Contract } from "ethers"
import { ethers } from "hardhat"

import type { RandomBeacon, RandomBeaconGovernance } from "../../typechain"

export const constants = {
  groupSize: 3,
  signatureThreshold: 2,
  timeDKG: 13,
}

export const params = {
  relayRequestFee: 0,
  relayEntrySubmissionEligibilityDelay: 10,
  relayEntryHardTimeout: 5760,
  callbackGasLimit: 200000,
  groupCreationFrequency: 10,
  groupLifeTime: 60 * 60 * 24 * 14,
  dkgResultChallengePeriodLength: 1440,
  dkgResultSubmissionEligibilityDelay: 10,
  dkgResultSubmissionReward: 0,
  sortitionPoolUnlockingReward: 0,
  relayEntrySubmissionFailureSlashingAmount: ethers.BigNumber.from(10)
    .pow(18)
    .mul(1000),
  maliciousDkgResultSlashingAmount: ethers.BigNumber.from(10)
    .pow(18)
    .mul(50000),
}

interface DeployedContracts {
  [key: string]: Contract
}

export async function randomBeaconDeployment(): Promise<DeployedContracts> {
  const DKG = await ethers.getContractFactory("DKG")
  const dkg = await DKG.deploy()

  const Groups = await ethers.getContractFactory("Groups")
  const groups = await Groups.deploy()

  const RandomBeacon = await ethers.getContractFactory("RandomBeacon", {
    libraries: {
      DKG: dkg.address,
      Groups: groups.address,
    },
  })

  const randomBeacon: RandomBeacon = await RandomBeacon.deploy(
    constants.groupSize,
    constants.signatureThreshold,
    constants.timeDKG
  )

  await randomBeacon.deployed()

  const contracts: DeployedContracts = { dkg, groups, randomBeacon }

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
