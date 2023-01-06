import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { governance } = await getNamedAccounts()
  const { execute } = deployments

  const RandomBeacon = await deployments.get("RandomBeacon")

  await execute(
    "WalletRegistryGovernance",
    { from: governance, log: true, waitConfirmations: 1 },
    "upgradeRandomBeacon",
    RandomBeacon.address
  )
}

export default func

func.tags = ["UpgradeRandomBeacon"]
func.dependencies = ["RandomBeacon", "WalletRegistryGovernance"]

// Skip for chaosnet deployments.
func.skip = async (hre: HardhatRuntimeEnvironment): Promise<boolean> =>
  hre.network.tags.chaosnet
