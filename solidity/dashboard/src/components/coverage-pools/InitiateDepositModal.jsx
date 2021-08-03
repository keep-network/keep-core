import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import TokenAmount from "../TokenAmount"
import Banner from "../Banner"
import * as Icons from "../Icons"
import ModalWithTimeline, {MODAL_WITH_TIMELINE_STEPS} from "./ModalWithTImeline";
import OnlyIf from "../OnlyIf";
import {KEEP} from "../../utils/token.utils";
import Divider from "../Divider";
import Button from "../Button";

const infoBannerTitle =
  "The cooldown period is 21 days.."

const infoBannerDescription =
  "A withdrawn deposit will be available to claim after 21 days. During cooldown, your funds will accumulate rewards but are also subject to risk to cover for a hit."

const InitiateDepositModal = ({
  amount,
  submitBtnText,
  onBtnClick,
  onCancel,
  transactionFinished = false,
}) => {
  return (
    <ModalWithTimeline
      className={`withdraw-modal__main-container`}
      step={
        transactionFinished ?
          MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT :
          MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS
      }
      withDescription={true}>
      <InitiateDepositModal.Container
        transactionFinished={transactionFinished}
        submitBtnText={submitBtnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}>
        <OnlyIf condition={transactionFinished}>
          <h3 className={"withdraw-modal__success-text"}><Icons.Time width={20} height={20} className="time-icon--yellow m-1" />Almost there!</h3>
          <h4 className={"text-gray-70 mb-1"}>After the <b>21 day cooldown</b> you can claim your tokens in the dashboard.</h4>
        </OnlyIf>
        <div className={"withdraw-modal__data"}>
          <TokenAmount
            amount={amount}
            wrapperClassName={"withdraw-modal__token-amount"}
            token={KEEP}
            withIcon
          />
        </div>
        <OnlyIf condition={!transactionFinished}>
          <Banner
            icon={Icons.Tooltip}
            className="withdraw-modal__banner banner--info mt-2 mb-2"
          >
            <Banner.Icon icon={Icons.Tooltip} className={`withdraw-modal__banner-icon mr-1`} backgroundColor={"transparent"} color={"black"}/>
            <div className={"withdraw-modal__banner-icon-text"}>
              <Banner.Title>
                {infoBannerTitle}
              </Banner.Title>
              <Banner.Description>
                {infoBannerDescription}
              </Banner.Description>
            </div>
          </Banner>
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
      </InitiateDepositModal.Container>
    </ModalWithTimeline>
  )
}

InitiateDepositModal.Container = ({transactionFinished = false, submitBtnText, onBtnClick, onCancel, children}) => {
  const container = transactionFinished ? <>{children}</> : <AcceptTermConfirmationModal
    title="You are about to withdraw:"
    termText="I confirm I have read the documentation and am aware of the risk."
    btnText={submitBtnText}
    onBtnClick={onBtnClick}
    onCancel={onCancel}
  >{children}</AcceptTermConfirmationModal>

  return container
}

export default InitiateDepositModal
