/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"
import type { Address } from "hardhat-deploy/types"
import { RandomBeacon } from "../../typechain"
import { constants } from "../fixtures"

export type OperatorID = number
export type Operator = { id: OperatorID; address: Address }

export async function registerOperators(
  randomBeacon: RandomBeacon,
  addresses: Address[]
): Promise<Operator[]> {
  const operators: Operator[] = []

  const sortitionPool = await ethers.getContractAt(
    "ISortitionPool",
    await randomBeacon.sortitionPool()
  )

  const staking = await ethers.getContractAt(
    "StakingStub",
    await randomBeacon.staking()
  )

  for (let i = 0; i < addresses.length; i++) {
    const address = addresses[i]

    await staking.setStake(address, constants.minimumStake)

    await randomBeacon
      .connect(await ethers.getSigner(address))
      .registerOperator()

    const id = await sortitionPool.getOperatorID(address)

    operators.push({ id, address })
  }

  return operators
}
