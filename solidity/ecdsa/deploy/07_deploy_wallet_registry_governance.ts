import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const WalletRegistry = await deployments.get("WalletRegistry")

  const GOVERNANCE_DELAY = 604800 // 1 week

  const WalletRegistryGovernance = await deployments.deploy(
    "WalletRegistryGovernance",
    {
      from: deployer,
      args: [WalletRegistry.address, GOVERNANCE_DELAY],
      log: true,
    }
  )

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "WalletRegistryGovernance",
      address: WalletRegistryGovernance.address,
    })
  }
}

export default func

func.tags = ["WalletRegistryGovernance"]
func.dependencies = ["WalletRegistry"]
