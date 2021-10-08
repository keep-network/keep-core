import React from "react"
import TokenAmount from "../TokenAmount"
import { covKEEP, KEEP } from "../../utils/token.utils"
import Divider from "../Divider"
import Button from "../Button"
import OnlyIf from "../OnlyIf"
import ModalWithTimeline, {
  MODAL_WITH_TIMELINE_STEPS,
} from "./ModalWithTImeline"
import { shortenAddress } from "../../utils/general.utils"
import { Keep } from "../../contracts"

const ClaimTokensModal = ({
  covAmount,
  collateralTokenAmount,
  submitBtnText,
  onBtnClick,
  onCancel,
  totalValueLocked,
  covTotalSupply,
  address,
  transactionFinished = false,
  transactionHash = "",
}) => {
  return (
    <ModalWithTimeline
      className={"claim-tokens-modal__main-container"}
      step={MODAL_WITH_TIMELINE_STEPS.CLAIM_TOKENS}
      withDescription={!transactionFinished}
    >
      <OnlyIf condition={!transactionFinished}>
        <h3 className={"mb-1"}>You are about to claim:</h3>
      </OnlyIf>
      <OnlyIf condition={transactionFinished}>
        <h3>Success!</h3>
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
      <div className={"claim-tokens-modal__data"}>
        <TokenAmount
          amount={covAmount}
          wrapperClassName={"claim-tokens-modal__token-amount"}
          token={covKEEP}
        />
        <div className={"claim-tokens-modal__data-row"}>
          <h4 className={"text-grey-50"}>Initial Withdrawal &nbsp;</h4>
          <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
            {KEEP.displayAmountWithSymbol(
              transactionFinished
                ? collateralTokenAmount
                : Keep.coveragePoolV1.estimatedBalanceFor(
                    covAmount,
                    covTotalSupply,
                    totalValueLocked
                  )
            )}
          </h4>
        </div>
        <div className={"claim-tokens-modal__data-row"}>
          <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
          <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
            {shortenAddress(address)}
          </h4>
        </div>
      </div>
      <Divider style={{ margin: "0.5rem 0", marginTop: "auto" }} />
      <div className="flex row center mt-2">
        <OnlyIf condition={!transactionFinished}>
          <Button
            className="btn btn-lg btn-primary"
            type="submit"
            disabled={false}
            onClick={onBtnClick}
          >
            {submitBtnText}
          </Button>
          <span onClick={onCancel} className="ml-1 text-link">
            Cancel
          </span>
        </OnlyIf>
        <OnlyIf condition={transactionFinished}>
          <Button
            className="btn btn-lg btn-secondary success-modal-close-button"
            disabled={false}
            onClick={onCancel}
          >
            close
          </Button>
        </OnlyIf>
      </div>
    </ModalWithTimeline>
  )
}

export default ClaimTokensModal
