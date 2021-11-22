import React from "react"
import { withBaseModal } from "../withBaseModal"
import ViewTransactionSuccessModal from "../staking/ViewTransactionSuccessModal"

const GrantTokensWithdrawnComponent = ({ txHash, onClose }) => {
  return (
    <ViewTransactionSuccessModal
      modalHeader={"Withdraw"}
      furtherDescription={"Go to the Threshold dapp to coplete upgrade"}
      txHash={txHash}
      onClose={onClose}
    />
  )
}

const GrantTokensWithdrawnWithBaseModal = withBaseModal(
  GrantTokensWithdrawnComponent
)

export const GrantTokensWithdrawn = (props) => (
  <GrantTokensWithdrawnWithBaseModal size="sm" {...props} />
)
