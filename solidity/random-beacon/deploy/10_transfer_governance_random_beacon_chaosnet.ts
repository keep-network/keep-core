import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, helpers } = hre
  const { deployer, chaosnetOwner } = await getNamedAccounts()

  await helpers.ownable.transferOwnership(
    "RandomBeaconChaosnet",
    chaosnetOwner,
    deployer
  )
}

export default func

func.tags = ["RandomBeaconChaosnetTransferGovernance"]
func.dependencies = ["RandomBeaconChaosnet"]

// Only execute for chaosnet deployments.
func.skip = async (hre: HardhatRuntimeEnvironment): Promise<boolean> =>
  !hre.network.tags.chaosnet
