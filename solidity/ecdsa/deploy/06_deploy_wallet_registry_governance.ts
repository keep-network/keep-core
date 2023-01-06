import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const WalletRegistry = await deployments.get("WalletRegistry")

  // 1 week for mainnet, 1 second for other test and dev networks.
  const GOVERNANCE_DELAY = hre.network.name === "mainnet" ? 604800 : 1

  const WalletRegistryGovernance = await deployments.deploy(
    "WalletRegistryGovernance",
    {
      from: deployer,
      args: [WalletRegistry.address, GOVERNANCE_DELAY],
      log: true,
      waitConfirmations: 1,
    }
  )

  if (hre.network.tags.etherscan) {
    await helpers.etherscan.verify(WalletRegistryGovernance)
  }

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
