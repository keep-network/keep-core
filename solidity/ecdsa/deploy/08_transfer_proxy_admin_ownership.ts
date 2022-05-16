import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { ethers, getNamedAccounts, upgrades, deployments, helpers } = hre
  const { esdm, deployer } = await getNamedAccounts()
  const { log } = deployments

  log(`transferring ProxyAdmin ownership to ${esdm}`)

  // TODO: Once a DAO is established we want to switch to ProxyAdminWithDeputy and
  // use the DAO as the proxy admin owner and ESDM as the deputy. Until then we
  // use ESDM as the owner of ProxyAdmin contract.
  const newProxyAdminOwner = esdm

  const proxyAdmin = await upgrades.admin.getInstance()

  const currentOwner = await proxyAdmin.owner()

  // The `@openzeppelin/hardhat-upgrades` plugin deploys a single ProxyAdmin
  // per network. We don't want to transfer the ownership if the owner is already
  // set to the desired address.
  if (!helpers.address.equal(currentOwner, newProxyAdminOwner)) {
    log(`transferring ownership of ProxyAdmin to ${newProxyAdminOwner}`)
    await proxyAdmin
      .connect(await ethers.getSigner(deployer))
      .transferOwnership(newProxyAdminOwner)
  }
}

export default func

func.tags = ["TransferProxyAdminOwnership"]
func.dependencies = ["WalletRegistry"]
