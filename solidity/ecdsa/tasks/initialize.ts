import { task, types } from "hardhat/config"
import {
  TASK_INITIALIZE,
  TASK_AUTHORIZE,
  TASK_REGISTER,
} from "@keep-network/random-beacon/tasks/initialize"
import { authorize, register } from "@keep-network/random-beacon/tasks/utils"

import type { HardhatRuntimeEnvironment } from "hardhat/types"

const TASK_AUTHORIZE_ECDSA = `${TASK_AUTHORIZE}:ecdsa`
const TASK_REGISTER_ECDSA = `${TASK_REGISTER}:ecdsa`

task(TASK_INITIALIZE, "Initializes staking for an operator").setAction(
  async (args, hre, runSuper) => {
    // Run initialization task from @keep-network/random-beacon.
    await runSuper(args)
    // Initialize ECDSA.
    await initializeEcdsa(hre, args)
  }
)

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
    await authorize(hre, "WalletRegistry", args)
  })

task(
  TASK_REGISTER_ECDSA,
  "Registers an operator for a staking provider in ECDSA"
)
  .addParam("provider", "Staking Provider", undefined, types.string)
  .addParam("operator", "Staking Operator", undefined, types.string)
  .setAction(async (args, hre) => {
    await register(hre, "WalletRegistry", args)
  })

// eslint-disable-next-line import/prefer-default-export
export async function initializeEcdsa(
  hre: HardhatRuntimeEnvironment,
  args
): Promise<void> {
  await hre.run(TASK_AUTHORIZE_ECDSA, args)
  await hre.run(TASK_REGISTER_ECDSA, args)
}
