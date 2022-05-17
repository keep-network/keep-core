/* eslint-disable no-console */
import { task } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("initialize-wallet-owner", "Initializes Wallet Owner for Wallet Registry")
  .addParam("walletOwnerAddress", "The Wallet Owner's address")
  .setAction(async (args, hre) => {
    const { walletOwnerAddress } = args

    await initializeWalletOwner(hre, walletOwnerAddress)
  })

async function initializeWalletOwner(
  hre: HardhatRuntimeEnvironment,
  walletOwnerAddress: string
): Promise<void> {
  const { getNamedAccounts, ethers, deployments } = hre
  const { governance } = await getNamedAccounts()

  if (!ethers.utils.isAddress(walletOwnerAddress)) {
    throw Error(`invalid address: ${walletOwnerAddress}`)
  }

  const tx = await deployments.execute(
    "WalletRegistryGovernance",
    { from: governance },
    "initializeWalletOwner",
    walletOwnerAddress
  )

  console.log(
    `Initialized Wallet Owner address: ${walletOwnerAddress} in transaction: ${tx.transactionHash}`
  )
}

export default initializeWalletOwner
