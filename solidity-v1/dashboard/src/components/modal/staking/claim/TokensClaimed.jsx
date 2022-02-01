import React from "react"
import { withBaseModal } from "../../withBaseModal"
import ViewTransactionSuccessModal from "../ViewTransactionSuccessModal"
import UpgradeToTButton from "../../../UpgradeToTButton"

const TokensClaimedComponent = ({ txHash, onClose }) => {
  return (
    <ViewTransactionSuccessModal
      modalHeader={"Claim"}
      furtherDescription={"Go to the Threshold dapp to coplete upgrade"}
      txHash={txHash}
      onClose={onClose}
      renderAdditionalButtons={
        <UpgradeToTButton className={"btn-primary btn-lg mr-1"} />
      }
    />
  )
}

const TokensClaimedWithBaseModal = withBaseModal(TokensClaimedComponent)

export const TokensClaimed = (props) => (
  <TokensClaimedWithBaseModal size="sm" {...props} />
)
