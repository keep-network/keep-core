/* eslint-disable no-console */
import type { BigNumberish, BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

// eslint-disable-next-line import/prefer-default-export
export async function mint(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    amount: BigNumberish
  }
): Promise<void> {
  const { ethers, helpers } = hre
  const { to1e18, from1e18 } = helpers.number
  const owner = ethers.utils.getAddress(args.owner)
  const stakeAmount = to1e18(args.amount)

  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")

  const tokenContractOwner = await t.owner()

  const currentBalance: BigNumber = await t.balanceOf(owner)

  console.log(`Account ${owner} balance is ${from1e18(currentBalance)} T`)

  if (currentBalance.lt(stakeAmount)) {
    const mintAmount = stakeAmount.sub(currentBalance)

    console.log(`Minting ${from1e18(mintAmount)} T for ${owner}...`)

    await (
      await t
        .connect(await ethers.getSigner(tokenContractOwner))
        .mint(owner, mintAmount)
    ).wait()
  }

  const currentAllowance: BigNumber = await t.allowance(owner, staking.address)

  console.log(
    `Account ${owner} allowance for ${staking.address} is ${from1e18(
      currentAllowance
    )} T`
  )

  if (currentAllowance.lt(stakeAmount)) {
    console.log(
      `Approving ${from1e18(stakeAmount)} T for ${staking.address}...`
    )
    await (
      await t
        .connect(await ethers.getSigner(owner))
        .approve(staking.address, stakeAmount)
    ).wait()
  }
}
