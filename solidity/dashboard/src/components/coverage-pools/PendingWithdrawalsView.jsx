import React from "react"
import { Column, DataTable } from "../DataTable"
import resourceTooltipProps from "../../constants/tooltips"
import TokenAmount from "../TokenAmount"
import moment from "moment"
import { SubmitButton } from "../Button"
import { useSelector } from "react-redux"
import * as Icons from "../Icons"
import Chip from "../Chip"
import ProgressBar from "../ProgressBar"
import { colors } from "../../constants/colors"

const PendingWithdrawalsView = ({
  onClaimTokensSubmitButtonClick,
  onReinitiateWithdrawal,
}) => {
  const {
    withdrawalDelay,
    withdrawalTimeout,
    pendingWithdrawals,
  } = useSelector((state) => state.coveragePool)

  const renderProgressBar = (
    withdrawalDate,
    endOfWithdrawalDelayDate,
    currentDate
  ) => {
    const progressBarValueInSeconds = currentDate.diff(
      withdrawalDate,
      "seconds"
    )
    const progressBarTotalInSeconds = endOfWithdrawalDelayDate.diff(
      withdrawalDate,
      "seconds"
    )
    return (
      <ProgressBar
        value={progressBarValueInSeconds}
        total={progressBarTotalInSeconds}
        color={colors.secondary}
        bgColor={colors.bgSecondary}
      >
        <ProgressBar.Inline
          height={20}
          className={"pending-withdrawal__progress-bar"}
        />
      </ProgressBar>
    )
  }

  const renderCooldownStatus = (timestamp) => {
    const withdrawalDate = moment.unix(timestamp)
    const currentDate = moment()
    const endOfWithdrawalDelayDate = moment
      .unix(timestamp)
      .add(withdrawalDelay, "seconds")
    const days = endOfWithdrawalDelayDate.diff(currentDate, "days")
    const hours = moment
      .duration(endOfWithdrawalDelayDate.diff(currentDate))
      .hours()
    const minutes = moment
      .duration(endOfWithdrawalDelayDate.diff(currentDate))
      .minutes()

    let cooldownStatus = <></>
    if (currentDate.isBefore(endOfWithdrawalDelayDate, "seconds")) {
      cooldownStatus = (
        <>
          {renderProgressBar(
            withdrawalDate,
            endOfWithdrawalDelayDate,
            currentDate
          )}
          <div className={"pending-withdrawal__cooldown-time-container"}>
            <Icons.Time
              width="16"
              height="16"
              className="time-icon time-icon--grey-30"
            />
            <span>
              {days}d {hours}h {minutes}m until available
            </span>
          </div>
        </>
      )
    } else {
      cooldownStatus = (
        <Chip
          className={"pending_withdrawal__cooldown-status-chip"}
          color="violet"
          text={"cooldown completed"}
          size="small"
        />
      )
    }

    return (
      <div className={"pending-withdrawal__cooldown-status"}>
        {cooldownStatus}
      </div>
    )
  }

  const isWithdrawalDelayOver = (pendingWithdrawalTimestamp) => {
    const currentDate = moment()
    const endOfWithdrawalDelayDate = moment
      .unix(pendingWithdrawalTimestamp)
      .add(withdrawalDelay, "seconds")

    return currentDate.isAfter(endOfWithdrawalDelayDate)
  }

  const isWithdrawalTimeoutOver = (pendingWithdrawalTimestamp) => {
    const currentDate = moment()
    const endOfWithdrawalTimeoutDate = moment
      .unix(pendingWithdrawalTimestamp)
      .add(withdrawalDelay, "seconds")
      .add(withdrawalTimeout, "seconds")

    return currentDate.isAfter(endOfWithdrawalTimeoutDate, "second")
  }

  const renderTimeLeftToClaimText = (pendingWithdrawalTimestamp) => {
    const currentDate = moment()
    const endOfWithdrawalDelayDate = moment
      .unix(pendingWithdrawalTimestamp)
      .add(withdrawalDelay, "seconds")

    if (currentDate.isBefore(endOfWithdrawalDelayDate, "second")) {
      return <></>
    }

    const endOfWithdrawalTimeoutDate = moment
      .unix(pendingWithdrawalTimestamp)
      .add(withdrawalDelay, "seconds")
      .add(withdrawalTimeout, "seconds")

    const days = endOfWithdrawalTimeoutDate.diff(currentDate, "days")
    const hours = moment
      .duration(endOfWithdrawalTimeoutDate.diff(currentDate))
      .hours()
    const minutes = moment
      .duration(endOfWithdrawalTimeoutDate.diff(currentDate))
      .minutes()

    let timeToClaim = <></>
    if (!isWithdrawalTimeoutOver(pendingWithdrawalTimestamp)) {
      timeToClaim = (
        <span>
          Time left to claim: {days}d {hours}h {minutes}m
        </span>
      )
    } else {
      timeToClaim = <span>Tokens went back to pool</span>
    }

    return timeToClaim
  }

  return (
    <section className={"tile pending-withdrawal"}>
      <DataTable
        data={pendingWithdrawals}
        itemFieldId="pendingWithdrawalId"
        title="Pending withdrawal"
        withTooltip
        tooltipProps={resourceTooltipProps.pendingWithdrawal}
        noDataMessage="No pending withdrawals."
      >
        <Column
          header="amount"
          field="covAmount"
          renderContent={({ covAmount }) => {
            return (
              <TokenAmount
                amount={covAmount}
                wrapperClassName={"pending-withdrawal__token-amount"}
                amountClassName={"h2 text-brand-violet-100"}
                symbolClassName={"h3 text-brand-violet-100"}
              />
            )
          }}
        />
        <Column
          header="withdrawal initiated"
          field="timestamp"
          renderContent={({ timestamp }) => {
            const withdrawalDate = moment.unix(timestamp)
            return (
              <div className={"pending-withdrawal__date"}>
                <span>{withdrawalDate.format("DD-MM-YYYY")}</span>
                <span>{withdrawalDate.format("HH:mm:ss")}</span>
              </div>
            )
          }}
        />
        <Column
          header="cooldown status"
          field="timestamp"
          tdClassName={"cooldown-status-column"}
          renderContent={({ timestamp }) => {
            return renderCooldownStatus(timestamp)
          }}
        />
        <Column
          header=""
          field="timestamp"
          renderContent={({ timestamp }) => (
            <div className={"pending-withdrawal__button-container"}>
              <SubmitButton
                className="btn btn-lg btn-primary"
                onSubmitAction={(awaitingPromise) => {
                  if (isWithdrawalTimeoutOver) {
                    onReinitiateWithdrawal(awaitingPromise)
                  } else {
                    onClaimTokensSubmitButtonClick(awaitingPromise)
                  }
                }}
                disabled={!isWithdrawalDelayOver(timestamp)}
              >
                <span
                  className={
                    "pending-withdrawal__button-container__button-text"
                  }
                >
                  {isWithdrawalTimeoutOver(timestamp)
                    ? "reinitiate"
                    : "claim tokens"}
                </span>
              </SubmitButton>
              <span
                className={
                  "pending-withdrawal__button-container__time-left-text"
                }
              >
                {renderTimeLeftToClaimText(timestamp)}
              </span>
            </div>
          )}
        />
      </DataTable>
    </section>
  )
}

export default PendingWithdrawalsView
