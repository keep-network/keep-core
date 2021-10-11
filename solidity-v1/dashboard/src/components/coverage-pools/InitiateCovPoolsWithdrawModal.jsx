import React from "react"
import ModalWithTimeline, {
  MODAL_WITH_TIMELINE_STEPS,
} from "./ModalWithTImeline"
import { covKEEP, KEEP } from "../../utils/token.utils"
import { shortenAddress } from "../../utils/general.utils"
import WithdrawalInfo from "./WithdrawalInfo"
import TokenAmount from "../TokenAmount"
import { useWeb3Address } from "../WithWeb3Context"
import { CoveragePoolV1ExchangeRate } from "./ExchangeRate"

const InitiateCovPoolsWithdrawModal = ({
  amount,
  covTokensAvailableToWithdraw,
  totalValueLocked,
  covTotalSupply,
  containerTitle,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  const yourAddress = useWeb3Address()

  return (
    <ModalWithTimeline
      className={`withdraw-modal__main-container ${className}`}
      step={
        transactionFinished
          ? MODAL_WITH_TIMELINE_STEPS.COOLDOWN
          : MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT
      }
      withDescription={true}
    >
      <WithdrawalInfo
        transactionFinished={transactionFinished}
        containerTitle={containerTitle}
        submitBtnText={submitBtnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
        amount={amount}
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
          <h4 className={"text-grey-50"}>Pool Balance &nbsp;</h4>
          <TokenAmount
            amount={covTokensAvailableToWithdraw}
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
    </ModalWithTimeline>
  )
}

export default InitiateCovPoolsWithdrawModal
