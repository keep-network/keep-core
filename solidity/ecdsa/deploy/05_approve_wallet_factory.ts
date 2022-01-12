import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()
  const { execute } = deployments

  const WalletFactory = await deployments.get("WalletFactory")

  await execute(
    "TokenStaking",
    { from: deployer },
    "approveApplication",
    WalletFactory.address
  )
}

export default func

func.tags = ["WalletFactoryApprove"]
func.dependencies = ["TokenStaking", "WalletFactory"]
