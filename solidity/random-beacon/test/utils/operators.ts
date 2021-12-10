/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"
import type { Address } from "hardhat-deploy/types"
import type { BigNumber } from "ethers"
import { constants } from "../fixtures"
import type { RandomBeacon } from "../../typechain"

export type OperatorID = number
export type Operator = { id: OperatorID; address: Address }

export async function registerOperators(
  randomBeacon: RandomBeacon,
  addresses: Address[],
  stakeAmount: BigNumber = constants.minimumStake
): Promise<Operator[]> {
  const operators: Operator[] = []

  const sortitionPool = await ethers.getContractAt(
    "SortitionPool",
    await randomBeacon.sortitionPool()
  )

  const staking = await ethers.getContractAt(
    "StakingStub",
    await randomBeacon.staking()
  )

  for (let i = 0; i < addresses.length; i++) {
    const address = addresses[i]

    await staking.stake(address, address, address, stakeAmount)
    await staking.increaseAuthorization(address, staking.address, stakeAmount)

    const id = await sortitionPool.getOperatorID(address)

    operators.push({ id, address })
  }

  return operators
}
