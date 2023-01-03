import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()
  const { execute } = deployments

  const WalletRegistry = await deployments.get("WalletRegistry")

  // TODO: Who should execute the transaction - the deployer or the governance?
  // TODO: This step should only be executed when both conditions are true:
  //       1) we are deploying for mainnet
  //       2) we are in the chaosnet phase
  await execute(
    "RandomBeaconChaosnet",
    { from: deployer, log: true, waitConfirmations: 1 },
    "setRequesterAuthorization",
    WalletRegistry.address,
    true
  )
}

export default func

func.tags = ["WalletRegistryAuthorizeInBeaconChaosnet"]
func.dependencies = ["WalletRegistry"]
