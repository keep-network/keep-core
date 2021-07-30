import React from "react"
import * as Icons from "../Icons"
import TokenAmount from "../TokenAmount"
import { KEEP } from "../../utils/token.utils"
import Divider from "../Divider"
import Button from "../Button"
import OnlyIf from "../OnlyIf"

const ClaimTokensModal = ({
  amount,
  submitBtnText,
  onBtnClick,
  onCancel,
  transactionFinished = false,
  transactionHash = "",
}) => {
  return (
    <div className={"claim-tokens-modal__content-container"}>
      <div className={"claim-tokens-modal__info"}>
        <OnlyIf condition={!transactionFinished}>
          <h3 className={"mb-1"}>You are about to claim:</h3>
        </OnlyIf>
        <OnlyIf condition={transactionFinished}>
          <h3>Success!</h3>
          <h4 className={"text-gray-70 mb-1"}>View your transaction here.</h4>
        </OnlyIf>
        <div className={"claim-tokens-modal__data"}>
          <TokenAmount
            amount={"20000"}
            wrapperClassName={"claim-tokens-modal__token-amount"}
            token={KEEP}
            withIcon
          />
          <div className={"claim-tokens-modal__data-row"}>
            <h4 className={"text-grey-50"}>Initial Withdrawal &nbsp;</h4>
            <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
          <div className={"claim-tokens-modal__data-row"}>
            <h4 className={"text-grey-50"}>Rewards earned &nbsp;</h4>
            <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
          <div className={"claim-tokens-modal__data-row"}>
            <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
            <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
        </div>
        <Divider style={{ margin: "0.5rem 0" }} />
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
              className="btn btn-lg btn-secondary"
              disabled={false}
              onClick={onCancel}
            >
              Close
            </Button>
          </OnlyIf>
        </div>
      </div>
      <div className={"claim-tokens-modal__timeline"}>
        You are about to claim:
      </div>
    </div>
  )
}

export default ClaimTokensModal
