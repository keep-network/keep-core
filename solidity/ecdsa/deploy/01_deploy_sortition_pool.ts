import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const POOL_WEIGHT_DIVISOR = 1 // TODO: Update value

  const TokenStaking = await deployments.get("TokenStaking")
  const T = await deployments.get("T")

  const SortitionPool = await deployments.deploy("SortitionPool", {
    from: deployer,
    args: [TokenStaking.address, T.address, POOL_WEIGHT_DIVISOR],
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "SortitionPool",
      address: SortitionPool.address,
    })
  }
}

export default func

func.tags = ["SortitionPool"]
// TokenStaking and T deployments are expected to be resolved from
// @threshold-network/solidity-contracts
func.dependencies = ["TokenStaking", "T"]
