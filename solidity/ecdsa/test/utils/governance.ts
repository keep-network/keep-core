import { ethers, helpers } from "hardhat"

import { constants } from "../fixtures"

import type { WalletRegistry, WalletRegistryGovernance } from "../../typechain"

// eslint-disable-next-line import/prefer-default-export
export async function updateRandomBeacon(
  walletRegistry: WalletRegistry,
  newRandomBeaconAddress: string
): Promise<void> {
  const walletRegistryGovernance: WalletRegistryGovernance =
    await ethers.getContractAt(
      "WalletRegistryGovernance",
      await walletRegistry.owner()
    )

  const governance = await ethers.getNamedSigner("governance")

  await walletRegistryGovernance
    .connect(governance)
    .beginRandomBeaconUpdate(newRandomBeaconAddress)

  await helpers.time.increaseTime(constants.governanceDelayCritical)

  await walletRegistryGovernance
    .connect(governance)
    .finalizeRandomBeaconUpdate()
}
