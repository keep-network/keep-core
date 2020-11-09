import React, { useMemo } from "react"
import { formatDate } from "../utils/general.utils"
import { SubmitButton } from "./Button"
import { colors } from "../constants/colors"
import ProgressBar from "./ProgressBar"

import moment from "moment"
import TokenAmount from "./TokenAmount"
import {
  displayAmount,
  displayAmountWithMetricSuffix,
} from "../utils/token.utils"
import * as Icons from "./Icons"

const TokenGrantOverview = ({ selectedGrant, selectedGrantStakedAmount }) => {
  return (
    <>
      <TokenGrantDetails selectedGrant={selectedGrant} />
      <hr />
      <div className="flex">
        <TokenGrantUnlockingdDetails selectedGrant={selectedGrant} />
      </div>
      <div className="flex mt-1">
        <TokenGrantStakedDetails
          selectedGrant={selectedGrant}
          stakedAmount={selectedGrantStakedAmount}
        />
      </div>
    </>
  )
}

export const TokenGrantDetails = ({ selectedGrant, availableAmount }) => {
  const cliffPeriod = useMemo(() => {
    return selectedGrant.cliff && selectedGrant.start
      ? moment
          .unix(selectedGrant.cliff)
          .from(moment.unix(selectedGrant.start), true)
      : null
  }, [selectedGrant.cliff, selectedGrant.start])

  const fullyUnlockedDate = useMemo(() => {
    return selectedGrant.start && selectedGrant.duration
      ? moment.unix(selectedGrant.start).add(selectedGrant.duration, "seconds")
      : null
  }, [selectedGrant.start, selectedGrant.duration])

  const totalAmount = useMemo(() => displayAmount(selectedGrant.amount), [
    selectedGrant.amount,
  ])

  return (
    <>
      <TokenAmount amount={availableAmount} currencySymbol="KEEP" />
      <h4 className="text-grey-70 mt-3 mb-1">Grant Details</h4>
      <section className="grant-details">
        <div className="flex row center mb-1">
          <Icons.Grant width={12} height={12} />
          <span className="text-small ml-1">Grant ID</span>
          <span className="text-small ml-a">
            {selectedGrant.id || "No data"}
          </span>
        </div>

        <div className="flex row center mb-1">
          <Icons.Calendar width={12} height={12} />
          <span className="text-small ml-1">Issued</span>
          <span className="text-small ml-a">
            {selectedGrant.start
              ? formatDate(moment.unix(selectedGrant.start))
              : "No data"}
          </span>
        </div>

        <div className="flex row center mb-1">
          <Icons.KeepToken width={12} height={12} />
          <span className="text-small ml-1">Issued Total</span>
          <span className="text-small ml-a">{totalAmount || "No Data"}</span>
        </div>

        <div className="flex row center">
          <Icons.Time width={12} height={12} className="time-icon--black" />
          <span className="text-small ml-1">Fully Unlocked</span>
          <span className="text-small ml-a">
            {fullyUnlockedDate ? formatDate(fullyUnlockedDate) : "No data"}
          </span>
        </div>
        {/* TODO tooltip */}
        {cliffPeriod && (
          <div
            className="text-caption text-grey-60"
            style={{ marginTop: "0.5rem", marginLeft: "1.75rem" }}
          >
            {cliffPeriod}&nbsp;cliff
          </div>
        )}
      </section>
    </>
  )
}

export default TokenGrantOverview

const TokenGrantUnlockingdDetailsComponent = ({ selectedGrant }) => {
  return (
    <ProgressBar
      value={selectedGrant.unlocked || 0}
      total={selectedGrant.amount || 0}
      color={colors.grey70}
      bgColor={colors.grey10}
    >
      <div className="circular-progress-bar-percentage-label-wrapper">
        <ProgressBar.Circular radius={82} barWidth={16} />
        <ProgressBar.PercentageLabel text="Unlocked" />
      </div>
      <ProgressBar.Legend
        leftValueLabel="Locked"
        valueLabel="Unlocked"
        displayLegendValuFn={displayAmountWithMetricSuffix}
      />
    </ProgressBar>
  )
}

export const TokenGrantUnlockingdDetails = TokenGrantUnlockingdDetailsComponent

export const TokenGrantStakedDetails = ({ selectedGrant, stakedAmount }) => {
  return (
    <ProgressBar
      value={stakedAmount || 0}
      total={selectedGrant.amount || 0}
      color={colors.mint80}
      bgColor={colors.mint20}
    >
      <div className="circular-progress-bar-percentage-label-wrapper">
        <ProgressBar.Circular radius={82} barWidth={16} />
        <ProgressBar.PercentageLabel text="Staked" />
      </div>
      <ProgressBar.Legend
        leftValueLabel="Unstaked"
        valueLabel="Staked"
        displayLegendValuFn={displayAmountWithMetricSuffix}
      />
    </ProgressBar>
  )
}

export const TokenGrantWithdrawnTokensDetails = ({
  selectedGrant,
  onWithdrawnBtn,
}) => {
  return (
    <>
      <ProgressBar
        value={selectedGrant.released || 0}
        total={selectedGrant.amount || 0}
        color={colors.secondary}
        bgColor={colors.bgSecondary}
      >
        <ProgressBar.Inline height={20} />
        <ProgressBar.Legend
          valueLabel="Withdrawn from Grant"
          leftValueLabel="Available to Withdraw"
          displayLegendValuFn={displayAmountWithMetricSuffix}
        />
      </ProgressBar>
      <SubmitButton
        className="btn btn-secondary btn-sm mt-2"
        onClick={onWithdrawnBtn}
      >
        withdraw tokens
      </SubmitButton>
    </>
  )
}

// const ConfirmWithdrawModal = ({ escrowAddress }) => {
//   return (
//     <>
//       <span>You have deposited tokens in the</span>&nbsp;
//       <ViewAddressInBlockExplorer
//         text="TokenStakingEscrow contract"
//         address={escrowAddress}
//       />
//       <p>
//         To withdraw all tokens it may be necessary to confirm more than one
//         transaction.
//       </p>
//     </>
//   )
// }
