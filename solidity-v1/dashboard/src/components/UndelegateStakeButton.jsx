import React from "react"
import { SubmitButton } from "./Button"
import { useModal } from "../hooks/useModal"
import { ContractsLoaded } from "../contracts"
import { cancelStake, undelegateStake } from "../actions/web3"
import { connect } from "react-redux"
import { MODAL_TYPES } from "../constants/constants"

const UndelegateStakeButton = (props) => {
  const { openConfirmationModal } = useModal()

  const undelegate = async (awaitingPromise) => {
    const {
      operator,
      isInInitializationPeriod,
      isFromGrant,
      cancelStake,
      undelegateStake,
    } = props

    if (isInInitializationPeriod && isFromGrant) {
      const { tokenStakingEscrow } = await ContractsLoaded
      await openConfirmationModal(
        MODAL_TYPES.ConfirmCancelDelegationFromGrant,
        {
          tokenStakingEscrowAddress: tokenStakingEscrow.options.address,
        }
      )
    } else if (!isInInitializationPeriod) {
      await openConfirmationModal(MODAL_TYPES.ConfirmUndelegation)
    }

    if (isInInitializationPeriod) {
      cancelStake(operator, awaitingPromise)
    } else {
      undelegateStake(operator, awaitingPromise)
    }
  }

  return (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={undelegate}
      pendingMessageTitle="Undelegate transaction is pending..."
      successCallback={props.successCallback}
      disabled={props.disabled}
    >
      {props.isInInitializationPeriod ? "cancel" : props.btnText}
    </SubmitButton>
  )
}

UndelegateStakeButton.defaultProps = {
  btnClassName: "btn btn-primary btn-sm",
  btnText: "undelegate",
  isInInitializationPeriod: false,
  successCallback: () => {},
  isFromGrant: false,
  disabled: false,
}

const mapDispatchToProps = {
  cancelStake,
  undelegateStake,
}

export default connect(null, mapDispatchToProps)(UndelegateStakeButton)
