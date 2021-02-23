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

export async function getTokenOwnership() {
  const context = await Context.initialize()
  console.debug("Context initialized")

  const inspector = new Inspector(context)

  // TODO: Register more truth sources.
  // inspector.registerTruthSource(KeepTokenTruthSource)
  inspector.registerTruthSource(TokenStakingTruthSource)

  return await inspector.getOwnershipsAtBlock(10958363) // 11909645 // TODO: Update to correct value
}

console.log(JSON.stringify(await getTokenOwnership()))

process.exit(0)
