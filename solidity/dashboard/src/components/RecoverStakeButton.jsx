import React, { useCallback } from "react"
import { SubmitButton } from "./Button"
import { useShowMessage, messageType } from "./Message"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { ContractsLoaded } from "../contracts"
import { useModal } from "../hooks/useModal"
import { withConfirmationModal } from "./ConfirmationModal"

const RecoverStakeButton = ({ operatorAddress, ...props }) => {
  const showMessage = useShowMessage()
  const { isFromGrant } = props
  const { openConfirmationModal } = useModal()

  const recoverStake = useCallback(
    async (onTransactionHashCallback) => {
      const { tokenStakingEscrow, stakingContract } = await ContractsLoaded
      try {
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
        await stakingContract.methods
          .recoverStake(operatorAddress)
          .send()
          .on("transactionHash", onTransactionHashCallback)
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Recover stake transaction successfully completed",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Recover stake action has failed ",
          content: error.message,
        })
        throw error
      }
    },
    [operatorAddress, showMessage, isFromGrant, openConfirmationModal]
  )

  return (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={recoverStake}
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

export default React.memo(RecoverStakeButton)

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
