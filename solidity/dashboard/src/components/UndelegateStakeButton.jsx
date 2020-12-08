import React from "react"
import { SubmitButton } from "./Button"
import { useModal } from "../hooks/useModal"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { ContractsLoaded } from "../contracts"
import { withConfirmationModal } from "./ConfirmationModal"
import { cancelStake, undelegateStake } from "../actions/web3"
import { connect } from "react-redux"

const confirmationModalOptions = {
  modalOptions: { title: "Are you sure?" },
  title: "You’re about to undelegate.",
  subtitle: `Undelegating will return all of your tokens to their owner. There is an undelegation period of 2 months until the tokens will be completely undelegated.`,
  btnText: "undelegate",
  confirmationText: "UNDELEGATE",
}

const confirmCancelModalOptions = {
  modalOptions: { title: "Are you sure?" },
  title: "You’re about to cancel tokens.",
  btnText: "cancel",
  confirmationText: "CANCEL",
}

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
        {
          ...confirmCancelModalOptions,
          tokenStakingEscrowAddress: tokenStakingEscrow.options.address,
        },
        withConfirmationModal(ConfirmCancelingFromGrant)
      )
    } else if (!isInInitializationPeriod) {
      await openConfirmationModal(confirmationModalOptions)
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

const ConfirmCancelingFromGrant = ({ tokenStakingEscrowAddress }) => {
  return (
    <>
      <span>Canceling will deposit delegated tokens in the</span>
      &nbsp;
      <span>
        <ViewAddressInBlockExplorer
          address={tokenStakingEscrowAddress}
          text="TokenStakingEscrow contract."
        />
      </span>
      <p>You can withdraw them via Release tokens.</p>
    </>
  )
}
