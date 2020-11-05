import React, { useMemo } from "react"
import { formatDate } from "../utils/general.utils"
import { SubmitButton } from "./Button"
import { colors } from "../constants/colors"
import ProgressBar from "./ProgressBar"

import moment from "moment"
import TokenAmount from "./TokenAmount"
import { displayAmount } from "../utils/token.utils"
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
  const cliffPeriod = useMemo(
    () =>
      moment
        .unix(selectedGrant.cliff)
        .from(moment.unix(selectedGrant.start), true),
    [selectedGrant.cliff, selectedGrant.start]
  )

  const fullyUnlockedDate = useMemo(
    () =>
      moment.unix(selectedGrant.start).add(selectedGrant.duration, "seconds"),
    [selectedGrant.start, selectedGrant.duration]
  )

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
          <span className="text-small ml-a">{selectedGrant.id}</span>
        </div>

        <div className="flex row center mb-1">
          <Icons.Calendar width={12} height={12} />
          <span className="text-small ml-1">Issued</span>
          <span className="text-small ml-a">
            {formatDate(selectedGrant.start)}
          </span>
        </div>

        <div className="flex row center mb-1">
          <Icons.KeepToken width={12} height={12} />
          <span className="text-small ml-1">Issued Total</span>
          <span className="text-small ml-a">{totalAmount}</span>
        </div>

        <div className="flex row center">
          <Icons.Time width={12} height={12} className="time-icon--black" />
          <span className="text-small ml-1">Fully Unlocked</span>
          <span className="text-small ml-a">
            {formatDate(fullyUnlockedDate)}
          </span>
        </div>
        {/* TODO tooltip */}
        <div
          className="ml-2 text-caption text-grey-60"
          style={{ marginTop: "0.5rem" }}
        >
          {cliffPeriod}&nbsp;cliff
        </div>
      </section>
    </>
  )
}

export default TokenGrantOverview

const TokenGrantUnlockingdDetailsComponent = ({ selectedGrant }) => {
  return (
    <ProgressBar
      value={selectedGrant.unlocked}
      total={selectedGrant.amount}
      color={colors.grey70}
      bgColor={colors.grey10}
    >
      <div className="circular-progress-bar-percentage-label-wrapper">
        <ProgressBar.Circular radius={82} barWidth={16} />
        <ProgressBar.PercentageLabel text="Unlocked" />
      </div>
      <ProgressBar.Legend leftValueLabel="Locked" valueLabel="Unlocked" />
    </ProgressBar>
  )
}

export const TokenGrantUnlockingdDetails = TokenGrantUnlockingdDetailsComponent

export const TokenGrantStakedDetails = ({ selectedGrant, stakedAmount }) => {
  return (
    <ProgressBar
      value={stakedAmount}
      total={selectedGrant.amount}
      color={colors.mint80}
      bgColor={colors.mint20}
    >
      <div className="circular-progress-bar-percentage-label-wrapper">
        <ProgressBar.Circular radius={82} barWidth={16} />
        <ProgressBar.PercentageLabel text="Staked" />
      </div>
      <ProgressBar.Legend leftValueLabel="Unstaked" valueLabel="Staked" />
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
        value={selectedGrant.released}
        total={selectedGrant.amount}
        color={colors.secondary}
        bgColor={colors.bgSecondary}
      >
        <ProgressBar.Inline height={20} />
        <ProgressBar.Legend
          valueLabel="Withdrawn from Grant"
          leftValueLabel="Available to Withdraw"
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
