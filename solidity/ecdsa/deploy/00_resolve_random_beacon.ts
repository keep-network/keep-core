import type {
  HardhatRuntimeEnvironment,
  HardhatNetworkConfig,
} from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, helpers } = hre
  const { log } = deployments

  const RandomBeacon = await deployments.getOrNull("RandomBeacon")

  const isRandomBeaconNeed = function () {
    if (hre.network.tags.chaosnet) {
      return false
    }
    if (!hre.network.tags.allowStubs) {
      return true
    }
    return (hre.network.config as HardhatNetworkConfig)?.forking?.enabled
  }

  if (RandomBeacon && helpers.address.isValid(RandomBeacon.address)) {
    log(`using existing RandomBeacon at ${RandomBeacon.address}`)
  } else if (isRandomBeaconNeed()) {
    throw new Error("deployed RandomBeacon contract not found")
  }
  // We don't deploy a stub of the RandomBeacon contract as unit tests mock
  // the IRandomBeacon.
}

export default func

func.tags = ["RandomBeacon"]
