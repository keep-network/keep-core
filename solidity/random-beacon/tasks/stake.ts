import { task, types } from "hardhat/config"

import type { BigNumberish } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("initialize:stake", "Stakes T token")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addOptionalParam("beneficiary", "Stake Beneficiary", undefined, types.string)
  .addOptionalParam("authorizer", "Stake Authorizer", undefined, types.string)
  .addOptionalParam("amount", "Stake amount", 1_000_000, types.int)
  .setAction(async (args, hre) => {
    await stake(hre, args)
  })

async function stake(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    beneficiary: string
    authorizer: string
    amount: BigNumberish
  }
) {
  const { ethers, helpers } = hre
  const { to1e18, from1e18 } = helpers.number
  const owner = ethers.utils.getAddress(args.owner)
  const provider = ethers.utils.getAddress(args.provider)
  const stakeAmount = to1e18(args.amount)

  // Beneficiary can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const beneficiary = args.beneficiary
    ? ethers.utils.getAddress(args.beneficiary)
    : owner

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const authorizer = args.authorizer
    ? ethers.utils.getAddress(args.authorizer)
    : owner

  const staking = await helpers.contracts.getContract("TokenStaking")

  const { tStake: currentStake } = await staking.callStatic.stakes(provider)

  console.log(`Current stake for ${provider} is ${from1e18(currentStake)} T`)

  if (currentStake.eq(0)) {
    console.log(
      `Staking ${from1e18(
        stakeAmount
      )} T to the staking provider ${provider}...`
    )

    await (
      await staking
        .connect(await ethers.getSigner(owner))
        .stake(provider, beneficiary, authorizer, stakeAmount)
    ).wait()
  } else if (currentStake.lt(stakeAmount)) {
    const topUpAmount = stakeAmount.sub(currentStake)

    console.log(
      `Topping up ${from1e18(
        topUpAmount
      )} T to the staking provider ${provider}...`
    )

    await (
      await staking
        .connect(await ethers.getSigner(owner))
        .topUp(provider, topUpAmount)
    ).wait()
  }
}
