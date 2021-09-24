import { HardhatRuntimeEnvironment } from "hardhat/types"
import { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const underwriterToken = await deployments.deploy("TestToken", {
    from: deployer,
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "UnderwriterToken",
      address: underwriterToken.address,
    })
  }
}

export default func

func.tags = ["UnderwriterToken"]