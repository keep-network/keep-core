import { writeFileSync } from "fs"
import path from "path"

import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const MasterWallet = await deployments.deploy("MasterWallet", {
    contract: "Wallet",
    from: deployer,
    log: true,
  })

  if (deployments.getNetworkName() !== "hardhat") {
    // Store plain Wallet artifact in the deployments directory for usage in other projects.
    writeFileSync(
      path.join(
        hre.config.paths.deployments,
        deployments.getNetworkName(),
        "Wallet.json"
      ),
      JSON.stringify(await deployments.getArtifact("Wallet"), null, 2)
    )
  }

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "MasterWallet",
      address: MasterWallet.address,
    })
  }
}

export default func

func.tags = ["MasterWallet"]
