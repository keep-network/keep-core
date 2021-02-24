#!/usr/bin/env node

// import yargs from "yargs"

// const options = usage("Usage: -n <name>").option("n", {
//   alias: "name",
//   describe: "Your name",
//   type: "string",
//   demandOption: true,
// }).argv

// const greeting = `Hello, ${options.name}!`

// import program from "commander"

import Context from "../src/lib/context.js"
import { Inspector } from "../src/inspector.js"
import { KeepTokenTruthSource } from "../src/truth-sources/keep-token.js"
import { TokenStakingTruthSource } from "../src/truth-sources/token-staking.js"
import { logger } from "../src/lib/winston.js"
import { mapToObject } from "../src/lib/map-helper.js"

export async function getTokenOwnership() {
  const context = await Context.initialize()
  logger.debug("Context initialized")

  const inspector = new Inspector(context)

  // TODO: Register more truth sources.
  // inspector.registerTruthSource(KeepTokenTruthSource)
  inspector.registerTruthSource(TokenStakingTruthSource)

  return await inspector.getOwnershipsAtBlock(10958363) // 11909645 // TODO: Update to correct value
}

logger.on("finish", function () {
  process.exit(0)
})

logger.info(JSON.stringify(mapToObject(await getTokenOwnership())))

logger.end()
