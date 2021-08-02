import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import TokenAmount from "../TokenAmount"
import Banner from "../Banner"
import * as Icons from "../Icons"
import ModalWithTimeline from "./ModalWithTImeline";
import OnlyIf from "../OnlyIf";
import {KEEP} from "../../utils/token.utils";
import {shortenAddress} from "../../utils/general.utils";

const warningBannerTitle =
  "Standard cooldown period is 21 days. During this cooldown period, your funds will continue to earn rewards and funds will be at risk of liquidation."

const WithdrawModal = ({
  amount,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  return (
    <ModalWithTimeline className={`withdraw-modal__main-container`}>
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
            amount={"20000"}
            wrapperClassName={"withdraw-modal__token-amount"}
            token={KEEP}
            withIcon
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
          inline
          icon={Icons.Tooltip}
          title={warningBannerTitle}
          className="banner--warning mt-2 mb-2"
        />
      </AcceptTermConfirmationModal>
    </ModalWithTimeline>
  )
}

export default WithdrawModal
