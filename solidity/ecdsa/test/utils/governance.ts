import { ethers } from "hardhat"

import type { WalletRegistry, WalletRegistryGovernance } from "../../typechain"

// eslint-disable-next-line import/prefer-default-export
export async function upgradeRandomBeacon(
  walletRegistry: WalletRegistry,
  newRandomBeaconAddress: string
): Promise<void> {
  const walletRegistryGovernance: WalletRegistryGovernance =
    await ethers.getContractAt(
      "WalletRegistryGovernance",
      await walletRegistry.governance()
    )

  const governance = await ethers.getNamedSigner("governance")

  await walletRegistryGovernance
    .connect(governance)
    .upgradeRandomBeacon(newRandomBeaconAddress)
}
