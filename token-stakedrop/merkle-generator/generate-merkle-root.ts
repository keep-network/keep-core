import { parseBalanceMap } from "../merkle-distributor/src/parse-balance-map"
import { program } from "commander"
import * as fs from "fs"

program
  .version("0.0.0")
  .requiredOption(
    "-i, --input <path>",
    "input JSON file location containing a map of account addresses to string balances"
  )

program.parse(process.argv)

const outputMerkleObject = "./output-merkle-object.json"

const merkleObject = JSON.parse(
  fs.readFileSync(program.input, { encoding: "utf8" })
)

if (typeof merkleObject !== "object") throw new Error("Invalid JSON")

fs.writeFileSync(
  outputMerkleObject,
  JSON.stringify(parseBalanceMap(merkleObject), null, 2)
)
