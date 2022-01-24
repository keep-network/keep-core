import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, helpers } = hre
  const { deployer, governance } = await getNamedAccounts()

  await helpers.ownable.transferOwnership(
    "WalletRegistry",
    governance,
    deployer
  )
}

export default func

func.tags = ["WalletRegistryTransferOwnership"]
func.dependencies = ["WalletRegistry"]
