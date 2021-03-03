#!/usr/bin/env NODE_BACKEND=js node --experimental-modules --experimental-json-modules

import { existsSync, mkdirSync } from "fs"
import path from "path"

import Context from "../src/lib/context.js"
import { logger } from "../src/lib/winston.js"
import { dumpDataToFile } from "../src/lib/file-helper.js"
import { getDeploymentBlockNumber } from "../src/lib/contract-helper.js"
import { Inspector } from "../src/inspector.js"

import { KeepTokenTruthSource } from "../src/truth-sources/keep-token.js"
import { TokenStakingTruthSource } from "../src/truth-sources/token-staking.js"
import { TokenGrantTruthSource } from "../src/truth-sources/token-grant.js"
import { LPTokenTruthSource } from "../src/truth-sources/lp-tokens.js"

import KeepTokenJson from "@keep-network/keep-core/artifacts/KeepToken.json"

import commander from "commander"
const program = new commander.Command()

program
  .requiredOption("--target-block <number>", "target block number")
  .parse(process.argv)

const TMP_DIR = "./tmp"
const OUT_DIR = "./output"
const RESULT_OUTPUT_PATH = path.join(OUT_DIR, "result.json")

async function initializeContext() {
  const context = await Context.initialize(
    process.env.ETH_HOSTNAME,
    // Initializes web3 in the read only mode. Provide actual private key to
    // interact with the chain.
    "01".repeat(32)
  )
  logger.debug("context initialized")

  // FIXME: We can get rid of global deployment block tracking and switch to
  // particular contracts deployment blocks if we use tbtc.js like functions
  // to get past events with a contract instance defining deployment block.
  context.deploymentBlock = await getDeploymentBlockNumber(
    KeepTokenJson,
    context.web3
  )

  logger.debug(`deployment block: ${context.deploymentBlock}`)

  return context
}

export async function getTokenOwnership(targetBlockNumber) {
  if (!targetBlockNumber) throw new Error("target block is not defined")

  logger.info(`Inspect token ownership at block ${targetBlockNumber}`)

  const context = await initializeContext()

  const inspector = new Inspector(context)

  // TODO: Register more truth sources.
  inspector.registerTruthSource(KeepTokenTruthSource)
  inspector.registerTruthSource(TokenStakingTruthSource)
  inspector.registerTruthSource(TokenGrantTruthSource)
  inspector.registerTruthSource(LPTokenTruthSource)

  return await inspector.getOwnershipsAtBlock(targetBlockNumber)
}

async function run() {
  if (!existsSync(TMP_DIR)) {
    mkdirSync(TMP_DIR)
  }
  if (!existsSync(OUT_DIR)) {
    mkdirSync(OUT_DIR)
  }

  const result = await getTokenOwnership(program.opts().targetBlock)

  dumpDataToFile(result, RESULT_OUTPUT_PATH)

  logger.info("DONE!")
}

logger.on("finish", function () {
  process.exit(0)
})

run()
  .then(() => logger.end())
  .catch((err) => {
    logger.error(err)
    logger.end()
  })
