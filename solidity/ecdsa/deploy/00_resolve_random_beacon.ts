import type {
  HardhatRuntimeEnvironment,
  HardhatNetworkConfig,
} from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments, helpers } = hre
  const { log } = deployments
  const { deployer } = await getNamedAccounts()

  const RandomBeacon = await deployments.getOrNull("RandomBeacon")

  if (RandomBeacon && helpers.address.isValid(RandomBeacon.address)) {
    log(`using external RandomBeacon at ${RandomBeacon.address}`)
  } else if (
    !hre.network.tags.allowStubs ||
    (hre.network.config as HardhatNetworkConfig)?.forking?.enabled
  ) {
    throw new Error("deployed RandomBeacon contract not found")
  } else {
    log("deploying RandomBeacon stub")

    await deployments.deploy("RandomBeacon", {
      contract: "RandomBeaconStub",
      from: deployer,
      log: true,
    })
  }
}

export default func

func.tags = ["RandomBeacon"]
