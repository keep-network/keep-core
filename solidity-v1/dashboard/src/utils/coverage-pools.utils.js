import moment from "moment"
import { PENDING_WITHDRAWAL_STATUS } from "../constants/constants"

/**
 *
 * @param {Number} withdrawalDelay - withdrawal delay in seconds
 * @param {Number} withdrawalTimeout - withdrawal timeout in seconds
 * @param {Number} withdrawalInitiatedTimestamp - unix timestamp of the
 * initiated withdrawal
 * @return {PENDING_WITHDRAWAL_STATUS} - returns the
 * status of the withdrawal
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
    return PENDING_WITHDRAWAL_STATUS.COMPLETED
  } else if (currentDate.isAfter(endOfWithdrawalTimeoutDate)) {
    return PENDING_WITHDRAWAL_STATUS.EXPIRED
  }

  return PENDING_WITHDRAWAL_STATUS.NONE
}
