import { task, types } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("initialize:register", "Registers an operator for a staking provider")
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .setAction(async (args, hre) => {
    await registerOperator(hre, args)
  })

async function registerOperator(
  hre: HardhatRuntimeEnvironment,
  args: {
    provider: string
    operator: string
  }
) {
  const { ethers, helpers } = hre

  const provider = ethers.utils.getAddress(args.provider)
  const operator = ethers.utils.getAddress(args.operator)

  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  const currentProvider = ethers.utils.getAddress(
    await randomBeacon.callStatic.operatorToStakingProvider(operator)
  )

  switch (currentProvider) {
    case provider: {
      console.log(
        `Current staking provider for operator ${operator} is ${currentProvider}`
      )
      return
    }
    case ethers.constants.AddressZero: {
      console.log(
        `Registering operator ${operator} for a staking provider ${provider}...`
      )

      await (
        await randomBeacon
          .connect(await ethers.getSigner(provider))
          .registerOperator(operator)
      ).wait()

      break
    }
    default: {
      throw new Error(
        `Operator [${operator}] has already been registered for another staking provider [${currentProvider}]`
      )
    }
  }
}
