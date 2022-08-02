import { task, types } from "hardhat/config"

import type { BigNumberish } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("initialize:authorize", "Sets authorization")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addOptionalParam("authorizer", "Stake Authorizer", undefined, types.string)
  .addOptionalParam(
    "authorization",
    "Authorization amount (default: minimumAuthorization)",
    undefined,
    types.int
  )
  .setAction(async (args, hre) => {
    await increaseAuthorization(hre, args)
  })

async function increaseAuthorization(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    authorizer: string
    authorization: BigNumberish
  }
) {
  const { ethers, helpers } = hre
  const owner = ethers.utils.getAddress(args.owner)
  const provider = ethers.utils.getAddress(args.provider)

  const walletRegistry = await helpers.contracts.getContract("WalletRegistry")

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const authorizer = args.authorizer
    ? ethers.utils.getAddress(args.authorizer)
    : owner

  const { to1e18, from1e18 } = helpers.number
  const staking = await helpers.contracts.getContract("TokenStaking")

  const authorization = args.authorization
    ? to1e18(args.authorization)
    : await walletRegistry.minimumAuthorization()

  const currentAuthorization = await staking.authorizedStake(
    provider,
    walletRegistry.address
  )

  if (currentAuthorization.gte(authorization)) {
    console.log(
      `Authorized stake for the Wallet Registry is ${from1e18(
        currentAuthorization
      )} T`
    )
    return
  }

  const increaseAmount = authorization.sub(currentAuthorization)

  console.log(
    `Increasing Wallet Registry's authorization by ${from1e18(
      increaseAmount
    )} T to ${from1e18(authorization)} T...`
  )

  await (
    await staking
      .connect(await ethers.getSigner(authorizer))
      .increaseAuthorization(provider, walletRegistry.address, increaseAmount)
  ).wait()
}
