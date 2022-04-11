import React from "react"
import { connect } from "react-redux"
import Button, { SubmitButton } from "./Button"
import { useModal } from "../hooks/useModal"
import { ContractsLoaded } from "../contracts"
import { LINK, MODAL_TYPES } from "../constants/constants"
import { cancelStake, undelegateStake } from "../actions/web3"
import ReactTooltip from "react-tooltip"

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

  const renderUndelegateStakeButton = (
    delegationOperator,
    btnClassName,
    disabled
  ) => {
    if (disabled) {
      return (
        <>
          <span
            data-tip
            data-for={`undelegate-button-for-operator-${delegationOperator}`}
          >
            <Button
              className={`undelegate-stake-button ${btnClassName}`}
              onClick={undelegate}
              disabled={true}
            >
              {props.btnText}
            </Button>
          </span>
          <ReactTooltip
            id={`undelegate-button-for-operator-${delegationOperator}`}
            place="top"
            type="dark"
            effect={"solid"}
            className={"react-tooltip-base"}
            clickable={true}
            delayHide={200}
          >
            <span>
              This stake is staked on Threshold. You first need to undelegate
              this stake from the{" "}
              <a
                href={LINK.thresholdDapp}
                rel="noopener noreferrer"
                target="_blank"
              >
                Threshold dashboard here
              </a>
            </span>
          </ReactTooltip>
        </>
      )
    }

    return (
      <Button
        className={props.btnClassName}
        onClick={undelegate}
        disabled={props.disabled}
      >
        {props.btnText}
      </Button>
    )
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
    renderUndelegateStakeButton(
      props.operator,
      props.btnClassName,
      props.disabled
    )
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
