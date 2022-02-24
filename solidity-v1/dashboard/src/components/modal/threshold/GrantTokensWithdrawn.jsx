import React from "react"
import { withBaseModal } from "../withBaseModal"
import ViewTransactionSuccessModal from "../staking/ViewTransactionSuccessModal"
import UpgradeToTButton from "../../UpgradeToTButton"

const GrantTokensWithdrawnComponent = ({ txHash, onClose }) => {
  return (
    <ViewTransactionSuccessModal
      modalHeader={"Withdraw"}
      furtherDescription={"Go to the Threshold dapp to coplete upgrade"}
      txHash={txHash}
      onClose={onClose}
      renderAdditionalButtons={
        <UpgradeToTButton className={"btn-primary btn-lg mr-1"} />
      }
    />
  )
}

const GrantTokensWithdrawnWithBaseModal = withBaseModal(
  GrantTokensWithdrawnComponent
)

export const GrantTokensWithdrawn = (props) => (
  <GrantTokensWithdrawnWithBaseModal size="sm" {...props} />
)
