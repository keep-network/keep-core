import BigNumber from "bignumber.js"
import { KEEP_TO_T_EXCHANGE_RATE_IN_WEI } from "../constants/constants"

export const toThresholdTokenAmount = (keepAmount) => {
  const floatingPointDivisor = new BigNumber(10).pow(15)
  const amountInBN = new BigNumber(keepAmount)
  const wrappedRemainder = amountInBN.modulo(floatingPointDivisor)
  const convertibleAmount = amountInBN.minus(wrappedRemainder)

  return convertibleAmount
    .multipliedBy(KEEP_TO_T_EXCHANGE_RATE_IN_WEI)
    .dividedBy(floatingPointDivisor)
    .toString()
}
