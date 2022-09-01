import { task, types } from "hardhat/config"
import {
  TASK_INITIALIZE,
  TASK_AUTHORIZE,
  TASK_REGISTER,
  TASK_INITIALIZE_STAKING,
} from "@keep-network/random-beacon/export/tasks/initialize"
import {
  authorize,
  register,
} from "@keep-network/random-beacon/export/tasks/utils"

// Tasks for the ECDSA application.
const TASK_INITIALIZE_ECDSA = `${TASK_INITIALIZE}:ecdsa`
const TASK_AUTHORIZE_ECDSA = `${TASK_AUTHORIZE}:ecdsa`
const TASK_REGISTER_ECDSA = `${TASK_REGISTER}:ecdsa`

task(
  TASK_INITIALIZE,
  "Initializes staking and the ECDSA application for a staking provider and an operator"
).setAction(async (args, hre) => {
  // Initialize staking
  await hre.run(TASK_INITIALIZE_STAKING, args)
  // Initialize ECDSA
  await hre.run(TASK_INITIALIZE_ECDSA, args)
})

task(TASK_INITIALIZE_ECDSA, "Initializes operator for ECDSA")
  .addParam("owner", "Stake Owner address", undefined, types.string)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .addOptionalParam("authorizer", "Stake Authorizer", undefined, types.string)
  .addOptionalParam(
    "authorization",
    "Authorization amount (default: minimumAuthorization)",
    undefined,
    types.int
  )
  .setAction(async (args, hre) => {
    await hre.run(TASK_AUTHORIZE_ECDSA, args)
    await hre.run(TASK_REGISTER_ECDSA, args)
  })

task(TASK_AUTHORIZE_ECDSA, "Sets authorization for ECDSA")
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
    await authorize(
      hre,
      "WalletRegistry",
      args.owner,
      args.provider,
      args.authorizer,
      args.authorization
    )
  })

task(
  TASK_REGISTER_ECDSA,
  "Registers an operator for a staking provider in ECDSA"
)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .setAction(async (args, hre) => {
    await register(hre, "WalletRegistry", args.provider, args.operator)
  })
