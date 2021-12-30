import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const SortitionPool = await deployments.get("SortitionPool")

  const DKGValidator = await deployments.deploy("DKGValidator", {
    from: deployer,
    args: [SortitionPool.address],
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "DKGValidator",
      address: DKGValidator.address,
    })
  }
}

export default func

func.tags = ["DKGValidator"]
func.dependencies = ["SortitionPool"]
