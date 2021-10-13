import React, { useState } from "react"
import OnlyIf from "../OnlyIf"
import * as Icons from "./../Icons"
import { colors } from "../../constants/colors"
import ProgressBar from "../ProgressBar"
import moment from "moment"
import { useDispatch } from "react-redux"
import List from "../List"
import Divider from "../Divider"
import Button from "../Button"
import { covKEEP, KEEP } from "../../utils/token.utils"
import { shortenAddress } from "../../utils/general.utils"
import WithdrawalInfo from "./WithdrawalInfo"
import { add } from "../../utils/arithmetics.utils"
import { useWeb3Address } from "../WithWeb3Context"
import { addAdditionalDataToModal } from "../../actions/modal"
import TokenAmount from "../TokenAmount"
import { CoveragePoolV1ExchangeRate } from "./ExchangeRate"
import { getPendingWithdrawalStatus } from "../../utils/coverage-pools.utils"
import { PENDING_WITHDRAWAL_STATUS } from "../../constants/constants"

const getItems = (covKeepAmount, pendingWithdrawalStatus) => {
  if (pendingWithdrawalStatus === PENDING_WITHDRAWAL_STATUS.PENDING) {
    return [
      {
        label: (
          <>
            Add&nbsp;
            <strong>{covKEEP.displayAmountWithSymbol(covKeepAmount)}</strong>
            &nbsp;to your existing withdrawal.
          </>
        ),
      },
      {
        label: (
          <>
            Reset the&nbsp;<strong>21 day</strong>&nbsp;cooldown period.
          </>
        ),
      },
    ]
  } else if (
    pendingWithdrawalStatus === PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
  ) {
    return [
      {
        label: (
          <>
            Add&nbsp;
            <strong>{covKEEP.displayAmountWithSymbol(covKeepAmount)}</strong>
            &nbsp;to your existing claimable tokens.
          </>
        ),
      },
      {
        label: (
          <>
            Reset the&nbsp;<strong>21 day</strong>&nbsp;cooldown period of your
            currently claimable tokens.
          </>
        ),
      },
    ]
  } else {
    return [
      {
        label: (
          <>
            Add&nbsp;
            <strong>{covKEEP.displayAmountWithSymbol(covKeepAmount)}</strong>
            &nbsp;to your existing expired withdrawal
          </>
        ),
      },
      {
        label: (
          <>
            Reset the&nbsp;<strong>21 day</strong>&nbsp;cooldown period.
          </>
        ),
      },
    ]
  }
}

const IncreaseWithdrawalModal = ({
  pendingWithdrawalBalance,
  amount, // amount addedd to withdrawal
  totalValueLocked,
  covTotalSupply,
  withdrawalDelay,
  withdrawalTimeout,
  withdrawalInitiatedTimestamp,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  const dispatch = useDispatch()
  const yourAddress = useWeb3Address()
  const [step, setStep] = useState(transactionFinished ? 2 : 1)

  const onSubmit = (values) => {
    if (step === 1) {
      setStep((prevStep) => prevStep + 1)
    } else if (step === 2) {
      dispatch(
        addAdditionalDataToModal({
          componentProps: {
            pendingWithdrawalBalance: pendingWithdrawalBalance,
            amount: amount,
          },
        })
      )
      onBtnClick()
    }
  }

  const pendingWithdrawalState = getPendingWithdrawalStatus(
    withdrawalDelay,
    withdrawalTimeout,
    withdrawalInitiatedTimestamp
  )

  return (
    <ModalWithOverview
      className={`${className} withdraw-modal__main-container`}
      pendingWithdrawalBalance={pendingWithdrawalBalance}
      addedAmount={amount}
      totalValueLocked={totalValueLocked}
      covTotalSupply={covTotalSupply}
      withdrawalDelay={withdrawalDelay}
      withdrawalTimeout={withdrawalTimeout}
      withdrawalInitiatedTimestamp={withdrawalInitiatedTimestamp}
      pendingWithdrawalState={pendingWithdrawalState}
    >
      <OnlyIf condition={step === 1}>
        <IncreaseWithdrawalModalStep1
          addedAmount={amount}
          totalValueLocked={totalValueLocked}
          covTotalSupply={covTotalSupply}
          onSubmit={onSubmit}
          onCancel={onCancel}
          pendingWithdrawalState={pendingWithdrawalState}
        />
      </OnlyIf>
      <OnlyIf condition={step === 2}>
        <WithdrawalInfo
          transactionFinished={transactionFinished}
          containerTitle={"Your new withdrawal amount"}
          submitBtnText={"withdraw"}
          onBtnClick={onSubmit}
          onCancel={onCancel}
          amount={add(pendingWithdrawalBalance, amount)}
          totalValueLocked={totalValueLocked}
          covTotalSupply={covTotalSupply}
        >
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Exchange Rate&nbsp;</h4>
            <CoveragePoolV1ExchangeRate
              htmlTag="h4"
              className="withdraw-modal__data__value text-grey-70"
              covToken={covKEEP}
              collateralToken={KEEP}
              covTotalSupply={covTotalSupply}
              totalValueLocked={totalValueLocked}
            />
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Expired withdrawal&nbsp;</h4>
            <TokenAmount
              amount={pendingWithdrawalBalance}
              wrapperClassName={"withdraw-modal__data__value"}
              amountClassName={"h4 text-grey-70"}
              symbolClassName={"h4 text-grey-70"}
              token={covKEEP}
            />
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Increase amount &nbsp;</h4>
            <TokenAmount
              amount={amount}
              wrapperClassName={"withdraw-modal__data__value"}
              amountClassName={"h4 text-grey-70"}
              symbolClassName={"h4 text-grey-70"}
              token={covKEEP}
            />
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              {shortenAddress(yourAddress)}
            </h4>
          </div>
        </WithdrawalInfo>
      </OnlyIf>
    </ModalWithOverview>
  )
}

const IncreaseWithdrawalModalStep1 = ({
  addedAmount,
  totalValueLocked,
  covTotalSupply,
  onSubmit,
  onCancel,
  pendingWithdrawalState,
}) => {
  const items = getItems(addedAmount, pendingWithdrawalState)

  const getContentTitle = (pendingWithdrawalState) => {
    if (
      pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
    ) {
      return "You have claimable tokens."
    }

    return "Take note!"
  }

  const getMainNote = (pendingWithdrawalState) => {
    if (pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.EXPIRED) {
      return "Your expired withdrawal needs to be re-initiated. This withdrawal will:"
    }

    return "This withdrawal will:"
  }

  return (
    <div>
      <h3 className={"mb-1"}>{getContentTitle(pendingWithdrawalState)}</h3>
      <h4 className={"text-grey-70"}>{getMainNote(pendingWithdrawalState)}:</h4>
      <List
        items={items}
        className="increase-withdrawal-modal-step1__list mt-1 mb-1"
      >
        <List.Content className="bullets text-grey-70" />
      </List>
      <h4 className={"text-grey-70 mb-3"}>Do you want to continue?</h4>
      <Divider className="divider divider--tile-fluid" />
      <div className={"flex row center"}>
        <Button className="btn btn-lg btn-primary" onClick={onSubmit}>
          continue
        </Button>
        <span onClick={onCancel} className="ml-2 text-link text-grey-70">
          Cancel
        </span>
      </div>
    </div>
  )
}

const ModalWithOverview = ({
  children,
  className = "",
  pendingWithdrawalBalance,
  totalValueLocked,
  covTotalSupply,
  withdrawalDelay,
  withdrawalTimeout,
  withdrawalInitiatedTimestamp,
  addedAmount,
  pendingWithdrawalState,
}) => {
  return (
    <div className={`modal-with-overview__content-container ${className}`}>
      <div className={"modal-with-overview-modal__info"}>{children}</div>
      <div className={"modal-with-overview__overview-container"}>
        <h4 className={"mb-1"}>Overview</h4>
        <IncreaseWithdrawalModal.Tile
          amount={pendingWithdrawalBalance}
          totalValueLocked={totalValueLocked}
          covTotalSupply={covTotalSupply}
          withdrawalDelay={withdrawalDelay}
          withdrawalTimeout={withdrawalTimeout}
          withdrawalInitiatedTimestamp={withdrawalInitiatedTimestamp}
          pendingWithdrawalState={pendingWithdrawalState}
          expired
        />
        <h4 className={"modal-with-overview__added-amount text-grey-70"}>
          <Icons.ArrowDown />
          &nbsp;
          <Icons.Add />
          &nbsp;
          <TokenAmount
            amount={addedAmount}
            token={covKEEP}
            amountClassName={"h4 text-grey-70"}
            symbolClassName={"h4 text-gray-70"}
          />
        </h4>
        <IncreaseWithdrawalModal.Tile
          amount={add(pendingWithdrawalBalance, addedAmount)}
          totalValueLocked={totalValueLocked}
          covTotalSupply={covTotalSupply}
          withdrawalDelay={withdrawalDelay}
          withdrawalTimeout={withdrawalTimeout}
          pendingWithdrawalState={PENDING_WITHDRAWAL_STATUS.NONE}
        />
      </div>
    </div>
  )
}

const IncreaseWithdrawalModalTile = ({
  title,
  amount,
  totalValueLocked,
  covTotalSupply,
  withdrawalDelay,
  withdrawalTimeout,
  /** if null then it is a new withdrawal */
  withdrawalInitiatedTimestamp = null,
  pendingWithdrawalState,
}) => {
  const endOfWithdrawalDelayDate = withdrawalInitiatedTimestamp
    ? moment.unix(withdrawalInitiatedTimestamp)
    : moment()
  endOfWithdrawalDelayDate.add(withdrawalDelay, "seconds")
  const endOfWithdrawalTimeoutDate = withdrawalInitiatedTimestamp
    ? moment.unix(withdrawalInitiatedTimestamp)
    : moment()
  endOfWithdrawalTimeoutDate
    .add(withdrawalDelay, "seconds")
    .add(withdrawalTimeout, "seconds")

  const getTitle = () => {
    if (!title) {
      if (!withdrawalInitiatedTimestamp) {
        return "new withdrawal"
      }

      if (
        pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.PENDING ||
        pendingWithdrawalState ===
          PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
      ) {
        return "existing withdrawal"
      } else if (pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.EXPIRED) {
        return "expired withdrawal"
      }
    }

    return title
  }

  const getMainClassModifier = () => {
    if (!withdrawalInitiatedTimestamp) {
      return "modal-with-overview__tile--new"
    } else {
      if (pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.PENDING) {
        return ""
      } else if (
        pendingWithdrawalState ===
        PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
      ) {
        return "modal-with-overview__tile--completed"
      } else {
        return "modal-with-overview__tile--expired"
      }
    }
  }

  const renderProgressBar = (
    withdrawalInitiatedTimestamp,
    endOfWithdrawalDelayDate,
    currentDate
  ) => {
    const progressBarValueInSeconds = withdrawalInitiatedTimestamp
      ? currentDate.diff(moment.unix(withdrawalInitiatedTimestamp), "seconds")
      : 0
    const progressBarTotalInSeconds = withdrawalInitiatedTimestamp
      ? endOfWithdrawalDelayDate.diff(
          moment.unix(withdrawalInitiatedTimestamp),
          "seconds"
        )
      : 100

    let mainColor = colors.yellowSecondary
    if (
      pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
    ) {
      mainColor = colors.mint80
    } else if (pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.EXPIRED) {
      mainColor = colors.error
    }

    return (
      <ProgressBar
        value={progressBarValueInSeconds}
        total={progressBarTotalInSeconds}
        color={mainColor}
        bgColor={colors.yellow30}
      >
        <ProgressBar.Inline
          height={10}
          className={"modal-with-overview__progress-bar"}
        />
      </ProgressBar>
    )
  }

  return (
    <div className={`modal-with-overview__tile ${getMainClassModifier()}`}>
      <h5 className={"modal-with-overview__tile-title text-grey-50"}>
        {getTitle()}
      </h5>
      <div className={"modal-with-overview__withdrawal-info"}>
        <h4 className={"modal-with-overview__amount text-grey-70"}>
          <TokenAmount
            amount={amount}
            token={covKEEP}
            amountClassName={"h4 text-grey-70"}
            symbolClassName={"h4 text-gray-70"}
          />
        </h4>
        <OnlyIf
          condition={
            pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.PENDING ||
            !withdrawalInitiatedTimestamp
          }
        >
          <div className={"modal-with-overview__delay text-grey-50"}>
            {Math.floor(withdrawalDelay / (60 * 60 * 24))} days:{" "}
            {endOfWithdrawalDelayDate.format("MM/DD")}
          </div>
        </OnlyIf>
        <OnlyIf
          condition={
            pendingWithdrawalState ===
            PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW
          }
        >
          <div className={"modal-with-overview__delay text-success"}>
            Completed
          </div>
        </OnlyIf>
        <OnlyIf
          condition={
            pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.EXPIRED
          }
        >
          <div className={"modal-with-overview__delay text-error"}>Expired</div>
        </OnlyIf>
      </div>
      {renderProgressBar(
        withdrawalInitiatedTimestamp,
        endOfWithdrawalDelayDate,
        moment()
      )}
    </div>
  )
}

IncreaseWithdrawalModal.Tile = IncreaseWithdrawalModalTile

export default IncreaseWithdrawalModal
