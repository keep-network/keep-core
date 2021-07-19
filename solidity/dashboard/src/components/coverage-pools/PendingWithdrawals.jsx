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

const PendingWithdrawals = (onClaimTokensSubmitButtonClick) => {
  const {
    withdrawalDelay,
    withdrawalTimeout,
    pendingWithdrawals,
  } = useSelector((state) => state.coveragePool)

  const renderProgressBar = (
    withdrawalDate,
    endOfCooldownDate,
    currentDate
  ) => {
    const progressBarValueInSeconds = currentDate.diff(
      withdrawalDate,
      "seconds"
    )
    const progressBarTotalInSeconds = endOfCooldownDate.diff(
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
    const endOfCooldownDate = moment
      .unix(timestamp)
      .add(withdrawalDelay, "seconds")
    const days = endOfCooldownDate.diff(currentDate, "days")
    const hours = moment.duration(endOfCooldownDate.diff(currentDate)).hours()
    const minutes = moment
      .duration(endOfCooldownDate.diff(currentDate))
      .minutes()

    let cooldownStatus = <></>
    if (days >= 0 && hours >= 0 && minutes >= 0) {
      cooldownStatus = (
        <>
          {renderProgressBar(withdrawalDate, endOfCooldownDate, currentDate)}
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
      cooldownStatus = <Chip text={"cooldown completed"} size="small" />
    }

    return (
      <div className={"pending-withdrawal__cooldown-status"}>
        {cooldownStatus}
      </div>
    )
  }

  const isWithdrawalCooldownOver = (pendingWithdrawalTimestamp) => {
    const currentDate = moment()
    const endOfCooldownDate = moment
      .unix(pendingWithdrawalTimestamp)
      .add(withdrawalDelay, "seconds")

    return currentDate.isAfter(endOfCooldownDate)
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
            return <TokenAmount amount={covAmount} />
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
                onSubmitAction={onClaimTokensSubmitButtonClick}
                disabled={!isWithdrawalCooldownOver(timestamp)}
              >
                claim tokens
              </SubmitButton>
            </div>
          )}
        />
      </DataTable>
    </section>
  )
}

export default PendingWithdrawals
