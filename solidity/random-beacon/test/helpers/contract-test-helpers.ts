import { ethers } from "hardhat"

export async function lastBlockTime(): Promise<number> {
  return (await ethers.provider.getBlock("latest")).timestamp
}

export async function increaseTime(time: number): Promise<void> {
  const now = await lastBlockTime()
  await ethers.provider.send("evm_setNextBlockTimestamp", [now + time])
  await ethers.provider.send("evm_mine", [])
}
