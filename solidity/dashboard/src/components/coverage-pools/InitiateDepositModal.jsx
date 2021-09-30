import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import TokenAmount from "../TokenAmount"
import * as Icons from "../Icons"
import ModalWithTimeline, {
  MODAL_WITH_TIMELINE_STEPS,
} from "./ModalWithTImeline"
import OnlyIf from "../OnlyIf"
import { covKEEP, KEEP } from "../../utils/token.utils"
import Divider from "../Divider"
import Button from "../Button"
import { CooldownPeriodBanner } from "../coverage-pools"

const InitiateDepositModal = ({
  amount, // amount of KEEP that user wants to deposit (in KEEP)
  balanceAmount, // total balance of the user after the deposit is done (in covKEEP)
  estimatedBalanceAmountInKeep, // estimated total balance of user in KEEP
  submitBtnText,
  onBtnClick,
  onCancel,
  transactionFinished = false,
  transactionHash = "",
}) => {
  return (
    <ModalWithTimeline
      className={`withdraw-modal__main-container`}
      step={MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS}
      withDescription={true}
    >
      <InitiateDepositModal.Container
        transactionFinished={transactionFinished}
        submitBtnText={submitBtnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
      >
        <OnlyIf condition={transactionFinished}>
          <h3 className={"withdraw-modal__success-text"}>
            <Icons.Success
              width={20}
              height={20}
              className="success-icon--green m-1"
            />
            Success!
          </h3>
          <h4 className={"text-gray-70 mb-1"}>
            View your transaction&nbsp;
            <a
              href={`https://etherscan.io/tx/${transactionHash}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              here
            </a>
            .
          </h4>
        </OnlyIf>
        <div className={"withdraw-modal__data"}>
          <TokenAmount
            amount={amount}
            wrapperClassName={"withdraw-modal__token-amount"}
            token={KEEP}
            withIcon
          />
          <OnlyIf condition={transactionFinished}>
            <div className={"withdraw-modal__data-row"}>
              <h4 className={"text-grey-50"}>Pool Balance &nbsp;</h4>
              <h4 className={"withdraw-modal__data__value text-grey-70"}>
                {KEEP.displayAmountWithMetricSuffix(
                  estimatedBalanceAmountInKeep
                )}{" "}
                KEEP
              </h4>
            </div>
            <div className={"withdraw-modal__data-row"}>
              <h4 className={"text-grey-50"}>Coverage Token &nbsp;</h4>
              <h4 className={"withdraw-modal__data__value text-grey-70"}>
                {covKEEP.displayAmountWithMetricSuffix(balanceAmount)} covKEEP
              </h4>
            </div>
          </OnlyIf>
        </div>
        <OnlyIf condition={!transactionFinished}>
          <CooldownPeriodBanner />
        </OnlyIf>
        <OnlyIf condition={transactionFinished}>
          <Divider
            className="divider divider--tile-fluid"
            style={{ margin: "0.5rem 0", marginTop: "auto" }}
          />
          <Button
            className="btn btn-lg btn-secondary success-modal-close-button"
            disabled={false}
            onClick={onCancel}
          >
            close
          </Button>
        </OnlyIf>
      </InitiateDepositModal.Container>
    </ModalWithTimeline>
  )
}

InitiateDepositModal.Container = ({
  transactionFinished = false,
  submitBtnText,
  onBtnClick,
  onCancel,
  children,
}) => {
  const container = transactionFinished ? (
    <>{children}</>
  ) : (
    <AcceptTermConfirmationModal
      title="You are about to deposit:"
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

export default InitiateDepositModal
