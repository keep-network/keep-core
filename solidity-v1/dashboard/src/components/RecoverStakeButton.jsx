import React from "react"
import Button from "./Button"
import { ContractsLoaded } from "../contracts"
import { useModal } from "../hooks/useModal"
import { MODAL_TYPES } from "../constants/constants"
import { useWeb3Address } from "./WithWeb3Context"

const RecoverStakeButton = ({
  operatorAddress,
  recoverStake,
  amount,
  isFromGrant,
  ...props
}) => {
  const { openConfirmationModal, openModal } = useModal()
  const address = useWeb3Address()

  const onRecoverStake = async () => {
    const { tokenStakingEscrow } = await ContractsLoaded
    let destinationAddress = address
    if (isFromGrant) {
      destinationAddress = tokenStakingEscrow.options.address
      await openConfirmationModal(MODAL_TYPES.ConfirmRecovering, {
        tokenStakingEscrowAddress: destinationAddress,
      })
    }

    openModal(MODAL_TYPES.ClaimStakingTokens, {
      amount,
      operator: operatorAddress,
      destinationAddress,
    })
  }

  return (
    <Button className={props.btnClassName} onClick={onRecoverStake}>
      {props.btnText}
    </Button>
  )
}

RecoverStakeButton.defaultProps = {
  btnClassName: "btn btn-sm btn-secondary",
  btnText: "recover",
  successCallback: () => {},
  isFromGrant: false,
}

export default React.memo(RecoverStakeButton)
