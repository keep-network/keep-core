/* eslint-disable import/no-extraneous-dependencies */
import { ethers, helpers, getNamedAccounts } from "hardhat"

import type { T, TokenStaking, RandomBeacon } from "../typechain"

const { to1e18 } = helpers.number

async function getAccounts() {
  const provider = new ethers.providers.JsonRpcProvider(process.env.NETWORK)
  return provider.listAccounts()
}

async function setup() {
  const t: T = await helpers.contracts.getContract<T>("T")
  const staking: TokenStaking =
    await helpers.contracts.getContract<TokenStaking>("TokenStaking")
  const randomBeacon: RandomBeacon = await helpers.contracts.getContract(
    "RandomBeacon"
  )

  const { deployer } = await getNamedAccounts()
  const accounts = await getAccounts()
  const stakeOwner = accounts[1]
  const stakingProvider = accounts[2]
  const operator = accounts[3]
  const beneficiary = stakeOwner
  const authorizer = stakeOwner

  const stakedAmount = to1e18(1_000_000) // 1MM T

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

  const minimumAuthorization = await randomBeacon.minimumAuthorization()
  const authorizerSigner = await ethers.getSigner(authorizer)

  console.log(
    `Increasing min authorization ${minimumAuthorization.toString()} for the Random Beacon ...`
  )

  await (
    await staking
      .connect(authorizerSigner)
      .increaseAuthorization(
        stakingProvider,
        randomBeacon.address,
        minimumAuthorization
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
    `Registering operator ${operator.toString()} for a staking provider ${stakingProvider.toString()}`
  )

  const stakingProviderSigner = await ethers.getSigner(stakingProvider)
  await (
    await randomBeacon.connect(stakingProviderSigner).registerOperator(operator)
  ).wait()
}

setup()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
