#!/usr/bin/env node --experimental-modules

import { getPastEvents } from "./lib/ethereum-helper.js"
import { writeFileSync } from "fs"
import Context from "./lib/context.js"
import { getPairData } from "./lib/uniswap.js"
import BigNumber from "bignumber.js"

// abi https://etherscan.io/address/0xe6f19dab7d43317344282f803f8e8d240708174a#code
import KEEPETHTokenJson from "../artifacts/KEEP-ETH-UNI-V2-Token.json"
// abi https://etherscan.io/address/0x38c8ffee49f286f25d25bad919ff7552e5daf081#code
import KEEPTBTCTokenJson from "../artifacts/KEEP-TBTC-UNI-V2-Token.json"

import LPRewardsKEEPETHJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"
import LPRewardsKEEPTBTCJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPTBTC.json"

const LIQUIDITY_POOLS = {
  KEEP_ETH: {
    pair: "0xe6f19dab7d43317344282f803f8e8d240708174a", // https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a
    lpRewardsContractAddress: "0x47A5f2ffdf66D13ED7e317581F458d09b49d6F44", // https://etherscan.io/address/0x47a5f2ffdf66d13ed7e317581f458d09b49d6f44
  },
  KEEP_TBTC: {
    pair: "0x38c8ffee49f286f25d25bad919ff7552e5daf081", // https://info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081
    lpRewardsContractAddress: "0xb3d03A5411261fC2094697C5e969D552eE55cF6B", // https://etherscan.io/address/0xb3d03a5411261fc2094697c5e969d552ee55cf6b
  },
}
const KEEP_AMOUNT_IN_KEEPETH_BY_STAKERS_PATH =
  "./tmp/keep-in-keepeth-by-stakers.json"
const KEEP_AMOUNT_IN_KEEPTBTC_BY_STAKERS_PATH =
  "./tmp/keep-in-keeptbtc-by-stakers.json"

// Block height when stakedrop happens. Change accordingly
const BLOCK_S = "latest"
const KEEPETH_CREATION_BLOCK = "10100034" // https://etherscan.io/tx/0xc64ac175846e719bb4f7f9b17a0b04bc365db3dda9d97ef70d7ede8f9c1a265b
const KEEPTBTC_CREATION_BLOCK = "11452642" // https://etherscan.io/tx/0x1592f9b235c602c87a5b8cc5f896164dc43d16b92664cb9c8b420d28b64ca4a0

async function run() {
  const context = await Context.initialize()
  const web3 = context.web3

  const KEEPETHPairTokenAbi = JSON.parse(KEEPETHTokenJson.result)
  const KEEPETHPairToken = new web3.eth.Contract(
    KEEPETHPairTokenAbi,
    LIQUIDITY_POOLS.KEEP_ETH.pair
  )

  const KEEPTBTCPairTokenAbi = JSON.parse(KEEPTBTCTokenJson.result)
  const KEEPTBTCPairToken = new web3.eth.Contract(
    KEEPTBTCPairTokenAbi,
    LIQUIDITY_POOLS.KEEP_TBTC.pair
  )

  const LPRewardsKEEPETHAbi = LPRewardsKEEPETHJson.abi
  const LPRewardsKEEPETH = new web3.eth.Contract(
    LPRewardsKEEPETHAbi,
    LIQUIDITY_POOLS.KEEP_ETH.lpRewardsContractAddress
  )

  const LPRewardsKEEPTBTCAbi = LPRewardsKEEPTBTCJson.abi
  const LPRewardsKEEPTBTC = new web3.eth.Contract(
    LPRewardsKEEPTBTCAbi,
    LIQUIDITY_POOLS.KEEP_TBTC.lpRewardsContractAddress
  )

  const keepInKEEPETHByStakers = await calculateKeepForStaker(
    LIQUIDITY_POOLS.KEEP_ETH.pair
  )
  const keepInKEEPTBTCByStakers = await calculateKeepForStaker(
    LIQUIDITY_POOLS.KEEP_TBTC.pair
  )
  writeToJSON(LIQUIDITY_POOLS.KEEP_ETH.pair, keepInKEEPETHByStakers)
  writeToJSON(LIQUIDITY_POOLS.KEEP_TBTC.pair, keepInKEEPTBTCByStakers)

  async function calculateKeepForStaker(pair) {
    const keepInLpByStakers = {}

    let stakersBalances = {}
    if (pair == LIQUIDITY_POOLS.KEEP_ETH.pair) {
      stakersBalances = await getLPBalanceForStakers(
        LIQUIDITY_POOLS.KEEP_ETH.lpRewardsContractAddress
      )
    } else {
      stakersBalances = await getLPBalanceForStakers(
        LIQUIDITY_POOLS.KEEP_TBTC.lpRewardsContractAddress
      )
    }

    // Retrieve current pair data. For historic data, need to pass a block number.
    const pairData = await getPairData(pair)
    const decimals = new BigNumber(10).pow(new BigNumber(18))
    for (const [stakerAddress, lpBalance] of Object.entries(stakersBalances)) {
      const lpBalanceBN = new BigNumber(lpBalance)
      const keepInLPToken = await calcKeepTokenfromLPToken(
        lpBalanceBN,
        pairData
      )
      keepInLpByStakers[stakerAddress] = keepInLPToken.toString()

      console.info(
        `Staker: ${stakerAddress} - LP Balance: ${lpBalanceBN.div(
          decimals
        )} - KEEP in LP: ${keepInLPToken}`
      )
    }

    return keepInLpByStakers
  }

  async function getLPBalanceForStakers(lpRewardContractAddress) {
    let totalSupply = new BigNumber(0)
    const stakersBalances = {}
    const lpStakers = await getLPStakers(lpRewardContractAddress)

    let lpRewardContract = {}
    if (
      lpRewardContractAddress ==
      LIQUIDITY_POOLS.KEEP_ETH.lpRewardsContractAddress
    ) {
      lpRewardContract = LPRewardsKEEPETH
    } else {
      lpRewardContract = LPRewardsKEEPTBTC
    }

    for (let i = 0; i < lpStakers.length; i++) {
      const lpBalance = new BigNumber(
        await lpRewardContract.methods.balanceOf(lpStakers[i]).call({}, BLOCK_S)
      )
      if (!lpBalance.isZero()) {
        stakersBalances[lpStakers[i]] = lpBalance.toString()
        totalSupply = totalSupply.plus(lpBalance)
      }
    }

    console.info(`Total supply of LP Token: ${totalSupply.toString()}`)

    return stakersBalances
  }

  async function getLPStakers(lpRewardContract) {
    let events = {}
    if (lpRewardContract == LIQUIDITY_POOLS.KEEP_ETH.lpRewardsContractAddress) {
      events = await getPastEvents(
        web3,
        KEEPETHPairToken,
        "Transfer",
        KEEPETH_CREATION_BLOCK,
        BLOCK_S
      )
    } else {
      events = await getPastEvents(
        web3,
        KEEPTBTCPairToken,
        "Transfer",
        KEEPTBTC_CREATION_BLOCK,
        BLOCK_S
      )
    }

    console.info(`found ${events.length} token transfer events`)

    const lpTokenStakersSet = new Set()
    events.forEach(function (event) {
      // include accounts that staked in LPRewards contracts only
      if (event.returnValues.to == lpRewardContract) {
        lpTokenStakersSet.add(event.returnValues.from)
      }
    })

    const lpTokenStakers = Array.from(lpTokenStakersSet)
    console.info(`found ${lpTokenStakers.length} unique historic LP stakers`)

    return lpTokenStakers
  }

  // Calculation is based on https://uniswap.org/docs/v2/advanced-topics/understanding-returns/
  async function calcKeepTokenfromLPToken(lpBalance, pairData) {
    const uniswapTotalSupply = web3.utils.toWei(pairData.totalSupply.toString())
    const keepLiquidityPool = new BigNumber(pairData.reserve0) // KEEP

    return lpBalance.multipliedBy(keepLiquidityPool).div(uniswapTotalSupply)
  }

  function writeToJSON(pair, data) {
    let keepInLpPath = ""
    if (pair == LIQUIDITY_POOLS.KEEP_ETH.pair) {
      console.info(
        `writing KEEP amount in LP token ${pair} for stakers to a file: ${KEEP_AMOUNT_IN_KEEPETH_BY_STAKERS_PATH}`
      )
      keepInLpPath = KEEP_AMOUNT_IN_KEEPETH_BY_STAKERS_PATH
    } else {
      console.info(
        `writing KEEP amount in LP token ${pair} for stakers to a file: ${KEEP_AMOUNT_IN_KEEPTBTC_BY_STAKERS_PATH}`
      )
      keepInLpPath = KEEP_AMOUNT_IN_KEEPTBTC_BY_STAKERS_PATH
    }

    writeFileSync(keepInLpPath, JSON.stringify(data, null, 2))
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
