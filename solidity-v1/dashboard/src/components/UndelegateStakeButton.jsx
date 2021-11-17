import React from "react"
import { connect } from "react-redux"
import Button, { SubmitButton } from "./Button"
import { useModal } from "../hooks/useModal"
import { ContractsLoaded } from "../contracts"
import { MODAL_TYPES } from "../constants/constants"
import { cancelStake, undelegateStake } from "../actions/web3"

const UndelegateStakeButton = (props) => {
  const { openConfirmationModal, openModal } = useModal()

  const onCancelStake = async (awaitingPromise) => {
    const { isFromGrant, operator, cancelStake } = props
    if (isFromGrant) {
      const { tokenStakingEscrow } = await ContractsLoaded
      await openConfirmationModal(
        MODAL_TYPES.ConfirmCancelDelegationFromGrant,
        {
          tokenStakingEscrowAddress: tokenStakingEscrow.options.address,
        }
      )
    }
    cancelStake(operator, awaitingPromise)
  }

  const undelegate = () => {
    const { operator, undelegationPeriod, amount, authorizer, beneficiary } =
      props

    openModal(MODAL_TYPES.UndelegateStake, {
      undelegationPeriod,
      amount,
      authorizer,
      operator,
      beneficiary,
    })
  }

  return props.isInInitializationPeriod ? (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={onCancelStake}
      successCallback={props.successCallback}
      disabled={props.disabled}
    >
      cancel
    </SubmitButton>
  ) : (
    <Button
      className={props.btnClassName}
      onClick={undelegate}
      disabled={props.disabled}
    >
      {props.btnText}
    </Button>
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
