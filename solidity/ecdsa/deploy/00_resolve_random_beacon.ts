import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments, helpers } = hre
  const { log } = deployments
  const { deployer } = await getNamedAccounts()

  // When TEST_USE_STUBS is set we deploy a RandomBeaconStub for unit tests.
  if (process.env.TEST_USE_STUBS === "true") {
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
