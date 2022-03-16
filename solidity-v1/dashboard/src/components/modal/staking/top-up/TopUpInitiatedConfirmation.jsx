import React from "react"
import { DelegationDetails } from "./components"
import { ModalHeader, ModalBody, ModalFooter } from "../../Modal"
import * as Icons from "../../../Icons"
import Button from "../../../Button"
import TokenAmount from "../../../TokenAmount"
import { withBaseModal } from "../../withBaseModal"
import { KEEP } from "../../../../utils/token.utils"
import { add } from "../../../../utils/arithmetics.utils"

export const TopUpInitiatedConfirmation = withBaseModal(
  ({
    addedAmount,
    currentAmount,
    authorizerAddress,
    beneficiary,
    operatorAddress,
    onClose,
  }) => {
    const newAmount = add(currentAmount, addedAmount)

    return (
      <>
        <ModalHeader>Add KEEP</ModalHeader>
        <ModalBody>
          <section className="text-center">
            <h3>Almost there...</h3>
            <p className="text-big text-grey-70">
              Return to the token dashboard <strong>in 12 hours</strong>&nbsp;to
              finalize this transaction!
            </p>
          </section>
          <TokenAmount
            amount={addedAmount}
            withIcon
            icon={Icons.Plus}
            iconProps={{ width: 24, height: 24, className: "plus-icon" }}
          />
          <h4 className="text-grey-70 mb-1">
            New delegation balance: {KEEP.displayAmountWithSymbol(newAmount)}
          </h4>
          <DelegationDetails
            authorizerAddress={authorizerAddress}
            beneficiary={beneficiary}
            operatorAddress={operatorAddress}
          />
        </ModalBody>
        <ModalFooter>
          <Button className="btn btn-lg btn-secondary" onClick={onClose}>
            close
          </Button>
        </ModalFooter>
      </>
    )
  }
)
