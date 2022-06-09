import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()
  const { execute } = deployments

  const RandomBeacon = await deployments.get("RandomBeacon")

  await execute(
    "TokenStaking",
    { from: deployer },
    "approveApplication",
    RandomBeacon.address
  )
}

export default func

func.tags = ["RandomBeaconApprove"]
func.dependencies = ["TokenStaking", "RandomBeacon"]
