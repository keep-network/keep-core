import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const BeaconSortitionPool = await deployments.get("BeaconSortitionPool")

  const BeaconDkgValidator = await deployments.deploy("BeaconDkgValidator", {
    from: deployer,
    args: [BeaconSortitionPool.address],
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "BeaconDkgValidator",
      address: BeaconDkgValidator.address,
    })
  }
}

export default func

func.tags = ["BeaconDkgValidator"]
func.dependencies = ["BeaconSortitionPool"]
