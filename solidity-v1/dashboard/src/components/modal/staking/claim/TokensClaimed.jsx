import React from "react"
import { withBaseModal } from "../../withBaseModal"
import ViewTransactionSuccessModal from "../ViewTransactionSuccessModal"

const TokensClaimedComponent = ({ txHash, onClose }) => {
  return (
    <ViewTransactionSuccessModal
      modalHeader={"Claim"}
      furtherDescription={"Go to the Threshold dapp to coplete upgrade"}
      txHash={txHash}
      onClose={onClose}
    />
  )
}

const TokensClaimedWithBaseModal = withBaseModal(TokensClaimedComponent)

export const TokensClaimed = (props) => (
  <TokensClaimedWithBaseModal size="sm" {...props} />
)
