import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"
import { useShowMessage, messageType } from "./Message"
import { SubmitButton } from "./Button"

const UndelegateStakeButton = (props) => {
  const web3Context = useContext(Web3Context)
  const { yourAddress, grantContract, stakingContract } = web3Context
  const showMessage = useShowMessage()

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
      showMessage({
        type: messageType.ERROR,
        title: "Undelegate action has been failed ",
        content: error.message,
      })
      throw error
    }
  }

  return (
    <SubmitButton
      className={props.btnClassName}
      onSubmitAction={undelegate}
      pendingMessageTitle="Undelegate transaction is pending..."
      successCallback={props.successCallback}
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
}

export default UndelegateStakeButton
