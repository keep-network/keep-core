import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer, thresholdCouncil } = await getNamedAccounts()

  await deployments.execute(
    "RandomBeaconChaosnet",
    { from: deployer, log: true, waitConfirmations: 1 },
    "transferGovernance",
    thresholdCouncil
  )
}

export default func

func.tags = ["RandomBeaconChaosnetTransferGovernance"]
func.dependencies = ["RandomBeaconChaosnet"]
