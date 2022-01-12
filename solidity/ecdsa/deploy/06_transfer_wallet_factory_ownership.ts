import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, helpers } = hre
  const { deployer, walletManager } = await getNamedAccounts()

  await helpers.ownable.transferOwnership(
    "WalletFactory",
    walletManager,
    deployer
  )
}

export default func

func.tags = ["WalletFactoryTransferOwnership"]
func.dependencies = ["WalletFactory"]
