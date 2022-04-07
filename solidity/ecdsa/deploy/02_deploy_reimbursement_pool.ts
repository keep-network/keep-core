/* eslint-disable import/no-extraneous-dependencies */
import { ethers } from "hardhat"

import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, deployments } = hre
  const { deployer } = await getNamedAccounts()
  const deployerSigner = await ethers.getSigner(deployer)

  const staticGas = 40800 // gas amount consumed by the refund() + tx cost
  const maxGasPrice = 500000000000 // 500 gwei

  const ReimbursementPool = await deployments.deploy("ReimbursementPool", {
    from: deployer,
    args: [staticGas, maxGasPrice],
    log: true,
  })

  await deployerSigner.sendTransaction({
    to: ReimbursementPool.address,
    value: ethers.utils.parseEther("100.0"), // Send 100.0 ETH
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
