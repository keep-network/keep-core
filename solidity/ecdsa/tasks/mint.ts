import { task, types } from "hardhat/config"

import type { BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

task("mint", "Mints and approves T tokens for staking")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addOptionalParam("amount", "Stake amount", 1_000_000, types.int)
  .setAction(async (args, hre) => {
    await setup(hre, args)
  })

async function setup(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    amount: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, amount } = args

  const { deployer } = await helpers.signers.getNamedSigners()
  const { to1e18, from1e18 } = helpers.number
  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")

  const toMintAndApprove = to1e18(amount)

  console.log(
    `Minting ${from1e18(toMintAndApprove)} T for the ${owner} and approving for staking...`
  )

  await (await t.connect(deployer).mint(owner, toMintAndApprove)).wait()
  const stakeOwnerSigner = await ethers.getSigner(owner)
  await (
    await t.connect(stakeOwnerSigner).approve(staking.address, toMintAndApprove)
  ).wait()
}
