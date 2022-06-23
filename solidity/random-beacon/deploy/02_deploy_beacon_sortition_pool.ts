import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments, helpers } = hre
  const { deployer } = await getNamedAccounts()
  const { to1e18 } = helpers.number

  const POOL_WEIGHT_DIVISOR = to1e18(1) // TODO: Update value

  const T = await deployments.get("T")

  const BeaconSortitionPool = await deployments.deploy("BeaconSortitionPool", {
    contract: "SortitionPool",
    from: deployer,
    args: [T.address, POOL_WEIGHT_DIVISOR],
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "BeaconSortitionPool",
      address: BeaconSortitionPool.address,
    })
  }
}

export default func

func.tags = ["BeaconSortitionPool"]
// TokenStaking and T deployments are expected to be resolved from
// @threshold-network/solidity-contracts
func.dependencies = ["TokenStaking", "T"]
