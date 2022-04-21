import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()
  const { execute } = deployments

  const WalletRegistry = await deployments.get("WalletRegistry")

  await execute(
    "ReimbursementPool",
    { from: deployer },
    "authorize",
    WalletRegistry.address
  )
}

export default func

func.tags = ["WalletRegistryAuthorize"]
func.dependencies = ["ReimbursementPool", "WalletRegistry"]
