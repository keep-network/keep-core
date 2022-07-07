import { HardhatRuntimeEnvironment, HardhatNetworkConfig } from "hardhat/types"
import { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments, helpers } = hre
  const { log } = deployments
  const { deployer } = await getNamedAccounts()

  const RandomBeacon = await deployments.getOrNull("RandomBeacon")

  if (RandomBeacon && helpers.address.isValid(RandomBeacon.address)) {
    log(`using existing RandomBeacon at ${RandomBeacon.address}`)

    // Save deployment artifact of external contract to include it in the package.
    await deployments.save("RandomBeacon", RandomBeacon)
  } else if (
    !hre.network.tags.allowStubs ||
    (hre.network.config as HardhatNetworkConfig)?.forking?.enabled
  ) {
    throw new Error("deployed RandomBeacon contract not found")
  }
  // We don't deploy a stub of the RandomBeacon contract as unit tests mock
  // the IRandomBeacon.
}

export default func

func.tags = ["RandomBeacon"]
