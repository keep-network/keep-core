import React, { useContext } from "react"
import moment from "moment"
import TokenAmount from "../../../TokenAmount"
import ProgressBar from "../../../ProgressBar"
import * as Icons from "../../../Icons"
import { covKEEP } from "../../../../utils/token.utils"
import { add } from "../../../../utils/arithmetics.utils"
import { PENDING_WITHDRAWAL_STATUS } from "../../../../constants/constants"
import { colors } from "../../../../constants/colors"

const WithdrawalOverviewContext = React.createContext({
  withdrawalDelay: 0,
  withdrawalInitiatedTimestamp: 0,
})

const useWithdrawalOverviewContext = () => {
  const context = useContext(WithdrawalOverviewContext)

  if (!context) {
    throw new Error(
      "useWithdrawalOverviewContext used outside of WithdrawalOverview component"
    )
  }

  return context
}

const tileBaseStyle = {
  backgroundColor: colors.white,
  padding: "1rem",
  borderRadius: "0.5rem",
  border: "1px solid",
  boxShadow: "0px 4px 4px rgba(196, 196, 196, 0.3)",
}

const styles = {
  expired: {
    label: {
      color: colors.error,
    },
    tile: {
      ...tileBaseStyle,
      borderColor: colors.red60,
    },
  },
  completed: {
    label: {
      color: colors.success,
    },
    tile: {
      ...tileBaseStyle,
      borderColor: colors.primary,
    },
  },
  pending: {
    label: {
      color: colors.grey50,
    },
    tile: {
      ...tileBaseStyle,
      borderColor: colors.yellowSecondary,
    },
  },
  new: {
    label: {
      color: colors.grey50,
    },
    tile: {
      ...tileBaseStyle,
      borderColor: colors.yellowSecondary,
    },
  },
  progressBar: {
    margin: "0",
    marginTop: "0.5rem",
  },
  title: {
    marginBottom: "0.5rem",
  },
}

const statusToProgressBarColor = {
  [PENDING_WITHDRAWAL_STATUS.COMPLETED]: colors.mint80,
  [PENDING_WITHDRAWAL_STATUS.PENDING]: colors.yellowSecondary,
  [PENDING_WITHDRAWAL_STATUS.EXPIRED]: colors.error,
  [PENDING_WITHDRAWAL_STATUS.NEW]: colors.grey20,
}

const WithdrawOverviewTile = ({ title, amount, label, status }) => {
  const { withdrawalDelay, withdrawalInitiatedTimestamp } =
    useWithdrawalOverviewContext()

  const getProgressBarProps = () => {
    if (status === PENDING_WITHDRAWAL_STATUS.PENDING) {
      const progressBarValueInSeconds = moment().diff(
        moment.unix(withdrawalInitiatedTimestamp),
        "seconds"
      )
      const progressBarTotalInSeconds = moment
        .unix(withdrawalInitiatedTimestamp)
        .add(withdrawalDelay, "seconds")
        .diff(moment.unix(withdrawalInitiatedTimestamp), "seconds")
      return {
        total: progressBarTotalInSeconds,
        value: progressBarValueInSeconds,
        color: statusToProgressBarColor[status],
      }
    }

    return {
      total: 100,
      value: 100,
      color: statusToProgressBarColor[status],
    }
  }

  return (
    <div style={styles[status].tile}>
      <h5 className="text-grey-50" style={styles.title}>
        {title}
      </h5>
      <h4 className={"flex row center"}>
        <TokenAmount
          amount={amount}
          token={covKEEP}
          amountClassName={"text-small text-grey-70"}
          symbolClassName={"text-small text-gray-70"}
        />
        <span className="text-small ml-a" style={styles[status].label}>
          {label}
        </span>
      </h4>
      <ProgressBar {...getProgressBarProps()} bgColor={colors.grey20}>
        <ProgressBar.Inline height={10} style={styles.progressBar} />
      </ProgressBar>
    </div>
  )
}

export const WithdrawalOverview = ({
  existingWithdrawalCovAmount,
  covAmountToAdd,
  withdrawalDelay,
  withdrawalInitiatedTimestamp,
  withdrawalStatus,
}) => {
  // TODO get title based on the current status
  const getTitle = () => {
    if (withdrawalStatus === PENDING_WITHDRAWAL_STATUS.EXPIRED) {
      return "expired withdrawal"
    }
    return "existing withdrawal"
  }

  const getLabel = () => {
    if (withdrawalStatus === PENDING_WITHDRAWAL_STATUS.PENDING) {
      return moment
        .unix(withdrawalInitiatedTimestamp)
        .add(withdrawalDelay, "seconds")
        .format("MM/DD")
    } else if (withdrawalDelay === PENDING_WITHDRAWAL_STATUS.COMPLETED) {
      return "Completed"
    }

    return "Expired"
  }

  const endOfTheNewWithdrawalDelayDate = moment().add(
    withdrawalDelay,
    "seconds"
  )

  return (
    <WithdrawalOverviewContext.Provider
      value={{
        withdrawalDelay,
        withdrawalInitiatedTimestamp,
      }}
    >
      <div className="modal__withdrawal-overview">
        <h4 className="mb-1">Overview</h4>
        <WithdrawOverviewTile
          title={getTitle()}
          amount={existingWithdrawalCovAmount}
          label={getLabel()}
          status={withdrawalStatus}
        />
        <h4 className={"flex row full-center text-grey-70 mt-1 mb-1"}>
          <Icons.ArrowDown />
          &nbsp;
          <Icons.Add />
          &nbsp;
          <TokenAmount
            amount={covAmountToAdd}
            token={covKEEP}
            amountClassName={"h4 text-grey-70"}
            symbolClassName={"h4 text-gray-70"}
          />
        </h4>
        <WithdrawOverviewTile
          title="new withdrawal"
          amount={add(existingWithdrawalCovAmount, covAmountToAdd)}
          label={endOfTheNewWithdrawalDelayDate.format("MM/DD")}
          status={PENDING_WITHDRAWAL_STATUS.NEW}
        />
      </div>
    </WithdrawalOverviewContext.Provider>
  )
}
