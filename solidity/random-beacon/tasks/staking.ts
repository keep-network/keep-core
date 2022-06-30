import { task, types } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

// TODO: split to 4 different tasks: mint, stake, increase authorization and
// register an operator
task(
  "stake",
  "Stakes T tokens, increases authorization and registers the operator "
)
  .addParam("owner", "Stake owner address", undefined, types.string)
  .addParam("provider", "Staking provider", undefined, types.string)
  .addParam("operator", "Staking operator", undefined, types.string)
  .addOptionalParam("beneficiary", "Stake beneficiary", undefined, types.string)
  .addOptionalParam("authorizer", "Stake authorizer", undefined, types.string)
  .addOptionalParam("amount", "Stake amount", 1_000_000, types.int)
  .addOptionalParam(
    "authorization",
    "Authorization amount",
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
    amount: BigInteger
    authorization: BigInteger
  }
) {

  const { ethers, helpers, getNamedAccounts } = hre
  let { owner, provider, operator, beneficiary, authorizer, amount, authorization } = args

  const { deployer } = await getNamedAccounts()
  const { to1e18 } = helpers.number
  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")
  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  const stakedAmount = to1e18(amount)

  const deployerSigner = await ethers.getSigner(deployer)
  await (await t.connect(deployerSigner).mint(owner, stakedAmount)).wait()
  const stakeOwnerSigner = await ethers.getSigner(owner)
  await (
    await t.connect(stakeOwnerSigner).approve(staking.address, stakedAmount)
  ).wait()

  // Beneficiary can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  if (beneficiary === undefined) {
    beneficiary = owner
  }

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  if (authorizer === undefined) {
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

  const stakingAuthorization =
    to1e18(authorization) || (await randomBeacon.minimumAuthorization())
  const authorizerSigner = await ethers.getSigner(authorizer)

  console.log(
    `Increasing authorization ${stakingAuthorization.toString()} for the Random Beacon ...`
  )

  await (
    await staking
      .connect(authorizerSigner)
      .increaseAuthorization(
        provider,
        randomBeacon.address,
        stakingAuthorization
      )
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
    await randomBeacon
      .connect(stakingProviderSigner)
      .registerOperator(operator)
  ).wait()
}
