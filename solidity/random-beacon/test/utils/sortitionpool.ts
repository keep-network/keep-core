/* eslint-disable no-await-in-loop */

import type { Address } from "hardhat-deploy/types"
import type { SortitionPoolStub } from "../../typechain/SortitionPoolStub"

export type OperatorID = number
export type Operator = { id: OperatorID; address: Address }

export async function registerOperators(
  sortitionPool: SortitionPoolStub,
  addresses: Address[]
): Promise<Operator[]> {
  const operators: Operator[] = []

  for (let i = 0; i < addresses.length; i++) {
    const address = addresses[i]

    await sortitionPool.insertOperator(address)
    const id = await sortitionPool.getOperatorID(address)

    operators.push({ id, address })
  }

  return operators
}
