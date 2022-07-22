import { task, types } from "hardhat/config"

import type { BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("increase-authorization", "Increases authorization")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam(
    "application",
    "Name of Application Contract",
    undefined,
    types.string
  )
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
    application: string
    authorizer: string
    authorization: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, provider, application } = args
  let { authorizer, authorization } = args

  const { to1e18 } = helpers.number
  const staking = await helpers.contracts.getContract("TokenStaking")
  const applicationContract = await helpers.contracts.getContract(application)

  // If not set, authorizer can be the owner. This simplification is used for
  // development purposes.
  if (!authorizer) {
    authorizer = owner
  }

  if (authorization) {
    authorization = to1e18(authorization)
  } else {
    authorization = await applicationContract.minimumAuthorization()
  }

  const authorizerSigner = await ethers.getSigner(authorizer)

  console.log(
    `Increasing authorization ${authorization.toString()} for the ${application.toString()} ...`
  )

  await (
    await staking
      .connect(authorizerSigner)
      .increaseAuthorization(
        provider,
        applicationContract.address,
        authorization
      )
  ).wait()

  const authorizedStaked = await staking.authorizedStake(
    provider,
    applicationContract.address
  )

  console.log(
    `Staked authorization ${authorizedStaked.toString()} was increased for the ${application.toString()}`
  )
}
