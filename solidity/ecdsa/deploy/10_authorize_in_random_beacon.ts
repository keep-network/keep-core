import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { governance } = await getNamedAccounts()
  const { execute } = deployments

  const WalletRegistry = await deployments.get("WalletRegistry")

  await execute(
    "RandomBeaconGovernance",
    { from: governance },
    "setRequesterAuthorization",
    WalletRegistry.address,
    true
  )
}

export default func

func.tags = ["WalletRegistryAuthorizeInBeacon"]
func.dependencies = ["RandomBeaconGovernance", "WalletRegistry"]
// When TEST_USE_STUBS is set we deploy a RandomBeaconStub for unit tests, hence
// there is no need for authorization.
func.skip = async function (hre: HardhatRuntimeEnvironment): Promise<boolean> {
  return process.env.TEST_USE_STUBS === "true"
}
