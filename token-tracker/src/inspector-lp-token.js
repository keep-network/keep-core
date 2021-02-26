#!/usr/bin/env node --experimental-modules

import { getPastEvents } from "./lib/ethereum-helper.js"
import { writeFileSync } from "fs"
import Context from "./lib/context.js"
import { getPairData } from "./lib/uniswap.js"
import BigNumber from "bignumber.js"

// abi https://etherscan.io/address/0xe6f19dab7d43317344282f803f8e8d240708174a#code
import KEEPETHTokenJson from "../artifacts/KEEP-ETH-UNI-V2-Token.json"
import LPRewardsKEEPETHJson from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"

const ONE_HUNDRED = new BigNumber(100)
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
const LP_KEEPETH_BALANCE_BY_STAKERS_PATH = "./tmp/lp-keepeth-stakers.json"
const KEEP_AMOUNT_IN_LP_BY_STAKERS_PATH = "./tmp/keep-in-lp-stakers.json"

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

  await calculateKeepForStaker()

  async function calculateKeepForStaker() {
    const keepInLpByStakers = {}
    const stakersBalances = await getBalanceForStakers()
    const pairData = await getPairData(LIQUIDITY_POOLS.KEEP_ETH.pair)
    for (const [stakerAddress, lpBalance] of Object.entries(stakersBalances)) {
      const keepInLPToken = await calcKeepTokenfromLPToken(
        new BigNumber(lpBalance),
        pairData
      )
      keepInLpByStakers[stakerAddress] = keepInLPToken.toString()
      console.info(`${stakerAddress}: ${lpBalance}: ${keepInLPToken}`)
    }

    writeFileSync(
      KEEP_AMOUNT_IN_LP_BY_STAKERS_PATH,
      JSON.stringify(keepInLpByStakers, null, 2)
    )
  }

  function toWei(amount) {
    return web3.utils.toWei(amount.toString())
  }

  // Calculation is based on https://uniswap.org/docs/v2/advanced-topics/understanding-returns/
  async function calcKeepTokenfromLPToken(lpBalance, pairData) {
    console.info(pairData)
    const uniswapTotalSupply = toWei(pairData.totalSupply)
    const keepLiquidityPool = new BigNumber(pairData.reserve0) // KEEP
    const ethLiquidityPool = new BigNumber(pairData.reserve1) // ETH

    console.info("reserve0 ", keepLiquidityPool.toString())
    console.info("reserve1 ", ethLiquidityPool.toString())

    const shareOfUniswapPool = ONE_HUNDRED.multipliedBy(lpBalance).div(
      uniswapTotalSupply
    )
    console.info("shareOfUniswapPool:: ", shareOfUniswapPool.toString())

    // const ethPrice =  reserve0.div(reserve1) // price in KEEP for 1 ETH
    // keep_liquidity_pool * eth_liquidity_pool from Uniswap
    // const constantProduct = reserve0.multipliedBy(reserve1)
    // console.info("constantProduct:: ", constantProduct.toString())

    // const ethLiquidityPool = (constantProduct.div(ethPrice)).squareRoot()
    // console.info("ethLiquidityPool:: ", ethLiquidityPool.toString())
    // const keepLiquidityPool = (constantProduct.multipliedBy(ethPrice)).squareRoot()
    // console.info("keepLiquidityPool:: ", keepLiquidityPool.toString())
    
    // const ethInLP = shareOfUniswapPool
    //   .div(ONE_HUNDRED)
    //   .multipliedBy(ethLiquidityPool)
    // console.info("ethInLP: ", ethInLP.toString())

    const keepInLP = shareOfUniswapPool
      .div(ONE_HUNDRED)
      .multipliedBy(keepLiquidityPool)
    console.info("keepInLP: ", keepInLP.toString())
    

    return keepInLP
  }

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
    let totalSupply = new BigNumber(0)
    const stakersBalances = {}
    const lpStakers = await getLPStakers(
      LIQUIDITY_POOLS.KEEP_ETH.lpRewardsKeepEth
    )

    // for (let i = 0; i < 10; i++) {
    for (let i = 0; i < lpStakers.length; i++) {
      const lpBalance = new BigNumber(
        await LPRewardsKEEPETH.methods.balanceOf(lpStakers[i]).call({}, blockS)
      )
      if (!lpBalance.isZero()) {
        stakersBalances[lpStakers[i]] = lpBalance.toString()
        totalSupply = totalSupply.plus(lpBalance)

        console.info(`Staker: ${lpStakers[i]} - balance: ${lpBalance}`)
      }
    }

    console.info(
      `writing token stakers balance to a file: ${LP_KEEPETH_BALANCE_BY_STAKERS_PATH}`
    )
    writeFileSync(
      LP_KEEPETH_BALANCE_BY_STAKERS_PATH,
      JSON.stringify(stakersBalances, null, 2)
    )

    console.info(`Total supply of LP KEEP-ETH: ${totalSupply.toString()}`)

    return stakersBalances
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
