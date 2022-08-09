import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { deployments, helpers } = hre
  const { log } = deployments

  const TokenStaking = await deployments.getOrNull("TokenStaking")

  if (TokenStaking && helpers.address.isValid(TokenStaking.address)) {
    log(`using existing TokenStaking at ${TokenStaking.address}`)
  } else {
    throw new Error("deployed TokenStaking contract not found")
  }
}

export default func

func.tags = ["TokenStaking"]
