import React from "react"
import { ModalBody, ModalFooter } from "../Modal"
import List from "../../List"
import Button from "../../Button"
import { covKEEP } from "../../../utils/token.utils"
import { withWithdrawalOverview } from "./withWithdrawalOverview"
import { PENDING_WITHDRAWAL_STATUS } from "../../../constants/constants"
import OnlyIf from "../../OnlyIf"

const ConfirmIncreaseWithdrawalComponent = ({
  covAmountToAdd,
  withdrawalStatus,
  onConfirm,
  onClose,
}) => {
  const getItemTitle = () => {
    if (withdrawalStatus === PENDING_WITHDRAWAL_STATUS.EXPIRED)
      return "expired withdrawal"
    else if (withdrawalStatus === PENDING_WITHDRAWAL_STATUS.COMPLETED)
      return "existing claimable tokens"

    return "existing withdrawal"
  }
  return (
    <>
      <ModalBody>
        <h3 className="mb-1">
          {withdrawalStatus === PENDING_WITHDRAWAL_STATUS.COMPLETED
            ? "You have claimable tokens."
            : "Take note!"}
        </h3>
        <List>
          <List.Title className="text-grey-70 mb-1">
            <OnlyIf
              condition={withdrawalStatus === PENDING_WITHDRAWAL_STATUS.EXPIRED}
            >
              Your expired withdrawal needs to be re-initiated.&nbsp;
            </OnlyIf>
            This withdrawal will:
          </List.Title>
          <List.Content className="bullets text-grey-70">
            <List.Item>
              Add{" "}
              <strong>{covKEEP.displayAmountWithSymbol(covAmountToAdd)}</strong>{" "}
              to your {getItemTitle()}.
            </List.Item>
            <List.Item>
              Reset the <strong>21 day cooldown period</strong>
              <OnlyIf
                condition={
                  withdrawalStatus === PENDING_WITHDRAWAL_STATUS.COMPLETED
                }
              >
                &nbsp;of your currently claimable tokens
              </OnlyIf>
              .
            </List.Item>
          </List.Content>
        </List>
        <p className="text-grey-70">Do you want to continue?</p>
      </ModalBody>

      <ModalFooter>
        <Button
          className="btn btn-primary btn-lg mr-2"
          type="submit"
          onClick={onConfirm}
        >
          continue
        </Button>
        <Button className="btn btn-unstyled text-link" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
}

export const ConfirmIncreaseWithdrawal = withWithdrawalOverview({
  title: "Increase Withdrawal",
})(ConfirmIncreaseWithdrawalComponent)
