/* eslint-disable import/prefer-default-export */

import { ethers } from "hardhat"
import type { BigNumber } from "ethers"

import { constants } from "../fixtures"
import type { Operator } from "./operators"
import type { SortitionPool } from "../../typechain"

const { keccak256, defaultAbiCoder } = ethers.utils

export async function selectGroup(
  sortitionPool: SortitionPool,
  seed: BigNumber
): Promise<Operator[]> {
  const identifiers = await sortitionPool.selectGroup(
    constants.groupSize,
    seed.toHexString()
  )
  const addresses = await sortitionPool.getIDOperators(identifiers)

  return identifiers.map((identifier, i) => ({
    id: identifier,
    address: addresses[i],
  }))
}

export function hashUint32Array(arrayToHash: number[]) {
  return keccak256(defaultAbiCoder.encode(["uint32[]"], [arrayToHash]))
}
