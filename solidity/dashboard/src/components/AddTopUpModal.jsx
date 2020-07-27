import React from "react"
import AmountForm from "./AmountForm"
import AddressShortcut from "./AddressShortcut"

const AddTopUpModal = ({
  authorizerAddress,
  beneficiary,
  operatorAddress,
  onCancel,
  submitBtnText,
  availableAmount,
}) => {
  return (
    <>
      <h3>Enter an amount of KEEP to add to this existing delegation.</h3>
      <p className="text-big text-grey-70">Delegation Details</p>
      <DelegationAddress address={authorizerAddress} label="authorizer" />
      <DelegationAddress address={operatorAddress} label="operator" />
      <DelegationAddress address={beneficiary} label="beneficiary" />
      <AmountForm
        onCancel={onCancel}
        submitBtnText={submitBtnText}
        availableAmount={availableAmount}
      />
    </>
  )
}

const DelegationAddress = ({ address, label }) => {
  return (
    <div className="flex row center">
      <div className="text-label text-small text-grey-70 flex-1">{label}</div>
      <div className="flex-2 self-start">
        <AddressShortcut
          address={address}
          classNames="text-label text-grey-50"
        />
      </div>
    </div>
  )
}

export default AddTopUpModal
