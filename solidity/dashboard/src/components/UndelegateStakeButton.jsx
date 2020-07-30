import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"
import { useShowMessage, messageType } from "./Message"
import { SubmitButton } from "./Button"
import { useModal } from "../hooks/useModal"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { contracts } from "../contracts"

const confirmationModalOptions = {
  title: "You’re about to undelegate.",
  subtitle:
    "Undelegating will return all of your tokens to their owner. There is an undelegation period of 1 week until the tokens will be completely undelegated.",
  btnText: "undelegate",
  confirmationText: "UNDELEGATE",
}

const UndelegateStakeButton = (props) => {
  const web3Context = useContext(Web3Context)
  const { yourAddress, grantContract, stakingContract } = web3Context
  const showMessage = useShowMessage()
  const { openConfirmationModal } = useModal()

  const undelegate = async (onTransactionHashCallback) => {
    const {
      operator,
      isInInitializationPeriod,
      isFromGrant,
      isManagedGrant,
      managedGrantContractInstance,
    } = props
    let contract
    if (isManagedGrant) {
      contract = managedGrantContractInstance
    } else if (isFromGrant) {
      contract = grantContract
    } else {
      contract = stakingContract
    }
    try {
      if (isInInitializationPeriod && isFromGrant) {
        await openConfirmationModal({
          title: "You’re about to cancel tokens.",
          subtitle: (
            <>
              <span>Canceling will deposit delegated tokens in the</span>
              &nbsp;
              <span>
                <ViewAddressInBlockExplorer
                  address={contracts.tokenStakingEscrow.options.address}
                  text="TokenStakingEscrow contract."
                />
              </span>
              <p>You can withdraw them via Release tokens.</p>
            </>
          ),
          btnText: "cancel",
          confirmationText: "CANCEL",
        })
      } else {
        await openConfirmationModal(confirmationModalOptions)
      }

      await contract.methods[
        isInInitializationPeriod ? "cancelStake" : "undelegate"
      ](operator)
        .send({ from: yourAddress })
        .on("transactionHash", onTransactionHashCallback)
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Undelegate transaction successfully completed",
      })
    } catch (error) {
      if (!error.type || error.type !== "canceled") {
        showMessage({
          type: messageType.ERROR,
          title: "Undelegate action has failed ",
          content: error.message,
        })
      }
      throw error
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

export default UndelegateStakeButton
