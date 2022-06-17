/* eslint-disable @typescript-eslint/no-unused-expressions, import/no-extraneous-dependencies */
import { ethers, helpers } from "hardhat"

import type { T, TokenStaking, RandomBeacon } from "../typechain"

const { to1e18 } = helpers.number

const BLOCK_TIME = 5000 // in ms

async function getAccounts() {
  const provider = new ethers.providers.JsonRpcProvider(process.env.NETWORK)
  return provider.listAccounts()
}

async function init() {
  const t: T = await helpers.contracts.getContract<T>("T")
  const staking: TokenStaking =
    await helpers.contracts.getContract<TokenStaking>("TokenStaking")
  const randomBeacon: RandomBeacon = await helpers.contracts.getContract(
    "RandomBeacon"
  )

  const accounts = await getAccounts()
  const deployer = accounts[0]
  const stakeOwner = accounts[1]
  const stakingProvider = accounts[2]
  const beneficiary = stakeOwner
  const authorizer = stakeOwner

  const stakedAmount = to1e18(1_000_000) // 1MM T

  const deployerSigner = await ethers.getSigner(deployer)
  await t.connect(deployerSigner).mint(stakeOwner, stakedAmount)
  const stakeOwnerSigner = await ethers.getSigner(stakeOwner)
  await t.connect(stakeOwnerSigner).approve(staking.address, stakedAmount)

  // We need to wait until the block is mined and that transactions in the
  // mempool were executed.
  await new Promise((f) => setTimeout(f, BLOCK_TIME))

  console.log(
    `Staking ${stakedAmount.toString()} to the staking provider ${stakingProvider} ...`
  )

  await staking
    .connect(stakeOwnerSigner)
    .stake(stakingProvider, beneficiary, authorizer, stakedAmount)

  await new Promise((f) => setTimeout(f, BLOCK_TIME))

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

  await staking
    .connect(authorizerSigner)
    .increaseAuthorization(
      stakingProvider,
      randomBeacon.address,
      minimumAuthorization
    )

  await new Promise((f) => setTimeout(f, BLOCK_TIME))

  const authorizedStaked = await staking.authorizedStake(
    stakingProvider,
    randomBeacon.address
  )

  console.log(
    `Staked authorization ${authorizedStaked.toString()} was increased for the Random Beacon`
  )
}

init()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error.error)
    process.exit(1)
  })
