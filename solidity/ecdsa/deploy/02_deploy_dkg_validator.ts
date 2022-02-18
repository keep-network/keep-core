import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const SortitionPool = await deployments.get("SortitionPool")

  const EcdsaDkgValidator = await deployments.deploy("EcdsaDkgValidator", {
    from: deployer,
    args: [SortitionPool.address],
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "EcdsaDkgValidator",
      address: EcdsaDkgValidator.address,
    })
  }
}

export default func

func.tags = ["EcdsaDkgValidator"]
func.dependencies = ["SortitionPool"]
