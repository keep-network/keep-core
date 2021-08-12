import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import TokenAmount from "../TokenAmount"
import Banner from "../Banner"
import * as Icons from "../Icons"
import ModalWithTimeline, {MODAL_WITH_TIMELINE_STEPS} from "./ModalWithTImeline";
import OnlyIf from "../OnlyIf";
import { covKEEP, KEEP } from "../../utils/token.utils";
import { shortenAddress } from "../../utils/general.utils";
import Divider from "../Divider";
import Button from "../Button";
import WithdrawalInfo from "./WithdrawalInfo";

const infoBannerTitle =
  "The cooldown period is 21 days.."

const infoBannerDescription =
  "A withdrawn deposit will be available to claim after 21 days. During cooldown, your funds will accumulate rewards but are also subject to risk to cover for a hit."

const InitiateCovPoolsWithdrawModal = ({
  amount,
  containerTitle,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  return (
    <ModalWithTimeline
      className={`withdraw-modal__main-container ${className}`}
      step={
        transactionFinished ?
          MODAL_WITH_TIMELINE_STEPS.COOLDOWN :
          MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT
      }
      withDescription={true}>
      <WithdrawalInfo
        transactionFinished={transactionFinished}
        containerTitle={containerTitle}
        submitBtnText={submitBtnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
        amount={amount}
        infoBannerTitle={infoBannerTitle}
        infoBannerDescription={infoBannerDescription}
      >
        <div className={"withdraw-modal__data-row"}>
          <h4 className={"text-grey-50"}>Pool Balance &nbsp;</h4>
          <h4 className={"withdraw-modal__data__value text-grey-70"}>
            1,000 KEEP
          </h4>
        </div>
        <div className={"withdraw-modal__data-row"}>
          <h4 className={"text-grey-50"}>Earned Balance &nbsp;</h4>
          <h4 className={"withdraw-modal__data__value text-grey-70"}>
            1,000 KEEP
          </h4>
        </div>
        <div className={"withdraw-modal__data-row"}>
          <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
          <h4 className={"withdraw-modal__data__value text-grey-70"}>
            {shortenAddress("0x254673e7c7d76e051e80d30FCc3EA6A9C2a22222")}
          </h4>
        </div>
      </WithdrawalInfo>

    </ModalWithTimeline>
  )
}

export default InitiateCovPoolsWithdrawModal
