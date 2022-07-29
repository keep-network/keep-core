/* eslint-disable no-console */
import { task, types } from "hardhat/config"

import type { BigNumberish, BigNumber } from "ethers"
import type { HardhatRuntimeEnvironment } from "hardhat/types"

const TASK_INITIALIZE = "initialize"
const TASK_MINT = `${TASK_INITIALIZE}:mint`
const TASK_STAKE = `${TASK_INITIALIZE}:stake`
const TASK_AUTHORIZE = `${TASK_INITIALIZE}:authorize`
const TASK_REGISTER = `${TASK_INITIALIZE}:register`

task(TASK_INITIALIZE, "Initializes staking for an operator")
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
    await initialize(hre, args)
  })

async function initialize(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    operator: string
    beneficiary: string
    authorizer: string
    amount: BigNumberish
    authorization: BigNumberish
  }
) {
  const tokensToMint = await calculateTokensNeededForStake(
    hre,
    args.provider,
    args.amount
  )

  if (!tokensToMint.isZero()) {
    await hre.run(TASK_MINT, { ...args, amount: tokensToMint.toNumber() })
  }

  await hre.run(TASK_STAKE, args)
  await hre.run(TASK_AUTHORIZE, args)
  await hre.run(TASK_REGISTER, args)
}

task(TASK_MINT, "Mints T tokens")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addOptionalParam("amount", "Stake amount", 1_000_000, types.int)
  .setAction(async (args, hre) => {
    await mint(hre, args)
  })

async function mint(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    amount: BigNumberish
  }
) {
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

task(TASK_STAKE, "Stakes T tokens")
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

task(TASK_AUTHORIZE, "Sets authorization")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addOptionalParam("authorizer", "Stake Authorizer", undefined, types.string)
  .addOptionalParam(
    "authorization",
    "Authorization amount (default: minimumAuthorization)",
    undefined,
    types.int
  )
  .setAction(async (args, hre) => {
    await authorize(hre, args)
  })

async function authorize(
  hre: HardhatRuntimeEnvironment,
  args: {
    owner: string
    provider: string
    authorizer: string
    authorization: BigNumberish
  }
) {
  const { ethers, helpers } = hre
  const owner = ethers.utils.getAddress(args.owner)
  const provider = ethers.utils.getAddress(args.provider)

  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const authorizer = args.authorizer
    ? ethers.utils.getAddress(args.authorizer)
    : owner

  const { to1e18, from1e18 } = helpers.number
  const staking = await helpers.contracts.getContract("TokenStaking")

  const authorization = args.authorization
    ? to1e18(args.authorization)
    : await randomBeacon.minimumAuthorization()

  const currentAuthorization = await staking.authorizedStake(
    provider,
    randomBeacon.address
  )

  if (currentAuthorization.gte(authorization)) {
    console.log(
      `Authorized stake for the Random Beacon is ${from1e18(
        currentAuthorization
      )} T`
    )
    return
  }

  const increaseAmount = authorization.sub(currentAuthorization)

  console.log(
    `Increasing Random Beacon's authorization by ${from1e18(
      increaseAmount
    )} T to ${from1e18(authorization)} T...`
  )

  await (
    await staking
      .connect(await ethers.getSigner(authorizer))
      .increaseAuthorization(provider, randomBeacon.address, increaseAmount)
  ).wait()
}

task(TASK_REGISTER, "Registers an operator for a staking provider")
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .setAction(async (args, hre) => {
    await register(hre, args)
  })

async function register(
  hre: HardhatRuntimeEnvironment,
  args: {
    provider: string
    operator: string
  }
) {
  const { ethers, helpers } = hre

  const provider = ethers.utils.getAddress(args.provider)
  const operator = ethers.utils.getAddress(args.operator)

  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  const currentProvider = ethers.utils.getAddress(
    await randomBeacon.callStatic.operatorToStakingProvider(operator)
  )

  switch (currentProvider) {
    case provider: {
      console.log(
        `Current staking provider for operator ${operator} is ${currentProvider}`
      )
      return
    }
    case ethers.constants.AddressZero: {
      console.log(
        `Registering operator ${operator} for a staking provider ${provider}...`
      )

      await (
        await randomBeacon
          .connect(await ethers.getSigner(provider))
          .registerOperator(operator)
      ).wait()

      break
    }
    default: {
      throw new Error(
        `Operator [${operator}] has already been registered for another staking provider [${currentProvider}]`
      )
    }
  }
}

async function calculateTokensNeededForStake(
  hre: HardhatRuntimeEnvironment,
  provider: string,
  amount: BigNumberish
): Promise<BigNumber> {
  const { ethers, helpers } = hre
  const { to1e18, from1e18 } = helpers.number

  const stakeAmount = to1e18(amount)

  const staking = await helpers.contracts.getContract("TokenStaking")

  const { tStake: currentStake } = await staking.callStatic.stakes(provider)

  if (currentStake.lt(stakeAmount)) {
    return ethers.BigNumber.from(from1e18(stakeAmount.sub(currentStake)))
  }

  return ethers.constants.Zero
}
