import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction, DeployOptions } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const deployOptions: DeployOptions = {
    from: deployer,
    log: true,
    waitConfirmations: 1,
  }

  const RandomBeaconChaosnet = await deployments.deploy(
    "RandomBeaconChaosnet",
    {
      ...deployOptions,
    }
  )

  if (hre.network.tags.etherscan) {
    await helpers.etherscan.verify(RandomBeaconChaosnet)
  }

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "RandomBeaconChaosnet",
      address: RandomBeaconChaosnet.address,
    })
  }
}

export default func

func.tags = ["RandomBeaconChaosnet"]

// Only execute for chaosnet deployments.
func.skip = async (hre: HardhatRuntimeEnvironment): Promise<boolean> =>
  !hre.network.tags.chaosnet
