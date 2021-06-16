import React from "react"
import { AcceptTermConfirmationModal } from "../ConfirmationModal"
import TokenAmount from "../TokenAmount"
import Banner from "../Banner"
import * as Icons from "../Icons"

const warningBannerTitle =
  "Standard cooldown period is 14 days. During this cooldown period, your funds will continue to earn rewards and funds will be at risk of liquidation."

const InitiateDepositModal = ({
  amount,
  submitBtnText,
  onBtnClick,
  onCancel,
}) => {
  return (
    <AcceptTermConfirmationModal
      title="You are about to deposit:"
      termText="I confirm I have read the documentation and am aware of the risk."
      btnText={submitBtnText}
      onBtnClick={onBtnClick}
      onCancel={onCancel}
    >
      <TokenAmount
        amount={amount}
        withIcon
        icon={Icons.Plus}
        iconProps={{
          width: 24,
          height: 24,
          className: "plus-icon plus-icon--mint-100",
        }}
      />
      <Banner
        inline
        icon={Icons.Time}
        title={warningBannerTitle}
        className="banner--warning mt-2 mb-2"
      ></Banner>
    </AcceptTermConfirmationModal>
  )
}

export default InitiateDepositModal
