#!/usr/bin/env node --experimental-modules

import BN from "bn.js"
import { getPastEvents } from "./lib/ethereum-helper.js"
import { writeFileSync } from "fs"
import Context from "./lib/context.js"

// abi https://etherscan.io/address/0xe6f19dab7d43317344282f803f8e8d240708174a#code
import KEEPETHTokenJson from "../artifacts/KEEP-ETH-UNI-V2-Token.json"
import LPRewardsKEEPETHJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"

const LIQUIDITY_POOLS = {
  KEEP_ETH: {
    pair: "0xe6f19dab7d43317344282f803f8e8d240708174a", // https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a
    lpRewardsKeepEth: "0x47A5f2ffdf66D13ED7e317581F458d09b49d6F44", // https://etherscan.io/address/0x47a5f2ffdf66d13ed7e317581f458d09b49d6f44
  },
  KEEP_TBTC: {
    address: "0x38c8ffee49f286f25d25bad919ff7552e5daf081", // https://info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081
    lpRewardsKeepTbtc: "0xb3d03a5411261fc2094697c5e969d552ee55cf6b", // https://etherscan.io/address/0xb3d03a5411261fc2094697c5e969d552ee55cf6b
  },
}
const LP_KEEPETH_STAKERS_BALANCE_PATH = "./tmp/lp-keepeth-stakers.json"

// Block height when stakedrop happens. Change accordingly
const blockS = "latest"
const KEEPETHCreationBlock = "10100034" // https://etherscan.io/tx/0xc64ac175846e719bb4f7f9b17a0b04bc365db3dda9d97ef70d7ede8f9c1a265b

async function run() {
  const context = await Context.initialize()
  const web3 = context.web3

  const KEEPETHPairTokenAbi = JSON.parse(KEEPETHTokenJson.result)
  const KEEPETHPairToken = new web3.eth.Contract(
    KEEPETHPairTokenAbi,
    LIQUIDITY_POOLS.KEEP_ETH.pair
  )

  const LPRewardsKEEPETHAbi = LPRewardsKEEPETHJson.abi
  const LPRewardsKEEPETH = new web3.eth.Contract(
    LPRewardsKEEPETHAbi,
    LIQUIDITY_POOLS.KEEP_ETH.lpRewardsKeepEth
  )

  await getBalanceForStakers()

  async function getLPStakers(lpRewardAddress) {
    const events = await getPastEvents(
      web3,
      KEEPETHPairToken,
      "Transfer",
      KEEPETHCreationBlock,
      blockS
    )

    console.info(`found ${events.length} token transfer events`)

    const lpTokenStakersSet = new Set()
    events.forEach(function (event) {
      // include accounts that staked in LPRewards contracts only
      if (event.returnValues.to == lpRewardAddress) {
        lpTokenStakersSet.add(event.returnValues.from)
      }
    })

    const lpTokenStakers = Array.from(lpTokenStakersSet)
    console.info(`found ${lpTokenStakers.length} unique historic LP stakers`)

    return lpTokenStakers
  }

  async function getBalanceForStakers() {
    let totalSupply = new BN(0)
    const stakerBalance = {}
    const lpStakers = await getLPStakers(
      LIQUIDITY_POOLS.KEEP_ETH.lpRewardsKeepEth
    )

    for (let i = 0; i < lpStakers.length; i++) {
      const lpBalance = new BN(
        await LPRewardsKEEPETH.methods.balanceOf(lpStakers[i]).call({}, blockS)
      )
      if (!lpBalance.isZero()) {
        stakerBalance[lpStakers[i]] = lpBalance.toString()
        totalSupply = totalSupply.add(lpBalance)

        console.info(`Staker: ${lpStakers[i]} - balance: ${lpBalance}`)
      }
    }

    console.info(
      `writing token stakers balance to a file: ${LP_KEEPETH_STAKERS_BALANCE_PATH}`
    )
    writeFileSync(
      LP_KEEPETH_STAKERS_BALANCE_PATH,
      JSON.stringify(stakerBalance, null, 2)
    )

    console.info(`Total supply LP KEEP-ETH: ${totalSupply.toString()}`)
  }
}

run()
  .then(() => {
    console.info("Retrieving of LP token staker balance completed successfully")

    process.exit(0)
  })
  .catch((error) => {
    console.error("Retrieving of LP token staker balance errored out: ", error)

    process.exit(1)
  })
