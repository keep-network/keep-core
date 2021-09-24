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

const WithdrawalInfo = ({
  transactionFinished,
  containerTitle,
  submitBtnText,
  onBtnClick,
  onCancel,
  amount,
  totalValueLocked,
  covTotalSupply,
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
        {children}
      </div>
      <OnlyIf condition={!transactionFinished}>
        <CooldownPeriodBanner />
      </OnlyIf>
      <OnlyIf condition={transactionFinished}>
        <Divider className="divider divider--tile-fluid" />
        <Button
          className="btn btn-lg btn-secondary"
          disabled={false}
          onClick={onCancel}
        >
          Close
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
