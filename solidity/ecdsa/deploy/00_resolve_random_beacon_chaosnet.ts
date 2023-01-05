import type {
  HardhatRuntimeEnvironment,
  HardhatNetworkConfig,
} from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, helpers } = hre
  const { log } = deployments

  const RandomBeaconChaosnet = await deployments.getOrNull(
    "RandomBeaconChaosnet"
  )

  const isRandomBeaconChaosnetNeeded = function () {
    if (hre.network.tags.chaosnet) {
      return true
    }
    if (!hre.network.tags.allowStubs) {
      return true
    }
    return (hre.network.config as HardhatNetworkConfig)?.forking?.enabled
  }

  if (
    RandomBeaconChaosnet &&
    helpers.address.isValid(RandomBeaconChaosnet.address)
  ) {
    log(
      `using existing RandomBeaconChaosnet at ${RandomBeaconChaosnet.address}`
    )
  } else if (isRandomBeaconChaosnetNeeded()) {
    throw new Error("deployed RandomBeaconChaosnet contract not found")
  }
}

export default func

func.tags = ["RandomBeaconChaosnet"]
