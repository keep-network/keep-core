/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"

import { params } from "../fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Address } from "hardhat-deploy/types"
import type { BigNumber } from "ethers"
import type { RandomBeacon } from "../../typechain"

export type OperatorID = number
export type Operator = { id: OperatorID; signer: SignerWithAddress }

export async function registerOperators(
  randomBeacon: RandomBeacon,
  addresses: Address[],
  stakeAmount: BigNumber = params.minimumAuthorization
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

    await staking.setStake(address, stakeAmount)

    await randomBeacon
      .connect(await ethers.getSigner(address))
      .registerOperator()

    const id = await sortitionPool.getOperatorID(address)

    operators.push({ id, signer: operator })
  }

  return operators
}
