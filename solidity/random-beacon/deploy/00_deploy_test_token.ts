import { HardhatRuntimeEnvironment } from "hardhat/types"
import { DeployFunction } from "hardhat-deploy/types"

// This is an example of deployment script.
// Should be removed once we have actual contracts to deploy.
const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const testToken = await deployments.deploy("TestToken", {
    from: deployer,
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "TestToken",
      address: testToken.address,
    })
  }
}

export default func

func.tags = ["TestToken"]