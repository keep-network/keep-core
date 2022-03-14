import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const staticGas = 41900 // gas amount consumed by the refund() + tx cost
  const maxGasPrice = 20000000000 // 20 gwei

  const ReimbursementPool = await deployments.deploy("ReimbursementPool", {
    from: deployer,
    args: [staticGas, maxGasPrice],
    log: true,
  })

  if (hre.network.tags.tenderly) {
    await hre.tenderly.verify({
      name: "ReimbursementPool",
      address: ReimbursementPool.address,
    })
  }
}

export default func

func.tags = ["ReimbursementPool"]
