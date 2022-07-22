import { task, types } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("register-operator", "Increases authorization and registers the operator ")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .addParam("application", "Name of Application Contract", undefined, types.string)
  .setAction(async (args, hre) => {
    await setup(hre, args)
  })

async function setup(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    operator: string
    application: string
  }
) {
  const { ethers, helpers } = hre
  const { provider, operator, application } = args

  const applicationContract = await helpers.contracts.getContract(application)

  console.log(
    `Registering operator ${operator.toString()} for a staking provider ${provider.toString()} ...`
  )

  const stakingProviderSigner = await ethers.getSigner(provider)
  await (
    await applicationContract
      .connect(stakingProviderSigner)
      .registerOperator(operator)
  ).wait()
}
