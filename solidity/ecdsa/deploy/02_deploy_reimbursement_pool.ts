/* eslint-disable import/no-extraneous-dependencies */
import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, ethers, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const staticGas = 48_000 // gas amount consumed by the refund() + tx cost
  const maxGasPrice = ethers.utils.parseUnits("500", "gwei")

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
