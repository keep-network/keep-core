import React, { useEffect } from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"
import { KeepLoadingIndicator } from "../Loadable"
import { colors } from "../../constants/colors"

const styles = {
  title: {
    marginBottom: "1rem",
  },
  tileTitle: {
    textAlign: "justify",
    marginBottom: "2rem",
  },
  addressWrapper: {
    marginBottom: "0.5rem",
  },
  delegationsList: { width: "100%" },
  delegationItem: {
    borderRadius: "10px",
  },
  selectedDelegationItem: {
    borderRadius: "10px",
    border: `1px solid ${colors.secondary}`,
  },
}

const CopyStakeStep1 = ({
  fetchDelegations,
  isFetching,
  delegations,
  decrementStep,
  incrementStep,
  onSelectDelegation,
  selectedDelegation,
}) => {
  useEffect(() => {
    fetchDelegations()
  }, [fetchDelegations])

  return (
    <>
      <h2 style={styles.title}>Stake balances to be upgraded.</h2>
      <h3 className="text-grey-70" style={styles.tileTitle}>
        Choose the stake delegation to be upgraded.
      </h3>
      {isFetching ? (
        <div className="flex flex-1 center">
          <KeepLoadingIndicator width={300} height={300} />
        </div>
      ) : (
        <ul style={styles.delegationsList}>
          {delegations.map((delegation) => (
            <DelegationItem
              key={delegation.operatorAddress}
              delegation={delegation}
              onSelect={onSelectDelegation}
              selectedOperatorAddress={
                selectedDelegation ? selectedDelegation.operatorAddress : null
              }
            />
          ))}
        </ul>
      )}
      <div className="flex row space-between self-end">
        <Button
          className="btn btn-transparent btn-lg mr-2"
          onClick={decrementStep}
        >
          back
        </Button>
        <Button
          disabled={!selectedDelegation}
          className="btn btn-primary btn-lg"
          onClick={incrementStep}
        >
          continue
        </Button>
      </div>
    </>
  )
}

const DelegationItem = ({ delegation, onSelect, selectedOperatorAddress }) => {
  const { amount, authorizerAddress, operatorAddress, beneficiary } = delegation
  const isSelected = operatorAddress === selectedOperatorAddress

  return (
    <li
      className="tile"
      style={isSelected ? styles.selectedDelegationItem : styles.delegationItem}
    >
      <input
        type="radio"
        name="option"
        value={operatorAddress}
        id={`option-${operatorAddress}`}
        checked={isSelected}
        onChange={() => onSelect(delegation)}
      />
      <label htmlFor={`option-${operatorAddress}`} style={{ width: "100%" }}>
        <div className="flex row">
          <TokenAmount
            amount={amount}
            currencySymbol="KEEP"
            wrapperClassName="self-start"
            amountClassName="text-primary h3"
          />
          <div style={{ marginLeft: "auto" }}>
            <Address address={authorizerAddress} label="authorizer" />
            <Address address={operatorAddress} label="operator" />
            <Address address={beneficiary} label="beneficiary" />
          </div>
        </div>
      </label>
    </li>
  )
}

const Address = ({ label, address }) => (
  <div className="flex row center" style={styles.addressWrapper}>
    <h5 className="text-grey-70 flex-1">{label}</h5>
    <div className="text-big text-grey-50">{address}</div>
  </div>
)

export default CopyStakeStep1
