import { PENDING_WITHDRAWAL_STATUS } from "../constants/constants"
import moment from "moment"

/**
 * @typedef {"none" | "pending" | "available_to_withdraw" | "expired"} PendingWithdrawalStatus
 */

/**
 *
 * @param {Number} withdrawalDelay - withdrawal delay in seconds
 * @param {Number} withdrawalTimeout - withdrawal timeout in seconds
 * @param {Number} withdrawalInitiatedTimestamp - unix timestamp of the
 * initiated withdrawal
 * @return {PendingWithdrawalStatus} - returns the status of the withdrawal
 */
export const getPendingWithdrawalStatus = (
  withdrawalDelay,
  withdrawalTimeout,
  withdrawalInitiatedTimestamp
) => {
  if (!withdrawalDelay || !withdrawalTimeout || !withdrawalInitiatedTimestamp) {
    return PENDING_WITHDRAWAL_STATUS.NONE
  }

  const currentDate = moment()
  const endOfWithdrawalDelayDate = moment
    .unix(withdrawalInitiatedTimestamp)
    .add(withdrawalDelay, "seconds")
  const endOfWithdrawalTimeoutDate = withdrawalInitiatedTimestamp
    ? moment.unix(withdrawalInitiatedTimestamp)
    : moment()
  endOfWithdrawalTimeoutDate
    .add(withdrawalDelay, "seconds")
    .add(withdrawalTimeout, "seconds")

  if (currentDate.isSameOrBefore(endOfWithdrawalDelayDate)) {
    return PENDING_WITHDRAWAL_STATUS.PENDING
  } else if (
    currentDate.isAfter(endOfWithdrawalDelayDate) &&
    currentDate.isSameOrBefore(endOfWithdrawalTimeoutDate)
  ) {
    return PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
  } else if (currentDate.isAfter(endOfWithdrawalTimeoutDate)) {
    return PENDING_WITHDRAWAL_STATUS.EXPIRED
  }

  return PENDING_WITHDRAWAL_STATUS.NONE
}
