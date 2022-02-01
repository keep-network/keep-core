import BigNumber from "bignumber.js"
import { KEEP_TO_T_EXCHANGE_RATE_IN_WEI } from "../constants/constants"

const floatingPointDivisor = new BigNumber(10).pow(15)

export const toThresholdTokenAmount = (keepAmount) => {
  const amountInBN = new BigNumber(keepAmount)
  const wrappedRemainder = amountInBN.modulo(floatingPointDivisor)
  const convertibleAmount = amountInBN.minus(wrappedRemainder)

  return convertibleAmount
    .multipliedBy(KEEP_TO_T_EXCHANGE_RATE_IN_WEI)
    .dividedBy(floatingPointDivisor)
    .toString()
}

export const fromThresholdTokenAmount = (tAmount) => {
  const amountInBN = new BigNumber(tAmount)
  const tRemainder = amountInBN.modulo(KEEP_TO_T_EXCHANGE_RATE_IN_WEI)
  const convertibleAmount = amountInBN.minus(tRemainder)

  return convertibleAmount
    .multipliedBy(floatingPointDivisor)
    .dividedBy(KEEP_TO_T_EXCHANGE_RATE_IN_WEI)
    .toString()
}
