import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import TokenAmount from "../TokenAmount"
import Banner from "../Banner"
import * as Icons from "../Icons"
import ModalWithTimeline, {MODAL_WITH_TIMELINE_STEPS} from "./ModalWithTImeline";
import OnlyIf from "../OnlyIf";
import { covKEEP, KEEP } from "../../utils/token.utils";
import { shortenAddress } from "../../utils/general.utils";

const infoBannerTitle =
  "The cooldown period is 21 days.."

const infoBannerDescription =
  "A withdrawn deposit will be available to claim after 21 days. During cooldown, your funds will accumulate rewards but are also subject to risk to cover for a hit."

const InitiateCovPoolsWithdrawModal = ({
  amount,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  return (
    <ModalWithTimeline className={`withdraw-modal__main-container`} step={MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT} withDescription={true}>
      <AcceptTermConfirmationModal
        title="You are about to withdraw:"
        termText="I confirm I have read the documentation and am aware of the risk."
        btnText={submitBtnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
      >
        <OnlyIf condition={transactionFinished}>
          <h3>Success!</h3>
          <h4 className={"text-gray-70 mb-1"}>View your transaction here.</h4>
        </OnlyIf>
        <div className={"withdraw-modal__data"}>
          <TokenAmount
            amount={"200000000000000000000"}
            wrapperClassName={"withdraw-modal__token-amount"}
            token={KEEP}
            withIcon
          />
          <TokenAmount
            wrapperClassName={"withdraw-modal__cov-token-amount"}
            amount={"200000000000000000000"}
            amountClassName={"h3 text-grey-60"}
            symbolClassName={"h3 text-grey-60"}
            token={covKEEP}
          />
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
        </div>
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
      </AcceptTermConfirmationModal>
    </ModalWithTimeline>
  )
}

export default InitiateCovPoolsWithdrawModal
