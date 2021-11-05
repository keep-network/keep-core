import React from "react"
import AddressShortcut from "../../../../AddressShortcut"

const DelegationAddress = ({ address, label }) => {
  return (
    <div className="flex row center text-grey-70">
      <h5 className="flex-1">{label}</h5>
      <AddressShortcut address={address} classNames="h6 text-grey-70" />
    </div>
  )
}

export const DelegationDetails = ({
  authorizerAddress,
  beneficiary,
  operatorAddress,
}) => (
  <>
    <DelegationAddress address={authorizerAddress} label="authorizer" />
    <DelegationAddress address={operatorAddress} label="operator" />
    <DelegationAddress address={beneficiary} label="beneficiary" />
  </>
)
