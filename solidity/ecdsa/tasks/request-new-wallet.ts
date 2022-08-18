import { task } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("request-new-wallet", "Requests for a new ECDSA wallet")
  .addParam("walletOwnerAddress", "The Wallet Owner's address")
  .setAction(async (args, hre) => {
    const { walletOwnerAddress } = args

    await requestNewWallet(hre, walletOwnerAddress)
  })

async function requestNewWallet(
  hre: HardhatRuntimeEnvironment,
  walletOwnerAddress: string
) {
  const { ethers, helpers } = hre

  if (!ethers.utils.isAddress(walletOwnerAddress)) {
    throw Error(`invalid address: ${walletOwnerAddress}`)
  }

  const walletOwner = await ethers.getSigner(walletOwnerAddress)

  const walletRegistry = await helpers.contracts.getContract("WalletRegistry")

  const tx = await walletRegistry.connect(walletOwner).requestNewWallet()
  await tx.wait()

  console.log("New ECDSA wallet was requested successfully")
}
