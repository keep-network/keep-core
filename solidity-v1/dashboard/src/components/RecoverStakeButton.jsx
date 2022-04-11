import React from "react"
import Button from "./Button"
import { ContractsLoaded } from "../contracts"
import { useModal } from "../hooks/useModal"
import { MODAL_TYPES } from "../constants/constants"
import { useWeb3Address } from "./WithWeb3Context"
import * as Icons from "./Icons"

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
      try {
        await openConfirmationModal(MODAL_TYPES.ConfirmRecovering, {
          tokenStakingEscrowAddress: destinationAddress,
        })
      } catch (err) {
        return
      }
    }

    openModal(MODAL_TYPES.ClaimStakingTokens, {
      amount,
      operator: operatorAddress,
      destinationAddress,
    })
  }

  return (
    <Button className={props.btnClassName} onClick={onRecoverStake}>
      <span className={"flex row center"}>
        <Icons.Refresh color={"#000000"} width={12} height={12} />
        &nbsp;{props.btnText}
      </span>
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
