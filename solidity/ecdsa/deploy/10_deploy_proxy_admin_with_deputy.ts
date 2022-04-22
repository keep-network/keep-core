import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"
import type { ProxyAdmin } from "../typechain"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { ethers, getNamedAccounts, upgrades, deployments } = hre
  const { esdm } = await getNamedAccounts()

  const ProxyAdminWithDeputy = await deployments.get("ProxyAdminWithDeputy")

  const WalletRegistry = await deployments.get("WalletRegistry")

  const proxyAdmin = (await upgrades.admin.getInstance()) as ProxyAdmin

  await proxyAdmin
    .connect(await ethers.getSigner(esdm))
    .changeProxyAdmin(WalletRegistry.address, ProxyAdminWithDeputy.address)
}

export default func

func.tags = ["ProxyAdminWithDeputy"]
func.dependencies = ["WalletRegistry"]

// For now we skip this script as DAO and ProxyAdminWithDeputy are not yet
// established.
func.skip = async () => true
