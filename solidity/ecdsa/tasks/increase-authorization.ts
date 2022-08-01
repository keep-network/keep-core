import { task, types } from "hardhat/config"

import type { BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("increase-authorization", "Increases authorization")
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
    await setup(hre, args)
  })

async function setup(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    authorizer: string
    authorization: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, provider } = args
  let { authorizer, authorization } = args

  const { to1e18, from1e18 } = helpers.number
  const staking = await helpers.contracts.getContract("TokenStaking")
  const walletRegistry = await helpers.contracts.getContract("WalletRegistry")

  // If not set, authorizer can be the owner. This simplification is used for
  // development purposes.
  if (!authorizer) {
    authorizer = owner
  }

  if (authorization) {
    authorization = to1e18(authorization)
  } else {
    authorization = await walletRegistry.minimumAuthorization()
  }

  const authorizerSigner = await ethers.getSigner(authorizer)

  console.log(
    `Increasing authorization ${from1e18(
      authorization
    )} for the Wallet Registry...`
  )

  await (
    await staking
      .connect(authorizerSigner)
      .increaseAuthorization(provider, walletRegistry.address, authorization)
  ).wait()

  const authorizedStaked = await staking.authorizedStake(
    provider,
    walletRegistry.address
  )

  console.log(
    `Authorization for Wallet Registry was increased to ${authorizedStaked.toString()}`
  )
}
