import { contractService } from "./contracts.service"
import {
  TBTC_SYSTEM_CONTRACT_NAME,
  DEPOSIT_STATES,
} from "../constants/constants"
import moment from "moment"
import web3Utils from "web3-utils"

/**
 * Import contract directly for constants
 */
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
// import TBTC from "tbtc.js"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  createDepositContractInstance,
  getTBTCConstantsContract,
  ContractsLoaded
} from "../contracts"

const getPastRedemptionRequestedEvents = async (web3Context) => {
  const pastRedemptionRequestedEvents = await contractService.getPastEvents(
    web3Context,
    TBTC_SYSTEM_CONTRACT_NAME,
    "RedemptionRequested",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_SYSTEM_CONTRACT_NAME],
      // filter: { attr : paramVal },
    }
  )
  const augmentedRedemRequestEvents = await Promise.all(
    pastRedemptionRequestedEvents.map(augmentEventDepositInfo(web3Context))
  )
  return augmentedRedemRequestEvents
}

const getPastRedemptionSignatureEvents = async (web3Context) => {
  const pastEvents = await contractService.getPastEvents(
    web3Context,
    TBTC_SYSTEM_CONTRACT_NAME,
    "GotRedemptionSignature",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_SYSTEM_CONTRACT_NAME],
      // filter: { attr : paramVal },
    }
  )
  const tBtcSystemEventsWithDepositInfo = await Promise.all(
    pastEvents.map(augmentEventDepositInfo(web3Context)),
    pastEvents.map(augmentEventPatchReturnValueWithUtxoVal(web3Context))
  )
  return tBtcSystemEventsWithDepositInfo
}

const getPastCourtesyCalledEvents = async (web3Context) => {
  const pastEvents = await contractService.getPastEvents(
    web3Context,
    TBTC_SYSTEM_CONTRACT_NAME,
    "CourtesyCalled",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_SYSTEM_CONTRACT_NAME],
      // filter: { attr : paramVal },
    }
  )
  const tBtcSystemEventsWithDepositInfo = await Promise.all(
    pastEvents.map(augmentEventDepositInfo(web3Context))
  )
  return tBtcSystemEventsWithDepositInfo
}

const getLastStartedLiquidationEvent = async (web3Context, prmDepositContractAdress) => {
  const pastStartedLiquidationEvents = await contractService.getPastEvents(
    web3Context,
    TBTC_SYSTEM_CONTRACT_NAME,
    "StartedLiquidation",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_SYSTEM_CONTRACT_NAME],
      filter: { _depositContractAddress: prmDepositContractAdress },
    }
  )
  const sortedStartedLiquidationEvent = pastStartedLiquidationEvents.sort((a, b) => a.returnValues._timestamp - b.returnValues._timestamp)
  return sortedStartedLiquidationEvent[0]
}

const augmentEventDepositInfo = (web3Context) => {
  return async (w3event) => {
    let evDepStatus = null
    try {
      evDepStatus = await getDepositState(
        web3Context,
        w3event.returnValues._depositContractAddress
      )
    } catch (error) {
      console.error(error)
      throw error
    }
    w3event.depositStatusObj = evDepStatus
    w3event.isAlignedToDeposit = isEventAlignedToDeposit(w3event, evDepStatus)
    w3event.timestamp = await eventTimestamp(web3Context, w3event)

    return w3event
  }
}
const augmentEventPatchReturnValueWithUtxoVal = (web3Context) => {
  return async (w3event) => {
    let evDepSizeSatoshis = null
    try {
      evDepSizeSatoshis = await getDepositSizeSatoshis(
        web3Context,
        w3event.returnValues._depositContractAddress
      )
    } catch (error) {
      console.error(error)
      throw error
    }
    w3event.returnValues._utxoValue = evDepSizeSatoshis
    return w3event
  }
}

const depositNotifySignatureTimeout = async (
  web3Context,
  depositContractAddress
) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context.web3,
    depositContractAddress
  )
  depositContractInstance.methods
    .notifyRedemptionSignatureTimedOut()
    .send({
      from: web3Context.yourAddress,
    })
    .on("receipt", function (receipt) {
      console.log(
        `notifySignatureTimeout on Deposit ${depositContractAddress} was successful.`
      )
      return receipt
    })
    .on("error", function (error, receipt) {
      console.log(error.message)
      console.log(
        `notifySignatureTimeout on Deposit ${depositContractAddress} failed.`
      )
      throw error
    })
}

const depositNotifyRedemptionProofTimedOut = async (
  web3Context,
  depositContractAddress
) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context.web3,
    depositContractAddress
  )
  depositContractInstance.methods
    .notifyRedemptionProofTimedOut()
    .send({
      from: web3Context.yourAddress,
    })
    .on("receipt", function (receipt) {
      console.log(
        `notifyRedemptionProofTimedOut on Deposit ${depositContractAddress} was successful.`
      )
      return receipt
    })
    .on("error", function (error, receipt) {
      console.log(error.message)
      console.log(
        `notifyRedemptionProofTimedOut on Deposit ${depositContractAddress} failed.`
      )
      throw error
    })
}

const eventTimestamp = async (web3Context, event) => {
  const { eth } = web3Context
  const block = await eth.getBlock(event.blockNumber)
  return moment.unix(block.timestamp)
}

const getDepositState = async (web3Context, depositContractAddress) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context.web3,
    depositContractAddress
  )
  const depStateCode = await depositContractInstance.methods
    .currentState()
    .call()
  const depStateObj = DEPOSIT_STATES.find(
    (obj) => obj.depositStatusId.toString() === depStateCode.toString()
  )
  return depStateObj
}

const getDepositSizeSatoshis = async (web3Context, depositContractAddress) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context.web3,
    depositContractAddress
  )
  const depositLotSizeSatoshis = await depositContractInstance.methods
    .lotSizeSatoshis()
    .call()
  const depositUtxoValue = await depositContractInstance.methods
    .utxoValue()
    .call()
  if (`${depositLotSizeSatoshis}` !== `${depositUtxoValue}`) {
    console.log(`--Diff in deposit size-----\n${depositLotSizeSatoshis}\n${depositUtxoValue}\n--Diff in deposit size-----\n`)
  }
  return depositUtxoValue
}

/**
 * 
 * @param {*} web3Context the Web3 context from the page
 * @param {*} depositContractAddress the address of a Deposit contract. This call will revert if the deposit is not in a state where an auction is currently in progress.
 * 
 * @returns a BN instance of the current auction value in ETH wei
 */
const getDepositCurrentAuctionValue = async (
  web3Context,
  depositContractAddress
) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context.web3,
    depositContractAddress
  )
  const currentAuctionValue = await depositContractInstance.methods
    .auctionValue()
    .call()
  return web3Utils.toBN(currentAuctionValue)
}

const purchaseDepositAtAuction = async (
  web3Context,
  depositContractAddress,
  onTransactionHashCallback = () => { }) => {

  const { web3, yourAddress } = web3Context

  const depositContractInstance = createDepositContractInstance(
    web3,
    depositContractAddress
  )
  console.log(depositContractInstance)

  try {
    const { tbtcTokenContract } = await ContractsLoaded
    const allowance = await tbtcTokenContract
      .methods.allowance(yourAddress, depositContractAddress)
      .call()
    const allowanceBN = web3Utils.toBN(allowance)
    const tBtcRequired = web3Utils.toBN(await depositContractInstance.methods.lotSizeTbtc().call())

    
    if (allowanceBN.lt(tBtcRequired)) {
      await tbtcTokenContract.methods.approve(depositContractAddress, tBtcRequired.toString())
        .send({ from: yourAddress })
    }
    //Estimate gas from tbtc.js sendSafely()
    const gasEstimate = await depositContractInstance.methods.purchaseSignerBondsAtAuction().estimateGas({})

    await depositContractInstance.methods
      .purchaseSignerBondsAtAuction()
      .send({
        from: yourAddress,
        gas: gasEstimate
      })
    
    await depositContractInstance.methods
      .withdrawFunds()
      .send({
        from: yourAddress
      })
      .on("transactionHash", onTransactionHashCallback)

  } catch (error) {
    console.log(error)
  }
}

/**
 * 
 * @param {*} web3Context 
 * @param {*} depositContractAddress the address of a Deposit contract.
 * 
 * @returns {*} Returns the the deposit's ETH bond balance in wei.
 */
const getDepositEthBalance = async (
  web3Context,
  depositContractAddress
) => {
  const { web3 } = web3Context
  const depositBondBalance = await web3.eth.getBalance(depositContractAddress)
  // console.log(`depositBondBalance: ${depositBondBalance}`)
  return depositBondBalance
}

/**
 * Inspiration: https://github.com/keep-network/tbtc/blob/master/solidity/contracts/deposit/DepositUtils.sol#L239-L252
 * Parameters (not for the function, but for calculation):
 * - Base Auction Percentage
 * - Auction Duration
 * Remarks:
 * (1) In the smart contract, percentage on auction is the result of an integer division 
 * rounding down. With an initialCollateralizedPercent of 150, _basePercentage should always be 66 (100/150).
 * 
 */
const getDepositAuctionOfferingSchedule = async (
  web3Context,
  depositContractAddress,
  startFromAuctionPct = null
) => {
  const startedLiquidationEvent = await getLastStartedLiquidationEvent(web3Context, depositContractAddress)
  const startedLiquidationTimestamp = startedLiquidationEvent.returnValues._timestamp
  const mmtStartedLiquidation = moment.unix(startedLiquidationTimestamp)

  //This could go in context
  const tBtcConstantsContract = await getTBTCConstantsContract(web3Context.web3)
  const baseAuctionDurationSeconds = await tBtcConstantsContract.methods.getAuctionDuration().call()

  const depositContractInstance = createDepositContractInstance(
    web3Context.web3,
    depositContractAddress
  )
  const initialCollateralizedPercent = await depositContractInstance.methods
    .initialCollateralizedPercent()
    .call() //uint16 
  const baseAuctionPercentage = web3Utils.toBN(10000).div(web3Utils.toBN(initialCollateralizedPercent)) // Remarks (1)
  console.log(`initialCollateralizedPercent: ${initialCollateralizedPercent}`)
  console.log(`baseAuctionPercentage: ${baseAuctionPercentage}`)

  const depositEthBalance = await getDepositEthBalance(web3Context, depositContractAddress)
  const bnHundred = web3Utils.toBN(100)
  let offeringSchedule = []

  offeringSchedule[0] = {
    depositPctOnOffer: baseAuctionPercentage,
    amountOnOffer: web3Utils.toBN(depositEthBalance).mul(baseAuctionPercentage).div(bnHundred),
    releasedInTimestamp: startedLiquidationTimestamp
  }
  console.log(`offeringSchedule[0]`)

  for (let index = baseAuctionPercentage.toNumber() + 1; index <= 100; index++) {
    const releaseSlotIndex = index - baseAuctionPercentage
    const jsNumberReleaseProgress = (releaseSlotIndex / (bnHundred - baseAuctionPercentage)).toPrecision(18) // specific precision aligned to power function to obtain whole number
    const intReleaseProgress = jsNumberReleaseProgress * Math.pow(10, 18) // to obtain whole number
    const bnReleaseProgress = web3Utils.toBN(intReleaseProgress)
    const bnElapsedFromPctOffered = bnReleaseProgress.mul(web3Utils.toBN(baseAuctionDurationSeconds))
    const secondsElapsedFromPctOffered = bnElapsedFromPctOffered.div(web3Utils.toBN(10).pow(web3Utils.toBN(18))) // rebase to real scale for seconds
    offeringSchedule[releaseSlotIndex] = {
      depositPctOnOffer: index,
      amountOnOffer: web3Utils.toBN(depositEthBalance).mul(web3Utils.toBN(index)).div(bnHundred),
      releasedInTimestamp: web3Utils.toBN(startedLiquidationTimestamp).add(secondsElapsedFromPctOffered).toString()
    }
  }

  return offeringSchedule
}

/**
 * 
 * @param {*} web3Context 
 * 
 * @returns {*} Returns the current user's tBTC balance in uint256.
 */
const getTBtcBalanceOf = async (
  web3Context
) => {
  const arbitrageurAddress = web3Context.yourAddress
  console.log(`arbitrageurAddress: ${arbitrageurAddress}`)
  const { tbtcTokenContract } = await ContractsLoaded
  const arbitrageurBalanceOf = await tbtcTokenContract.methods
    .balanceOf(arbitrageurAddress)
    .call()
  console.log(`arbitrageurBalanceOf: ${arbitrageurBalanceOf}`)
  return arbitrageurBalanceOf
}

const isEventAlignedToDeposit = (tbtcSystemContractEvent, evDepStatus) => {
  const contract = TBTCSystem
  let isAligned = false
  switch (evDepStatus.name) {
    case DEPOSIT_STATES.find((dso) => dso.depositStatusId === 5).name: // AWAITING_WITHDRAWAL_SIGNATURE
      if (tbtcSystemContractEvent.event === contract.abi[20].name)
        // RedemptionRequested
        isAligned = true
      break
    case DEPOSIT_STATES.find((dso) => dso.depositStatusId === 6).name: // AWAITING_WITHDRAWAL_PROOF
      if (tbtcSystemContractEvent.event === contract.abi[12].name)
        // GotRedemptionSignature
        isAligned = true
      break
    default:
      isAligned = true
  }
  return isAligned
}

export const liquidationService = {
  getPastRedemptionRequestedEvents,
  getPastRedemptionSignatureEvents,
  getPastCourtesyCalledEvents,
  getLastStartedLiquidationEvent,
  depositNotifySignatureTimeout,
  depositNotifyRedemptionProofTimedOut,
  getDepositState,
  getDepositSizeSatoshis,
  getDepositCurrentAuctionValue,
  getDepositEthBalance,
  getTBtcBalanceOf,
  getDepositAuctionOfferingSchedule,
  purchaseDepositAtAuction
}
