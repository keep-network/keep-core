import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const EcdsaSortitionPool = await deployments.get("EcdsaSortitionPool")

  const EcdsaDkgValidator = await deployments.deploy("EcdsaDkgValidator", {
    from: deployer,
    args: [EcdsaSortitionPool.address],
    log: true,
    waitConfirmations: 1,
  })

  if (hre.network.tags.etherscan) {
    await helpers.etherscan.verify(EcdsaDkgValidator)
  }

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "EcdsaDkgValidator",
      address: EcdsaDkgValidator.address,
    })
  }
}

export default func

func.tags = ["EcdsaDkgValidator"]
func.dependencies = ["EcdsaSortitionPool"]
