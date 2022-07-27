import { task, types } from "hardhat/config"

import type { BigNumber } from "ethers"
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
  .addOptionalParam("force", "Force initialization", false, types.boolean)
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
    amount: BigNumber
    authorization: BigNumber
    force: boolean
  }
) {
  if (await isAlreadyStaked(hre, args.operator)) {
    console.log(`Operator ${args.operator} is already staked`)
    return
  }

  await hre.run(TASK_MINT, args)
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
    amount: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, amount } = args

  const { to1e18, from1e18 } = helpers.number
  const t = await helpers.contracts.getContract("T")
  const staking = await helpers.contracts.getContract("TokenStaking")

  const stakeAmount = to1e18(amount)

  // TODO: Check if he already have tokens
  const tokenContractOwner = await t.owner()

  const currentBalance: BigNumber = await t.balanceOf(owner)
  if (currentBalance.lt(stakeAmount)) {
    const mintAmount = stakeAmount.sub(currentBalance)

    console.log(`Minting ${from1e18(mintAmount)} T for ${owner}...`)

    await (
      await t
        .connect(await ethers.getSigner(tokenContractOwner))
        .mint(owner, mintAmount)
    ).wait()
  } else {
    console.log(`Account ${owner} already holds ${from1e18(currentBalance)} T`)
  }

  console.log(`Approving ${from1e18(stakeAmount)} T for ${staking.address}...`)
  await (
    await t
      .connect(await ethers.getSigner(owner))
      .approve(staking.address, stakeAmount)
  ).wait()
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
    amount: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, provider } = args

  // Beneficiary can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const beneficiary = args.beneficiary ?? owner

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const authorizer = args.authorizer ?? owner

  const { to1e18, from1e18 } = helpers.number
  const staking = await helpers.contracts.getContract("TokenStaking")

  const stakedAmount = to1e18(args.amount)

  console.log(
    `Staking ${from1e18(stakedAmount)} T to the staking provider ${provider}...`
  )

  await (
    await staking
      .connect(await ethers.getSigner(owner))
      .stake(provider, beneficiary, authorizer, stakedAmount)
  ).wait()
}

task(TASK_AUTHORIZE, "Sets authorization")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addOptionalParam("beneficiary", "Stake Beneficiary", undefined, types.string)
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
    authorization: BigNumber
  }
) {
  const { ethers, helpers } = hre
  const { owner, provider } = args

  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  // Authorizer can equal to the owner if not set otherwise. This simplification
  // is used for development purposes.
  const authorizer = args.authorizer ?? owner

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
  const { provider, operator } = args

  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")

  console.log(
    `Registering operator ${operator.toString()} for a staking provider ${provider.toString()}...`
  )

  await (
    await randomBeacon
      .connect(await ethers.getSigner(provider))
      .registerOperator(operator)
  ).wait()
}

async function isAlreadyStaked(
  hre: HardhatRuntimeEnvironment,
  operator: string
): Promise<boolean> {
  const { ethers, helpers } = hre
  const { from1e18 } = helpers.number

  const randomBeacon = await helpers.contracts.getContract("RandomBeacon")
  const staking = await helpers.contracts.getContract("TokenStaking")

  const currentStakingProvider =
    await randomBeacon.callStatic.operatorToStakingProvider(operator)
  console.log(
    `Current staking provider for operator ${operator} is ${currentStakingProvider}`
  )

  if (currentStakingProvider === ethers.constants.AddressZero) {
    return false
  }

  const currentAuthorization = await staking.authorizedStake(
    currentStakingProvider,
    randomBeacon.address
  )
  console.log(
    `Current authorization for ${
      randomBeacon.address
    } application is ${from1e18(currentAuthorization)} T`
  )
  if (currentAuthorization === 0) {
    return false
  }

  return true
}
