import { task, types } from "hardhat/config"

import type { BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

// TODO: split to 4 different tasks: mint, stake, increase authorization and
// register an operator
task(
  "stake",
  "Stakes T tokens, increases authorization and registers the operator "
)
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .addOptionalParam("beneficiary", "Stake Beneficiary", undefined, types.string)
  .addOptionalParam("authorizer", "Stake Authorizer", undefined, types.string)
  .addOptionalParam("amount", "Stake amount", 1_000_000, types.int)
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
    operator: string
    beneficiary: string
    authorizer: string
    amount: BigNumber
    authorization: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, provider, operator, amount } = args
  let { beneficiary, authorizer, authorization } = args

  const { deployer } = await helpers.signers.getNamedSigners()
  const { to1e18 } = helpers.number
  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")
  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  const stakedAmount = to1e18(amount)

  await (await t.connect(deployer).mint(owner, stakedAmount)).wait()
  const stakeOwnerSigner = await ethers.getSigner(owner)
  await (
    await t.connect(stakeOwnerSigner).approve(staking.address, stakedAmount)
  ).wait()

  // Beneficiary can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  if (!beneficiary) {
    beneficiary = owner
  }

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  if (!authorizer) {
    authorizer = owner
  }

  console.log(
    `Staking ${stakedAmount.toString()} to the staking provider ${provider} ...`
  )

  await (
    await staking
      .connect(stakeOwnerSigner)
      .stake(provider, beneficiary, authorizer, stakedAmount)
  ).wait()

  console.log(
    `T balance of the staking contract: ${(
      await t.balanceOf(staking.address)
    ).toString()}`
  )

  if (authorization) {
    authorization = to1e18(authorization)
  } else {
    authorization = await randomBeacon.minimumAuthorization()
  }

  const authorizerSigner = await ethers.getSigner(authorizer)

  console.log(
    `Increasing authorization ${authorization.toString()} for the Random Beacon ...`
  )

  await (
    await staking
      .connect(authorizerSigner)
      .increaseAuthorization(provider, randomBeacon.address, authorization)
  ).wait()

  const authorizedStaked = await staking.authorizedStake(
    provider,
    randomBeacon.address
  )

  console.log(
    `Staked authorization ${authorizedStaked.toString()} was increased for the Random Beacon`
  )

  console.log(
    `Registering operator ${operator.toString()} for a staking provider ${provider.toString()} ...`
  )

  const stakingProviderSigner = await ethers.getSigner(provider)
  await (
    await randomBeacon.connect(stakingProviderSigner).registerOperator(operator)
  ).wait()
}
