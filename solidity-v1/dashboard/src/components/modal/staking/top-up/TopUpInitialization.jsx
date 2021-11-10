import React from "react"
import { DelegationDetails } from "./components"
import { ModalHeader, ModalBody } from "../../Modal"
import { withBaseModal } from "../../withBaseModal"
import AmountForm from "../../../AmountForm"
import Tag from "../../../Tag"
import * as Icons from "../../../Icons"
import TokenAmount from "../../../TokenAmount"

export const TopUpInitialization = withBaseModal(
  ({
    authorizerAddress,
    beneficiary,
    operatorAddress,
    currentAmount,
    onConfirm,
    onClose,
    availableAmount,
  }) => {
    return (
      <>
        <ModalHeader>Add KEEP</ModalHeader>
        <ModalBody>
          <h3 className="mb-1">
            Enter an amount of KEEP to add to this existing delegation.
          </h3>
          <h4 className="text-grey-70 mb-1">Delegation Details</h4>
          <DelegationDetails
            authorizerAddress={authorizerAddress}
            beneficiary={beneficiary}
            operatorAddress={operatorAddress}
          />
          <div className="flex row center mt-2 mb-2">
            <div className="flex-1">
              <Tag text="Current" IconComponent={Icons.KeepToken} />
            </div>
            <TokenAmount amount={currentAmount} />
          </div>
          <AmountForm
            onCancel={onClose}
            submitBtnText="add keep"
            availableAmount={availableAmount}
            currentAmount={currentAmount}
            onBtnClick={onConfirm}
          />
        </ModalBody>
      </>
    )
  }
)
