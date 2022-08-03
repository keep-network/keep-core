/* eslint-disable no-console */
import type { BigNumberish } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

// eslint-disable-next-line import/prefer-default-export
export async function authorize(
  hre: HardhatRuntimeEnvironment,
  deploymentName: string, // TODO: Change to IApplication
  args: {
    owner: string
    provider: string
    authorizer: string
    authorization: BigNumberish
  }
): Promise<void> {
  const { ethers, helpers } = hre
  const owner = ethers.utils.getAddress(args.owner)
  const provider = ethers.utils.getAddress(args.provider)

  const application = await helpers.contracts.getContract(deploymentName)

  console.log(
    `Authorizing provider's ${provider} stake in ${deploymentName} application (${application.address})`
  )

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const authorizer = args.authorizer
    ? ethers.utils.getAddress(args.authorizer)
    : owner

  const { to1e18, from1e18 } = helpers.number
  const staking = await helpers.contracts.getContract("TokenStaking")

  const authorization = args.authorization
    ? to1e18(args.authorization)
    : await application.minimumAuthorization()

  const currentAuthorization = await staking.authorizedStake(
    provider,
    application.address
  )

  if (currentAuthorization.gte(authorization)) {
    console.log(
      `Authorized stake is already ${from1e18(currentAuthorization)} T`
    )
    return
  }

  const increaseAmount = authorization.sub(currentAuthorization)

  console.log(
    `Increasing authorization by ${from1e18(increaseAmount)} T to ${from1e18(
      authorization
    )} T...`
  )

  await (
    await staking
      .connect(await ethers.getSigner(authorizer))
      .increaseAuthorization(provider, application.address, increaseAmount)
  ).wait()
}
