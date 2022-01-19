import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, helpers } = hre
  const { deployer, walletRegistry } = await getNamedAccounts()

  await helpers.ownable.transferOwnership(
    "WalletFactory",
    walletRegistry,
    deployer
  )
}

export default func

func.tags = ["WalletFactoryTransferOwnership"]
func.dependencies = ["WalletFactory"]
