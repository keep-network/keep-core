import React from "react"
import { useDispatch } from "react-redux"
import { ModalBody, ModalHeader, ModalFooter } from "../../Modal"
import TokenAmount from "../../../TokenAmount"
import AddressShortcut from "../../../AddressShortcut"
import Button, { SubmitButton } from "../../../Button"
import { KEEP } from "../../../../utils/token.utils"
import { recoverStake } from "../../../../actions/web3"
import { withBaseModal } from "../../withBaseModal"

const ClaimTokensComponent = ({
  amount,
  operator,
  destinationAddress,
  onClose,
  //   TODO: Render an additional info indicating the tokens go to the
  //   `TokenStakingEscrow` contract
  isFromGrant = false,
}) => {
  const dispatch = useDispatch()
  return (
    <>
      <ModalHeader>Claim</ModalHeader>
      <ModalBody>
        <h3 className="mb-1">Claim tokens</h3>
        <TokenAmount amount={amount} token={KEEP} withIcon />
        <div className="flex row center mt-1">
          <div className="text-grey-50 mr-a">Wallet</div>
          <div className="text-grey-70">
            <AddressShortcut address={destinationAddress} />
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <SubmitButton
          className="btn btn-primary btn-lg mr-2"
          onSubmitAction={(awaitingPromise) =>
            dispatch(recoverStake(operator, awaitingPromise))
          }
        >
          claim
        </SubmitButton>
        <Button className="btn btn-unstyled" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
}

const ClaimTokensWithBaseModal = withBaseModal(ClaimTokensComponent)

export const ClaimTokens = (props) => (
  <ClaimTokensWithBaseModal size="sm" {...props} />
)
