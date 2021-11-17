import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../../Modal"
import { ViewInBlockExplorer } from "../../../ViewInBlockExplorer"
import * as Icons from "../../../Icons"
import Button from "../../../Button"
import { colors } from "../../../../constants/colors"
import { withBaseModal } from "../../withBaseModal"

const TokensClaimedComponent = ({ txHash, onClose }) => {
  return (
    <>
      <ModalHeader>Claim</ModalHeader>
      <ModalBody>
        <h3 className="flex row center mb-1">
          <Icons.OK color={colors.success} />
          <span className="ml-1">Success!</span>
        </h3>
        <h4 className="text-grey-70 mb-2">
          View your transaction&nbsp;
          <ViewInBlockExplorer
            type="tx"
            className="text-grey-70"
            id={txHash}
            text="here"
          />
          . Go to the Threshold dapp to complete upgrade.
        </h4>
      </ModalBody>
      <ModalFooter>
        <Button className="btn btn-secondary btn-lg" onClick={onClose}>
          close
        </Button>
      </ModalFooter>
    </>
  )
}

const TokensClaimedWithBaseModal = withBaseModal(TokensClaimedComponent)

export const TokensClaimed = (props) => (
  <TokensClaimedWithBaseModal size="sm" {...props} />
)
