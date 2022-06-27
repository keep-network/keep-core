import { task, types } from "hardhat/config"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

task(
  "staking",
  "Stakes T tokens, increases authorization and registers the operator "
)
  .addOptionalParam("owner", "Staking owner address", undefined, types.string)
  .addOptionalParam("provider", "Staking provider", undefined, types.string)
  .addOptionalParam("operator", "Staking operator", undefined, types.string)
  .addOptionalParam("amount", "Staking amount", 1_000_000, types.int)
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
    amount: BigInteger
    authorization: BigInteger
  }
) {
  const { ethers, helpers, getNamedAccounts } = hre
  const { owner, provider, operator, amount, authorization } = args

  const { deployer } = await getNamedAccounts()
  const { to1e18 } = helpers.number
  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")
  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")
  const accounts = await new ethers.providers.JsonRpcProvider().listAccounts()
  const stakeOwner = owner || accounts[1]
  const stakingProvider = provider || accounts[2]
  const stakingOperator = operator || accounts[3]
  const beneficiary = stakeOwner
  const authorizer = stakeOwner

  const stakedAmount = to1e18(amount)

  const deployerSigner = await ethers.getSigner(deployer)
  await (await t.connect(deployerSigner).mint(stakeOwner, stakedAmount)).wait()
  const stakeOwnerSigner = await ethers.getSigner(stakeOwner)
  await (
    await t.connect(stakeOwnerSigner).approve(staking.address, stakedAmount)
  ).wait()

  console.log(
    `Staking ${stakedAmount.toString()} to the staking provider ${stakingProvider} ...`
  )

  await (
    await staking
      .connect(stakeOwnerSigner)
      .stake(stakingProvider, beneficiary, authorizer, stakedAmount)
  ).wait()

  console.log(
    `T balance of the staking contract: ${(
      await t.balanceOf(staking.address)
    ).toString()}`
  )

  const stakingAuthorization =
    authorization || (await randomBeacon.minimumAuthorization())
  const authorizerSigner = await ethers.getSigner(authorizer)

  console.log(
    `Increasing authorization ${stakingAuthorization.toString()} for the Random Beacon ...`
  )

  await (
    await staking
      .connect(authorizerSigner)
      .increaseAuthorization(
        stakingProvider,
        randomBeacon.address,
        stakingAuthorization
      )
  ).wait()

  const authorizedStaked = await staking.authorizedStake(
    stakingProvider,
    randomBeacon.address
  )

  console.log(
    `Staked authorization ${authorizedStaked.toString()} was increased for the Random Beacon`
  )

  console.log(
    `Registering operator ${stakingOperator.toString()} for a staking provider ${stakingProvider.toString()} ...`
  )

  const stakingProviderSigner = await ethers.getSigner(stakingProvider)
  await (
    await randomBeacon
      .connect(stakingProviderSigner)
      .registerOperator(stakingOperator)
  ).wait()
}
