import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()

  const SortitionPool = await deployments.get("SortitionPool")
  const DKGValidator = await deployments.get("DKGValidator")

  const DKG = await deployments.deploy("DKG", {
    from: deployer,
    log: true,
  })

  const WalletFactory = await deployments.deploy("WalletFactory", {
    from: deployer,
    args: [SortitionPool.address, DKGValidator.address],
    libraries: { DKG: DKG.address },
    log: true,
  })

  await helpers.ownable.transferOwnership(
    "SortitionPool",
    WalletFactory.address,
    deployer
  )

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "WalletFactory",
      address: WalletFactory.address,
    })
  }
}

export default func

func.tags = ["WalletFactory"]
func.dependencies = ["SortitionPool", "DKGValidator"]
