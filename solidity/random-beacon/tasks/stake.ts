import { task, types } from "hardhat/config"

import type { BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("stake", "Stakes T token")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addOptionalParam("beneficiary", "Stake Beneficiary", undefined, types.string)
  .addOptionalParam("authorizer", "Stake Authorizer", undefined, types.string)
  .addOptionalParam("amount", "Stake amount", 1_000_000, types.int)
  .setAction(async (args, hre) => {
    await setup(hre, args)
  })

async function setup(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    beneficiary: string
    authorizer: string
    amount: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, provider, amount } = args
  let { beneficiary, authorizer } = args

  const { deployer } = await helpers.signers.getNamedSigners()
  const { to1e18 } = helpers.number
  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")

  const stakedAmount = to1e18(amount)

  await (await t.connect(deployer).mint(owner, stakedAmount)).wait()
  const stakeOwnerSigner = await ethers.getSigner(owner)
  await (
    await t.connect(stakeOwnerSigner).approve(staking.address, stakedAmount)
  ).wait()

  // If not set, beneficiary can be the owner. This simplification is used for
  // development purposes.
  if (!beneficiary) {
    beneficiary = owner
  }

  // If not set, authorizer can be the owner. This simplification is used for
  // development purposes.
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
}
