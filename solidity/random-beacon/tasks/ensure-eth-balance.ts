import { task, types } from "hardhat/config"

import { parseValue, TASK_SEND_ETH } from "./send-eth"

const TASK_ENSURE_ETH_BALANCE = "ensure-eth-balance"

task(
  TASK_ENSURE_ETH_BALANCE,
  "Ensure addresses hold a minimum ether balance, top-up if needed"
)
  .addOptionalParam(
    "from",
    "Address to send value from",
    undefined,
    types.string
  )
  .addPositionalParam(
    "balance",
    'Minimum ether balance the addresses are expected to hold, e.g. "0.5 ether", "100 gwei"',
    undefined,
    types.string
  )
  .addVariadicPositionalParam(
    "addresses",
    "Addresses for which balance should be validated",
    undefined,
    types.string
  )
  .setAction(async (args, hre) => {
    const { ethers } = hre

    // FIXME: `validate` will fail for badly checksummed addresses
    // see: https://github.com/ethers-io/ethers.js/discussions/3261
    const addresses: Set<string> = new Set(
      Array.from(args.addresses).map(hre.helpers.address.validate)
    )

    // eslint-disable-next-line no-restricted-syntax
    for (const address of addresses) {
      const expectedBalance = parseValue(args.balance, hre)
      const currentBalance = await ethers.provider.getBalance(address)

      console.log(
        `current balance of ${address} is ${ethers.utils.formatEther(
          currentBalance
        )} ether`
      )

      if (currentBalance.lt(expectedBalance)) {
        const topUpAmount = expectedBalance.sub(currentBalance)

        await hre.run(TASK_SEND_ETH, {
          from: args.from,
          value: topUpAmount.toString(),
          to: address,
        })
      }
    }
  })
