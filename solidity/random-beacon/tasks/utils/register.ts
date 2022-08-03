/* eslint-disable no-console */
import type { HardhatRuntimeEnvironment } from "hardhat/types"

// eslint-disable-next-line import/prefer-default-export
export async function register(
  hre: HardhatRuntimeEnvironment,
  deploymentName: string,
  args: {
    provider: string
    operator: string
  }
): Promise<void> {
  const { ethers, helpers } = hre

  const provider = ethers.utils.getAddress(args.provider)
  const operator = ethers.utils.getAddress(args.operator)

  const application = await helpers.contracts.getContract(deploymentName)

  console.log(
    `Registering operator ${operator} in ${deploymentName} application (${application.address})`
  )

  const currentProvider = ethers.utils.getAddress(
    await application.callStatic.operatorToStakingProvider(operator)
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
        await application
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
