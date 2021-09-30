// TODO: This helper is placed in this repo just temporarily. It is expected to
// be extracted to a separate common repo.
import type { BigNumber } from "ethers"
import { ethers } from "hardhat"

export async function lastBlockTime(): Promise<number> {
  return (await ethers.provider.getBlock("latest")).timestamp
}

export async function increaseTime(time: number | BigNumber): Promise<number> {
  const lastBlock = await lastBlockTime()

  const expectedTime = ethers.BigNumber.from(lastBlock).add(time)

  await ethers.provider.send("evm_setNextBlockTimestamp", [
    expectedTime.toNumber(),
  ])
  await ethers.provider.send("evm_mine", [])

  return expectedTime.toNumber()
}
