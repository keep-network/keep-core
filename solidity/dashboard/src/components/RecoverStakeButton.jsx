import React, { useCallback } from "react"
import { SubmitButton } from "./Button"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { ContractsLoaded } from "../contracts"
import { useModal } from "../hooks/useModal"
import { withConfirmationModal } from "./ConfirmationModal"
import { connect } from "react-redux"
import { recoverStake } from "../actions/web3"

const RecoverStakeButton = ({ operatorAddress, recoverStake, ...props }) => {
  const { isFromGrant } = props
  const { openConfirmationModal } = useModal()

  const onRecoverStake = useCallback(
    async (awaitingPromise) => {
      const { tokenStakingEscrow } = await ContractsLoaded

      if (isFromGrant) {
        await openConfirmationModal(
          {
            modalOptions: { title: "Are you sure?" },
            title: "You’re about to recover tokens.",
            address: tokenStakingEscrow.options.address,
            btnText: "recover",
            confirmationText: "RECOVER",
          },
          withConfirmationModal(ConfirmRecoveringModal)
        )
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

const ConfirmRecoveringModal = ({ address }) => {
  return (
    <>
      <span>Recovering will deposit delegated tokens in the</span>
      &nbsp;
      <span>
        <ViewAddressInBlockExplorer
          address={address}
          text="TokenStakingEscrow contract."
        />
      </span>
      <p>You can withdraw them via Release tokens.</p>
    </>
  )
}
