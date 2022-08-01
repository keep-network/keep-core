import { task, types } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("register-operator", "Registers an operator")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .setAction(async (args, hre) => {
    await setup(hre, args)
  })

async function setup(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    operator: string
  }
) {
  const { ethers, helpers } = hre
  const { provider, operator } = args

  const walletRegistry = await helpers.contracts.getContract("WalletRegistry")

  console.log(
    `Registering operator ${operator.toString()} for a staking provider ${provider.toString()}...`
  )

  const stakingProviderSigner = await ethers.getSigner(provider)
  await (
    await walletRegistry
      .connect(stakingProviderSigner)
      .registerOperator(operator)
  ).wait()
}
