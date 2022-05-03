import type { HardhatRuntimeEnvironment } from "hardhat/types"
import type { DeployFunction } from "hardhat-deploy/types"

const func: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { getNamedAccounts, ethers, deployments } = hre
  const { deployer } = await getNamedAccounts()

  const staticGas = 40_800 // gas amount consumed by the refund() + tx cost
  // FIXME: As a workaround for a bug in hardhat-gas-reporter #86 we need to provide
  // alternative deployment script to obtain a gas report.
  // #86: https://github.com/cgewecke/hardhat-gas-reporter/issues/86
  const maxGasPrice =
    process.env.GAS_REPORTER_BUG_WORKAROUND === "true"
      ? 500_000_000_000
      : ethers.utils.parseUnits("500", "gwei")

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
