import { ethers } from "hardhat"
import { smock } from "@defi-wonderland/smock"
import chai from "chai"

import type { BigNumber } from "ethers"
import type { WalletRegistry, IRandomBeacon } from "../../typechain"
import type { FakeContract } from "@defi-wonderland/smock"

chai.use(smock.matchers)

export async function fakeRandomBeacon(
  walletRegistry: WalletRegistry
): Promise<FakeContract<IRandomBeacon>> {
  const randomBeacon = await smock.fake<IRandomBeacon>("IRandomBeacon", {
    address: await walletRegistry.callStatic.randomBeacon(),
  })

  await (
    await ethers.getSigners()
  )[0].sendTransaction({
    to: randomBeacon.address,
    value: ethers.utils.parseEther("1000"),
  })

  return randomBeacon
}

export function resetMock(randomBeacon: FakeContract<IRandomBeacon>): void {
  randomBeacon.requestRelayEntry.reset()
}

export async function submitRelayEntry(
  walletRegistry: WalletRegistry,
  randomBeacon?: FakeContract<IRandomBeacon>
): Promise<{
  startBlock: number
  dkgSeed: BigNumber
}> {
  if (!randomBeacon) {
    // eslint-disable-next-line no-param-reassign
    randomBeacon = await fakeRandomBeacon(walletRegistry)
  }

  const relayEntry: BigNumber = ethers.BigNumber.from(
    ethers.utils.randomBytes(32)
  )

  // eslint-disable-next-line no-underscore-dangle
  const tx = await walletRegistry
    .connect(randomBeacon.wallet)
    .__beaconCallback(relayEntry, 0)

  return {
    startBlock: tx.blockNumber,
    dkgSeed: relayEntry,
  }
}
