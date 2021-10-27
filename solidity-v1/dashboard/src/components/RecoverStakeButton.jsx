import React, { useCallback } from "react"
import { SubmitButton } from "./Button"
import { ContractsLoaded } from "../contracts"
import { useModal } from "../hooks/useModal"
import { connect } from "react-redux"
import { recoverStake } from "../actions/web3"
import { MODAL_TYPES } from "../constants/constants"

const RecoverStakeButton = ({ operatorAddress, recoverStake, ...props }) => {
  const { isFromGrant } = props
  const { openConfirmationModal } = useModal()

  const onRecoverStake = useCallback(
    async (awaitingPromise) => {
      const { tokenStakingEscrow } = await ContractsLoaded

      if (isFromGrant) {
        await openConfirmationModal(MODAL_TYPES.ConfirmRecovering, {
          tokenStakingEscrowAddress: tokenStakingEscrow.options.address,
        })
      }

      recoverStake(operatorAddress, awaitingPromise)
    },
    [operatorAddress, recoverStake, isFromGrant, openConfirmationModal]
  )

  return (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={onRecoverStake}
      pendingMessageTitle="Recover stake transaction is pending..."
      successCallback={props.successCallback}
    >
      {props.btnText}
    </SubmitButton>
  )
}

RecoverStakeButton.defaultProps = {
  btnClassName: "btn btn-sm btn-secondary",
  btnText: "recover",
  successCallback: () => {},
  isFromGrant: false,
}

const mapDispatchToProps = {
  recoverStake,
}

const ConnectedWithRedux = connect(null, mapDispatchToProps)(RecoverStakeButton)

export default React.memo(ConnectedWithRedux)
