// eslint-disable-next-line import/no-extraneous-dependencies
import { getNamedAccounts, getUnnamedAccounts, ethers } from "hardhat"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

// TODO: Move these utils to hardhat-helpers plugin.

export async function getNamedSigners(): Promise<{
  [name: string]: SignerWithAddress
}> {
  const namedSigners: { [name: string]: SignerWithAddress } = {}

  await Promise.all(
    Object.entries(await getNamedAccounts()).map(async ([name, address]) => {
      namedSigners[name] = await ethers.getSigner(address)
    })
  )

  return namedSigners
}

export async function getUnnamedSigners(): Promise<SignerWithAddress[]> {
  const unnamedSigners: SignerWithAddress[] = []

  await Promise.all(
    (
      await getUnnamedAccounts()
    ).map(async (address) => {
      unnamedSigners.push(await ethers.getSigner(address))
    })
  )

  return unnamedSigners
}
