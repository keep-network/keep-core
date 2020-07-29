import React, { useContext, useCallback, useMemo } from "react"
import { SubmitButton } from "./Button"
import { Web3Context } from "./WithWeb3Context"
import { useShowMessage, messageType } from "./Message"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { contracts } from "../contracts"

const RecoverStakeButton = ({ operatorAddress, ...props }) => {
  const web3Context = useContext(Web3Context)
  const { yourAddress, grantContract, stakingContract } = web3Context
  const showMessage = useShowMessage()
  const { isFromGrant, isManagedGrant, managedGrantContractInstance } = props
  const contract = useMemo(() => {
    if (isManagedGrant) {
      return managedGrantContractInstance
    } else if (isFromGrant) {
      return grantContract
    } else {
      return stakingContract
    }
  }, [
    grantContract,
    isFromGrant,
    isManagedGrant,
    managedGrantContractInstance,
    stakingContract,
  ])

  const recoverStake = useCallback(
    async (
      onTransactionHashCallback,
      openMessageInfo,
      setFetching,
      openConfirmationModal
    ) => {
      try {
        if (isFromGrant) {
          await openConfirmationModal({
            title: "You’re about to recover tokens.",
            subtitle: (
              <>
                <span>Recovering will deposit delegated tokens in the</span>
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
            btnText: "recover",
            confirmationText: "RECOVER",
          })
        }
        await contract.methods
          .recoverStake(operatorAddress)
          .send({ from: yourAddress })
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
    [operatorAddress, yourAddress, contract.methods, showMessage, isFromGrant]
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
