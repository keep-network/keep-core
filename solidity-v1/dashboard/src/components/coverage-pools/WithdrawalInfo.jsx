import OnlyIf from "../OnlyIf"
import * as Icons from "../Icons"
import TokenAmount from "../TokenAmount"
import { covKEEP, KEEP } from "../../utils/token.utils"
import Divider from "../Divider"
import Button from "../Button"
import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import { Keep } from "../../contracts"
import { CooldownPeriodBanner } from "../coverage-pools"
import Chip from "../Chip"
import { PENDING_WITHDRAWAL_STATUS } from "../../constants/constants"

const WithdrawalInfo = ({
  transactionFinished,
  containerTitle,
  submitBtnText,
  onBtnClick,
  onCancel,
  amount,
  pendingWithdrawalBalance,
  addedAmount,
  totalValueLocked,
  covTotalSupply,
  pendingWithdrawalState = PENDING_WITHDRAWAL_STATUS.NONE,
  children,
}) => {
  return (
    <WithdrawalInfo.Container
      transactionFinished={transactionFinished}
      containerTitle={containerTitle}
      submitBtnText={submitBtnText}
      onBtnClick={onBtnClick}
      onCancel={onCancel}
    >
      <OnlyIf condition={transactionFinished}>
        <h3 className={"withdraw-modal__success-text"}>
          <Icons.Time
            width={20}
            height={20}
            className="time-icon--yellow m-1"
          />
          Almost there!
        </h3>
        <h4 className={"text-gray-70 mb-1"}>
          After the <b>21 day cooldown</b> you can claim your tokens in the
          dashboard.
        </h4>
      </OnlyIf>
      <div className={"withdraw-modal__data"}>
        <OnlyIf
          condition={
            transactionFinished ||
            pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.EXPIRED ||
            pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.NONE
          }
        >
          <TokenAmount
            amount={amount}
            wrapperClassName={"withdraw-modal__token-amount"}
            token={covKEEP}
          />
          <TokenAmount
            wrapperClassName={"withdraw-modal__cov-token-amount"}
            amount={Keep.coveragePoolV1.estimatedBalanceFor(
              amount,
              covTotalSupply,
              totalValueLocked
            )}
            amountClassName={"h3 text-grey-60"}
            symbolClassName={"h3 text-grey-60"}
            token={KEEP}
          />
        </OnlyIf>
        <OnlyIf
          condition={
            !transactionFinished &&
            (pendingWithdrawalState === PENDING_WITHDRAWAL_STATUS.PENDING ||
              pendingWithdrawalState ===
                PENDING_WITHDRAWAL_STATUS.AVAILABLE_TO_WITHDRAW)
          }
        >
          <div
            className={
              "withdraw-modal__data-row withdraw-modal__data-row--baseline-align"
            }
          >
            <TokenAmount
              amount={addedAmount}
              wrapperClassName={"withdraw-modal__data-row__token-amount"}
              token={covKEEP}
            />
            <TokenAmount
              wrapperClassName={"withdraw-modal__data-row__cov-token-amount"}
              amount={Keep.coveragePoolV1.estimatedBalanceFor(
                addedAmount,
                covTotalSupply,
                totalValueLocked
              )}
              amountClassName={"h4 text-grey-50"}
              symbolClassName={"h4 text-grey-50"}
              token={KEEP}
            />
          </div>
          <div
            className={
              "withdraw-modal__data-row withdraw-modal__data-row--baseline-align mb-3"
            }
          >
            <span className={"h4 text-gray-70"}>+</span>
            <TokenAmount
              amount={pendingWithdrawalBalance}
              wrapperClassName={"withdraw-modal__data-row__token-amount"}
              amountClassName={"h4 text-grey-70"}
              symbolClassName={"h4 text-grey-70"}
              token={covKEEP}
            />
            <Chip
              text={`existing`}
              size="small"
              className={"withdraw-modal_existing-withdrawal-chip"}
              color="yellow"
            />
            <TokenAmount
              wrapperClassName={"withdraw-modal__data-row__cov-token-amount"}
              amount={Keep.coveragePoolV1.estimatedBalanceFor(
                pendingWithdrawalBalance,
                covTotalSupply,
                totalValueLocked
              )}
              amountClassName={"h4 text-grey-50"}
              symbolClassName={"h4 text-grey-50"}
              token={KEEP}
            />
          </div>
        </OnlyIf>

        {children}
      </div>
      <OnlyIf condition={!transactionFinished}>
        <CooldownPeriodBanner />
      </OnlyIf>
      <OnlyIf condition={transactionFinished}>
        <Divider
          className="divider divider--tile-fluid"
          style={{ marginTop: "auto" }}
        />
        <Button
          className="btn btn-lg btn-secondary success-modal-close-button"
          disabled={false}
          onClick={onCancel}
        >
          close
        </Button>
      </OnlyIf>
    </WithdrawalInfo.Container>
  )
}

WithdrawalInfo.Container = ({
  transactionFinished = false,
  containerTitle = "You are about to withdraw:",
  submitBtnText,
  onBtnClick,
  onCancel,
  children,
}) => {
  const container = transactionFinished ? (
    <>{children}</>
  ) : (
    <AcceptTermConfirmationModal
      title={containerTitle}
      termText="I confirm I have read the documentation and am aware of the risk."
      btnText={submitBtnText}
      onBtnClick={onBtnClick}
      onCancel={onCancel}
    >
      {children}
    </AcceptTermConfirmationModal>
  )

  return container
}

export default WithdrawalInfo
