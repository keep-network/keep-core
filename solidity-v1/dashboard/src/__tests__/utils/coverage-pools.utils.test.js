import moment from "moment"
import { getPendingWithdrawalStatus } from "../../utils/coverage-pools.utils"
import { PENDING_WITHDRAWAL_STATUS } from "../../constants/constants"

describe("Test getPendingWithdrawalStatus function", () => {
  let currentDate
  let withdrawalDelay
  let withdrawalTimeout

  beforeEach(() => {
    currentDate = moment()
    withdrawalDelay = 1814400 // 21 days
    withdrawalTimeout = 259200 // 3 days
  })

  it("Should return status pending when current date is before end of withdrawal delay date", () => {
    const withdrawalInitiatedTimestamp = currentDate.subtract(1, "days").unix()

    const result = getPendingWithdrawalStatus(
      withdrawalDelay,
      withdrawalTimeout,
      withdrawalInitiatedTimestamp
    )

    expect(result).toBe(PENDING_WITHDRAWAL_STATUS.PENDING)
  })

  it("Should return status available to withdraw when current date is after end of withdrawal delay date but before end of withdrawal timeout date", () => {
    const withdrawalInitiatedTimestamp = currentDate
      .subtract(withdrawalDelay, "seconds")
      .subtract(1, "days")
      .unix()

    const result = getPendingWithdrawalStatus(
      withdrawalDelay,
      withdrawalTimeout,
      withdrawalInitiatedTimestamp
    )

    expect(result).toBe(PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW)
  })

  it("Should return status expire when current date is after end of timeout date", () => {
    const withdrawalInitiatedTimestamp = currentDate
      .subtract(withdrawalDelay, "seconds")
      .subtract(withdrawalTimeout, "seconds")
      .subtract(1, "days")
      .unix()

    const result = getPendingWithdrawalStatus(
      withdrawalDelay,
      withdrawalTimeout,
      withdrawalInitiatedTimestamp
    )

    expect(result).toBe(PENDING_WITHDRAWAL_STATUS.EXPIRED)
  })

  it("Should return status none when withdrawal initiated timestamp is not provided", () => {
    const result = getPendingWithdrawalStatus(
      withdrawalDelay,
      withdrawalTimeout
    )

    expect(result).toBe(PENDING_WITHDRAWAL_STATUS.NONE)
  })
})
